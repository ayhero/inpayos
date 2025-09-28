package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// 环境常量
const (
	DevEnv     = "dev"
	ProdEnv    = "prod"
	DefaultEnv = DevEnv
)

// 全局配置实例
var (
	config *Config
)

type Config struct {
	Debug    bool            `mapstructure:"debug"`
	Env      string          `mapstructure:"env"`
	Server   *ServerConfig   `mapstructure:"server"`
	Database *DatabaseConfig `mapstructure:"database"`
	Redis    *RedisConfig    `mapstructure:"redis"`
	Log      *LogConfig      `mapstructure:"log"`
	JWT      *JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	API   *ServiceConfig `mapstructure:"api"`
	Admin *ServiceConfig `mapstructure:"admin"`
}

type ServiceConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type DatabaseConfig struct {
	DSN             string        `mapstructure:"dsn"` // DSN连接字符串
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	Debug           bool          `mapstructure:"debug"`
}

type RedisConfig struct {
	DSN string `mapstructure:"dsn"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

type JWTConfig struct {
	Secret          string        `mapstructure:"secret"`
	ExpireDuration  time.Duration `mapstructure:"expire_duration"`
	RefreshDuration time.Duration `mapstructure:"refresh_duration"`
	Issuer          string        `mapstructure:"issuer"`
	SkipPaths       []string      `mapstructure:"skip_paths"`
}

func (c *Config) IsSandbox() bool {
	return c.Env != ProdEnv
}

// Validate 验证并设置日志配置默认值
func (l *LogConfig) Validate() {
	if l.Level == "" {
		l.Level = "info"
	}
}

// Validate 验证并设置JWT配置默认值
func (j *JWTConfig) Validate() {
	if j.Secret == "" {
		j.Secret = "inpayos-default-jwt-secret-key"
	}
	if j.ExpireDuration == 0 {
		j.ExpireDuration = 24 * time.Hour // 24小时
	}
	if j.RefreshDuration == 0 {
		j.RefreshDuration = 7 * 24 * time.Hour // 7天
	}
	if j.Issuer == "" {
		j.Issuer = "inpayos"
	}
	if len(j.SkipPaths) == 0 {
		j.SkipPaths = []string{"/health", "/ping", "/api/v1/webhook"}
	}
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

	// Validate other configs
	c.validateDatabaseConfig()
	c.validateRedisConfig()
}

func (s *ServerConfig) Validate() {
	if s.API == nil {
		s.API = &ServiceConfig{Port: "8080", Host: "0.0.0.0"}
	}
	if s.API.Port == "" {
		s.API.Port = "8080"
	}
	if s.API.Host == "" {
		s.API.Host = "0.0.0.0"
	}
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

// Get 获取配置单例
func Get() *Config {
	return config
}

// Set 设置配置单例
func Set(cfg *Config) {
	config = cfg
}
