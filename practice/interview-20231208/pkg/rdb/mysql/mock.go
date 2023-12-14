package mysql

import (
	"context"
	"interview20231208/model"
	"interview20231208/pkg/rdb"
)

var _ rdb.RDB = (*mockMysqlClient)(nil)

type mockMysqlClient struct {
	createOrderCallback                              func(ctx context.Context, order *model.Order)
	updateOrdersAndCreateTransactionLogsCallbackfunc func(ctx context.Context, orders map[int]*model.Order, logs []*model.TransactionLog)
}

func NewMockMysqlClient() *mockMysqlClient {
	return &mockMysqlClient{}
}

func (c *mockMysqlClient) Shutdown(ctx context.Context) {

}

func (c *mockMysqlClient) SetCreateOrderCallback(f func(ctx context.Context, order *model.Order)) {
	c.createOrderCallback = f
}

func (c *mockMysqlClient) CreateOrder(ctx context.Context, order *model.Order) {
	c.createOrderCallback(ctx, order)
}

func (c *mockMysqlClient) SetUpdateOrdersAndCreateTransactionLogs(f func(ctx context.Context, orders map[int]*model.Order, logs []*model.TransactionLog)) {
	c.updateOrdersAndCreateTransactionLogsCallbackfunc = f
}

func (c *mockMysqlClient) UpdateOrdersAndCreateTransactionLogs(ctx context.Context, orders map[int]*model.Order, logs []*model.TransactionLog) {
	c.updateOrdersAndCreateTransactionLogsCallbackfunc(ctx, orders, logs)
}

func (c *mockMysqlClient) GetOrders(ctx context.Context, opts *model.OrderQueryOpts) []*model.Order {
	result := []*model.Order{}
	return result
}

func (c *mockMysqlClient) GetTransactionLogs(ctx context.Context, opts *model.TransactionLogQueryOpts) []*model.TransactionLog {
	result := []*model.TransactionLog{}
	return result
}
