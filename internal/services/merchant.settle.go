package services

import (
	"context"
	"fmt"
	"inpayos/internal/config"
	"inpayos/internal/log"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
	"strconv"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// MerchantSettleService 结算服务
type MerchantSettleService struct {
	*config.SettleConfig
}

var (
	merchantSettleServiceInstance *MerchantSettleService
	merchantSettleServiceOnce     sync.Once
)

// GetMerchantSettleService 获取结算服务单例
func GetMerchantSettleService() *MerchantSettleService {
	merchantSettleServiceOnce.Do(func() {
		SetupSettleService()
	})
	return merchantSettleServiceInstance
}

// SetupSettleService 设置结算服务
func SetupSettleService() {
	merchantSettleServiceInstance = &MerchantSettleService{
		SettleConfig: config.Get().Settle,
	}
}

// GetSettleStrategiesWithPeriodCache 基于结算周期记录使用缓存获取结算策略
// 这是新的缓存策略，优先从结算周期记录获取策略，实现周期级别的缓存
func (s *MerchantSettleService) GetSettleStrategiesWithPeriodCache(ctx context.Context, settleLog *models.MerchantSettleLog) []*protocol.SettleStrategy {
	if settleLog == nil {
		return nil
	}

	// 从 context 获取缓存和统计信息
	cache, stats := protocol.GetCacheFromContext(ctx)

	// 使用结算周期记录的唯一标识作为缓存key
	// 格式: settleperiod_{商户ID}_{周期}
	cacheKey := fmt.Sprintf("settleperiod_%d_%d",
		settleLog.MID,
		settleLog.Period)

	// 首先检查缓存
	strategies, exists := cache.Get(cacheKey)
	if exists {
		stats.RecordHit()
		log.Get().Debugf("GetSettleStrategiesWithPeriodCache: cache hit for settle log %s, found %d strategies", settleLog.SettleID, len(strategies))
		return strategies
	}

	// 缓存未命中，从结算周期记录或数据库查询
	stats.RecordMiss()
	log.Get().Debugf("GetSettleStrategiesWithPeriodCache: cache miss for settle log %s, retrieving strategies", settleLog.SettleID)

	var strategies_result []*protocol.SettleStrategy

	// 优先从结算周期记录中获取预设的策略
	if len(settleLog.GetStrategyCodes()) > 0 {
		strategyCodes := settleLog.GetStrategyCodes()
		log.Get().Debugf("GetSettleStrategiesWithPeriodCache: using preset strategy codes from settle log: %v", strategyCodes)
		strategies_result = s.GetSettleStrategiesByCodes(strategyCodes, "")
	} else {
		// 回退：查询在周期时间有效的合同
		log.Get().Debugf("GetSettleStrategiesWithPeriodCache: no preset strategies, falling back to contract query for merchant %d at time %d", settleLog.MID, settleLog.TrxStartAt)
		contract := models.GetValidContractsAtTime(strconv.FormatInt(settleLog.MID, 10), settleLog.TrxStartAt)
		if contract != nil {
			strategies_result = s.GenerateStrategiesFromContract(contract)
		}
	}

	if len(strategies_result) == 0 {
		log.Get().Warnf("GetSettleStrategiesWithPeriodCache: no settlement strategies found for settle log %s", settleLog.SettleID)
		return nil
	}

	// 存入缓存，使用周期级别的缓存key
	cache.Set(cacheKey, strategies_result)
	log.Get().Debugf("GetSettleStrategiesWithPeriodCache: cached %d strategies for settle log %s with key %s", len(strategies_result), settleLog.SettleID, cacheKey)

	return strategies_result
}

// GetSettleStrategiesWithCache 使用缓存获取结算策略（保留原有逻辑，兼容旧代码）
// 根据商户合同和交易时间确定适用的结算策略
// 注意：此函数已不推荐使用，请使用 GetSettleStrategiesWithPeriodCache
func (s *MerchantSettleService) GetSettleStrategiesWithCache(ctx context.Context, mid string, trxTime int64) []*protocol.SettleStrategy {
	if mid == "" {
		return nil
	}

	// 从 context 获取缓存和统计信息
	cache, stats := protocol.GetCacheFromContext(ctx)

	// 生成包含交易时间的缓存key
	cacheKey := fmt.Sprintf("%s_%d", mid, trxTime/86400000) // 按天缓存，避免缓存过于细粒度

	// 首先检查缓存
	strategies, exists := cache.Get(cacheKey)
	if exists {
		stats.RecordHit()
		log.Get().Debugf("GetSettleStrategiesWithCache: cache hit for merchant %s at time %d, found %d strategies", mid, trxTime, len(strategies))
		return strategies
	}

	// 缓存未命中，从数据库查询
	stats.RecordMiss()
	log.Get().Debugf("GetSettleStrategiesWithCache: cache miss for merchant %s at time %d, querying database", mid, trxTime)

	// 查询在指定时间有效的合同
	contract := models.GetValidContractsAtTime(mid, trxTime)
	if contract == nil {
		log.Get().Warnf("GetSettleStrategiesWithCache: no valid contracts found for merchant %s at time %d", mid, trxTime)
		return nil
	}

	// 根据合同配置生成结算策略
	strategies = s.GenerateStrategiesFromContract(contract)

	// 如果合同中没有配置结算策略，回退到传统方式
	if len(strategies) == 0 {
		log.Get().Warnf("GetSettleStrategiesWithCache: no settlement strategies found in contracts for merchant %s, falling back to traditional query", mid)
		return nil
	}

	// 转换为 protocol 类型并存入缓存
	cache.Set(cacheKey, strategies)
	log.Get().Debugf("GetSettleStrategiesWithCache: cached %d strategies for merchant %s at time %d", len(strategies), mid, trxTime)

	return strategies
}

// SettleResult 结算结果统计
type SettleResult struct {
	TotalTransactions   int64
	SuccessTransactions int64
	FailedTransactions  int64
	ProcessedPages      int
	TotalPages          int
	StartTime           time.Time
	EndTime             time.Time
	Duration            time.Duration
}

// SettleByTimeRange 时间范围结算
func (s *MerchantSettleService) SettleByTimeRange(ctx context.Context, trx_type string, startAt, endAt int64) *SettleResult {
	result := &SettleResult{
		StartTime: time.Now(),
	}

	// 创建包含缓存的 context
	ctx = protocol.CreateSettleContext(ctx)
	cache, stats := protocol.GetCacheFromContext(ctx)

	defer func() {
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)

		// 输出缓存使用情况
		hits, misses := stats.GetStats()
		cacheSize := cache.Size()
		hitRate := stats.GetHitRate()

		log.Get().Infof("SettleByTimeRange: completed settling %d transactions (success: %d, failed: %d, duration: %v)",
			result.TotalTransactions, result.SuccessTransactions, result.FailedTransactions, result.Duration)
		log.Get().Infof("SettleByTimeRange: cache stats - hits: %d, misses: %d, hit rate: %.2f%%, cache size: %d",
			hits, misses, hitRate, cacheSize)
	}()

	if startAt <= 0 || endAt <= 0 || startAt >= endAt {
		log.Get().Error("SettleByTimeRange: invalid time range")
		return result
	}

	log.Get().Infof("SettleByTimeRange: settling transactions from %d to %d", startAt, endAt)

	query := &models.TrxQuery{
		TrxType:        trx_type,
		SettleStatus:   protocol.StatusPending, // 只处理待结算的交易
		SettledAtStart: startAt,
		SettledAtEnd:   endAt,
	}

	// 获取总数
	total, err := models.CountTransactionByQuery(query)
	if err != nil {
		log.Get().Errorf("SettleByTimeRange: failed to count transactions: %v", err)
		return result
	}

	result.TotalTransactions = total

	if total == 0 {
		log.Get().Warn("SettleByTimeRange: no transactions found in the specified time range")
		return result
	}

	log.Get().Infof("SettleByTimeRange: found %d transactions to settle", total)

	// 分页处理配置
	const pageSize = 500
	const concurrency = 10
	const maxConcurrentPages = 5
	totalPages := int((total + int64(pageSize) - 1) / pageSize)
	result.TotalPages = totalPages

	// 并发控制
	semaphore := make(chan struct{}, concurrency)
	pageSemaphore := make(chan struct{}, maxConcurrentPages)
	var wg sync.WaitGroup
	var pageWg sync.WaitGroup

	// 用于收集结算结果的channel
	resultChan := make(chan bool, int(total))

	// 用于同步处理页面计数的mutex
	var processedPagesMutex sync.Mutex

	log.Get().Infof("SettleByTimeRange: starting settlement with %d pages", totalPages)

	// 分页查询，每个页面在独立协程中处理
	for page := 1; page <= totalPages; page++ {
		pageWg.Add(1)
		pageSemaphore <- struct{}{}

		go func(currentPage int) {
			defer func() {
				<-pageSemaphore
				pageWg.Done()
				if r := recover(); r != nil {
					log.Get().Errorf("SettleByTimeRange: panic during page %d processing: %v", currentPage, r)
				}
			}()

			pageQuery := &models.TrxQuery{
				Mid:            query.Mid,
				CreatedAtStart: query.CreatedAtStart,
				CreatedAtEnd:   query.CreatedAtEnd,
				TrxType:        query.TrxType,
				SettleStatus:   query.SettleStatus,
			}

			offset := (currentPage - 1) * pageSize
			trxs, err := models.ListTransactionByQuery(pageQuery, offset, pageSize)
			if err != nil {
				log.Get().Errorf("SettleByTimeRange: failed to list transactions for page %d: %v", currentPage, err)
				return
			}

			if len(trxs) == 0 {
				log.Get().Warnf("SettleByTimeRange: no transactions found on page %d", currentPage)
				return
			}

			processedPagesMutex.Lock()
			result.ProcessedPages++
			processedPagesMutex.Unlock()

			log.Get().Infof("SettleByTimeRange: processing page %d/%d with %d transactions", currentPage, totalPages, len(trxs))

			// 为当前页的每个交易分发线程
			for _, trx := range trxs {
				wg.Add(1)
				semaphore <- struct{}{}
				go func(transaction *models.Transaction) {
					defer func() {
						<-semaphore
						wg.Done()
					}()
					success := s.SettleTransaction(ctx, transaction)
					resultChan <- success
				}(trx)
			}
		}(page)
	}

	// 等待所有页面处理完成
	pageWg.Wait()
	// 等待所有交易处理完成
	wg.Wait()
	close(resultChan)

	// 收集结果
	for success := range resultChan {
		if success {
			result.SuccessTransactions++
		} else {
			result.FailedTransactions++
		}
	}

	return result
}

