package services

import (
	"fmt"
	"inpayos/internal/models"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ConfigService 配置服务
type ConfigService struct {
	db *gorm.DB
}

// NewConfigService 创建配置服务
func NewConfigService() *ConfigService {
	return &ConfigService{
		db: models.GetDB(),
	}
}

// GetTrxConfigByMerchantID 根据商户ID获取交易配置
func (s *ConfigService) GetTrxConfigByMerchantID(merchantID, trxType string) *models.TrxConfig {
	config := s.DefaultTrxConfig()

	// 全局配置优先
	globalConfig := s.getMerchantConfig(models.GlobalMerchantID, trxType)
	if globalConfig != nil {
		if trxConfig, err := globalConfig.GetConfig(); err == nil {
			config.Copy(trxConfig)
		}
	}

	// 商户配置覆盖全局配置
	if merchantID != "" && merchantID != models.GlobalMerchantID {
		merchantConfig := s.getMerchantConfig(merchantID, trxType)
		if merchantConfig != nil {
			if trxConfig, err := merchantConfig.GetConfig(); err == nil {
				config.Copy(trxConfig)
			}
		}
	}

	return config
}

// SaveMerchantConfig 保存商户配置
func (s *ConfigService) SaveMerchantConfig(merchantID, trxType string, config *models.TrxConfig) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var record models.MerchantConfig
		err := tx.Where("merchant_id = ? AND type = ?", merchantID, trxType).First(&record).Error

		if err == gorm.ErrRecordNotFound {
			// 新建记录
			record = models.MerchantConfig{
				Mid:  merchantID,
				Type: trxType,
			}
			if err := record.SetConfig(config); err != nil {
				return fmt.Errorf("failed to set config: %w", err)
			}
			return tx.Create(&record).Error
		} else if err != nil {
			return fmt.Errorf("failed to query config: %w", err)
		}

		// 更新记录
		if err := record.SetConfig(config); err != nil {
			return fmt.Errorf("failed to set config: %w", err)
		}
		record.UpdatedAt = time.Now().UnixMilli()
		return tx.Save(&record).Error
	})
}

// GetMerchantConfig 获取商户配置记录
func (s *ConfigService) GetMerchantConfig(merchantID, trxType string) (*models.MerchantConfig, error) {
	var record models.MerchantConfig
	err := s.db.Where("merchant_id = ? AND type = ?", merchantID, trxType).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get merchant config: %w", err)
	}
	return &record, nil
}

// DeleteMerchantConfig 删除商户配置
func (s *ConfigService) DeleteMerchantConfig(merchantID, trxType string) error {
	result := s.db.Where("merchant_id = ? AND type = ?", merchantID, trxType).Delete(&models.MerchantConfig{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete config: %w", result.Error)
	}
	return nil
}

// ListMerchantConfigs 列表商户配置
func (s *ConfigService) ListMerchantConfigs(merchantID string) ([]*models.MerchantConfig, error) {
	var configs []*models.MerchantConfig
	query := s.db.Model(&models.MerchantConfig{})

	if merchantID != "" {
		query = query.Where("merchant_id = ?", merchantID)
	}

	err := query.Order("created_at DESC").Find(&configs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list configs: %w", err)
	}

	return configs, nil
}

// DefaultTrxConfig 默认交易配置
func (s *ConfigService) DefaultTrxConfig() *models.TrxConfig {
	return &models.TrxConfig{
		MinAmount: map[string]string{
			"CNY": "0.01",
			"USD": "0.01",
		},
		MaxAmount: map[string]string{
			"CNY": "100000.00",
			"USD": "10000.00",
		},
		DailyLimit: map[string]string{
			"CNY": "1000000.00",
			"USD": "100000.00",
		},
		MonthlyLimit: map[string]string{
			"CNY": "10000000.00",
			"USD": "1000000.00",
		},
		FeeRate: map[string]string{
			"CNY": "0.006",
			"USD": "0.006",
		},
		Status:        models.ConfigStatusOn,
		AutoConfirm:   models.ConfigStatusOff,
		TimeoutMinute: 30,
	}
}

// InitDefaultConfigs 初始化默认配置
func (s *ConfigService) InitDefaultConfigs() error {
	trxTypes := []string{
		models.TrxTypeReceipt,
		models.TrxTypePayment,
		models.TrxTypeDeposit,
		models.TrxTypeWithdraw,
		models.TrxTypeRefund,
	}

	for _, trxType := range trxTypes {
		// 检查全局配置是否存在
		existing := s.getMerchantConfig(models.GlobalMerchantID, trxType)
		if existing == nil {
			// 创建默认全局配置
			config := s.DefaultTrxConfig()

			// 针对不同交易类型调整默认配置
			switch trxType {
			case models.TrxTypePayment:
				config.AutoConfirm = models.ConfigStatusOff // 代付需要人工确认
				config.TimeoutMinute = 60
			case models.TrxTypeDeposit:
				config.AutoConfirm = models.ConfigStatusOn // 充值自动确认
				config.TimeoutMinute = 15
			case models.TrxTypeWithdraw:
				config.AutoConfirm = models.ConfigStatusOff // 提现需要人工确认
				config.TimeoutMinute = 60
			case models.TrxTypeRefund:
				config.AutoConfirm = models.ConfigStatusOff // 退款需要人工确认
				config.TimeoutMinute = 120
			}

			if err := s.SaveMerchantConfig(models.GlobalMerchantID, trxType, config); err != nil {
				return fmt.Errorf("failed to init default config for %s: %w", trxType, err)
			}
		}
	}

	return nil
}

// CalculateFee 计算手续费
func (s *ConfigService) CalculateFee(merchantID, trxType, currency string, amount decimal.Decimal) decimal.Decimal {
	config := s.GetTrxConfigByMerchantID(merchantID, trxType)

	// 固定费用优先
	if config.FeeFixed != nil {
		if fixedFee, exists := config.FeeFixed[currency]; exists {
			if fee, err := decimal.NewFromString(fixedFee); err == nil {
				return fee
			}
		}
	}

	// 费率计算
	feeRate := config.GetFeeRate(currency)
	if feeRate.IsZero() {
		return decimal.Zero
	}

	return amount.Mul(feeRate)
}

// ValidateAmount 验证金额
func (s *ConfigService) ValidateAmount(merchantID, trxType, currency string, amount decimal.Decimal) error {
	config := s.GetTrxConfigByMerchantID(merchantID, trxType)

	// 检查最小金额
	minAmount := config.GetMinAmount(currency)
	if !minAmount.IsZero() && amount.LessThan(minAmount) {
		return fmt.Errorf("amount %s is less than minimum %s", amount.String(), minAmount.String())
	}

	// 检查最大金额
	maxAmount := config.GetMaxAmount(currency)
	if !maxAmount.IsZero() && amount.GreaterThan(maxAmount) {
		return fmt.Errorf("amount %s is greater than maximum %s", amount.String(), maxAmount.String())
	}

	return nil
}

// getMerchantConfig 获取商户配置记录（内部方法）
func (s *ConfigService) getMerchantConfig(merchantID, trxType string) *models.MerchantConfig {
	var record models.MerchantConfig
	err := s.db.Where("merchant_id = ? AND type = ?", merchantID, trxType).First(&record).Error
	if err != nil {
		return nil
	}
	return &record
}

// 全局服务实例
var configService *ConfigService

// GetConfigService 获取配置服务实例
func GetConfigService() *ConfigService {
	if configService == nil {
		configService = NewConfigService()
	}
	return configService
}
