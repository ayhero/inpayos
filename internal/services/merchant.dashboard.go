package services

import (
	"errors"
	"fmt"
	"inpayos/internal/models"
	"time"
)

// DashboardTodayStats 今日统计数据
type DashboardTodayStats struct {
	TodayCollection     string `json:"today_collection"`      // 今日代收金额
	TodayCollectionRate string `json:"today_collection_rate"` // 今日代收增长率
	TodayPayout         string `json:"today_payout"`          // 今日代付金额
	TodayPayoutRate     string `json:"today_payout_rate"`     // 今日代付增长率
	SuccessRate         string `json:"success_rate"`          // 成功率
	SuccessRateChange   string `json:"success_rate_change"`   // 成功率变化
}

// DashboardTransactionTrend 交易趋势数据
type DashboardTransactionTrend struct {
	Date       string `json:"date"`       // 日期 (格式: MM-DD)
	Collection int64  `json:"collection"` // 代收金额
	Payout     int64  `json:"payout"`     // 代付金额
}

// DashboardSettlementTrend 结算趋势数据
type DashboardSettlementTrend struct {
	Week   string  `json:"week"`   // 周期 (第X周)
	Amount float64 `json:"amount"` // 结算金额
}

// DashboardAccountBalance 账户余额
type DashboardAccountBalance struct {
	Currency     string `json:"currency"`      // 币种
	Balance      string `json:"balance"`       // 余额
	FrozenAmt    string `json:"frozen_amt"`    // 冻结金额
	AvailableAmt string `json:"available_amt"` // 可用金额
}

// DashboardOverview Dashboard概览数据
type DashboardOverview struct {
	TodayStats       *DashboardTodayStats        `json:"today_stats"`       // 今日统计
	TransactionTrend []DashboardTransactionTrend `json:"transaction_trend"` // 交易趋势
}

