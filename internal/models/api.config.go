package models

import (
	"inpayos/internal/protocol"
)

// APIConfig API配置表
type APIConfig struct {
	ID      uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Sid     string `gorm:"column:sid;type:varchar(64);not null;index" json:"sid"`
	SType   string `gorm:"column:s_type;type:varchar(32);not null;default:'api'" json:"s_type"` // 配置类型，默认为 'api'
	APIName string `gorm:"column:api_name;type:varchar(100);not null" json:"api_name"`          // API名称，如 "create_receipt", "query_balance"
	*APIConfigValues
	CreatedAt int64 `gorm:"column:created_at;type:bigint;autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64 `gorm:"column:updated_at;type:bigint;autoUpdateTime:milli" json:"updated_at"`
}

type APIConfigValues struct {
	IsEnabled    *bool             `gorm:"column:is_enabled;type:boolean;default:true" json:"is_enabled"`    // 是否启用
	RateLimit    *int              `gorm:"column:rate_limit;type:integer;default:1000" json:"rate_limit"`    // 每分钟请求限制
	DailyLimit   *int              `gorm:"column:daily_limit;type:integer;default:0" json:"daily_limit"`     // 每日请求限制，0表示无限制
	MonthlyLimit *int              `gorm:"column:monthly_limit;type:integer;default:0" json:"monthly_limit"` // 每月请求限制，0表示无限制
	IPWhitelist  *string           `gorm:"column:ip_whitelist;type:text" json:"ip_whitelist"`                // IP白名单，JSON数组格式
	Config       *protocol.MapData `gorm:"column:config;type:text" json:"config"`                            // API特定配置，JSON格式
	Description  *string           `gorm:"column:description;type:varchar(512)" json:"description"`          // 描述
}

// TableName 返回表名
func (APIConfig) TableName() string {
	return "t_api_configs"
}

// IsAPIEnabled 检查API是否启用
func (ac *APIConfig) IsAPIEnabled() bool {
	return ac.GetIsEnabled()
}

// CheckRateLimit 检查是否超过速率限制（这里只是结构，实际需要配合Redis等实现）
func (ac *APIConfig) CheckRateLimit(currentCount int) bool {
	rateLimit := ac.GetRateLimit()
	if rateLimit <= 0 {
		return true // 无限制
	}
	return currentCount < rateLimit
}

// CheckDailyLimit 检查是否超过每日限制
func (ac *APIConfig) CheckDailyLimit(currentCount int) bool {
	dailyLimit := ac.GetDailyLimit()
	if dailyLimit <= 0 {
		return true // 无限制
	}
	return currentCount < dailyLimit
}

// CheckMonthlyLimit 检查是否超过每月限制
func (ac *APIConfig) CheckMonthlyLimit(currentCount int) bool {
	monthlyLimit := ac.GetMonthlyLimit()
	if monthlyLimit <= 0 {
		return true // 无限制
	}
	return currentCount < monthlyLimit
}

// GetAPIConfigByMerchantAndAPI 根据商户ID和API名称获取配置
func GetAPIConfigByMerchantAndAPI(sid, apiName string) (*APIConfig, error) {
	var config APIConfig
	err := WriteDB.Where("sid = ? AND api_name = ? and s_type='merchant'", sid, apiName).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetAPIConfigsByMerchant 获取商户的所有API配置
func GetAPIConfigsByMerchant(sid string) ([]*APIConfig, error) {
	var configs []*APIConfig
	err := WriteDB.Where("sid = ? and s_type='merchant'", sid).Find(&configs).Error
	if err != nil {
		return nil, err
	}
	return configs, nil
}

// UpdateAPIConfig 更新API配置
func (ac *APIConfig) UpdateConfig(updates map[string]interface{}) error {
	updates["updated_at"] = getCurrentTimeMillis()
	return WriteDB.Model(ac).Updates(updates).Error
}

// EnableAPI 启用API
func (ac *APIConfig) EnableAPI() error {
	return ac.UpdateConfig(map[string]interface{}{
		"is_enabled": true,
	})
}

// DisableAPI 禁用API
func (ac *APIConfig) DisableAPI() error {
	return ac.UpdateConfig(map[string]interface{}{
		"is_enabled": false,
	})
}

// SetRateLimit 设置速率限制
func (ac *APIConfig) SetRateLimit(rateLimit int) error {
	return ac.UpdateConfig(map[string]interface{}{
		"rate_limit": rateLimit,
	})
}

// SetDailyLimit 设置每日限制
func (ac *APIConfig) SetDailyLimit(dailyLimit int) error {
	return ac.UpdateConfig(map[string]interface{}{
		"daily_limit": dailyLimit,
	})
}

// SetMonthlyLimit 设置每月限制
func (ac *APIConfig) SetMonthlyLimit(monthlyLimit int) error {
	return ac.UpdateConfig(map[string]interface{}{
		"monthly_limit": monthlyLimit,
	})
}

// SetIPWhitelist 设置IP白名单
func (ac *APIConfig) SetIPWhitelist(ipWhitelist string) error {
	return ac.UpdateConfig(map[string]interface{}{
		"ip_whitelist": ipWhitelist,
	})
}

// APIConfigValues Getter Methods
// GetIsEnabled returns the IsEnabled value
func (acv *APIConfigValues) GetIsEnabled() bool {
	if acv.IsEnabled == nil {
		return false
	}
	return *acv.IsEnabled
}

// GetRateLimit returns the RateLimit value
func (acv *APIConfigValues) GetRateLimit() int {
	if acv.RateLimit == nil {
		return 0
	}
	return *acv.RateLimit
}

// GetDailyLimit returns the DailyLimit value
func (acv *APIConfigValues) GetDailyLimit() int {
	if acv.DailyLimit == nil {
		return 0
	}
	return *acv.DailyLimit
}

// GetMonthlyLimit returns the MonthlyLimit value
func (acv *APIConfigValues) GetMonthlyLimit() int {
	if acv.MonthlyLimit == nil {
		return 0
	}
	return *acv.MonthlyLimit
}

// GetIPWhitelist returns the IPWhitelist value
func (acv *APIConfigValues) GetIPWhitelist() string {
	if acv.IPWhitelist == nil {
		return ""
	}
	return *acv.IPWhitelist
}

// GetConfig returns the Config value
func (acv *APIConfigValues) GetConfig() protocol.MapData {
	if acv.Config == nil {
		return nil
	}
	return *acv.Config
}

// GetDescription returns the Description value
func (acv *APIConfigValues) GetDescription() string {
	if acv.Description == nil {
		return ""
	}
	return *acv.Description
}

// APIConfigValues Setter Methods (support method chaining)
// SetIsEnabled sets the IsEnabled value
func (acv *APIConfigValues) SetIsEnabled(value bool) *APIConfigValues {
	acv.IsEnabled = &value
	return acv
}

// SetRateLimit sets the RateLimit value
func (acv *APIConfigValues) SetRateLimit(value int) *APIConfigValues {
	acv.RateLimit = &value
	return acv
}

// SetDailyLimit sets the DailyLimit value
func (acv *APIConfigValues) SetDailyLimit(value int) *APIConfigValues {
	acv.DailyLimit = &value
	return acv
}

// SetMonthlyLimit sets the MonthlyLimit value
func (acv *APIConfigValues) SetMonthlyLimit(value int) *APIConfigValues {
	acv.MonthlyLimit = &value
	return acv
}

// SetIPWhitelist sets the IPWhitelist value
func (acv *APIConfigValues) SetIPWhitelist(value string) *APIConfigValues {
	acv.IPWhitelist = &value
	return acv
}

// SetConfig sets the Config value
func (acv *APIConfigValues) SetConfig(value protocol.MapData) *APIConfigValues {
	acv.Config = &value
	return acv
}

// SetDescription sets the Description value
func (acv *APIConfigValues) SetDescription(value string) *APIConfigValues {
	acv.Description = &value
	return acv
}

// SetValues sets multiple APIConfigValues fields at once
func (ac *APIConfig) SetValues(values *APIConfigValues) *APIConfig {
	if values == nil {
		return ac
	}

	if ac.APIConfigValues == nil {
		ac.APIConfigValues = &APIConfigValues{}
	}

	// Set all fields from the provided values
	if values.IsEnabled != nil {
		ac.APIConfigValues.SetIsEnabled(*values.IsEnabled)
	}
	if values.RateLimit != nil {
		ac.APIConfigValues.SetRateLimit(*values.RateLimit)
	}
	if values.DailyLimit != nil {
		ac.APIConfigValues.SetDailyLimit(*values.DailyLimit)
	}
	if values.MonthlyLimit != nil {
		ac.APIConfigValues.SetMonthlyLimit(*values.MonthlyLimit)
	}
	if values.IPWhitelist != nil {
		ac.APIConfigValues.SetIPWhitelist(*values.IPWhitelist)
	}
	if values.Config != nil {
		ac.APIConfigValues.SetConfig(*values.Config)
	}
	if values.Description != nil {
		ac.APIConfigValues.SetDescription(*values.Description)
	}

	return ac
}
