package protocol

// CreateTransactionRequest 创建交易请求
type CreateTransactionRequest struct {
	ReqID         string `json:"req_id" binding:"required"`
	TrxType       string `json:"trx_type"` // 交易类型：payin, payout
	Ccy           string `json:"ccy" binding:"required"`
	Amount        string `json:"amount" binding:"required"`
	Country       string `json:"country"`
	PaymentMethod string `json:"payment_method" binding:"required"`
	NotifyURL     string `json:"notify_url"`
	ReturnURL     string `json:"return_url"`
}

// ConfirmTransactionRequest 确认交易请求（支付完成）
type ConfirmTransactionRequest struct {
	TrxID  string `json:"trx_id" ` // 交易ID
	Ccy    string `json:"ccy" binding:"required"`
	Amount string `json:"amount" binding:"required"`
}

// CancelTransactionRequest 取消/退款交易请求
type CancelTransactionRequest struct {
	TrxID   string `json:"trx_id" `  // 交易ID
	TrxType string `json:"trx_type"` // 交易类型
	ReqID   string `json:"req_id"`   // 请求ID
}

// QueryTransactionRequest 查询交易请求
type QueryTransactionRequest struct {
	TrxID   string `json:"trx_id" `  // 交易ID
	TrxType string `json:"trx_type"` // 交易类型
	ReqID   string `json:"req_id"`   // 请求ID
}

// Transaction 交易响应
type Transaction struct {
	TrxID         string `json:"trx_id"`
	Mid           string `json:"mid"`
	ReqID         string `json:"req_id"`
	Amount        string `json:"amount"`
	ActualAmount  string `json:"actual_amount,omitempty"`
	FeeAmount     string `json:"fee_amount,omitempty"`
	Ccy           string `json:"ccy"`
	Country       string `json:"country"`
	PaymentMethod string `json:"payment_method"`
	Status        string `json:"status"`
	NotifyURL     string `json:"notify_url,omitempty"`
	ReturnURL     string `json:"return_url,omitempty"`
	FailureReason string `json:"failure_reason,omitempty"`
	Metadata      string `json:"metadata,omitempty"`
	ExpiredAt     int64  `json:"expired_at"`
	CompletedAt   int64  `json:"completed_at"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}

// TransactionsResponse 交易列表响应
type TransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
	Total        int64         `json:"total"`
	Page         int           `json:"page"`
	Size         int           `json:"size"`
}
