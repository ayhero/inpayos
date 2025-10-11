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

type MerchantPayinService struct {
}

var (
	merchantPayinService     *MerchantPayinService
	merchantPayinServiceOnce sync.Once
)

func SetupMerchantPayinService() {
	merchantPayinServiceOnce.Do(func() {
		merchantPayinService = &MerchantPayinService{}
	})
}

// GetMerchantPayinService 获取Payin服务单例
func GetMerchantPayinService() *MerchantPayinService {
	if merchantPayinService == nil {
		SetupMerchantPayinService()
	}
	return merchantPayinService
}

func (s *MerchantPayinService) Create(ctx *gin.Context, req *protocol.MerchantPayinRequest) (info *protocol.Transaction, code protocol.ErrorCode) {

	// 检查是否已存在相同的请求ID
	payin := models.GetMerchantPayinByReqID(req.Mid, req.ReqID)
	if payin != nil {
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
	payin = &models.MerchantPayin{
		Mid:                 req.Mid,
		TrxType:             protocol.TrxTypePayin,
		ReqID:               req.ReqID,
		TrxID:               utils.GeneratePayinID(),
		Ccy:                 req.Ccy,
		Amount:              &amount,
		TrxMethod:           req.TrxMethod,
		TrxMode:             req.TrxMode,
		TrxApp:              req.TrxApp,
		Pkg:                 req.Pkg,
		Did:                 req.Did,
		ProductID:           req.ProductID,
		UserIP:              req.UserIP,
		ReturnURL:           req.ReturnURL,
		MerchantPayinValues: &models.MerchantPayinValues{},
	}
	payinCfg := config.Get().MerchantPayin
	payin.SetStatus(protocol.StatusPending).
		SetExpiredAt(now.Add(time.Duration(payinCfg.ExpiryMinutes) * time.Minute).UnixMilli()) //过期时间

	// 设置通知URL
	if payin.GetNotifyURL() == "" && m.GetNotifyURL() != "" {
		payin.SetNotifyURL(m.GetNotifyURL())
	}

	trans := payin.ToTransaction()
	// 获取可用渠道
	routerInfo := GetChannelRouterByMerchant(trans)
	if routerInfo == nil {
		code = protocol.ChannelNotFound
		return
	}
	// 设置渠道信息
	payin.SetChannelGroup(routerInfo.ChannelGroup)
	trans.SetChannelGroup(routerInfo.ChannelGroup)
	values := models.NewTrxValues()
	er := models.WriteDB.Transaction(func(tx *gorm.DB) error {
		// 创建代收订单
		if err := tx.Create(payin).Error; err != nil {
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
		if result.ExpiredAt > 0 {
			values.SetExpiredAt(result.ExpiredAt)
		}
		values.ChannelFeeAmount = result.ChannelFeeAmount
		if result.Status == protocol.StatusFailed || result.ChannelStatus == protocol.StatusSuccess {
			values.SetCompletedAt(utils.TimeNowMilli())
		}
		return nil
	})
	if er != nil {
		return
	}
	if _err := models.SaveTransactionValues(models.WriteDB, trans, values); _err != nil {
		log.Get().Errorf("SaveTransactionValues error: %v", _err)
	}
	AfterTransactionCreate(trans)
	return
}
