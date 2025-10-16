package services

import (
	"errors"
	"fmt"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"time"
)

// GetCashierTeamTodayStats 获取今日统计数据
func GetCashierTeamTodayStats(tid string, currency string) (*DashboardTodayStats, error) {
	if tid == "" {
		return nil, errors.New("invalid cashier team")
	}

	fmt.Printf("GetTodayStats called with merchantID: %v\n", tid)

	db := models.GetDB()

	// 检查数据库中是否有任何CashierPayin记录
	var totalPayinCount int64
	db.Model(&models.CashierPayin{}).Count(&totalPayinCount)
	fmt.Printf("Total payins in database: %d\n", totalPayinCount)

	// 检查数据库中是否有任何CashierPayout记录
	var totalPayoutCount int64
	db.Model(&models.CashierPayout{}).Count(&totalPayoutCount)
	fmt.Printf("Total payouts in database: %d\n", totalPayoutCount)

	// 首先获取商户的Tid
	team := &models.CashierTeam{}
	err := db.Where("tid = ?", tid).First(team).Error
	if err != nil {
		return nil, fmt.Errorf("merchant not found: %v", err)
	}
	fmt.Printf("Found merchant: ID=%d, Tid=%s\n", team.ID, team.Tid)

	// 检查特定商户的记录
	var merchantPayinCount int64
	db.Model(&models.CashierPayin{}).Where("tid = ?", team.Tid).Count(&merchantPayinCount)
	fmt.Printf("Payins for merchant %s: %d\n", team.Tid, merchantPayinCount)

	var merchantPayoutCount int64
	db.Model(&models.CashierPayout{}).Where("tid = ?", team.Tid).Count(&merchantPayoutCount)
	fmt.Printf("Payouts for merchant %s: %d\n", team.Tid, merchantPayoutCount)

	// 计算今日和昨日的时间范围（毫秒时间戳）
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UnixMilli()
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999, now.Location()).UnixMilli()
	yesterdayStart := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location()).UnixMilli()
	yesterdayEnd := time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 999, now.Location()).UnixMilli()

	// 获取今日代收数据
	var todayPayin, yesterdayPayin float64
	var todayPayinCount int64

	// 今日代收 - 先不限制时间，查询所有数据
	fmt.Printf("Querying today payin: tid=%s, start=%d, end=%d\n", team.Tid, todayStart, todayEnd)

	// 临时：查询所有数据不限时间
	result := db.Model(&models.CashierPayin{}).
		Where("tid = ?", team.Tid).
		Select("COALESCE(SUM(amount), 0) as total, COUNT(*) as count").
		Row()
	scanErr := result.Scan(&todayPayin, &todayPayinCount)
	if scanErr != nil {
		fmt.Printf("Error scanning today payin: %v\n", scanErr)
	}
	fmt.Printf("Today payin result (all time): total=%.2f, count=%d\n", todayPayin, todayPayinCount)

	// 查询状态分布
	var statusCount int64
	db.Model(&models.CashierPayin{}).Where("tid = ? AND status = ?", team.Tid, protocol.StatusSuccess).Count(&statusCount)
	fmt.Printf("Success status payins: %d\n", statusCount)

	// 查询最新的几条记录看时间戳
	var recentPayins []models.CashierPayin
	db.Where("tid = ?", team.Tid).Order("created_at DESC").Limit(3).Find(&recentPayins)
	for i, payin := range recentPayins {
		fmt.Printf("Recent payin %d: ID=%d, Amount=%.2f, Status=%s, CreatedAt=%d\n",
			i+1, payin.ID, payin.Amount, payin.Status, payin.CreatedAt)
	}

	// 昨日代收 - 临时设为0，专注今日数据
	yesterdayPayin = 0

	// 获取今日代付数据
	var todayPayout, yesterdayPayout float64
	var todayPayoutCount int64

	// 今日代付 - 查询所有数据不限时间
	db.Model(&models.CashierPayout{}).
		Where("tid = ?", team.Tid).
		Select("COALESCE(SUM(amount), 0) as total, COUNT(*) as count").
		Row().Scan(&todayPayout, &todayPayoutCount)
	fmt.Printf("Today payout result (all time): total=%.2f, count=%d\n", todayPayout, todayPayoutCount)

	// 昨日代付 - 临时设为0
	yesterdayPayout = 0

	// 计算成功率
	var todayTotal, todaySuccess int64

	// 今日总交易数 = 代收总数 + 代付总数
	var todayPayinTotal, todayPayoutTotal int64
	db.Model(&models.CashierPayin{}).
		Where("tid = ? AND created_at >= ? AND created_at <= ?", team.Tid, todayStart, todayEnd).
		Count(&todayPayinTotal)
	db.Model(&models.CashierPayout{}).
		Where("tid = ? AND created_at >= ? AND created_at <= ?", team.Tid, todayStart, todayEnd).
		Count(&todayPayoutTotal)
	todayTotal = todayPayinTotal + todayPayoutTotal

	// 今日成功交易数 = 代收成功数 + 代付成功数
	var todayPayinSuccess, todayPayoutSuccess int64
	db.Model(&models.CashierPayin{}).
		Where("tid = ? AND created_at >= ? AND created_at <= ? AND status = ?", team.Tid, todayStart, todayEnd, "success").
		Count(&todayPayinSuccess)
	db.Model(&models.CashierPayout{}).
		Where("tid = ? AND created_at >= ? AND created_at <= ? AND status = ?", team.Tid, todayStart, todayEnd, "success").
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
	db.Model(&models.CashierPayin{}).
		Where("tid = ? AND created_at >= ? AND created_at <= ?", team.Tid, yesterdayStart, yesterdayEnd).
		Count(&yesterdayPayinTotal)
	db.Model(&models.CashierPayout{}).
		Where("tid = ? AND created_at >= ? AND created_at <= ?", team.Tid, yesterdayStart, yesterdayEnd).
		Count(&yesterdayPayoutTotal)
	yesterdayTotal = yesterdayPayinTotal + yesterdayPayoutTotal

	// 昨日成功交易数
	var yesterdayPayinSuccess, yesterdayPayoutSuccess int64
	db.Model(&models.CashierPayin{}).
		Where("tid = ? AND created_at >= ? AND created_at <= ? AND status = ?", team.Tid, yesterdayStart, yesterdayEnd, protocol.StatusSuccess).
		Count(&yesterdayPayinSuccess)
	db.Model(&models.CashierPayout{}).
		Where("tid = ? AND created_at >= ? AND created_at <= ? AND status = ?", team.Tid, yesterdayStart, yesterdayEnd, protocol.StatusSuccess).
		Count(&yesterdayPayoutSuccess)
	yesterdaySuccess = yesterdayPayinSuccess + yesterdayPayoutSuccess

	var yesterdaySuccessRate float64
	if yesterdayTotal > 0 {
		yesterdaySuccessRate = float64(yesterdaySuccess) / float64(yesterdayTotal) * 100
	}

	// 计算增长率
	payinRate := calculateRate(todayPayin, yesterdayPayin)
	payoutRate := calculateRate(todayPayout, yesterdayPayout)

	return &DashboardTodayStats{
		TodayPayin:        fmt.Sprintf("%.2f", todayPayin),
		TodayPayinRate:    payinRate,
		TodayPayout:       fmt.Sprintf("%.2f", todayPayout),
		TodayPayoutRate:   payoutRate,
		SuccessRate:       fmt.Sprintf("%.1f", successRate),
		SuccessRateChange: fmt.Sprintf("%.1f", successRate-yesterdaySuccessRate),
	}, nil
}

