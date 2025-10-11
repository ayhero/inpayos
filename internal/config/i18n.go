package config

type I18nConfig struct {
	LocalesDir      string   `mapstructure:"locales_dir"`      // locales文件目录
	DefaultLanguage string   `mapstructure:"default_language"` // 默认语言
	SupportedLangs  []string `mapstructure:"supported_langs"`  // 支持的语言列表
}

// Validate 验证并设置I18n配置默认值
func (i *I18nConfig) Validate() {
	// 设置默认路径为容器中的绝对路径
	if i.LocalesDir == "" {
		i.LocalesDir = "/app/internal/locales" // 容器中的绝对路径
	}

	if i.DefaultLanguage == "" {
		i.DefaultLanguage = "en"
	}

	if len(i.SupportedLangs) == 0 {
		i.SupportedLangs = []string{"en", "zh"}
	}
}