// GetTodayStats 获取今日统计数据
func GetTodayStats(merchantID int64, currency string) (*DashboardTodayStats, error) {
	if merchantID == 0 {
		return nil, errors.New("invalid merchant ID")
	}

	fmt.Printf("GetTodayStats called with merchantID: %d\n", merchantID)

	db := models.GetDB()

	// 检查数据库中是否有任何MerchantPayin记录
	var totalPayinCount int64
	db.Model(&models.MerchantPayin{}).Count(&totalPayinCount)
	fmt.Printf("Total payins in database: %d\n", totalPayinCount)

	// 检查数据库中是否有任何MerchantPayout记录
	var totalPayoutCount int64
	db.Model(&models.MerchantPayout{}).Count(&totalPayoutCount)
	fmt.Printf("Total payouts in database: %d\n", totalPayoutCount)

	// 首先获取商户的Mid
	merchant := &models.Merchant{}
	err := db.Where("id = ?", merchantID).First(merchant).Error
	if err != nil {
		return nil, fmt.Errorf("merchant not found: %v", err)
	}
	fmt.Printf("Found merchant: ID=%d, Mid=%s\n", merchant.ID, merchant.Mid)

	// 检查特定商户的记录
	var merchantPayinCount int64
	db.Model(&models.MerchantPayin{}).Where("mid = ?", merchant.Mid).Count(&merchantPayinCount)
	fmt.Printf("Payins for merchant %s: %d\n", merchant.Mid, merchantPayinCount)

	var merchantPayoutCount int64
	db.Model(&models.MerchantPayout{}).Where("mid = ?", merchant.Mid).Count(&merchantPayoutCount)
	fmt.Printf("Payouts for merchant %s: %d\n", merchant.Mid, merchantPayoutCount)

	// 计算今日和昨日的时间范围（毫秒时间戳）
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UnixMilli()
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999, now.Location()).UnixMilli()
	yesterdayStart := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location()).UnixMilli()
	yesterdayEnd := time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 999, now.Location()).UnixMilli()

	// 获取今日代收数据
	var todayCollection, yesterdayCollection float64
	var todayCollectionCount int64

	// 今日代收 - 先不限制时间，查询所有数据
	fmt.Printf("Querying today collection: mid=%s, start=%d, end=%d\n", merchant.Mid, todayStart, todayEnd)

	// 临时：查询所有数据不限时间
	result := db.Model(&models.MerchantPayin{}).
		Where("mid = ?", merchant.Mid).
		Select("COALESCE(SUM(amount), 0) as total, COUNT(*) as count").
		Row()
	scanErr := result.Scan(&todayCollection, &todayCollectionCount)
	if scanErr != nil {
		fmt.Printf("Error scanning today collection: %v\n", scanErr)
	}
	fmt.Printf("Today collection result (all time): total=%.2f, count=%d\n", todayCollection, todayCollectionCount)

	// 查询状态分布
	var statusCount int64
	db.Model(&models.MerchantPayin{}).Where("mid = ? AND status = ?", merchant.Mid, "success").Count(&statusCount)
	fmt.Printf("Success status payins: %d\n", statusCount)

	// 查询最新的几条记录看时间戳
	var recentPayins []models.MerchantPayin
	db.Where("mid = ?", merchant.Mid).Order("created_at DESC").Limit(3).Find(&recentPayins)
	for i, payin := range recentPayins {
		fmt.Printf("Recent payin %d: ID=%d, Amount=%.2f, Status=%s, CreatedAt=%d\n",
			i+1, payin.ID, payin.Amount, payin.Status, payin.CreatedAt)
	}

	// 昨日代收 - 临时设为0，专注今日数据
	yesterdayCollection = 0

	// 获取今日代付数据
	var todayPayout, yesterdayPayout float64
	var todayPayoutCount int64

	// 今日代付 - 查询所有数据不限时间
	db.Model(&models.MerchantPayout{}).
		Where("mid = ?", merchant.Mid).
		Select("COALESCE(SUM(amount), 0) as total, COUNT(*) as count").
		Row().Scan(&todayPayout, &todayPayoutCount)
	fmt.Printf("Today payout result (all time): total=%.2f, count=%d\n", todayPayout, todayPayoutCount)

	// 昨日代付 - 临时设为0
	yesterdayPayout = 0

	// 计算成功率
	var todayTotal, todaySuccess int64

	// 今日总交易数 = 代收总数 + 代付总数
	var todayPayinTotal, todayPayoutTotal int64
	db.Model(&models.MerchantPayin{}).
		Where("mid = ? AND created_at >= ? AND created_at <= ?", merchant.Mid, todayStart, todayEnd).
		Count(&todayPayinTotal)
	db.Model(&models.MerchantPayout{}).
		Where("mid = ? AND created_at >= ? AND created_at <= ?", merchant.Mid, todayStart, todayEnd).
		Count(&todayPayoutTotal)
	todayTotal = todayPayinTotal + todayPayoutTotal

	// 今日成功交易数 = 代收成功数 + 代付成功数
	var todayPayinSuccess, todayPayoutSuccess int64
	db.Model(&models.MerchantPayin{}).
		Where("mid = ? AND created_at >= ? AND created_at <= ? AND status = ?", merchant.Mid, todayStart, todayEnd, "success").
		Count(&todayPayinSuccess)
	db.Model(&models.MerchantPayout{}).
		Where("mid = ? AND created_at >= ? AND created_at <= ? AND status = ?", merchant.Mid, todayStart, todayEnd, "success").
		Count(&todayPayoutSuccess)
	todaySuccess = todayPayinSuccess + todayPayoutSuccess

	var successRate float64
	if todayTotal > 0 {
		successRate = float64(todaySuccess) / float64(todayTotal) * 100
	}

	// 计算昨日成功率用于对比
	var yesterdayTotal, yesterdaySuccess int64

	// 昨日总交易数
	var yesterdayPayinTotal, yesterdayPayoutTotal int64
	db.Model(&models.MerchantPayin{}).
		Where("mid = ? AND created_at >= ? AND created_at <= ?", merchant.Mid, yesterdayStart, yesterdayEnd).
		Count(&yesterdayPayinTotal)
	db.Model(&models.MerchantPayout{}).
		Where("mid = ? AND created_at >= ? AND created_at <= ?", merchant.Mid, yesterdayStart, yesterdayEnd).
		Count(&yesterdayPayoutTotal)
	yesterdayTotal = yesterdayPayinTotal + yesterdayPayoutTotal

	// 昨日成功交易数
	var yesterdayPayinSuccess, yesterdayPayoutSuccess int64
	db.Model(&models.MerchantPayin{}).
		Where("mid = ? AND created_at >= ? AND created_at <= ? AND status = ?", merchant.Mid, yesterdayStart, yesterdayEnd, "success").
		Count(&yesterdayPayinSuccess)
	db.Model(&models.MerchantPayout{}).
		Where("mid = ? AND created_at >= ? AND created_at <= ? AND status = ?", merchant.Mid, yesterdayStart, yesterdayEnd, "success").
		Count(&yesterdayPayoutSuccess)
	yesterdaySuccess = yesterdayPayinSuccess + yesterdayPayoutSuccess

	var yesterdaySuccessRate float64
	if yesterdayTotal > 0 {
		yesterdaySuccessRate = float64(yesterdaySuccess) / float64(yesterdayTotal) * 100
	}

	// 计算增长率
	collectionRate := calculateRate(todayCollection, yesterdayCollection)
	payoutRate := calculateRate(todayPayout, yesterdayPayout)

	return &DashboardTodayStats{
		TodayCollection:     fmt.Sprintf("%.2f", todayCollection),
		TodayCollectionRate: collectionRate,
		TodayPayout:         fmt.Sprintf("%.2f", todayPayout),
		TodayPayoutRate:     payoutRate,
		SuccessRate:         fmt.Sprintf("%.1f", successRate),
		SuccessRateChange:   fmt.Sprintf("%.1f", successRate-yesterdaySuccessRate),
	}, nil
}

