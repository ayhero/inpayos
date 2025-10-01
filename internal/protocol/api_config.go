package protocol

// CreateAPIConfigRequest 创建API配置请求
type CreateAPIConfigRequest struct {
	MerchantID  string `json:"merchant_id" binding:"required"`
	ConfigKey   string `json:"config_key" binding:"required"`
	ConfigValue string `json:"config_value" binding:"required"`
	ConfigType  string `json:"config_type"` // json, string, number, boolean
	Description string `json:"description"`
	IsEnabled   *bool  `json:"is_enabled"`  // 是否启用
	Environment string `json:"environment"` // test, production
	ValidFrom   int64  `json:"valid_from"`  // 生效时间
	ValidTo     int64  `json:"valid_to"`    // 失效时间
	Metadata    string `json:"metadata"`    // JSON格式的元数据
}

// UpdateAPIConfigRequest 更新API配置请求
type UpdateAPIConfigRequest struct {
	ID          uint64 `json:"id" binding:"required"`
	ConfigValue string `json:"config_value"`
	ConfigType  string `json:"config_type"`
	Description string `json:"description"`
	IsEnabled   *bool  `json:"is_enabled"`
	Environment string `json:"environment"`
	ValidFrom   int64  `json:"valid_from"`
	ValidTo     int64  `json:"valid_to"`
	Metadata    string `json:"metadata"`
}

// QueryAPIConfigRequest 查询API配置请求
type QueryAPIConfigRequest struct {
	MerchantID  string `json:"merchant_id" form:"merchant_id"`
	ConfigKey   string `json:"config_key" form:"config_key"`
	ConfigType  string `json:"config_type" form:"config_type"`
	Environment string `json:"environment" form:"environment"`
	IsEnabled   *bool  `json:"is_enabled" form:"is_enabled"`
	Page        int    `json:"page" form:"page" binding:"min=1"`
	Size        int    `json:"size" form:"size" binding:"min=1,max=100"`
}

// GetAPIConfigRequest 获取单个API配置请求
type GetAPIConfigRequest struct {
	MerchantID  string `json:"merchant_id" binding:"required"`
	ConfigKey   string `json:"config_key" binding:"required"`
	Environment string `json:"environment"`
}

// DeleteAPIConfigRequest 删除API配置请求
type DeleteAPIConfigRequest struct {
	ID uint64 `json:"id" binding:"required"`
}

// APIConfigResponse API配置响应
type APIConfigResponse struct {
	ID          uint64 `json:"id"`
	MerchantID  string `json:"merchant_id"`
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	ConfigType  string `json:"config_type"`
	Description string `json:"description"`
	IsEnabled   bool   `json:"is_enabled"`
	Environment string `json:"environment"`
	ValidFrom   int64  `json:"valid_from"`
	ValidTo     int64  `json:"valid_to"`
	Metadata    string `json:"metadata"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

// APIConfigsResponse API配置列表响应
type APIConfigsResponse struct {
	Configs []APIConfigResponse `json:"configs"`
	Total   int64               `json:"total"`
	Page    int                 `json:"page"`
	Size    int                 `json:"size"`
}

// BatchUpdateAPIConfigRequest 批量更新API配置请求
type BatchUpdateAPIConfigRequest struct {
	MerchantID string                `json:"merchant_id" binding:"required"`
	Configs    []APIConfigUpdateItem `json:"configs" binding:"required,min=1"`
}

// APIConfigUpdateItem 批量更新配置项
type APIConfigUpdateItem struct {
	ConfigKey   string `json:"config_key" binding:"required"`
	ConfigValue string `json:"config_value" binding:"required"`
	ConfigType  string `json:"config_type"`
	Description string `json:"description"`
	IsEnabled   *bool  `json:"is_enabled"`
	Environment string `json:"environment"`
	ValidFrom   int64  `json:"valid_from"`
	ValidTo     int64  `json:"valid_to"`
	Metadata    string `json:"metadata"`
}

// BatchDeleteAPIConfigRequest 批量删除API配置请求
type BatchDeleteAPIConfigRequest struct {
	IDs []uint64 `json:"ids" binding:"required,min=1"`
}

// ExportAPIConfigRequest 导出API配置请求
type ExportAPIConfigRequest struct {
	MerchantID  string `json:"merchant_id" form:"merchant_id"`
	Environment string `json:"environment" form:"environment"`
	Format      string `json:"format" form:"format"` // json, yaml, env
}

// ImportAPIConfigRequest 导入API配置请求
type ImportAPIConfigRequest struct {
	MerchantID  string `json:"merchant_id" binding:"required"`
	Environment string `json:"environment" binding:"required"`
	Format      string `json:"format" binding:"required"` // json, yaml, env
	Data        string `json:"data" binding:"required"`   // 配置数据
	Overwrite   bool   `json:"overwrite"`                 // 是否覆盖已存在的配置
}

// ValidateAPIConfigRequest 验证API配置请求
type ValidateAPIConfigRequest struct {
	MerchantID  string            `json:"merchant_id" binding:"required"`
	ConfigKey   string            `json:"config_key" binding:"required"`
	ConfigValue string            `json:"config_value" binding:"required"`
	ConfigType  string            `json:"config_type" binding:"required"`
	Rules       map[string]string `json:"rules"` // 验证规则
}

// ValidateAPIConfigResponse 验证API配置响应
type ValidateAPIConfigResponse struct {
	IsValid  bool     `json:"is_valid"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}
