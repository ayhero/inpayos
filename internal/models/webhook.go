package models

import (
	"github.com/shopspring/decimal"
)

// Webhook Webhook通知记录表
type Webhook struct {
	ID        uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	WebhookID string `json:"webhook_id" gorm:"column:webhook_id;type:varchar(64);uniqueIndex"`
	*WebhookValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type WebhookValues struct {
	UserID        *string          `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	UserType      *string          `json:"user_type" gorm:"column:user_type;type:varchar(16);index"`
	TransactionID *string          `json:"transaction_id" gorm:"column:transaction_id;type:varchar(64);index"`
	BillID        *string          `json:"bill_id" gorm:"column:bill_id;type:varchar(64);index"`
	Type          *string          `json:"type" gorm:"column:type;type:varchar(16);index"`
	Status        *string          `json:"status" gorm:"column:status;type:varchar(16);index;default:'pending'"`
	Amount        *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(20,8)"`
	Fee           *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(20,8)"`
	Currency      *string          `json:"currency" gorm:"column:currency;type:varchar(8);index"`
	NotifyURL     *string          `json:"notify_url" gorm:"column:notify_url;type:varchar(512)"`
	NotifyStatus  *string          `json:"notify_status" gorm:"column:notify_status;type:varchar(16);default:'pending'"`
	NotifyTimes   *int32           `json:"notify_times" gorm:"column:notify_times;default:0"`
	MaxRetryTimes *int32           `json:"max_retry_times" gorm:"column:max_retry_times;default:5"`
	NextNotifyAt  *int64           `json:"next_notify_at" gorm:"column:next_notify_at;index"`
	LastNotifyAt  *int64           `json:"last_notify_at" gorm:"column:last_notify_at"`
	ResponseCode  *string          `json:"response_code" gorm:"column:response_code;type:varchar(16)"`
	ResponseBody  *string          `json:"response_body" gorm:"column:response_body;type:text"`
	RequestBody   *string          `json:"request_body" gorm:"column:request_body;type:text"`
	Remark        *string          `json:"remark" gorm:"column:remark;type:varchar(512)"`
}

// 表名
func (Webhook) TableName() string {
	return "t_webhooks"
}

// Chainable setters
func (v *WebhookValues) SetUserID(userID string) *WebhookValues {
	v.UserID = &userID
	return v
}

func (v *WebhookValues) SetUserType(userType string) *WebhookValues {
	v.UserType = &userType
	return v
}

func (v *WebhookValues) SetTransactionID(transactionID string) *WebhookValues {
	v.TransactionID = &transactionID
	return v
}

func (v *WebhookValues) SetBillID(billID string) *WebhookValues {
	v.BillID = &billID
	return v
}

func (v *WebhookValues) SetType(txType string) *WebhookValues {
	v.Type = &txType
	return v
}

func (v *WebhookValues) SetStatus(status string) *WebhookValues {
	v.Status = &status
	return v
}

func (v *WebhookValues) SetAmount(amount decimal.Decimal) *WebhookValues {
	v.Amount = &amount
	return v
}

func (v *WebhookValues) SetFee(fee decimal.Decimal) *WebhookValues {
	v.Fee = &fee
	return v
}

func (v *WebhookValues) SetCurrency(currency string) *WebhookValues {
	v.Currency = &currency
	return v
}

func (v *WebhookValues) SetNotifyURL(notifyURL string) *WebhookValues {
	v.NotifyURL = &notifyURL
	return v
}

func (v *WebhookValues) SetNotifyStatus(status string) *WebhookValues {
	v.NotifyStatus = &status
	return v
}

func (v *WebhookValues) SetNotifyTimes(times int32) *WebhookValues {
	v.NotifyTimes = &times
	return v
}

func (v *WebhookValues) SetMaxRetryTimes(times int32) *WebhookValues {
	v.MaxRetryTimes = &times
	return v
}

func (v *WebhookValues) SetNextNotifyAt(timestamp int64) *WebhookValues {
	v.NextNotifyAt = &timestamp
	return v
}

func (v *WebhookValues) SetLastNotifyAt(timestamp int64) *WebhookValues {
	v.LastNotifyAt = &timestamp
	return v
}

func (v *WebhookValues) SetResponseCode(code string) *WebhookValues {
	v.ResponseCode = &code
	return v
}

func (v *WebhookValues) SetResponseBody(body string) *WebhookValues {
	v.ResponseBody = &body
	return v
}

func (v *WebhookValues) SetRequestBody(body string) *WebhookValues {
	v.RequestBody = &body
	return v
}

func (v *WebhookValues) SetRemark(remark string) *WebhookValues {
	v.Remark = &remark
	return v
}

// Chainable getters
func (v *WebhookValues) GetUserID() string {
	if v.UserID == nil {
		return ""
	}
	return *v.UserID
}

func (v *WebhookValues) GetUserType() string {
	if v.UserType == nil {
		return ""
	}
	return *v.UserType
}

func (v *WebhookValues) GetTransactionID() string {
	if v.TransactionID == nil {
		return ""
	}
	return *v.TransactionID
}

func (v *WebhookValues) GetBillID() string {
	if v.BillID == nil {
		return ""
	}
	return *v.BillID
}

func (v *WebhookValues) GetType() string {
	if v.Type == nil {
		return ""
	}
	return *v.Type
}

func (v *WebhookValues) GetStatus() string {
	if v.Status == nil {
		return "pending"
	}
	return *v.Status
}

func (v *WebhookValues) GetAmount() decimal.Decimal {
	if v.Amount == nil {
		return decimal.Zero
	}
	return *v.Amount
}

func (v *WebhookValues) GetFee() decimal.Decimal {
	if v.Fee == nil {
		return decimal.Zero
	}
	return *v.Fee
}

func (v *WebhookValues) GetCurrency() string {
	if v.Currency == nil {
		return ""
	}
	return *v.Currency
}

func (v *WebhookValues) GetNotifyURL() string {
	if v.NotifyURL == nil {
		return ""
	}
	return *v.NotifyURL
}

func (v *WebhookValues) GetNotifyStatus() string {
	if v.NotifyStatus == nil {
		return "pending"
	}
	return *v.NotifyStatus
}

func (v *WebhookValues) GetNotifyTimes() int32 {
	if v.NotifyTimes == nil {
		return 0
	}
	return *v.NotifyTimes
}

func (v *WebhookValues) GetMaxRetryTimes() int32 {
	if v.MaxRetryTimes == nil {
		return 5
	}
	return *v.MaxRetryTimes
}

func (v *WebhookValues) GetNextNotifyAt() int64 {
	if v.NextNotifyAt == nil {
		return 0
	}
	return *v.NextNotifyAt
}

func (v *WebhookValues) GetLastNotifyAt() int64 {
	if v.LastNotifyAt == nil {
		return 0
	}
	return *v.LastNotifyAt
}

func (v *WebhookValues) GetResponseCode() string {
	if v.ResponseCode == nil {
		return ""
	}
	return *v.ResponseCode
}

func (v *WebhookValues) GetResponseBody() string {
	if v.ResponseBody == nil {
		return ""
	}
	return *v.ResponseBody
}

func (v *WebhookValues) GetRequestBody() string {
	if v.RequestBody == nil {
		return ""
	}
	return *v.RequestBody
}

func (v *WebhookValues) GetRemark() string {
	if v.Remark == nil {
		return ""
	}
	return *v.Remark
}
