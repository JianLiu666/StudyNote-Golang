package trading

import (
	"container/heap"
	"context"
	"interview20231208/model"
	"interview20231208/pkg/e"
	"interview20231208/pkg/rdb"
	"interview20231208/pkg/trading"
	"time"
)

type tradingPool struct {
	rdb        rdb.RDB
	orderChan  chan *model.Order
	buyerHeap  *CustomHeap
	sellerHeap *CustomHeap
}

func newTradingPool(rdb rdb.RDB) *tradingPool {
	return NewTradingPool(rdb).(*tradingPool)
}

func NewTradingPool(rdb rdb.RDB) trading.TradingPool {
	return &tradingPool{
		rdb:       rdb,
		orderChan: make(chan *model.Order, 1024),
		buyerHeap: NewCustomHeap(func(i, j *model.Order) bool {
			if i.Price == j.Price {
				return i.Timestamp.Unix() < j.Timestamp.Unix()
			}
			return i.Price > j.Price
		}),
		sellerHeap: NewCustomHeap(func(i, j *model.Order) bool {
			if i.Price == j.Price {
				return i.Timestamp.Unix() < j.Timestamp.Unix()
			}
			return i.Price < j.Price
		}),
	}
}

func (t *tradingPool) Enable(ctx context.Context) {
	go t.schedule(ctx)
}

func (t *tradingPool) AddOrder(order *model.Order) e.CODE {
	t.rdb.CreateOrder(context.TODO(), order)
	t.orderChan <- order

	return e.SUCCESS
}

func (t *tradingPool) schedule(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// TODO
			break

		case order := <-t.orderChan:
			t.consume(order)
		}
	}
}

func (t *tradingPool) consume(order *model.Order) {
	var logs []*model.TransactionLog
	orderSet := map[int]*model.Order{}

	if order.OrderType == e.ORDER_LIMIT {
		switch order.DurationType {
		case e.DURATION_ROD:
			logs = t.processLimitROD(order)
		case e.DURATION_IOC:
			t.processLimitIOC(order)
		case e.DURATION_FOK:
			logs = t.processLimitFOK(order)
		}

	} else if order.OrderType == e.ORDER_MARKET {
		switch order.DurationType {
		case e.DURATION_ROD:
			t.processMarketROD(order)
		case e.DURATION_IOC:
			t.processMarketIOC(order)
		case e.DURATION_FOK:
			logs = t.processMarketFOK(order)
		}
	}

	// 如果期限是 ROD 的交易單, 即使沒有撮合成功也不需要更新交易單的狀態
	if len(logs) == 0 && order.DurationType == e.DURATION_ROD {
		return
	}

	// 期限是 IOC/FOK 的交易單, 根據是否有撮合成功來決定如何更新 order 狀態
	if len(logs) == 0 {
		orderSet[order.ID] = &model.Order{
			ID:     order.ID,
			Status: e.STATUS_CANCEL,
		}
		order.Status = e.STATUS_CANCEL

	} else {
		for _, log := range logs {
			orderSet[log.BuyerOrderID] = &model.Order{
				ID:     log.BuyerOrderID,
				Status: e.STATUS_SUCCESS,
			}
			orderSet[log.SellerOrderID] = &model.Order{
				ID:     log.SellerOrderID,
				Status: e.STATUS_SUCCESS,
			}
		}
	}

	t.rdb.UpdateOrdersAndCreateTransactionLogs(context.TODO(), orderSet, logs)
}

func (t *tradingPool) processLimitROD(order *model.Order) []*model.TransactionLog {
	result := []*model.TransactionLog{}

	// 加入新的 order 到對應 PQ
	if order.RoleType == e.ROLE_BUYER {
		heap.Push(t.buyerHeap, order)
	} else if order.RoleType == e.ROLE_SELLER {
		heap.Push(t.sellerHeap, order)
	}

	// 不斷撮合直到買賣方出現價差為止
	for t.buyerHeap.Len() > 0 && t.sellerHeap.Len() > 0 && t.buyerHeap.Peek().Price >= t.sellerHeap.Peek().Price {

		// step.1 產生一筆新的交易紀錄
		result = append(result, &model.TransactionLog{
			BuyerOrderID:  t.buyerHeap.Peek().ID,
			SellerOrderID: t.sellerHeap.Peek().ID,
			Price:         t.buyerHeap.Peek().Price,
			Quantity:      min(t.buyerHeap.Peek().Quantity, t.sellerHeap.Peek().Quantity),
			Timestamp:     time.Now(),
		})

		// step.2 根據買賣方的 quantity 決定如何更新 PQ
		if t.buyerHeap.Peek().Quantity < t.sellerHeap.Peek().Quantity {
			t.sellerHeap.Peek().Quantity -= t.buyerHeap.Peek().Quantity
			heap.Pop(t.buyerHeap)
		} else if t.buyerHeap.Peek().Quantity > t.sellerHeap.Peek().Quantity {
			t.buyerHeap.Peek().Quantity -= t.sellerHeap.Peek().Quantity
			heap.Pop(t.sellerHeap)
		} else {
			heap.Pop(t.buyerHeap)
			heap.Pop(t.sellerHeap)
		}
	}

	return result
}

func (t *tradingPool) processLimitIOC(order *model.Order) {
	// TODO
}

