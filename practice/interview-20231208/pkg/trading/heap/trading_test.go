// go:build: unittest

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

// TestAddOrder_LimitRODvsLimitROD_case1
// 買方限價單完全成交, 賣方限價單部分成交
// Goal:
//   - 依序寫入買方限價單, 測試 PQ 排序規則(price & timestamp)
//   - 依序寫入賣方限價單, 測試撮合機制是否正確
//   - 測試 transaction logs 數量是否正確
//   - 測試買賣方 PQs 剩餘數量與 quantity 是否正確
func TestAddOrder_LimitRODvsLimitROD_case1(t *testing.T) {
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
			UserID:       1,
			RoleType:     e.ROLE_BUYER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        100,
			Quantity:     20,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           3,
			UserID:       1,
			RoleType:     e.ROLE_BUYER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        110,
			Quantity:     30,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           4,
			UserID:       2,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        100,
			Quantity:     100,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           5,
			UserID:       2,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        100,
			Quantity:     40,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           6,
			UserID:       2,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        120,
			Quantity:     10,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           7,
			UserID:       2,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        90,
			Quantity:     10,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
	}

	dst_transactionLogsIdx := 0
	dst_transactionLogs := make([][]*model.TransactionLog, len(src_orders))
	dst_transactionLogs[0] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[1] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[2] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[3] = []*model.TransactionLog{
		{
			BuyerOrderID:  3,
			SellerOrderID: 4,
			Price:         100,
			Quantity:      30,
			Timestamp:     time.Now(),
		},
		{
			BuyerOrderID:  1,
			SellerOrderID: 4,
			Price:         100,
			Quantity:      70,
			Timestamp:     time.Now(),
		},
	}
	dst_transactionLogs[4] = []*model.TransactionLog{
		{
			BuyerOrderID:  1,
			SellerOrderID: 5,
			Price:         100,
			Quantity:      30,
			Timestamp:     time.Now(),
		},
		{
			BuyerOrderID:  2,
			SellerOrderID: 5,
			Price:         100,
			Quantity:      10,
			Timestamp:     time.Now(),
		},
	}
	dst_transactionLogs[5] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[6] = []*model.TransactionLog{
		{
			BuyerOrderID:  1,
			SellerOrderID: 5,
			Price:         100,
			Quantity:      30,
			Timestamp:     time.Now(),
		},
		{
			BuyerOrderID:  2,
			SellerOrderID: 7,
			Price:         90,
			Quantity:      10,
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

	time.Sleep(100 * time.Millisecond)

	// other validation
	assert.Equal(t, 0, pool.buyerHeap.Len())

	assert.Equal(t, 1, pool.sellerHeap.Len())
	assert.Equal(t, 6, pool.sellerHeap.Peek().ID)
	assert.Equal(t, 10, pool.sellerHeap.Peek().Quantity)
}

// TestAddOrder_LimitRODvsLimitFOK_case1
// 預先建立買方限價單, 測試賣方市價單撮合邏輯
// Goal:
//   - 依序寫入買方限價單, 測試 PQ 排序規則(price & timestamp)
//   - 依序寫入賣方市價單, 測試撮合機制是否正確
//   - 測試 transaction logs 數量是否正確
//   - 測試買賣方 PQs 剩餘數量與 quantity 是否正確
func TestAddOrder_LimitRODvsLimitFOK_case1(t *testing.T) {
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
			UserID:       1,
			RoleType:     e.ROLE_BUYER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        100,
			Quantity:     20,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           3,
			UserID:       1,
			RoleType:     e.ROLE_BUYER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        110,
			Quantity:     30,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           4,
			UserID:       2,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_FOK,
			Price:        120,
			Quantity:     150,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           5,
			UserID:       2,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_FOK,
			Price:        110,
			Quantity:     40,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           6,
			UserID:       2,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_FOK,
			Price:        110,
			Quantity:     20,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           7,
			UserID:       2,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_FOK,
			Price:        100,
			Quantity:     90,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
	}

	dst_transactionLogsIdx := 0
	dst_transactionLogs := make([][]*model.TransactionLog, len(src_orders))
	dst_transactionLogs[0] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[1] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[2] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[3] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[4] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[5] = []*model.TransactionLog{
		{
			BuyerOrderID:  3,
			SellerOrderID: 6,
			Price:         110,
			Quantity:      20,
			Timestamp:     time.Now(),
		},
	}
	dst_transactionLogs[6] = []*model.TransactionLog{
		{
			BuyerOrderID:  3,
			SellerOrderID: 7,
			Price:         110,
			Quantity:      10,
			Timestamp:     time.Now(),
		},
		{
			BuyerOrderID:  1,
			SellerOrderID: 7,
			Price:         100,
			Quantity:      80,
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

	time.Sleep(100 * time.Millisecond)

	// other validation
	assert.Equal(t, 2, pool.buyerHeap.Len())
	assert.Equal(t, 1, pool.buyerHeap.Peek().ID)
	assert.Equal(t, 20, pool.buyerHeap.Peek().Quantity)

	assert.Equal(t, 0, pool.sellerHeap.Len())
}

// TestAddOrder_LimitFOKvsLimitROD_case1
// 預先建立賣方限價單, 測試買方市價單撮合邏輯
// Goal:
//   - 依序寫入賣方限價單, 測試 PQ 排序規則(price & timestamp)
//   - 依序寫入買方市價單, 測試撮合機制是否正確
//   - 測試 transaction logs 數量是否正確
//   - 測試買賣方 PQs 剩餘數量與 quantity 是否正確
func TestAddOrder_LimitFOKvsLimitROD_case1(t *testing.T) {
	// prepare data
	src_ordersIdx := 0
	src_orders := []*model.Order{
		{
			ID:           1,
			UserID:       1,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        110,
			Quantity:     100,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           2,
			UserID:       1,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        110,
			Quantity:     20,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           3,
			UserID:       1,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        100,
			Quantity:     30,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           4,
			UserID:       2,
			RoleType:     e.ROLE_BUYER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_FOK,
			Price:        100,
			Quantity:     160,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           5,
			UserID:       2,
			RoleType:     e.ROLE_BUYER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_FOK,
			Price:        100,
			Quantity:     40,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           6,
			UserID:       2,
			RoleType:     e.ROLE_BUYER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_FOK,
			Price:        100,
			Quantity:     20,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           7,
			UserID:       2,
			RoleType:     e.ROLE_BUYER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_FOK,
			Price:        110,
			Quantity:     90,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
	}

	dst_transactionLogsIdx := 0
	dst_transactionLogs := make([][]*model.TransactionLog, len(src_orders))
	dst_transactionLogs[0] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[1] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[2] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[3] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[4] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[5] = []*model.TransactionLog{
		{
			BuyerOrderID:  6,
			SellerOrderID: 3,
			Price:         100,
			Quantity:      20,
			Timestamp:     time.Now(),
		},
	}
	dst_transactionLogs[6] = []*model.TransactionLog{
		{
			BuyerOrderID:  7,
			SellerOrderID: 3,
			Price:         100,
			Quantity:      10,
			Timestamp:     time.Now(),
		},
		{
			BuyerOrderID:  7,
			SellerOrderID: 1,
			Price:         110,
			Quantity:      80,
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

	time.Sleep(100 * time.Millisecond)

	// other validation
	assert.Equal(t, 0, pool.buyerHeap.Len())

	assert.Equal(t, 2, pool.sellerHeap.Len())
	assert.Equal(t, 1, pool.sellerHeap.Peek().ID)
	assert.Equal(t, 20, pool.sellerHeap.Peek().Quantity)
}

// TestAddOrder_LimitRODvsMarketFOK_case1
// 預先建立買方限價單, 測試賣方市價單撮合邏輯
// Goal:
//   - 依序寫入買方限價單, 測試 PQ 排序規則(price & timestamp)
//   - 依序寫入賣方市價單, 測試撮合機制是否正確
//   - 測試 transaction logs 數量是否正確
//   - 測試買賣方 PQs 剩餘數量與 quantity 是否正確
func TestAddOrder_LimitRODvsMarketFOK_case1(t *testing.T) {
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
			UserID:       1,
			RoleType:     e.ROLE_BUYER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        100,
			Quantity:     20,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           3,
			UserID:       1,
			RoleType:     e.ROLE_BUYER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        110,
			Quantity:     30,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           4,
			UserID:       2,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_MARKET,
			DurationType: e.DURATION_FOK,
			Price:        0,
			Quantity:     160,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           5,
			UserID:       2,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_MARKET,
			DurationType: e.DURATION_FOK,
			Price:        0,
			Quantity:     50,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           6,
			UserID:       2,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_MARKET,
			DurationType: e.DURATION_FOK,
			Price:        0,
			Quantity:     100,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
	}

	dst_transactionLogsIdx := 0
	dst_transactionLogs := make([][]*model.TransactionLog, len(src_orders))
	dst_transactionLogs[0] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[1] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[2] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[3] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[4] = []*model.TransactionLog{
		{
			BuyerOrderID:  3,
			SellerOrderID: 5,
			Price:         110,
			Quantity:      30,
			Timestamp:     time.Now(),
		},
		{
			BuyerOrderID:  1,
			SellerOrderID: 5,
			Price:         100,
			Quantity:      20,
			Timestamp:     time.Now(),
		},
	}
	dst_transactionLogs[5] = []*model.TransactionLog{
		{
			BuyerOrderID:  1,
			SellerOrderID: 6,
			Price:         100,
			Quantity:      80,
			Timestamp:     time.Now(),
		},
		{
			BuyerOrderID:  2,
			SellerOrderID: 6,
			Price:         100,
			Quantity:      20,
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

	time.Sleep(100 * time.Millisecond)

	// other validation
	assert.Equal(t, 0, pool.buyerHeap.Len())

	assert.Equal(t, 0, pool.sellerHeap.Len())
}

// TestAddOrder_MarketFOKvsLimitROD_case1
// 預先建立賣方限價單, 測試買方市價單撮合邏輯
// Goal:
//   - 依序寫入賣方限價單, 測試 PQ 排序規則(price & timestamp)
//   - 依序寫入買方市價單, 測試撮合機制是否正確
//   - 測試 transaction logs 數量是否正確
//   - 測試買賣方 PQs 剩餘數量與 quantity 是否正確
func TestAddOrder_MarketFOKvsLimitROD_case1(t *testing.T) {
	// prepare data
	src_ordersIdx := 0
	src_orders := []*model.Order{
		{
			ID:           1,
			UserID:       1,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        110,
			Quantity:     100,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           2,
			UserID:       1,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        110,
			Quantity:     20,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           3,
			UserID:       1,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        100,
			Quantity:     30,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           4,
			UserID:       2,
			RoleType:     e.ROLE_BUYER,
			OrderType:    e.ORDER_MARKET,
			DurationType: e.DURATION_FOK,
			Price:        0,
			Quantity:     160,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           5,
			UserID:       2,
			RoleType:     e.ROLE_BUYER,
			OrderType:    e.ORDER_MARKET,
			DurationType: e.DURATION_FOK,
			Price:        0,
			Quantity:     50,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
		{
			ID:           6,
			UserID:       2,
			RoleType:     e.ROLE_BUYER,
			OrderType:    e.ORDER_MARKET,
			DurationType: e.DURATION_FOK,
			Price:        0,
			Quantity:     100,
			Status:       e.STATUS_PENDING,
			Timestamp:    time.Now(),
		},
	}

	dst_transactionLogsIdx := 0
	dst_transactionLogs := make([][]*model.TransactionLog, len(src_orders))
	dst_transactionLogs[0] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[1] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[2] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[3] = make([]*model.TransactionLog, 0)
	dst_transactionLogs[4] = []*model.TransactionLog{
		{
			BuyerOrderID:  5,
			SellerOrderID: 3,
			Price:         100,
			Quantity:      30,
			Timestamp:     time.Now(),
		},
		{
			BuyerOrderID:  5,
			SellerOrderID: 1,
			Price:         110,
			Quantity:      20,
			Timestamp:     time.Now(),
		},
	}
	dst_transactionLogs[5] = []*model.TransactionLog{
		{
			BuyerOrderID:  6,
			SellerOrderID: 1,
			Price:         110,
			Quantity:      80,
			Timestamp:     time.Now(),
		},
		{
			BuyerOrderID:  6,
			SellerOrderID: 2,
			Price:         110,
			Quantity:      20,
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

	time.Sleep(100 * time.Millisecond)

	// other validation
	assert.Equal(t, 0, pool.buyerHeap.Len())

	assert.Equal(t, 0, pool.sellerHeap.Len())
}
