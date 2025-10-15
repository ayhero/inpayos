package services

import (
	"fmt"
	"inpayos/internal/log"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var (
	accountServiceInstance *AccountService
	accountServiceOnce     sync.Once
)

type AccountService struct{}

// GetAccountService 获取账户服务单例
func GetAccountService() *AccountService {
	if accountServiceInstance == nil {
		SetupAccountService()
	}
	return accountServiceInstance
}

func SetupAccountService() {
	accountServiceOnce.Do(func() {
		accountServiceInstance = &AccountService{}
	})
}

// CreateAccount 创建账户
func (s *AccountService) CreateAccount(req *protocol.CreateAccountRequest) (*models.Account, error) {
	// 检查账户是否已存在
	existingAccount, err := models.GetAccountByUserIDAndCurrency(req.UserID, req.UserType, req.Ccy)
	if err == nil && existingAccount != nil {
		return nil, fmt.Errorf("account already exists for user %s with type %s and currency %s", req.UserID, req.UserType, req.Ccy)
	}

	// 创建新账户
	account := models.NewAccount()
	account.AccountID = utils.GenerateAccountID()
	account.UserID = req.UserID
	account.UserType = req.UserType
	account.Ccy = req.Ccy
	account.SetStatus(protocol.StatusActive).
		SetVersion(1).
		SetLastActiveAt(time.Now().UnixMilli())

	// 初始化资产
	asset := &models.Asset{
		Balance:                decimal.Zero,
		AvailableBalance:       decimal.Zero,
		FrozenBalance:          decimal.Zero,
		MarginBalance:          decimal.Zero,
		AvailableMarginBalance: decimal.Zero,
		FrozenMarginBalance:    decimal.Zero,
		Ccy:                    req.Ccy,
		UpdatedAt:              time.Now().UnixMilli(),
	}
	account.AccountValues.SetAsset(asset)

	// 保存到数据库
	err = models.WriteDB.Create(account).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

func (s AccountService) GetMerchantAccountBalance(merchantID string) (balance []*protocol.Account) {
	acct := models.GetAccountsByUserID(merchantID, protocol.UserTypeMerchant)
	if acct == nil {
		return []*protocol.Account{}
	}
	return acct.Protocol()
}

// UpdateBalance 更新账户余额
func (s *AccountService) UpdateBalance(req *protocol.UpdateBalanceRequest) (err_code protocol.ErrorCode) {

	direction, ok := protocol.AccountDirectionMap[req.TrxType]
	if !ok {
		return protocol.AccountErrorInvalidTrxType
	}
	err := models.WriteDB.Transaction(func(tx *gorm.DB) error {
		// 锁定账户
		account, err := models.GetAccountForUpdate(tx, req.UserID, req.UserType, req.Ccy)
		if err != nil {
			err_code = protocol.AccountErrorAccountNotFound
			return err
		}
		// 初始化资产如果为空
		if account.Asset == nil {
			account.Asset = &models.Asset{
				Balance:             decimal.Zero,
				FrozenBalance:       decimal.Zero,
				MarginBalance:       decimal.Zero,
				FrozenMarginBalance: decimal.Zero,
				Ccy:                 req.Ccy,
				UpdatedAt:           time.Now().UnixMilli(),
			}
		}
		// 记录操作前余额
		beforeAssert := *account.Asset
		afterAssert := account.Asset
		afterAssert.UpdatedAt = utils.TimeNowMilli()

		values := &models.AccountValues{
			Asset: afterAssert,
		}
		// 执行余额操作
		switch req.TrxType {
		case protocol.TrxTypePayin:
			afterAssert.Balance = afterAssert.Balance.Add(req.Amount)
		case protocol.TrxTypePayout:
			if afterAssert.Balance.LessThan(req.Amount) {
				err_code = protocol.AccountErrorInsufficientBalance
				return fmt.Errorf("insufficient available balance")
			}
			direction = protocol.DirectionOut
			afterAssert.Balance = afterAssert.Balance.Sub(req.Amount)
			afterAssert.Balance = afterAssert.Balance.Sub(req.Amount)
		case protocol.TrxTypeFreeze:
			if afterAssert.Balance.LessThan(req.Amount) {
				err_code = protocol.AccountErrorInsufficientBalance
				return fmt.Errorf("insufficient available balance to freeze")
			}
			direction = protocol.DirectionOut
			afterAssert.Balance = afterAssert.Balance.Sub(req.Amount)
			afterAssert.FrozenBalance = afterAssert.FrozenBalance.Add(req.Amount)
		case protocol.TrxTypeUnfreeze:
			if afterAssert.FrozenBalance.LessThan(req.Amount) {
				err_code = protocol.AccountErrorInsufficientFrozenBalance
				return fmt.Errorf("insufficient frozen balance to unfreeze")
			}
			afterAssert.FrozenBalance = afterAssert.FrozenBalance.Sub(req.Amount)
			afterAssert.Balance = afterAssert.Balance.Add(req.Amount)
		case protocol.TrxTypeMarginDeposit:
			if afterAssert.Balance.LessThan(req.Amount) {
				err_code = protocol.AccountErrorInsufficientBalance
				return fmt.Errorf("insufficient available balance for margin")
			}
			afterAssert.MarginBalance = afterAssert.MarginBalance.Add(req.Amount)
		case protocol.TrxTypeMarginRelease:
			if afterAssert.MarginBalance.LessThan(req.Amount) {
				err_code = protocol.AccountErrorInsufficientMarginBalance
				return fmt.Errorf("insufficient margin balance to release")
			}
			direction = protocol.DirectionOut
			afterAssert.MarginBalance = afterAssert.MarginBalance.Sub(req.Amount)
		default:
			err_code = protocol.AccountErrorUnsupportedTrxType
			return fmt.Errorf("unsupported trx_type: %s", req.TrxType)
		}

		// 更新账户信息
		afterAssert.AvailableBalance = afterAssert.Balance.Sub(afterAssert.FrozenBalance)
		values.SetVersion(account.GetVersion() + 1)

		// 创建资金流水记录
		fundFlow := &models.FundFlow{
			FlowNo:         utils.GenerateFlowNo(),
			Direction:      direction,
			UserID:         req.UserID,
			UserType:       req.UserType,
			AccountID:      account.AccountID,
			AccountVersion: account.GetVersion(),
			TrxID:          req.TrxID,
			TrxType:        req.TrxType,
			Ccy:            req.Ccy,
			Amount:         &req.Amount,
			BeforeAsset:    &beforeAssert,
			AfterAsset:     afterAssert,
			CreatedAt:      utils.TimeNowMilli(),
		}
		if err := tx.Create(fundFlow).Error; err != nil {
			return err
		}
		// 保存账户
		if err := tx.Model(account).UpdateColumns(values).Error; err != nil {
			err_code = protocol.AccountErrorUpdateFailed
			return err
		}

		return nil
	})
	if err != nil {
		log.Get().Error("UpdateBalance error:", err)
	}
	return protocol.Success
}

// GetAccountList 获取账户列表
func (s *AccountService) GetAccountList(userID, userType string) ([]*protocol.Account, protocol.ErrorCode) {
	// 查询账户列表
	accounts := models.GetAccountsByUserID(userID, userType)
	// 转换为响应格式
	return accounts.Protocol(), protocol.Success
}

// ListAccountFlowByQuery 根据查询条件获取账户流水列表
func (s *AccountService) ListAccountFlowByQuery(userID, userType string, req *protocol.AccountFlowListRequest) ([]*models.FundFlow, int64, protocol.ErrorCode) {
	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	size := req.Size
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	// 构建查询条件
	query := models.WriteDB.Model(&models.FundFlow{}).
		Where("user_id = ? AND user_type = ?", userID, userType)

	// 添加可选查询条件
	if req.FlowNo != "" {
		query = query.Where("flow_no = ?", req.FlowNo)
	}
	if req.TrxID != "" {
		query = query.Where("trx_id = ?", req.TrxID)
	}
	if req.TrxType != "" {
		query = query.Where("trx_type = ?", req.TrxType)
	}
	if req.Direction != "" {
		query = query.Where("direction = ?", req.Direction)
	}
	if req.Ccy != "" {
		query = query.Where("ccy = ?", req.Ccy)
	}
	if req.CreatedAtStart > 0 {
		query = query.Where("created_at >= ?", req.CreatedAtStart)
	}
	if req.CreatedAtEnd > 0 {
		query = query.Where("created_at <= ?", req.CreatedAtEnd)
	}

	// 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, protocol.DatabaseError
	}

	// 分页查询
	var flows []*models.FundFlow
	offset := (page - 1) * size
	if err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&flows).Error; err != nil {
		return nil, 0, protocol.DatabaseError
	}

	return flows, total, protocol.Success
}
