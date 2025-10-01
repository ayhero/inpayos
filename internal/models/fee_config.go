package models

// 移除未使用的import

// FeeConfig 费率配置表
type FeeConfig struct {
	ID              uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	MerchantID      string `gorm:"column:merchant_id;type:varchar(64);not null;index" json:"merchant_id"`
	TransactionType string `gorm:"column:transaction_type;type:varchar(32);not null" json:"transaction_type"` // receipt, payment, deposit, withdraw
	Country         string `gorm:"column:country;type:varchar(3)" json:"country"`                             // ISO 3166-1 alpha-3
	PaymentMethod   string `gorm:"column:payment_method;type:varchar(32)" json:"payment_method"`              // bank_transfer, upi, card, etc.
	CoinID          string `gorm:"column:coin_id;type:varchar(34)" json:"coin_id"`                            // USD, EUR, INR, etc.
	FeePercent      string `gorm:"column:fee_percent;type:decimal(10,4);default:0" json:"fee_percent"`        // 百分比费率
	FeeFixed        string `gorm:"column:fee_fixed;type:decimal(20,8);default:0" json:"fee_fixed"`            // 固定费率
	MinFee          string `gorm:"column:min_fee;type:decimal(20,8);default:0" json:"min_fee"`                // 最小费用
	MaxFee          string `gorm:"column:max_fee;type:decimal(20,8);default:0" json:"max_fee"`                // 最大费用
	IsActive        bool   `gorm:"column:is_active;type:boolean;default:true" json:"is_active"`               // 是否启用
	CreatedAt       int64  `gorm:"column:created_at;type:bigint;autoCreateTime:milli" json:"created_at"`
	UpdatedAt       int64  `gorm:"column:updated_at;type:bigint;autoUpdateTime:milli" json:"updated_at"`
	DeletedAt       int64  `gorm:"column:deleted_at;type:bigint;index" json:"deleted_at,omitempty"`
}

// TableName 返回表名
func (FeeConfig) TableName() string {
	return "t_fee_configs"
}

// IsExpired 检查是否过期（如果有过期逻辑的话）
func (fc *FeeConfig) IsExpired() bool {
	// 费率配置通常不过期，这里可以根据业务需求扩展
	return false
}

// GetFeePercentDecimal 获取百分比费率的 decimal 值
func (fc *FeeConfig) GetFeePercentDecimal() float64 {
	// 这里可以使用 shopspring/decimal 进行精确计算
	// 暂时返回简单的字符串转换，实际使用时建议用 decimal
	return 0.0
}

// CalculateFee 计算费用
func (fc *FeeConfig) CalculateFee(amount string) (string, error) {
	// TODO: 实现费用计算逻辑
	// 1. 百分比费用 = amount * fee_percent / 100
	// 2. 总费用 = max(min_fee, min(max_fee, 百分比费用 + 固定费用))
	return "0", nil
}

// FeeConfigResponse 费率配置响应结构
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

// ToResponse 转换为响应结构
func (fc *FeeConfig) ToResponse() *FeeConfigResponse {
	return &FeeConfigResponse{
		ID:              fc.ID,
		MerchantID:      fc.MerchantID,
		TransactionType: fc.TransactionType,
		Country:         fc.Country,
		PaymentMethod:   fc.PaymentMethod,
		CoinID:          fc.CoinID,
		FeePercent:      fc.FeePercent,
		FeeFixed:        fc.FeeFixed,
		MinFee:          fc.MinFee,
		MaxFee:          fc.MaxFee,
		IsActive:        fc.IsActive,
		CreatedAt:       fc.CreatedAt,
		UpdatedAt:       fc.UpdatedAt,
	}
}

// GetFeeConfigByCondition 根据条件获取费率配置
func GetFeeConfigByCondition(merchantID, transactionType, country, paymentMethod, coinID string) (*FeeConfig, error) {
	var config FeeConfig
	query := DB.Where("merchant_id = ? AND transaction_type = ? AND is_active = ?", merchantID, transactionType, true)

	if country != "" {
		query = query.Where("country = ? OR country = ''", country)
	}
	if paymentMethod != "" {
		query = query.Where("payment_method = ? OR payment_method = ''", paymentMethod)
	}
	if coinID != "" {
		query = query.Where("coin_id = ? OR coin_id = ''", coinID)
	}

	// 按照优先级排序：具体条件优先于通用条件
	err := query.Order("country DESC, payment_method DESC, coin_id DESC").First(&config).Error
	if err != nil {
		return nil, err
	}

	return &config, nil
}
