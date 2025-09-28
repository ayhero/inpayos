package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Transaction 交易记录表
type Transaction struct {
	ID            uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TransactionID string `json:"transaction_id" gorm:"column:transaction_id;type:varchar(64);uniqueIndex"`
	Salt          string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*TransactionValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type TransactionValues struct {
	UserID        *string          `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	UserType      *string          `json:"user_type" gorm:"column:user_type;type:varchar(16);index"`
	BillID        *string          `json:"bill_id" gorm:"column:bill_id;type:varchar(64);index"`
	Type          *string          `json:"type" gorm:"column:type;type:varchar(16);index"`     // receipt, payment, refund, transfer
	Status        *string          `json:"status" gorm:"column:status;type:varchar(16);index"` // pending, processing, success, failed
	Amount        *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Fee           *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(36,18);default:0"`
	Currency      *string          `json:"currency" gorm:"column:currency;type:varchar(16)"`
	ChannelCode   *string          `json:"channel_code" gorm:"column:channel_code;type:varchar(32)"`
	PaymentMethod *string          `json:"payment_method" gorm:"column:payment_method;type:varchar(32)"`
	NotifyURL     *string          `json:"notify_url" gorm:"column:notify_url;type:varchar(512)"`
	ReturnURL     *string          `json:"return_url" gorm:"column:return_url;type:varchar(512)"`
	NotifyStatus  *string          `json:"notify_status" gorm:"column:notify_status;type:varchar(16);default:'pending'"`
	NotifyTimes   *int             `json:"notify_times" gorm:"column:notify_times;type:int;default:0"`
	SourceTxID    *string          `json:"source_tx_id" gorm:"column:source_tx_id;type:varchar(64)"` // 原交易ID(退款使用)
	Metadata      *string          `json:"metadata" gorm:"column:metadata;type:json"`
	Remark        *string          `json:"remark" gorm:"column:remark;type:varchar(512)"`
	ExpiredAt     *int64           `json:"expired_at" gorm:"column:expired_at"`
	ConfirmedAt   *int64           `json:"confirmed_at" gorm:"column:confirmed_at"`

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Transaction) TableName() string {
	return "t_transactions"
}

// NewTransaction 创建新交易
func NewTransaction() *Transaction {
	return &Transaction{
		TransactionValues: &TransactionValues{},
	}
}

// SetValues 设置TransactionValues
func (t *TransactionValues) SetValues(values *TransactionValues) {
	if values == nil {
		return
	}
	if values.UserID != nil {
		t.UserID = values.UserID
	}
	if values.UserType != nil {
		t.UserType = values.UserType
	}
	if values.BillID != nil {
		t.BillID = values.BillID
	}
	if values.Type != nil {
		t.Type = values.Type
	}
	if values.Status != nil {
		t.Status = values.Status
	}
	if values.Amount != nil {
		t.Amount = values.Amount
	}
	if values.Fee != nil {
		t.Fee = values.Fee
	}
	if values.Currency != nil {
		t.Currency = values.Currency
	}
	if values.ChannelCode != nil {
		t.ChannelCode = values.ChannelCode
	}
	if values.PaymentMethod != nil {
		t.PaymentMethod = values.PaymentMethod
	}
	if values.NotifyURL != nil {
		t.NotifyURL = values.NotifyURL
	}
	if values.ReturnURL != nil {
		t.ReturnURL = values.ReturnURL
	}
	if values.NotifyStatus != nil {
		t.NotifyStatus = values.NotifyStatus
	}
	if values.NotifyTimes != nil {
		t.NotifyTimes = values.NotifyTimes
	}
	if values.SourceTxID != nil {
		t.SourceTxID = values.SourceTxID
	}
	if values.Metadata != nil {
		t.Metadata = values.Metadata
	}
	if values.Remark != nil {
		t.Remark = values.Remark
	}
	if values.ExpiredAt != nil {
		t.ExpiredAt = values.ExpiredAt
	}
	if values.ConfirmedAt != nil {
		t.ConfirmedAt = values.ConfirmedAt
	}
}

// Getter方法
func (t *TransactionValues) GetUserID() string {
	if t.UserID == nil {
		return ""
	}
	return *t.UserID
}

func (t *TransactionValues) GetType() string {
	if t.Type == nil {
		return ""
	}
	return *t.Type
}

func (t *TransactionValues) GetStatus() string {
	if t.Status == nil {
		return ""
	}
	return *t.Status
}

func (t *TransactionValues) GetAmount() decimal.Decimal {
	if t.Amount == nil {
		return decimal.Zero
	}
	return *t.Amount
}

func (t *TransactionValues) GetCurrency() string {
	if t.Currency == nil {
		return ""
	}
	return *t.Currency
}

func (t *TransactionValues) GetUserType() string {
	if t.UserType == nil {
		return ""
	}
	return *t.UserType
}

func (t *TransactionValues) GetBillID() string {
	if t.BillID == nil {
		return ""
	}
	return *t.BillID
}

func (t *TransactionValues) GetFee() decimal.Decimal {
	if t.Fee == nil {
		return decimal.Zero
	}
	return *t.Fee
}

func (t *TransactionValues) GetChannelCode() string {
	if t.ChannelCode == nil {
		return ""
	}
	return *t.ChannelCode
}

func (t *TransactionValues) GetPaymentMethod() string {
	if t.PaymentMethod == nil {
		return ""
	}
	return *t.PaymentMethod
}

func (t *TransactionValues) GetNotifyStatus() string {
	if t.NotifyStatus == nil {
		return ""
	}
	return *t.NotifyStatus
}

func (t *TransactionValues) GetNotifyTimes() int {
	if t.NotifyTimes == nil {
		return 0
	}
	return *t.NotifyTimes
}

func (t *TransactionValues) GetSourceTxID() string {
	if t.SourceTxID == nil {
		return ""
	}
	return *t.SourceTxID
}

func (t *TransactionValues) GetRemark() string {
	if t.Remark == nil {
		return ""
	}
	return *t.Remark
}

func (t *TransactionValues) GetExpiredAt() int64 {
	if t.ExpiredAt == nil {
		return 0
	}
	return *t.ExpiredAt
}

func (t *TransactionValues) GetConfirmedAt() int64 {
	if t.ConfirmedAt == nil {
		return 0
	}
	return *t.ConfirmedAt
}

// Setter方法(支持链式调用)
func (t *TransactionValues) SetUserID(userID string) *TransactionValues {
	t.UserID = &userID
	return t
}

func (t *TransactionValues) SetUserType(userType string) *TransactionValues {
	t.UserType = &userType
	return t
}

func (t *TransactionValues) SetBillID(billID string) *TransactionValues {
	t.BillID = &billID
	return t
}

func (t *TransactionValues) SetType(txType string) *TransactionValues {
	t.Type = &txType
	return t
}

func (t *TransactionValues) SetStatus(status string) *TransactionValues {
	t.Status = &status
	return t
}

func (t *TransactionValues) SetAmount(amount decimal.Decimal) *TransactionValues {
	t.Amount = &amount
	return t
}

func (t *TransactionValues) SetFee(fee decimal.Decimal) *TransactionValues {
	t.Fee = &fee
	return t
}

func (t *TransactionValues) SetCurrency(currency string) *TransactionValues {
	t.Currency = &currency
	return t
}

func (t *TransactionValues) SetChannelCode(channelCode string) *TransactionValues {
	t.ChannelCode = &channelCode
	return t
}

func (t *TransactionValues) SetPaymentMethod(paymentMethod string) *TransactionValues {
	t.PaymentMethod = &paymentMethod
	return t
}

func (t *TransactionValues) SetNotifyURL(notifyURL string) *TransactionValues {
	t.NotifyURL = &notifyURL
	return t
}

func (t *TransactionValues) SetSourceTxID(sourceTxID string) *TransactionValues {
	t.SourceTxID = &sourceTxID
	return t
}

func (t *TransactionValues) SetExpiredAt(expiredAt int64) *TransactionValues {
	t.ExpiredAt = &expiredAt
	return t
}

func (t *TransactionValues) SetConfirmedAt(confirmedAt int64) *TransactionValues {
	t.ConfirmedAt = &confirmedAt
	return t
}

// Transaction模型的方法
func (t *Transaction) GetUserType() string {
	return t.TransactionValues.GetUserType()
}

// 查询方法
func GetTransactionByID(id uint64) (*Transaction, error) {
	var transaction Transaction
	err := DB.Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func GetTransactionByTransactionID(transactionID string) (*Transaction, error) {
	var transaction Transaction
	err := DB.Where("transaction_id = ?", transactionID).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func GetTransactionsByUserID(userID, userType string, limit, offset int) ([]*Transaction, error) {
	var transactions []*Transaction
	err := DB.Where("user_id = ? AND user_type = ?", userID, userType).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func GetTransactionForUpdate(tx *gorm.DB, transactionID string) (*Transaction, error) {
	var transaction Transaction
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("transaction_id = ?", transactionID).
		First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}
