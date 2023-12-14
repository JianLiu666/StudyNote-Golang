package order

import (
	"interview20231208/model"
	"interview20231208/pkg/e"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (o *orderRouter) GetOrders(c *gin.Context) {
	// step.1 validation
	limitStr := c.DefaultQuery("limit", "10") // TODO: magic number
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.String(http.StatusBadRequest, e.GetMsg(e.INVALID_PARAMS))
		return
	}

	userIdStr := c.DefaultQuery("userId", strconv.Itoa(e.UNDIFIED_USERID))
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.String(http.StatusBadRequest, e.GetMsg(e.INVALID_PARAMS))
		return
	}

	statusStr := c.DefaultQuery("status", strconv.Itoa(e.STATUS_PLACEHOLDER))
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		c.String(http.StatusBadRequest, e.GetMsg(e.INVALID_PARAMS))
		return
	}

	var startTimestamp, endTimestamp time.Time
	startTimestampStr := c.Query("startTimestamp")
	if startTimestampStr != "" {
		startTimestamp, err = time.Parse(time.RFC3339, startTimestampStr)
		if err != nil {
			c.String(http.StatusBadRequest, e.GetMsg(e.INVALID_PARAMS))
			return
		}
	}
	endTimestampStr := c.Query("endTimestamp")
	if endTimestampStr != "" {
		endTimestamp, err = time.Parse(time.RFC3339, endTimestampStr)
		if err != nil {
			c.String(http.StatusBadRequest, e.GetMsg(e.INVALID_PARAMS))
			return
		}
	}

	opts := &model.OrderQueryOpts{
		UserID:         userId,
		Status:         e.ORDER_STATUS(status),
		StartTimestamp: startTimestamp,
		EndTimestamp:   endTimestamp,
	}

	// step.2 business logic
	result, code := o.tradingPool.GetOrders(limit, opts)
	if code != e.SUCCESS {
		c.String(http.StatusInternalServerError, e.GetMsg(e.ERROR))
		return
	}

	c.JSON(http.StatusOK, result)
}
