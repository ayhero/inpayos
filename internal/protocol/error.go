package protocol

// ErrorCode 错误码类型
type ErrorCode string

// 系统级错误码 (1000-1999)
const (
	Success         ErrorCode = "0000" // 成功
	SystemError     ErrorCode = "1000" // 系统错误
	DatabaseError   ErrorCode = "1001" // 数据库错误
	CacheError      ErrorCode = "1002" // 缓存错误
	NetworkError    ErrorCode = "1003" // 网络错误
	ServiceUnavail  ErrorCode = "1004" // 服务不可用
	InternalError   ErrorCode = "1005" // 内部错误
	ConfigError     ErrorCode = "1006" // 配置错误
	FileError       ErrorCode = "1007" // 文件操作错误
	JSONParseError  ErrorCode = "1008" // JSON解析错误
	ThirdPartyError ErrorCode = "1009" // 第三方服务错误
)

// 请求相关错误码 (2000-2999)
const (
	InvalidRequest     ErrorCode = "2000" // 无效请求
	InvalidParams      ErrorCode = "2001" // 参数错误
	MissingParams      ErrorCode = "2002" // 缺少参数
	InvalidJSON        ErrorCode = "2003" // JSON格式错误
	InvalidMethod      ErrorCode = "2004" // 请求方法错误
	RequestTooLarge    ErrorCode = "2005" // 请求体过大
	RateLimitExceeded  ErrorCode = "2006" // 请求频率限制
	InvalidContentType ErrorCode = "2007" // 内容类型错误
	InvalidEncoding    ErrorCode = "2008" // 编码错误
	RequestTimeout     ErrorCode = "2009" // 请求超时
)

// 认证相关错误码 (3000-3999)
const (
	AuthenticationFailed    ErrorCode = "3000" // 认证失败
	InvalidToken            ErrorCode = "3001" // 无效令牌
	TokenExpired            ErrorCode = "3002" // 令牌过期
	InvalidCredentials      ErrorCode = "3003" // 凭据无效
	AccessDenied            ErrorCode = "3004" // 访问被拒绝
	PermissionDenied        ErrorCode = "3005" // 权限不足
	AccountDisabled         ErrorCode = "3006" // 账户被禁用
	AccountLocked           ErrorCode = "3007" // 账户被锁定
	LoginRequired           ErrorCode = "3008" // 需要登录
	RefreshTokenExpired     ErrorCode = "3009" // 刷新令牌过期
	InvalidSignature        ErrorCode = "3010" // 签名无效
	InvalidAPIKey           ErrorCode = "3011" // API密钥无效
	InsufficientPermissions ErrorCode = "3012" // 权限不足
	TwoFactorRequired       ErrorCode = "3013" // 需要双因子认证
	InvalidTwoFactorCode    ErrorCode = "3014" // 双因子认证码无效
	AccountSuspended        ErrorCode = "3015" // 账户被暂停
	IPNotAllowed            ErrorCode = "3016" // IP地址不被允许
	SessionLimitExceeded    ErrorCode = "3017" // 会话限制超出
)

// 商户相关错误码 (4000-4999)
const (
	MerchantNotFound      ErrorCode = "4000" // 商户不存在
	MerchantAlreadyExists ErrorCode = "4001" // 商户已存在
	MerchantDisabled      ErrorCode = "4002" // 商户被禁用
	MerchantSuspended     ErrorCode = "4003" // 商户被暂停
	InvalidMerchantID     ErrorCode = "4004" // 商户ID无效
	MerchantNotActive     ErrorCode = "4005" // 商户未激活
	InvalidSecretKey      ErrorCode = "4006" // 密钥无效
	SecretKeyExpired      ErrorCode = "4007" // 密钥过期
	InvalidAppID          ErrorCode = "4008" // 应用ID无效
	AppIDNotFound         ErrorCode = "4009" // 应用ID不存在
)

