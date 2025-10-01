package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/oklog/ulid/v2"
)

const (
	SALT_KEY = "jK9$mP2#nL5@qR8*"

	// ID前缀常量
	ID_PREFIX_USER        = "U"
	ID_PREFIX_ACCOUNT     = "AC"
	ID_PREFIX_FUNDFLOW    = "FF"
	ID_PREFIX_TRANSACTION = "TX"
	ID_PREFIX_PAYIN       = "PI"
	ID_PREFIX_PAYOUT      = "PO"
	ID_PREFIX_CHECKOUT    = "CKO"
	ID_PREFIX_WEBHOOK     = "WH"
	ID_PREFIX_SANDBOX     = "sandbox_"
)

func GenerateID() string {
	return ulid.Make().String()
}

// GenerateShortID 生成短ID
func GenerateShortID() string {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return strconv.FormatInt(time.Now().UnixMilli(), 36)
	}
	return hex.EncodeToString(b)
}
func GenerateWebhookID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_WEBHOOK, GenerateID())
}

func GeneratePayinID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_PAYIN, GenerateID())
}

func GeneratePayoutID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_PAYOUT, GenerateID())
}

// GenerateFlowID 生成流水ID
func GenerateFlowID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_FUNDFLOW, GenerateID())
}

// GenerateTransactionID 生成交易ID
func GenerateTransactionID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_TRANSACTION, GenerateID())
}

// GenerateAccountID 生成账户ID
func GenerateAccountID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_ACCOUNT, GenerateID())
}

func GenerateCheckoutID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_CHECKOUT, GenerateID())
}

// 用户相关ID生成
func GenerateUserID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_USER, GenerateID())
}

// GenerateSandboxChannelPaymentID 生成沙盒渠道支付ID
func GenerateSandboxChannelPaymentID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_SANDBOX, GenerateID())
}

// GenerateSalt 生成加密盐值
func GenerateSalt() string {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 16)
	}
	return hex.EncodeToString(salt)
}

// GenerateInviteCode 生成10位随机字母数字组合的邀请码
func GenerateInviteCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 10

	b := make([]byte, codeLength)
	for i := range b {
		randNum, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// 如果获取随机数失败，使用时间作为备用
			return fmt.Sprintf("%010d", time.Now().UnixNano())[:10]
		}
		b[i] = charset[randNum.Int64()]
	}
	return string(b)
}

// GenerateVerifyCode 生成6位数字验证码
func GenerateVerifyCode() string {
	return generateRandomNumber(6)
}

// generateRandomNumber 生成指定长度的随机数字字符串
func generateRandomNumber(length int) string {
	const charset = "0123456789"
	result := make([]byte, length)
	for i := range result {
		randNum, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// 如果获取随机数失败，使用简单的循环
			result[i] = charset[i%len(charset)]
		} else {
			result[i] = charset[randNum.Int64()]
		}
	}
	return string(result)
}
