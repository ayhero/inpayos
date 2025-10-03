package models

// APIConfig API配置表
type APIConfig struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	MerchantID   string `gorm:"column:merchant_id;type:varchar(64);not null;index" json:"merchant_id"`
	APIName      string `gorm:"column:api_name;type:varchar(100);not null" json:"api_name"`       // API名称，如 "create_receipt", "query_balance"
	IsEnabled    bool   `gorm:"column:is_enabled;type:boolean;default:true" json:"is_enabled"`    // 是否启用
	RateLimit    int    `gorm:"column:rate_limit;type:integer;default:1000" json:"rate_limit"`    // 每分钟请求限制
	DailyLimit   int    `gorm:"column:daily_limit;type:integer;default:0" json:"daily_limit"`     // 每日请求限制，0表示无限制
	MonthlyLimit int    `gorm:"column:monthly_limit;type:integer;default:0" json:"monthly_limit"` // 每月请求限制，0表示无限制
	IPWhitelist  string `gorm:"column:ip_whitelist;type:text" json:"ip_whitelist"`                // IP白名单，JSON数组格式
	Permissions  string `gorm:"column:permissions;type:text" json:"permissions"`                  // 权限列表，JSON数组格式
	Config       string `gorm:"column:config;type:text" json:"config"`                            // API特定配置，JSON格式
	Description  string `gorm:"column:description;type:varchar(512)" json:"description"`          // 描述
	CreatedAt    int64  `gorm:"column:created_at;type:bigint;autoCreateTime:milli" json:"created_at"`
	UpdatedAt    int64  `gorm:"column:updated_at;type:bigint;autoUpdateTime:milli" json:"updated_at"`
	DeletedAt    int64  `gorm:"column:deleted_at;type:bigint;index" json:"deleted_at,omitempty"`
}

// TableName 返回表名
func (APIConfig) TableName() string {
	return "t_api_configs"
}

// IsAPIEnabled 检查API是否启用
func (ac *APIConfig) IsAPIEnabled() bool {
	return ac.IsEnabled
}

// CheckRateLimit 检查是否超过速率限制（这里只是结构，实际需要配合Redis等实现）
func (ac *APIConfig) CheckRateLimit(currentCount int) bool {
	if ac.RateLimit <= 0 {
		return true // 无限制
	}
	return currentCount < ac.RateLimit
}

// CheckDailyLimit 检查是否超过每日限制
func (ac *APIConfig) CheckDailyLimit(currentCount int) bool {
	if ac.DailyLimit <= 0 {
		return true // 无限制
	}
	return currentCount < ac.DailyLimit
}

// CheckMonthlyLimit 检查是否超过每月限制
func (ac *APIConfig) CheckMonthlyLimit(currentCount int) bool {
	if ac.MonthlyLimit <= 0 {
		return true // 无限制
	}
	return currentCount < ac.MonthlyLimit
}

// APIConfigResponse API配置响应结构
type APIConfigResponse struct {
	ID           uint64 `json:"id"`
	MerchantID   string `json:"merchant_id"`
	APIName      string `json:"api_name"`
	IsEnabled    bool   `json:"is_enabled"`
	RateLimit    int    `json:"rate_limit"`
	DailyLimit   int    `json:"daily_limit"`
	MonthlyLimit int    `json:"monthly_limit"`
	IPWhitelist  string `json:"ip_whitelist,omitempty"`
	Permissions  string `json:"permissions,omitempty"`
	Config       string `json:"config,omitempty"`
	Description  string `json:"description,omitempty"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

// ToResponse 转换为响应结构
func (ac *APIConfig) ToResponse() *APIConfigResponse {
	return &APIConfigResponse{
		ID:           ac.ID,
		MerchantID:   ac.MerchantID,
		APIName:      ac.APIName,
		IsEnabled:    ac.IsEnabled,
		RateLimit:    ac.RateLimit,
		DailyLimit:   ac.DailyLimit,
		MonthlyLimit: ac.MonthlyLimit,
		IPWhitelist:  ac.IPWhitelist,
		Permissions:  ac.Permissions,
		Config:       ac.Config,
		Description:  ac.Description,
		CreatedAt:    ac.CreatedAt,
		UpdatedAt:    ac.UpdatedAt,
	}
}

// GetAPIConfigByMerchantAndAPI 根据商户ID和API名称获取配置
func GetAPIConfigByMerchantAndAPI(merchantID, apiName string) (*APIConfig, error) {
	var config APIConfig
	err := WriteDB.Where("merchant_id = ? AND api_name = ?", merchantID, apiName).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetAPIConfigsByMerchant 获取商户的所有API配置
func GetAPIConfigsByMerchant(merchantID string) ([]*APIConfig, error) {
	var configs []*APIConfig
	err := WriteDB.Where("merchant_id = ?", merchantID).Find(&configs).Error
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

// SetPermissions 设置权限列表
func (ac *APIConfig) SetPermissions(permissions string) error {
	return ac.UpdateConfig(map[string]interface{}{
		"permissions": permissions,
	})
}
