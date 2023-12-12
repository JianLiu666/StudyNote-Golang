package trading

import (
	"context"
	"interview20231208/model"
	"interview20231208/pkg/e"
	"interview20231208/pkg/rdb/mysql"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestAddOrder_RODvsROD(t *testing.T) {
	// prepare data
	src_ordersIdx := 0
	src_orders := []*model.Order{
		{
			ID:           1,
			UserID:       1,
			RoleType:     e.ROLE_BUYER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        100,
			Quantity:     100,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           2,
			UserID:       2,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        100,
			Quantity:     30,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
	}

	dst_transactionLogsIdx := 0
	dst_transactionLogs := make([][]*model.TransactionLog, 2)
	dst_transactionLogs[0] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[1] = []*model.TransactionLog{
		{
			ID:            0,
			BuyerOrderID:  1,
			SellerOrderID: 2,
			Price:         100,
			Quantity:      30,
			Timestamp:     time.Now(),
		},
	}

	// init
	mockrdb := mysql.NewMockMysqlClient()
	mockrdb.SetCreateOrderCallback(func(ctx context.Context, order *model.Order) {
		assert.DeepEqual(t, src_orders[src_ordersIdx], order)
		src_ordersIdx++
	})
	mockrdb.SetUpdateOrdersAndCreateTransactionLogs(func(ctx context.Context, orders map[int]*model.Order, logs []*model.TransactionLog) {
		for i, dst := range dst_transactionLogs[dst_transactionLogsIdx] {
			logs[i].ID = 0
			logs[i].Timestamp = time.Now()
			assert.DeepEqual(t, dst, logs[i])
		}
		dst_transactionLogsIdx++
	})

	pool := newTradingPool(mockrdb)

	// testcase flow
	pool.Enable(context.Background())

	for _, order := range src_orders {
		pool.AddOrder(order)
	}

	time.Sleep(1 * time.Second)

	// other validation
	if pool.buyerHeap.Len() != 1 {
		t.Error()
	}
	assert.Equal(t, 70, pool.buyerHeap.Peek().Quantity)

	if pool.sellerHeap.Len() != 0 {
		t.Error()
	}
}
