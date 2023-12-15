package order

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	NewOrderRouter(infra.TradingPool).Init(api)

	os.Exit(m.Run())
}

func TestPendingOrder_Success(t *testing.T) {
	// init environment
	_sqlServer.Clear()

	// prepare data
	order := &model.Order{
		UserID:       1,
		RoleType:     e.ROLE_BUYER,
		OrderType:    e.ORDER_LIMIT,
		DurationType: e.DURATION_ROD,
		Price:        100,
		Quantity:     100,
	}
	data, err := json.Marshal(order)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// testcase
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	_router.ServeHTTP(w, req)

	// validation
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestPendingOrder_BadRequest_case1(t *testing.T) {
	// init environment
	_sqlServer.Clear()

	// prepare data
	order := &model.Order{
		UserID:       1,
		RoleType:     3,
		OrderType:    e.ORDER_LIMIT,
		DurationType: e.DURATION_ROD,
		Price:        100,
		Quantity:     100,
	}
	data, err := json.Marshal(order)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// testcase
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	_router.ServeHTTP(w, req)

	// validation
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestGetOrder_Success(t *testing.T) {
	// init environment
	_sqlServer.Clear()

	// prepare data
	order := &model.Order{
		UserID:       1,
		RoleType:     e.ROLE_BUYER,
		OrderType:    e.ORDER_LIMIT,
		DurationType: e.DURATION_ROD,
		Price:        100,
		Quantity:     100,
	}
	data, err := json.Marshal(order)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// testcase
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	_router.ServeHTTP(w, req)

	time.Sleep(100 * time.Millisecond)

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/orders?userId=%v", order.UserID), nil)
	w = httptest.NewRecorder()
	_router.ServeHTTP(w, req)

	result := []*model.Order{}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// validation
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, order.UserID, result[0].UserID)
	assert.Equal(t, order.RoleType, result[0].RoleType)
	assert.Equal(t, order.OrderType, result[0].OrderType)
	assert.Equal(t, order.DurationType, result[0].DurationType)
	assert.Equal(t, order.Price, result[0].Price)
	assert.Equal(t, order.Quantity, result[0].Quantity)
	assert.Equal(t, e.ORDER_STATUS(e.STATUS_PENDING), result[0].Status)
}

func TestGetOrder_BadRequest_case1(t *testing.T) {
	// init environment
	_sqlServer.Clear()

	// prepare data
	order := &model.Order{
		UserID:       1,
		RoleType:     e.ROLE_BUYER,
		OrderType:    e.ORDER_LIMIT,
		DurationType: e.DURATION_ROD,
		Price:        100,
		Quantity:     100,
	}
	data, err := json.Marshal(order)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// testcase
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	_router.ServeHTTP(w, req)

	time.Sleep(100 * time.Millisecond)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/orders?status=aaa", nil)
	w = httptest.NewRecorder()
	_router.ServeHTTP(w, req)

	// validation
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}
