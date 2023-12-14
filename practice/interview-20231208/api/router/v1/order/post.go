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

	// validation
	if err := c.BindJSON(order); err != nil {
		logrus.Errorf("Failed to execute c.BindJSON: %v", err)
		c.String(http.StatusBadRequest, e.GetMsg(e.INVALID_PARAMS))
		return
	}

	// population fields
	order.Timestamp = time.Now()

	// business logic
	if code := o.tradingPool.AddOrder(order); code != e.SUCCESS {
		c.String(http.StatusBadRequest, e.GetMsg(code))
		return
	}

	c.Status(http.StatusOK)
}
