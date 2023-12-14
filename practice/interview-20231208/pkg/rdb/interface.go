package rdb

import (
	"context"
	"interview20231208/model"
)

const (
	TbOrders          = "orders"
	TbTransactionLogs = "transactionLogs"
)

type RDB interface {
	Shutdown(ctx context.Context)

	CreateOrder(ctx context.Context, order *model.Order)

	GetOrders(ctx context.Context, opts *model.OrderQueryOpts) []*model.Order

	UpdateOrdersAndCreateTransactionLogs(ctx context.Context, orders map[int]*model.Order, logs []*model.TransactionLog)

	GetTransactionLogs(ctx context.Context, opts *model.TransactionLogQueryOpts) []*model.TransactionLog
}
