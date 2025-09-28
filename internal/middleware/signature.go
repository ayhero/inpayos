package middleware

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"inpayos/internal/protocol"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// 常量定义
const (
	HeaderAppID = "X-App-Id"
	HeaderSign  = "X-Sign"
	FieldAppID  = "app_id"
	FieldSign   = "sign"
)

// SignatureConfig 签名中间件配置
type SignatureConfig struct {
	SecretKey string
	Required  bool
	SkipPaths []string
}

// SignatureMiddleware 签名验证中间件
func SignatureMiddleware(config SignatureConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过指定路径
		for _, path := range config.SkipPaths {
			if strings.HasPrefix(c.Request.URL.Path, path) {
				c.Next()
				return
			}
		}

		// 获取app_id和签名
		appID := c.GetHeader(HeaderAppID)
		if appID == "" {
			appID = c.Query("app_id")
		}

		sign := c.GetHeader(HeaderSign)
		if sign == "" {
			sign = c.Query("sign")
		}

		// 如果不是必需的且没有签名，跳过验证
		if !config.Required && (appID == "" || sign == "") {
			c.Next()
			return
		}

		// 如果是必需的但缺少参数，返回错误
		if appID == "" || sign == "" {
			c.JSON(http.StatusUnauthorized, protocol.Response{
				Code:    401,
				Message: "Missing app_id or signature",
			})
			c.Abort()
			return
		}

		// 生成请求参数映射
		reqMap := make(map[string]interface{})

		// 添加查询参数
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				reqMap[key] = values[0]
			}
		}

		// 添加POST参数（如果是表单提交）
		if strings.Contains(c.GetHeader("Content-Type"), "application/x-www-form-urlencoded") {
			c.Request.ParseForm()
			for key, values := range c.Request.PostForm {
				if len(values) > 0 {
					reqMap[key] = values[0]
				}
			}
		}

		// 添加必要字段
		reqMap[FieldAppID] = appID
		reqMap[FieldSign] = sign

		// 验证签名
		if !verifySign(reqMap, config.SecretKey) {
			c.JSON(http.StatusUnauthorized, protocol.Response{
				Code:    401,
				Message: "Signature verification failed",
			})
			c.Abort()
			return
		}

		// 设置上下文
		c.Set("app_id", appID)
		c.Set("signature_verified", true)
		c.Next()
	}
}

// verifySign 验证签名
func verifySign(data map[string]interface{}, secretKey string) bool {
	// 获取所有键并排序
	keys := make([]string, 0)
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建签名字符串
	vals := make([]string, 0)
	for _, k := range keys {
		v := strings.TrimSpace(cast.ToString(data[k]))
		if k != FieldSign && len(v) > 0 {
			vals = append(vals, fmt.Sprintf("%s=%s", k, v))
		}
	}

	signStr := strings.Join(vals, "&")
	if secretKey != "" {
		signStr += "&key=" + secretKey
	}

	// 计算MD5签名
	expectedSign := fmt.Sprintf("%x", md5.Sum([]byte(signStr)))
	receivedSign := strings.TrimSpace(cast.ToString(data[FieldSign]))

	return strings.EqualFold(expectedSign, receivedSign)
}

// DefaultConfig 默认配置
func DefaultConfig() SignatureConfig {
	return SignatureConfig{
		Required: false,
		SkipPaths: []string{
			"/health",
			"/ping",
			"/metrics",
		},
	}
}
