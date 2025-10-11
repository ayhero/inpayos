package protocol

import "github.com/shopspring/decimal"

type OpenApiRequest struct {
	Mid          string `json:"mid"`
	TrxType      string `json:"trx_type"`
	ReqID        string `json:"req_id" binding:"required"`
	Ccy          string `json:"ccy" binding:"required"`
	Amount       string `json:"amount" binding:"required"`
	TrxMethod    string `json:"trx_method" binding:"required"`
	TrxMode      string `json:"trx_mode"`
	TrxApp       string `json:"trx_app"`
	Pkg          string `json:"pkg"`
	Did          string `json:"did"`
	ProductID    string `json:"product_id"`
	UserIP       string `json:"user_ip"`
	NotifyURL    string `json:"notify_url"`
	ReturnURL    string `json:"return_url"`
	ChannelCode  string `json:"channel_code"`
	ChannelGroup string `json:"channel_group"`
}

type MerchantPayinRequest struct {
	Mid          string `json:"mid"`
	ReqID        string `json:"req_id" binding:"required"`
	Ccy          string `json:"ccy" binding:"required"`
	Amount       string `json:"amount" binding:"required"`
	TrxMethod    string `json:"trx_method" binding:"required"`
	TrxMode      string `json:"trx_mode"`
	TrxApp       string `json:"trx_app"`
	Pkg          string `json:"pkg"`
	Did          string `json:"did"`
	ProductID    string `json:"product_id"`
	UserIP       string `json:"user_ip"`
	NotifyURL    string `json:"notify_url"`
	ReturnURL    string `json:"return_url"`
	ChannelCode  string `json:"channel_code"`
	ChannelGroup string `json:"channel_group"`
}

type MerchantPayoutRequest struct {
	Mid          string `json:"mid"`
	ReqID        string `json:"req_id" binding:"required"`
	Ccy          string `json:"ccy" binding:"required"`
	Amount       string `json:"amount" binding:"required"`
	TrxMethod    string `json:"trx_method" binding:"required"`
	TrxMode      string `json:"trx_mode"`
	TrxApp       string `json:"trx_app"`
	Pkg          string `json:"pkg"`
	Did          string `json:"did"`
	ProductID    string `json:"product_id"`
	UserIP       string `json:"user_ip"`
	NotifyURL    string `json:"notify_url"`
	ReturnURL    string `json:"return_url"`
	ChannelCode  string `json:"channel_code"`
	ChannelGroup string `json:"channel_group"`
}

type MerchantCancelRequest struct {
	Mid     string `json:"mid"`
	ReqID   string `json:"req_id" binding:"required"`
	TrxID   string `json:"trx_id" binding:"required"`
	TrxType string `json:"trx_type" binding:"required"`
}

type MerchantQueryRequest struct {
	Mid     string `json:"mid"`
	ReqID   string `json:"req_id" binding:"required"`
	TrxID   string `json:"trx_id" binding:"required"`
	TrxType string `json:"trx_type" binding:"required"`
}

// MerchantRefundRequest 退款请求
type MerchantRefundRequest struct {
	Mid       string           `json:"mid"`                  // 商户ID
	ReqID     string           `json:"req_id"`               // 退款请求ID
	OriTrxID  string           `json:"ori_trx_id"`           // 原支付交易ID
	Amount    *decimal.Decimal `json:"amount,omitempty"`     // 退款金额，不传则全额退款
	Reason    string           `json:"reason,omitempty"`     // 退款原因，可选
	NotifyURL string           `json:"notify_url,omitempty"` // 退款通知地址，可选，不传则使用原支付的通知地址
}
