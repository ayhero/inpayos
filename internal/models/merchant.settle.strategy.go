package models

import (
	"inpayos/internal/log"
	"inpayos/internal/protocol"
)

// 结算策略
type MerchantSettleStrategy struct {
	ID                            int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Code                          string `json:"code" gorm:"column:code;uniqueIndex:uk_strategy_code"`                            // Code 策略代码
	Mid                           string `json:"mid" gorm:"column:mid;type:varchar(64);index;uniqueIndex:uk_mid_settle_strategy"` // MerchantID 商户ID
	SettleCcy                     string `json:"settle_ccy" gorm:"column:settle_ccy;index;uniqueIndex:uk_mid_settle_strategy"`    // Currency 币种
	*MerchantSettleStrategyValues `gorm:"embedded"`
	CreatedAt                     int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt                     int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

type MerchantSettleStrategyValues struct {
	TrxType   string               `json:"trx_type" gorm:"column:trx_type;index;uniqueIndex:uk_mid_settle_strategy"`     // TrxType 交易类型
	TrxMode   string               `json:"trx_mode" gorm:"column:trx_mode;index;uniqueIndex:uk_mid_settle_strategy"`     // TrxMode 交易模式
	TrxMethod string               `json:"trx_method" gorm:"column:trx_method;index;uniqueIndex:uk_mid_settle_strategy"` // TrxMethod 交易方式
	TrxCcy    string               `json:"trx_ccy" gorm:"column:trx_ccy;index;uniqueIndex:uk_mid_settle_strategy"`       // TrxCcy 交易币种
	Country   string               `json:"country" gorm:"column:country;index;uniqueIndex:uk_mid_settle_strategy"`       // Country 国家
	Status    int64                `json:"status" gorm:"column:status;index"`                                            // Status 状态
	Rules     protocol.SettleRules `json:"rules" gorm:"column:rules;type:json;serializer:json"`                          // SettleRules 结算规则
}

func (t MerchantSettleStrategy) TableName() string {
	return "t_merchant_settle_strategies"
}

type MerchantSettleStrategies []*MerchantSettleStrategy

func (t MerchantSettleStrategies) Protocol() []*protocol.SettleStrategy {
	strategies := make([]*protocol.SettleStrategy, 0, len(t))
	for _, s := range t {
		strategies = append(strategies, s.Protocol())
	}
	return strategies
}

// Getters for MerchantSettleStrategy

// MerchantSettleStrategyValues Getter Methods
func (v *MerchantSettleStrategyValues) GetTrxType() string {
	return v.TrxType
}

func (v *MerchantSettleStrategyValues) GetTrxMode() string {
	return v.TrxMode
}

func (v *MerchantSettleStrategyValues) GetTrxMethod() string {
	return v.TrxMethod
}

func (v *MerchantSettleStrategyValues) GetTrxCcy() string {
	return v.TrxCcy
}

func (v *MerchantSettleStrategyValues) GetCountry() string {
	return v.Country
}

func (v *MerchantSettleStrategyValues) GetStatus() int64 {
	return v.Status
}

func (v *MerchantSettleStrategyValues) GetRules() protocol.SettleRules {
	return v.Rules
}

// MerchantSettleStrategyValues Setter Methods (support method chaining)
func (v *MerchantSettleStrategyValues) SetTrxType(trxType string) *MerchantSettleStrategyValues {
	v.TrxType = trxType
	return v
}

func (v *MerchantSettleStrategyValues) SetTrxMode(trxMode string) *MerchantSettleStrategyValues {
	v.TrxMode = trxMode
	return v
}

func (v *MerchantSettleStrategyValues) SetTrxMethod(trxMethod string) *MerchantSettleStrategyValues {
	v.TrxMethod = trxMethod
	return v
}

func (v *MerchantSettleStrategyValues) SetTrxCcy(trxCcy string) *MerchantSettleStrategyValues {
	v.TrxCcy = trxCcy
	return v
}

func (v *MerchantSettleStrategyValues) SetCountry(country string) *MerchantSettleStrategyValues {
	v.Country = country
	return v
}

func (v *MerchantSettleStrategyValues) SetStatus(status int64) *MerchantSettleStrategyValues {
	v.Status = status
	return v
}

func (v *MerchantSettleStrategyValues) SetRules(rules protocol.SettleRules) *MerchantSettleStrategyValues {
	v.Rules = rules
	return v
}

// SetValues sets multiple MerchantSettleStrategyValues fields at once
func (t *MerchantSettleStrategy) SetValues(values *MerchantSettleStrategyValues) *MerchantSettleStrategy {
	if values == nil {
		return t
	}

	if t.MerchantSettleStrategyValues == nil {
		t.MerchantSettleStrategyValues = &MerchantSettleStrategyValues{}
	}

	// Set all fields from the provided values
	if values.TrxType != "" {
		t.MerchantSettleStrategyValues.SetTrxType(values.TrxType)
	}
	if values.TrxMode != "" {
		t.MerchantSettleStrategyValues.SetTrxMode(values.TrxMode)
	}
	if values.TrxMethod != "" {
		t.MerchantSettleStrategyValues.SetTrxMethod(values.TrxMethod)
	}
	if values.TrxCcy != "" {
		t.MerchantSettleStrategyValues.SetTrxCcy(values.TrxCcy)
	}
	if values.Country != "" {
		t.MerchantSettleStrategyValues.SetCountry(values.Country)
	}
	if values.Status != 0 {
		t.MerchantSettleStrategyValues.SetStatus(values.Status)
	}
	if len(values.Rules) > 0 {
		t.MerchantSettleStrategyValues.SetRules(values.Rules)
	}

	return t
}

func GetSettleConfigByMid(mid, trx_type string, trx_time int64) *MerchantSettleStrategy {
	if mid == "" {
		return nil
	}
	// 查询数据库获取结算配置
	var cfg *MerchantSettleStrategy
	err := ReadDB.
		Where("mid = ? ", mid).
		Where("trx_type = ?", trx_type).
		Where("start_at <= ?", trx_time).
		Where("expired_at == 0 or expired_at >= ?", trx_time).
		First(&cfg).Error
	if err != nil {
		return nil
	}
	return cfg
}

func ListSettleStrategiesByMid(mid string) MerchantSettleStrategies {
	if mid == "" {
		return nil
	}
	var strategies []*MerchantSettleStrategy
	err := ReadDB.
		Where("mid = ?", mid).
		Find(&strategies).Error
	if err != nil {
		return nil
	}
	return strategies
}

func (t *MerchantSettleStrategy) Protocol() *protocol.SettleStrategy {
	v := &protocol.SettleStrategy{
		ID:        t.ID,
		Mid:       t.Mid,
		SettleCcy: t.SettleCcy,
		TrxType:   t.TrxType,
		TrxMode:   t.TrxMode,
		TrxMethod: t.TrxMethod,
		Country:   t.Country,
		TrxCcy:    t.TrxCcy,
		Status:    t.Status,
		Rules:     t.Rules,
	}
	return v
}

// GetSettleStrategiesByCodes 根据策略代码和交易类型批量查询结算策略
func GetSettleStrategiesByCodes(codes []string) []*MerchantSettleStrategy {
	if len(codes) == 0 {
		return nil
	}

	var strategies []*MerchantSettleStrategy
	err := ReadDB.
		Where("code IN ? AND status = ?", codes, protocol.StatusActive).
		Find(&strategies).Error
	if err != nil {
		log.Get().Errorf("GetSettleStrategiesByCodes failed, codes: %v, err: %v", codes, err)
		return nil
	}

	log.Get().Debugf("GetSettleStrategiesByCodes: found %d strategies for codes %v ", len(strategies), codes)
	return strategies
}
