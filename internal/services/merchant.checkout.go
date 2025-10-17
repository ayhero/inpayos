package services

import (
	"inpayos/internal/config"
	"inpayos/internal/log"
	"inpayos/internal/middleware"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
	"slices"
	"sync"
	"time"

	"gorm.io/gorm"
)

type CheckoutService struct {
}

var (
	checkoutService     *CheckoutService
	checkoutServiceOnce sync.Once
)

func SetupCheckoutService() {
	checkoutServiceOnce.Do(func() {
		checkoutService = &CheckoutService{}
	})
}

// GetCheckoutService 获取Checkout服务单例
func GetCheckoutService() *CheckoutService {
	if checkoutService == nil {
		SetupCheckoutService()
	}
	return checkoutService
}

// CreateCheckout 创建收银台会话
func (s *CheckoutService) Create(req *protocol.CreateCheckoutRequest) (info *protocol.Checkout, code protocol.ErrorCode) {
	code = protocol.Success
	// 2. 验证金额
	amt, err := utils.ValidateAmount(req.Amount)
	if err != nil {
		return nil, protocol.InvalidParams
	}

	// 3. 检查重复请求
	checkout := models.GetMerchantCheckoutByReqID(req.Mid, req.ReqID)
	if checkout != nil {
		return nil, protocol.DuplicateTransaction
	}

	// 4. 生成收银台ID和过期时间
	checkoutID := utils.GenerateCheckoutID()
	ckCfg := config.Get().MerchantCheckout
	expiredAt := time.Now().Add(time.Duration(ckCfg.ExpiryMinutes) * time.Minute).UnixMilli()

	// 6. 创建收银台记录
	checkout = &models.MerchantCheckout{
		CheckoutID: checkoutID,
		Mid:        req.Mid,
		ReqID:      req.ReqID,
		TrxID:      utils.GeneratePayinID(),
		TrxType:    protocol.TrxTypePayin, // 默认代收类型
		MerchantCheckoutValues: &models.MerchantCheckoutValues{
			Ccy:          &req.Ccy,
			Amount:       amt,
			TrxMethod:    &req.TrxMethod,
			ExpiredAt:    &expiredAt,
			Country:      &req.Country,
			NotifyURL:    &req.NotifyURL,
			ReturnURL:    &req.ReturnURL,
			Metadata:     &protocol.MapData{},
			Transactions: []*models.Transaction{},
		},
	}
	checkout.SetStatus(protocol.StatusPending)

	// 7. 保存到数据库
	dbErr := models.WriteDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(checkout).Error; err != nil {
			return err
		}
		return nil
	})

	if dbErr != nil {
		return nil, protocol.SystemError
	}

	// 8. 返回协议格式的响应
	info = checkout.Protocol()
	return
}

