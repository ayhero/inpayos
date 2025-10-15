package services

import (
	"context"
	"fmt"
	"inpayos/internal/log"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/task"
	"time"
)

var (
	// 近期统计指标（3天内）：每分钟统计一次
	recentTimeRanges = []*protocol.TimeIndex{
		// 最近时间范围（分钟级）
		{Range: protocol.IDX_RECENT, Target: 1, Unit: protocol.IDX_TIME_MINUTE},  // 最近1分钟
		{Range: protocol.IDX_RECENT, Target: 5, Unit: protocol.IDX_TIME_MINUTE},  // 最近5分钟
		{Range: protocol.IDX_RECENT, Target: 10, Unit: protocol.IDX_TIME_MINUTE}, // 最近10分钟
		{Range: protocol.IDX_RECENT, Target: 15, Unit: protocol.IDX_TIME_MINUTE}, // 最近15分钟
		{Range: protocol.IDX_RECENT, Target: 30, Unit: protocol.IDX_TIME_MINUTE}, // 最近30分钟
		{Range: protocol.IDX_RECENT, Target: 45, Unit: protocol.IDX_TIME_MINUTE}, // 最近45分钟

		// 最近时间范围（小时级）
		{Range: protocol.IDX_RECENT, Target: 1, Unit: protocol.IDX_TIME_HOUR},  // 最近1小时
		{Range: protocol.IDX_RECENT, Target: 2, Unit: protocol.IDX_TIME_HOUR},  // 最近2小时
		{Range: protocol.IDX_RECENT, Target: 3, Unit: protocol.IDX_TIME_HOUR},  // 最近3小时
		{Range: protocol.IDX_RECENT, Target: 6, Unit: protocol.IDX_TIME_HOUR},  // 最近6小时
		{Range: protocol.IDX_RECENT, Target: 12, Unit: protocol.IDX_TIME_HOUR}, // 最近12小时
		{Range: protocol.IDX_RECENT, Target: 24, Unit: protocol.IDX_TIME_HOUR}, // 最近24小时

		// 最近时间范围（天级）- 3天内
		{Range: protocol.IDX_RECENT, Target: 2, Unit: protocol.IDX_TIME_DAY}, // 最近2天
		{Range: protocol.IDX_RECENT, Target: 3, Unit: protocol.IDX_TIME_DAY}, // 最近3天

		// 上一时间段范围（用于环比分析）
		{Range: protocol.IDX_LAST, Target: 1, Unit: protocol.IDX_TIME_MINUTE},  // 上1分钟
		{Range: protocol.IDX_LAST, Target: 5, Unit: protocol.IDX_TIME_MINUTE},  // 上5分钟
		{Range: protocol.IDX_LAST, Target: 10, Unit: protocol.IDX_TIME_MINUTE}, // 上10分钟
		{Range: protocol.IDX_LAST, Target: 15, Unit: protocol.IDX_TIME_MINUTE}, // 上15分钟
		{Range: protocol.IDX_LAST, Target: 30, Unit: protocol.IDX_TIME_MINUTE}, // 上30分钟
		{Range: protocol.IDX_LAST, Target: 45, Unit: protocol.IDX_TIME_MINUTE}, // 上45分钟
		{Range: protocol.IDX_LAST, Target: 1, Unit: protocol.IDX_TIME_HOUR},    // 上1小时
		{Range: protocol.IDX_LAST, Target: 2, Unit: protocol.IDX_TIME_HOUR},    // 上2小时
		{Range: protocol.IDX_LAST, Target: 3, Unit: protocol.IDX_TIME_HOUR},    // 上3小时
	}

	// 历史统计指标（3天以上）：每小时统计一次
	historicalTimeRanges = []*protocol.TimeIndex{
		// 最近时间范围（天级）- 3天以上
		{Range: protocol.IDX_RECENT, Target: 7, Unit: protocol.IDX_TIME_DAY},  // 最近7天
		{Range: protocol.IDX_RECENT, Target: 15, Unit: protocol.IDX_TIME_DAY}, // 最近15天
		{Range: protocol.IDX_RECENT, Target: 30, Unit: protocol.IDX_TIME_DAY}, // 最近30天

		// 历史时间段范围（用于同比分析）
		{Range: protocol.IDX_LAST, Target: 7, Unit: protocol.IDX_TIME_DAY},  // 上7天
		{Range: protocol.IDX_LAST, Target: 15, Unit: protocol.IDX_TIME_DAY}, // 上15天
		{Range: protocol.IDX_LAST, Target: 30, Unit: protocol.IDX_TIME_DAY}, // 上30天
	}
)

func init() {
	// 注册近期统计任务处理器
	task.RegisterHandler(protocol.StatisticsRecent, HandleRecentStats)
	// 注册历史统计任务处理器
	task.RegisterHandler(protocol.StatisticsHistorical, HandleHistoricalStats)
}

func RegisterSummaryTasks() {
	log.Get().Info("注册统计任务...")

	// 定义统计相关的系统任务
	tasks := []*models.Task{
		{
			TaskID:     "statistics_recent",
			Type:       protocol.StatisticsRecent,
			HandlerKey: protocol.StatisticsRecent,
			Name:       "近期统计数据更新",
			TaskValues: &models.TaskValues{
				Cron:    &[]string{"@every 30s"}[0], // 每30秒执行一次
				Timeout: &[]int{60}[0],              // 1分钟超时
				Status:  &[]string{protocol.StatusEnabled}[0],
				Params:  map[string]any{},
			},
		},
		{
			TaskID:     "statistics_historical",
			Type:       protocol.StatisticsHistorical,
			HandlerKey: protocol.StatisticsHistorical,
			Name:       "历史统计数据更新",
			TaskValues: &models.TaskValues{
				Cron:    &[]string{"@every 1h"}[0], // 每小时执行一次
				Timeout: &[]int{1800}[0],           // 30分钟超时
				Status:  &[]string{protocol.StatusEnabled}[0],
				Params:  map[string]any{},
			},
		},
	}

	task.InitTasks(tasks)
	log.Get().Infof("统计任务注册完成，共 %d 个任务", len(tasks))
}

// HandleRecentStats 处理近期统计任务
func HandleRecentStats(ctx context.Context, params protocol.MapData) error {
	startTime := time.Now()
	if err := UpdateTransactionStatistics(ctx, recentTimeRanges); err != nil {
		return fmt.Errorf("更新近期统计失败: %v", err)
	}
	log.Get().Infof("更新近期统计完成，耗时：%v", time.Since(startTime))
	return nil
}

// HandleHistoricalStats 处理历史统计任务
func HandleHistoricalStats(ctx context.Context, params protocol.MapData) error {
	startTime := time.Now()
	if err := UpdateTransactionStatistics(ctx, historicalTimeRanges); err != nil {
		return fmt.Errorf("更新历史统计失败: %v", err)
	}
	log.Get().Infof("更新历史统计完成，耗时：%v", time.Since(startTime))
	return nil
}