// 交易相关错误码 (5000-5999)
const (
	TransactionNotFound      ErrorCode = "5000" // 交易不存在
	TransactionExpired       ErrorCode = "5001" // 交易过期
	TransactionCancelled     ErrorCode = "5002" // 交易已取消
	TransactionCompleted     ErrorCode = "5003" // 交易已完成
	TransactionFailed        ErrorCode = "5004" // 交易失败
	TransactionProcessing    ErrorCode = "5005" // 交易处理中
	InvalidTransactionID     ErrorCode = "5006" // 交易ID无效
	DuplicateTransaction     ErrorCode = "5007" // 重复交易
	InvalidAmount            ErrorCode = "5008" // 金额无效
	AmountTooLarge           ErrorCode = "5009" // 金额过大
	AmountTooSmall           ErrorCode = "5010" // 金额过小
	InvalidCurrency          ErrorCode = "5011" // 货币无效
	TransactionLimitExceeded ErrorCode = "5012" // 交易限额超出
)

// 代收相关错误码 (5100-5199)
const (
	ReceiptNotFound    ErrorCode = "5100" // 代收订单不存在
	ReceiptExpired     ErrorCode = "5101" // 代收订单过期
	ReceiptCancelled   ErrorCode = "5102" // 代收订单已取消
	ReceiptCompleted   ErrorCode = "5103" // 代收订单已完成
	ReceiptFailed      ErrorCode = "5104" // 代收失败
	ReceiptProcessing  ErrorCode = "5105" // 代收处理中
	InvalidReceiptData ErrorCode = "5106" // 代收数据无效
)

// 代付相关错误码 (5200-5299)
const (
	PaymentNotFound    ErrorCode = "5200" // 代付订单不存在
	PaymentExpired     ErrorCode = "5201" // 代付订单过期
	PaymentCancelled   ErrorCode = "5202" // 代付订单已取消
	PaymentCompleted   ErrorCode = "5203" // 代付订单已完成
	PaymentFailed      ErrorCode = "5204" // 代付失败
	PaymentProcessing  ErrorCode = "5205" // 代付处理中
	InvalidPaymentData ErrorCode = "5206" // 代付数据无效
	InsufficientFunds  ErrorCode = "5207" // 余额不足
)

// 收银台相关错误码 (5300-5399)
const (
	CheckoutNotFound    ErrorCode = "5300" // 收银台订单不存在
	CheckoutExpired     ErrorCode = "5301" // 收银台订单过期
	CheckoutCancelled   ErrorCode = "5302" // 收银台订单已取消
	CheckoutCompleted   ErrorCode = "5303" // 收银台订单已完成
	CheckoutFailed      ErrorCode = "5304" // 收银台订单失败
	CheckoutProcessing  ErrorCode = "5305" // 收银台订单处理中
	InvalidCheckoutData ErrorCode = "5306" // 收银台数据无效
	CheckoutConfigError ErrorCode = "5307" // 收银台配置错误
)

// 余额相关错误码 (5400-5499)
const (
	BalanceNotFound     ErrorCode = "5400" // 余额记录不存在
	BalanceInsufficient ErrorCode = "5401" // 余额不足
	BalanceFrozen       ErrorCode = "5402" // 余额被冻结
	BalanceQueryFailed  ErrorCode = "5403" // 余额查询失败
	InvalidBalanceType  ErrorCode = "5404" // 余额类型无效
)

// 渠道相关错误码 (6000-6999)
const (
	ChannelNotFound     ErrorCode = "6000" // 渠道不存在
	ChannelDisabled     ErrorCode = "6001" // 渠道被禁用
	ChannelMaintenance  ErrorCode = "6002" // 渠道维护中
	ChannelError        ErrorCode = "6003" // 渠道错误
	ChannelTimeout      ErrorCode = "6004" // 渠道超时
	InvalidChannelID    ErrorCode = "6005" // 渠道ID无效
	ChannelNotSupported ErrorCode = "6006" // 渠道不支持
)

// Webhook相关错误码 (7000-7999)
const (
	WebhookNotFound       ErrorCode = "7000" // Webhook不存在
	WebhookFailed         ErrorCode = "7001" // Webhook发送失败
	WebhookTimeout        ErrorCode = "7002" // Webhook超时
	InvalidWebhookURL     ErrorCode = "7003" // Webhook URL无效
	WebhookConfigError    ErrorCode = "7004" // Webhook配置错误
	WebhookSignatureError ErrorCode = "7005" // Webhook签名错误
)