func (t *tradingPool) processLimitFOK(order *model.Order) []*model.TransactionLog {
	result := []*model.TransactionLog{}

	if order.RoleType == e.ROLE_BUYER {
		// 確認賣方數量必須大於 order
		if t.sellerHeap.TotalQuantity() < order.Quantity {
			return result
		}

		candidates := []*model.Order{}

		// 需要考慮 limit order 可能還是會遇到 range 內數量不夠最後要 rollback 的問題
		// 如果 iteration 結束後 order.Quantity != 0 時, 那 PQ[0].Quantity 絕對不會被異動到
		for order.Price >= t.sellerHeap.Peek().Price && order.Quantity > 0 {
			if order.Quantity >= t.sellerHeap.Peek().Quantity {
				order.Quantity -= t.sellerHeap.Peek().Quantity
				candidates = append(candidates, heap.Pop(t.sellerHeap).(*model.Order))
			} else {
				result = append(result, &model.TransactionLog{
					BuyerOrderID:  order.ID,
					SellerOrderID: t.sellerHeap.Peek().ID,
					Price:         t.sellerHeap.Peek().Price,
					Quantity:      order.Quantity,
					Timestamp:     time.Now(), // FIXME: 這邊可能存在最後一筆交易紀錄的時間早於其他交易紀錄的問題
				})
				t.sellerHeap.Peek().Quantity -= order.Quantity
				order.Quantity = 0
			}
		}

		if order.Quantity != 0 {
			// rollback
			for len(candidates) > 0 {
				heap.Push(t.sellerHeap, candidates[0])
				candidates = candidates[1:]
			}
			return result
		}

		for len(candidates) > 0 {
			result = append(result, &model.TransactionLog{
				BuyerOrderID:  order.ID,
				SellerOrderID: candidates[0].ID,
				Price:         candidates[0].Price,
				Quantity:      candidates[0].Quantity,
				Timestamp:     time.Now(),
			})
			candidates = candidates[1:]
		}

	} else if order.RoleType == e.ROLE_SELLER {
		// 確認買方數量必須大於 order
		if t.buyerHeap.TotalQuantity() < order.Quantity {
			return result
		}

		candidates := []*model.Order{}

		// 需要考慮 limit order 可能還是會遇到 range 內數量不夠最後要 rollback 的問題
		// 如果 iteration 結束後 order.Quantity != 0 時, 那 PQ[0].Quantity 絕對不會被異動到
		for t.buyerHeap.Peek().Price >= order.Price && order.Quantity > 0 {
			if order.Quantity >= t.buyerHeap.Peek().Quantity {
				order.Quantity -= t.buyerHeap.Peek().Quantity
				candidates = append(candidates, heap.Pop(t.buyerHeap).(*model.Order))
			} else {
				result = append(result, &model.TransactionLog{
					BuyerOrderID:  t.buyerHeap.Peek().ID,
					SellerOrderID: order.ID,
					Price:         t.buyerHeap.Peek().Price,
					Quantity:      order.Quantity,
					Timestamp:     time.Now(), // FIXME: 這邊可能存在最後一筆交易紀錄的時間早於其他交易紀錄的問題
				})
				t.buyerHeap.Peek().Quantity -= order.Quantity
				order.Quantity = 0
			}
		}

		if order.Quantity != 0 {
			// rollback
			for len(candidates) > 0 {
				heap.Push(t.buyerHeap, candidates[0])
				candidates = candidates[1:]
			}
			return result
		}

		for len(candidates) > 0 {
			result = append(result, &model.TransactionLog{
				BuyerOrderID:  candidates[0].ID,
				SellerOrderID: order.ID,
				Price:         candidates[0].Price,
				Quantity:      candidates[0].Quantity,
				Timestamp:     time.Now(),
			})
			candidates = candidates[1:]
		}
	}

	return result
}

func (t *tradingPool) processMarketROD(order *model.Order) {
	// TODO
}

func (t *tradingPool) processMarketIOC(order *model.Order) {
	// TODO
}

func (t *tradingPool) processMarketFOK(order *model.Order) []*model.TransactionLog {
	result := []*model.TransactionLog{}

	if order.RoleType == e.ROLE_BUYER {
		// 確認賣方數量必須大於 order
		if t.sellerHeap.TotalQuantity() < order.Quantity {
			return result
		}

		// 不斷撮合直到 market order 收滿為止
		for order.Quantity > 0 {
			// step.1 產生一筆新的交易紀錄
			result = append(result, &model.TransactionLog{
				BuyerOrderID:  order.ID,
				SellerOrderID: t.sellerHeap.Peek().ID,
				Price:         t.sellerHeap.Peek().Price,
				Quantity:      min(order.Quantity, t.sellerHeap.Peek().Quantity),
				Timestamp:     time.Now(),
			})

			// step.2 根據 quantity 決定如何更新賣方的 PQ
			if order.Quantity >= t.sellerHeap.Peek().Quantity {
				order.Quantity -= t.sellerHeap.Peek().Quantity
				heap.Pop(t.sellerHeap)
			} else {
				t.sellerHeap.Peek().Quantity -= order.Quantity
				order.Quantity = 0
			}
		}

	} else if order.RoleType == e.ROLE_SELLER {
		// 確認買方數量必須大於 order
		if t.buyerHeap.TotalQuantity() < order.Quantity {
			return result
		}

		for order.Quantity > 0 {
			// step.1 產生一筆新的交易紀錄
			result = append(result, &model.TransactionLog{
				BuyerOrderID:  t.buyerHeap.Peek().ID,
				SellerOrderID: order.ID,
				Price:         t.buyerHeap.Peek().Price,
				Quantity:      min(order.Quantity, t.buyerHeap.Peek().Quantity),
				Timestamp:     time.Now(),
			})

			// step.2 根據 quantity 決定如何更新買方的 PQ
			if order.Quantity >= t.buyerHeap.Peek().Quantity {
				order.Quantity -= t.buyerHeap.Peek().Quantity
				heap.Pop(t.buyerHeap)
			} else {
				t.buyerHeap.Peek().Quantity -= order.Quantity
				order.Quantity = 0
			}
		}
	}

	return result
}
