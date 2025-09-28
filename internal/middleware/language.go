package middleware

import (
	"inpayos/internal/i18n"

	"github.com/gin-gonic/gin"
)

// LanguageMiddleware 语言处理中间件
func LanguageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取语言设置
		acceptLang := c.GetHeader("Accept-Language")
		if acceptLang == "" {
			acceptLang = c.GetHeader("Content-Language")
		}
		if acceptLang == "" {
			acceptLang = c.Query("lang")
		}

		// 解析语言代码
		lang := i18n.GetLanguageFromAcceptLanguage(acceptLang)

		// 验证语言是否支持
		lang = i18n.GetValidLanguage(lang)

		// 存储到上下文中
		c.Set("language", lang)

		c.Next()
	}
}

// GetLanguage 从上下文获取语言代码
func GetLanguage(c *gin.Context) string {
	if lang, exists := c.Get("language"); exists {
		return lang.(string)
	}
	return i18n.DefaultLanguage
}
