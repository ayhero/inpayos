package config

import "time"

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
