package transaction

import (
	"interview20231208/model"
	"interview20231208/pkg/e"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (t *transactionRouter) GetTransactionLogs(c *gin.Context) {
	// step.1 validation
	limitStr := c.DefaultQuery("limit", "10") // TODO: magic number
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.String(http.StatusBadRequest, e.GetMsg(e.INVALID_PARAMS))
		return
	}

	buyerOrderIdStr := c.DefaultQuery("buyerOrderId", strconv.Itoa(e.UNDIFIED_USERID))
	buyerOrderId, err := strconv.Atoi(buyerOrderIdStr)
	if err != nil {
		c.String(http.StatusBadRequest, e.GetMsg(e.INVALID_PARAMS))
		return
	}

	sellerOrderIdStr := c.DefaultQuery("sellerOrderId", strconv.Itoa(e.UNDIFIED_USERID))
	sellerOrderId, err := strconv.Atoi(sellerOrderIdStr)
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

	opts := &model.TransactionLogQueryOpts{
		BuyerOrderID:   buyerOrderId,
		SellerOrderID:  sellerOrderId,
		StartTimestamp: startTimestamp,
		EndTimestamp:   endTimestamp,
	}

	// step.2 business logic
	result, code := t.tradingPool.GetTransactionLogs(limit, opts)
	if code != e.SUCCESS {
		c.String(http.StatusInternalServerError, e.GetMsg(e.ERROR))
		return
	}

	c.JSON(http.StatusOK, result)
}
