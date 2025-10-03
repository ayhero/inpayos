package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// AuthorizationHeader 认证请求头
	AuthorizationHeader = "Authorization"
	// BearerPrefix Bearer前缀
	BearerPrefix = "Bearer "
	// TokenKey 请求参数中的token键名
	TokenKey = "token"
)

// GetContextData 从上下文中获取指定key的数据
func GetContextData[T any](c *gin.Context, key string) *T {
	value, exists := c.Get(key)
	if !exists {
		return nil
	}
	if data, ok := value.(*T); ok {
		return data
	}
	return nil
}

// 辅助函数

// 从请求获取Token
func GetTokenFromRequest(c *gin.Context) string {
	// 优先从请求头获取token
	authHeader := c.GetHeader(AuthorizationHeader)
	if authHeader != "" {
		// 检查 Authorization 格式是否为 Bearer <token>
		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString != authHeader {
			return tokenString
		}
	}

	// 从请求参数获取token
	return c.Query(TokenKey)
}
