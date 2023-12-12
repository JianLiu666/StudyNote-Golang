package trading

import (
	"container/heap"
	"context"
	"encoding/json"
	"fmt"
	"interview20231208/model"
	"interview20231208/pkg/e"
	"interview20231208/pkg/trading"
	"time"

	"github.com/rs/xid"
)

type tradingPool struct {
	orderChan  chan *model.Order
	buyerHeap  *CustomHeap
	sellerHeap *CustomHeap
}

func newTradingPool() *tradingPool {
	return NewTradingPool().(*tradingPool)
}

func NewTradingPool() trading.TradingPool {
	return &tradingPool{
		orderChan: make(chan *model.Order, 1024),
		buyerHeap: NewCustomHeap(func(i, j *model.Order) bool {
			if i.Price == j.Price {
				return i.Timestamp < j.Timestamp
			}
			return i.Price > j.Price
		}),
		sellerHeap: NewCustomHeap(func(i, j *model.Order) bool {
			if i.Price == j.Price {
				return i.Timestamp < j.Timestamp
			}
			return i.Price < j.Price
		}),
	}
}

func (t *tradingPool) Enable(ctx context.Context) {
	go t.schedule(ctx)
}

func (t *tradingPool) AddOrder(order *model.Order) {
	t.orderChan <- order
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
	var result []*model.TransactionLog

	if order.OrderType == e.ORDER_LIMIT {
		switch order.DurationType {
		case e.DURATION_ROD:
			result = t.processLimitROD(order)
		case e.DURATION_IOC:
			t.processLimitIOC(order)
		case e.DURATION_FOK:
			result = t.processLimitFOK(order)
		}

	} else if order.OrderType == e.ORDER_MARKET {
		switch order.DurationType {
		case e.DURATION_ROD:
			t.processMarketROD(order)
		case e.DURATION_IOC:
			t.processMarketIOC(order)
		case e.DURATION_FOK:
			result = t.processMarketFOK(order)
		}
	}

	// TODO
	b, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
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
			UUID:          xid.New().String(),
			BuyerOrderID:  t.buyerHeap.Peek().UUID,
			SellerOrderID: t.sellerHeap.Peek().UUID,
			Price:         t.buyerHeap.Peek().Price,
			Quantity:      min(t.buyerHeap.Peek().Quantity, t.sellerHeap.Peek().Quantity),
			Timestamp:     int(time.Now().Unix()),
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
					UUID:          xid.New().String(),
					BuyerOrderID:  order.UUID,
					SellerOrderID: t.sellerHeap.Peek().UUID,
					Price:         t.sellerHeap.Peek().Price,
					Quantity:      order.Quantity,
					Timestamp:     int(time.Now().Unix()), // FIXME: 這邊可能存在最後一筆交易紀錄的時間早於其他交易紀錄的問題
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
				UUID:          xid.New().String(),
				BuyerOrderID:  order.UUID,
				SellerOrderID: candidates[0].UUID,
				Price:         candidates[0].Price,
				Quantity:      candidates[0].Quantity,
				Timestamp:     int(time.Now().Unix()),
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
					UUID:          xid.New().String(),
					BuyerOrderID:  t.buyerHeap.Peek().UUID,
					SellerOrderID: order.UUID,
					Price:         t.buyerHeap.Peek().Price,
					Quantity:      order.Quantity,
					Timestamp:     int(time.Now().Unix()), // FIXME: 這邊可能存在最後一筆交易紀錄的時間早於其他交易紀錄的問題
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
				UUID:          xid.New().String(),
				BuyerOrderID:  candidates[0].UUID,
				SellerOrderID: order.UUID,
				Price:         candidates[0].Price,
				Quantity:      candidates[0].Quantity,
				Timestamp:     int(time.Now().Unix()),
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
				UUID:          xid.New().String(),
				BuyerOrderID:  order.UUID,
				SellerOrderID: t.sellerHeap.Peek().UUID,
				Price:         t.sellerHeap.Peek().Price,
				Quantity:      min(order.Quantity, t.sellerHeap.Peek().Quantity),
				Timestamp:     int(time.Now().Unix()),
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
				UUID:          xid.New().String(),
				BuyerOrderID:  t.buyerHeap.Peek().UUID,
				SellerOrderID: order.UUID,
				Price:         t.buyerHeap.Peek().Price,
				Quantity:      min(order.Quantity, t.buyerHeap.Peek().Quantity),
				Timestamp:     int(time.Now().Unix()),
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
