package models

import (
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Account 统一账户表
// 时间字段均为int64，结尾为_at
// UserID + UserType + Currency 唯一
// Asset为JSON字段，包含多资金属性
type Account struct {
	ID        uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	AccountID string `json:"account_id" gorm:"column:account_id;type:varchar(32);uniqueIndex"`
	Salt      string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	UserID    string `json:"user_id" gorm:"column:user_id;type:varchar(32);uniqueIndex:uk_userid_usertype_ccy"`
	UserType  string `json:"user_type" gorm:"column:user_type;type:varchar(16);uniqueIndex:uk_userid_usertype_ccy"`
	Ccy       string `json:"ccy" gorm:"column:ccy;type:varchar(16);uniqueIndex:uk_userid_usertype_ccy"`
	*AccountValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type AccountValues struct {
	Asset        *Asset  `json:"asset" gorm:"column:asset;serializer:json;type:json"`
	Status       *string `json:"status" gorm:"column:status;type:varchar(16);default:''"`
	Version      *int64  `json:"version" gorm:"column:version;type:bigint;default:1"`
	LastActiveAt *int64  `json:"last_active_at" gorm:"column:last_active_at;type:bigint"`
}

func (Account) TableName() string {
	return "t_accounts"
}

type Accounts []*Account

func (t Accounts) Protocol() []*protocol.Account {
	result := make([]*protocol.Account, 0, len(t))
	for _, account := range t {
		result = append(result, account.Protocol())
	}
	return result
}

// Asset 资产模型，支持多资金属性
type Asset struct {
	Balance                decimal.Decimal `json:"balance"`                  // 总余额
	AvailableBalance       decimal.Decimal `json:"available_balance"`        // 可用余额
	FrozenBalance          decimal.Decimal `json:"frozen_balance"`           // 冻结余额
	MarginBalance          decimal.Decimal `json:"margin_balance"`           // 保证金余额
	AvailableMarginBalance decimal.Decimal `json:"available_margin_balance"` // 可用保证金余额
	FrozenMarginBalance    decimal.Decimal `json:"frozen_margin_balance"`    // 冻结保证金余额
	Ccy                    string          `json:"ccy"`                      // 币种
	UpdatedAt              int64           `json:"updated_at"`               // 更新时间
}

// NewAccount 创建新账户
func NewAccount() *Account {
	return &Account{
		AccountValues: &AccountValues{},
	}
}

func (t *Asset) Protocol() *protocol.Assert {
	if t == nil {
		return nil
	}
	return &protocol.Assert{
		Balance:                t.Balance.String(),
		AvailableBalance:       t.Balance.Sub(t.FrozenBalance).String(),
		FrozenBalance:          t.FrozenBalance.String(),
		MarginBalance:          t.MarginBalance.String(),
		AvailableMarginBalance: t.MarginBalance.Sub(t.FrozenMarginBalance).String(),
		FrozenMarginBalance:    t.FrozenMarginBalance.String(),
		Ccy:                    t.Ccy,
		UpdatedAt:              t.UpdatedAt,
	}
}

func (t *Account) Protocol() *protocol.Account {
	if t == nil {
		return nil
	}
	info := &protocol.Account{
		AccountID:    t.AccountID,
		UserID:       t.UserID,
		UserType:     t.UserType,
		Ccy:          t.Ccy,
		Balance:      t.Asset.Protocol(),
		Status:       t.GetStatus(),
		Version:      t.GetVersion(),
		LastActiveAt: t.GetLastActiveAt(),
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
	if t.Asset != nil {
		info.Balance = t.Asset.Protocol()
	}
	return info
}

// SetValues 设置AccountValues
func (a *AccountValues) SetValues(values *AccountValues) {
	if values == nil {
		return
	}
	if values.Asset != nil {
		a.Asset = values.Asset
	}
	if values.Status != nil {
		a.Status = values.Status
	}
	if values.Version != nil {
		a.Version = values.Version
	}
	if values.LastActiveAt != nil {
		a.LastActiveAt = values.LastActiveAt
	}
}

func (a *AccountValues) GetStatus() string {
	if a.Status == nil {
		return ""
	}
	return *a.Status
}

func (a *AccountValues) GetVersion() int64 {
	if a.Version == nil {
		return 0
	}
	return *a.Version
}

func (a *AccountValues) GetLastActiveAt() int64 {
	if a.LastActiveAt == nil {
		return 0
	}
	return *a.LastActiveAt
}

func (a *AccountValues) SetAsset(asset *Asset) *AccountValues {
	a.Asset = asset
	return a
}

func (a *AccountValues) SetStatus(status string) *AccountValues {
	a.Status = &status
	return a
}

func (a *AccountValues) SetVersion(version int64) *AccountValues {
	a.Version = &version
	return a
}

func (a *AccountValues) SetLastActiveAt(lastActiveAt int64) *AccountValues {
	a.LastActiveAt = &lastActiveAt
	return a
}

// 查询方法
func GetAccountByID(id uint64) (*Account, error) {
	var account Account
	err := WriteDB.Where("id = ?", id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func GetAccountByAccountID(accountID string) (*Account, error) {
	var account Account
	err := WriteDB.Where("account_id = ?", accountID).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func GetAccountByUserIDAndCurrency(userID, userType, currency string) (*Account, error) {
	var account Account
	err := WriteDB.Where("user_id = ? AND user_type = ? AND currency = ?", userID, userType, currency).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func GetAccountForUpdate(tx *gorm.DB, userID, userType, currency string) (*Account, error) {
	var account Account
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND user_type = ? AND currency = ?", userID, userType, currency).
		First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func GetAccountsByUserID(userID, userType string) Accounts {
	var accounts Accounts
	err := WriteDB.Where("user_id = ? AND user_type = ?", userID, userType).Find(&accounts).Error
	if err != nil {
		return nil
	}
	return accounts
}

func GetAccountsByIDs(ids []string) ([]*Account, error) {
	var accounts []*Account
	err := WriteDB.Where("account_id IN ?", ids).Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
