package protocol

import "github.com/shopspring/decimal"

// ChannelResult 支付渠道返回结果
type ChannelResult struct {
	Status           string           `json:"status"`             // 渠道支付状态
	ResCode          string           `json:"res_code"`           // 渠道返回码
	ResMsg           string           `json:"res_msg"`            // 渠道返回信息
	Reason           string           `json:"reason"`             // 渠道原因描述
	Link             string           `json:"link"`               // 渠道链接
	ChannelStatus    string           `json:"channel_status"`     // 渠道状态
	ChannelCode      string           `json:"channel_code"`       // 渠道代码
	ChannelTrxID     string           `json:"channel_trx_id"`     // 渠道交易ID
	DealID           string           `json:"deal_id"`            // 上游交易ID
	ChannelAccountID string           `json:"channel_account_id"` // 渠道账户ID
	ChannelFeeCcy    string           `json:"channel_fee_ccy"`    // 渠道费用币种
	ChannelFeeAmount *decimal.Decimal `json:"channel_fee_amount"` // 渠道费用金额
	CompletedAt      int64            `json:"completed_at"`       // 渠道完成时间
	CreatedAt        int64            `json:"created_at"`         // 渠道创建时间
	UpdatedAt        int64            `json:"updated_at"`         // 渠道更新时间
}

type ChannelPayinRequest struct {
	Mid            string           `json:"mid"`
	ReqID          string           `json:"req_id"`
	TrxID          string           `json:"trx_id"`
	Ccy            string           `json:"ccy"`
	Amount         *decimal.Decimal `json:"amount"`
	ChannelCode    string           `json:"channel_code"`
	ChannelAccount string           `json:"channel_account"`
	ChannelGroup   string           `json:"channel_group"`
	TrxMethod      string           `json:"trx_method"`
	TrxMode        string           `json:"trx_mode"`
	TrxApp         string           `json:"trx_app"`
	Did            string           `json:"did"`
	Pkg            string           `json:"pkg"`
	UserID         string           `json:"user_id"`
	ProductID      string           `json:"product_id"`
	UserIP         string           `json:"user_ip"`
	CountryCode    string           `json:"country_code"`
	AccountID      string           `json:"account_id"`
	AccountName    string           `json:"account_name"`
	BankCode       string           `json:"bank_code"`
	Email          string           `json:"email"`
	NotifyURL      string           `json:"notify_url"`
	ReturnURL      string           `json:"return_url"`
}

type ChannelRefundRequest struct {
	Mid            string           `json:"mid"`
	ReqID          string           `json:"req_id"`
	TrxID          string           `json:"trx_id"`
	ChannelTrxID   string           `json:"channel_trx_id"`
	OriTrxID       string           `json:"ori_trx_id"`
	Ccy            string           `json:"ccy"`
	Amount         *decimal.Decimal `json:"amount"`
	ChannelCode    string           `json:"channel_code"`
	ChannelAccount string           `json:"channel_account"`
	NotifyURL      string           `json:"notify_url"`
}

type ChannelPayoutRequest struct {
	Mid            string           `json:"mid"`
	ReqID          string           `json:"req_id"`
	TrxID          string           `json:"trx_id"`
	Ccy            string           `json:"ccy"`
	Amount         *decimal.Decimal `json:"amount"`
	ChannelCode    string           `json:"channel_code"`
	ChannelAccount string           `json:"channel_account"`
	ChannelGroup   string           `json:"channel_group"`
	TrxMethod      string           `json:"trx_method"`
	TrxMode        string           `json:"trx_mode"`
	TrxApp         string           `json:"trx_app"`
	Did            string           `json:"did"`
	Pkg            string           `json:"pkg"`
	UserID         string           `json:"user_id"`
	ProductID      string           `json:"product_id"`
	UserIP         string           `json:"user_ip"`
	CountryCode    string           `json:"country_code"`
	AccountID      string           `json:"account_id"`
	AccountName    string           `json:"account_name"`
	BankCode       string           `json:"bank_code"`
	Phone          string           `json:"phone"`
	Email          string           `json:"email"`
	NotifyURL      string           `json:"notify_url"`
	ReturnURL      string           `json:"return_url"`
}

type ChannelQueryQuest struct {
	Mid          string `json:"mid"`
	TrxType      string `json:"trx_type"`
	ReqID        string `json:"req_id"`
	TrxID        string `json:"trx_id"`
	ChannelTrxID string `json:"channel_trx_id"`
}

const (
	ChannnelBravo = "bravo"
)
