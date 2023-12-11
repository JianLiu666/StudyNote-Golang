package e

type ROLE_TYPE int // 交易角色

const (
	ROLE_BUYER  = iota // 買方
	ROLE_SELLER        // 賣方
)

type ORDER_TYPE int // 訂單類型

const (
	ORDER_MARKET = iota // 市價單
	ORDER_LIMIT         // 限價單
)

type DURATION_TYPE int // 有效期限

const (
	DURATION_ROD = iota // 當日委託有效單 (Rest of Day)
	DURATION_IOC        // 立即成交否則取消 (Immediate or Cancel)
	DURATION_FOK        // 立即全部成交否則取消 (Fill or Kill)
)
