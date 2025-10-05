package services

import (
	"context"
	"fmt"
	"inpayos/internal/log"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"time"

	"github.com/shopspring/decimal"
)

// calculateTimeRange 根据 TimeIndex 计算统计的开始和结束时间
func calculateTimeRange(tr *protocol.TimeIndex, currentTime time.Time) (startTime, endTime time.Time) {
	params := protocol.MapData{
		protocol.QUERY_IDX_REFRESH_AT:      currentTime.Unix(),
		protocol.QUERY_IDX_TODAY_ZERO_UNIX: time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Unix(),
		protocol.QUERY_IDX_TODAY_AGE_UNIX:  currentTime.Unix() - time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Unix(),
	}

	tr.GenTimeRange(params)

	startUnix := params.GetInt64(protocol.QUEYRY_IDX_QUERY_START_UNIX)
	endUnix := params.GetInt64(protocol.QUERY_IDX_QUERY_END_UNIX)

	return time.Unix(startUnix, 0), time.Unix(endUnix, 0)
}

// initializeEmptyStats 初始化空的统计数据
func initializeEmptyStats(target string, timeZone string, timeIndex string, startTime, endTime time.Time) *models.SummaryStats {
	return &models.SummaryStats{
		Target:         target,
		TimeZone:       timeZone,
		TimeIndex:      timeIndex,
		StartUnix:      startTime.Unix(),
		EndUnix:        endTime.Unix(),
		RefreshAt:      time.Now().Unix(),
		TotalCount:     0,
		SuccCount:      0,
		FailCount:      0,
		PndCount:       0,
		TotalUsdAmount: decimal.Zero,
		SuccUsdAmount:  decimal.Zero,
		FailUsdAmount:  decimal.Zero,
		PndUsdAmount:   decimal.Zero,
		SuccRate:       decimal.Zero,
		FailRate:       decimal.Zero,
		PndRate:        decimal.Zero,
	}
}

// getTargets 根据类型获取需要统计的目标列表
func getTargets(targetType string) (map[string]bool, error) {
	targetsMap := make(map[string]bool)

	switch targetType {
	case "merchant":
		merchants, err := models.GetMerchantsByStatus(protocol.StatusActive)
		if err != nil {
			return nil, fmt.Errorf("获取商户列表失败: %v", err)
		}
		for _, m := range merchants {
			targetsMap[m.Mid] = true
		}
	case "channel_account":
		accounts := models.GetChannelAccounts()
		for _, a := range accounts {
			targetsMap[a.GetAccountID()] = true
		}
	case "channel_group":
		groups, err := models.GetActiveChannelGroups()
		if err != nil {
			return nil, fmt.Errorf("获取渠道组列表失败: %v", err)
		}
		for _, g := range groups {
			targetsMap[g.Code] = true
		}
	default:
		return nil, fmt.Errorf("未知的目标类型: %s", targetType)
	}

	return targetsMap, nil
}

// mergeStats 合并实时统计数据到基础统计数据
func mergeStats(baseStats, realtimeStats *models.SummaryStats) {
	baseStats.TotalCount = realtimeStats.TotalCount
	baseStats.SuccCount = realtimeStats.SuccCount
	baseStats.FailCount = realtimeStats.FailCount
	baseStats.PndCount = realtimeStats.PndCount
	baseStats.TotalUsdAmount = realtimeStats.TotalUsdAmount
	baseStats.SuccUsdAmount = realtimeStats.SuccUsdAmount
	baseStats.FailUsdAmount = realtimeStats.FailUsdAmount
	baseStats.PndUsdAmount = realtimeStats.PndUsdAmount
	baseStats.SuccRate = realtimeStats.SuccRate
	baseStats.FailRate = realtimeStats.FailRate
	baseStats.PndRate = realtimeStats.PndRate
}

// updateStats 处理指定类型的统计数据
func updateStats(targetType, timeZone, timeIndex string, startTime, endTime time.Time) error {
	// 1. 获取所有统计对象
	targetsMap, err := getTargets(targetType)
	if err != nil {
		return err
	}

	// 2. 获取实时统计数据
	realtimeStats, err := models.CalculateStatistics(models.ReadDB, targetType, startTime, endTime)
	if err != nil {
		return fmt.Errorf("计算%s统计数据失败: %v", targetType, err)
	}

	statsMap := make(map[string]*models.SummaryStats)
	// 3. 初始化所有目标的统计数据
	for target := range targetsMap {
		statsMap[target] = initializeEmptyStats(target, timeZone, timeIndex, startTime, endTime)
	}

	// 4. 合并实时统计数据
	for _, stat := range realtimeStats {
		if baseStats, exists := statsMap[stat.Target]; exists {
			mergeStats(baseStats, stat)
		}
	}

	// 5. 转换为数组并保存
	stats := make([]*models.SummaryStats, 0, len(statsMap))
	for _, stat := range statsMap {
		stats = append(stats, stat)
	}

	if err := models.SaveStatistics(models.WriteDB, stats...); err != nil {
		return fmt.Errorf("保存%s统计数据失败: %v", targetType, err)
	}

	return nil
}

// UpdateTransactionStatistics 更新指定时间范围的交易统计
func UpdateTransactionStatistics(ctx context.Context, timeRanges []*protocol.TimeIndex) error {
	currentTime := time.Now()
	timeZones := []string{protocol.Asia_ShangHai.String()}
	targetTypes := []string{"merchant", "channel_account", "channel_group"}

	// 如果未指定时间范围，则使用默认的时间范围
	if len(timeRanges) == 0 {
		return fmt.Errorf("未指定时间范围")
	}

	for _, tz := range timeZones {
		loc := protocol.GetTimeLocation(tz)
		if loc == nil {
			continue
		}
		tzCurrentTime := currentTime.In(loc)

		for _, tr := range timeRanges {
			startTime, endTime := calculateTimeRange(tr, tzCurrentTime)
			timeIndex := fmt.Sprintf("%s_%d_%s", tr.Range, tr.Target, tr.Unit)

			for _, targetType := range targetTypes {
				if err := updateStats(targetType, tz, timeIndex, startTime, endTime); err != nil {
					log.Get().Error(err)
				}
			}
		}
	}

	return nil
}
