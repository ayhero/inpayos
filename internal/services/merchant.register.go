package services

import (
	"errors"
	"inpayos/internal/log"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
	"regexp"
	"slices"
	"strings"
)

type MerchantRegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`    // 邮箱
	Password    string `json:"password" binding:"required,min=8"` // 密码
	Nickname    string `json:"nickname" binding:"required,min=2"` // 昵称
	VerifyCode  string `json:"verify_code" binding:"required"`    // 验证码
	CompanyName string `json:"company_name" binding:"required"`   // 公司名称
	Phone       string `json:"phone" binding:"required"`          // 联系电话
	Region      string `json:"region"`                            // 地区
	RegIP       string `json:"reg_ip,omitempty"`                  // 注册IP
}

// ValidateRegisterRequest 验证注册请求
func (req *MerchantRegisterRequest) ValidateRegisterRequest() error {
	// 验证昵称长度
	if len(req.Nickname) > 32 {
		return errors.New("nickname too long, max length is 32")
	}

	// 验证手机号格式
	if !IsValidPhone(req.Phone) {
		return errors.New("invalid phone number format")
	}

	// 验证邮箱域名
	if !IsValidEmailDomain(req.Email) {
		return errors.New("email domain not allowed")
	}

	return nil
}

// RegisterMerchant 注册商户
func RegisterMerchant(req *MerchantRegisterRequest) error {
	// 参数验证
	if err := req.ValidateRegisterRequest(); err != nil {
		return err
	}

	// 验证码校验
	if !VerifyEmailCode(protocol.VerifyCodeTypeRegister, req.Email, req.VerifyCode) {
		return errors.New("invalid verify code")
	}

	// 检查邮箱是否已注册
	if models.CheckMerchantEmail(req.Email) {
		return errors.New("email already registered")
	}

	// 创建商户
	merchant := &models.Merchant{
		Mid:            utils.GenerateUserID(),
		Salt:           utils.GenerateSalt(),
		MerchantValues: &models.MerchantValues{},
	}

	// 设置商户基本信息
	merchant.SetEmail(req.Email).
		SetPassword(req.Password).
		SetName(req.Nickname).
		SetType("merchant").
		SetStatus(protocol.StatusActive).
		//SetRegIP(req.RegIP)
		// 如果有公司信息，设置公司相关字段
		SetName(req.CompanyName).
		SetPhone(req.Phone).
		SetRegion(req.Region).
		SetRegIP(req.RegIP)

	// 加密密码
	merchant.Encrypt()

	// 保存到数据库
	if err := models.WriteDB.Create(merchant).Error; err != nil {
		log.Get().Errorf("Failed to create merchant: %v", err)
		return errors.New("failed to create merchant")
	}

	// 发送注册成功邮件
	//SendRegisterSuccessEmail(req.Email)
	return nil
}

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
