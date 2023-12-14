package transaction

import (
	"interview20231208/api/router"
	"interview20231208/pkg/trading"

	"github.com/gin-gonic/gin"
)

type transactionRouter struct {
	tradingPool trading.TradingPool
}

func NewTransactionRouter(tradingPool trading.TradingPool) router.Router {
	return &transactionRouter{
		tradingPool: tradingPool,
	}
}

func (t *transactionRouter) Init(r *gin.RouterGroup) {
	v1 := r.Group("/v1")

	v1.GET("/transactions", t.GetTransactionLogs)
}