// GetTransactionTrend 获取交易趋势数据
func GetCashierTeamTransactionTrend(tid string, days int) ([]DashboardTransactionTrend, error) {
	if tid == "" {
		return nil, errors.New("invalid cashier team")
	}

	db := models.GetDB()

	// 获取商户的Tid
	team := &models.CashierTeam{}
	err := db.Where("tid = ?", tid).First(team).Error
	if err != nil {
		return nil, fmt.Errorf("merchant not found: %v", err)
	}

	var trends []DashboardTransactionTrend

	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		displayDate := date.Format("01-02")

		var payin, payout float64

		// 临时：获取所有代收数据（不限时间）
		if i == 0 { // 只在第一天显示数据，其他天为0
			db.Model(&models.CashierPayin{}).
				Where("tid = ?", team.Tid).
				Select("COALESCE(SUM(amount), 0)").
				Row().Scan(&payin)

			// 获取所有代付数据（不限时间）
			db.Model(&models.CashierPayout{}).
				Where("tid = ?", team.Tid).
				Select("COALESCE(SUM(amount), 0)").
				Row().Scan(&payout)
		} else {
			payin = 0
			payout = 0
		}

		trends = append(trends, DashboardTransactionTrend{
			Date:   displayDate,
			Payin:  int64(payin),
			Payout: int64(payout),
		})
	}

	return trends, nil
}

// GetAccountBalance 获取账户余额
func GetCashierTeamAccountBalance(merchantID int64) ([]DashboardAccountBalance, error) {
	if merchantID == 0 {
		return nil, errors.New("invalid cashier team")
	}

	db := models.GetDB()
	var accounts []models.Account

	if err := db.Where("user_id = ? AND user_type = ?", fmt.Sprintf("%d", merchantID), "merchant").Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get account balance: %v", err)
	}

	var balances []DashboardAccountBalance
	for _, account := range accounts {
		if account.Asset == nil || account.Ccy == "" {
			continue
		}

		// 从 Asset 中获取余额信息
		balance := account.Asset.Balance.InexactFloat64()
		frozenAmt := account.Asset.FrozenBalance.InexactFloat64()
		availableAmt := account.Asset.AvailableBalance.InexactFloat64()

		balances = append(balances, DashboardAccountBalance{
			Currency:     account.Ccy,
			Balance:      fmt.Sprintf("%.2f", balance),
			FrozenAmt:    fmt.Sprintf("%.2f", frozenAmt),
			AvailableAmt: fmt.Sprintf("%.2f", availableAmt),
		})
	}

	return balances, nil
}

// GetDashboardOverview 获取Dashboard概览数据
func GetCashierTeamDashboardOverview(tid string) (*DashboardOverview, error) {
	if tid == "" {
		return nil, errors.New("invalid cashier team")
	}

	// 获取今日统计
	todayStats, err := GetCashierTeamTodayStats(tid, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get today stats: %v", err)
	}
	return &DashboardOverview{
		TodayStats: todayStats,
	}, nil
}
