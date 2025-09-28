package models

import (
	"github.com/shopspring/decimal"
)

// Receipt 代收记录表
type Receipt struct {
	ID       uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	RecordID string `json:"record_id" gorm:"column:record_id;type:varchar(64);uniqueIndex"`
	Salt     string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*ReceiptValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type ReceiptValues struct {
	UserID        *string          `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	UserType      *string          `json:"user_type" gorm:"column:user_type;type:varchar(16);index"`
	BillID        *string          `json:"bill_id" gorm:"column:bill_id;type:varchar(64);index"`
	Status        *string          `json:"status" gorm:"column:status;type:varchar(16);index;default:'pending'"`
	Amount        *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Fee           *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(36,18);default:0"`
	Currency      *string          `json:"currency" gorm:"column:currency;type:varchar(16)"`
	ChannelCode   *string          `json:"channel_code" gorm:"column:channel_code;type:varchar(32)"`
	PaymentMethod *string          `json:"payment_method" gorm:"column:payment_method;type:varchar(32)"`
	NotifyURL     *string          `json:"notify_url" gorm:"column:notify_url;type:varchar(512)"`
	ReturnURL     *string          `json:"return_url" gorm:"column:return_url;type:varchar(512)"`
	Country       *string          `json:"country" gorm:"column:country;type:varchar(8)"`
	ExpiredAt     *int64           `json:"expired_at" gorm:"column:expired_at"`
	ConfirmedAt   *int64           `json:"confirmed_at" gorm:"column:confirmed_at"`

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Receipt) TableName() string {
	return "t_receipts"
}

// ToTransaction 转换为Transaction实体
func (r *Receipt) ToTransaction() *Transaction {
	transaction := NewTransaction()
	transaction.TransactionID = r.RecordID

	if r.ReceiptValues != nil {
		transaction.TransactionValues.SetUserID(r.GetUserID()).
			SetUserType(r.GetUserType()).
			SetBillID(r.GetBillID()).
			SetType("receipt").
			SetStatus(r.GetStatus()).
			SetAmount(r.GetAmount()).
			SetFee(r.GetFee()).
			SetCurrency(r.GetCurrency()).
			SetChannelCode(r.GetChannelCode()).
			SetPaymentMethod(r.GetPaymentMethod()).
			SetNotifyURL(r.GetNotifyURL()).
			SetExpiredAt(r.GetExpiredAt()).
			SetConfirmedAt(r.GetConfirmedAt())
	}

	transaction.CreatedAt = r.CreatedAt
	return transaction
}

// Getter方法
func (r *ReceiptValues) GetUserID() string {
	if r.UserID == nil {
		return ""
	}
	return *r.UserID
}

func (r *ReceiptValues) GetUserType() string {
	if r.UserType == nil {
		return ""
	}
	return *r.UserType
}

func (r *ReceiptValues) GetBillID() string {
	if r.BillID == nil {
		return ""
	}
	return *r.BillID
}

func (r *ReceiptValues) GetStatus() string {
	if r.Status == nil {
		return ""
	}
	return *r.Status
}

func (r *ReceiptValues) GetAmount() decimal.Decimal {
	if r.Amount == nil {
		return decimal.Zero
	}
	return *r.Amount
}

func (r *ReceiptValues) GetFee() decimal.Decimal {
	if r.Fee == nil {
		return decimal.Zero
	}
	return *r.Fee
}

func (r *ReceiptValues) GetCurrency() string {
	if r.Currency == nil {
		return ""
	}
	return *r.Currency
}

func (r *ReceiptValues) GetChannelCode() string {
	if r.ChannelCode == nil {
		return ""
	}
	return *r.ChannelCode
}

func (r *ReceiptValues) GetPaymentMethod() string {
	if r.PaymentMethod == nil {
		return ""
	}
	return *r.PaymentMethod
}

func (r *ReceiptValues) GetNotifyURL() string {
	if r.NotifyURL == nil {
		return ""
	}
	return *r.NotifyURL
}

func (r *ReceiptValues) GetExpiredAt() int64 {
	if r.ExpiredAt == nil {
		return 0
	}
	return *r.ExpiredAt
}

func (r *ReceiptValues) GetConfirmedAt() int64 {
	if r.ConfirmedAt == nil {
		return 0
	}
	return *r.ConfirmedAt
}

// Setter方法(支持链式调用)
func (r *ReceiptValues) SetUserID(userID string) *ReceiptValues {
	r.UserID = &userID
	return r
}

func (r *ReceiptValues) SetUserType(userType string) *ReceiptValues {
	r.UserType = &userType
	return r
}

func (r *ReceiptValues) SetBillID(billID string) *ReceiptValues {
	r.BillID = &billID
	return r
}

func (r *ReceiptValues) SetStatus(status string) *ReceiptValues {
	r.Status = &status
	return r
}

func (r *ReceiptValues) SetAmount(amount decimal.Decimal) *ReceiptValues {
	r.Amount = &amount
	return r
}

func (r *ReceiptValues) SetFee(fee decimal.Decimal) *ReceiptValues {
	r.Fee = &fee
	return r
}

func (r *ReceiptValues) SetCurrency(currency string) *ReceiptValues {
	r.Currency = &currency
	return r
}

func (r *ReceiptValues) SetChannelCode(channelCode string) *ReceiptValues {
	r.ChannelCode = &channelCode
	return r
}

func (r *ReceiptValues) SetPaymentMethod(paymentMethod string) *ReceiptValues {
	r.PaymentMethod = &paymentMethod
	return r
}

func (r *ReceiptValues) SetNotifyURL(notifyURL string) *ReceiptValues {
	r.NotifyURL = &notifyURL
	return r
}

func (r *ReceiptValues) SetExpiredAt(expiredAt int64) *ReceiptValues {
	r.ExpiredAt = &expiredAt
	return r
}

func (r *ReceiptValues) SetConfirmedAt(confirmedAt int64) *ReceiptValues {
	r.ConfirmedAt = &confirmedAt
	return r
}

// 查询方法
func GetReceiptByRecordID(recordID string) (*Receipt, error) {
	var receipt Receipt
	err := DB.Where("record_id = ?", recordID).First(&receipt).Error
	if err != nil {
		return nil, err
	}
	return &receipt, nil
}

func GetReceiptsByUserID(userID, userType string, limit, offset int) ([]*Receipt, error) {
	var receipts []*Receipt
	err := DB.Where("user_id = ? AND user_type = ?", userID, userType).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&receipts).Error
	if err != nil {
		return nil, err
	}
	return receipts, nil
}
