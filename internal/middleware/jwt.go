package middleware

import (
	"net/http"
	"strings"
	"time"

	"inpayos/internal/protocol"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	// AuthorizationHeader 认证请求头
	AuthorizationHeader = "Authorization"
	// BearerPrefix Bearer前缀
	BearerPrefix = "Bearer "
	// TokenKey 请求参数中的token键名
	TokenKey = "token"
)

// JWTClaims JWT载荷
type JWTClaims struct {
	UserID     string `json:"user_id"`
	UserType   string `json:"user_type"` // merchant, admin, cashier
	Email      string `json:"email"`
	MerchantID string `json:"merchant_id"` // 商户ID（对于admin可能为空）
	Role       string `json:"role"`        // 角色权限
	jwt.RegisteredClaims
}

// JWTMiddleware JWT认证中间件
func JWTMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ValidToken(c, []byte(jwtSecret))
		if token == nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
			c.Abort()
			return
		}

		// 提取用户信息
		if claims, ok := token.Claims.(*JWTClaims); ok {
			// 将用户信息存储到上下文中
			c.Set("user_id", claims.UserID)
			c.Set("user_type", claims.UserType)
			c.Set("email", claims.Email)
			c.Set("merchant_id", claims.MerchantID)
			c.Set("role", claims.Role)
			c.Set("jwt_claims", claims)
		}

		c.Next()
	}
}

// AdminPermissionMiddleware 管理员权限中间件
func AdminPermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
			c.Abort()
			return
		}

		// 检查是否为管理员
		if userType != "admin" {
			c.JSON(http.StatusForbidden, protocol.NewBusinessErrorResult("Admin access required"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// MerchantPermissionMiddleware 商户权限中间件
func MerchantPermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
			c.Abort()
			return
		}

		// 检查是否为商户
		if userType != "merchant" {
			c.JSON(http.StatusForbidden, protocol.NewBusinessErrorResult("Merchant access required"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// CashierPermissionMiddleware 出纳员权限中间件
func CashierPermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
			c.Abort()
			return
		}

		// 检查是否为出纳员或管理员
		if userType != "cashier" && userType != "admin" {
			c.JSON(http.StatusForbidden, protocol.NewBusinessErrorResult("Cashier access required"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireUserType 用户类型验证中间件
func RequireUserType(allowedTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
			c.Abort()
			return
		}

		userTypeStr := userType.(string)
		for _, allowedType := range allowedTypes {
			if userTypeStr == allowedType {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, protocol.NewBusinessErrorResult("Access denied for user type"))
		c.Abort()
	}
}

// RoleMiddleware 角色权限中间件
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
			c.Abort()
			return
		}

		roleStr := role.(string)
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, protocol.NewBusinessErrorResult("Insufficient permissions"))
		c.Abort()
	}
}

// ValidToken 验证Token
func ValidToken(c *gin.Context, jwtSecret []byte) *jwt.Token {
	// 从请求获取Token
	tokenString := GetTokenFromRequest(c)
	if tokenString == "" {
		return nil
	}

	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil
	}

	return token
}

// GetTokenFromRequest 从请求获取Token
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

// GenerateToken 生成JWT Token
func GenerateToken(userID, userType, email, merchantID, role string, expiresAt time.Time, jwtSecret string) (string, error) {
	claims := &JWTClaims{
		UserID:     userID,
		UserType:   userType,
		Email:      email,
		MerchantID: merchantID,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// RequestLogger 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// 计算请求处理时间
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// 简单日志记录（生产环境建议使用专业日志库）
		method := c.Request.Method
		path := c.Request.URL.Path
		ip := c.ClientIP()

		// 记录请求信息
		if statusCode >= 400 {
			// 错误请求使用更高级别日志
			println("ERROR:", method, path, statusCode, latency.String(), ip)
		} else {
			println("INFO:", method, path, statusCode, latency.String(), ip)
		}
	}
}

// CORSMiddleware CORS跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-App-Id, X-Sign")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理在处理请求过程中可能出现的错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(http.StatusInternalServerError, protocol.NewBusinessErrorResult("Internal server error: "+err.Error()))
		}
	}
}
