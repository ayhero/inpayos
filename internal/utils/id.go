package utils

import (
	"crypto/rand"
	"fmt"
	"time"
)

// GenerateID 生成带前缀的唯一ID
func GenerateID(prefix string) string {
	timestamp := time.Now().UnixMilli()

	// 生成4字节随机数
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)

	return fmt.Sprintf("%s_%d_%X", prefix, timestamp, randomBytes)
}

// GenerateFlowID 生成流水ID
func GenerateFlowID() string {
	return GenerateID("FLOW")
}

// GenerateTransactionID 生成交易ID
func GenerateTransactionID() string {
	return GenerateID("TX")
}

// GenerateAccountID 生成账户ID
func GenerateAccountID() string {
	return GenerateID("ACC")
}
