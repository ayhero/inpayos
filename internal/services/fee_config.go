package services

import (
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type FeeConfigService struct {
	db *gorm.DB
}

func NewFeeConfigService(db *gorm.DB) *FeeConfigService {
	return &FeeConfigService{db: db}
}

// CreateFeeConfig 创建费用配置
func (s *FeeConfigService) CreateFeeConfig(req *protocol.CreateFeeConfigRequest) protocol.ErrorCode {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	feeConfig := &models.FeeConfig{
		MerchantID:      req.MerchantID,
		TransactionType: req.TransactionType,
		Country:         req.Country,
		PaymentMethod:   req.PaymentMethod,
		CoinID:          req.CoinID,
		FeePercent:      req.FeePercent,
		FeeFixed:        req.FeeFixed,
		MinFee:          req.MinFee,
		MaxFee:          req.MaxFee,
		IsActive:        isActive,
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
	}

	if err := s.db.Create(feeConfig).Error; err != nil {
		return protocol.InternalError
	}

	return protocol.Success
}

// UpdateFeeConfig 更新费用配置
func (s *FeeConfigService) UpdateFeeConfig(req *protocol.UpdateFeeConfigRequest) protocol.ErrorCode {
	updates := make(map[string]interface{})

	if req.TransactionType != "" {
		updates["transaction_type"] = req.TransactionType
	}
	if req.Country != "" {
		updates["country"] = req.Country
	}
	if req.PaymentMethod != "" {
		updates["payment_method"] = req.PaymentMethod
	}
	if req.CoinID != "" {
		updates["coin_id"] = req.CoinID
	}
	if req.FeePercent != "" {
		updates["fee_percent"] = req.FeePercent
	}
	if req.FeeFixed != "" {
		updates["fee_fixed"] = req.FeeFixed
	}
	if req.MinFee != "" {
		updates["min_fee"] = req.MinFee
	}
	if req.MaxFee != "" {
		updates["max_fee"] = req.MaxFee
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	updates["updated_at"] = time.Now().Unix()

	result := s.db.Model(&models.FeeConfig{}).Where("id = ?", req.ID).Updates(updates)
	if result.Error != nil {
		return protocol.InternalError
	}
	if result.RowsAffected == 0 {
		return protocol.ConfigNotFound
	}

	return protocol.Success
}

// GetFeeConfig 获取费用配置
func (s *FeeConfigService) GetFeeConfig(id uint64) (*protocol.FeeConfigResponse, protocol.ErrorCode) {
	var feeConfig models.FeeConfig
	if err := s.db.Where("id = ?", id).First(&feeConfig).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, protocol.ConfigNotFound
		}
		return nil, protocol.InternalError
	}

	return s.convertToResponse(&feeConfig), protocol.Success
}

// QueryFeeConfigs 查询费用配置列表
func (s *FeeConfigService) QueryFeeConfigs(req *protocol.QueryFeeConfigRequest) (*protocol.FeeConfigsResponse, protocol.ErrorCode) {
	query := s.db.Model(&models.FeeConfig{})

	// 构建查询条件
	if req.MerchantID != "" {
		query = query.Where("merchant_id = ?", req.MerchantID)
	}
	if req.TransactionType != "" {
		query = query.Where("transaction_type = ?", req.TransactionType)
	}
	if req.Country != "" {
		query = query.Where("country = ?", req.Country)
	}
	if req.PaymentMethod != "" {
		query = query.Where("payment_method = ?", req.PaymentMethod)
	}
	if req.CoinID != "" {
		query = query.Where("coin_id = ?", req.CoinID)
	}
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, protocol.InternalError
	}

	// 分页查询
	var feeConfigs []models.FeeConfig
	offset := (req.Page - 1) * req.Size
	if err := query.Order("created_at DESC").
		Offset(offset).Limit(req.Size).Find(&feeConfigs).Error; err != nil {
		return nil, protocol.InternalError
	}

	// 转换响应
	configs := make([]protocol.FeeConfigResponse, len(feeConfigs))
	for i, config := range feeConfigs {
		configs[i] = *s.convertToResponse(&config)
	}

	return &protocol.FeeConfigsResponse{
		Configs: configs,
		Total:   total,
		Page:    req.Page,
		Size:    req.Size,
	}, protocol.Success
}

