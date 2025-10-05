package models

// 移除未使用的import

// MerchantFeeConfig 费率配置表
type MerchantFeeConfig struct {
	ID      uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Mid     string `gorm:"column:mid;type:varchar(64);not null;index" json:"mid"`
	TrxType string `gorm:"column:trx_type;type:varchar(32);not null" json:"trx_type"` // receipt, payment, deposit, withdraw
	*FeeConfigValues
	CreatedAt int64 `gorm:"column:created_at;type:bigint;autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64 `gorm:"column:updated_at;type:bigint;autoUpdateTime:milli" json:"updated_at"`
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
func (MerchantFeeConfig) TableName() string {
	return "t_fee_configs"
}

// FeeConfigValues Getter Methods
// GetStatus returns the Status value
func (fcv *FeeConfigValues) GetStatus() string {
	if fcv.Status == nil {
		return ""
	}
	return *fcv.Status
}

// GetCountry returns the Country value
func (fcv *FeeConfigValues) GetCountry() string {
	if fcv.Country == nil {
		return ""
	}
	return *fcv.Country
}

// GetPaymentMethod returns the PaymentMethod value
func (fcv *FeeConfigValues) GetPaymentMethod() string {
	if fcv.PaymentMethod == nil {
		return ""
	}
	return *fcv.PaymentMethod
}

// GetCcy returns the Ccy value
func (fcv *FeeConfigValues) GetCcy() string {
	if fcv.Ccy == nil {
		return ""
	}
	return *fcv.Ccy
}

// GetPercent returns the Percent value
func (fcv *FeeConfigValues) GetPercent() string {
	if fcv.Percent == nil {
		return ""
	}
	return *fcv.Percent
}

// GetFixed returns the Fixed value
func (fcv *FeeConfigValues) GetFixed() string {
	if fcv.Fixed == nil {
		return ""
	}
	return *fcv.Fixed
}

// GetMinFee returns the MinFee value
func (fcv *FeeConfigValues) GetMinFee() string {
	if fcv.MinFee == nil {
		return ""
	}
	return *fcv.MinFee
}

// GetMaxFee returns the MaxFee value
func (fcv *FeeConfigValues) GetMaxFee() string {
	if fcv.MaxFee == nil {
		return ""
	}
	return *fcv.MaxFee
}

// FeeConfigValues Setter Methods (support method chaining)
// SetStatus sets the Status value
func (fcv *FeeConfigValues) SetStatus(value string) *FeeConfigValues {
	fcv.Status = &value
	return fcv
}

// SetCountry sets the Country value
func (fcv *FeeConfigValues) SetCountry(value string) *FeeConfigValues {
	fcv.Country = &value
	return fcv
}

// SetPaymentMethod sets the PaymentMethod value
func (fcv *FeeConfigValues) SetPaymentMethod(value string) *FeeConfigValues {
	fcv.PaymentMethod = &value
	return fcv
}

// SetCcy sets the Ccy value
func (fcv *FeeConfigValues) SetCcy(value string) *FeeConfigValues {
	fcv.Ccy = &value
	return fcv
}

// SetPercent sets the Percent value
func (fcv *FeeConfigValues) SetPercent(value string) *FeeConfigValues {
	fcv.Percent = &value
	return fcv
}

// SetFixed sets the Fixed value
func (fcv *FeeConfigValues) SetFixed(value string) *FeeConfigValues {
	fcv.Fixed = &value
	return fcv
}

// SetMinFee sets the MinFee value
func (fcv *FeeConfigValues) SetMinFee(value string) *FeeConfigValues {
	fcv.MinFee = &value
	return fcv
}

// SetMaxFee sets the MaxFee value
func (fcv *FeeConfigValues) SetMaxFee(value string) *FeeConfigValues {
	fcv.MaxFee = &value
	return fcv
}

// SetValues sets multiple FeeConfigValues fields at once
func (fc *MerchantFeeConfig) SetValues(values *FeeConfigValues) *MerchantFeeConfig {
	if values == nil {
		return fc
	}

	if fc.FeeConfigValues == nil {
		fc.FeeConfigValues = &FeeConfigValues{}
	}

	// Set all fields from the provided values
	if values.Status != nil {
		fc.FeeConfigValues.SetStatus(*values.Status)
	}
	if values.Country != nil {
		fc.FeeConfigValues.SetCountry(*values.Country)
	}
	if values.PaymentMethod != nil {
		fc.FeeConfigValues.SetPaymentMethod(*values.PaymentMethod)
	}
	if values.Ccy != nil {
		fc.FeeConfigValues.SetCcy(*values.Ccy)
	}
	if values.Percent != nil {
		fc.FeeConfigValues.SetPercent(*values.Percent)
	}
	if values.Fixed != nil {
		fc.FeeConfigValues.SetFixed(*values.Fixed)
	}
	if values.MinFee != nil {
		fc.FeeConfigValues.SetMinFee(*values.MinFee)
	}
	if values.MaxFee != nil {
		fc.FeeConfigValues.SetMaxFee(*values.MaxFee)
	}

	return fc
}
