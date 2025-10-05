package services

import (
	"encoding/json"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"sync"
	"time"

	"gorm.io/gorm"
)

type APIConfigService struct {
}

var (
	apiConfigService     *APIConfigService
	apiConfigServiceOnce sync.Once
)

func SetupAPIConfigService() {
	apiConfigServiceOnce.Do(func() {
		apiConfigService = &APIConfigService{}
	})
}

// GetAPIConfigService 获取APIConfig服务单例
func GetAPIConfigService() *APIConfigService {
	if apiConfigService == nil {
		SetupAPIConfigService()
	}
	return apiConfigService
}

// CreateAPIConfig 创建API配置
func (s *APIConfigService) CreateAPIConfig(req *protocol.CreateAPIConfigRequest) protocol.ErrorCode {
	isEnabled := true
	if req.IsEnabled != nil {
		isEnabled = *req.IsEnabled
	}

	// 将protocol字段映射到model字段
	// 解析ConfigValue为MapData
	var configData protocol.MapData
	if req.ConfigValue != "" {
		if err := json.Unmarshal([]byte(req.ConfigValue), &configData); err != nil {
			// 如果解析失败，创建简单的MapData
			configData = protocol.MapData{"value": req.ConfigValue}
		}
	}

	apiConfig := &models.APIConfig{
		Sid:     req.Mid,
		APIName: req.ConfigKey, // ConfigKey映射为APIName
		APIConfigValues: &models.APIConfigValues{
			IsEnabled:    &isEnabled,
			RateLimit:    &[]int{1000}[0], // 默认值
			DailyLimit:   &[]int{0}[0],    // 默认无限制
			MonthlyLimit: &[]int{0}[0],    // 默认无限制
			Config:       &configData,
			Description:  &req.Description,
		},
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}

	if err := models.ReadDB.Create(apiConfig).Error; err != nil {
		return protocol.InternalError
	}

	return protocol.Success
}

// UpdateAPIConfig 更新API配置
func (s *APIConfigService) UpdateAPIConfig(req *protocol.UpdateAPIConfigRequest) protocol.ErrorCode {
	updates := make(map[string]interface{})

	if req.ConfigValue != "" {
		updates["config"] = req.ConfigValue
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}

	updates["updated_at"] = time.Now().Unix()

	result := models.ReadDB.Model(&models.APIConfig{}).Where("id = ?", req.ID).Updates(updates)
	if result.Error != nil {
		return protocol.InternalError
	}
	if result.RowsAffected == 0 {
		return protocol.ConfigNotFound
	}

	return protocol.Success
}

// GetAPIConfig 获取API配置
func (s *APIConfigService) GetAPIConfig(req *protocol.GetAPIConfigRequest) (*protocol.APIConfigResponse, protocol.ErrorCode) {
	var apiConfig models.APIConfig
	query := models.ReadDB.Where("merchant_id = ? AND api_name = ?", req.Mid, req.ConfigKey)

	if err := query.Where("is_enabled = ?", true).First(&apiConfig).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, protocol.ConfigNotFound
		}
		return nil, protocol.InternalError
	}

	return s.convertToResponse(&apiConfig), protocol.Success
}

// QueryAPIConfigs 查询API配置列表
func (s *APIConfigService) QueryAPIConfigs(req *protocol.QueryAPIConfigRequest) (*protocol.APIConfigsResponse, protocol.ErrorCode) {
	query := models.ReadDB.Model(&models.APIConfig{})

	// 构建查询条件
	if req.Mid != "" {
		query = query.Where("merchant_id = ?", req.Mid)
	}
	if req.ConfigKey != "" {
		query = query.Where("api_name LIKE ?", "%"+req.ConfigKey+"%")
	}
	if req.IsEnabled != nil {
		query = query.Where("is_enabled = ?", *req.IsEnabled)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, protocol.InternalError
	}

	// 分页查询
	var apiConfigs []models.APIConfig
	offset := (req.Page - 1) * req.Size
	if err := query.Order("created_at DESC").
		Offset(offset).Limit(req.Size).Find(&apiConfigs).Error; err != nil {
		return nil, protocol.InternalError
	}

	// 转换响应
	configs := make([]protocol.APIConfigResponse, len(apiConfigs))
	for i, config := range apiConfigs {
		configs[i] = *s.convertToResponse(&config)
	}

	return &protocol.APIConfigsResponse{
		Configs: configs,
		Total:   total,
		Page:    req.Page,
		Size:    req.Size,
	}, protocol.Success
}

// DeleteAPIConfig 删除API配置
func (s *APIConfigService) DeleteAPIConfig(req *protocol.DeleteAPIConfigRequest) protocol.ErrorCode {
	result := models.ReadDB.Delete(&models.APIConfig{}, req.ID)
	if result.Error != nil {
		return protocol.InternalError
	}
	if result.RowsAffected == 0 {
		return protocol.ConfigNotFound
	}

	return protocol.Success
}

