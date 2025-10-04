package models

import (
	"fmt"
	"inpayos/internal/protocol"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// SummaryStats 统计数据结构
type SummaryStats struct {
	ID             int64           `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Target         string          `json:"target" gorm:"column:target;index"`               // 统计目标（商户ID/渠道账户/渠道组）
	TimeZone       string          `json:"time_zone" gorm:"column:time_zone"`               // 时区
	TimeIndex      string          `json:"time_index" gorm:"column:time_index;index"`       // 时间范围索引
	StartUnix      int64           `json:"start_unix" gorm:"column:start_unix;index"`       // 开始时间戳
	EndUnix        int64           `json:"end_unix" gorm:"column:end_unix;index"`           // 结束时间戳
	RefreshAt      int64           `json:"refresh_at" gorm:"column:refresh_at"`             // 刷新时间戳
	TotalCount     int64           `json:"total_count" gorm:"column:total_count"`           // 总笔数
	SuccCount      int64           `json:"succ_count" gorm:"column:succ_count"`             // 成功笔数
	FailCount      int64           `json:"fail_count" gorm:"column:fail_count"`             // 失败笔数
	PndCount       int64           `json:"pnd_count" gorm:"column:pnd_count"`               // 处理中笔数
	TotalUsdAmount decimal.Decimal `json:"total_usd_amount" gorm:"column:total_usd_amount"` // 总金额(USD)
	SuccUsdAmount  decimal.Decimal `json:"succ_usd_amount" gorm:"column:succ_usd_amount"`   // 成功金额(USD)
	FailUsdAmount  decimal.Decimal `json:"fail_usd_amount" gorm:"column:fail_usd_amount"`   // 失败金额(USD)
	PndUsdAmount   decimal.Decimal `json:"pnd_usd_amount" gorm:"column:pnd_usd_amount"`     // 处理中金额(USD)
	SuccRate       decimal.Decimal `json:"succ_rate" gorm:"column:succ_rate"`               // 成功率
	FailRate       decimal.Decimal `json:"fail_rate" gorm:"column:fail_rate"`               // 失败率
	PndRate        decimal.Decimal `json:"pnd_rate" gorm:"column:pnd_rate"`                 // 处理中率
	CreatedAt      int64           `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt      int64           `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

// TableName 返回数据库表名
func (s *SummaryStats) TableName() string {
	return "t_summary_stats"
}

// CalculateStatistics 只负责实时数据查询
func CalculateStatistics(db *gorm.DB, targetType string, startTime, endTime time.Time) ([]*SummaryStats, error) {
	var paymentStats, withdrawStats []*SummaryStats
	var err error

	// 获取支付统计
	paymentStats, err = CalculatePaymentStatistics(db, targetType, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("calculate payment stats error: %v", err)
	}

	// 获取提现统计
	withdrawStats, err = CalculateWithdrawStatistics(db, targetType, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("calculate withdraw stats error: %v", err)
	}

	// 合并同一目标的支付和提现统计
	mergedMap := make(map[string]*SummaryStats)

	// 处理支付统计
	for _, ps := range paymentStats {
		mergedMap[ps.Target] = ps
	}

	// 合并提现统计
	for _, ws := range withdrawStats {
		if ms, exists := mergedMap[ws.Target]; exists {
			// 合并已存在目标的统计数据
			ms.TotalCount += ws.TotalCount
			ms.SuccCount += ws.SuccCount
			ms.FailCount += ws.FailCount
			ms.PndCount += ws.PndCount
			ms.TotalUsdAmount = ms.TotalUsdAmount.Add(ws.TotalUsdAmount)
			ms.SuccUsdAmount = ms.SuccUsdAmount.Add(ws.SuccUsdAmount)
			ms.FailUsdAmount = ms.FailUsdAmount.Add(ws.FailUsdAmount)
			ms.PndUsdAmount = ms.PndUsdAmount.Add(ws.PndUsdAmount)
		} else {
			// 添加新的统计数据
			mergedMap[ws.Target] = ws
		}
	}

	// 转换为数组
	result := make([]*SummaryStats, 0, len(mergedMap))
	for _, stat := range mergedMap {
		result = append(result, stat)
	}

	return result, nil
}

// CalculatePaymentStatistics 计算支付统计数据
func CalculatePaymentStatistics(db *gorm.DB, targetType string, startTime, endTime time.Time) ([]*SummaryStats, error) {
	var stats []*SummaryStats

	// 使用TrxTypeTableMap获取表名
	tableName, ok := TrxTypeTableMap[protocol.TrxTypePayin]
	if !ok {
		return nil, fmt.Errorf("unsupported transaction type: %s", protocol.TrxTypePayin)
	}

	query := db.Table(tableName)

	groupByField := ""
	switch targetType {
	case "merchant":
		groupByField = "mid"
	case "channel_account":
		groupByField = "channel_account"
	case "channel_group":
		groupByField = "channel_group"
	default:
		return nil, fmt.Errorf("unknown target type: %s", targetType)
	}

	query = query.Where("created_at BETWEEN ? AND ?", startTime.Unix(), endTime.Unix())

	err := query.Select(fmt.Sprintf(`
		%s as target,
		COUNT(*) as %s,
		COALESCE(SUM(usd_amount), 0) as %s,
		COUNT(CASE WHEN status = 1 THEN 1 END) as %s,
		COALESCE(SUM(CASE WHEN status = 1 THEN usd_amount ELSE 0 END), 0) as %s,
		COUNT(CASE WHEN status = 2 THEN 1 END) as %s,
		COALESCE(SUM(CASE WHEN status = 2 THEN usd_amount ELSE 0 END), 0) as %s,
		COUNT(CASE WHEN status = 0 THEN 1 END) as %s,
		COALESCE(SUM(CASE WHEN status = 0 THEN usd_amount ELSE 0 END), 0) as %s,
		CAST(COUNT(CASE WHEN status = 1 THEN 1 END) AS DECIMAL(20,8)) / NULLIF(COUNT(*), 0) as %s,
		CAST(COUNT(CASE WHEN status = 2 THEN 1 END) AS DECIMAL(20,8)) / NULLIF(COUNT(*), 0) as %s,
		CAST(COUNT(CASE WHEN status = 0 THEN 1 END) AS DECIMAL(20,8)) / NULLIF(COUNT(*), 0) as %s
	`, groupByField,
		protocol.STAT_IDX_TOTAL_COUNT,
		protocol.STAT_IDX_TOTAL_USD_AMOUNT,
		protocol.STAT_IDX_SUCC_COUNT,
		protocol.STAT_IDX_SUCC_USD_AMOUNT,
		protocol.STAT_IDX_FAIL_COUNT,
		protocol.STAT_IDX_FAIL_USD_AMOUNT,
		protocol.STAT_IDX_PND_COUNT,
		protocol.STAT_IDX_PND_USD_AMOUNT,
		protocol.STAT_IDX_SUCC_RATE,
		protocol.STAT_IDX_FAIL_RATE,
		protocol.STAT_IDX_PND_RATE,
	)).Group(groupByField).Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// CalculateWithdrawStatistics 计算提现统计数据
func CalculateWithdrawStatistics(db *gorm.DB, targetType string, startTime, endTime time.Time) ([]*SummaryStats, error) {
	var stats []*SummaryStats
	query := db.Table("t_withdraw")

	groupByField := ""
	switch targetType {
	case "merchant":
		groupByField = "mid"
	case "channel_account":
		groupByField = "channel_account_id"
	case "channel_group":
		groupByField = "channel_group"
	default:
		return nil, fmt.Errorf("unknown target type: %s", targetType)
	}

	query = query.Where("created_at BETWEEN ? AND ?", startTime.Unix(), endTime.Unix())

	err := query.Select(fmt.Sprintf(`
		%s as target,
		COUNT(*) as %s,
		COALESCE(SUM(usd_amount), 0) as %s,
		COUNT(CASE WHEN status = 1 THEN 1 END) as %s,
		COALESCE(SUM(CASE WHEN status = 1 THEN usd_amount ELSE 0 END), 0) as %s,
		COUNT(CASE WHEN status = 2 THEN 1 END) as %s,
		COALESCE(SUM(CASE WHEN status = 2 THEN usd_amount ELSE 0 END), 0) as %s,
		COUNT(CASE WHEN status = 0 THEN 1 END) as %s,
		COALESCE(SUM(CASE WHEN status = 0 THEN usd_amount ELSE 0 END), 0) as %s,
		CAST(COUNT(CASE WHEN status = 1 THEN 1 END) AS DECIMAL(20,8)) / NULLIF(COUNT(*), 0) as %s,
		CAST(COUNT(CASE WHEN status = 2 THEN 1 END) AS DECIMAL(20,8)) / NULLIF(COUNT(*), 0) as %s,
		CAST(COUNT(CASE WHEN status = 0 THEN 1 END) AS DECIMAL(20,8)) / NULLIF(COUNT(*), 0) as %s
	`, groupByField,
		protocol.STAT_IDX_TOTAL_COUNT,
		protocol.STAT_IDX_TOTAL_USD_AMOUNT,
		protocol.STAT_IDX_SUCC_COUNT,
		protocol.STAT_IDX_SUCC_USD_AMOUNT,
		protocol.STAT_IDX_FAIL_COUNT,
		protocol.STAT_IDX_FAIL_USD_AMOUNT,
		protocol.STAT_IDX_PND_COUNT,
		protocol.STAT_IDX_PND_USD_AMOUNT,
		protocol.STAT_IDX_SUCC_RATE,
		protocol.STAT_IDX_FAIL_RATE,
		protocol.STAT_IDX_PND_RATE,
	)).Group(groupByField).Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// SaveStatistics 批量保存或更新统计数据
func SaveStatistics(db *gorm.DB, stats ...*SummaryStats) error {
	if len(stats) == 0 {
		return nil
	}

	// 使用事务批量保存或更新
	return db.Transaction(func(tx *gorm.DB) error {
		for _, stat := range stats {
			// 尝试根据唯一索引更新记录
			result := tx.Model(&SummaryStats{}).Where(
				"target = ? AND time_zone = ? AND time_index = ?",
				stat.Target,
				stat.TimeZone,
				stat.TimeIndex,
			).Updates(stat)

			if result.Error != nil {
				return result.Error
			}

			// 如果没有更新任何记录，则创建新记录
			if result.RowsAffected == 0 {
				if err := tx.Create(stat).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
