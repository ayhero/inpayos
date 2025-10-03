package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"inpayos/internal/models"
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
)

type ChannelService struct{}

func NewChannelService() *ChannelService {
	return &ChannelService{}
}

// 生成渠道ID
func (s *ChannelService) generateChannelID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("生成渠道ID失败: %v", err)
	}
	return "channel_" + hex.EncodeToString(bytes), nil
}

// CreateChannel 创建渠道
func (s *ChannelService) CreateChannel(req *protocol.CreateChannelRequest) (*models.Channel, error) {
	channelID, err := s.generateChannelID()
	if err != nil {
		return nil, err
	}

	channel := &models.Channel{
		ChannelID:     channelID,
		ChannelValues: &models.ChannelValues{},
	}

	// 设置渠道参数
	channel.ChannelValues.SetCode(req.Code).
		SetName(req.Name).
		SetType(req.Type).
		SetStatus(req.Status)

	// 设置费用
	if req.FeeType == "fixed" {
		channel.ChannelValues.SetFeeFixed(req.FeeValue)
	} else if req.FeeType == "percent" {
		channel.ChannelValues.SetFeePercent(req.FeeValue)
	}

	// 设置金额限制
	if !req.MinAmount.IsZero() {
		channel.ChannelValues.SetMinAmount(req.MinAmount)
	}
	if !req.MaxAmount.IsZero() {
		channel.ChannelValues.SetMaxAmount(req.MaxAmount)
	}
	if !req.DailyLimit.IsZero() {
		channel.ChannelValues.SetDailyLimit(req.DailyLimit)
	}
	if !req.MonthlyLimit.IsZero() {
		channel.ChannelValues.SetMonthlyLimit(req.MonthlyLimit)
	}

	// 设置配置和备注
	if req.Config != nil {
		configStr, _ := json.Marshal(req.Config)
		channel.ChannelValues.SetConfig(string(configStr))
	}
	if req.Remark != "" {
		channel.ChannelValues.SetRemark(req.Remark)
	}

	if err := models.WriteDB.Create(channel).Error; err != nil {
		return nil, fmt.Errorf("创建渠道失败: %v", err)
	}

	return channel, nil
}

// GetChannel 获取渠道
func (s *ChannelService) GetChannel(channelID string) (*models.Channel, error) {
	var channel models.Channel
	if err := models.WriteDB.Where("channel_id = ?", channelID).First(&channel).Error; err != nil {
		return nil, fmt.Errorf("渠道不存在")
	}
	return &channel, nil
}

// GetChannelByCode 根据代码获取渠道
func (s *ChannelService) GetChannelByCode(code string) (*models.Channel, error) {
	var channel models.Channel
	if err := models.WriteDB.Where("code = ?", code).First(&channel).Error; err != nil {
		return nil, fmt.Errorf("渠道不存在")
	}
	return &channel, nil
}

