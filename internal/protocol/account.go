package protocol

import "github.com/shopspring/decimal"

// 账户相关请求/响应
type CreateAccountRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	UserType string `json:"user_type" binding:"required,oneof=merchant cashier bank"`
	Currency string `json:"currency" binding:"required"`
}

type UpdateBalanceRequest struct {
	UserID        string          `json:"user_id" binding:"required"`
	UserType      string          `json:"user_type" binding:"required"`
	Currency      string          `json:"currency" binding:"required"`
	Operation     string          `json:"operation" binding:"required,oneof=add subtract freeze unfreeze margin release_margin"`
	Amount        decimal.Decimal `json:"amount" binding:"required"`
	TransactionID string          `json:"transaction_id"`
	BillID        string          `json:"bill_id"`
	BusinessType  string          `json:"business_type"`
	Description   string          `json:"description"`
}

type Balance struct {
	Balance          string `json:"balance"`
	AvailableBalance string `json:"available_balance"`
	FrozenBalance    string `json:"frozen_balance"`
	MarginBalance    string `json:"margin_balance"`
	ReserveBalance   string `json:"reserve_balance"`
	Currency         string `json:"currency"`
	UpdatedAt        int64  `json:"updated_at"`
}

// 流水相关响应
type FundFlowResponse struct {
	ID            uint64          `json:"id"`
	FlowID        string          `json:"flow_id"`
	UserID        string          `json:"user_id"`
	UserType      string          `json:"user_type"`
	AccountID     string          `json:"account_id"`
	TransactionID string          `json:"transaction_id"`
	BillID        string          `json:"bill_id"`
	FlowType      string          `json:"flow_type"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	BeforeBalance decimal.Decimal `json:"before_balance"`
	AfterBalance  decimal.Decimal `json:"after_balance"`
	BusinessType  string          `json:"business_type"`
	Description   string          `json:"description"`
	FlowAt        int64           `json:"flow_at"`
	CreatedAt     int64           `json:"created_at"`
}
