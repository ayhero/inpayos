package protocol

// TrxMethod 交易方式常量定义
const (
	TrxMethodUPI          = "upi"
	TrxMethodBankCard     = "bank_card"
	TrxMethodBankTransfer = "bank_transfer"
	TrxMethodUSDT         = "usdt"
	TrxMethodWallet       = "wallet"
	TrxMethodQRCode       = "qr_code"
)

// MerchantCheckoutConfig 商户收银台配置信息
type MerchantCheckoutConfig struct {
	MerchantID string                            `json:"merchant_id"`
	Countries  []string                          `json:"countries"`         // 支持的国家列表
	Configs    map[string]*CountryCheckoutConfig `json:"configs,omitempty"` // 按国家分组的配置
}

// CountryCheckoutConfig 国家级收银台配置信息
type CountryCheckoutConfig struct {
	Country    string                      `json:"country"`
	TrxMethods []string                    `json:"trx_methods"`       // 支持的交易方式列表
	Configs    map[string]*TrxMethodConfig `json:"configs,omitempty"` // 按交易方式分组的配置
}

// TrxMethodConfig 交易方式配置信息
type TrxMethodConfig struct {
	Country   string   `json:"country"`
	TrxMethod string   `json:"trx_method"`         // 交易方式
	Ccy       []string `json:"ccy"`                // 支持的币种列表
	LogoURL   string   `json:"logo_url,omitempty"` // Logo地址
}

// CheckoutServiceListRequest 获取收银台服务列表请求
type CheckoutServiceListRequest struct {
	CheckoutID string `json:"checkout_id" form:"checkout_id" binding:"required"`
	Country    string `json:"country" form:"country"`   // 可选的国家过滤
	Currency   string `json:"currency" form:"currency"` // 可选的币种过滤
}

type CheckoutInfoRequest struct {
	CheckoutID string `json:"checkout_id" binding:"required"`
}

type CreateCheckoutRequest struct {
	Mid       string `json:"mid"`
	ReqID     string `json:"req_id" binding:"required"`
	Ccy       string `json:"ccy" binding:"required"`
	Amount    string `json:"amount" binding:"required"`
	Country   string `json:"country"`
	TrxMethod string `json:"trx_method" `
	NotifyURL string `json:"notify_url"`
	ReturnURL string `json:"return_url"`
}
type ConfirmCheckoutRequest struct {
	Mid        string   `json:"mid"`
	CheckoutID string   `json:"checkout_id" binding:"required"`
	TrxID      string   `json:"trx_id" binding:"required"`   // 要确认的交易ID
	ProofID    string   `json:"proof_id" binding:"required"` // 支付凭证ID/UTR
	TrxApp     string   `json:"trx_app"`                     // 支付平台应用
	ProofUrls  []string `json:"proof_urls,omitempty"`        // 支付凭证截图URLs
}

type CancelCheckoutRequest struct {
	Mid        string `json:"mid"`
	CheckoutID string `json:"checkout_id" binding:"required"`
}

type SubmitCheckoutRequest struct {
	Mid        string `json:"mid"`
	CheckoutID string `json:"checkout_id" binding:"required"`
	TrxMethod  string `json:"trx_method" binding:"required"`
	AccountNo  string `json:"account_no" binding:"required"`
}

type Checkout struct {
	CheckoutID  string       `json:"checkout_id"`
	Mid         string       `json:"mid"`
	ReqID       string       `json:"req_id"`
	Amount      string       `json:"amount"`
	Ccy         string       `json:"ccy"`
	Country     string       `json:"country"`
	TrxMethod   string       `json:"trx_method"`
	Status      string       `json:"status"`
	Transaction *Transaction `json:"transaction,omitempty"`
	NotifyURL   string       `json:"notify_url,omitempty"`
	ReturnURL   string       `json:"return_url,omitempty"`
	CreatedAt   int64        `json:"created_at"`
	UpdatedAt   int64        `json:"updated_at"`
	ExpiredAt   int64        `json:"expired_at"`
	CheckoutURL string       `json:"checkout_url,omitempty"`
	ErrorCode   string       `json:"error_code,omitempty"`
	ErrorMsg    string       `json:"error_msg,omitempty"`
	SubmitedAt  int64        `json:"submited_at,omitempty"`
	ConfirmedAt int64        `json:"confirmed_at,omitempty"`
	CanceledAt  int64        `json:"canceled_at,omitempty"`
	CompletedAt int64        `json:"completed_at,omitempty"`
	Token       string       `json:"token,omitempty"`
}
