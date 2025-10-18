package services

import (
	"fmt"
	"inpayos/internal/config"
	"inpayos/internal/log"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
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

// SettleByTimeRange 按完成时间范围进行结算
func (s *MerchantSettleService) SettleByTimeRange(trx_type string, startAt, endAt int64) *protocol.SettleResult {
	result := &protocol.SettleResult{
		StartTime: time.Now(),
	}
	// 创建包含缓存的 context
	ctx := protocol.NewMerchantSettleContext()
	defer func() {
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)

		log.Get().Infof("SettleByCompletedTime: completed settling %d transactions (success: %d, failed: %d, duration: %v)",
			result.TotalCount, result.SuccessCount, result.FailedCount, result.Duration)
	}()

	if startAt <= 0 || endAt <= 0 || startAt >= endAt {
		result.Result = "invalid time range"
		return result
	}

	log.Get().Infof("SettleByCompletedTime: settling transactions completed from %d to %d", startAt, endAt)

	query := &models.TrxQuery{
		TrxType:          trx_type,
		Status:           protocol.StatusSuccess, // 只处理成功的交易
		SettleStatus:     protocol.StatusPending, // 只处理待结算的交易
		CompletedAtStart: startAt,                //以完成时间为准
		CompletedAtEnd:   endAt,
	}

	// 获取总数
	total := models.CountTransactionByQuery(query)
	if total == 0 {
		result.Result = "no transactions to settle"
		return result
	}

	result.TotalCount = total

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
	for isSuccess := range resultChan {
		if isSuccess {
			result.SuccessCount++
		} else {
			result.FailedCount++
		}
	}

	return result
}

