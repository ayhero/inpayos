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

	// 3. 检查是否已有该支付方式的交易记录（从Checkout.Transactions字段中检查）
	existingTrx := checkout.FindTransactionByTrxMethod(req.TrxMethod)
	if existingTrx != nil {
		// 如果已存在该支付方式的交易，直接返回
		trx = checkout.Protocol()
		trx.Transaction = existingTrx.Protocol()
		return
	}

	// 4. 创建交易记录（暂存在Checkout.Transactions中，不写入代收表）
	reqID := checkout.CheckoutID + "-" + req.TrxMethod
	trxID := utils.GeneratePayinID()

	// 创建Transaction记录
	status := protocol.StatusPending
	notifyURL := checkout.GetNotifyURL()
	transaction := &models.Transaction{
		Mid:       req.Mid,
		TrxType:   protocol.TrxTypePayin,
		ReqID:     reqID,
		TrxID:     trxID,
		Ccy:       checkout.GetCcy(),
		Amount:    checkout.Amount,
		TrxMethod: req.TrxMethod,
		ReturnURL: checkout.GetReturnURL(),
		AccountNo: req.AccountNo,
		TransactionValues: &models.TransactionValues{
			Status:    &status,
			NotifyURL: &notifyURL,
		},
	}
	transaction.SetCountry(req.Country)
	// 5. 将交易记录添加到Checkout.Transactions字段中
	checkout.AddTransaction(transaction)

	nowtime := utils.TimeNowMilli()
	// 6. 更新收银台记录（包含新的交易数据）
	checkoutValues := &models.MerchantCheckoutValues{}
	checkoutValues.SetSubmitedAt(nowtime).
		SetTransactions(checkout.GetTransactions())
	err := models.SaveMerchantCheckout(models.WriteDB, checkout, checkoutValues)
	if err != nil {
		log.Get().Errorf("Update checkout submit status failed: %v", err)
		code = protocol.SystemError
		return
	}

	// 7. 返回更新后的收银台信息
	trx = checkout.Protocol()
	trx.Transaction = transaction.Protocol()
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
		code = protocol.TransactionNotFound
		return
	}

	// 验证收银台状态
	if checkout.GetStatus() != protocol.StatusPending {
		code = protocol.TransactionCompleted
		return
	}

	// 3. 从Checkout.Transactions中查找对应的交易记录
	transaction := checkout.FindTransactionByTrxID(req.TrxID)
	if transaction == nil {
		code = protocol.TransactionNotFound
		return
	}

	// 验证交易状态
	if transaction.GetStatus() != protocol.StatusPending {
		code = protocol.TransactionCompleted
		return
	}
	// 4. 现在才将交易数据写入代收表
	now := utils.TimeNowMilli()
	//更新交易记录
	transaction.SetStatus(protocol.StatusConfirming).
		SetSubmitedAt(now)

	// 更新收银台状态
	checkoutValues := &models.MerchantCheckoutValues{}
	checkoutValues.SetStatus(protocol.StatusConfirming).
		SetTrxID(transaction.TrxID).
		SetCountry(transaction.GetCountry()).
		SetTrxMethod(transaction.TrxMethod).
		SetSubmitedAt(now).
		SetTransactions(checkout.GetTransactions())

	// 使用数据库事务确保数据一致性
	err := models.WriteDB.Transaction(func(tx *gorm.DB) error {
		// 先检查代收表中是否已存在该交易
		if trx := models.GetTransactionByMidAndTrxID(req.Mid, transaction.TrxID, transaction.TrxType); trx == nil {
			// 保存代收记录到数据库
			if err := tx.Table(models.TrxTypeTableMap[transaction.TrxType]).Create(transaction.ToTrxByType()).Error; err != nil {
				return err
			}
		}
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
	routers := models.ListActiveRouterByMerchant(checkout.Mid, protocol.TrxTypePayin)
	if len(routers) == 0 {
		code = protocol.ChannelNotSupported
		return
	}

	// 4. 创建商户收银台配置结构
	result = &protocol.MerchantCheckoutConfig{
		Mid:       checkout.Mid,
		Countries: []string{},
		Configs:   make(map[string]*protocol.CountryCheckoutConfig),
	}

	// 5. 按国家分组处理路由配置
	countryMap := make(map[string]*protocol.CountryCheckoutConfig)

	for _, router := range routers {
		trxMethod := router.GetTrxMethod()
		if trxMethod == "" {
			continue
		}

		country := router.GetCountry()
		if country == "" {
			continue
		}

		// 如果国家配置不存在，创建新的
		if _, exists := countryMap[country]; !exists {
			countryMap[country] = &protocol.CountryCheckoutConfig{
				Country:    country,
				TrxMethods: []string{},
				Configs:    map[string]*protocol.TrxMethodConfig{},
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
