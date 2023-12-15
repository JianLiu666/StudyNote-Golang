package transaction

import (
	"bytes"
	"context"
	"encoding/json"
	"interview20231208/api/router/v1/order"
	"interview20231208/model"
	"interview20231208/pkg/accessor"
	"interview20231208/pkg/e"
	"interview20231208/pkg/rdb"
	trading "interview20231208/pkg/trading/heap"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gotest.tools/assert"
)

var _router *gin.Engine
var _sqlServer *rdb.SqlServer

func TestMain(m *testing.M) {
	infra := accessor.Build()

	_sqlServer = rdb.NewSqlServer(infra.Config.MySQL.Address)
	_sqlServer.Enable()

	infra.InitRDB(context.TODO())
	infra.InitTradingPool(context.TODO(), trading.NewTradingPool(infra.RDB))

	_router = gin.Default()
	api := _router.Group("/api")

	order.NewOrderRouter(infra.TradingPool).Init(api)
	NewTransactionRouter(infra.TradingPool).Init(api)

	os.Exit(m.Run())
}

func TestGetTransaction_Success(t *testing.T) {
	// init environment
	_sqlServer.Clear()

	// prepare data
	orders := []*model.Order{
		{
			UserID:       1,
			RoleType:     e.ROLE_BUYER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        100,
			Quantity:     100,
		},
		{
			UserID:       2,
			RoleType:     e.ROLE_SELLER,
			OrderType:    e.ORDER_LIMIT,
			DurationType: e.DURATION_ROD,
			Price:        100,
			Quantity:     100,
		},
	}

	// testcase
	for _, order := range orders {
		data, err := json.Marshal(order)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		_router.ServeHTTP(w, req)
	}
	time.Sleep(100 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions", nil)
	w := httptest.NewRecorder()
	_router.ServeHTTP(w, req)

	result := []*model.TransactionLog{}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// validation
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, orders[0].Price, result[0].Price)
	assert.Equal(t, orders[0].Quantity, result[0].Quantity)
}
