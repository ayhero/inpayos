package services

import (
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
	existingTrx := models.GetMerchantPayinByReqID(req.Mid, req.ReqID)
	if existingTrx != nil {
		code = protocol.DuplicateTransaction
		return
	}

	m := middleware.GetMerchantFromContext(ctx)
	amount, _err := decimal.NewFromString(req.Amount)
	if _err != nil {
		code = protocol.InvalidParams
		return
	}

	// 直接创建Transaction实体
	payin := &models.MerchantPayin{
		Mid:                 req.Mid,
		TrxType:             protocol.TrxTypePayin,
		ReqID:               req.ReqID,
		TrxID:               utils.GeneratePayinID(),
		MerchantPayinValues: &models.MerchantPayinValues{},
	}
	// 设置金额
	payin.Amount = &amount

	// 设置通知URL
	if payin.GetNotifyURL() == "" && m.GetNotifyURL() != "" {
		payin.SetNotifyURL(m.GetNotifyURL())
	}

	// 获取可用渠道
	routerInfo := GetChannelRouterByMerchant(payin.ToTransaction())
	if routerInfo == nil {
		code = protocol.ChannelNotFound
		return
	}

	// 设置渠道信息
	payin.SetChannelGroup(routerInfo.ChannelGroup)
	values := models.NewTrxValues()
	er := models.WriteDB.Transaction(func(tx *gorm.DB) error {
		// 创建代收订单
		if err := tx.Create(payin).Error; err != nil {
			return err
		}
		trx := payin.ToTransaction()
		// 执行渠道请求
		result, err := RequestByRouter(ctx, tx, trx, routerInfo)
		if err != nil {
			return err
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
		if result.Status == protocol.StatusFailed {
			values.SetCompletedAt(time.Now().UnixMilli())
		}
		return nil
	})
	if er != nil {
		code = protocol.InternalError
		return
	}
	models.SaveTransactionValues(models.WriteDB, payin.ToTransaction(), values)
	AfterTransactionCreate(payin.ToTransaction())
	return
}
