package trading

import (
	"context"
	"interview20231208/model"
	"interview20231208/pkg/e"
	"interview20231208/pkg/rdb"
	"interview20231208/pkg/rdb/mysql"
	"os"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestMain(m *testing.M) {
	sqlServer := rdb.NewSqlServer("localhost:3306")
	sqlServer.Enable()

	os.Exit(m.Run())
}

// TestAddOrder_LimitRODvsLimitROD_case1
// 買方限價單完全成交, 賣方限價單部分成交
// Goal:
//   - 依序寫入買方限價單, 測試 PQ 排序規則(price & timestamp)
//   - 依序寫入賣方限價單, 測試撮合機制是否正確
//   - 測試 transaction logs 數量是否正確
//   - 測試買賣方 PQs 剩餘數量與 quantity 是否正確
func TestAddOrder_LimitRODvsLimitROD_case1(t *testing.T) {
	// prepare data
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

	// prepare environment
	dsn := "root:0@tcp(localhost:3306)/trading?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlClient := mysql.NewMySqlClient(context.TODO(), dsn, 60, 100, 10)
	pool := newTradingPool(mysqlClient)
	pool.Enable(context.Background())

	// testcase flow
	for _, order := range src_orders {
		pool.AddOrder(order)
	}
	time.Sleep(100 * time.Millisecond)

	// validation
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

	// prepare environment
	dsn := "root:0@tcp(localhost:3306)/trading?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlClient := mysql.NewMySqlClient(context.TODO(), dsn, 60, 100, 10)
	pool := newTradingPool(mysqlClient)
	pool.Enable(context.Background())

	// testcase flow
	for _, order := range src_orders {
		pool.AddOrder(order)
	}
	time.Sleep(100 * time.Millisecond)

	// validation
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

	// prepare environment
	dsn := "root:0@tcp(localhost:3306)/trading?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlClient := mysql.NewMySqlClient(context.TODO(), dsn, 60, 100, 10)
	pool := newTradingPool(mysqlClient)
	pool.Enable(context.Background())

	// testcase flow
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

	// prepare environment
	dsn := "root:0@tcp(localhost:3306)/trading?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlClient := mysql.NewMySqlClient(context.TODO(), dsn, 60, 100, 10)
	pool := newTradingPool(mysqlClient)
	pool.Enable(context.Background())

	// testcase flow
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

	// prepare environment
	dsn := "root:0@tcp(localhost:3306)/trading?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlClient := mysql.NewMySqlClient(context.TODO(), dsn, 60, 100, 10)
	pool := newTradingPool(mysqlClient)
	pool.Enable(context.Background())

	// testcase flow
	for _, order := range src_orders {
		pool.AddOrder(order)
	}
	time.Sleep(100 * time.Millisecond)

	// other validation
	assert.Equal(t, 0, pool.buyerHeap.Len())

	assert.Equal(t, 0, pool.sellerHeap.Len())
}
