package models

import (
	"encoding/json"
	"time"

	"inpayos/internal/protocol"
)

// MerchantSecret 商户密钥表
type MerchantSecret struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Mid       string `gorm:"column:mid;type:varchar(64);not null;index" json:"mid"`
	AppID     string `gorm:"column:app_id;type:varchar(64);not null;uniqueIndex" json:"app_id"`
	AppName   string `gorm:"column:app_name;type:varchar(128);not null" json:"app_name"`
	SecretKey string `gorm:"column:secret_key;type:varchar(128);not null;uniqueIndex" json:"secret_key"`
	*MerchantSecretValues
	CreatedAt int64 `gorm:"column:created_at;type:bigint;autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64 `gorm:"column:updated_at;type:bigint;autoUpdateTime:milli" json:"updated_at"`
	DeletedAt int64 `gorm:"column:deleted_at;type:bigint;index" json:"deleted_at,omitempty"`
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
	return ms.GetStatus() == protocol.StatusActive && !ms.IsExpired()
}

// IsExpired 检查是否已过期
func (ms *MerchantSecret) IsExpired() bool {
	expiresAt := ms.GetExpiresAt()
	if expiresAt == 0 {
		return false // 没有设置过期时间表示永不过期
	}
	return expiresAt < time.Now().UnixMilli()
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
		Status:    ms.GetStatus(),
		ExpiresAt: ms.GetExpiresAt(),
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

// MerchantSecretValues Getter Methods
// GetPermissions returns the Permissions value
func (msv *MerchantSecretValues) GetPermissions() string {
	return msv.Permissions
}

// GetStatus returns the Status value
func (msv *MerchantSecretValues) GetStatus() string {
	if msv.Status == nil {
		return ""
	}
	return *msv.Status
}

// GetExpiresAt returns the ExpiresAt value
func (msv *MerchantSecretValues) GetExpiresAt() int64 {
	if msv.ExpiresAt == nil {
		return 0
	}
	return *msv.ExpiresAt
}

// MerchantSecretValues Setter Methods (support method chaining)
// SetPermissions sets the Permissions value
func (msv *MerchantSecretValues) SetPermissions(value string) *MerchantSecretValues {
	msv.Permissions = value
	return msv
}

// SetStatus sets the Status value
func (msv *MerchantSecretValues) SetStatus(value string) *MerchantSecretValues {
	msv.Status = &value
	return msv
}

// SetExpiresAt sets the ExpiresAt value
func (msv *MerchantSecretValues) SetExpiresAt(value int64) *MerchantSecretValues {
	msv.ExpiresAt = &value
	return msv
}

// SetValues sets multiple MerchantSecretValues fields at once
func (ms *MerchantSecret) SetValues(values *MerchantSecretValues) *MerchantSecret {
	if values == nil {
		return ms
	}

	if ms.MerchantSecretValues == nil {
		ms.MerchantSecretValues = &MerchantSecretValues{}
	}

	// Set all fields from the provided values
	// Permissions is not a pointer, so we always set it
	ms.MerchantSecretValues.SetPermissions(values.Permissions)
	if values.Status != nil {
		ms.MerchantSecretValues.SetStatus(*values.Status)
	}
	if values.ExpiresAt != nil {
		ms.MerchantSecretValues.SetExpiresAt(*values.ExpiresAt)
	}

	return ms
}
