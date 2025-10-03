package utils

import (
	"inpayos/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// InitJWT 初始化 JWT 配置
func InitJWT() {
	cfg := config.Get().JWT
	jwtSecret = []byte(cfg.Secret)
}

// GetJWTSecret 获取 JWT 密钥
func GetJWTSecret() []byte {
	return jwtSecret
}

// GenerateJWT 生成 JWT
func GenerateMerchantTokenWithExpire(mid string, expire time.Duration) (string, error) {
	if expire <= 0 {
		expire = time.Duration(config.Get().JWT.ExpireDuration) * time.Second // 默认过期时间 24 小时
	}
	claims := jwt.MapClaims{
		"mid": mid,
		"exp": time.Now().Add(expire).Unix(), // 过期时间根据参数设置
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateJWT 生成 JWT
func GenerateTokenWithExpire(userID string, expire time.Duration) (string, error) {
	if expire <= 0 {
		expire = time.Duration(config.Get().JWT.ExpireDuration) * time.Second // 默认过期时间 24 小时
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expire).Unix(), // 过期时间根据参数设置
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateJWT 生成 JWT
func GenerateJWT(userID string) (string, error) {
	return GenerateTokenWithExpire(userID, 0)
}

// ParseJWT 解析 JWT
func ParseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}
