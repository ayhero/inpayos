package protocol

// 用户类型常量
const (
	UserTypeMerchant = "merchant"
	UserTypeCashier  = "cashier"
	UserTypeBank     = "bank"
)

// 交易类型常量
const (
	TxTypeReceipt  = "receipt"
	TxTypePayment  = "payment"
	TxTypeRefund   = "refund"
	TxTypeDeposit  = "deposit"
	TxTypeWithdraw = "withdraw"
)

// 通用状态常量
const (
	StatusOn  = "on"
	StatusOff = "off"
)

// 业务状态常量
const (
	StatusActive     = "active"
	StatusInactive   = "inactive"
	StatusSuspended  = "suspended"
	StatusDeleted    = "deleted"
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusSuccess    = "success"
	StatusFailed     = "failed"
	StatusCancelled  = "cancelled"
	StatusExpired    = "expired"
	StatusCompleted  = "completed"
	StatusApproved   = "approved"
	StatusRejected   = "rejected"
)

// 流水类型常量
const (
	FlowTypeIncome   = "income"
	FlowTypeExpense  = "expense"
	FlowTypeFreeze   = "freeze"
	FlowTypeUnfreeze = "unfreeze"
	FlowTypeMargin   = "margin"
)

// 渠道相关常量
const (
	ChannelStatusActive   = "active"
	ChannelStatusInactive = "inactive"
	ChannelStatusMaintain = "maintain"
)
