package protocol

// CreateFeeConfigRequest 创建费率配置请求
type CreateFeeConfigRequest struct {
	MerchantID      string `json:"merchant_id" binding:"required"`
	TransactionType string `json:"transaction_type" binding:"required"` // receipt, payment, deposit, withdraw
	Country         string `json:"country"`                             // ISO 3166-1 alpha-3
	PaymentMethod   string `json:"payment_method"`                      // bank_transfer, upi, card, etc.
	CoinID          string `json:"coin_id"`                             // USD, EUR, INR, etc.
	FeePercent      string `json:"fee_percent"`                         // 百分比费率
	FeeFixed        string `json:"fee_fixed"`                           // 固定费率
	MinFee          string `json:"min_fee"`                             // 最小费用
	MaxFee          string `json:"max_fee"`                             // 最大费用
	IsActive        *bool  `json:"is_active"`                           // 是否启用
	Description     string `json:"description"`                         // 描述
}

// UpdateFeeConfigRequest 更新费率配置请求
type UpdateFeeConfigRequest struct {
	ID              uint64 `json:"id" binding:"required"`
	TransactionType string `json:"transaction_type"`
	Country         string `json:"country"`
	PaymentMethod   string `json:"payment_method"`
	CoinID          string `json:"coin_id"`
	FeePercent      string `json:"fee_percent"`
	FeeFixed        string `json:"fee_fixed"`
	MinFee          string `json:"min_fee"`
	MaxFee          string `json:"max_fee"`
	IsActive        *bool  `json:"is_active"`
	Description     string `json:"description"`
}

// QueryFeeConfigRequest 查询费率配置请求
type QueryFeeConfigRequest struct {
	MerchantID      string `json:"merchant_id" form:"merchant_id"`
	TransactionType string `json:"transaction_type" form:"transaction_type"`
	Country         string `json:"country" form:"country"`
	PaymentMethod   string `json:"payment_method" form:"payment_method"`
	CoinID          string `json:"coin_id" form:"coin_id"`
	IsActive        *bool  `json:"is_active" form:"is_active"`
	Page            int    `json:"page" form:"page" binding:"min=1"`
	Size            int    `json:"size" form:"size" binding:"min=1,max=100"`
}

// FeeConfigResponse 费率配置响应
type FeeConfigResponse struct {
	ID              uint64 `json:"id"`
	MerchantID      string `json:"merchant_id"`
	TransactionType string `json:"transaction_type"`
	Country         string `json:"country"`
	PaymentMethod   string `json:"payment_method"`
	CoinID          string `json:"coin_id"`
	FeePercent      string `json:"fee_percent"`
	FeeFixed        string `json:"fee_fixed"`
	MinFee          string `json:"min_fee"`
	MaxFee          string `json:"max_fee"`
	IsActive        bool   `json:"is_active"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

// CalculateFeeRequest 计算费用请求
type CalculateFeeRequest struct {
	MerchantID      string `json:"merchant_id" binding:"required"`
	TransactionType string `json:"transaction_type" binding:"required"`
	Amount          string `json:"amount" binding:"required"`
	Currency        string `json:"currency" binding:"required"`
	Country         string `json:"country"`
	PaymentMethod   string `json:"payment_method"`
}

// CalculateFeeResponse 计算费用响应
type CalculateFeeResponse struct {
	Amount       string            `json:"amount"`        // 原始金额
	Fee          string            `json:"fee"`           // 计算出的费用
	TotalAmount  string            `json:"total_amount"`  // 总金额（原始金额+费用）
	FeeBreakdown map[string]string `json:"fee_breakdown"` // 费用明细
	Currency     string            `json:"currency"`      // 货币
	FeeConfigID  uint64            `json:"fee_config_id"` // 使用的费率配置ID
	CalculatedAt int64             `json:"calculated_at"` // 计算时间
}

// FeeConfigsResponse 费用配置列表响应
type FeeConfigsResponse struct {
	Configs []FeeConfigResponse `json:"configs"`
	Total   int64               `json:"total"`
	Page    int                 `json:"page"`
	Size    int                 `json:"size"`
}
