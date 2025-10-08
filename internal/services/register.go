package services

import (
	"errors"
	"inpayos/internal/log"
	"inpayos/internal/protocol"
	"regexp"
	"slices"
	"strings"
)

// SendRegisterSuccessEmail 发送注册成功邮件
func SendRegisterSuccessEmail(email string) error {
	msgService := GetMessageService()
	if msgService == nil {
		return errors.New("message service not initialized")
	}

	msg := &Message{
		Type: protocol.MsgTypeRegisterSuccess,
		Params: map[string]any{
			"to": email,
		},
	}
	if err := msgService.SendEmailMessage(msg); err != nil {
		log.Get().Errorf("Failed to send register success email: %v", err)
		return err
	}
	return nil
}

// 验证手机号格式
func IsValidPhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	return phoneRegex.MatchString(phone)
}

// 验证邮箱域名
func IsValidEmailDomain(email string) bool {
	// 可以在这里添加允许的邮箱域名列表
	allowedDomains := []string{}
	if len(allowedDomains) > 0 {
		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			return false
		}
		domain := parts[1]
		// 如果没有设置允许的域名，则默认允许所有
		if slices.Contains(allowedDomains, domain) {
			return true
		}
	}
	return true
}