// BatchUpdateAPIConfig 批量更新API配置
func (s *APIConfigService) BatchUpdateAPIConfig(req *protocol.BatchUpdateAPIConfigRequest) protocol.ErrorCode {
	tx := models.ReadDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, config := range req.Configs {
		var apiConfig models.APIConfig
		err := tx.Where("merchant_id = ? AND api_name = ?",
			req.MerchantID, config.ConfigKey).First(&apiConfig).Error

		if err == gorm.ErrRecordNotFound {
			// 创建新配置
			isEnabled := true
			if config.IsEnabled != nil {
				isEnabled = *config.IsEnabled
			}

			newConfig := &models.APIConfig{
				Sid:     req.MerchantID,
				APIName: config.ConfigKey,
				APIConfigValues: &models.APIConfigValues{
					IsEnabled:    &isEnabled,
					RateLimit:    &[]int{1000}[0],
					DailyLimit:   &[]int{0}[0],
					MonthlyLimit: &[]int{0}[0],
					Description:  &config.Description,
				},
				CreatedAt: time.Now().UnixMilli(),
				UpdatedAt: time.Now().UnixMilli(),
			}

			if err := tx.Create(newConfig).Error; err != nil {
				tx.Rollback()
				return protocol.InternalError
			}
		} else if err != nil {
			tx.Rollback()
			return protocol.InternalError
		} else {
			// 更新现有配置
			updates := map[string]interface{}{
				"config":     config.ConfigValue,
				"updated_at": time.Now().Unix(),
			}

			if config.Description != "" {
				updates["description"] = config.Description
			}
			if config.IsEnabled != nil {
				updates["is_enabled"] = *config.IsEnabled
			}

			if err := tx.Model(&apiConfig).Updates(updates).Error; err != nil {
				tx.Rollback()
				return protocol.InternalError
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return protocol.InternalError
	}

	return protocol.Success
}

// BatchDeleteAPIConfig 批量删除API配置
func (s *APIConfigService) BatchDeleteAPIConfig(req *protocol.BatchDeleteAPIConfigRequest) protocol.ErrorCode {
	result := models.ReadDB.Delete(&models.APIConfig{}, req.IDs)
	if result.Error != nil {
		return protocol.InternalError
	}

	return protocol.Success
}

// ValidateAPIConfig 验证API配置
func (s *APIConfigService) ValidateAPIConfig(req *protocol.ValidateAPIConfigRequest) (*protocol.ValidateAPIConfigResponse, protocol.ErrorCode) {
	response := &protocol.ValidateAPIConfigResponse{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// 基本验证
	if req.ConfigKey == "" {
		response.IsValid = false
		response.Errors = append(response.Errors, "config_key is required")
	}
	if req.ConfigValue == "" {
		response.IsValid = false
		response.Errors = append(response.Errors, "config_value is required")
	}

	// 类型验证
	switch req.ConfigType {
	case "json":
		if req.ConfigValue != "" && !isValidJSON(req.ConfigValue) {
			response.IsValid = false
			response.Errors = append(response.Errors, "invalid JSON format")
		}
	case "number":
		if req.ConfigValue != "" && !isValidNumber(req.ConfigValue) {
			response.IsValid = false
			response.Errors = append(response.Errors, "invalid number format")
		}
	case "boolean":
		if req.ConfigValue != "" && !isValidBoolean(req.ConfigValue) {
			response.IsValid = false
			response.Errors = append(response.Errors, "invalid boolean format, expected true or false")
		}
	}

	// 检查是否已存在相同配置
	var existingConfig models.APIConfig
	err := models.ReadDB.Where("merchant_id = ? AND api_name = ?",
		req.Mid, req.ConfigKey).First(&existingConfig).Error
	if err == nil {
		response.Warnings = append(response.Warnings, "configuration with the same key already exists")
	}

	return response, protocol.Success
}

// convertToResponse 转换为响应格式
func (s *APIConfigService) convertToResponse(apiConfig *models.APIConfig) *protocol.APIConfigResponse {
	// 安全获取Config值
	configValue := ""
	if apiConfig.GetConfig() != nil {
		if configBytes, err := json.Marshal(apiConfig.GetConfig()); err == nil {
			configValue = string(configBytes)
		}
	}

	return &protocol.APIConfigResponse{
		ID:          apiConfig.ID,
		Mid:         apiConfig.Sid,
		ConfigKey:   apiConfig.APIName, // APIName映射为ConfigKey
		ConfigValue: configValue,       // Config映射为ConfigValue
		ConfigType:  "string",          // 默认类型
		Description: apiConfig.GetDescription(),
		IsEnabled:   apiConfig.GetIsEnabled(),
		Environment: "production", // 默认环境
		ValidFrom:   0,            // 默认值
		ValidTo:     0,            // 默认值
		Metadata:    "",           // 默认值
		CreatedAt:   apiConfig.CreatedAt,
		UpdatedAt:   apiConfig.UpdatedAt,
	}
}

// 辅助验证函数
func isValidJSON(value string) bool {
	return value != "" && (value[0] == '{' || value[0] == '[')
}

func isValidNumber(value string) bool {
	for _, char := range value {
		if char < '0' || char > '9' {
			if char != '.' && char != '-' {
				return false
			}
		}
	}
	return true
}

func isValidBoolean(value string) bool {
	return value == "true" || value == "false"
}
