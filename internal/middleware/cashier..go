package middleware

import (
	"inpayos/internal/config"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTidFromContext(c *gin.Context) string {
	if _v, exists := c.Get(UserIDKey); exists {
		if strTid, ok := _v.(string); ok {
			return strTid
		}
	}
	return ""
}

func GetCashierTeamFromContext(c *gin.Context) *models.CashierTeam {
	if _v, exists := c.Get(UserKey); exists {
		if v, ok := _v.(*models.CashierTeam); ok {
			return v
		}
	}
	return nil
}

// CashierTeamJWTAuth JWT认证中间件
func CashierTeamJWTAuth() gin.HandlerFunc {
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
		team := models.GetCashierTeamByTid(claims.UserID)
		// 临时设置用户信息
		c.Set(UserIDKey, team.ID)
		c.Set(UserKey, team)
		c.Next()
	}
}
