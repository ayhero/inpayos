package protocol

const (
	System  = "system"
	Default = "default"
)

const (
	EnvSandbox    = "sandbox"
	EnvProduction = "production"
)

// 用户类型常量
const (
	UserTypeMerchant    = "merchant"
	UserTypeCashier     = "cashier"
	UserTypeCashierTeam = "cashier_team"
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
	StatusEnabled    = "enabled"  // 启用
	StatusDisabled   = "disabled" // 禁用
	StatusOnline     = "online"
	StatusOffline    = "offline"
	StatusBusy       = "busy"
	StatusLocked     = "locked"
)

// 流水类型常量
const (
	FlowTypeIncome   = "income"
	FlowTypeExpense  = "expense"
	FlowTypeFreeze   = "freeze"
	FlowTypeUnfreeze = "unfreeze"
	FlowTypeMargin   = "margin"
)

const (
	VerifyCodeTypeRegister      = "register"       // 注册验证码
	VerifyCodeTypeLogin         = "login"          // 登录验证码
	VerifyCodeTypeResetPassword = "reset_password" // 重置密码验证码
	VerifyCodeTypeResetG2FA     = "reset_g2fa"     // 重置G2FA验证码
	VerifyCodeTypeReset         = "reset"          // 重置密码验证码
)

// 订单类型
const (
	TrxTypePayin     = "payin"      // 代收订单
	TrxTypePayout    = "payout"     // 代付订单
	TrxTypeWithdraw  = "withdraw"   // 提现订单
	TrxTypeRefund    = "refund"     // 退款订单
	TrxTypeRecharge  = "recharge"   // 充值订单
	TrxTypeFreeze    = "freeze"     // 冻结订单
	TrxTypeUnfreeze  = "unfreeze"   // 解冻订单
	TrxTypeRfRecover = "rf_recover" // 退款回撤订单
	TrxTypeWdRecover = "wd_recover" // 提现回撤订单
)

const (
	WeightTypePriority = "POT"
	WeightTypePercent  = "PCT"
)

// Admin 状态常量
const (
	AdminStatusActive    = StatusActive
	AdminStatusInactive  = StatusInactive
	AdminStatusSuspended = StatusSuspended
	AdminStatusLocked    = StatusLocked
)

// Admin 活跃状态常量
const (
	AdminActiveStatusOnline  = StatusOnline
	AdminActiveStatusOffline = StatusOffline
	AdminActiveStatusBusy    = StatusBusy
)
