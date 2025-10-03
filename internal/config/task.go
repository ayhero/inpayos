package config

const (
	DefaultTaskScanInterval = 3  // 默认任务扫描间隔(秒)
	DefaultCacheDuration    = 10 // 默认缓存时间(分钟)
)

// TaskConfig 任务配置
type TaskConfig struct {
	ScanInterval int  `mapstructure:"scan_interval"`          // 扫描间隔(秒)
	Enabled      bool `mapstructure:"enabled" default:"true"` // 定时任务开关，默认开启
}

// Validate 验证并设置任务配置默认值
func (c *TaskConfig) Validate() {
	if c.ScanInterval <= 0 {
		c.ScanInterval = DefaultTaskScanInterval
	}
}
