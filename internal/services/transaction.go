package services

import (
	"context"
	"inpayos/internal/channels"
	"inpayos/internal/log"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MerchantTransactionService struct {
	PayinService  *MerchantPayinService
	PayoutService *MerchantPayoutService
}

var (
	transactionService     *MerchantTransactionService
	transactionServiceOnce sync.Once
)

func SetupTransactionService() {
	transactionServiceOnce.Do(func() {
		transactionService = &MerchantTransactionService{
			PayinService:  GetMerchantPayinService(),
			PayoutService: GetMerchantPayoutService(),
		}
	})
}

// GetMerchantTransactionService 获取Transaction服务单例
func GetMerchantTransactionService() *MerchantTransactionService {
	if transactionService == nil {
		SetupTransactionService()
	}
	return transactionService
}

func (s *MerchantTransactionService) CreatePayin(ctx *gin.Context, req *protocol.MerchantPayinRequest) (trx *protocol.Transaction, code protocol.ErrorCode) {
	return s.PayinService.Create(ctx, req)
}

func (s *MerchantTransactionService) CreatePayout(ctx *gin.Context, req *protocol.MerchantPayoutRequest) (trx *protocol.Transaction, code protocol.ErrorCode) {
	return s.PayoutService.Create(ctx, req)
}

func (s *MerchantTransactionService) Cancel(req *protocol.MerchantCancelRequest) (trx *protocol.Transaction, code protocol.ErrorCode) {
	return &protocol.Transaction{}, protocol.Success
}
func (s *MerchantTransactionService) Query(req *protocol.MerchantQueryRequest) (trx *protocol.Transaction, code protocol.ErrorCode) {
	return &protocol.Transaction{}, protocol.Success
}

func GetChannelRouterByMerchant(req *models.Transaction) *protocol.RouterInfo {
	in := &RouterRequest{
		Mid:       req.Mid,
		TrxType:   req.TrxType,
		ReqID:     req.ReqID,
		Ccy:       req.Ccy,
		Amount:    req.Amount,
		TrxMethod: req.TrxMethod,
		TrxMode:   req.TrxMode,
		TrxApp:    req.TrxApp,
		Pkg:       req.Pkg,
		Did:       req.Did,
		ProductID: req.ProductID,
	}
	return GetChannelByMerchant(in)
}

func RequestByRouter(ctx context.Context, tx *gorm.DB, trx *models.Transaction, routerInfo *protocol.RouterInfo) (result *protocol.ChannelResult, err protocol.ErrorCode) {
	isAll := routerInfo.Strategy == protocol.RouterStrategyAll
	err = protocol.ChannelNotSupported
	for _, account := range routerInfo.ChannelAccounts {
		trx.SetChannelCode(routerInfo.ChannelCodeLib[account]).
			SetChannelAccount(account)
		switch trx.TrxType {
		case protocol.TrxTypePayin:
			result, err = RequestChannelPayin(ctx, trx)
		case protocol.TrxTypePayout:
			result, err = RequestChannelPayout(ctx, trx)
		}
		if err != protocol.Success {
			continue
		}
		if !isAll || result.Status != protocol.StatusFailed {
			result.ChannelAccountID = trx.GetChannelCode()
			result.ChannelAccountID = trx.GetChannelAccount()
			break
		}
	}
	return
}

func RequestChannelPayin(ctx context.Context, trx *models.Transaction) (result *protocol.ChannelResult, err protocol.ErrorCode) {
	err = protocol.Success
	if trx.GetChannelAccount() == "" {
		return nil, protocol.InvalidParams
	}
	err = protocol.Success
	svc, ok := channels.GetOpenApiChannelService(trx.GetChannelAccount())
	if !ok {
		err = protocol.ChannelNotSupported
		return
	}
	in := &channels.ChannelTrxRequest{
		Transaction: trx,
	}
	result = svc.Payin(in)
	// 7. 处理支付结果
	if result == nil {
		result = &protocol.ChannelResult{
			Status:  protocol.StatusPending,
			ResCode: protocol.ResCodeResponseError,
		}
	}
	log.Get().Infof("RequestChannelPayin: trxID=%s, channelAccount=%s", trx.TrxID, trx.GetChannelAccount())
	return
}

func RequestChannelPayout(ctx context.Context, trx *models.Transaction) (result *protocol.ChannelResult, err protocol.ErrorCode) {
	err = protocol.Success
	if trx.GetChannelAccount() == "" {
		return nil, protocol.InvalidParams
	}
	svc, ok := channels.GetOpenApiChannelService(trx.GetChannelAccount())
	if !ok {
		err = protocol.ChannelNotSupported
		return
	}
	in := &channels.ChannelTrxRequest{
		Transaction: trx,
	}
	result = svc.Payout(in)
	// 7. 处理支付结果
	if result == nil {
		result = &protocol.ChannelResult{
			Status:  protocol.StatusPending,
			ResCode: protocol.ResCodeResponseError,
		}
	}
	log.Get().Infof("RequestChannelPayout: trxID=%s, channelAccount=%s", trx.TrxID, trx.GetChannelAccount())
	return
}

func RefreshTrxFlag(trx *models.Transaction) {

}

func AfterTransactionCreate(trx *models.Transaction) {
	go func() {
		RefreshTrxFlag(trx)
		history := models.NewTrxHistoryByTransaction(trx)
		if _err := models.CreateHistory(history); _err != nil {
			log.Get().Error("Save Transaction history error: ", _err)
		}
	}()
}

// ListTransactionByQuery 统一查询交易列表
func ListTransactionByQuery(query *models.TrxQuery) ([]*models.Transaction, int64, protocol.ErrorCode) {
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

func AfterPayinSuccess(trx *models.Transaction) {
	// 处理支付成功后的逻辑
	// 例如，更新商户配置、发送通知等
	log.Get().Infof("Payin transaction %s succeeded for merchant %s", trx.TrxID, trx.Mid)
	// 设置结算状态为待结算
	trx.SetSettleStatus(protocol.StatusPending)
}

// TodayStats 今日统计数据
type TodayStats struct {
	TotalAmount  string  `json:"totalAmount"`  // 总金额
	TotalCount   int64   `json:"totalCount"`   // 总交易数
	SuccessCount int64   `json:"successCount"` // 成功交易数
	SuccessRate  float64 `json:"successRate"`  // 成功率（百分比）
	PendingCount int64   `json:"pendingCount"` // 待处理交易数
}

// GetTransactionTodayStats 获取今日交易统计数据
func GetTransactionTodayStats(mid string, trxType string) (*TodayStats, protocol.ErrorCode) {
	// 计算今日时间范围（毫秒时间戳）
	todayStart := utils.TodayZeroTimeMilli()
	todayEnd := todayStart + 86400000

	var stats TodayStats
	err := models.GetTransactionQueryByType(trxType).
		Where("mid=?", mid).
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
