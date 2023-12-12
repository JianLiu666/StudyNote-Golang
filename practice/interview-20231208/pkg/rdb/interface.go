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

	UpdateOrdersAndCreateTransactionLogs(ctx context.Context, orders map[string]*model.Order, logs []*model.TransactionLog)
}
