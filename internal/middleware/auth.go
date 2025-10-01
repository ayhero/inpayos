package middleware

import (
	"inpayos/internal/protocol"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,x-token,x-app-id,x-signature,x-timestamp,accept-language")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// APIKeyAuth API Key认证中间件
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		appID := c.GetHeader("x-app-id")
		// signature := c.GetHeader("x-signature")
		// timestamp := c.GetHeader("x-timestamp")

		// 简化验证逻辑，实际应该验证signature
		if appID == "" {
			lang := GetLanguage(c)
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidAPIKey, lang))
			c.Abort()
			return
		}

		// TODO: 实际验证API Key和签名
		// 这里应该验证appID是否存在于merchant_secrets表中
		// 以及验证signature是否正确

		// 设置商户ID到上下文（从API Key获取）
		c.Set("mid", appID)
		c.Next()
	}
}

func GetMidFromContext(c *gin.Context) string {
	if mid, exists := c.Get("mid"); exists {
		if strMid, ok := mid.(string); ok {
			return strMid
		}
	}
	return ""
}

// WebhookSignatureVerification Webhook签名验证中间件
func WebhookSignatureVerification() gin.HandlerFunc {
	return func(c *gin.Context) {
		signature := c.GetHeader("x-signature")
		if signature == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing signature"})
			c.Abort()
			return
		}

		// TODO: 实现实际的签名验证逻辑
		c.Next()
	}
}

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			lang := GetLanguage(c)
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.AuthenticationFailed, lang))
			c.Abort()
			return
		}

		// TODO: 实现JWT验证逻辑
		// 解析token，验证有效性，提取用户信息

		// 临时设置用户信息
		c.Set("merchant_id", "test_merchant")
		c.Set("user_id", "test_user")
		c.Next()
	}
}

// PermissionCheck 权限检查中间件
func PermissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现权限检查逻辑
		// 检查用户是否有访问特定资源的权限
		c.Next()
	}
}
