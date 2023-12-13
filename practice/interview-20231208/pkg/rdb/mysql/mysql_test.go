// go:build: !unittest

package mysql

import (
	"context"
	"interview20231208/model"
	"interview20231208/pkg/e"
	"testing"
	"time"
)

func TestCreateOrder(t *testing.T) {
	// TODO: remove magic number
	dsn := "root:0@tcp(localhost:3306)/trading?charset=utf8mb4&parseTime=True&loc=Local"
	client := NewMySqlClient(context.TODO(), dsn, 60, 100, 10)

	client.CreateOrder(context.TODO(), &model.Order{
		UserID:       1,
		RoleType:     e.ROLE_BUYER,
		OrderType:    e.ORDER_LIMIT,
		DurationType: e.DURATION_ROD,
		Price:        100,
		Quantity:     100,
		Status:       e.STATUS_PENDING,
		Timestamp:    time.Now(),
	})
}

func TestUpdateOrdersAndCreateTransactionLogs(t *testing.T) {
	// prepare data
	orderSet := map[int]*model.Order{
		1: {ID: 1, Status: e.STATUS_SUCCESS},
		2: {ID: 2, Status: e.STATUS_SUCCESS},
	}
	logs := []*model.TransactionLog{
		{
			BuyerOrderID:  1,
			SellerOrderID: 2,
			Price:         100,
			Quantity:      100,
			Timestamp:     time.Now(),
		},
	}

	// TODO: remove magic number
	dsn := "root:0@tcp(localhost:3306)/trading?charset=utf8mb4&parseTime=True&loc=Local"
	client := NewMySqlClient(context.TODO(), dsn, 60, 100, 10)

	// gen fake data
	client.CreateOrder(context.TODO(), &model.Order{
		UserID:       1,
		RoleType:     e.ROLE_BUYER,
		OrderType:    e.ORDER_LIMIT,
		DurationType: e.DURATION_ROD,
		Price:        100,
		Quantity:     100,
		Status:       e.STATUS_PENDING,
		Timestamp:    time.Now(),
	})

	client.CreateOrder(context.TODO(), &model.Order{
		UserID:       2,
		RoleType:     e.ROLE_SELLER,
		OrderType:    e.ORDER_LIMIT,
		DurationType: e.DURATION_ROD,
		Price:        100,
		Quantity:     100,
		Status:       e.STATUS_PENDING,
		Timestamp:    time.Now(),
	})

	// testcase
	client.UpdateOrdersAndCreateTransactionLogs(context.TODO(), orderSet, logs)
}
