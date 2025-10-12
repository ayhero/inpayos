package middleware

import (
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth API Key认证中间件
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := GetLanguage(c)
		appID := c.GetHeader("x-app-id")
		secret := c.GetHeader("x-secret")
		timestamp := c.GetHeader("x-timestamp")

		if timestamp == "" {
			timestamp = c.Query("x-timestamp")
		}
		if appID == "" {
			appID = c.Query("x-app-id")
		}
		if secret == "" {
			secret = c.Query("x-secret")
		}

		// 简单检查时间戳，防止重放攻击（允许5分钟误差）
		if timestamp == "" {
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidRequest, lang))
			c.Abort()
			return
		}
		if !IsTimestampValid(timestamp, 300) {
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidRequest, lang))
			c.Abort()
			return
		}

		// 验证 App ID 和 Secret
		// 简化验证逻辑，实际应该验证signature
		if appID == "" {
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidAPIKey, lang))
			c.Abort()
			return
		}

		app := models.GetAppByID(appID, secret)
		if app == nil || app.GetStatus() != protocol.StatusActive {
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidAPIKey, lang))
			c.Abort()
			return
		}
		merchant := models.GetMerchantByMID(app.Mid)
		if merchant == nil || merchant.GetStatus() != protocol.StatusActive {
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidAPIKey, lang))
			c.Abort()
			return
		}
		// 设置商户ID到上下文（从API Key获取）
		c.Set(UserIDKey, merchant.Mid)
		c.Set(UserKey, merchant)
		c.Next()
	}
}

// IsTimestampValid 验证时间戳是否有效
func IsTimestampValid(timestampStr string, allowedSkewSeconds int64) bool {
	if timestampStr == "" {
		return false
	}

	// 将字符串转换为整数时间戳
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return false
	}

	// 获取当前时间戳（秒）
	now := utils.TimeNowMilli()

	// 检查时间戳是否在允许的误差范围内
	diff := (now - timestamp) / 1000
	if diff < 0 {
		return false
	}

	return diff <= allowedSkewSeconds
}
