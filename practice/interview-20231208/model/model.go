package model

import (
	"interview20231208/pkg/e"
	"time"
)

type Order struct {
	ID           int             `json:"id"           gorm:"column:id;primaryKey"` // 交易單唯一識別碼
	UserID       int             `json:"userId"       gorm:"column:userId"`        // 用戶唯一識別碼
	RoleType     e.ROLE_TYPE     `json:"roleType"     gorm:"column:roleType"`      // 掛單角色(e.g. 買方/賣方)
	OrderType    e.ORDER_TYPE    `json:"orderType"    gorm:"column:orderType"`     // 交易單類型(e.g. 市價單/限價單)
	DurationType e.DURATION_TYPE `json:"durationType" gorm:"column:durationType"`  // 交易單期限(e.g. ROD/IOC/FOK)
	Price        int             `json:"price"        gorm:"column:price"`         // 交易單價格
	Quantity     int             `json:"quantity"     gorm:"column:quantity"`      // 交易單數量
	Status       e.ORDER_STATUS  `json:"status"       gorm:"column:status"`        // 交易單狀態
	Timestamp    time.Time       `json:"timestamp"    gorm:"column:timestamp"`     // 交易單時間戳
}

type OrderQueryOpts struct {
	UserID         int
	Status         e.ORDER_STATUS
	StartTimestamp time.Time
	EndTimestamp   time.Time
}

type TransactionLog struct {
	ID            int       `json:"id"            gorm:"column:id;primaryKey"` // 成交紀錄唯一識別碼
	BuyerOrderID  int       `json:"buyerOrderId"  gorm:"column:buyerOrderId"`  // 買方唯一識別碼
	SellerOrderID int       `json:"sellerOrderId" gorm:"column:sellerOrderId"` // 賣方唯一識別罵
	Price         int       `json:"price"         gorm:"column:price"`         // 成交價格
	Quantity      int       `json:"quantity"      gorm:"column:quantity"`      // 成交數量
	Timestamp     time.Time `json:"timestamp"     gorm:"column:timestamp"`     // 成交時間戳
}

type TransactionLogQueryOpts struct {
	BuyerOrderID   int
	SellerOrderID  int
	StartTimestamp time.Time
	EndTimestamp   time.Time
}
