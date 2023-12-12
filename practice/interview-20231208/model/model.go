package model

import "interview20231208/pkg/e"

type Order struct {
	ID           string          `json:"id"`           // 交易單唯一識別碼
	UserID       int             `json:"userId"`       // 用戶唯一識別碼
	RoleType     e.ROLE_TYPE     `json:"roleType"`     // 掛單角色(e.g. 買方/賣方)
	OrderType    e.ORDER_TYPE    `json:"orderType"`    // 交易單類型(e.g. 市價單/限價單)
	DurationType e.DURATION_TYPE `json:"durationType"` // 交易單期限(e.g. ROD/IOC/FOK)
	Price        int             `json:"price"`        // 交易單價格
	Quantity     int             `json:"quantity"`     // 交易單數量
	Status       e.ORDER_STATUS  `json:"status"`       // 交易單狀態
	Timestamp    int             `json:"timestamp"`    // 交易單時間戳
}

type TransactionLog struct {
	ID            string `json:"id"`            // 成交紀錄唯一識別碼
	BuyerOrderID  string `json:"buyerOrderId"`  // 買方唯一識別碼
	SellerOrderID string `json:"sellerOrderId"` // 賣方唯一識別罵
	Price         int    `json:"price"`         // 成交價格
	Quantity      int    `json:"quantity"`      // 成交數量
	Timestamp     int    `json:"timstamp"`      // 成交時間戳
}
