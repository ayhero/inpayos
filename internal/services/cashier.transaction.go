package services

import (
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
	"sync"
)

type CashierTransactionService struct {
}

var (
	cashierTransactionService     *CashierTransactionService
	cashierTransactionServiceOnce sync.Once
)

func SetupCashierTransactionService() {
	cashierTransactionServiceOnce.Do(func() {
		cashierTransactionService = &CashierTransactionService{}
	})
}

// GetCashierTransactionService 获取Transaction服务单例
func GetCashierTransactionService() *CashierTransactionService {
	if cashierTransactionService == nil {
		SetupCashierTransactionService()
	}
	return cashierTransactionService
}

// ListTransactionByQuery 统一查询交易列表
func (t *CashierTransactionService) ListTransactionByQuery(query *models.TrxQuery) ([]*models.Transaction, int64, protocol.ErrorCode) {
	var transactions []*models.Transaction
	var total int64
	var err error

	// 创建临时Transaction对象来获取表名
	db := models.GetTransactionQueryByType(query.TrxType)
	// 应用查询条件
	db = query.BuildQuery(db)

	// 统计总数
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, protocol.DatabaseError
	}

	// 查询列表
	err = db.Order("created_at desc").
		Offset(query.GetOffset()).
		Limit(query.GetLimit()).
		Find(&transactions).Error
	if err != nil {
		return nil, 0, protocol.DatabaseError
	}

	return transactions, total, protocol.Success
}

// GetTransactionTodayStats 获取今日交易统计数据
func (t *CashierTransactionService) GetTransactionTodayStats(tid string, trxType string) (*TodayStats, protocol.ErrorCode) {
	// 计算今日时间范围（毫秒时间戳）
	todayStart := utils.TodayZeroTimeMilli()
	todayEnd := todayStart + 86400000

	var stats TodayStats
	err := models.GetTransactionQueryByType(trxType).
		Where("tid=?", tid).
		Where("created_at >= ? AND created_at <= ?", todayStart, todayEnd).
		Select(`
			COUNT(*) as total_count,
			ROUND(COALESCE(SUM(CASE WHEN amount IS NOT NULL THEN amount ELSE 0 END), 0)::numeric, 2) as total_amount,
			COUNT(CASE WHEN status = 'success' THEN 1 END) as success_count,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_count,
			ROUND(COALESCE(CASE WHEN COUNT(*) > 0 THEN COUNT(CASE WHEN status = 'success' THEN 1 END) * 100.0 / COUNT(*) ELSE 0 END, 0)::numeric, 2) as success_rate
		`).Find(&stats).Error

	if err != nil {
		return nil, protocol.DatabaseError
	}

	return &stats, protocol.Success
}
