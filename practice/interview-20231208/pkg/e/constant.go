package e

const UNDIFIED_USERID = -1

type ROLE_TYPE int // 交易角色

const (
	ROLE_BUYER       ROLE_TYPE = iota // 買方
	ROLE_SELLER                       // 賣方
	ROLE_PLACEHOLDER                  // placeholder
)

func IsRoleTypeValid(n ROLE_TYPE) bool {
	return n >= 0 && n < ROLE_PLACEHOLDER
}

type ORDER_TYPE int // 訂單類型

const (
	ORDER_MARKET      ORDER_TYPE = iota // 市價單
	ORDER_LIMIT                         // 限價單
	ORDER_PLACEHOLDER                   // placeholder
)

func IsOrderTypeValid(n ORDER_TYPE) bool {
	return n >= 0 && n < ORDER_PLACEHOLDER
}

type DURATION_TYPE int // 有效期限

const (
	DURATION_ROD         DURATION_TYPE = iota // 當日委託有效單 (Rest of Day)
	DURATION_IOC                              // 立即成交否則取消 (Immediate or Cancel)
	DURATION_FOK                              // 立即全部成交否則取消 (Fill or Kill)
	DURATION_PLACEHOLDER                      // placeholder
)

func IsDurationTypeValid(n DURATION_TYPE) bool {
	return n >= 0 && n < DURATION_PLACEHOLDER
}

type ORDER_STATUS int

const (
	STATUS_PENDING     ORDER_STATUS = iota // 掛單中
	STATUS_CANCEL                          // 已取消
	STATUS_SUCCESS                         // 已成交
	STATUS_PLACEHOLDER                     // placeholder
)

func IsOrderStatusValid(n ORDER_STATUS) bool {
	return n >= 0 && n < STATUS_PLACEHOLDER
}
