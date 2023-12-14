package order

import (
	"interview20231208/api/router"
	"interview20231208/pkg/trading"

	"github.com/gin-gonic/gin"
)

type orderRouter struct {
	tradingPool trading.TradingPool
}

func NewOrderRouter(tradingPool trading.TradingPool) router.Router {
	return &orderRouter{
		tradingPool: tradingPool,
	}
}

func (o *orderRouter) Init(r *gin.RouterGroup) {
	v1 := r.Group("/v1")

	v1.POST("/orders", o.PendingOrder)
	v1.GET("/orders", o.GetOrders)
}