// GetTransactionTrend 获取交易趋势数据
func GetTransactionTrend(merchantID int64, days int) ([]DashboardTransactionTrend, error) {
	if merchantID == 0 {
		return nil, errors.New("invalid merchant ID")
	}

	db := models.GetDB()

	// 获取商户的Mid
	merchant := &models.Merchant{}
	err := db.Where("id = ?", merchantID).First(merchant).Error
	if err != nil {
		return nil, fmt.Errorf("merchant not found: %v", err)
	}

	var trends []DashboardTransactionTrend

	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		displayDate := date.Format("01-02")

		var collection, payout float64

		// 临时：获取所有代收数据（不限时间）
		if i == 0 { // 只在第一天显示数据，其他天为0
			db.Model(&models.MerchantPayin{}).
				Where("mid = ?", merchant.Mid).
				Select("COALESCE(SUM(amount), 0)").
				Row().Scan(&collection)

			// 获取所有代付数据（不限时间）
			db.Model(&models.MerchantPayout{}).
				Where("mid = ?", merchant.Mid).
				Select("COALESCE(SUM(amount), 0)").
				Row().Scan(&payout)
		} else {
			collection = 0
			payout = 0
		}

		trends = append(trends, DashboardTransactionTrend{
			Date:       displayDate,
			Collection: int64(collection),
			Payout:     int64(payout),
		})
	}

	return trends, nil
}

