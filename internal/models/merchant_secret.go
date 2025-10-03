package models

import (
	"encoding/json"
	"time"

	"inpayos/internal/protocol"
)

// MerchantSecret 商户密钥表
type MerchantSecret struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Mid         string `gorm:"column:mid;type:varchar(64);not null;index" json:"mid"`
	AppID       string `gorm:"column:app_id;type:varchar(64);not null;uniqueIndex" json:"app_id"`
	AppName     string `gorm:"column:app_name;type:varchar(128);not null" json:"app_name"`
	SecretKey   string `gorm:"column:secret_key;type:varchar(128);not null;uniqueIndex" json:"secret_key"`
	Permissions string `gorm:"column:permissions;type:text" json:"permissions"`                        // JSON 格式存储权限列表
	Status      string `gorm:"column:status;type:varchar(20);not null;default:'active'" json:"status"` // active, inactive, suspended
	ExpiresAt   int64  `gorm:"column:expires_at;type:bigint" json:"expires_at"`                        // 过期时间戳
	CreatedAt   int64  `gorm:"column:created_at;type:bigint;autoCreateTime:milli" json:"created_at"`
	UpdatedAt   int64  `gorm:"column:updated_at;type:bigint;autoUpdateTime:milli" json:"updated_at"`
	DeletedAt   int64  `gorm:"column:deleted_at;type:bigint;index" json:"deleted_at,omitempty"`
}

type MerchantSecretValues struct {
	Permissions string  `gorm:"column:permissions;type:text" json:"permissions"`                        // JSON 格式存储权限列表
	Status      *string `gorm:"column:status;type:varchar(20);not null;default:'active'" json:"status"` // active, inactive, suspended
	ExpiresAt   *int64  `gorm:"column:expires_at;type:bigint" json:"expires_at"`                        // 过期时间戳
}

// TableName 返回表名
func (MerchantSecret) TableName() string {
	return "t_merchant_secrets"
}

// IsActive 检查是否为活跃状态
func (ms *MerchantSecret) IsActive() bool {
	return ms.Status == "active"
}

// IsExpired 检查是否已过期
func (ms *MerchantSecret) IsExpired() bool {
	if ms.ExpiresAt == 0 {
		return false // 没有设置过期时间表示永不过期
	}
	return ms.ExpiresAt < time.Now().UnixMilli()
}

// GetPermissionList 获取权限列表
func (ms *MerchantSecret) GetPermissionList() []string {
	if ms.Permissions == "" {
		return []string{"*"} // 默认所有权限
	}

	var permissions []string
	err := json.Unmarshal([]byte(ms.Permissions), &permissions)
	if err != nil {
		return []string{"*"} // 解析失败时给予所有权限
	}

	return permissions
}

// SetPermissionList 设置权限列表
func (ms *MerchantSecret) SetPermissionList(permissions []string) error {
	data, err := json.Marshal(permissions)
	if err != nil {
		return err
	}
	ms.Permissions = string(data)
	return nil
}

// ToProtocol 转换为协议结构
func (ms *MerchantSecret) ToProtocol(includeSecret bool) *protocol.MerchantSecret {
	resp := &protocol.MerchantSecret{
		ID:        ms.ID,
		Mid:       ms.Mid,
		AppID:     ms.AppID,
		AppName:   ms.AppName,
		Status:    ms.Status,
		ExpiresAt: ms.ExpiresAt,
		CreatedAt: ms.CreatedAt,
		UpdatedAt: ms.UpdatedAt,
	}

	if includeSecret {
		resp.SecretKey = ms.SecretKey
	}

	return resp
}

// GetBySecretKey 根据 secret key 获取商户密钥信息
func GetBySecretKey(secretKey string) *MerchantSecret {
	if secretKey == "" {
		return nil
	}

	var secret MerchantSecret
	now := time.Now().UnixMilli()
	err := WriteDB.Where("secret_key = ? AND status = ? AND (expires_at = 0 OR expires_at > ?)", secretKey, "active", now).First(&secret).Error
	if err != nil {
		return nil
	}

	return &secret
}

// GetByAppIDAndSecret 根据 app_id 和 secret key 获取商户密钥信息
func GetByAppIDAndSecret(appID, secretKey string) *MerchantSecret {
	if appID == "" || secretKey == "" {
		return nil
	}

	var secret MerchantSecret
	now := time.Now().UnixMilli()
	err := WriteDB.Where("app_id = ? AND secret_key = ? AND status = ? AND (expires_at = 0 OR expires_at > ?)", appID, secretKey, "active", now).First(&secret).Error
	if err != nil {
		return nil
	}

	return &secret
}
