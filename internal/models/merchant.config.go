package models

import (
	"encoding/json"
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
)

// MerchantConfig 商户配置表
type MerchantConfig struct {
	ID         uint64           `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Mid        string           `json:"mid" gorm:"column:mid;type:varchar(64);index"`
	Type       string           `json:"type" gorm:"column:type;type:varchar(32);index"` // receipt, payment, deposit, withdraw等
	ConfigData protocol.MapData `json:"config_data" gorm:"column:config_data;type:text"`
	CreatedAt  int64            `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt  int64            `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

// TrxConfig 交易配置
type TrxConfig struct {
	MinAmount     map[string]string `json:"min_amount,omitempty"`     // 最小金额限制
	MaxAmount     map[string]string `json:"max_amount,omitempty"`     // 最大金额限制
	DailyLimit    map[string]string `json:"daily_limit,omitempty"`    // 每日限额
	MonthlyLimit  map[string]string `json:"monthly_limit,omitempty"`  // 每月限额
	FeeRate       map[string]string `json:"fee_rate,omitempty"`       // 费率配置
	FeeFixed      map[string]string `json:"fee_fixed,omitempty"`      // 固定费用
	Status        string            `json:"status,omitempty"`         // 状态：on, off
	AutoConfirm   string            `json:"auto_confirm,omitempty"`   // 自动确认：on, off
	NotifyURL     string            `json:"notify_url,omitempty"`     // 通知地址
	TimeoutMinute int               `json:"timeout_minute,omitempty"` // 超时时间（分钟）
}

// 表名
func (MerchantConfig) TableName() string {
	return "t_merchant_configs"
}

// GetConfig 获取配置对象
func (m *MerchantConfig) GetConfig() (*TrxConfig, error) {
	if m.ConfigData == nil || len(m.ConfigData) == 0 {
		return &TrxConfig{}, nil
	}

	var config TrxConfig
	// 将MapData转换为TrxConfig
	data, err := json.Marshal(m.ConfigData)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// SetConfig 设置配置对象
func (m *MerchantConfig) SetConfig(config *TrxConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	// 将JSON字符串转换为MapData
	var mapData protocol.MapData
	err = json.Unmarshal(data, &mapData)
	if err != nil {
		return err
	}

	m.ConfigData = mapData
	return nil
}

// Copy 复制配置
func (c *TrxConfig) Copy(source *TrxConfig) {
	if source == nil {
		return
	}

	if source.MinAmount != nil {
		if c.MinAmount == nil {
			c.MinAmount = make(map[string]string)
		}
		for k, v := range source.MinAmount {
			c.MinAmount[k] = v
		}
	}

	if source.MaxAmount != nil {
		if c.MaxAmount == nil {
			c.MaxAmount = make(map[string]string)
		}
		for k, v := range source.MaxAmount {
			c.MaxAmount[k] = v
		}
	}

	if source.DailyLimit != nil {
		if c.DailyLimit == nil {
			c.DailyLimit = make(map[string]string)
		}
		for k, v := range source.DailyLimit {
			c.DailyLimit[k] = v
		}
	}

	if source.MonthlyLimit != nil {
		if c.MonthlyLimit == nil {
			c.MonthlyLimit = make(map[string]string)
		}
		for k, v := range source.MonthlyLimit {
			c.MonthlyLimit[k] = v
		}
	}

	if source.FeeRate != nil {
		if c.FeeRate == nil {
			c.FeeRate = make(map[string]string)
		}
		for k, v := range source.FeeRate {
			c.FeeRate[k] = v
		}
	}

	if source.FeeFixed != nil {
		if c.FeeFixed == nil {
			c.FeeFixed = make(map[string]string)
		}
		for k, v := range source.FeeFixed {
			c.FeeFixed[k] = v
		}
	}

	if source.Status != "" {
		c.Status = source.Status
	}

	if source.AutoConfirm != "" {
		c.AutoConfirm = source.AutoConfirm
	}

	if source.NotifyURL != "" {
		c.NotifyURL = source.NotifyURL
	}

	if source.TimeoutMinute > 0 {
		c.TimeoutMinute = source.TimeoutMinute
	}
}

// GetMinAmount 获取最小金额
func (c *TrxConfig) GetMinAmount(currency string) decimal.Decimal {
	if c.MinAmount == nil {
		return decimal.Zero
	}

	if amount, exists := c.MinAmount[currency]; exists {
		if amt, err := decimal.NewFromString(amount); err == nil {
			return amt
		}
	}

	return decimal.Zero
}

// GetMaxAmount 获取最大金额
func (c *TrxConfig) GetMaxAmount(currency string) decimal.Decimal {
	if c.MaxAmount == nil {
		return decimal.Zero
	}

	if amount, exists := c.MaxAmount[currency]; exists {
		if amt, err := decimal.NewFromString(amount); err == nil {
			return amt
		}
	}

	return decimal.Zero
}

// GetFeeRate 获取费率
func (c *TrxConfig) GetFeeRate(currency string) decimal.Decimal {
	if c.FeeRate == nil {
		return decimal.Zero
	}

	if rate, exists := c.FeeRate[currency]; exists {
		if r, err := decimal.NewFromString(rate); err == nil {
			return r
		}
	}

	return decimal.Zero
}

// IsEnabled 检查是否启用
func (c *TrxConfig) IsEnabled() bool {
	return c.Status == "on"
}

// IsAutoConfirm 检查是否自动确认
func (c *TrxConfig) IsAutoConfirm() bool {
	return c.AutoConfirm == "on"
}

// 交易类型常量
const (
	TrxTypeReceipt  = "receipt"  // 代收
	TrxTypePayment  = "payment"  // 代付
	TrxTypeDeposit  = "deposit"  // 充值
	TrxTypeWithdraw = "withdraw" // 提现
	TrxTypeRefund   = "refund"   // 退款
)

// 配置状态常量
const (
	ConfigStatusOn  = "on"
	ConfigStatusOff = "off"
)

// 全局商户ID（用于全局配置）
const GlobalMerchantID = "*"