// 配置相关错误码 (8000-8999)
const (
	ConfigNotFound     ErrorCode = "8000" // 配置不存在
	ConfigInvalid      ErrorCode = "8001" // 配置无效
	ConfigUpdateFailed ErrorCode = "8002" // 配置更新失败
	ConfigLoadFailed   ErrorCode = "8003" // 配置加载失败
)

// GetMessage 获取错误码对应的英文消息 (默认)
func (code ErrorCode) GetMessage() string {
	messages := map[ErrorCode]string{
		// 系统级错误码
		Success:         "Success",
		SystemError:     "System error",
		DatabaseError:   "Database error",
		CacheError:      "Cache error",
		NetworkError:    "Network error",
		ServiceUnavail:  "Service unavailable",
		InternalError:   "Internal error",
		ConfigError:     "Configuration error",
		FileError:       "File operation error",
		JSONParseError:  "JSON parsing error",
		ThirdPartyError: "Third-party service error",

		// 请求相关错误码
		InvalidRequest:     "Invalid request",
		InvalidParams:      "Invalid parameters",
		MissingParams:      "Missing required parameters",
		InvalidJSON:        "Invalid JSON format",
		InvalidMethod:      "Invalid request method",
		RequestTooLarge:    "Request body too large",
		RateLimitExceeded:  "Rate limit exceeded",
		InvalidContentType: "Invalid content type",
		InvalidEncoding:    "Invalid encoding",
		RequestTimeout:     "Request timeout",

		// 认证相关错误码
		AuthenticationFailed:    "Authentication failed",
		InvalidToken:            "Invalid token",
		TokenExpired:            "Token expired",
		InvalidCredentials:      "Invalid credentials",
		AccessDenied:            "Access denied",
		PermissionDenied:        "Permission denied",
		AccountDisabled:         "Account disabled",
		AccountLocked:           "Account locked",
		LoginRequired:           "Login required",
		RefreshTokenExpired:     "Refresh token expired",
		InvalidSignature:        "Invalid signature",
		InvalidAPIKey:           "Invalid API key",
		InsufficientPermissions: "Insufficient permissions",
		TwoFactorRequired:       "Two-factor authentication required",
		InvalidTwoFactorCode:    "Invalid two-factor authentication code",
		AccountSuspended:        "Account suspended",
		IPNotAllowed:            "IP address not allowed",
		SessionLimitExceeded:    "Session limit exceeded",

		// 商户相关错误码
		MerchantNotFound:      "Merchant not found",
		MerchantAlreadyExists: "Merchant already exists",
		MerchantDisabled:      "Merchant disabled",
		MerchantSuspended:     "Merchant suspended",
		InvalidMerchantID:     "Invalid merchant ID",
		MerchantNotActive:     "Merchant not active",
		InvalidSecretKey:      "Invalid secret key",
		SecretKeyExpired:      "Secret key expired",
		InvalidAppID:          "Invalid app ID",
		AppIDNotFound:         "App ID not found",

		// 交易相关错误码
		TransactionNotFound:      "Transaction not found",
		TransactionExpired:       "Transaction expired",
		TransactionCancelled:     "Transaction cancelled",
		TransactionCompleted:     "Transaction completed",
		TransactionFailed:        "Transaction failed",
		TransactionProcessing:    "Transaction processing",
		InvalidTransactionID:     "Invalid transaction ID",
		DuplicateTransaction:     "Duplicate transaction",
		InvalidAmount:            "Invalid amount",
		AmountTooLarge:           "Amount too large",
		AmountTooSmall:           "Amount too small",
		InvalidCurrency:          "Invalid currency",
		TransactionLimitExceeded: "Transaction limit exceeded",

		// 代收相关错误码
		ReceiptNotFound:    "Receipt not found",
		ReceiptExpired:     "Receipt expired",
		ReceiptCancelled:   "Receipt cancelled",
		ReceiptCompleted:   "Receipt completed",
		ReceiptFailed:      "Receipt failed",
		ReceiptProcessing:  "Receipt processing",
		InvalidReceiptData: "Invalid receipt data",

		// 代付相关错误码
		PaymentNotFound:    "Payment not found",
		PaymentExpired:     "Payment expired",
		PaymentCancelled:   "Payment cancelled",
		PaymentCompleted:   "Payment completed",
		PaymentFailed:      "Payment failed",
		PaymentProcessing:  "Payment processing",
		InvalidPaymentData: "Invalid payment data",
		InsufficientFunds:  "Insufficient funds",

		// 收银台相关错误码
		CheckoutNotFound:    "Checkout not found",
		CheckoutExpired:     "Checkout expired",
		CheckoutCancelled:   "Checkout cancelled",
		CheckoutCompleted:   "Checkout completed",
		CheckoutFailed:      "Checkout failed",
		CheckoutProcessing:  "Checkout processing",
		InvalidCheckoutData: "Invalid checkout data",
		CheckoutConfigError: "Checkout configuration error",

		// 余额相关错误码
		BalanceNotFound:     "Balance not found",
		BalanceInsufficient: "Insufficient balance",
		BalanceFrozen:       "Balance frozen",
		BalanceQueryFailed:  "Balance query failed",
		InvalidBalanceType:  "Invalid balance type",

		// 渠道相关错误码
		ChannelNotFound:     "Channel not found",
		ChannelDisabled:     "Channel disabled",
		ChannelMaintenance:  "Channel under maintenance",
		ChannelError:        "Channel error",
		ChannelTimeout:      "Channel timeout",
		InvalidChannelID:    "Invalid channel ID",
		ChannelNotSupported: "Channel not supported",

		// Webhook相关错误码
		WebhookNotFound:       "Webhook not found",
		WebhookFailed:         "Webhook failed",
		WebhookTimeout:        "Webhook timeout",
		InvalidWebhookURL:     "Invalid webhook URL",
		WebhookConfigError:    "Webhook configuration error",
		WebhookSignatureError: "Webhook signature error",

		// 配置相关错误码
		ConfigNotFound:     "Configuration not found",
		ConfigInvalid:      "Invalid configuration",
		ConfigUpdateFailed: "Configuration update failed",
		ConfigLoadFailed:   "Configuration load failed",
	}

	if msg, exists := messages[code]; exists {
		return msg
	}
	return "Unknown error"
}

