package services

import (
	"fmt"
	"strconv"
	"time"

	"inpayos/internal/models"
	"inpayos/internal/protocol"

	"gorm.io/gorm"
)

// CashierService 收银员服务
type CashierService struct {
	db *gorm.DB
}

// NewCashierService 创建收银员服务
func NewCashierService() *CashierService {
	return &CashierService{
		db: models.GetDB(),
	}
}

// CreateCashier 创建收银员
func (s *CashierService) CreateCashier(req *protocol.CreateCashierRequest) (*models.Cashier, error) {
	now := getCurrentTimeMillis()
	cashier := &models.Cashier{
		CashierID:     generateCashierID(),
		AccountID:     req.AccountID,
		CashierValues: &models.CashierValues{},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// 设置基础信息
	cashier.CashierValues.SetType(req.Type).
		SetCardNumber(req.CardNumber).
		SetHolderName(req.HolderName).
		SetBankCode(req.BankCode).
		SetBankName(req.BankName).
		SetStatus("active")

	// 设置可选信息
	if req.HolderPhone != "" {
		cashier.CashierValues.SetHolderPhone(req.HolderPhone)
	}
	if req.HolderEmail != "" {
		cashier.CashierValues.SetHolderEmail(req.HolderEmail)
	}
	if req.Country != "" {
		cashier.CashierValues.SetCountry(req.Country)
	}
	if req.CountryCode != "" {
		cashier.CashierValues.SetCountryCode(req.CountryCode)
	}
	if req.Province != "" {
		cashier.CashierValues.SetProvince(req.Province)
	}
	if req.City != "" {
		cashier.CashierValues.SetCity(req.City)
	}
	if req.Currency != "" {
		cashier.CashierValues.SetCurrency(req.Currency)
	}
	if req.Usage > 0 {
		cashier.CashierValues.SetUsage(req.Usage)
	}
	if !req.DailyLimit.IsZero() {
		cashier.CashierValues.SetDailyLimit(req.DailyLimit)
	}
	if !req.MonthlyLimit.IsZero() {
		cashier.CashierValues.SetMonthlyLimit(req.MonthlyLimit)
	}
	if req.ExpireAt > 0 {
		cashier.CashierValues.SetExpireAt(req.ExpireAt)
	}
	if req.Logo != "" {
		cashier.CashierValues.SetLogo(req.Logo)
	}
	if req.Remark != "" {
		cashier.CashierValues.SetRemark(req.Remark)
	}

	if err := s.db.Create(cashier).Error; err != nil {
		return nil, fmt.Errorf("failed to create cashier: %v", err)
	}

	return cashier, nil
}

// GetCashier 获取收银员详情
func (s *CashierService) GetCashier(cashierID string) (*protocol.CashierResponse, error) {
	var cashier models.Cashier
	if err := s.db.Where("cashier_id = ?", cashierID).First(&cashier).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("cashier not found")
		}
		return nil, fmt.Errorf("failed to get cashier: %v", err)
	}

	return s.buildCashierResponse(&cashier), nil
}

// ListCashiers 获取收银员列表
func (s *CashierService) ListCashiers(req *protocol.ListCashiersRequest) ([]*protocol.CashierResponse, int64, error) {
	var cashiers []models.Cashier
	var total int64

	// 构建查询条件
	query := s.db.Model(&models.Cashier{})

	if req.AccountID != "" {
		query = query.Where("account_id = ?", req.AccountID)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.BankCode != "" {
		query = query.Where("bank_code = ?", req.BankCode)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count cashiers: %v", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.Size
	if err := query.Offset(offset).Limit(req.Size).Order("created_at DESC").Find(&cashiers).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list cashiers: %v", err)
	}

	// 转换响应
	responses := make([]*protocol.CashierResponse, len(cashiers))
	for i, cashier := range cashiers {
		responses[i] = s.buildCashierResponse(&cashier)
		// 对于列表，遮蔽卡号
		responses[i].CardNumber = maskCardNumber(responses[i].CardNumber)
	}

	return responses, total, nil
}

// UpdateCashier 更新收银员
func (s *CashierService) UpdateCashier(cashierID string, req *protocol.UpdateCashierRequest) (*protocol.CashierResponse, error) {
	var cashier models.Cashier
	if err := s.db.Where("cashier_id = ?", cashierID).First(&cashier).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("cashier not found")
		}
		return nil, fmt.Errorf("failed to get cashier: %v", err)
	}

	// 更新字段
	if req.HolderName != nil && *req.HolderName != "" {
		cashier.CashierValues.SetHolderName(*req.HolderName)
	}
	if req.BankName != nil && *req.BankName != "" {
		cashier.CashierValues.SetBankName(*req.BankName)
	}
	if req.HolderPhone != nil {
		cashier.CashierValues.SetHolderPhone(*req.HolderPhone)
	}
	if req.HolderEmail != nil {
		cashier.CashierValues.SetHolderEmail(*req.HolderEmail)
	}
	if req.Province != nil {
		cashier.CashierValues.SetProvince(*req.Province)
	}
	if req.City != nil {
		cashier.CashierValues.SetCity(*req.City)
	}
	if req.Currency != nil {
		cashier.CashierValues.SetCurrency(*req.Currency)
	}
	if req.Usage != nil {
		cashier.CashierValues.SetUsage(*req.Usage)
	}
	if req.Status != nil && *req.Status != "" {
		cashier.CashierValues.SetStatus(*req.Status)
	}
	if req.DailyLimit != nil {
		cashier.CashierValues.SetDailyLimit(*req.DailyLimit)
	}
	if req.MonthlyLimit != nil {
		cashier.CashierValues.SetMonthlyLimit(*req.MonthlyLimit)
	}
	if req.ExpireAt != nil {
		cashier.CashierValues.SetExpireAt(*req.ExpireAt)
	}
	if req.Logo != nil {
		cashier.CashierValues.SetLogo(*req.Logo)
	}
	if req.Remark != nil {
		cashier.CashierValues.SetRemark(*req.Remark)
	}

	cashier.UpdatedAt = getCurrentTimeMillis()

	if err := s.db.Save(&cashier).Error; err != nil {
		return nil, fmt.Errorf("failed to update cashier: %v", err)
	}

	return s.GetCashier(cashierID)
}

