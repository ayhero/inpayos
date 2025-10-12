package services

import (
	"fmt"
	"inpayos/internal/config"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

// VerifyCodeService handles verification code operations
type VerifyCodeService struct {
	config     *config.VerifyCodeConfig
	msgService *MessageService
}

var (
	verifyCodeService *VerifyCodeService
	serviceLock       sync.Once
)

// SetupVerifyCodeService initializes verify code service with proper email service
func SetupVerifyCodeService() {
	cfg := config.Get()

	// Setup email service
	SetupEmailService()

	verifyCodeService = &VerifyCodeService{
		config:     cfg.VerifyCode,
		msgService: GetMessageService(),
	}
}

func GetVerifyCodeService() *VerifyCodeService {
	serviceLock.Do(func() {
		SetupVerifyCodeService()
	})
	return verifyCodeService
}

// GenerateCode generates a random verification code
func (s *VerifyCodeService) GenerateCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var result string
	for i := 0; i < s.config.Length; i++ {
		result += fmt.Sprintf("%d", r.Intn(10))
	}
	return result
}

// getInt64FromCache gets int64 value from cache
func (s *VerifyCodeService) getInt64FromCache(key string) (int64, error) {
	value, err := models.GetCache(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(value, 10, 64)
}

// setInt64ToCache sets int64 value to cache
func (s *VerifyCodeService) setInt64ToCache(key string, value int64, expiration time.Duration) error {
	return models.SetCache(key, strconv.FormatInt(value, 10), expiration)
}

// SendVerifyCode sends verification code via email or SMS
func (s *VerifyCodeService) SendVerifyCode(contactType, contact, purpose, language string) (protocol.ErrorCode, int) {
	// Validate contact type
	if contactType != protocol.MsgChannelEmail && contactType != protocol.MsgChannelSms {
		return protocol.InvalidVerificationMethod, 0
	}

	// Check send frequency
	lastTimeKey := fmt.Sprintf("%s_verify_code_%v_%s_time", contactType, purpose, contact)
	lastTime, err := s.getInt64FromCache(lastTimeKey)
	if err == nil && lastTime > 0 {
		if time.Now().Unix()-lastTime < int64(s.config.SendInterval) {
			remainingSeconds := s.config.SendInterval - int(time.Now().Unix()-lastTime)
			return protocol.VerificationCooldown, remainingSeconds
		}
	}
	isSandbox := config.Get().IsSandbox()
	if contactType == protocol.MsgChannelSms && strings.HasPrefix(contact, "+86") {
		isSandbox = true
	}
	// Generate verification code
	var code string
	if isSandbox {
		// For sandbox users, always use "123456"
		code = "123456"
	} else {
		// For real users, generate a random code
		code = s.GenerateCode()
	}

	// Store verification code in cache (convert minutes to seconds)
	codeKey := fmt.Sprintf("%s_verify_code_%v_%s", contactType, purpose, contact)
	if err := models.SetCache(codeKey, code, time.Duration(s.config.Expiration)*time.Minute); err != nil {
		return protocol.CacheError, 0
	}
	if !isSandbox {
		msg := &Message{
			Type:     protocol.MsgTypeVerifyCode,
			Channels: []string{contactType},
			Language: language,
			To:       contact,
			Params: map[string]any{
				"code":       code,
				"expiration": s.config.Expiration,
			},
		}
		if err := s.msgService.SendMessage(msg); err != nil {
			return protocol.VerificationCodeSendFailed, 0
		}
	}

	// Update send records
	s.setInt64ToCache(lastTimeKey, time.Now().Unix(), 24*time.Hour)

	return protocol.Success, 0
}

// SendEmailCode sends verification code via email
func (s *VerifyCodeService) SendEmailCode(email, purpose, language string) (protocol.ErrorCode, int) {
	return s.SendVerifyCode(protocol.MsgChannelEmail, email, purpose, language)
}

// SendSMSCode sends verification code via SMS
func (s *VerifyCodeService) SendSMSCode(phone, purpose, language string) (protocol.ErrorCode, int) {
	return s.SendVerifyCode(protocol.MsgChannelSms, phone, purpose, language)
}

// VerifyCode verifies a verification code for either email or SMS
func (s *VerifyCodeService) VerifyCode(contactType, purpose, contact, code string) bool {
	// Validate contact type
	if contactType != protocol.MsgChannelEmail && contactType != protocol.MsgChannelSms {
		return false
	}

	codeKey := fmt.Sprintf("%s_verify_code_%v_%v", contactType, purpose, contact)
	var storedCode string
	storedCode, err := models.GetCache(codeKey)
	if err != nil {
		return false
	}

	if storedCode != code {
		return false
	}

	// Delete the code after successful verification
	models.DelCache(codeKey)
	return true
}

// VerifyEmailCode verifies email verification code
func (s *VerifyCodeService) VerifyEmailCode(purpose, email, code string) bool {
	return s.VerifyCode(protocol.MsgChannelEmail, purpose, email, code)
}

// VerifySMSCode verifies SMS verification code
func (s *VerifyCodeService) VerifySMSCode(purpose, phone, code string) bool {
	return s.VerifyCode(protocol.MsgChannelSms, purpose, phone, code)
}
