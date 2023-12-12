package rdb

import (
	"context"
	"interview20231208/model"
)

type RDB interface {
	Shutdown(ctx context.Context)

	CreateOrder(ctx context.Context, order *model.Order)

	UpdateOrdersAndCreateTransactionLogs(ctx context.Context, logs []*model.TransactionLog)
}
