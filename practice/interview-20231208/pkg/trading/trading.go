package trading

import (
	"interview20231208/model"
	"interview20231208/pkg/e"
)

type tradingPool struct {
	buyerQuantities  map[int]int
	sellerQuantities map[int]int
	buyerHeap        *CustomHeap
	sellerHeap       *CustomHeap
}

func NewTradingPool() *tradingPool {
	return &tradingPool{
		buyerQuantities:  make(map[int]int, 0),
		sellerQuantities: make(map[int]int, 0),
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

func (t *tradingPool) AddOrder(order *model.Order) {
	if order.DurationType != e.DURATION_FOK {
		// TODO: error handler
		return
	}

	if order.OrderType == e.ORDER_MARKET {
		// TODO
		return
	}

	if order.RoleType == e.ROLE_BUYER {
		t.buyerQuantities[order.Price]++
	} else if order.RoleType == e.ROLE_SELLER {
		t.sellerQuantities[order.Price]++
	}

	// TODO: matching
}