// DeleteFeeConfig 删除费用配置
func (s *FeeConfigService) DeleteFeeConfig(id uint64) protocol.ErrorCode {
	result := s.db.Delete(&models.FeeConfig{}, id)
	if result.Error != nil {
		return protocol.InternalError
	}
	if result.RowsAffected == 0 {
		return protocol.ConfigNotFound
	}

	return protocol.Success
}

// CalculateFee 计算费用
func (s *FeeConfigService) CalculateFee(req *protocol.CalculateFeeRequest) (*protocol.CalculateFeeResponse, protocol.ErrorCode) {
	// 查找最匹配的费用配置
	query := s.db.Model(&models.FeeConfig{}).
		Where("merchant_id = ? AND transaction_type = ?", req.MerchantID, req.TransactionType).
		Where("is_active = ?", true)

	// 条件匹配优先级
	if req.Country != "" {
		query = query.Where("(country = '' OR country = ?)", req.Country)
	}
	if req.PaymentMethod != "" {
		query = query.Where("(payment_method = '' OR payment_method = ?)", req.PaymentMethod)
	}
	if req.Currency != "" {
		query = query.Where("(coin_id = '' OR coin_id = ?)", req.Currency)
	}

	var feeConfig models.FeeConfig
	if err := query.Order("country DESC, payment_method DESC, coin_id DESC").First(&feeConfig).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, protocol.ConfigNotFound
		}
		return nil, protocol.InternalError
	}

	// 解析金额
	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return nil, protocol.InvalidAmount
	}

	// 计算费用
	var feeAmount float64
	feePercent, _ := strconv.ParseFloat(feeConfig.FeePercent, 64)
	feeFixed, _ := strconv.ParseFloat(feeConfig.FeeFixed, 64)

	// 百分比费用 + 固定费用
	feeAmount = amount*feePercent/100 + feeFixed

	// 应用费用限制
	if feeConfig.MinFee != "" {
		minFee, _ := strconv.ParseFloat(feeConfig.MinFee, 64)
		if feeAmount < minFee {
			feeAmount = minFee
		}
	}
	if feeConfig.MaxFee != "" {
		maxFee, _ := strconv.ParseFloat(feeConfig.MaxFee, 64)
		if feeAmount > maxFee {
			feeAmount = maxFee
		}
	}

	return &protocol.CalculateFeeResponse{
		Amount:       req.Amount,
		Fee:          strconv.FormatFloat(feeAmount, 'f', 2, 64),
		TotalAmount:  strconv.FormatFloat(amount+feeAmount, 'f', 2, 64),
		Currency:     req.Currency,
		FeeConfigID:  feeConfig.ID,
		CalculatedAt: time.Now().Unix(),
		FeeBreakdown: map[string]string{
			"percentage_fee": strconv.FormatFloat(amount*feePercent/100, 'f', 2, 64),
			"fixed_fee":      feeConfig.FeeFixed,
			"total_fee":      strconv.FormatFloat(feeAmount, 'f', 2, 64),
		},
	}, protocol.Success
}

// convertToResponse 转换为响应格式
func (s *FeeConfigService) convertToResponse(feeConfig *models.FeeConfig) *protocol.FeeConfigResponse {
	return &protocol.FeeConfigResponse{
		ID:              feeConfig.ID,
		MerchantID:      feeConfig.MerchantID,
		TransactionType: feeConfig.TransactionType,
		Country:         feeConfig.Country,
		PaymentMethod:   feeConfig.PaymentMethod,
		CoinID:          feeConfig.CoinID,
		FeePercent:      feeConfig.FeePercent,
		FeeFixed:        feeConfig.FeeFixed,
		MinFee:          feeConfig.MinFee,
		MaxFee:          feeConfig.MaxFee,
		IsActive:        feeConfig.IsActive,
		CreatedAt:       feeConfig.CreatedAt,
		UpdatedAt:       feeConfig.UpdatedAt,
	}
}
