package protocol

import "github.com/shopspring/decimal"

// CreateWebhookRequest 创建Webhook请求
type CreateWebhookRequest struct {
	UserID        string          `json:"user_id"`
	UserType      string          `json:"user_type"`
	TransactionID string          `json:"trx_id"`
	BillID        string          `json:"bill_id"`
	Type          string          `json:"type"`
	Status        string          `json:"status"`
	Amount        decimal.Decimal `json:"amount"`
	Fee           decimal.Decimal `json:"fee"`
	Currency      string          `json:"currency"`
	NotifyURL     string          `json:"notify_url"`
}

// WebhookResponse Webhook响应
type WebhookResponse struct {
	ID            uint64          `json:"id"`
	WebhookID     string          `json:"webhook_id"`
	UserID        string          `json:"user_id"`
	UserType      string          `json:"user_type"`
	TransactionID string          `json:"trx_id"`
	BillID        string          `json:"bill_id"`
	Type          string          `json:"type"`
	Status        string          `json:"status"`
	Amount        decimal.Decimal `json:"amount"`
	Fee           decimal.Decimal `json:"fee"`
	Currency      string          `json:"currency"`
	NotifyURL     string          `json:"notify_url"`
	NotifyStatus  string          `json:"notify_status"`
	NotifyTimes   int32           `json:"notify_times"`
	CreatedAt     int64           `json:"created_at"`
	UpdatedAt     int64           `json:"updated_at"`
}
