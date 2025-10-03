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