// ListChannels 列出渠道
func (s *ChannelService) ListChannels(req *protocol.ListChannelsRequest) ([]*models.Channel, int64, error) {
	db := models.WriteDB.Model(&models.Channel{})

	// 过滤条件
	if req.Code != "" {
		db = db.Where("code = ?", req.Code)
	}
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Type != "" {
		db = db.Where("type = ?", req.Type)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.Currency != "" {
		db = db.Where("JSON_CONTAINS(currencies, ?)", fmt.Sprintf(`"%s"`, req.Currency))
	}
	if req.PayMethod != "" {
		db = db.Where("JSON_CONTAINS(pay_methods, ?)", fmt.Sprintf(`"%s"`, req.PayMethod))
	}

	// 计算总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("计算渠道总数失败: %v", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.Size
	var channels []*models.Channel
	if err := db.Offset(offset).
		Limit(req.Size).
		Order("priority DESC, created_at DESC").
		Find(&channels).Error; err != nil {
		return nil, 0, fmt.Errorf("查询渠道列表失败: %v", err)
	}

	return channels, total, nil
}

// UpdateChannel 更新渠道
func (s *ChannelService) UpdateChannel(channelID string, req *protocol.UpdateChannelRequest) (*models.Channel, error) {
	channel, err := s.GetChannel(channelID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Name != nil {
		channel.ChannelValues.SetName(*req.Name)
	}
	if req.Status != nil {
		channel.ChannelValues.SetStatus(*req.Status)
	}
	if req.Type != nil {
		channel.ChannelValues.SetType(*req.Type)
	}

	// 更新费用
	if req.FeeType != nil && req.FeeValue != nil {
		if *req.FeeType == "fixed" {
			channel.ChannelValues.SetFeeFixed(*req.FeeValue)
			channel.ChannelValues.SetFeePercent(decimal.Zero)
		} else if *req.FeeType == "percent" {
			channel.ChannelValues.SetFeePercent(*req.FeeValue)
			channel.ChannelValues.SetFeeFixed(decimal.Zero)
		}
	}

	// 更新金额限制
	if req.MinAmount != nil {
		channel.ChannelValues.SetMinAmount(*req.MinAmount)
	}
	if req.MaxAmount != nil {
		channel.ChannelValues.SetMaxAmount(*req.MaxAmount)
	}
	if req.DailyLimit != nil {
		channel.ChannelValues.SetDailyLimit(*req.DailyLimit)
	}
	if req.MonthlyLimit != nil {
		channel.ChannelValues.SetMonthlyLimit(*req.MonthlyLimit)
	}

	// 更新配置和备注
	if req.Config != nil {
		configStr, _ := json.Marshal(req.Config)
		channel.ChannelValues.SetConfig(string(configStr))
	}
	if req.Remark != nil {
		channel.ChannelValues.SetRemark(*req.Remark)
	}

	if err := models.WriteDB.Save(channel).Error; err != nil {
		return nil, fmt.Errorf("更新渠道失败: %v", err)
	}

	return channel, nil
}

// DeleteChannel 删除渠道
func (s *ChannelService) DeleteChannel(channelID string) error {
	if err := models.WriteDB.Where("channel_id = ?", channelID).Delete(&models.Channel{}).Error; err != nil {
		return fmt.Errorf("删除渠道失败: %v", err)
	}
	return nil
}

// GetAvailableChannels 获取可用渠道
func (s *ChannelService) GetAvailableChannels(txType, currency, payMethod string, amount decimal.Decimal) ([]*models.Channel, error) {
	db := models.WriteDB.Model(&models.Channel{}).Where("status = ?", "active")

	// 过滤条件
	if txType != "" {
		db = db.Where("type = ?", txType)
	}
	if currency != "" {
		db = db.Where("JSON_CONTAINS(currencies, ?)", fmt.Sprintf(`"%s"`, currency))
	}
	if payMethod != "" {
		db = db.Where("JSON_CONTAINS(pay_methods, ?)", fmt.Sprintf(`"%s"`, payMethod))
	}

	var channels []*models.Channel
	if err := db.Order("priority DESC").Find(&channels).Error; err != nil {
		return nil, fmt.Errorf("查询可用渠道失败: %v", err)
	}

	// 过滤可用渠道
	var availableChannels []*models.Channel
	for _, channel := range channels {
		if s.isChannelAvailable(channel, amount) {
			availableChannels = append(availableChannels, channel)
		}
	}

	return availableChannels, nil
}

// CalculateFee 计算手续费
func (s *ChannelService) CalculateFee(channelID string, amount decimal.Decimal) (decimal.Decimal, error) {
	channel, err := s.GetChannel(channelID)
	if err != nil {
		return decimal.Zero, err
	}

	return channel.ChannelValues.CalculateFee(amount), nil
}

// SelectBestChannel 选择最佳渠道
func (s *ChannelService) SelectBestChannel(txType, currency, payMethod string, amount decimal.Decimal) (*models.Channel, error) {
	channels, err := s.GetAvailableChannels(txType, currency, payMethod, amount)
	if err != nil {
		return nil, err
	}

	if len(channels) == 0 {
		return nil, fmt.Errorf("没有可用渠道")
	}

	// 选择第一个可用渠道（已按优先级排序）
	return channels[0], nil
}

// GetChannelStats 获取渠道统计
func (s *ChannelService) GetChannelStats() (*protocol.ChannelStatsResponse, error) {
	var totalChannels, activeChannels, disabledChannels, maintainChannels int64

	// 总渠道数
	if err := models.WriteDB.Model(&models.Channel{}).Count(&totalChannels).Error; err != nil {
		return nil, fmt.Errorf("查询总渠道数失败: %v", err)
	}

	// 活跃渠道数
	if err := models.WriteDB.Model(&models.Channel{}).Where("status = ?", "active").Count(&activeChannels).Error; err != nil {
		return nil, fmt.Errorf("查询活跃渠道数失败: %v", err)
	}

	// 禁用渠道数
	if err := models.WriteDB.Model(&models.Channel{}).Where("status = ?", "inactive").Count(&disabledChannels).Error; err != nil {
		return nil, fmt.Errorf("查询禁用渠道数失败: %v", err)
	}

	// 维护中渠道数
	if err := models.WriteDB.Model(&models.Channel{}).Where("status = ?", "maintain").Count(&maintainChannels).Error; err != nil {
		return nil, fmt.Errorf("查询维护中渠道数失败: %v", err)
	}

	return &protocol.ChannelStatsResponse{
		TotalChannels:    totalChannels,
		ActiveChannels:   activeChannels,
		DisabledChannels: disabledChannels,
		MaintainChannels: maintainChannels,
	}, nil
}

// isChannelAvailable 检查渠道是否可用
func (s *ChannelService) isChannelAvailable(channel *models.Channel, amount decimal.Decimal) bool {
	// 检查状态
	if !channel.ChannelValues.IsActive() {
		return false
	}

	// 检查金额范围
	if !channel.ChannelValues.IsAmountValid(amount) {
		return false
	}

	// 检查每日限额
	dailyUsed := s.getDailyUsedAmount(channel.ChannelValues.GetCode())
	dailyLimit := channel.ChannelValues.GetDailyLimit()
	if dailyLimit.GreaterThan(decimal.Zero) && dailyUsed.Add(amount).GreaterThan(dailyLimit) {
		return false
	}

	// 检查每月限额
	monthlyUsed := s.getMonthlyUsedAmount(channel.ChannelValues.GetCode())
	monthlyLimit := channel.ChannelValues.GetMonthlyLimit()
	if monthlyLimit.GreaterThan(decimal.Zero) && monthlyUsed.Add(amount).GreaterThan(monthlyLimit) {
		return false
	}

	return true
}

// getDailyUsedAmount 获取每日已使用金额
func (s *ChannelService) getDailyUsedAmount(channelCode string) decimal.Decimal {
	today := time.Now().Truncate(24 * time.Hour).UnixMilli()
	tomorrow := today + 24*60*60*1000

	var amount decimal.Decimal
	models.WriteDB.Raw(`
		SELECT COALESCE(SUM(amount), 0) FROM (
			SELECT amount FROM t_receipts WHERE channel_code = ? AND status = 'success' AND created_at >= ? AND created_at < ?
			UNION ALL
			SELECT amount FROM t_payments WHERE channel_code = ? AND status = 'success' AND created_at >= ? AND created_at < ?
		) AS total
	`, channelCode, today, tomorrow, channelCode, today, tomorrow).Scan(&amount)

	return amount
}

// getMonthlyUsedAmount 获取每月已使用金额
func (s *ChannelService) getMonthlyUsedAmount(channelCode string) decimal.Decimal {
	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).UnixMilli()
	nextMonth := firstDay + int64(now.AddDate(0, 1, -now.Day()+1).Sub(now).Milliseconds())

	var amount decimal.Decimal
	models.WriteDB.Raw(`
		SELECT COALESCE(SUM(amount), 0) FROM (
			SELECT amount FROM t_receipts WHERE channel_code = ? AND status = 'success' AND created_at >= ? AND created_at < ?
			UNION ALL
			SELECT amount FROM t_payments WHERE channel_code = ? AND status = 'success' AND created_at >= ? AND created_at < ?
		) AS total
	`, channelCode, firstDay, nextMonth, channelCode, firstDay, nextMonth).Scan(&amount)

	return amount
}
