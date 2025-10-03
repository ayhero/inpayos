package services

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// ValidateG2FAKey validates the provided G2FA key.
func ValidateG2FAKey(g2faKey string) bool {
	// G2FA key validation logic
	return len(g2faKey) >= 16 && len(g2faKey) <= 32
}

// VerifyG2FACode verifies the G2FA code against the merchant's G2FA key.
func VerifyG2FACode(g2faKey, g2faCode string) bool {
	valid, err := totp.ValidateCustom(
		g2faCode,
		g2faKey,
		time.Now(),
		totp.ValidateOpts{
			Period:    30,
			Digits:    6,
			Algorithm: otp.AlgorithmSHA1,
		},
	)

	if err != nil {
		return false
	}

	return valid
}

// GenerateG2FAKey generates a new G2FA key.
func GenerateG2FAKey() string {
	// Generate a random 20-byte key
	key := make([]byte, 20)
	_, err := rand.Read(key)
	if err != nil {
		// If there's an error, return a fallback key (this should not happen in production)
		return ""
	}
	// Encode the key to base32 (which is the standard for TOTP)
	encoded := base32.StdEncoding.EncodeToString(key)

	// Remove padding and return
	return strings.TrimRight(encoded, "=")
}

// GenerateG2FAQRCode 生成 G2FA 的二维码内容
func GenerateG2FAQRCode(merchantID, g2faKey string) string {
	return fmt.Sprintf("otpauth://totp/inpay:%s?secret=%s&issuer=inpay", merchantID, g2faKey)
}