func (s *CheckoutService) Submit(req *protocol.SubmitCheckoutRequest) (trx *protocol.Checkout, code protocol.ErrorCode) {
	code = protocol.Success
	// 2. 获取收银台记录
	checkout := models.GetMerchantCheckoutByCheckoutID(req.CheckoutID)
	if checkout == nil || checkout.Mid != req.Mid {
		code = protocol.TransactionNotFound
		return
	}
	// 验证收银台状态
	if checkout.GetStatus() != protocol.StatusPending {
		code = protocol.TransactionCompleted
		return
	}

	// 3. 检查是否已有该支付方式的交易记录
	// 使用checkout的CheckoutID + TrxMethod作为ReqID，确保同一收银台的同一支付方式唯一
	reqID := checkout.CheckoutID + "-" + req.TrxMethod

	// 检查是否已存在相同支付方式的交易
	payin := models.GetMerchantPayinByReqID(req.Mid, reqID)
	if payin != nil {
		trx = checkout.Protocol()
		trx.Transaction = payin.ToTransaction().Protocol()
		return
	}

	// 4. 创建支付交易请求
	// ReqID在第3步已经定义，这里直接使用

	payinReq := &protocol.MerchantPayinRequest{
		Mid:       req.Mid,
		ReqID:     reqID, // 设置ReqID用于重复检查
		Ccy:       checkout.GetCcy(),
		Amount:    checkout.GetAmount().String(),
		TrxMethod: req.TrxMethod,
		TrxMode:   "", // 基于 TrxMethod 自动推导
		TrxApp:    "", // 在具体支付时再指定
		ReturnURL: checkout.GetReturnURL(),
		NotifyURL: checkout.GetNotifyURL(),
	}

	// 5. 创建支付订单
	var transaction *protocol.Transaction

	// 使用数据库事务
	err := models.WriteDB.Transaction(func(tx *gorm.DB) error {
		// 先检查重复性
		payin := models.GetMerchantPayinByReqID(payinReq.Mid, payinReq.ReqID)
		if payin != nil {
			transaction = payin.ToTransaction().Protocol()
			return nil
		}
		// 创建代收记录
		payin = &models.MerchantPayin{
			Mid:                 payinReq.Mid,
			TrxType:             protocol.TrxTypePayin,
			ReqID:               payinReq.ReqID,
			TrxID:               utils.GeneratePayinID(),
			Ccy:                 payinReq.Ccy,
			Amount:              checkout.Amount,
			TrxMethod:           payinReq.TrxMethod,
			ReturnURL:           payinReq.ReturnURL,
			AccountNo:           req.AccountNo,
			MerchantPayinValues: &models.MerchantPayinValues{},
		}
		payin.SetStatus(protocol.StatusPending).
			SetNotifyURL(payinReq.NotifyURL)

		// 保存到数据库先创建基础记录
		return tx.Create(payin).Error
	})

	if err != nil {
		code = protocol.SystemError
		return
	}
	nowtime := utils.TimeNowMilli()
	// 6. 更新收银台记录
	checkoutValues := &models.MerchantCheckoutValues{}
	checkoutValues.SetSubmitedAt(nowtime)

	err = models.SaveMerchantCheckout(models.WriteDB, checkout, checkoutValues)
	if err != nil {
		// 更新失败不影响主流程，记录日志即可
		log.Get().Errorf("Update checkout submit status failed: %v", err)
	}

	// 7. 返回更新后的收银台信息
	trx = checkout.Protocol()
	// 这里还需要包含新创建的交易信息
	trx.Transaction = transaction
	return
}

func (s *CheckoutService) Confirm(req *protocol.ConfirmCheckoutRequest) (trx *protocol.Checkout, code protocol.ErrorCode) {
	code = protocol.Success

	// 2. 获取收银台记录
	checkout := models.GetMerchantCheckoutByCheckoutID(req.CheckoutID)
	if checkout == nil {
		code = protocol.TransactionNotFound
		return
	}

	// 验证商户匹配
	if checkout.Mid != req.Mid {
		code = protocol.InvalidParams
		return
	}

	// 验证收银台状态
	if checkout.GetStatus() != protocol.StatusPending {
		code = protocol.TransactionCompleted
		return
	}

	// 3. 获取对应的交易记录
	var transaction *models.Transaction

	// 根据TrxID查找对应的交易记录
	if checkout.TrxType == protocol.TrxTypePayin {
		transaction = models.GetMerchantPayinByTrxID(req.Mid, req.TrxID)
		if transaction == nil {
			code = protocol.TransactionNotFound
			return
		}
	}
	if transaction == nil {
		code = protocol.TransactionNotFound
		return
	}

	// 验证交易状态
	if transaction.GetStatus() != protocol.StatusSubmitted {
		code = protocol.TransactionCompleted
		return
	}

	// 4. 更新交易状态为确认中
	now := time.Now().UnixMilli()
	// 使用数据库事务确保数据一致性
	err := models.WriteDB.Transaction(func(tx *gorm.DB) error {
		// 更新交易记录状态
		values := models.NewTrxValues()
		values.SetStatus(protocol.StatusConfirming).
			SetConfirmedAt(now).
			SetProofID(req.ProofID)

		// 保存交易状态更新
		if err := models.SaveTransactionValues(tx, transaction, values); err != nil {
			return err
		}

		// 更新收银台状态
		checkoutValues := &models.MerchantCheckoutValues{}
		checkoutValues.SetStatus(protocol.StatusConfirming).
			SetConfirmedAt(now)

		return models.SaveMerchantCheckout(tx, checkout, checkoutValues)
	})

	if err != nil {
		log.Get().Errorf("Confirm checkout failed: %v", err)
		code = protocol.SystemError
		return
	}
	trx = checkout.Protocol()
	return
}

