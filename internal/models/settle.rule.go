package models

import (
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
)

type SettleRule struct {
	ID                int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	RuleID            string `json:"rule_id" gorm:"column:rule_id;index"` // RuleID 规则ID
	*SettleRuleValues `gorm:"embedded"`
	CreatedAt         int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt         int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

type SettleRuleValues struct {
	Status      int64            `json:"status" gorm:"column:status;index"`         // Status 状态
	MinAmount   *decimal.Decimal `json:"min_amount" gorm:"column:min_amount"`       // MinSettleAmount 最小结算金额
	MaxAmount   *decimal.Decimal `json:"max_amount" gorm:"column:max_amount"`       // MaxSettleAmount 最大结算金额
	MinFee      *decimal.Decimal `json:"min_fee" gorm:"column:min_fee"`             // MinSettleFee 最小手续费
	MaxFee      *decimal.Decimal `json:"max_fee" gorm:"column:max_fee"`             // MaxSettleFee 最大手续费
	MinRate     *decimal.Decimal `json:"min_rate" gorm:"column:min_rate"`           // MinSettleRate 最小费率
	MaxRate     *decimal.Decimal `json:"max_rate" gorm:"column:max_rate"`           // MaxSettleRate 最大费率
	MinUsdFee   *decimal.Decimal `json:"min_usd_fee" gorm:"column:min_usd_fee"`     // MinSettleUsdFee 最小美元手续费
	MaxUsdFee   *decimal.Decimal `json:"max_usd_fee" gorm:"column:max_usd_fee"`     // MaxSettleUsdFee 最大美元手续费
	MinUsdRate  *decimal.Decimal `json:"min_usd_rate" gorm:"column:min_usd_rate"`   // MinSettleUsdRate 最小美元费率
	MaxUsdRate  *decimal.Decimal `json:"max_usd_rate" gorm:"column:max_usd_rate"`   // MaxSettleUsdRate 最大美元费率
	Ccy         string           `json:"ccy" gorm:"column:ccy"`                     // CurrencyType 币种
	FixedFee    *decimal.Decimal `json:"fixed_fee" gorm:"column:fixed_fee"`         // FixedSettleFee 固定手续费
	Rate        *decimal.Decimal `json:"rate" gorm:"column:rate"`                   // SettleRate 费率
	FixedUsdFee *decimal.Decimal `json:"fixed_usd_fee" gorm:"column:fixed_usd_fee"` // FixedSettleUsdFee 固定美元手续费
	UsdRate     *decimal.Decimal `json:"usd_rate" gorm:"column:usd_rate"`           // SettleUsdRate 美元费率
}

type SettleRules []*SettleRule

func (t SettleRules) ToProtocol() []*protocol.SettleRule {
	rules := make([]*protocol.SettleRule, 0, len(t))
	for _, rule := range t {
		rules = append(rules, rule.ToProtocol())
	}
	return rules
}

func (t SettleRule) TableName() string {
	return "t_settle_rules"
}

func (t SettleRule) ToProtocol() *protocol.SettleRule {
	return &protocol.SettleRule{
		RuleID:      t.RuleID,
		MinAmount:   t.MinAmount,
		MaxAmount:   t.MaxAmount,
		MinFee:      t.MinFee,
		MaxFee:      t.MaxFee,
		MinRate:     t.MinRate,
		MaxRate:     t.MaxRate,
		MinUsdFee:   t.MinUsdFee,
		MaxUsdFee:   t.MaxUsdFee,
		MinUsdRate:  t.MinUsdRate,
		MaxUsdRate:  t.MaxUsdRate,
		Ccy:         t.Ccy,
		FixedFee:    t.FixedFee,
		Rate:        t.Rate,
		FixedUsdFee: t.FixedUsdFee,
		UsdRate:     t.UsdRate,
	}
}

// SetValues sets the SettleRuleValues fields, keeping existing values if new values are nil
func (t *SettleRule) SetValues(values *SettleRuleValues) *SettleRule {
	if values == nil {
		return t
	}

	if t.SettleRuleValues == nil {
		t.SettleRuleValues = &SettleRuleValues{}
	}

	t.Status = values.Status

	if values.MinAmount != nil {
		t.MinAmount = values.MinAmount
	}
	if values.MaxAmount != nil {
		t.MaxAmount = values.MaxAmount
	}
	if values.MinFee != nil {
		t.MinFee = values.MinFee
	}
	if values.MaxFee != nil {
		t.MaxFee = values.MaxFee
	}
	if values.MinRate != nil {
		t.MinRate = values.MinRate
	}
	if values.MaxRate != nil {
		t.MaxRate = values.MaxRate
	}
	if values.MinUsdFee != nil {
		t.MinUsdFee = values.MinUsdFee
	}
	if values.MaxUsdFee != nil {
		t.MaxUsdFee = values.MaxUsdFee
	}
	if values.MinUsdRate != nil {
		t.MinUsdRate = values.MinUsdRate
	}
	if values.MaxUsdRate != nil {
		t.MaxUsdRate = values.MaxUsdRate
	}
	if values.Ccy != "" {
		t.Ccy = values.Ccy
	}
	if values.FixedFee != nil {
		t.FixedFee = values.FixedFee
	}
	if values.Rate != nil {
		t.Rate = values.Rate
	}
	if values.FixedUsdFee != nil {
		t.FixedUsdFee = values.FixedUsdFee
	}
	if values.UsdRate != nil {
		t.UsdRate = values.UsdRate
	}

	return t
}

// Getters
func (t *SettleRuleValues) GetStatus() int64 {
	return t.Status
}
func (t *SettleRuleValues) GetMinAmount() decimal.Decimal {
	if t.MinAmount == nil {
		return decimal.Zero
	}
	return *t.MinAmount
}
func (t *SettleRuleValues) GetMaxAmount() decimal.Decimal {
	if t.MaxAmount == nil {
		return decimal.Zero
	}
	return *t.MaxAmount
}
func (t *SettleRuleValues) GetMinFee() decimal.Decimal {
	if t.MinFee == nil {
		return decimal.Zero
	}
	return *t.MinFee
}
func (t *SettleRuleValues) GetMaxFee() decimal.Decimal {
	if t.MaxFee == nil {
		return decimal.Zero
	}
	return *t.MaxFee
}
func (t *SettleRuleValues) GetMinRate() decimal.Decimal {
	if t.MinRate == nil {
		return decimal.Zero
	}
	return *t.MinRate
}
func (t *SettleRuleValues) GetMaxRate() decimal.Decimal {
	if t.MaxRate == nil {
		return decimal.Zero
	}
	return *t.MaxRate
}
func (t *SettleRuleValues) GetMinUsdFee() decimal.Decimal {
	if t.MinUsdFee == nil {
		return decimal.Zero
	}
	return *t.MinUsdFee
}
func (t *SettleRuleValues) GetMaxUsdFee() decimal.Decimal {
	if t.MaxUsdFee == nil {
		return decimal.Zero
	}
	return *t.MaxUsdFee
}
func (t *SettleRuleValues) GetMinUsdRate() decimal.Decimal {
	if t.MinUsdRate == nil {
		return decimal.Zero
	}
	return *t.MinUsdRate
}
func (t *SettleRuleValues) GetMaxUsdRate() decimal.Decimal {
	if t.MaxUsdRate == nil {
		return decimal.Zero
	}
	return *t.MaxUsdRate
}
func (t *SettleRuleValues) GetCcy() string {
	return t.Ccy
}
func (t *SettleRuleValues) GetFixedFee() decimal.Decimal {
	if t.FixedFee == nil {
		return decimal.Zero
	}
	return *t.FixedFee
}
func (t *SettleRuleValues) GetRate() decimal.Decimal {
	if t.Rate == nil {
		return decimal.Zero
	}
	return *t.Rate
}
func (t *SettleRuleValues) GetFixedUsdFee() decimal.Decimal {
	if t.FixedUsdFee == nil {
		return decimal.Zero
	}
	return *t.FixedUsdFee
}
func (t *SettleRuleValues) GetUsdRate() decimal.Decimal {
	if t.UsdRate == nil {
		return decimal.Zero
	}
	return *t.UsdRate
}

// Setters
func (t *SettleRuleValues) SetStatus(status int64) *SettleRuleValues {
	t.Status = status
	return t
}
func (t *SettleRuleValues) SetMinAmount(amount decimal.Decimal) *SettleRuleValues {
	t.MinAmount = &amount
	return t
}
func (t *SettleRuleValues) SetMaxAmount(amount decimal.Decimal) *SettleRuleValues {
	t.MaxAmount = &amount
	return t
}
func (t *SettleRuleValues) SetMinFee(fee decimal.Decimal) *SettleRuleValues {
	t.MinFee = &fee
	return t
}
func (t *SettleRuleValues) SetMaxFee(fee decimal.Decimal) *SettleRuleValues {
	t.MaxFee = &fee
	return t
}
func (t *SettleRuleValues) SetMinRate(rate decimal.Decimal) *SettleRuleValues {
	t.MinRate = &rate
	return t
}
func (t *SettleRuleValues) SetMaxRate(rate decimal.Decimal) *SettleRuleValues {
	t.MaxRate = &rate
	return t
}
func (t *SettleRuleValues) SetMinUsdFee(fee decimal.Decimal) *SettleRuleValues {
	t.MinUsdFee = &fee
	return t
}
func (t *SettleRuleValues) SetMaxUsdFee(fee decimal.Decimal) *SettleRuleValues {
	t.MaxUsdFee = &fee
	return t
}
func (t *SettleRuleValues) SetMinUsdRate(rate decimal.Decimal) *SettleRuleValues {
	t.MinUsdRate = &rate
	return t
}
func (t *SettleRuleValues) SetMaxUsdRate(rate decimal.Decimal) *SettleRuleValues {
	t.MaxUsdRate = &rate
	return t
}
func (t *SettleRuleValues) SetCcy(ccy string) *SettleRuleValues {
	t.Ccy = ccy
	return t
}
func (t *SettleRuleValues) SetFixedFee(fee decimal.Decimal) *SettleRuleValues {
	t.FixedFee = &fee
	return t
}
func (t *SettleRuleValues) SetRate(rate decimal.Decimal) *SettleRuleValues {
	t.Rate = &rate
	return t
}
func (t *SettleRuleValues) SetFixedUsdFee(fee decimal.Decimal) *SettleRuleValues {
	t.FixedUsdFee = &fee
	return t
}
func (t *SettleRuleValues) SetUsdRate(rate decimal.Decimal) *SettleRuleValues {
	t.UsdRate = &rate
	return t
}
