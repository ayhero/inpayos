package middleware

import (
	"inpayos/internal/config"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMidFromContext(c *gin.Context) string {
	if _v, exists := c.Get(UserIDKey); exists {
		if strMid, ok := _v.(string); ok {
			return strMid
		}
	}
	return ""
}

func GetMerchantFromContext(c *gin.Context) *models.Merchant {
	if _v, exists := c.Get(UserKey); exists {
		if v, ok := _v.(*models.Merchant); ok {
			return v
		}
	}
	return nil
}

// MerchantJWTAuth JWT认证中间件
func MerchantJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ValidToken(c, []byte(config.Get().Server.Merchant.Jwt.Secret))
		if token == nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
			c.Abort()
			return
		}
		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
			c.Abort()
			return
		}
		if claims.UserID == "" {
			c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
			c.Abort()
			return
		}
		merchant := models.GetMerchantByMID(claims.UserID)
		// 临时设置用户信息
		c.Set(UserIDKey, merchant.Mid)
		c.Set(UserKey, merchant)
		c.Next()
	}
}