func (s *CheckoutService) Configs(checkoutID string) (result *protocol.MerchantCheckoutConfig, code protocol.ErrorCode) {
	code = protocol.Success

	// 1. 参数验证
	if checkoutID == "" {
		code = protocol.InvalidParams
		return
	}

	// 2. 获取收银台记录
	checkout := models.GetMerchantCheckoutByCheckoutID(checkoutID)
	if checkout == nil {
		code = protocol.TransactionNotFound
		return
	}

	// 3. 获取商户的支付路由配置
	routers := models.ListMerchantRouterByMerchant(checkout.Mid, protocol.TrxTypePayin)
	if len(routers) == 0 {
		code = protocol.ChannelNotSupported
		return
	}

	// 4. 创建商户收银台配置结构
	result = &protocol.MerchantCheckoutConfig{
		MerchantID: checkout.Mid,
		Countries:  []string{},
		Configs:    make(map[string]*protocol.CountryCheckoutConfig),
	}

	// 5. 按国家分组处理路由配置
	countryMap := make(map[string]*protocol.CountryCheckoutConfig)

	for _, router := range routers {
		// 跳过非活跃状态的路由
		if router.GetStatus() != protocol.StatusActive {
			continue
		}

		trxMethod := router.GetTrxMethod()
		if trxMethod == "" {
			continue
		}

		country := router.GetCountry()
		if country == "" {
			country = "*" // 默认国家
		}

		// 如果国家配置不存在，创建新的
		if _, exists := countryMap[country]; !exists {
			countryMap[country] = &protocol.CountryCheckoutConfig{
				Country:    country,
				TrxMethods: []string{},
				Configs:    make(map[string]*protocol.TrxMethodConfig),
			}
			result.Countries = append(result.Countries, country)
		}

		countryConfig := countryMap[country]

		// 添加交易方式到国家配置
		if !slices.Contains(countryConfig.TrxMethods, trxMethod) {
			countryConfig.TrxMethods = append(countryConfig.TrxMethods, trxMethod)
		}

		// 创建或更新交易方式配置
		if methodConfig, exists := countryConfig.Configs[trxMethod]; exists {
			// 合并币种支持
			ccy := router.GetCcy()
			if ccy != "" && !slices.Contains(methodConfig.Ccy, ccy) {
				methodConfig.Ccy = append(methodConfig.Ccy, ccy)
			}
		} else {
			// 创建新的交易方式配置
			methodConfig := &protocol.TrxMethodConfig{
				Country:   country,
				TrxMethod: trxMethod,
				Ccy:       []string{},
			}

			// 添加币种支持
			ccy := router.GetCcy()
			if ccy != "" {
				methodConfig.Ccy = append(methodConfig.Ccy, ccy)
			}

			countryConfig.Configs[trxMethod] = methodConfig
		}
	}

	// 6. 将国家配置添加到结果中
	result.Configs = countryMap

	return
}

func (s *CheckoutService) Info(checkoutID string) (trx *protocol.Checkout, code protocol.ErrorCode) {
	if checkoutID == "" {
		return nil, protocol.InvalidParams
	}
	checkout := models.GetMerchantCheckoutByCheckoutID(checkoutID)
	if checkout == nil {
		return nil, protocol.TransactionNotFound
	}
	if checkout.GetExpiredAt() < utils.TimeNowMilli() {
		return nil, protocol.TransactionExpired
	}
	expiryTime := time.UnixMilli(checkout.GetExpiredAt())
	token, err := middleware.GenerateToken(checkout.Mid, expiryTime, config.Get().Server.Merchant.Jwt.Secret)
	if err != nil {
		return nil, protocol.TransactionNotFound
	}
	trx = checkout.Protocol()
	trx.Token = token
	return trx, protocol.Success
}

func (s *CheckoutService) Cancel(checkoutID string) (trx *protocol.Checkout, code protocol.ErrorCode) {
	if checkoutID == "" {
		return nil, protocol.InvalidParams
	}

	checkout := models.GetMerchantCheckoutByCheckoutID(checkoutID)
	if checkout == nil {
		return nil, protocol.TransactionNotFound
	}

	if checkout.GetStatus() == protocol.StatusCancelled {
		return checkout.Protocol(), protocol.Success
	}

	// 检查是否可以取消
	if !checkout.CanCancel() {
		return nil, protocol.TransactionCompleted
	}

	// 更新状态为已取消
	values := &models.MerchantCheckoutValues{}
	values.SetStatus(protocol.StatusCancelled).
		SetCanceledAt(utils.TimeNowMilli())

	if _err := models.SaveMerchantCheckout(models.WriteDB, checkout, values); _err != nil {
		return nil, protocol.SystemError
	}

	return checkout.Protocol(), protocol.Success
}