// GetSettlementTrend 获取结算趋势数据
func GetSettlementTrend(merchantID int64, weeks int) ([]DashboardSettlementTrend, error) {
	if merchantID == 0 {
		return nil, errors.New("invalid merchant ID")
	}

	db := models.GetDB()

	// 获取商户的Mid
	merchant := &models.Merchant{}
	err := db.Where("id = ?", merchantID).First(merchant).Error
	if err != nil {
		return nil, fmt.Errorf("merchant not found: %v", err)
	}

	var trends []DashboardSettlementTrend

	for i := weeks - 1; i >= 0; i-- {
		// 计算周的开始和结束日期
		now := time.Now()
		weekStart := now.AddDate(0, 0, -int(now.Weekday())-i*7)
		weekEnd := weekStart.AddDate(0, 0, 6)

		var amount float64
		// 这里暂时使用已结算的交易作为结算趋势数据
		// 可以根据实际业务需求调整查询逻辑
		var payinAmount, payoutAmount float64

		db.Model(&models.MerchantPayin{}).
			Where("mid = ? AND created_at >= ? AND created_at <= ? AND settle_status = ?",
				merchant.Mid, weekStart.UnixMilli(), weekEnd.UnixMilli(), "settled").
			Select("COALESCE(SUM(amount), 0)").
			Row().Scan(&payinAmount)

		db.Model(&models.MerchantPayout{}).
			Where("mid = ? AND created_at >= ? AND created_at <= ? AND settle_status = ?",
				merchant.Mid, weekStart.UnixMilli(), weekEnd.UnixMilli(), "settled").
			Select("COALESCE(SUM(amount), 0)").
			Row().Scan(&payoutAmount)

		amount = payinAmount + payoutAmount

		trends = append(trends, DashboardSettlementTrend{
			Week:   fmt.Sprintf("第%d周", weeks-i),
			Amount: amount,
		})
	}

	return trends, nil
}

// GetAccountBalance 获取账户余额
func GetAccountBalance(merchantID int64) ([]DashboardAccountBalance, error) {
	if merchantID == 0 {
		return nil, errors.New("invalid merchant ID")
	}

	db := models.GetDB()
	var accounts []models.Account

	if err := db.Where("user_id = ? AND user_type = ?", fmt.Sprintf("%d", merchantID), "merchant").Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get account balance: %v", err)
	}

	var balances []DashboardAccountBalance
	for _, account := range accounts {
		if account.Asset == nil || account.Ccy == nil {
			continue
		}

		// 从 Asset 中获取余额信息
		balance := account.Asset.Balance.InexactFloat64()
		frozenAmt := account.Asset.FrozenBalance.InexactFloat64()
		availableAmt := account.Asset.AvailableBalance.InexactFloat64()

		balances = append(balances, DashboardAccountBalance{
			Currency:     *account.Ccy,
			Balance:      fmt.Sprintf("%.2f", balance),
			FrozenAmt:    fmt.Sprintf("%.2f", frozenAmt),
			AvailableAmt: fmt.Sprintf("%.2f", availableAmt),
		})
	}

	return balances, nil
}

// GetDashboardOverview 获取Dashboard概览数据
func GetDashboardOverview(merchantID int64) (*DashboardOverview, error) {
	if merchantID == 0 {
		return nil, errors.New("invalid merchant ID")
	}

	// 获取今日统计
	todayStats, err := GetTodayStats(merchantID, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get today stats: %v", err)
	}

	// 获取交易趋势（最近7天）
	transactionTrend, err := GetTransactionTrend(merchantID, 7)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction trend: %v", err)
	}

	return &DashboardOverview{
		TodayStats:       todayStats,
		TransactionTrend: transactionTrend,
	}, nil
}

// 辅助函数：计算增长率
func calculateRate(current, previous float64) string {
	if previous == 0 {
		if current > 0 {
			return "+100.0"
		}
		return "0.0"
	}
	rate := (current - previous) / previous * 100
	if rate >= 0 {
		return fmt.Sprintf("+%.1f", rate)
	}
	return fmt.Sprintf("%.1f", rate)
}

// 辅助函数：解析float64字符串
func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}
