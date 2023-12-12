package mysql

import (
	"context"
	"interview20231208/model"
)

type mockMysqlClient struct {
	createOrderCallback                              func(ctx context.Context, order *model.Order)
	updateOrdersAndCreateTransactionLogsCallbackfunc func(ctx context.Context, logs []*model.TransactionLog)
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

func (c *mockMysqlClient) SetUpdateOrdersAndCreateTransactionLogs(f func(ctx context.Context, logs []*model.TransactionLog)) {
	c.updateOrdersAndCreateTransactionLogsCallbackfunc = f
}

func (c *mockMysqlClient) UpdateOrdersAndCreateTransactionLogs(ctx context.Context, logs []*model.TransactionLog) {
	c.updateOrdersAndCreateTransactionLogsCallbackfunc(ctx, logs)
}
