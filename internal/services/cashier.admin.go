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

// ListCashiersByQuery 根据查询条件获取出纳员列表
func (s *CashierAdminService) ListCashiersByQuery(query *models.CashierTeamQuery) ([]*models.CashierTeam, int64, protocol.ErrorCode) {
	if query == nil {
		return nil, 0, protocol.InvalidParams
	}

	// 设置默认分页参数
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 20
	}
	if query.Size > 100 {
		query.Size = 100
	}

	// 构建查询条件
	db := query.BuildQuery()

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		log.Get().Errorf("ListCashiersByQuery: Count error: %v", err)
		return nil, 0, protocol.InternalError
	}

	// 分页查询
	var cashiers []*models.CashierTeam
	offset := (query.Page - 1) * query.Size
	if err := db.Offset(offset).Limit(query.Size).Order("created_at DESC").Find(&cashiers).Error; err != nil {
		log.Get().Errorf("ListCashiersByQuery: Find error: %v", err)
		return nil, 0, protocol.InternalError
	}

	return cashiers, total, protocol.Success
}

// GetCashierDetail 根据团队ID获取出纳员详情
func (s *CashierAdminService) GetCashierDetail(tid string) (*models.CashierTeam, protocol.ErrorCode) {
	if tid == "" {
		return nil, protocol.InvalidParams
	}

	// 根据Tid查询出纳员团队
	cashier := models.GetCashierTeamByTid(tid)
	if cashier == nil {
		return nil, protocol.CashierNotFound
	}

	return cashier, protocol.Success
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
	if models.CheckCashierTeamEmail(req.Email) {
		return protocol.InvalidParams
	}

	// 创建商户
	salt := utils.GenerateSalt()
	team := &models.CashierTeam{
		Tid: utils.GenerateCashierTeamID(),
		CashierTeamValues: &models.CashierTeamValues{
			Salt: &salt,
		},
	}

	// 设置商户基本信息
	team.SetEmail(req.Email).
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
	team.Encrypt()

	// 保存到数据库
	if err := models.WriteDB.Create(team).Error; err != nil {
		log.Get().Errorf("Failed to create merchant: %v", err)
		return protocol.InternalError
	}

	// 发送注册成功邮件
	//SendRegisterSuccessEmail(req.Email)
	return protocol.Success
}

// ResetCashierTeamPassword 重置商户密码
func ResetCashierTeamPassword(email string) (string, error) {
	// 检查商户是否存在
	merchant := models.GetMerchantByEmail(email)
	if merchant == nil {
		return "", errors.New("商户不存在")
	}

	// 生成新密码
	newPassword, err := utils.GenerateRandomPassword(12)
	if err != nil {
		log.Get().Error("Failed to generate password: ", err)
		return "", errors.New("生成新密码失败")
	}

	// 更新商户密码
	merchant.SetPassword(newPassword)
	merchant.Encrypt()

	if err := models.WriteDB.Model(merchant).Update("password", merchant.Password).Error; err != nil {
		log.Get().Error("Failed to update password: ", err)
		return "", errors.New("密码更新失败")
	}

	return newPassword, nil
}
