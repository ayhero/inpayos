package models

import (
	"github.com/shopspring/decimal"
)

// Payment 代付记录表
type Payment struct {
	ID       uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	RecordID string `json:"record_id" gorm:"column:record_id;type:varchar(64);uniqueIndex"`
	Salt     string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*PaymentValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type PaymentValues struct {
	UserID        *string          `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	UserType      *string          `json:"user_type" gorm:"column:user_type;type:varchar(16);index"`
	BillID        *string          `json:"bill_id" gorm:"column:bill_id;type:varchar(64);index"`
	Status        *string          `json:"status" gorm:"column:status;type:varchar(16);index;default:'pending'"`
	Amount        *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Fee           *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(36,18);default:0"`
	Currency      *string          `json:"currency" gorm:"column:currency;type:varchar(16)"`
	ChannelCode   *string          `json:"channel_code" gorm:"column:channel_code;type:varchar(32)"`
	PaymentMethod *string          `json:"payment_method" gorm:"column:payment_method;type:varchar(32)"`
	RecipientInfo *string          `json:"recipient_info" gorm:"column:recipient_info;type:json"` // 收款方信息
	NotifyURL     *string          `json:"notify_url" gorm:"column:notify_url;type:varchar(512)"`
	Country       *string          `json:"country" gorm:"column:country;type:varchar(8)"`
	ExpiredAt     *int64           `json:"expired_at" gorm:"column:expired_at"`
	ProcessedAt   *int64           `json:"processed_at" gorm:"column:processed_at"`

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Payment) TableName() string {
	return "t_payments"
}

// ToTransaction 转换为Transaction实体
func (p *Payment) ToTransaction() *Transaction {
	transaction := NewTransaction()
	transaction.TransactionID = p.RecordID

	if p.PaymentValues != nil {
		transaction.TransactionValues.SetUserID(p.GetUserID()).
			SetUserType(p.GetUserType()).
			SetBillID(p.GetBillID()).
			SetType("payment").
			SetStatus(p.GetStatus()).
			SetAmount(p.GetAmount()).
			SetFee(p.GetFee()).
			SetCurrency(p.GetCurrency()).
			SetChannelCode(p.GetChannelCode()).
			SetPaymentMethod(p.GetPaymentMethod()).
			SetNotifyURL(p.GetNotifyURL()).
			SetExpiredAt(p.GetExpiredAt()).
			SetConfirmedAt(p.GetProcessedAt())
	}

	transaction.CreatedAt = p.CreatedAt
	return transaction
}

// Getter方法
func (p *PaymentValues) GetUserID() string {
	if p.UserID == nil {
		return ""
	}
	return *p.UserID
}

func (p *PaymentValues) GetUserType() string {
	if p.UserType == nil {
		return ""
	}
	return *p.UserType
}

func (p *PaymentValues) GetBillID() string {
	if p.BillID == nil {
		return ""
	}
	return *p.BillID
}

func (p *PaymentValues) GetStatus() string {
	if p.Status == nil {
		return ""
	}
	return *p.Status
}

func (p *PaymentValues) GetAmount() decimal.Decimal {
	if p.Amount == nil {
		return decimal.Zero
	}
	return *p.Amount
}

func (p *PaymentValues) GetFee() decimal.Decimal {
	if p.Fee == nil {
		return decimal.Zero
	}
	return *p.Fee
}

func (p *PaymentValues) GetCurrency() string {
	if p.Currency == nil {
		return ""
	}
	return *p.Currency
}

func (p *PaymentValues) GetChannelCode() string {
	if p.ChannelCode == nil {
		return ""
	}
	return *p.ChannelCode
}

func (p *PaymentValues) GetPaymentMethod() string {
	if p.PaymentMethod == nil {
		return ""
	}
	return *p.PaymentMethod
}

func (p *PaymentValues) GetNotifyURL() string {
	if p.NotifyURL == nil {
		return ""
	}
	return *p.NotifyURL
}

func (p *PaymentValues) GetExpiredAt() int64 {
	if p.ExpiredAt == nil {
		return 0
	}
	return *p.ExpiredAt
}

func (p *PaymentValues) GetProcessedAt() int64 {
	if p.ProcessedAt == nil {
		return 0
	}
	return *p.ProcessedAt
}

// Setter方法(支持链式调用)
func (p *PaymentValues) SetUserID(userID string) *PaymentValues {
	p.UserID = &userID
	return p
}

func (p *PaymentValues) SetUserType(userType string) *PaymentValues {
	p.UserType = &userType
	return p
}

func (p *PaymentValues) SetBillID(billID string) *PaymentValues {
	p.BillID = &billID
	return p
}

func (p *PaymentValues) SetStatus(status string) *PaymentValues {
	p.Status = &status
	return p
}

func (p *PaymentValues) SetAmount(amount decimal.Decimal) *PaymentValues {
	p.Amount = &amount
	return p
}

func (p *PaymentValues) SetFee(fee decimal.Decimal) *PaymentValues {
	p.Fee = &fee
	return p
}

func (p *PaymentValues) SetCurrency(currency string) *PaymentValues {
	p.Currency = &currency
	return p
}

func (p *PaymentValues) SetChannelCode(channelCode string) *PaymentValues {
	p.ChannelCode = &channelCode
	return p
}

func (p *PaymentValues) SetPaymentMethod(paymentMethod string) *PaymentValues {
	p.PaymentMethod = &paymentMethod
	return p
}

func (p *PaymentValues) SetNotifyURL(notifyURL string) *PaymentValues {
	p.NotifyURL = &notifyURL
	return p
}

func (p *PaymentValues) SetExpiredAt(expiredAt int64) *PaymentValues {
	p.ExpiredAt = &expiredAt
	return p
}

func (p *PaymentValues) SetProcessedAt(processedAt int64) *PaymentValues {
	p.ProcessedAt = &processedAt
	return p
}

// 查询方法
func GetPaymentByRecordID(recordID string) (*Payment, error) {
	var payment Payment
	err := DB.Where("record_id = ?", recordID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func GetPaymentsByUserID(userID, userType string, limit, offset int) ([]*Payment, error) {
	var payments []*Payment
	err := DB.Where("user_id = ? AND user_type = ?", userID, userType).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}
