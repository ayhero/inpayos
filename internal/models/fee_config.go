package models

// 移除未使用的import

// FeeConfig 费率配置表
type FeeConfig struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Mid       string `gorm:"column:mid;type:varchar(64);not null;index" json:"mid"`
	TrxType   string `gorm:"column:trx_type;type:varchar(32);not null" json:"trx_type"` // receipt, payment, deposit, withdraw
	CreatedAt int64  `gorm:"column:created_at;type:bigint;autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint;autoUpdateTime:milli" json:"updated_at"`
}
type FeeConfigValues struct {
	Status        *string `gorm:"column:status;type:varchar(32);not null;default:'active'" json:"status"` // active, inactive
	Country       *string `gorm:"column:country;type:varchar(3)" json:"country"`                          // ISO 3166-1 alpha-3
	PaymentMethod *string `gorm:"column:payment_method;type:varchar(32)" json:"payment_method"`           // bank_transfer, upi, card, etc.
	Ccy           *string `gorm:"column:ccy;type:varchar(34)" json:"ccy"`                                 // USD, EUR, INR, etc.
	Percent       *string `gorm:"column:percent;type:decimal(10,4);default:0" json:"percent"`             // 百分比费率
	Fixed         *string `gorm:"column:fixed;type:decimal(20,8);default:0" json:"fixed"`                 // 固定费率
	MinFee        *string `gorm:"column:min_fee;type:decimal(20,8);default:0" json:"min_fee"`             // 最小费用
	MaxFee        *string `gorm:"column:max_fee;type:decimal(20,8);default:0" json:"max_fee"`             // 最大费用
}

// TableName 返回表名
func (FeeConfig) TableName() string {
	return "t_fee_configs"
}
