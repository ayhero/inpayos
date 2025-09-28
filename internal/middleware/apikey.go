package middleware

import (
	"net/http"
	"strings"

	"inpayos/internal/models"
	"inpayos/internal/protocol"

	"github.com/gin-gonic/gin"
)

// APIKeyMiddleware API Key 认证中间件
// 支持以下方式传递认证信息：
// 1. Header: X-App-ID 和 X-Secret
// 2. Authorization: Bearer appid:secret
func APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 header 中获取 App ID 和 Secret
		appID := c.GetHeader("X-App-ID")
		secret := c.GetHeader("X-Secret")

		// 兼容 Authorization header 格式
		if appID == "" || secret == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				// 格式: Bearer appid:secret
				token := strings.TrimPrefix(authHeader, "Bearer ")
				parts := strings.SplitN(token, ":", 2)
				if len(parts) == 2 {
					appID = parts[0]
					secret = parts[1]
				}
			}
		}

		if appID == "" || secret == "" {
			lang := GetLanguage(c)
			c.JSON(http.StatusUnauthorized, protocol.NewErrorResultWithCode(protocol.InvalidAPIKey, lang))
			c.Abort()
			return
		}

		// 验证 App ID 和 Secret
		merchantSecret := models.GetByAppIDAndSecret(appID, secret)
		if merchantSecret == nil {
			lang := GetLanguage(c)
			c.JSON(http.StatusUnauthorized, protocol.NewErrorResultWithCode(protocol.InvalidAPIKey, lang))
			c.Abort()
			return
		}

		// 将商户信息存储到上下文中
		c.Set("mid", merchantSecret.Mid)
		c.Set("app_id", merchantSecret.AppID)
		c.Set("app_name", merchantSecret.AppName)
		c.Set("merchant_secret", merchantSecret)
		c.Set("api_key_verified", true)

		c.Next()
	}
}

// RequirePermission 权限验证中间件
func RequirePermission(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查API Key是否已验证
		verified, exists := c.Get("api_key_verified")
		if !exists || !verified.(bool) {
			lang := GetLanguage(c)
			c.JSON(http.StatusUnauthorized, protocol.NewErrorResultWithCode(protocol.AuthenticationFailed, lang))
			c.Abort()
			return
		}

		// 获取商户密钥信息
		merchantSecret, exists := c.Get("merchant_secret")
		if !exists {
			lang := GetLanguage(c)
			c.JSON(http.StatusForbidden, protocol.NewErrorResultWithCode(protocol.AuthenticationFailed, lang))
			c.Abort()
			return
		}

		// 检查是否有所需权限
		secret := merchantSecret.(*models.MerchantSecret)
		permissions := secret.GetPermissionList()

		hasPermission := false
		for _, permission := range permissions {
			if permission == requiredPermission || permission == "*" {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			lang := GetLanguage(c)
			c.JSON(http.StatusForbidden, protocol.NewErrorResultWithCode(protocol.InsufficientPermissions, lang))
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetMerchantSecret 从上下文获取商户密钥信息
func GetMerchantSecret(c *gin.Context) *models.MerchantSecret {
	merchantSecret, exists := c.Get("merchant_secret")
	if !exists {
		return nil
	}
	return merchantSecret.(*models.MerchantSecret)
}
