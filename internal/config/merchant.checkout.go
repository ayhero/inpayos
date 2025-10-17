package config

const (
	DefaultCheckoutExpiryMinutes = 30 // 默认支付订单过期时间，单位：分钟
)

type MerchantCheckoutConfig struct {
	ExpiryMinutes int `mapstructure:"expiry_minutes"` // 支付订单过期时间，单位：分钟
}

func (c *MerchantCheckoutConfig) Validate() {
	if c.ExpiryMinutes <= 0 {
		c.ExpiryMinutes = DefaultMerchantPayoutExpiryMinutes
	}
}
