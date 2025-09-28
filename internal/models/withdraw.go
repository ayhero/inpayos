package models

import (
	"github.com/shopspring/decimal"
)

// Withdraw 提现记录表
type Withdraw struct {
	ID       uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	RecordID string `json:"record_id" gorm:"column:record_id;type:varchar(64);uniqueIndex"`
	Salt     string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*WithdrawValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type WithdrawValues struct {
	UserID          *string          `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	UserType        *string          `json:"user_type" gorm:"column:user_type;type:varchar(16);index"`
	BillID          *string          `json:"bill_id" gorm:"column:bill_id;type:varchar(64);index"`
	Status          *string          `json:"status" gorm:"column:status;type:varchar(16);index;default:'pending'"`
	Amount          *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Fee             *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(36,18);default:0"`
	Currency        *string          `json:"currency" gorm:"column:currency;type:varchar(16)"`
	ChannelCode     *string          `json:"channel_code" gorm:"column:channel_code;type:varchar(32)"`
	WithdrawAddress *string          `json:"withdraw_address" gorm:"column:withdraw_address;type:varchar(256)"` // 提现目标地址
	TxHash          *string          `json:"tx_hash" gorm:"column:tx_hash;type:varchar(128)"`                   // 区块链交易哈希
	NotifyURL       *string          `json:"notify_url" gorm:"column:notify_url;type:varchar(512)"`
	Country         *string          `json:"country" gorm:"column:country;type:varchar(8)"`
	ProcessedAt     *int64           `json:"processed_at" gorm:"column:processed_at"`

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Withdraw) TableName() string {
	return "t_withdraws"
}

// ToTransaction 转换为Transaction实体
func (w *Withdraw) ToTransaction() *Transaction {
	transaction := NewTransaction()
	transaction.TransactionID = w.RecordID

	if w.WithdrawValues != nil {
		transaction.TransactionValues.SetUserID(w.GetUserID()).
			SetUserType(w.GetUserType()).
			SetBillID(w.GetBillID()).
			SetType("withdraw").
			SetStatus(w.GetStatus()).
			SetAmount(w.GetAmount()).
			SetFee(w.GetFee()).
			SetCurrency(w.GetCurrency()).
			SetChannelCode(w.GetChannelCode()).
			SetNotifyURL(w.GetNotifyURL()).
			SetConfirmedAt(w.GetProcessedAt())
	}

	transaction.CreatedAt = w.CreatedAt
	return transaction
}

// Getter方法
func (w *WithdrawValues) GetUserID() string {
	if w.UserID == nil {
		return ""
	}
	return *w.UserID
}

func (w *WithdrawValues) GetUserType() string {
	if w.UserType == nil {
		return ""
	}
	return *w.UserType
}

func (w *WithdrawValues) GetBillID() string {
	if w.BillID == nil {
		return ""
	}
	return *w.BillID
}

func (w *WithdrawValues) GetStatus() string {
	if w.Status == nil {
		return ""
	}
	return *w.Status
}

func (w *WithdrawValues) GetAmount() decimal.Decimal {
	if w.Amount == nil {
		return decimal.Zero
	}
	return *w.Amount
}

func (w *WithdrawValues) GetFee() decimal.Decimal {
	if w.Fee == nil {
		return decimal.Zero
	}
	return *w.Fee
}

func (w *WithdrawValues) GetCurrency() string {
	if w.Currency == nil {
		return ""
	}
	return *w.Currency
}

func (w *WithdrawValues) GetChannelCode() string {
	if w.ChannelCode == nil {
		return ""
	}
	return *w.ChannelCode
}

func (w *WithdrawValues) GetNotifyURL() string {
	if w.NotifyURL == nil {
		return ""
	}
	return *w.NotifyURL
}

func (w *WithdrawValues) GetProcessedAt() int64 {
	if w.ProcessedAt == nil {
		return 0
	}
	return *w.ProcessedAt
}

// Setter方法(支持链式调用)
func (w *WithdrawValues) SetUserID(userID string) *WithdrawValues {
	w.UserID = &userID
	return w
}

func (w *WithdrawValues) SetUserType(userType string) *WithdrawValues {
	w.UserType = &userType
	return w
}

func (w *WithdrawValues) SetBillID(billID string) *WithdrawValues {
	w.BillID = &billID
	return w
}

func (w *WithdrawValues) SetStatus(status string) *WithdrawValues {
	w.Status = &status
	return w
}

func (w *WithdrawValues) SetAmount(amount decimal.Decimal) *WithdrawValues {
	w.Amount = &amount
	return w
}

func (w *WithdrawValues) SetFee(fee decimal.Decimal) *WithdrawValues {
	w.Fee = &fee
	return w
}

func (w *WithdrawValues) SetCurrency(currency string) *WithdrawValues {
	w.Currency = &currency
	return w
}

func (w *WithdrawValues) SetChannelCode(channelCode string) *WithdrawValues {
	w.ChannelCode = &channelCode
	return w
}

func (w *WithdrawValues) SetWithdrawAddress(address string) *WithdrawValues {
	w.WithdrawAddress = &address
	return w
}

func (w *WithdrawValues) SetTxHash(txHash string) *WithdrawValues {
	w.TxHash = &txHash
	return w
}

func (w *WithdrawValues) SetNotifyURL(notifyURL string) *WithdrawValues {
	w.NotifyURL = &notifyURL
	return w
}

func (w *WithdrawValues) SetProcessedAt(processedAt int64) *WithdrawValues {
	w.ProcessedAt = &processedAt
	return w
}

// 查询方法
func GetWithdrawByRecordID(recordID string) (*Withdraw, error) {
	var withdraw Withdraw
	err := DB.Where("record_id = ?", recordID).First(&withdraw).Error
	if err != nil {
		return nil, err
	}
	return &withdraw, nil
}

func GetWithdrawsByUserID(userID, userType string, limit, offset int) ([]*Withdraw, error) {
	var withdraws []*Withdraw
	err := DB.Where("user_id = ? AND user_type = ?", userID, userType).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&withdraws).Error
	if err != nil {
		return nil, err
	}
	return withdraws, nil
}