// DeleteCashier 删除收银员
func (s *CashierService) DeleteCashier(cashierID string) error {
	result := s.db.Where("cashier_id = ?", cashierID).Delete(&models.Cashier{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete cashier: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("cashier not found")
	}

	return nil
}

// GetCashierByAccountID 根据账户ID获取收银员
func (s *CashierService) GetCashierByAccountID(accountID string, cashierType string) (*protocol.CashierResponse, error) {
	var cashier models.Cashier
	query := s.db.Where("account_id = ? AND status = ?", accountID, "active")

	if cashierType != "" {
		query = query.Where("type = ?", cashierType)
	}

	if err := query.Order("created_at ASC").First(&cashier).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no active cashier found for account")
		} else {
			return nil, fmt.Errorf("failed to get cashier: %v", err)
		}
	}

	return s.buildCashierResponse(&cashier), nil
}

// buildCashierResponse 构建收银员响应
func (s *CashierService) buildCashierResponse(cashier *models.Cashier) *protocol.CashierResponse {
	return &protocol.CashierResponse{
		ID:           cashier.ID,
		CashierID:    cashier.CashierID,
		AccountID:    cashier.AccountID,
		Type:         cashier.CashierValues.GetType(),
		CardNumber:   cashier.CashierValues.GetCardNumber(),
		HolderName:   cashier.CashierValues.GetHolderName(),
		BankCode:     cashier.CashierValues.GetBankCode(),
		BankName:     cashier.CashierValues.GetBankName(),
		HolderPhone:  cashier.CashierValues.GetHolderPhone(),
		HolderEmail:  cashier.CashierValues.GetHolderEmail(),
		Country:      cashier.CashierValues.GetCountry(),
		CountryCode:  cashier.CashierValues.GetCountryCode(),
		Province:     cashier.CashierValues.GetProvince(),
		City:         cashier.CashierValues.GetCity(),
		Currency:     cashier.CashierValues.GetCurrency(),
		Usage:        cashier.CashierValues.GetUsage(),
		Status:       cashier.CashierValues.GetStatus(),
		DailyLimit:   cashier.CashierValues.GetDailyLimit(),
		MonthlyLimit: cashier.CashierValues.GetMonthlyLimit(),
		DailyUsed:    cashier.CashierValues.GetDailyUsed(),
		MonthlyUsed:  cashier.CashierValues.GetMonthlyUsed(),
		ExpireAt:     cashier.CashierValues.GetExpireAt(),
		Logo:         cashier.CashierValues.GetLogo(),
		Remark:       cashier.CashierValues.GetRemark(),
		CreatedAt:    cashier.CreatedAt,
		UpdatedAt:    cashier.UpdatedAt,
	}
}

// maskCardNumber 遮蔽卡号
func maskCardNumber(cardNumber string) string {
	if len(cardNumber) <= 8 {
		return cardNumber
	}

	// 显示前4位和后4位，中间用*遮蔽
	return cardNumber[:4] + "****" + cardNumber[len(cardNumber)-4:]
}

// 工具函数

// getCurrentTimeMillis 获取当前时间毫秒数
func getCurrentTimeMillis() int64 {
	return time.Now().UnixMilli()
}

// generateCashierID 生成收银员ID
func generateCashierID() string {
	return "cashier_" + strconv.FormatInt(time.Now().UnixNano(), 36)
}

// 全局服务实例
var cashierService *CashierService

// GetCashierService 获取收银员服务实例
func GetCashierService() *CashierService {
	if cashierService == nil {
		cashierService = NewCashierService()
	}
	return cashierService
}
