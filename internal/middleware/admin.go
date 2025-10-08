package middleware

import (
	"inpayos/internal/config"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	UserKey     = "user"
	UserIDKey   = "user_id"
	MerchantKey = "merchant"
)

func GetAdminIdFromContext(c *gin.Context) string {
	if _v, exists := c.Get(UserIDKey); exists {
		if strAdminId, ok := _v.(string); ok {
			return strAdminId
		}
	}
	return ""
}

func GetAdminFromContext(c *gin.Context) *models.Admin {
	if _v, exists := c.Get(UserKey); exists {
		if v, ok := _v.(*models.Admin); ok {
			return v
		}
	}
	return nil
}

// AdminJWTAuth JWT认证中间件
func AdminJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ValidToken(c, []byte(config.Get().Server.CashierAdmin.Jwt.Secret))
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
		team := models.GetAdminUserByID(claims.UserID)
		// 临时设置用户信息
		c.Set(UserIDKey, team.ID)
		c.Set(UserKey, team)
		c.Next()
	}
}