// SettleTransaction 处理单个交易结算
func (s *MerchantSettleService) SettleTransaction(ctx *protocol.MerchantSettleContext, trx *models.Transaction) (isSuccess bool) {
	isSuccess = false
	defer func() {
		if r := recover(); r != nil {
			log.Get().Errorf("SettleTransactionWithPeriod: panic while settling transaction %s: %v", trx.TrxID, r)
		}
	}()

	if trx == nil {
		log.Get().Error("SettleTransactionWithPeriod: transaction is nil")
		return
	}

	mid := trx.Mid
	log.Get().Infof("SettleTransactionWithPeriod: processing transaction %s for merchant %s", trx.TrxID, mid)

	// 1. 获取或创建结算周期记录（最优先级）
	settleLog, err := s.GetOrCreateTransactionSettleLog(trx, ctx.SettledAt)
	if err != nil {
		log.Get().Errorf("SettleTransactionWithPeriod: %v", err)
		return
	}

	// 2. 使用周期级别缓存获取结算策略
	strategies := s.GetSettleStrategies(ctx, settleLog)

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

	// 7. 检查交易是否已有结算记录
	settleTransaction := models.GetExistingSettleRecord(trx.TrxID)
	newSettle := false
	updateSettle := false
	if settleTransaction == nil {
		// 8. 创建新的结算交易记录，直接关联结算周期记录
		settleTransaction = models.NewSettleTransaction(trx, settleLog.SettleID)
		newSettle = true
	} else if settleTransaction.GetStatus() != protocol.StatusSuccess {
		// 8. 只有当结算状态为非成功的，才进行更新
		updateSettle = true
		settleTransaction.SetSettleLogID(settleLog.SettleID)
	} else {
		isSuccess = true
		return // 已成功结算，无需重复处理
	}
	// 6. 计算结算金额和手续费
	settlementResult := s.CalculateSettlement(trx, matchedRule)
	if settlementResult == nil {
		log.Get().Errorf("SettleTransactionWithPeriod: failed to calculate settlement for transaction %s", trx.TrxID)
		return
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
	settleTransaction.SetStatus(protocol.StatusSuccess).
		SetSettledAt(utils.TimeNowMilli())

	// 10. 在事务中保存结算记录和更新交易状态
	err = models.WriteDB.Transaction(func(tx *gorm.DB) error {
		// 保存或更新结算交易记录
		if newSettle {
			if err := tx.Create(settleTransaction).Error; err != nil {
				return err
			}
		}
		if updateSettle {
			if err := tx.Model(settleTransaction).UpdateColumns(settleTransaction.MerchantSettleTransactionValues).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if newSettle || updateSettle {
		// 更新交易状态
		values := models.NewTrxValues().
			SetSettleID(settleLog.SettleID).
			SetSettledAt(settleTransaction.GetSettledAt()).
			SetSettleStatus(protocol.StatusSuccess)
		if err := models.SaveTransactionValues(models.WriteDB, trx, values); err != nil {
			log.Get().Errorf("SettleTransactionWithPeriod: failed to update transaction %s: %v", trx.TrxID, err)
		}
	}
	if err != nil {
		log.Get().Errorf("SettleTransactionWithPeriod: failed to settle transaction %s: %v", trx.TrxID, err)
		return
	}
	log.Get().Infof("SettleTransactionWithPeriod: successfully settled transaction %s with amount %s, settle_log_id: %s",
		trx.TrxID, settlementResult.SettleAmount.String(), settleLog.SettleID)

	return true
}

// GetSettleStrategies 基于结算周期记录使用缓存获取结算策略
func (s *MerchantSettleService) GetSettleStrategies(ctx *protocol.MerchantSettleContext, settleLog *models.MerchantSettle) []*protocol.SettleStrategy {
	if settleLog == nil {
		return nil
	}

	// 使用结算周期记录的唯一标识作为缓存key
	// 格式: settle.period_{商户ID}_{周期}
	cacheKey := fmt.Sprintf("settle.period_%v_%v_%v",
		settleLog.Mid,
		settleLog.PeriodType,
		settleLog.Period)

	// 首先检查缓存
	strategies := ctx.Get(cacheKey)
	if len(strategies) > 0 {
		log.Get().Debugf("GetSettleStrategiesWithPeriodCache: cache hit for settle log %s, found %d strategies", settleLog.SettleID, len(strategies))
		return strategies
	}
	log.Get().Debugf("GetSettleStrategiesWithPeriodCache: cache miss for settle log %s, retrieving strategies", settleLog.SettleID)

	// 优先从结算周期记录中获取预设的策略
	if len(settleLog.GetStrategyCodes()) > 0 {
		strategyCodes := settleLog.GetStrategyCodes()
		log.Get().Debugf("GetSettleStrategiesWithPeriodCache: using preset strategy codes from settle log: %v", strategyCodes)
		strategies = s.GetSettleStrategiesByCodes(strategyCodes)
	} else {
		// 回退：查询在周期时间有效的合同
		log.Get().Debugf("GetSettleStrategiesWithPeriodCache: no preset strategies, falling back to contract query for merchant %d at time %d", settleLog.Mid, settleLog.TrxStartAt)
		contract := models.GetValidContractsAtTime(settleLog.Mid, settleLog.TrxStartAt)
		if contract != nil {
			strategies = s.GenerateStrategiesFromContract(contract)
		}
	}

	if len(strategies) == 0 {
		log.Get().Warnf("GetSettleStrategiesWithPeriodCache: no settlement strategies found for settle log %s", settleLog.SettleID)
		return nil
	}

	// 存入缓存，使用周期级别的缓存key
	ctx.Set(cacheKey, strategies)
	log.Get().Debugf("GetSettleStrategiesWithPeriodCache: cached %d strategies for settle log %s with key %s", len(strategies), settleLog.SettleID, cacheKey)

	return strategies
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
	if strategy.TrxCcy != "" && strategy.TrxCcy != trx.Ccy {
		return false
	}

	// 检查结算币种匹配
	if strategy.SettleCcy != "" && strategy.SettleCcy != trx.Ccy {
		return false
	}

	// 检查商户匹配 - 使用字符串比较
	if strategy.Mid != trx.Mid {
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
	if rule.Ccy != "" && rule.Ccy != trx.Ccy {
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
func (s *MerchantSettleService) CalculateSettlement(trx *models.Transaction, rule *protocol.SettleRule) *protocol.SettlementResult {
	if trx == nil || rule == nil {
		return nil
	}

	result := &protocol.SettlementResult{
		FeeCcy: trx.Ccy,
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

// GenerateStrategiesFromContract 根据合同配置生成结算策略
func (s *MerchantSettleService) GenerateStrategiesFromContract(contract *models.MerchantContract) []*protocol.SettleStrategy {
	if contract == nil {
		return nil
	}

	// 收集所有的结算策略代码
	var payinCodes []string
	var payoutCodes []string

	codeLib := map[string]bool{}

	// 遍历合同，收集所有的策略代码
	// 收集充值结算策略代码
	if contract.Payin != nil {
		for _, code := range contract.Payin.Settle.Strategies {
			if !codeLib[code] {
				payinCodes = append(payinCodes, code)
				codeLib[code] = true
			}
		}
	}

	// 收集提现结算策略代码
	if contract.Payout != nil {
		for _, code := range contract.Payout.Settle.Strategies {
			if !codeLib[code] {
				payoutCodes = append(payoutCodes, code)
				codeLib[code] = true
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
func (s *MerchantSettleService) GetSettleStrategiesByCodes(strategyCodes []string) []*protocol.SettleStrategy {
	if len(strategyCodes) == 0 {
		return nil
	}

	strategies := []*protocol.SettleStrategy{}
	strategyList := models.GetSettleStrategiesByCodes(strategyCodes)
	if len(strategyList) > 0 {
		for _, item := range strategyList {
			strategy := item.Protocol()
			strategies = append(strategies, strategy)
		}
	}

	return strategies
}

// GetOrCreateTransactionSettleLog 获取或创建交易对应的结算周期记录
func (s *MerchantSettleService) GetOrCreateTransactionSettleLog(trx *models.Transaction, settleAt int64) (*models.MerchantSettle, error) {
	mid := trx.Mid
	trxTime := trx.CreatedAt

	// 1. 以交易时间获得有效合同
	contract := models.GetValidContractsAtTime(mid, trxTime)
	if contract == nil {
		return nil, fmt.Errorf("no valid contracts found for merchant %s at time %d", mid, trxTime)
	}

	// 2. 从合同配置获取结算周期类型和策略
	settleConfig := contract.GetSettleConfig(trx.TrxType)
	if settleConfig == nil {
		return nil, fmt.Errorf("invalid contract settle config for merchant %s", mid)
	}

	// 4. 获取或创建结算周期记录
	settleLog := models.GetOrCreateSettleLog(mid, trx.GetCompletedAt(), settleAt, settleConfig)
	if settleLog == nil {
		return nil, fmt.Errorf("failed to get or create settle log")
	}

	return settleLog, nil
}

// ProcessSettleLogAccounting 处理单个结算记录的记账操作
func (s *MerchantSettleService) ProcessSettleLogAccounting(settleLog *models.MerchantSettle) error {
	if settleLog == nil {
		return fmt.Errorf("settle log is nil")
	}

	// 获取结算金额
	settleAmount := settleLog.GetSettleAmount()
	if settleAmount.IsZero() {
		log.Get().Warnf("ProcessSettleLogAccounting: settle amount is zero for settle_id %s", settleLog.SettleID)
		return nil
	}

	// 调用账户服务更新余额
	accountService := GetAccountService()
	balanceReq := &protocol.UpdateBalanceRequest{
		UserID:      cast.ToString(settleLog.Mid),
		UserType:    protocol.UserTypeMerchant,
		Ccy:         settleLog.SettleCcy,
		TrxType:     protocol.TrxTypeDeposit, // 结算入账使用充值类型
		Amount:      settleAmount,
		TrxID:       settleLog.SettleID, // 使用结算ID作为交易ID
		Description: fmt.Sprintf("结算周期记账，周期: %d, 类型: %s", settleLog.Period, settleLog.PeriodType),
	}

	errCode := accountService.UpdateBalance(balanceReq)
	if errCode != protocol.Success {
		return fmt.Errorf("failed to update merchant balance for settle_id %s: %s", settleLog.SettleID, errCode)
	}

	// 更新结算记录的完成时间，标记为已记账
	err := models.WriteDB.Model(settleLog).Updates(map[string]interface{}{
		"completed_at": time.Now().UnixMilli(),
		"updated_at":   time.Now().UnixMilli(),
	}).Error

	if err != nil {
		return fmt.Errorf("failed to update settle log completed_at for settle_id %s: %v", settleLog.SettleID, err)
	}

	log.Get().Infof("ProcessSettleLogAccounting: successfully processed accounting for settle_id %s, merchant %d, amount %s %s",
		settleLog.SettleID, settleLog.Mid, settleAmount.String(), settleLog.SettleCcy)

	return nil
}

// ProcessBatchSettleAccounting 批量处理结算记账
func (s *MerchantSettleService) ProcessBatchSettleAccounting(currentTime int64) *protocol.SettleAccountingResult {
	result := &protocol.SettleAccountingResult{
		StartTime: time.Now(),
	}

	defer func() {
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
	}()

	// 获取待记账的结算记录
	settleLogs, err := models.GetPendingAccountingSettleLogs(currentTime)
	if err != nil {

		return result
	}

	result.TotalCount = int64(len(settleLogs))

	if result.TotalCount == 0 {
		result.Result = "no pending accounting settle logs"
		return result
	}

	log.Get().Infof("ProcessBatchSettleAccounting: found %d pending accounting settle logs", result.TotalCount)

	// 逐个处理记账
	for _, settleLog := range settleLogs {
		err := s.ProcessSettleLogAccounting(settleLog)
		if err != nil {
			result.FailedCount++
			continue
		}
		result.SuccessCount++
	}

	log.Get().Infof("ProcessBatchSettleAccounting: completed processing %d settle logs (success: %d, failed: %d, duration: %v)",
		result.TotalCount, result.SuccessCount, result.FailedCount, result.Duration)

	return result
}

// FixTransactionSettleID 修复交易记录的settle_id字段
// 对于结算交易记录中对应的交易记录settle_id为空的情况，更新为对应的settle_log_id
// 使用并发处理，最多5个批次，每批次最多10笔交易
func (s *MerchantSettleService) FixTransactionSettleID(startTime, endTime int64) *protocol.SettleFixResult {
	result := &protocol.SettleFixResult{
		StartTime: time.Now(),
	}

	defer func() {
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
	}()

	// 获取总记录数
	count := models.CountSettleTransactionsByTimeRange(startTime, endTime)
	if count == 0 {
		result.Result = "no settle transactions found in time range"
		return result
	}

	result.TotalCount = count
	log.Get().Infof("FixTransactionSettleID: found %d settle transactions to process", result.TotalCount)

	// 并发控制参数
	const maxBatches = 5                               // 最多5个批次同时进行
	const batchSize = 10                               // 每批次最多10笔交易
	const transactionsPerPage = maxBatches * batchSize // 每页处理的交易数量

	// 计算总页数
	totalPages := int((count + int64(transactionsPerPage) - 1) / int64(transactionsPerPage))

	// 用于收集处理结果
	var resultMutex sync.Mutex

	// 逐页处理
	for page := 0; page < totalPages; page++ {
		offset := page * transactionsPerPage
		limit := transactionsPerPage

		// 获取当前页的结算交易记录（只查询必要字段）
		settleTransactions := models.ListSettleTransactionsWithLimitedFields(startTime, endTime, offset, limit)
		if len(settleTransactions) == 0 {
			continue
		}

		log.Get().Infof("FixTransactionSettleID: processing page %d with %d transactions", page+1, len(settleTransactions))

		// 将当前页的交易分批处理
		batches := s.createBatches(settleTransactions, batchSize)

		// 使用WaitGroup控制批次并发
		var batchWG sync.WaitGroup
		batchSemaphore := make(chan struct{}, maxBatches) // 限制并发批次数

		for batchIndex, batch := range batches {
			batchWG.Add(1)
			batchSemaphore <- struct{}{} // 获取批次信号量

			go func(batchNum int, transactions []*models.MerchantSettleTransaction) {
				defer func() {
					<-batchSemaphore // 释放批次信号量
					batchWG.Done()
				}()

				log.Get().Debugf("FixTransactionSettleID: starting batch %d with %d transactions", batchNum+1, len(transactions))

				batchSuccessCount := int64(0)
				batchFailedCount := int64(0)

				// 在批次内并发处理每个交易
				var transactionWG sync.WaitGroup
				transactionSemaphore := make(chan struct{}, batchSize) // 限制批次内并发数

				for _, settleTransaction := range transactions {
					transactionWG.Add(1)
					transactionSemaphore <- struct{}{}

					go func(st *models.MerchantSettleTransaction) {
						defer func() {
							<-transactionSemaphore
							transactionWG.Done()
						}()

						err := s.ProcessTransactionSettleIDFix(st)
						if err != nil {
							log.Get().Errorf("FixTransactionSettleID: failed to process transaction %s: %v", st.TrxID, err)
							batchFailedCount++
						} else {
							batchSuccessCount++
						}
					}(settleTransaction)
				}

				// 等待批次内所有交易处理完成
				transactionWG.Wait()

				// 更新总体结果（加锁保护）
				resultMutex.Lock()
				result.SuccessCount += batchSuccessCount
				result.FailedCount += batchFailedCount
				resultMutex.Unlock()

				log.Get().Debugf("FixTransactionSettleID: completed batch %d (success: %d, failed: %d)",
					batchNum+1, batchSuccessCount, batchFailedCount)

			}(batchIndex, batch)
		}

		// 等待当前页的所有批次完成
		batchWG.Wait()
		log.Get().Infof("FixTransactionSettleID: completed page %d", page+1)
	}

	log.Get().Infof("FixTransactionSettleID: completed processing %d transactions (success: %d, failed: %d, duration: %v)",
		result.TotalCount, result.SuccessCount, result.FailedCount, result.Duration)

	return result
}

// createBatches 将交易列表分割成批次
func (s *MerchantSettleService) createBatches(transactions []*models.MerchantSettleTransaction, batchSize int) [][]*models.MerchantSettleTransaction {
	var batches [][]*models.MerchantSettleTransaction

	for i := 0; i < len(transactions); i += batchSize {
		end := min(i+batchSize, len(transactions))
		batches = append(batches, transactions[i:end])
	}

	return batches
}

// ProcessTransactionSettleIDFix 处理单个交易记录的settle_id修复
func (s *MerchantSettleService) ProcessTransactionSettleIDFix(settleTransaction *models.MerchantSettleTransaction) error {
	if settleTransaction == nil {
		return fmt.Errorf("settle transaction is nil")
	}

	// 更新交易记录的settle_id为对应的settle_log_id
	if settleTransaction.GetSettleLogID() == "" {
		log.Get().Warnf("ProcessTransactionSettleIDFix: settle_log_id is empty for transaction %s", settleTransaction.TrxID)
		return nil
	}

	// 查询对应的交易记录
	var transaction models.Transaction
	err := models.ReadDB.Where("trx_id = ? && (settle_id is null OR settle_id = '')", settleTransaction.TrxID).First(&transaction).Error
	if err != nil || transaction.GetSettleID() != "" {
		return nil // 交易不存在或不需要更新
	}

	// 更新交易的settle_id
	values := models.NewTrxValues().SetSettleID(settleTransaction.GetSettleLogID())
	err = models.SaveTransactionValues(models.WriteDB, &transaction, values)
	if err != nil {
		return fmt.Errorf("failed to update transaction settle_id for %s: %v", settleTransaction.TrxID, err)
	}

	log.Get().Infof("ProcessTransactionSettleIDFix: successfully updated settle_id for transaction %s to %s",
		settleTransaction.TrxID, settleTransaction.GetSettleLogID())

	return nil
}