// SettleByCompletedTime 按交易完成时间进行结算
func (s *MerchantSettleService) SettleByCompletedTime(ctx context.Context, trx_type string, startAt, endAt int64) *SettleResult {
	result := &SettleResult{
		StartTime: time.Now(),
	}

	// 创建包含缓存的 context
	ctx = protocol.CreateSettleContext(ctx)
	cache, stats := protocol.GetCacheFromContext(ctx)

	defer func() {
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)

		// 输出缓存使用情况
		hits, misses := stats.GetStats()
		cacheSize := cache.Size()
		hitRate := stats.GetHitRate()

		log.Get().Infof("SettleByCompletedTime: completed settling %d transactions (success: %d, failed: %d, duration: %v)",
			result.TotalTransactions, result.SuccessTransactions, result.FailedTransactions, result.Duration)
		log.Get().Infof("SettleByCompletedTime: cache stats - hits: %d, misses: %d, hit rate: %.2f%%, cache size: %d",
			hits, misses, hitRate, cacheSize)
	}()

	if startAt <= 0 || endAt <= 0 || startAt >= endAt {
		log.Get().Error("SettleByCompletedTime: invalid time range")
		return result
	}

	log.Get().Infof("SettleByCompletedTime: settling transactions completed from %d to %d", startAt, endAt)

	query := &models.TrxQuery{
		TrxType:          trx_type,
		Status:           protocol.StatusSuccess, // 只处理成功的交易
		SettleStatus:     protocol.StatusPending, // 只处理待结算的交易
		CompletedAtStart: startAt,
		CompletedAtEnd:   endAt,
	}

	// 获取总数
	total, err := models.CountTransactionByQuery(query)
	if err != nil {
		log.Get().Errorf("SettleByCompletedTime: failed to count transactions: %v", err)
		return result
	}

	result.TotalTransactions = total

	if total == 0 {
		log.Get().Warn("SettleByCompletedTime: no transactions found in the specified time range")
		return result
	}

	log.Get().Infof("SettleByCompletedTime: found %d transactions to settle", total)

	// 分页处理配置
	const pageSize = 500
	const concurrency = 10
	const maxConcurrentPages = 5
	totalPages := int((total + int64(pageSize) - 1) / pageSize)
	result.TotalPages = totalPages

	// 并发控制
	semaphore := make(chan struct{}, concurrency)
	pageSemaphore := make(chan struct{}, maxConcurrentPages)
	var wg sync.WaitGroup
	var pageWg sync.WaitGroup

	// 用于收集结算结果的channel
	resultChan := make(chan bool, int(total))

	// 用于同步处理页面计数的mutex
	var processedPagesMutex sync.Mutex

	log.Get().Infof("SettleByCompletedTime: starting settlement with %d pages", totalPages)

	// 分页查询，每个页面在独立协程中处理
	for page := 1; page <= totalPages; page++ {
		pageWg.Add(1)
		pageSemaphore <- struct{}{}

		go func(currentPage int) {
			defer func() {
				<-pageSemaphore
				pageWg.Done()
				processedPagesMutex.Lock()
				result.ProcessedPages++
				processedPagesMutex.Unlock()
			}()

			// 分页查询交易
			offset := (currentPage - 1) * pageSize
			transactions, err := models.ListTransactionByQuery(query, offset, pageSize)
			if err != nil {
				log.Get().Errorf("SettleByCompletedTime: failed to query transactions for page %d: %v", currentPage, err)
				return
			}

			log.Get().Infof("SettleByCompletedTime: processing page %d with %d transactions", currentPage, len(transactions))

			// 为当前页的每个交易分发线程
			for _, trx := range transactions {
				wg.Add(1)
				semaphore <- struct{}{}

				go func(transaction *models.Transaction) {
					defer func() {
						<-semaphore
						wg.Done()
					}()

					// 处理单个交易的结算，使用新的流程
					success := s.SettleTransactionWithPeriod(ctx, transaction)
					resultChan <- success
				}(trx)
			}
		}(page)
	}

	// 等待所有页面处理完成
	pageWg.Wait()
	// 等待所有交易处理完成
	wg.Wait()
	close(resultChan)

	// 收集结果
	for success := range resultChan {
		if success {
			result.SuccessTransactions++
		} else {
			result.FailedTransactions++
		}
	}

	return result
}

