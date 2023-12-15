package trading

import (
	"context"
	"interview20231208/model"
	"interview20231208/pkg/e"
)

type TradingPool interface {
	// Enable 開始啟動股票撮合模組
	//
	// @param ctx context for tracing
	Enable(ctx context.Context)

	// AddOrder 建立一筆新的掛單
	//
	// @param order 欲建立的交易單
	//
	// @return e.Code 執行結果
	AddOrder(order *model.Order) e.CODE

	// GetOrders 取得交易單
	// 根據傳入的 filters 取回符合條件的交易單
	//
	// @param limit 取回的交易單數量上限
	//
	// @param opts 查詢條件
	//
	// @return []*model.Order 符合條件的交易單
	//
	// @return e.Code 執行結果
	GetOrders(limit int, opts *model.OrderQueryOpts) ([]*model.Order, e.CODE)

	// GetTransactionLogs 取得撮合紀錄
	// 根據傳入的 filters 取回符合條件的撮合紀錄
	//
	// @param limit 取回的撮合紀錄數量上限
	//
	// @param opts 查詢條件
	//
	// @return []*model.TransactionLog 符合條件的撮合紀錄
	//
	// @return e.Code 執行結果
	GetTransactionLogs(limit int, opts *model.TransactionLogQueryOpts) ([]*model.TransactionLog, e.CODE)
}
