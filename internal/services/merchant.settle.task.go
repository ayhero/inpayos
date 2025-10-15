package services

import (
	"context"
	"fmt"
	"inpayos/internal/log"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/task"
	"time"

	"github.com/spf13/cast"
)

// SettleTaskParams 结算任务参数结构体
type SettleTaskParams struct {
	TrxType string `json:"trx_type"` // 交易类型，默认为空表示所有类型
	Days    int    `json:"days"`     // 处理几天内的交易，默认3天
}

func init() {
	// 注册结算任务处理器
	task.RegisterHandler(protocol.MerchantSettleProcess, HandleSettleProcess)
	task.RegisterHandler(protocol.MerchantSettleAccounting, HandleSettleAccounting)
	task.RegisterHandler(protocol.MerchantTransactionSettleIDFix, HandleTransactionSettleIDFix)
}

func RegisterSettleTasks() {
	log.Get().Info("注册结算任务...")
	// 定义结算相关的系统任务
	tasks := []*models.Task{
		{
			TaskID:     "merchant_settle_process_payin",
			Type:       protocol.MerchantSettle,
			HandlerKey: protocol.MerchantSettleProcess,
			Name:       "商户入金结算处理",
			TaskValues: &models.TaskValues{
				Cron:    &[]string{"0 1 * * *"}[0], // 每天凌晨1点执行
				Timeout: &[]int{3600}[0],           // 1小时超时
				Status:  &[]string{protocol.StatusEnabled}[0],
				Params: map[string]any{
					"trx_type": protocol.TrxTypePayin,
					"days":     3,
				},
			},
		},
		{
			TaskID:     "merchant_settle_process_payout",
			Type:       protocol.MerchantSettle,
			HandlerKey: protocol.MerchantSettleProcess,
			Name:       "商户出金结算处理",
			TaskValues: &models.TaskValues{
				Cron:    &[]string{"0 2 * * *"}[0], // 每天凌晨2点执行
				Timeout: &[]int{3600}[0],           // 1小时超时
				Status:  &[]string{protocol.StatusEnabled}[0],
				Params: map[string]any{
					"trx_type": protocol.TrxTypePayout,
					"days":     3,
				},
			},
		},
		{
			TaskID:     "merchant_settle_accounting",
			Type:       protocol.MerchantSettle,
			HandlerKey: protocol.MerchantSettleAccounting,
			Name:       "商户结算记账处理",
			TaskValues: &models.TaskValues{
				Cron:    &[]string{"@every 15m"}[0], // 每15分钟执行一次
				Timeout: &[]int{600}[0],             // 10分钟超时
				Status:  &[]string{protocol.StatusEnabled}[0],
				Params:  map[string]any{},
			},
		},
		{
			TaskID:     "merchant_transaction_settle_id_fix",
			Type:       protocol.MerchantSettle,
			HandlerKey: protocol.MerchantTransactionSettleIDFix,
			Name:       "商户交易结算ID修复",
			TaskValues: &models.TaskValues{
				Cron:    &[]string{"0 3 * * *"}[0], // 每天凌晨3点执行
				Timeout: &[]int{1800}[0],           // 30分钟超时
				Status:  &[]string{protocol.StatusEnabled}[0],
				Params: map[string]any{
					"days": 7, // 处理最近7天的数据
				},
			},
		},
	}

	task.InitTasks(tasks)
	log.Get().Infof("结算任务注册完成，共 %d 个任务", len(tasks))
}

// HandleSettleProcess 处理结算任务
func HandleSettleProcess(ctx context.Context, params protocol.MapData) error {
	startTime := time.Now()

	// 解析任务参数
	taskParams := &SettleTaskParams{
		TrxType: params.Get("trx_type"),
		Days:    cast.ToInt(params.Get("days")),
	}

	// 设置默认值
	if taskParams.Days <= 0 {
		taskParams.Days = 3
	}
	if taskParams.TrxType == "" {
		return fmt.Errorf("交易类型不能为空")
	}

	// 计算时间范围：处理指定天数内完成的交易
	trxEndTime := time.Now()
	trxStartTime := trxEndTime.AddDate(0, 0, -taskParams.Days)

	log.Get().Infof("开始执行结算任务，交易完成时间范围: %v - %v",
		trxStartTime.Format(time.DateOnly),
		trxEndTime.Format(time.DateOnly),
	)

	// 执行结算处理，按交易完成时间查询
	result := GetMerchantSettleService().SettleByTimeRange(
		taskParams.TrxType,
		trxStartTime.UnixMilli(),
		trxEndTime.UnixMilli(),
	)

	// 记录执行结果
	log.Get().Infof("结算任务执行完成，耗时: %v，总交易数: %d，成功: %d，失败: %d \n",
		time.Since(startTime),
		result.TotalCount,
		result.SuccessCount,
		result.FailedCount,
	)

	return nil
}

// HandleSettleAccounting 处理结算记账任务
func HandleSettleAccounting(ctx context.Context, params protocol.MapData) error {
	log.Get().Info("HandleSettleAccounting: starting settle accounting task")

	// 从参数中获取时间戳，如果没有则使用当前时间
	var currentTime int64
	if ts, exists := params["timestamp"]; exists && ts != nil {
		currentTime = cast.ToInt64(ts)
	} else {
		currentTime = time.Now().UnixMilli()
	}

	log.Get().Infof("HandleSettleAccounting: processing with timestamp %d", currentTime)

	// 获取结算服务实例
	settleService := GetMerchantSettleService()
	if settleService == nil {
		return fmt.Errorf("HandleSettleAccounting: failed to get merchant settle service")
	}

	// 执行批量结算记账
	result := settleService.ProcessBatchSettleAccounting(currentTime)
	// 记录处理结果
	log.Get().Infof("HandleSettleAccounting: task completed - total: %d, success: %d, failed: %d, duration: %v",
		result.TotalCount, result.SuccessCount, result.FailedCount, result.Duration)

	return nil
}

// HandleTransactionSettleIDFix 处理交易结算ID修复任务
func HandleTransactionSettleIDFix(ctx context.Context, params protocol.MapData) error {
	log.Get().Info("HandleTransactionSettleIDFix: starting transaction settle_id fix task")

	// 从参数中获取处理天数，默认为7天
	days := cast.ToInt(params.Get("days"))
	if days <= 0 {
		days = 7
	}

	// 计算时间范围
	endTime := time.Now().UnixMilli()
	startTime := time.Now().AddDate(0, 0, -days).UnixMilli()

	log.Get().Infof("HandleTransactionSettleIDFix: processing time range from %d to %d (%d days)",
		startTime, endTime, days)

	// 获取结算服务实例
	settleService := GetMerchantSettleService()
	if settleService == nil {
		return fmt.Errorf("HandleTransactionSettleIDFix: failed to get merchant settle service")
	}

	// 执行修复操作
	result := settleService.FixTransactionSettleID(startTime, endTime)

	// 记录处理结果
	log.Get().Infof("HandleTransactionSettleIDFix: task completed - total: %d, success: %d, failed: %d, duration: %v",
		result.TotalCount, result.SuccessCount, result.FailedCount, result.Duration)
	return nil
}
