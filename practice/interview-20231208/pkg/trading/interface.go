package trading

import (
	"context"
	"interview20231208/model"
)

type TradingPool interface {
	Enable(ctx context.Context)
	AddOrder(order *model.Order)
}
