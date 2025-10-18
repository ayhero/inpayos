package services

import (
	"inpayos/internal/log"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
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
func (req *MerchantRegisterRequest) ValidateRegisterRequest() protocol.ErrorCode {
	// 验证昵称长度
	if len(req.Nickname) > 32 {
		return protocol.InvalidParams
	}

	// 验证手机号格式
	if !IsValidPhone(req.Phone) {
		return protocol.InvalidParams
	}

	// 验证邮箱域名
	if !IsValidEmailDomain(req.Email) {
		return protocol.InvalidParams
	}

	return protocol.Success
}

// RegisterMerchant 注册商户
func RegisterMerchant(req *MerchantRegisterRequest) protocol.ErrorCode {
	// 参数验证
	if err := req.ValidateRegisterRequest(); err != protocol.Success {
		return err
	}

	// 验证码校验
	if !GetVerifyCodeService().VerifyCode(protocol.MsgChannelEmail, protocol.VerifyCodeTypeRegister, req.Email, req.VerifyCode) {
		return protocol.InvalidVerificationCode
	}

	// 检查邮箱是否已注册
	if models.CheckMerchantEmail(req.Email) {
		return protocol.InvalidParams
	}

	// 创建商户
	salt := utils.GenerateSalt()
	merchant := &models.Merchant{
		Mid: utils.GenerateMerchantID(),
		MerchantValues: &models.MerchantValues{
			Salt: &salt,
		},
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
		return protocol.InternalError
	}

	// 发送注册成功邮件
	//SendRegisterSuccessEmail(req.Email)
	return protocol.Success
}
