package config

import (
	"net/http"
	"time"
)

// 环境常量
const (
	DevEnv     = "dev"
	ProdEnv    = "prod"
	DefaultEnv = DevEnv
)

const (
	DefaultReadTimeout  = 60
	DefaultWriteTimeout = 60
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
	OpenAPI  *ServiceConfig `mapstructure:"openapi"`
	Merchant *ServiceConfig `mapstructure:"merchant"`
	Admin    *ServiceConfig `mapstructure:"admin"`
}

// Validate 验证并设置服务配置默认值
func (c *ServerConfig) Validate() {
	if c.OpenAPI != nil {
		c.OpenAPI.Validate()
	}
	if c.Merchant != nil {
		c.Merchant.Validate()
	}
	if c.Admin != nil {
		c.Admin.Validate()
	}
}

type ServiceConfig struct {
	Prefix       string     `mapstructure:"prefix"`
	Name         string     `mapstructure:"name"`
	Port         string     `mapstructure:"port"`
	Version      string     `mapstructure:"version"`
	ReadTimeout  int        `mapstructure:"read_timeout"`  // 读取超时时间(秒)
	WriteTimeout int        `mapstructure:"write_timeout"` // 写入超时时间(秒)
	Jwt          *JWTConfig `mapstructure:"jwt"`           // JWT配置
}

func (s *ServiceConfig) Validate() {
	if s.Port == "" {
		panic("Service port cannot be empty")
	}
	if s.ReadTimeout <= 0 {
		s.ReadTimeout = DefaultReadTimeout
	}
	if s.WriteTimeout <= 0 {
		s.WriteTimeout = DefaultWriteTimeout
	}
	if s.Jwt == nil {
		s.Jwt = &JWTConfig{}
	}
	s.Jwt.Validate()
}

func (s *ServiceConfig) ToServer() *http.Server {
	return &http.Server{
		Addr:         ":" + s.Port,
		ReadTimeout:  DefaultReadTimeout * time.Second,
		WriteTimeout: DefaultWriteTimeout * time.Second,
	}
}