// GetCode 获取错误码字符串值
func (code ErrorCode) GetCode() string {
	return string(code)
}

// IsSuccess 判断是否成功
func (code ErrorCode) IsSuccess() bool {
	return code == Success
}

// IsSystemError 判断是否系统错误 (1000-1999)
func (code ErrorCode) IsSystemError() bool {
	return code >= "1000" && code < "2000"
}

// IsRequestError 判断是否请求错误 (2000-2999)
func (code ErrorCode) IsRequestError() bool {
	return code >= "2000" && code < "3000"
}

// IsAuthError 判断是否认证错误 (3000-3999)
func (code ErrorCode) IsAuthError() bool {
	return code >= "3000" && code < "4000"
}

// IsMerchantError 判断是否商户错误 (4000-4999)
func (code ErrorCode) IsMerchantError() bool {
	return code >= "4000" && code < "5000"
}

// IsTransactionError 判断是否交易错误 (5000-5999)
func (code ErrorCode) IsTransactionError() bool {
	return code >= "5000" && code < "6000"
}

// IsChannelError 判断是否渠道错误 (6000-6999)
func (code ErrorCode) IsChannelError() bool {
	return code >= "6000" && code < "7000"
}

// IsWebhookError 判断是否Webhook错误 (7000-7999)
func (code ErrorCode) IsWebhookError() bool {
	return code >= "7000" && code < "8000"
}

// IsConfigError 判断是否配置错误 (8000-8999)
func (code ErrorCode) IsConfigError() bool {
	return code >= "8000" && code < "9000"
}

// ServiceError 服务层错误类型
type ServiceError struct {
	Code    ErrorCode
	Message string
}

// Error 实现error接口
func (e *ServiceError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return string(e.Code)
}

// NewServiceError 创建服务错误
func NewServiceError(code ErrorCode, message string) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
	}
}

// IsServiceError 判断是否为ServiceError类型
func IsServiceError(err error) (*ServiceError, bool) {
	if serviceErr, ok := err.(*ServiceError); ok {
		return serviceErr, true
	}
	return nil, false
}
