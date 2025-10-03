package models

import (
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
	*AccountValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type AccountValues struct {
	UserID       *string `json:"user_id" gorm:"column:user_id;type:varchar(32);uniqueIndex:uk_userid_usertype_ccy"`
	UserType     *string `json:"user_type" gorm:"column:user_type;type:varchar(16);uniqueIndex:uk_userid_usertype_ccy"`
	Currency     *string `json:"currency" gorm:"column:currency;type:varchar(16);uniqueIndex:uk_userid_usertype_ccy"`
	Asset        *Asset  `json:"asset" gorm:"column:asset;serializer:json;type:json"`
	Status       *int    `json:"status" gorm:"column:status;type:int;default:1"`
	Version      *int64  `json:"version" gorm:"column:version;type:bigint;default:1"`
	LastActiveAt *int64  `json:"last_active_at" gorm:"column:last_active_at;type:bigint"`

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Account) TableName() string {
	return "t_accounts"
}

// Asset 资产模型，支持多资金属性
type Asset struct {
	Balance          decimal.Decimal `json:"balance"`           // 总余额
	AvailableBalance decimal.Decimal `json:"available_balance"` // 可用余额
	FrozenBalance    decimal.Decimal `json:"frozen_balance"`    // 冻结余额
	MarginBalance    decimal.Decimal `json:"margin_balance"`    // 保证金余额
	ReserveBalance   decimal.Decimal `json:"reserve_balance"`   // 预留余额
	Currency         string          `json:"currency"`          // 币种
	UpdatedAt        int64           `json:"updated_at"`        // 更新时间
}

// NewAccount 创建新账户
func NewAccount() *Account {
	return &Account{
		AccountValues: &AccountValues{},
	}
}

// SetValues 设置AccountValues
func (a *AccountValues) SetValues(values *AccountValues) {
	if values == nil {
		return
	}
	if values.UserID != nil {
		a.UserID = values.UserID
	}
	if values.UserType != nil {
		a.UserType = values.UserType
	}
	if values.Currency != nil {
		a.Currency = values.Currency
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

// Getter方法
func (a *AccountValues) GetUserID() string {
	if a.UserID == nil {
		return ""
	}
	return *a.UserID
}

func (a *AccountValues) GetUserType() string {
	if a.UserType == nil {
		return ""
	}
	return *a.UserType
}

func (a *AccountValues) GetCurrency() string {
	if a.Currency == nil {
		return ""
	}
	return *a.Currency
}

func (a *AccountValues) GetStatus() int {
	if a.Status == nil {
		return 0
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

// Setter方法(支持链式调用)
func (a *AccountValues) SetUserID(userID string) *AccountValues {
	a.UserID = &userID
	return a
}

func (a *AccountValues) SetUserType(userType string) *AccountValues {
	a.UserType = &userType
	return a
}

func (a *AccountValues) SetCurrency(currency string) *AccountValues {
	a.Currency = &currency
	return a
}

func (a *AccountValues) SetAsset(asset *Asset) *AccountValues {
	a.Asset = asset
	return a
}

func (a *AccountValues) SetStatus(status int) *AccountValues {
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

func GetAccountsByUserID(userID, userType string) ([]*Account, error) {
	var accounts []*Account
	err := WriteDB.Where("user_id = ? AND user_type = ?", userID, userType).Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func GetAccountsByIDs(ids []uint64) ([]*Account, error) {
	var accounts []*Account
	err := WriteDB.Where("id IN ?", ids).Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
