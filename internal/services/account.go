package services

import (
	"fmt"
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
	account.AccountValues.SetUserID(req.UserID).
		SetUserType(req.UserType).
		SetCcy(req.Ccy).
		SetStatus(1).
		SetVersion(1).
		SetLastActiveAt(time.Now().UnixMilli())

	// 初始化资产
	asset := &models.Asset{
		Balance:          decimal.Zero,
		AvailableBalance: decimal.Zero,
		FrozenBalance:    decimal.Zero,
		MarginBalance:    decimal.Zero,
		ReserveBalance:   decimal.Zero,
		Ccy:              req.Ccy,
		UpdatedAt:        time.Now().UnixMilli(),
	}
	account.AccountValues.SetAsset(asset)

	// 保存到数据库
	err = models.WriteDB.Create(account).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

func (s AccountService) GetMerchantAccountBalance(merchantID string) (balance *protocol.Balance, code protocol.ErrorCode) {
	return &protocol.Balance{}, protocol.Success
}

// GetBalance 获取账户余额
func (s *AccountService) GetBalance(userID, userType, currency string) (*protocol.Balance, error) {
	account, err := models.GetAccountByUserIDAndCurrency(userID, userType, currency)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	if account.Asset == nil {
		return &protocol.Balance{
			Balance:          "0",
			AvailableBalance: "0",
			FrozenBalance:    "0",
			MarginBalance:    "0",
			ReserveBalance:   "0",
			Currency:         currency,
			UpdatedAt:        account.UpdatedAt,
		}, nil
	}

	return &protocol.Balance{
		Balance:          account.Asset.Balance.String(),
		AvailableBalance: account.Asset.AvailableBalance.String(),
		FrozenBalance:    account.Asset.FrozenBalance.String(),
		MarginBalance:    account.Asset.MarginBalance.String(),
		ReserveBalance:   account.Asset.ReserveBalance.String(),
		Currency:         account.Asset.Ccy,
		UpdatedAt:        account.Asset.UpdatedAt,
	}, nil
}

// UpdateBalance 更新账户余额
func (s *AccountService) UpdateBalance(req *protocol.UpdateBalanceRequest) error {
	return models.WriteDB.Transaction(func(tx *gorm.DB) error {
		// 锁定账户
		account, err := models.GetAccountForUpdate(tx, req.UserID, req.UserType, req.Ccy)
		if err != nil {
			return fmt.Errorf("account not found: %w", err)
		}

		// 初始化资产如果为空
		if account.Asset == nil {
			account.Asset = &models.Asset{
				Balance:          decimal.Zero,
				AvailableBalance: decimal.Zero,
				FrozenBalance:    decimal.Zero,
				MarginBalance:    decimal.Zero,
				ReserveBalance:   decimal.Zero,
				Ccy:              req.Ccy,
				UpdatedAt:        time.Now().UnixMilli(),
			}
		}

		// 记录操作前余额
		//beforeBalance := account.Asset.Balance

		// 执行余额操作
		switch req.Operation {
		case "add":
			account.Asset.Balance = account.Asset.Balance.Add(req.Amount)
			account.Asset.AvailableBalance = account.Asset.AvailableBalance.Add(req.Amount)
		case "subtract":
			if account.Asset.AvailableBalance.LessThan(req.Amount) {
				return fmt.Errorf("insufficient available balance")
			}
			account.Asset.Balance = account.Asset.Balance.Sub(req.Amount)
			account.Asset.AvailableBalance = account.Asset.AvailableBalance.Sub(req.Amount)
		case "freeze":
			if account.Asset.AvailableBalance.LessThan(req.Amount) {
				return fmt.Errorf("insufficient available balance to freeze")
			}
			account.Asset.AvailableBalance = account.Asset.AvailableBalance.Sub(req.Amount)
			account.Asset.FrozenBalance = account.Asset.FrozenBalance.Add(req.Amount)
		case "unfreeze":
			if account.Asset.FrozenBalance.LessThan(req.Amount) {
				return fmt.Errorf("insufficient frozen balance to unfreeze")
			}
			account.Asset.FrozenBalance = account.Asset.FrozenBalance.Sub(req.Amount)
			account.Asset.AvailableBalance = account.Asset.AvailableBalance.Add(req.Amount)
		case "margin":
			if account.Asset.AvailableBalance.LessThan(req.Amount) {
				return fmt.Errorf("insufficient available balance for margin")
			}
			account.Asset.AvailableBalance = account.Asset.AvailableBalance.Sub(req.Amount)
			account.Asset.MarginBalance = account.Asset.MarginBalance.Add(req.Amount)
		case "release_margin":
			if account.Asset.MarginBalance.LessThan(req.Amount) {
				return fmt.Errorf("insufficient margin balance to release")
			}
			account.Asset.MarginBalance = account.Asset.MarginBalance.Sub(req.Amount)
			account.Asset.AvailableBalance = account.Asset.AvailableBalance.Add(req.Amount)
		default:
			return fmt.Errorf("unsupported operation: %s", req.Operation)
		}

		// 更新账户信息
		account.Asset.UpdatedAt = time.Now().UnixMilli()
		account.AccountValues.SetLastActiveAt(time.Now().UnixMilli())
		account.AccountValues.SetVersion(account.GetVersion() + 1)

		// 保存账户
		if err := tx.Save(account).Error; err != nil {
			return fmt.Errorf("failed to update account: %w", err)
		}

		// 创建资金流水记录
		fundFlow := &models.FundFlow{}
		fundFlow.FlowNo = utils.GenerateFlowID()
		/*
			fundFlow.SetUserID(req.UserID).
				SetUserType(req.UserType).
				SetAccountID(account.AccountID).
				SetTransactionID(req.TransactionID).
				SetBillID(req.BillID).
				SetFlowType(req.Operation).
				SetAmount(req.Amount).
				SetCurrency(req.Currency).
				SetBeforeBalance(beforeBalance).
				SetAfterBalance(account.Asset.Balance).
				SetBusinessType(req.BusinessType).
				SetDescription(req.Description).
				SetFlowAt(time.Now().UnixMilli())
		*/
		if err := tx.Create(fundFlow).Error; err != nil {
			return fmt.Errorf("failed to create fund flow: %w", err)
		}

		return nil
	})
}
