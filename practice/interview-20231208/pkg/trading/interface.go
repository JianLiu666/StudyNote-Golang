package trading

import (
	"context"
	"interview20231208/model"
	"interview20231208/pkg/e"
)

type TradingPool interface {
	Enable(ctx context.Context)
	AddOrder(order *model.Order) e.CODE
}
