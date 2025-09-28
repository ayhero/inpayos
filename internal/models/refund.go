package models

import (
	"github.com/shopspring/decimal"
)

// Refund 退款记录表
type Refund struct {
	ID       uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	RecordID string `json:"record_id" gorm:"column:record_id;type:varchar(64);uniqueIndex"`
	Salt     string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*RefundValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type RefundValues struct {
	UserID         *string          `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	UserType       *string          `json:"user_type" gorm:"column:user_type;type:varchar(16);index"`
	BillID         *string          `json:"bill_id" gorm:"column:bill_id;type:varchar(64);index"`
	OriginalBillID *string          `json:"original_bill_id" gorm:"column:original_bill_id;type:varchar(64);index"` // 原交易订单ID
	Status         *string          `json:"status" gorm:"column:status;type:varchar(16);index;default:'pending'"`
	Amount         *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Fee            *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(36,18);default:0"`
	Currency       *string          `json:"currency" gorm:"column:currency;type:varchar(16)"`
	RefundReason   *string          `json:"refund_reason" gorm:"column:refund_reason;type:varchar(256)"`
	NotifyURL      *string          `json:"notify_url" gorm:"column:notify_url;type:varchar(512)"`
	ProcessedAt    *int64           `json:"processed_at" gorm:"column:processed_at"`

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Refund) TableName() string {
	return "t_refunds"
}

// ToTransaction 转换为Transaction实体
func (r *Refund) ToTransaction() *Transaction {
	transaction := NewTransaction()
	transaction.TransactionID = r.RecordID

	if r.RefundValues != nil {
		transaction.TransactionValues.SetUserID(r.GetUserID()).
			SetUserType(r.GetUserType()).
			SetBillID(r.GetBillID()).
			SetType("refund").
			SetStatus(r.GetStatus()).
			SetAmount(r.GetAmount()).
			SetFee(r.GetFee()).
			SetCurrency(r.GetCurrency()).
			SetSourceTxID(r.GetOriginalBillID()).
			SetNotifyURL(r.GetNotifyURL()).
			SetConfirmedAt(r.GetProcessedAt())
	}

	transaction.CreatedAt = r.CreatedAt
	return transaction
}

// Getter方法
func (r *RefundValues) GetUserID() string {
	if r.UserID == nil {
		return ""
	}
	return *r.UserID
}

func (r *RefundValues) GetUserType() string {
	if r.UserType == nil {
		return ""
	}
	return *r.UserType
}

func (r *RefundValues) GetBillID() string {
	if r.BillID == nil {
		return ""
	}
	return *r.BillID
}

func (r *RefundValues) GetOriginalBillID() string {
	if r.OriginalBillID == nil {
		return ""
	}
	return *r.OriginalBillID
}

func (r *RefundValues) GetStatus() string {
	if r.Status == nil {
		return ""
	}
	return *r.Status
}

func (r *RefundValues) GetAmount() decimal.Decimal {
	if r.Amount == nil {
		return decimal.Zero
	}
	return *r.Amount
}

func (r *RefundValues) GetFee() decimal.Decimal {
	if r.Fee == nil {
		return decimal.Zero
	}
	return *r.Fee
}

func (r *RefundValues) GetCurrency() string {
	if r.Currency == nil {
		return ""
	}
	return *r.Currency
}

func (r *RefundValues) GetNotifyURL() string {
	if r.NotifyURL == nil {
		return ""
	}
	return *r.NotifyURL
}

func (r *RefundValues) GetProcessedAt() int64 {
	if r.ProcessedAt == nil {
		return 0
	}
	return *r.ProcessedAt
}

// Setter方法(支持链式调用)
func (r *RefundValues) SetUserID(userID string) *RefundValues {
	r.UserID = &userID
	return r
}

func (r *RefundValues) SetUserType(userType string) *RefundValues {
	r.UserType = &userType
	return r
}

func (r *RefundValues) SetBillID(billID string) *RefundValues {
	r.BillID = &billID
	return r
}

func (r *RefundValues) SetOriginalBillID(originalBillID string) *RefundValues {
	r.OriginalBillID = &originalBillID
	return r
}

func (r *RefundValues) SetStatus(status string) *RefundValues {
	r.Status = &status
	return r
}

func (r *RefundValues) SetAmount(amount decimal.Decimal) *RefundValues {
	r.Amount = &amount
	return r
}

func (r *RefundValues) SetFee(fee decimal.Decimal) *RefundValues {
	r.Fee = &fee
	return r
}

func (r *RefundValues) SetCurrency(currency string) *RefundValues {
	r.Currency = &currency
	return r
}

func (r *RefundValues) SetRefundReason(reason string) *RefundValues {
	r.RefundReason = &reason
	return r
}

func (r *RefundValues) SetNotifyURL(notifyURL string) *RefundValues {
	r.NotifyURL = &notifyURL
	return r
}

func (r *RefundValues) SetProcessedAt(processedAt int64) *RefundValues {
	r.ProcessedAt = &processedAt
	return r
}

// 查询方法
func GetRefundByRecordID(recordID string) (*Refund, error) {
	var refund Refund
	err := DB.Where("record_id = ?", recordID).First(&refund).Error
	if err != nil {
		return nil, err
	}
	return &refund, nil
}

func GetRefundsByOriginalBillID(originalBillID string) ([]*Refund, error) {
	var refunds []*Refund
	err := DB.Where("original_bill_id = ?", originalBillID).
		Order("created_at DESC").
		Find(&refunds).Error
	if err != nil {
		return nil, err
	}
	return refunds, nil
}
