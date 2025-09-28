package i18n

import (
	"log"
	"path/filepath"
)

// InitI18n 初始化国际化系统
func InitI18n(localesDir string) {
	translator := NewFileTranslator()

	// 如果localesDir为空，使用默认路径
	if localesDir == "" {
		localesDir = "internal/locales"
	}

	// 转换为绝对路径
	absPath, err := filepath.Abs(localesDir)
	if err != nil {
		log.Printf("Warning: Could not resolve locales directory path: %v", err)
		absPath = localesDir
	}

	// 加载翻译文件
	if err := translator.LoadTranslations(absPath); err != nil {
		log.Printf("Warning: Failed to load translations from %s: %v", absPath, err)
		// 即使加载失败也要设置翻译器，这样至少能返回原始错误码
	}

	// 设置全局翻译器
	SetGlobalTranslator(translator)

	log.Printf("I18n system initialized with locales from: %s", absPath)
}
