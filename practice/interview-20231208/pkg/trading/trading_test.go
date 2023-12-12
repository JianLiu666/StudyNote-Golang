package trading

import (
	"context"
	"interview20231208/model"
	"interview20231208/pkg/e"
	"testing"
	"time"

	"github.com/rs/xid"
)

func TestAddOrder_RODvsROD(t *testing.T) {
	pool := NewTradingPool()
	pool.Enable(context.Background())

	pool.AddOrder(&model.Order{
		UUID:         xid.New().String(),
		UserID:       1,
		RoleType:     e.ROLE_BUYER,
		OrderType:    e.ORDER_LIMIT,
		DurationType: e.DURATION_ROD,
		Price:        100,
		Quantity:     100,
		Status:       e.STATUS_PENDING,
		Timestamp:    int(time.Now().Unix()),
	})

	pool.AddOrder(&model.Order{
		UUID:         xid.New().String(),
		UserID:       1,
		RoleType:     e.ROLE_SELLER,
		OrderType:    e.ORDER_LIMIT,
		DurationType: e.DURATION_ROD,
		Price:        100,
		Quantity:     30,
		Status:       e.STATUS_PENDING,
		Timestamp:    int(time.Now().Unix()),
	})

	time.Sleep(1 * time.Second)
}
