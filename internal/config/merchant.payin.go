package config

const (
	DefaultMerchantPayinExpiryMinutes = 30 // 默认支付订单过期时间，单位：分钟
)

type MerchantPayinConfig struct {
	ExpiryMinutes int `mapstructure:"expiry_minutes"` // 支付订单过期时间，单位：分钟
}

func (c *MerchantPayinConfig) Validate() {
	if c.ExpiryMinutes <= 0 {
		c.ExpiryMinutes = DefaultMerchantPayinExpiryMinutes
	}
}
