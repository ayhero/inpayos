package protocol

import (
	"context"
	"sync"

	"github.com/shopspring/decimal"
)

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

// 结算策略缓存的 context key
type contextKey string

const (
	settleStrategyCacheKey contextKey = "settle_strategy_cache"
	cacheStatsKey          contextKey = "cache_stats"
)

// SettleStrategyCache 结算策略缓存
type SettleStrategyCache struct {
	strategies map[string][]*SettleStrategy
	mutex      sync.RWMutex
}

// CacheStats 缓存统计信息
type CacheStats struct {
	hits   int64
	misses int64
	mutex  sync.RWMutex
}

// NewSettleStrategyCache 创建新的结算策略缓存
func NewSettleStrategyCache() *SettleStrategyCache {
	return &SettleStrategyCache{
		strategies: make(map[string][]*SettleStrategy),
	}
}

// NewCacheStats 创建新的缓存统计
func NewCacheStats() *CacheStats {
	return &CacheStats{}
}

// Get 从缓存中获取策略
func (c *SettleStrategyCache) Get(mid string) ([]*SettleStrategy, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	strategies, exists := c.strategies[mid]
	return strategies, exists
}

// Set 将策略存入缓存
func (c *SettleStrategyCache) Set(mid string, strategies []*SettleStrategy) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.strategies[mid] = strategies
}

// Size 获取缓存大小
func (c *SettleStrategyCache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.strategies)
}

// Clear 清空缓存
func (c *SettleStrategyCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.strategies = make(map[string][]*SettleStrategy)
}

// Delete 从缓存中删除指定 mid 的策略
func (c *SettleStrategyCache) Delete(mid string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.strategies, mid)
}

// Keys 获取缓存中所有的 mid
func (c *SettleStrategyCache) Keys() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	keys := make([]string, 0, len(c.strategies))
	for key := range c.strategies {
		keys = append(keys, key)
	}
	return keys
}

// RecordHit 记录缓存命中
func (s *CacheStats) RecordHit() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.hits++
}

// RecordMiss 记录缓存未命中
func (s *CacheStats) RecordMiss() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.misses++
}

// Reset 重置缓存统计
func (s *CacheStats) Reset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.hits = 0
	s.misses = 0
}

// GetStats 获取统计信息
func (s *CacheStats) GetStats() (hits, misses int64) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.hits, s.misses
}

// GetHitRate 获取缓存命中率
func (s *CacheStats) GetHitRate() float64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.hits+s.misses == 0 {
		return 0
	}
	return float64(s.hits) / float64(s.hits+s.misses) * 100
}

// CreateSettleContext 创建包含缓存的 context
func CreateSettleContext(ctx context.Context) context.Context {
	cache := NewSettleStrategyCache()
	stats := NewCacheStats()

	ctx = context.WithValue(ctx, settleStrategyCacheKey, cache)
	ctx = context.WithValue(ctx, cacheStatsKey, stats)

	return ctx
}

// GetCacheFromContext 从 context 中获取缓存和统计信息
func GetCacheFromContext(ctx context.Context) (*SettleStrategyCache, *CacheStats) {
	cache, _ := ctx.Value(settleStrategyCacheKey).(*SettleStrategyCache)
	stats, _ := ctx.Value(cacheStatsKey).(*CacheStats)

	if cache == nil {
		cache = NewSettleStrategyCache()
	}
	if stats == nil {
		stats = NewCacheStats()
	}

	return cache, stats
}
