package protocol

// MerchantSecret 商户密钥协议结构
type MerchantSecret struct {
	ID        uint64 `json:"id"`
	Mid       string `json:"mid"`
	AppID     string `json:"app_id"`
	AppName   string `json:"app_name"`
	SecretKey string `json:"secret_key,omitempty"` // 可选返回，用于创建时
	Status    string `json:"status"`
	ExpiresAt int64  `json:"expires_at"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// Merchant 商户信息
type Merchant struct {
	Mid     string `json:"mid"`              // 商户ID
	Name    string `json:"name"`             // 商户名称
	Type    string `json:"type"`             // 商户类型
	Email   string `json:"email"`            // 商户邮箱
	Phone   string `json:"phone"`            // 商户电话
	Status  string `json:"status"`           // 商户状态
	Region  string `json:"region,omitempty"` // 商户区域
	Avatar  string `json:"avatar,omitempty"` // 商户头像
	HasG2FA bool   `json:"has_g2fa"`         // 是否启用二次验证
}
