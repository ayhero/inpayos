package protocol

import "github.com/shopspring/decimal"

// 交易相关请求/响应
type CreateTransactionRequest struct {
	UserID        string          `json:"user_id" binding:"required"`
	UserType      string          `json:"user_type" binding:"required"`
	BillID        string          `json:"bill_id" binding:"required"`
	Type          string          `json:"type" binding:"required,oneof=receipt payment refund deposit withdraw"`
	Amount        decimal.Decimal `json:"amount" binding:"required"`
	Fee           decimal.Decimal `json:"fee"`
	Currency      string          `json:"currency" binding:"required"`
	ChannelCode   string          `json:"channel_code"`
	PaymentMethod string          `json:"payment_method"`
	NotifyURL     string          `json:"notify_url"`
	ReturnURL     string          `json:"return_url"`
	SourceTxID    string          `json:"source_tx_id"` // 退款时的原交易ID
	Remark        string          `json:"remark"`
	ExpiredAt     int64           `json:"expired_at"`
}

type TransactionResponse struct {
	ID            uint64          `json:"id"`
	TransactionID string          `json:"transaction_id"`
	UserID        string          `json:"user_id"`
	UserType      string          `json:"user_type"`
	BillID        string          `json:"bill_id"`
	Type          string          `json:"type"`
	Status        string          `json:"status"`
	Amount        decimal.Decimal `json:"amount"`
	Fee           decimal.Decimal `json:"fee"`
	Currency      string          `json:"currency"`
	ChannelCode   string          `json:"channel_code"`
	PaymentMethod string          `json:"payment_method"`
	NotifyStatus  string          `json:"notify_status"`
	NotifyTimes   int             `json:"notify_times"`
	SourceTxID    string          `json:"source_tx_id"`
	Remark        string          `json:"remark"`
	ExpiredAt     int64           `json:"expired_at"`
	ConfirmedAt   int64           `json:"confirmed_at"`
	CreatedAt     int64           `json:"created_at"`
	UpdatedAt     int64           `json:"updated_at"`
}
