package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// 全局配置实例
var (
	config *Config
)

type Config struct {
	Debug            bool                    `mapstructure:"debug"`
	Env              string                  `mapstructure:"env"`
	Server           *ServerConfig           `mapstructure:"server"`
	WriteDB          *DatabaseConfig         `mapstructure:"write_db"`
	ReadDB           *DatabaseConfig         `mapstructure:"read_db"`
	Redis            *RedisConfig            `mapstructure:"redis"`
	I18n             *I18nConfig             `mapstructure:"i18n"`
	Log              *LogConfig              `mapstructure:"log"`
	Email            *EmailConfig            `mapstructure:"email"`       // 邮件服务配置
	SMS              *SMSConfig              `mapstructure:"sms"`         // SMS服务配置
	VerifyCode       *VerifyCodeConfig       `mapstructure:"verify_code"` // 验证码配置
	Settle           *SettleConfig           `mapstructure:"settle"`      // 结算配置
	Task             *TaskConfig             `mapstructure:"task"`        // 任务调度配置
	MerchantPayin    *MerchantPayinConfig    `mapstructure:"payin"`       // 支付配置
	MerchantPayout   *MerchantPayoutConfig   `mapstructure:"payout"`      // 支付配置
	MerchantCheckout *MerchantCheckoutConfig `mapstructure:"checkout"`    // 结账配置
}

// Get 获取配置单例
func Get() *Config {
	return config
}

// Set 设置配置单例
func Set(cfg *Config) {
	config = cfg
}

func (c *Config) ValidateDB() {
	if c.WriteDB == nil {
		panic("WriteDB config is required")
	}
	if c.ReadDB == nil {
		panic("ReadDB config is required")
	}
	if c.Redis == nil {
		panic("Redis config is required")
	}
}

// Validate 验证并设置所有配置默认值
func (c *Config) Validate() {
	c.ValidateDB()
	c.Env = strings.ToLower(c.Env)
	if c.Env == "" || (c.Env != DevEnv && c.Env != ProdEnv) {
		c.Env = DefaultEnv
	}
	if c.Server != nil {
		c.Server.Validate()
	}
	if c.Log == nil {
		c.Log = &LogConfig{}
	}
	c.Log.Validate()
	if c.Settle != nil {
		c.Settle.Validate()
	}
	if c.I18n == nil {
		c.I18n = &I18nConfig{}
	}
	c.I18n.Validate()

	if c.Email != nil {
		c.Email.Validate()
	}
	if c.SMS != nil {
		c.SMS.Validate()
	}
	if c.VerifyCode == nil {
		c.VerifyCode = &VerifyCodeConfig{}
	}
	c.VerifyCode.Validate()
	if c.Task != nil {
		c.Task.Validate()
	}
	if c.MerchantPayin == nil {
		c.MerchantPayin = &MerchantPayinConfig{}
	}
	c.MerchantPayin.Validate()
	if c.MerchantPayout == nil {
		c.MerchantPayout = &MerchantPayoutConfig{}
	}
	c.MerchantPayout.Validate()
	if c.MerchantCheckout == nil {
		c.MerchantCheckout = &MerchantCheckoutConfig{}
	}
	c.MerchantCheckout.Validate()
}

// LoadConfig 加载配置
func LoadConfig() (err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(&config); err != nil {
		return
	}
	fmt.Println("Configuration loaded successfully")
	// 验证并设置默认值
	config.Validate()
	return
}
