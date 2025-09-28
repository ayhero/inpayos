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
