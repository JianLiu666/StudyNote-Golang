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
