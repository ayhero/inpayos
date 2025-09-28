package models

import (
	"github.com/shopspring/decimal"
)

// Deposit 充值记录表
type Deposit struct {
	ID       uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	RecordID string `json:"record_id" gorm:"column:record_id;type:varchar(64);uniqueIndex"`
	Salt     string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*DepositValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type DepositValues struct {
	UserID         *string          `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	UserType       *string          `json:"user_type" gorm:"column:user_type;type:varchar(16);index"`
	BillID         *string          `json:"bill_id" gorm:"column:bill_id;type:varchar(64);index"`
	Status         *string          `json:"status" gorm:"column:status;type:varchar(16);index;default:'pending'"`
	Amount         *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Fee            *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(36,18);default:0"`
	Currency       *string          `json:"currency" gorm:"column:currency;type:varchar(16)"`
	ChannelCode    *string          `json:"channel_code" gorm:"column:channel_code;type:varchar(32)"`
	ChannelAddress *string          `json:"channel_address" gorm:"column:channel_address;type:varchar(256)"` // 充值地址
	TxHash         *string          `json:"tx_hash" gorm:"column:tx_hash;type:varchar(128)"`                 // 区块链交易哈希
	NotifyURL      *string          `json:"notify_url" gorm:"column:notify_url;type:varchar(512)"`
	Country        *string          `json:"country" gorm:"column:country;type:varchar(8)"`
	ProcessedAt    *int64           `json:"processed_at" gorm:"column:processed_at"`

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Deposit) TableName() string {
	return "t_deposits"
}

// ToTransaction 转换为Transaction实体
func (d *Deposit) ToTransaction() *Transaction {
	transaction := NewTransaction()
	transaction.TransactionID = d.RecordID

	if d.DepositValues != nil {
		transaction.TransactionValues.SetUserID(d.GetUserID()).
			SetUserType(d.GetUserType()).
			SetBillID(d.GetBillID()).
			SetType("deposit").
			SetStatus(d.GetStatus()).
			SetAmount(d.GetAmount()).
			SetFee(d.GetFee()).
			SetCurrency(d.GetCurrency()).
			SetChannelCode(d.GetChannelCode()).
			SetNotifyURL(d.GetNotifyURL()).
			SetConfirmedAt(d.GetProcessedAt())
	}

	transaction.CreatedAt = d.CreatedAt
	return transaction
}

// Getter方法
func (d *DepositValues) GetUserID() string {
	if d.UserID == nil {
		return ""
	}
	return *d.UserID
}

func (d *DepositValues) GetUserType() string {
	if d.UserType == nil {
		return ""
	}
	return *d.UserType
}

func (d *DepositValues) GetBillID() string {
	if d.BillID == nil {
		return ""
	}
	return *d.BillID
}

func (d *DepositValues) GetStatus() string {
	if d.Status == nil {
		return ""
	}
	return *d.Status
}

func (d *DepositValues) GetAmount() decimal.Decimal {
	if d.Amount == nil {
		return decimal.Zero
	}
	return *d.Amount
}

func (d *DepositValues) GetFee() decimal.Decimal {
	if d.Fee == nil {
		return decimal.Zero
	}
	return *d.Fee
}

func (d *DepositValues) GetCurrency() string {
	if d.Currency == nil {
		return ""
	}
	return *d.Currency
}

func (d *DepositValues) GetChannelCode() string {
	if d.ChannelCode == nil {
		return ""
	}
	return *d.ChannelCode
}

func (d *DepositValues) GetNotifyURL() string {
	if d.NotifyURL == nil {
		return ""
	}
	return *d.NotifyURL
}

func (d *DepositValues) GetProcessedAt() int64 {
	if d.ProcessedAt == nil {
		return 0
	}
	return *d.ProcessedAt
}

// Setter方法(支持链式调用)
func (d *DepositValues) SetUserID(userID string) *DepositValues {
	d.UserID = &userID
	return d
}

func (d *DepositValues) SetUserType(userType string) *DepositValues {
	d.UserType = &userType
	return d
}

func (d *DepositValues) SetBillID(billID string) *DepositValues {
	d.BillID = &billID
	return d
}

func (d *DepositValues) SetStatus(status string) *DepositValues {
	d.Status = &status
	return d
}

func (d *DepositValues) SetAmount(amount decimal.Decimal) *DepositValues {
	d.Amount = &amount
	return d
}

func (d *DepositValues) SetFee(fee decimal.Decimal) *DepositValues {
	d.Fee = &fee
	return d
}

func (d *DepositValues) SetCurrency(currency string) *DepositValues {
	d.Currency = &currency
	return d
}

func (d *DepositValues) SetChannelCode(channelCode string) *DepositValues {
	d.ChannelCode = &channelCode
	return d
}

func (d *DepositValues) SetChannelAddress(address string) *DepositValues {
	d.ChannelAddress = &address
	return d
}

func (d *DepositValues) SetTxHash(txHash string) *DepositValues {
	d.TxHash = &txHash
	return d
}

func (d *DepositValues) SetNotifyURL(notifyURL string) *DepositValues {
	d.NotifyURL = &notifyURL
	return d
}

func (d *DepositValues) SetProcessedAt(processedAt int64) *DepositValues {
	d.ProcessedAt = &processedAt
	return d
}

// 查询方法
func GetDepositByRecordID(recordID string) (*Deposit, error) {
	var deposit Deposit
	err := DB.Where("record_id = ?", recordID).First(&deposit).Error
	if err != nil {
		return nil, err
	}
	return &deposit, nil
}

func GetDepositsByUserID(userID, userType string, limit, offset int) ([]*Deposit, error) {
	var deposits []*Deposit
	err := DB.Where("user_id = ? AND user_type = ?", userID, userType).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&deposits).Error
	if err != nil {
		return nil, err
	}
	return deposits, nil
}
