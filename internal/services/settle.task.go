package services

import (
	"context"
	"fmt"
	"inpayos/internal/log"
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
	task.RegisterHandler("settle.process", HandleSettleProcess)
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
	result := GetSettleService().SettleByCompletedTime(
		ctx,
		taskParams.TrxType,
		trxStartTime.UnixMilli(),
		trxEndTime.UnixMilli(),
	)

	// 记录执行结果
	log.Get().Infof("结算任务执行完成，耗时: %v，总交易数: %d，成功: %d，失败: %d \n",
		time.Since(startTime),
		result.TotalTransactions,
		result.SuccessTransactions,
		result.FailedTransactions,
	)

	return nil
}
