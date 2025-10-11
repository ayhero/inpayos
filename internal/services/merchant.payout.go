package services

import (
	"inpayos/internal/config"
	"inpayos/internal/log"
	"inpayos/internal/middleware"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type MerchantPayoutService struct {
}

var (
	merchantPayoutService     *MerchantPayoutService
	merchantPayoutServiceOnce sync.Once
)

func SetupMerchantPayoutService() {
	merchantPayoutServiceOnce.Do(func() {
		merchantPayoutService = &MerchantPayoutService{}
	})
}

// GetMerchantPayoutService 获取Payout服务单例
func GetMerchantPayoutService() *MerchantPayoutService {
	if merchantPayoutService == nil {
		SetupMerchantPayoutService()
	}
	return merchantPayoutService
}

func (s *MerchantPayoutService) Create(ctx *gin.Context, req *protocol.MerchantPayoutRequest) (info *protocol.Transaction, code protocol.ErrorCode) {
	// 检查是否已存在相同的请求ID
	payout := models.GetMerchantPayoutByReqID(req.Mid, req.ReqID)
	if payout != nil {
		code = protocol.DuplicateTransaction
		return
	}

	m := middleware.GetMerchantFromContext(ctx)
	amount, _err := decimal.NewFromString(req.Amount)
	if _err != nil {
		code = protocol.InvalidParams
		return
	}
	now := time.Now()
	// 直接创建Transaction实体
	payout = &models.MerchantPayout{
		Mid:                  req.Mid,
		TrxType:              protocol.TrxTypePayout,
		ReqID:                req.ReqID,
		TrxID:                utils.GeneratePayinID(),
		Ccy:                  req.Ccy,
		Amount:               &amount,
		TrxMethod:            req.TrxMethod,
		TrxMode:              req.TrxMode,
		TrxApp:               req.TrxApp,
		Pkg:                  req.Pkg,
		Did:                  req.Did,
		ProductID:            req.ProductID,
		UserIP:               req.UserIP,
		ReturnURL:            req.ReturnURL,
		MerchantPayoutValues: &models.MerchantPayoutValues{},
	}
	payoutCfg := config.Get().MerchantPayout
	payout.SetStatus(protocol.StatusPending).
		SetExpiredAt(now.Add(time.Duration(payoutCfg.ExpiryMinutes) * time.Minute).UnixMilli()) //过期时间

	// 设置通知URL
	if payout.GetNotifyURL() == "" && m.GetNotifyURL() != "" {
		payout.SetNotifyURL(m.GetNotifyURL())
	}

	trans := payout.ToTransaction()
	// 获取可用渠道
	routerInfo := GetChannelRouterByMerchant(trans)
	if routerInfo == nil {
		code = protocol.ChannelNotFound
		return
	}

	// 设置渠道信息
	payout.SetChannelGroup(routerInfo.ChannelGroup)
	trans.SetChannelGroup(routerInfo.ChannelGroup)
	values := models.NewTrxValues()
	er := models.WriteDB.Transaction(func(tx *gorm.DB) error {
		// 创建代收订单
		if err := tx.Create(payout).Error; err != nil {
			return err
		}
		// 执行渠道请求
		result, errCode := RequestByRouter(ctx, tx, trans, routerInfo)
		if errCode != protocol.Success {
			code = errCode
			return protocol.NewServiceError(errCode, "channel request error")
		}
		values.SetStatus(result.Status).
			SetChannelStatus(result.ChannelStatus).
			SetChannelCode(result.ChannelCode).
			SetChannelAccount(result.ChannelAccountID).
			SetChannelTrxID(result.ChannelTrxID).
			SetLink(result.Link).
			SetResCode(result.ResCode).
			SetResMsg(result.ResMsg).
			SetChannelFeeCcy(result.ChannelFeeCcy)
		values.ChannelFeeAmount = result.ChannelFeeAmount
		if result.Status == protocol.StatusFailed || result.ChannelStatus == protocol.StatusSuccess {
			values.SetCompletedAt(utils.TimeNowMilli())
		}
		return nil
	})
	if er != nil {
		return
	}

	if _err := models.SaveTransactionValues(models.WriteDB, payout.ToTransaction(), values); _err != nil {
		log.Get().Errorf("SaveTransactionValues error: %v", _err)
	}
	AfterTransactionCreate(payout.ToTransaction())
	return
}
