package protocol

// 消息渠道常量
const (
	MsgChannelEmail = "email"
	MsgChannelFcm   = "fcm"
	MsgChannelSms   = "sms"
)

// 消息类型常量
const (
	// 通用消息类型
	MsgTypeGeneric             = "generic"
	MsgTypeVerifyCode          = "verify_code"
	MsgTypeRegisterSuccess     = "register_success"
	MsgTypePasswordReset       = "password_reset"
	MsgTypeAccountVerification = "account_verification"
	MsgTypePasswordUpdate      = "password_update"
	MsgTypeNewPassword         = "new_password" // 新密码邮件
)

// 语言常量
const (
	LangEnglish     = "en"
	LangFrench      = "fr"
	LangChinese     = "zh"
	LangKinyarwanda = "rw"
)

// 地区常量
const (
	RegionRwanda   = "RW"
	RegionKenya    = "KE"
	RegionTanzania = "TZ"
	RegionUganda   = "UG"
)
