package protocol

import (
	"inpayos/internal/utils"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

const (
	MerchantSettle                 = "merchant.settle"
	MerchantSettleProcess          = "merchant.settle.process"
	MerchantSettleAccounting       = "merchant.settle.accounting"
	MerchantTransactionSettleIDFix = "merchant.transaction.settle.id.fix"
)

const (
	T0 = "T+0" // T+0 立即结算
	T1 = "T+1" // T+1 次日结算
	T2 = "T+2" // T+2 两日结算
	T3 = "T+3" // T+3 三日结算
	W1 = "W+1" // W+1 下周结算
	M1 = "M+1" // M+1 下月结算
)

// SettleAccountingResult 结算记账结果统计
type SettleAccountingResult struct {
	TotalCount   int64
	SuccessCount int64
	FailedCount  int64
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
	Result       string
}

// SettleFixResult 交易settle_id修复结果统计
type SettleFixResult struct {
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	Duration     time.Duration `json:"duration"`
	TotalCount   int64         `json:"total_count"`
	SuccessCount int64         `json:"success_count"`
	FailedCount  int64         `json:"failed_count"`
	Result       string        `json:"result"`
}

// SettleResult 结算结果统计
type SettleResult struct {
	TotalCount     int64
	SuccessCount   int64
	FailedCount    int64
	ProcessedPages int
	TotalPages     int
	StartTime      time.Time
	EndTime        time.Time
	Duration       time.Duration
	Result         string
}

// SettlementResult represents the result of settlement calculation
type SettlementResult struct {
	SettleAmount    decimal.Decimal
	SettleUsdAmount decimal.Decimal
	Fee             decimal.Decimal
	UsdFee          decimal.Decimal
	FixedFee        decimal.Decimal
	FixedUsdFee     decimal.Decimal
	Rate            decimal.Decimal
	UsdRate         decimal.Decimal
	FeeCcy          string
}

type SettleStrategy struct {
	ID        int64       `json:"id" gorm:"primaryKey;autoIncrement;comment:自增ID"`
	MID       int64       `json:"mid" gorm:"index;comment:商户ID"`        // MerchantID 商户ID
	Period    int64       `json:"period" gorm:"index;comment:结算周期"`     // SettlePeriod 结算周期
	SettleCcy string      `json:"settle_ccy" gorm:"index;comment:币种"`   // Currency 币种
	TrxType   string      `json:"trx_type" gorm:"index;comment:交易类型"`   // TrxType 交易类型
	TrxMode   string      `json:"trx_mode" gorm:"index;comment:交易模式"`   // TrxMode 交易模式
	TrxMethod string      `json:"trx_method" gorm:"index;comment:交易方式"` // TrxMethod 交易方式
	Country   string      `json:"country" gorm:"index;comment:国家"`      // Country 国家
	TrxCcy    string      `json:"trx_ccy" gorm:"index;comment:交易币种"`    // TrxCcy 交易币种
	Status    int64       `json:"status" gorm:"index;comment:状态"`       // Status 状态
	Rules     SettleRules `json:"rules" gorm:"comment:结算规则"`            // SettleRules 结算规则
	CreatedAt int64       `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt int64       `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}
type SettleStrategies []*SettleStrategy

type SettleRule struct {
	Period      int64            `json:"period" gorm:"index;comment:结算周期"`     // Period 结算周期
	RuleID      string           `json:"rule_id" gorm:"index;comment:规则ID"`    // RuleID 规则ID
	TrxType     string           `json:"trx_type" gorm:"index;comment:交易类型"`   // TrxType 交易类型
	TrxMode     string           `json:"trx_mode" gorm:"index;comment:交易模式"`   // TrxMode 交易模式
	TrxMethod   string           `json:"trx_method" gorm:"index;comment:交易方式"` // TrxMethod 交易方式
	Country     string           `json:"country" gorm:"index;comment:国家"`      // Country 国家
	MinAmount   *decimal.Decimal `json:"min_amount" gorm:"comment:最小结算金额"`     // MinSettleAmount 最小结算金额
	MaxAmount   *decimal.Decimal `json:"max_amount" gorm:"comment:最大结算金额"`     // MaxSettleAmount 最大结算金额
	MinFee      *decimal.Decimal `json:"min_fee" gorm:"comment:最小手续费"`         // MinSettleFee 最小手续费
	MaxFee      *decimal.Decimal `json:"max_fee" gorm:"comment:最大手续费"`         // MaxSettleFee 最大手续费
	MinRate     *decimal.Decimal `json:"min_rate" gorm:"comment:最小费率"`         // MinSettleRate 最小费率
	MaxRate     *decimal.Decimal `json:"max_rate" gorm:"comment:最大费率"`         // MaxSettleRate 最大费率
	MinUsdFee   *decimal.Decimal `json:"min_usd_fee" gorm:"comment:最小美元手续费"`   // MinSettleUsdFee 最小美元手续费
	MaxUsdFee   *decimal.Decimal `json:"max_usd_fee" gorm:"comment:最大美元手续费"`   // MaxSettleUsdFee 最大美元手续费
	MinUsdRate  *decimal.Decimal `json:"min_usd_rate" gorm:"comment:最小美元费率"`   // MinSettleUsdRate 最小美元费率
	MaxUsdRate  *decimal.Decimal `json:"max_usd_rate" gorm:"comment:最大美元费率"`   // MaxSettleUsdRate 最大美元费率
	Ccy         string           `json:"ccy" gorm:"comment:币种"`                // CurrencyType 币种
	FixedFee    *decimal.Decimal `json:"fixed_fee" gorm:"comment:固定手续费"`       // FixedSettleFee 固定手续费
	Rate        *decimal.Decimal `json:"rate" gorm:"comment:费率"`               // SettleRate 费率
	FixedUsdFee *decimal.Decimal `json:"fixed_usd_fee" gorm:"comment:固定美元手续费"` // FixedSettleUsdFee 固定美元手续费
	UsdRate     *decimal.Decimal `json:"usd_rate" gorm:"comment:美元费率"`         // SettleUsdRate 美元费率
}

type SettleRules []*SettleRule

// MerchantSettleContext 结算策略缓存（使用 sync.Map 实现线程安全）
type MerchantSettleContext struct {
	strategies sync.Map // map[string][]*SettleStrategy
	SettledAt  int64    // 当前结算时间
}

// NewMerchantSettleContext 创建新的结算策略缓存
func NewMerchantSettleContext() *MerchantSettleContext {
	return &MerchantSettleContext{
		strategies: sync.Map{},
		SettledAt:  utils.TimeNowMilli(),
	}
}

// Get 从缓存中获取策略
func (c *MerchantSettleContext) Get(mid string) []*SettleStrategy {
	value, exists := c.strategies.Load(mid)
	if !exists {
		return nil
	}
	strategies, ok := value.([]*SettleStrategy)
	if !ok {
		return nil
	}
	return strategies
}

// Set 将策略存入缓存
func (c *MerchantSettleContext) Set(mid string, strategies []*SettleStrategy) {
	c.strategies.Store(mid, strategies)
}

// Size 获取缓存大小
func (c *MerchantSettleContext) Size() int {
	count := 0
	c.strategies.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// Clear 清空缓存
func (c *MerchantSettleContext) Clear() {
	c.strategies.Range(func(key, value interface{}) bool {
		c.strategies.Delete(key)
		return true
	})
}

// Delete 从缓存中删除指定 mid 的策略
func (c *MerchantSettleContext) Delete(mid string) {
	c.strategies.Delete(mid)
}

// Keys 获取缓存中所有的 mid
func (c *MerchantSettleContext) Keys() []string {
	keys := make([]string, 0)
	c.strategies.Range(func(key, value interface{}) bool {
		if mid, ok := key.(string); ok {
			keys = append(keys, mid)
		}
		return true
	})
	return keys
}

// LoadOrStore 原子性地加载或存储策略
func (c *MerchantSettleContext) LoadOrStore(mid string, strategies []*SettleStrategy) ([]*SettleStrategy, bool) {
	value, loaded := c.strategies.LoadOrStore(mid, strategies)
	if loaded {
		if existingStrategies, ok := value.([]*SettleStrategy); ok {
			return existingStrategies, true
		}
	}
	return strategies, false
}

// Range 遍历所有缓存项
func (c *MerchantSettleContext) Range(fn func(mid string, strategies []*SettleStrategy) bool) {
	c.strategies.Range(func(key, value interface{}) bool {
		if mid, ok := key.(string); ok {
			if strategies, ok := value.([]*SettleStrategy); ok {
				return fn(mid, strategies)
			}
		}
		return true
	})
}
