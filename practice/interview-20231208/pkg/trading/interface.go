package trading

import (
	"context"
	"interview20231208/model"
	"interview20231208/pkg/e"
)

type TradingPool interface {
	Enable(ctx context.Context)

	AddOrder(order *model.Order) e.CODE

	GetOrders(limit int, opts *model.OrderQueryOpts) ([]*model.Order, e.CODE)

	GetTransactionLogs(limit int, opts *model.TransactionLogQueryOpts) ([]*model.TransactionLog, e.CODE)
}