// SettleTransaction 结算单个交易并返回是否成功
func (s *MerchantSettleService) SettleTransaction(ctx context.Context, trx *models.Transaction) (isSuccess bool) {
	isSuccess = false
	defer func() {
		if r := recover(); r != nil {
			log.Get().Errorf("SettleTransaction: panic during settlement of transaction %s: %v", trx.TrxID, r)
		}
	}()

	if trx == nil {
		log.Get().Error("SettleTransaction: transaction is nil")
		return
	}

	mid := trx.Mid
	log.Get().Infof("SettleTransaction: processing transaction %s for merchant %s", trx.TrxID, mid)

	newSettle := true
	// 0. Check if transaction already has settlement record
	settleTransaction, err := models.GetExistingSettleRecord(trx.TrxID)
	if err != nil {
		log.Get().Errorf("SettleTransaction: failed to check existing settlement for transaction %s: %v", trx.TrxID, err)
		return
	}

	if settleTransaction != nil {
		newSettle = false
	} else {
		// 创建新的结算交易记录时需要settleLogID，这个方法需要重构
		// 暂时使用空字符串，后续会修复
		settleTransaction = NewSettleTransaction(trx, "")
	}

	// Check settlement status
	status := settleTransaction.GetStatus()
	switch status {
	case protocol.StatusSuccess:
		log.Get().Infof("SettleTransaction: transaction %s already settled successfully (settle_id: %s)", trx.TrxID, settleTransaction.SettleID)
		return true
	case protocol.StatusPending:
		log.Get().Infof("SettleTransaction: transaction %s settlement is pending (settle_id: %s)", trx.TrxID, settleTransaction.SettleID)
		return false
	case protocol.StatusFailed:
		log.Get().Infof("SettleTransaction: transaction %s has failed settlement record (settle_id: %s), will retry", trx.TrxID, settleTransaction.SettleID)
	default:
		log.Get().Warnf("SettleTransaction: transaction %s has unknown settlement status %d (settle_id: %s), will retry", trx.TrxID, status, settleTransaction.SettleID)
	}

	// 使用 context 缓存获取结算策略，传入交易时间
	trxTime := trx.CreatedAt // 使用交易创建时间来确定适用的合同
	strategies := s.GetSettleStrategiesWithCache(ctx, mid, trxTime)
	if len(strategies) == 0 {
		log.Get().Warnf("SettleTransaction: no settlement strategies found for merchant %s", mid)
		return
	}

	// 2. Find matching strategy based on transaction properties
	var matchedStrategy *protocol.SettleStrategy
	for _, strategy := range strategies {
		if s.IsStrategyMatched(trx, strategy) {
			matchedStrategy = strategy
			break
		}
	}

	if matchedStrategy == nil {
		log.Get().Warnf("SettleTransaction: no matching strategy found for transaction %s", trx.TrxID)
		return
	}

	log.Get().Infof("SettleTransaction: matched strategy %d for transaction %s", matchedStrategy.ID, trx.TrxID)

	// 3. Get settlement rules from strategy
	settleRules := matchedStrategy.Rules
	if len(settleRules) == 0 {
		log.Get().Warnf("SettleTransaction: no settlement rules found in strategy %d", matchedStrategy.ID)
		return
	}

	// 4. Find the best matching rule
	var matchedRule *protocol.SettleRule
	for _, rule := range settleRules {
		if s.IsRuleApplicable(trx, rule) {
			matchedRule = rule
			break
		}
	}

	if matchedRule == nil {
		log.Get().Warnf("SettleTransaction: no applicable settlement rule found for transaction %s", trx.TrxID)
		return
	}

	log.Get().Infof("SettleTransaction: selected rule %s for transaction %s", matchedRule.RuleID, trx.TrxID)

	// 5. Calculate settlement amounts and fees
	settlementResult := s.CalculateSettlement(trx, matchedRule)
	if settlementResult == nil {
		log.Get().Errorf("SettleTransaction: failed to calculate settlement for transaction %s", trx.TrxID)
		return
	}

	settleTransaction.MerchantSettleTransactionValues = &models.MerchantSettleTransactionValues{
		SettleAmount:    &settlementResult.SettleAmount,
		SettleUsdAmount: &settlementResult.SettleUsdAmount,
		FeeCcy:          &settlementResult.FeeCcy,
		Fee:             &settlementResult.Fee,
		UsdFee:          &settlementResult.UsdFee,
		FixedFee:        &settlementResult.FixedFee,
		FixedUsdFee:     &settlementResult.FixedUsdFee,
		Rate:            &settlementResult.Rate,
		UsdRate:         &settlementResult.UsdRate,
		SettleStrategy:  matchedStrategy,
		SettleRule:      matchedRule,
	}

	err = models.WriteDB.Transaction(func(tx *gorm.DB) error {
		// 6. Create or update settlement transaction record
		if newSettle {
			// 创建新的结算记录
			if err := tx.Create(settleTransaction).Error; err != nil {
				return err
			}
		} else {
			// 更新现有的失败记录
			if err := tx.Model(settleTransaction).UpdateColumns(settleTransaction.MerchantSettleTransactionValues).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Get().Errorf("SettleTransaction: failed to settle transaction %s: %v", trx.TrxID, err)
		return
	}
	values := models.NewTrxValues().
		SetSettleID(settleTransaction.SettleID).
		SetSettledAt(settleTransaction.CreatedAt)
	// 7. Update transaction status to settled
	if err := models.SaveTransactionValues(models.WriteDB, trx, values); err != nil {
		log.Get().Errorf("SettleTransaction: failed to update transaction %s status: %v", trx.TrxID, err)
	}
	log.Get().Infof("SettleTransaction: successfully settled transaction %s with amount %s",
		trx.TrxID, settlementResult.SettleAmount.String())
	return true
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

// IsStrategyMatched 检查交易是否匹配指定的结算策略
func (s *MerchantSettleService) IsStrategyMatched(trx *models.Transaction, strategy *protocol.SettleStrategy) bool {
	if strategy == nil || trx == nil {
		return false
	}

	// 检查交易类型匹配
	if strategy.TrxType != "" && strategy.TrxType != trx.TrxType {
		return false
	}

	// 检查交易币种匹配
	if strategy.TrxCcy != "" && strategy.TrxCcy != trx.GetCcy() {
		return false
	}

	// 检查结算币种匹配
	if strategy.SettleCcy != "" && strategy.SettleCcy != trx.GetCcy() {
		return false
	}

	// 检查商户匹配 - 使用字符串比较
	if strategy.MID > 0 && strconv.FormatInt(strategy.MID, 10) != trx.Mid {
		return false
	}

	return true
}

// IsRuleApplicable 检查结算规则是否适用于指定交易
func (s *MerchantSettleService) IsRuleApplicable(trx *models.Transaction, rule *protocol.SettleRule) bool {
	if rule == nil || trx == nil {
		return false
	}

	// 检查交易类型匹配
	if rule.TrxType != "" && rule.TrxType != trx.TrxType {
		return false
	}

	// 检查币种匹配
	if rule.Ccy != "" && rule.Ccy != trx.GetCcy() {
		return false
	}

	// 检查金额范围
	if rule.MinAmount != nil && trx.Amount != nil && trx.Amount.LessThan(*rule.MinAmount) {
		return false
	}
	if rule.MaxAmount != nil && trx.Amount != nil && trx.Amount.GreaterThan(*rule.MaxAmount) {
		return false
	}

	return true
}

// CalculateSettlement 计算结算金额和费用
func (s *MerchantSettleService) CalculateSettlement(trx *models.Transaction, rule *protocol.SettleRule) *SettlementResult {
	if trx == nil || rule == nil {
		return nil
	}

	result := &SettlementResult{
		FeeCcy: trx.GetCcy(),
	}

	// 获取交易金额
	var trxAmount, trxUsdAmount decimal.Decimal
	if trx.Amount != nil {
		trxAmount = *trx.Amount
	}
	if trx.UsdAmount != nil {
		trxUsdAmount = *trx.UsdAmount
	}

	// 计算费率
	if rule.Rate != nil {
		result.Rate = *rule.Rate
		result.Fee = trxAmount.Mul(*rule.Rate).Div(decimal.NewFromInt(100))
	}

	if rule.UsdRate != nil {
		result.UsdRate = *rule.UsdRate
		result.UsdFee = trxUsdAmount.Mul(*rule.UsdRate).Div(decimal.NewFromInt(100))
	}

	// 计算固定费用
	if rule.FixedFee != nil {
		result.FixedFee = *rule.FixedFee
	}

	if rule.FixedUsdFee != nil {
		result.FixedUsdFee = *rule.FixedUsdFee
	}

	// 计算总费用
	totalFee := result.Fee.Add(result.FixedFee)
	totalUsdFee := result.UsdFee.Add(result.FixedUsdFee)

	// 计算结算金额 = 交易金额 - 费用
	result.SettleAmount = trxAmount.Sub(totalFee)
	result.SettleUsdAmount = trxUsdAmount.Sub(totalUsdFee)

	return result
}

// NewSettleTransaction 创建新的结算交易记录，挂靠到结算周期记录
func NewSettleTransaction(trx *models.Transaction, settleLogID string) *models.MerchantSettleTransaction {
	settleTransaction := &models.MerchantSettleTransaction{
		TrxID:                           trx.TrxID,
		SettleID:                        utils.GenerateSettleTrxID(),
		SettleLogID:                     &settleLogID, // 直接关联结算周期记录
		MID:                             trx.Mid,
		TrxType:                         trx.TrxType,
		TrxCcy:                          trx.GetCcy(),
		TrxAmount:                       trx.Amount,
		TrxUsdAmount:                    trx.UsdAmount,
		TrxAt:                           trx.CreatedAt,
		SettleCcy:                       trx.GetCcy(),
		MerchantSettleTransactionValues: &models.MerchantSettleTransactionValues{},
	}
	settleTransaction.SetStatus(protocol.StatusPending)
	return settleTransaction
}

// NewSettleTransactionRecord 创建包含完整信息的结算交易记录
func NewSettleTransactionRecord(trx *models.Transaction, settleLogID string, strategy *protocol.SettleStrategy, rule *protocol.SettleRule, result *SettlementResult) *models.MerchantSettleTransaction {
	settleTransaction := NewSettleTransaction(trx, settleLogID)

	settleTransaction.MerchantSettleTransactionValues = &models.MerchantSettleTransactionValues{
		SettleAmount:    &result.SettleAmount,
		SettleUsdAmount: &result.SettleUsdAmount,
		FeeCcy:          &result.FeeCcy,
		Fee:             &result.Fee,
		UsdFee:          &result.UsdFee,
		FixedFee:        &result.FixedFee,
		FixedUsdFee:     &result.FixedUsdFee,
		Rate:            &result.Rate,
		UsdRate:         &result.UsdRate,
		SettleStrategy:  strategy,
		SettleRule:      rule,
	}

	return settleTransaction
}

// GenerateStrategiesFromContract 根据合同配置生成结算策略
func (s *MerchantSettleService) GenerateStrategiesFromContract(contract *models.Contract) []*protocol.SettleStrategy {
	if contract == nil {
		return nil
	}

	// 收集所有的结算策略代码
	var payinCodes []string
	var payoutCodes []string

	codeLib := map[string]bool{}

	// 遍历合同，收集所有的策略代码
	// 收集充值结算策略代码
	if contract.SettleConfig.Payin != nil {
		for _, setting := range contract.SettleConfig.Payin {
			for _, code := range setting.Strategies {
				if !codeLib[code] {
					payinCodes = append(payinCodes, code)
					codeLib[code] = true
				}
			}
		}
	}

	// 收集提现结算策略代码
	if contract.SettleConfig.Payout != nil {
		for _, setting := range contract.SettleConfig.Payout {
			for _, code := range setting.Strategies {
				if !codeLib[code] {
					payoutCodes = append(payoutCodes, code)
					codeLib[code] = true
				}
			}
		}
	}

	// 批量查询结算策略
	strategies := []*protocol.SettleStrategy{}

	if len(payinCodes) > 0 {
		_list := models.GetSettleStrategiesByCodes(payinCodes)
		if len(_list) > 0 {
			for _, strategy := range _list {
				if strategy.TrxType == "" {
					strategy.TrxType = protocol.TrxTypePayin
				}
				strategies = append(strategies, strategy.Protocol())
			}
		}
	}

	if len(payoutCodes) > 0 {
		_list := models.GetSettleStrategiesByCodes(payoutCodes)
		if len(_list) > 0 {
			for _, strategy := range _list {
				if strategy.TrxType == "" {
					strategy.TrxType = protocol.TrxTypePayout
				}
				strategies = append(strategies, strategy.Protocol())
			}
		}
	}
	return strategies
}

// GetSettleStrategiesByCodes 根据策略代码获取结算策略
func (s *MerchantSettleService) GetSettleStrategiesByCodes(strategyCodes []string, trxType string) []*protocol.SettleStrategy {
	if len(strategyCodes) == 0 {
		return nil
	}

	strategies := []*protocol.SettleStrategy{}
	strategyList := models.GetSettleStrategiesByCodes(strategyCodes)
	if len(strategyList) > 0 {
		for _, strategy := range strategyList {
			protocolStrategy := strategy.Protocol()
			if protocolStrategy.TrxType == "" {
				protocolStrategy.TrxType = trxType
			}
			strategies = append(strategies, protocolStrategy)
		}
	}

	return strategies
}

// GetOrCreateTransactionSettleLog 获取或创建交易对应的结算周期记录
func (s *MerchantSettleService) GetOrCreateTransactionSettleLog(trx *models.Transaction) (*models.MerchantSettleLog, error) {
	mid := trx.Mid
	trxTime := trx.CreatedAt

	// 1. 获取合同配置以确定结算周期
	contract := models.GetValidContractsAtTime(mid, trxTime)
	if contract == nil {
		return nil, fmt.Errorf("no valid contracts found for merchant %s at time %d", mid, trxTime)
	}

	// 2. 获取结算配置
	var settlePeriodType string = "D1" // 默认按天结算
	var settleCcy string = trx.GetCcy()
	var settleStrategies []string

	// 从合同配置获取结算周期类型和策略
	if contract.SettleConfig != nil {
		if trx.TrxType == "payin" && len(contract.SettleConfig.Payin) > 0 {
			settlePeriodType = contract.SettleConfig.Payin[0].Type
			if contract.SettleConfig.Payin[0].Ccy != "" {
				settleCcy = contract.SettleConfig.Payin[0].Ccy
			}
			settleStrategies = contract.SettleConfig.Payin[0].Strategies
		} else if trx.TrxType == "payout" && len(contract.SettleConfig.Payout) > 0 {
			settlePeriodType = contract.SettleConfig.Payout[0].Type
			if contract.SettleConfig.Payout[0].Ccy != "" {
				settleCcy = contract.SettleConfig.Payout[0].Ccy
			}
			settleStrategies = contract.SettleConfig.Payout[0].Strategies
		}
	}

	// 3. 根据交易时间和结算周期类型计算结算周期
	period := models.GetSettlePeriodFromTime(trx.CreatedAt, settlePeriodType)

	// 4. 获取或创建结算周期记录
	settleLog, isNewLog, err := models.GetOrCreateSettleLog(mid, period, trx.TrxType, settlePeriodType, settleCcy, settleStrategies)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create settle log: %v", err)
	}

	if isNewLog {
		log.Get().Infof("GetOrCreateTransactionSettleLog: created new settle log %s for merchant %s, period %d", settleLog.SettleID, mid, period)
	} else {
		log.Get().Infof("GetOrCreateTransactionSettleLog: using existing settle log %s for merchant %s, period %d", settleLog.SettleID, mid, period)
	}

	return settleLog, nil
}

// SettleTransactionWithPeriod 使用结算周期记录处理单个交易结算
func (s *MerchantSettleService) SettleTransactionWithPeriod(ctx context.Context, trx *models.Transaction) (isSuccess bool) {
	isSuccess = false
	defer func() {
		if r := recover(); r != nil {
			log.Get().Errorf("SettleTransactionWithPeriod: panic while settling transaction %s: %v", trx.TrxID, r)
			isSuccess = false
		}
	}()

	if trx == nil {
		log.Get().Error("SettleTransactionWithPeriod: transaction is nil")
		return
	}

	mid := trx.Mid
	log.Get().Infof("SettleTransactionWithPeriod: processing transaction %s for merchant %s", trx.TrxID, mid)

	// 1. 获取或创建结算周期记录（最优先级）
	settleLog, err := s.GetOrCreateTransactionSettleLog(trx)
	if err != nil {
		log.Get().Errorf("SettleTransactionWithPeriod: %v", err)
		return
	}

	// 2. 使用周期级别缓存获取结算策略
	strategies := s.GetSettleStrategiesWithPeriodCache(ctx, settleLog)

	if len(strategies) == 0 {
		log.Get().Warnf("SettleTransactionWithPeriod: no settlement strategies found for transaction %s", trx.TrxID)
		return
	}

	// 3. 找到匹配的策略
	var matchedStrategy *protocol.SettleStrategy
	for _, strategy := range strategies {
		if s.IsStrategyMatched(trx, strategy) {
			matchedStrategy = strategy
			break
		}
	}

	if matchedStrategy == nil {
		log.Get().Warnf("SettleTransactionWithPeriod: no matching strategy found for transaction %s", trx.TrxID)
		return
	}

	log.Get().Infof("SettleTransactionWithPeriod: matched strategy %d for transaction %s", matchedStrategy.ID, trx.TrxID)

	// 4. 获取结算规则
	settleRules := matchedStrategy.Rules
	if len(settleRules) == 0 {
		log.Get().Warnf("SettleTransactionWithPeriod: no settlement rules found in strategy %d", matchedStrategy.ID)
		return
	}

	// 5. 找到适用的规则
	var matchedRule *protocol.SettleRule
	for _, rule := range settleRules {
		if s.IsRuleApplicable(trx, rule) {
			matchedRule = rule
			break
		}
	}

	if matchedRule == nil {
		log.Get().Warnf("SettleTransactionWithPeriod: no applicable settlement rule found for transaction %s", trx.TrxID)
		return
	}

	log.Get().Infof("SettleTransactionWithPeriod: selected rule %s for transaction %s", matchedRule.RuleID, trx.TrxID)

	// 6. 计算结算金额和手续费
	settlementResult := s.CalculateSettlement(trx, matchedRule)
	if settlementResult == nil {
		log.Get().Errorf("SettleTransactionWithPeriod: failed to calculate settlement for transaction %s", trx.TrxID)
		return
	}

	// 7. 检查交易是否已有结算记录
	existingSettleRecord, err := models.GetExistingSettleRecord(trx.TrxID)
	if err != nil {
		log.Get().Errorf("SettleTransactionWithPeriod: failed to check existing settlement for transaction %s: %v", trx.TrxID, err)
		return
	}

	newSettle := existingSettleRecord == nil
	var settleTransaction *models.MerchantSettleTransaction

	if newSettle {
		// 8. 创建新的结算交易记录，直接关联结算周期记录
		settleTransaction = NewSettleTransaction(trx, settleLog.SettleID)
	} else {
		settleTransaction = existingSettleRecord
		// 更新关联的结算周期记录
		settleTransaction.SetSettleLogID(settleLog.SettleID)
	}

	// 9. 设置结算计算结果
	settleTransaction.MerchantSettleTransactionValues = &models.MerchantSettleTransactionValues{
		SettleAmount:    &settlementResult.SettleAmount,
		SettleUsdAmount: &settlementResult.SettleUsdAmount,
		FeeCcy:          &settlementResult.FeeCcy,
		Fee:             &settlementResult.Fee,
		UsdFee:          &settlementResult.UsdFee,
		FixedFee:        &settlementResult.FixedFee,
		FixedUsdFee:     &settlementResult.FixedUsdFee,
		Rate:            &settlementResult.Rate,
		UsdRate:         &settlementResult.UsdRate,
		SettleStrategy:  matchedStrategy,
		SettleRule:      matchedRule,
	}
	settleTransaction.SetStatus(protocol.StatusSuccess)
	settleTransaction.SetSettledAt(time.Now().UnixMilli())

	// 10. 在事务中保存结算记录和更新交易状态
	err = models.WriteDB.Transaction(func(tx *gorm.DB) error {
		// 保存或更新结算交易记录
		if newSettle {
			if err := tx.Create(settleTransaction).Error; err != nil {
				return fmt.Errorf("failed to create settlement transaction: %v", err)
			}
		} else {
			if err := tx.Save(settleTransaction).Error; err != nil {
				return fmt.Errorf("failed to update settlement transaction: %v", err)
			}
		}

		// 更新结算周期记录的统计信息
		if err := models.UpdateSettleLogWithTransaction(settleLog.SettleID, trx,
			settlementResult.SettleAmount, settlementResult.SettleUsdAmount,
			settlementResult.Fee, settlementResult.UsdFee); err != nil {
			return fmt.Errorf("failed to update settle log: %v", err)
		}

		// 更新交易状态
		values := models.NewTrxValues().
			SetSettleID(settleTransaction.SettleID).
			SetSettledAt(settleTransaction.GetSettledAt()).
			SetSettleStatus(protocol.StatusSuccess)

		if err := models.SaveTransactionValues(tx, trx, values); err != nil {
			return fmt.Errorf("failed to update transaction status: %v", err)
		}

		return nil
	})

	if err != nil {
		log.Get().Errorf("SettleTransactionWithPeriod: failed to settle transaction %s: %v", trx.TrxID, err)
		return
	}

	log.Get().Infof("SettleTransactionWithPeriod: successfully settled transaction %s with amount %s, settle_log_id: %s",
		trx.TrxID, settlementResult.SettleAmount.String(), settleLog.SettleID)

	return true
}
