package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	"net/url"
	"sort"
	"strings"
	"time"
)

// SignatureType 签名类型
type SignatureType string

const (
	SignatureTypeMD5    SignatureType = "MD5"
	SignatureTypeSHA256 SignatureType = "SHA256"
	SignatureTypeHMAC   SignatureType = "HMAC"
)

// SignatureConfig 签名配置
type SignatureConfig struct {
	Type      SignatureType // 签名类型
	SecretKey string        // 密钥
	Timeout   int64         // 超时时间(秒)
}

// SignatureValidator 签名验证器
type SignatureValidator struct {
	config *SignatureConfig
}

// NewSignatureValidator 创建签名验证器
func NewSignatureValidator(config *SignatureConfig) *SignatureValidator {
	return &SignatureValidator{
		config: config,
	}
}

// GenerateSignature 生成签名
func (v *SignatureValidator) GenerateSignature(params map[string]string, timestamp int64) (string, error) {
	// 添加时间戳
	if timestamp > 0 {
		params["timestamp"] = fmt.Sprintf("%d", timestamp)
	}

	// 参数排序并拼接
	queryString := v.buildQueryString(params)

	// 根据不同类型生成签名
	switch v.config.Type {
	case SignatureTypeMD5:
		return v.generateMD5Signature(queryString), nil
	case SignatureTypeSHA256:
		return v.generateSHA256Signature(queryString), nil
	case SignatureTypeHMAC:
		return v.generateHMACSignature(queryString), nil
	default:
		return "", fmt.Errorf("unsupported signature type: %s", v.config.Type)
	}
}

// ValidateSignature 验证签名
func (v *SignatureValidator) ValidateSignature(params map[string]string, signature string, timestamp int64) error {
	// 检查时间戳是否过期
	if v.config.Timeout > 0 && timestamp > 0 {
		now := time.Now().Unix()
		if now-timestamp > v.config.Timeout {
			return fmt.Errorf("signature expired")
		}
	}

	// 生成期望的签名
	expectedSignature, err := v.GenerateSignature(params, timestamp)
	if err != nil {
		return fmt.Errorf("failed to generate signature: %v", err)
	}

	// 比较签名
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return fmt.Errorf("signature verification failed")
	}

	return nil
}

// buildQueryString 构建查询字符串
func (v *SignatureValidator) buildQueryString(params map[string]string) string {
	// 排序参数
	keys := make([]string, 0, len(params))
	for k := range params {
		// 跳过签名字段
		if k == "signature" || k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建查询字符串
	var parts []string
	for _, k := range keys {
		v := params[k]
		if v != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
		}
	}

	queryString := strings.Join(parts, "&")
	// 添加密钥
	if v.config.SecretKey != "" {
		queryString += "&key=" + v.config.SecretKey
	}

	return queryString
}

// generateMD5Signature 生成MD5签名
func (v *SignatureValidator) generateMD5Signature(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// generateSHA256Signature 生成SHA256签名
func (v *SignatureValidator) generateSHA256Signature(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// generateHMACSignature 生成HMAC签名
func (v *SignatureValidator) generateHMACSignature(data string) string {
	var h hash.Hash
	h = hmac.New(sha256.New, []byte(v.config.SecretKey))
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// ParseTimestamp 解析时间戳
func ParseTimestamp(timestampStr string) (int64, error) {
	if timestampStr == "" {
		return 0, nil
	}

	// 尝试解析不同格式的时间戳
	formats := []string{
		"1641024000",    // Unix timestamp (10位)
		"1641024000000", // Unix timestamp (13位毫秒)
	}

	for _, format := range formats {
		if len(timestampStr) == len(format) {
			var timestamp int64
			if _, err := fmt.Sscanf(timestampStr, "%d", &timestamp); err == nil {
				// 如果是13位时间戳，转换为10位
				if len(timestampStr) == 13 {
					timestamp = timestamp / 1000
				}
				return timestamp, nil
			}
		}
	}

	return 0, fmt.Errorf("invalid timestamp format: %s", timestampStr)
}

// ExtractParamsFromQuery 从查询字符串提取参数
func ExtractParamsFromQuery(queryString string) map[string]string {
	params := make(map[string]string)

	if queryString == "" {
		return params
	}

	pairs := strings.Split(queryString, "&")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key, _ := url.QueryUnescape(parts[0])
			value, _ := url.QueryUnescape(parts[1])
			params[key] = value
		}
	}

	return params
}
