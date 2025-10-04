package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// 全局配置实例
var (
	config *Config
)

type Config struct {
	Debug      bool              `mapstructure:"debug"`
	Env        string            `mapstructure:"env"`
	Server     *ServerConfig     `mapstructure:"server"`
	Database   *DatabaseConfig   `mapstructure:"database"`
	WriteDB    *DatabaseConfig   `mapstructure:"write_db"`
	ReadDB     *DatabaseConfig   `mapstructure:"read_db"`
	Redis      *RedisConfig      `mapstructure:"redis"`
	Log        *LogConfig        `mapstructure:"log"`
	JWT        *JWTConfig        `mapstructure:"jwt"`
	Email      *EmailConfig      `mapstructure:"email"`       // 邮件服务配置
	SMS        *SMSConfig        `mapstructure:"sms"`         // SMS服务配置
	VerifyCode *VerifyCodeConfig `mapstructure:"verify_code"` // 验证码配置
	Settle     *SettleConfig     `mapstructure:"settle"`      // 结算配置
	Task       *TaskConfig       `mapstructure:"task"`        // 任务调度配置
}

// Get 获取配置单例
func Get() *Config {
	return config
}

// Set 设置配置单例
func Set(cfg *Config) {
	config = cfg
}

// Validate 验证并设置所有配置默认值
func (c *Config) Validate() {
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

	if c.JWT == nil {
		c.JWT = &JWTConfig{}
	}
	c.JWT.Validate()
	if c.Settle != nil {
		c.Settle.Validate()
	}

	// Validate other configs
	c.validateDatabaseConfig()
	c.validateRedisConfig()
}
func (c *Config) validateDatabaseConfig() {
	if c.Database.MaxIdleConns == 0 {
		c.Database.MaxIdleConns = 5
	}
	if c.Database.MaxOpenConns == 0 {
		c.Database.MaxOpenConns = 25
	}
	if c.Database.ConnMaxLifetime == 0 {
		c.Database.ConnMaxLifetime = 300 * time.Second
	}
}

func (c *Config) validateRedisConfig() {
	// Redis配置验证
	// DSN格式验证可以在连接时进行
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
