package services

import (
	"errors"
	"inpayos/internal/log"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
	"sync"
)

type CashierAdminService struct {
}

var (
	cashierAdminService     *CashierAdminService
	cashierAdminServiceOnce sync.Once
)

func SetupCashierAdminService() {
	cashierAdminServiceOnce.Do(func() {
		cashierAdminService = &CashierAdminService{}
	})
}

// GetCashierAdminService 获取CashierAdmin服务单例
func GetCashierAdminService() *CashierAdminService {
	if cashierAdminService == nil {
		SetupCashierAdminService()
	}
	return cashierAdminService
}

// ChangeCashierTeamPassword 修改商户密码
func ChangeCashierTeamPassword(email, newPassword string) error {
	// 根据邮箱获取商户信息
	team := models.GetCashierTeamByEmail(email)
	if team == nil {
		//log.Get().Errorf("ChangeCashierTeamPassword: CashierTeam not found for email: %s", email)
		return errors.New("Team 不存在")
	}
	salt := utils.GenerateSalt()
	team.Salt = &salt
	// 设置新密码和盐
	team.SetPassword(newPassword)

	// 加密密码
	team.Encrypt()

	// 更新密码和盐
	if err := models.GetDB().Model(&models.CashierTeam{}).Where("email = ?", email).Updates(map[string]interface{}{
		"password": team.GetPassword(),
		"salt":     team.Salt,
	}).Error; err != nil {
		//log.Get().Errorf("ChangeCashierTeamPassword: Update password error: %v", err)
		return errors.New("密码更新失败")
	}

	//log.Get().Infof("ChangeCashierTeamPassword: Password changed successfully for email: %s", email)
	return nil
}

type CashierTeamRegisterRequest struct {
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
func (req *CashierTeamRegisterRequest) ValidateRegisterRequest() protocol.ErrorCode {
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

// RegisterCashierTeam 注册商户
func RegisterCashierTeam(req *CashierTeamRegisterRequest) protocol.ErrorCode {
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
		Mid: utils.GenerateUserID(),
		MerchantValues: &models.MerchantValues{
			Salt: &salt,
		},
	}

	// 设置商户基本信息
	merchant.SetEmail(req.Email).
		SetPassword(req.Password).
		SetName(req.Nickname).
		SetType("cashier-team").
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
