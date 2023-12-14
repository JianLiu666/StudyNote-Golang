package order

import (
	"interview20231208/model"
	"interview20231208/pkg/e"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (o *orderRouter) PendingOrder(c *gin.Context) {
	order := &model.Order{}

	// step.1 validation
	if err := c.BindJSON(order); err != nil {
		logrus.Errorf("Failed to execute c.BindJSON: %v", err)
		c.String(http.StatusBadRequest, e.GetMsg(e.INVALID_PARAMS))
		return
	}
	if !e.IsRoleTypeValid(int(order.RoleType)) {
		c.String(http.StatusBadRequest, e.GetMsg(e.INVALID_PARAMS))
		return
	}
	if !e.IsOrderTypeValid(int(order.OrderType)) {
		c.String(http.StatusBadRequest, e.GetMsg(e.INVALID_PARAMS))
		return
	}
	if !e.IsDurationTypeValid(int(order.DurationType)) {
		c.String(http.StatusBadRequest, e.GetMsg(e.INVALID_PARAMS))
		return
	}

	// FIXME: 交易尚未實作, 暫時保留
	if order.OrderType == e.ORDER_LIMIT && order.DurationType == e.DURATION_IOC {
		c.String(http.StatusBadRequest, e.GetMsg(e.INVALID_PARAMS))
		return
	}
	if order.OrderType == e.ORDER_MARKET && order.DurationType == e.DURATION_ROD {
		c.String(http.StatusBadRequest, e.GetMsg(e.INVALID_PARAMS))
		return
	}
	if order.OrderType == e.ORDER_MARKET && order.DurationType == e.DURATION_IOC {
		c.String(http.StatusBadRequest, e.GetMsg(e.INVALID_PARAMS))
		return
	}

	// step.2 population fields
	order.Timestamp = time.Now()

	// step.3 business logic
	if code := o.tradingPool.AddOrder(order); code != e.SUCCESS {
		c.String(http.StatusBadRequest, e.GetMsg(code))
		return
	}

	c.Status(http.StatusOK)
}
