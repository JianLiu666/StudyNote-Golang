package rdb

import (
	"context"
	"interview20231208/model"
)

const (
	TbOrders          = "orders"
	TbTransactionLogs = "transactionLogs"
)

type RDB interface {
	// Shutdown 提供外部調用 graceful shutdown 流程
	//
	// @param ctx context for tracing
	Shutdown(ctx context.Context)

	// CreateOrder 建立交易單
	// 對 RDB 寫入交易單
	//
	// @param ctx context for tracing
	//
	// @param order 欲建立的交易單
	CreateOrder(ctx context.Context, order *model.Order)

	// GetOrders 取得交易單
	// 根據傳入的 filters 組裝 SQL 後取回符合條件的交易單
	//
	// @param ctx context for tracing
	//
	// @param opts 查詢條件
	//
	// @return []*model.Order 符合條件的交易單
	GetOrders(ctx context.Context, opts *model.OrderQueryOpts) []*model.Order

	// UpdateOrdersAndCreateTransactionLogs 更新交易紀錄
	// 事務應包含以下操作:
	//  - 根據給定的交易單更新最新狀態(e.g. 已完成, 已取消)
	//  - 根據給定的撮合紀錄寫入 RDB
	//
	// @param ctx context for tracing
	//
	// @param orders 欲更新的交易單與交易單狀態
	//
	// @param logs 欲寫入的撮合紀錄
	UpdateOrdersAndCreateTransactionLogs(ctx context.Context, orders map[int]*model.Order, logs []*model.TransactionLog)

	// GetTransactionLogs 取得撮合紀錄
	// 根據傳入的 filters 組裝 SQL 後取回符合條件的撮合紀錄
	//
	// @param ctx context for tracing
	//
	// @param opts 查詢條件
	//
	// @return []*model.TransactionLog 符合條件的撮合紀錄
	GetTransactionLogs(ctx context.Context, opts *model.TransactionLogQueryOpts) []*model.TransactionLog
}
