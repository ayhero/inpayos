package services

import (
	"context"
	"errors"
	"inpayos/internal/channels"
	"inpayos/internal/log"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
	"sync"

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

func (s *MerchantTransactionService) CreatePayin(req *protocol.MerchantPayinRequest) (trx *protocol.Transaction, code protocol.ErrorCode) {
	return &protocol.Transaction{}, protocol.Success
}

func (s *MerchantTransactionService) CreatePayout(req *protocol.MerchantPayoutRequest) (trx *protocol.Transaction, code protocol.ErrorCode) {
	return &protocol.Transaction{}, protocol.Success
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

func RequestByRouter(ctx context.Context, tx *gorm.DB, trx *models.Transaction, routerInfo *protocol.RouterInfo) (result *protocol.ChannelResult, err error) {
	isAll := routerInfo.Strategy == protocol.RouterStrategyAll
	for _, account := range routerInfo.ChannelAccounts {
		trx.SetChannelCode(routerInfo.ChannelCodeLib[account]).SetChannelAccount(account)
		switch trx.TrxType {
		case protocol.TrxTypePayin:
			result, err = RequestChannelPayin(ctx, tx, trx)
		case protocol.TrxTypePayout:
			result, err = RequestChannelPayout(ctx, tx, trx)
		}
		if err == nil {
			//只要是失败，则继续循环
			if result.Status == protocol.StatusFailed {
				continue
			}
			break
		}
		//非全部轮询，则终止
		if !isAll {
			break
		}
	}
	return
}

func RequestChannelPayin(ctx context.Context, tx *gorm.DB, trx *models.Transaction) (result *protocol.ChannelResult, err error) {
	if trx.GetChannelAccount() == "" {
		return nil, errors.New("no channel account")
	}
	_, ok := channels.GetOpenApiChannelService(trx.GetChannelAccount())
	if !ok {
		err = errors.New("channel not supported")
		return
	}
	//result = svc.Payin(ctx, trx)
	// 暂时简化实现，直接返回成功结果
	// TODO: 实现完整的渠道支付逻辑
	result = &protocol.ChannelResult{
		Status:        protocol.StatusSuccess,
		ChannelStatus: protocol.StatusSuccess,
		ResCode:       "0000",
		ResMsg:        "Success",
		ChannelTrxID:  utils.GenerateID(),
		Link:          "",
	}

	log.Get().Infof("RequestChannelPayin: trxID=%s, channelAccount=%s", trx.TrxID, trx.GetChannelAccount())
	return
}

func RequestChannelPayout(ctx context.Context, db *gorm.DB, trx *models.Transaction) (result *protocol.ChannelResult, err error) {
	if trx.GetChannelAccount() == "" {
		return nil, errors.New("no channel account")
	}
	_, ok := channels.GetOpenApiChannelService(trx.GetChannelAccount())
	if !ok {
		err = errors.New("channel not supported")
		return
	}

	// 暂时简化实现，直接返回成功结果
	// TODO: 实现完整的渠道代付逻辑
	result = &protocol.ChannelResult{
		Status:        protocol.StatusSuccess,
		ChannelStatus: "SUCCESS",
		ResCode:       "0000",
		ResMsg:        "Success",
		ChannelTrxID:  trx.TrxID + "_channel",
		Link:          "",
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
	tempTrx := &models.Transaction{TrxType: query.TrxType}
	if tempTrx.TableName() == "" {
		return nil, 0, protocol.InvalidParams
	}

	db := models.ReadDB.Model(tempTrx)
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

	// 设置交易类型
	for i := range transactions {
		transactions[i].TrxType = query.TrxType
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
