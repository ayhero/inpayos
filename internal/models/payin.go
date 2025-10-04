package models

import (
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
)

// Payin 代收记录表
type Payin struct {
	ID     uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TrxID  string `json:"transaction_id" gorm:"column:transaction_id;type:varchar(64);uniqueIndex"`
	Mid    string `json:"mid" gorm:"column:mid;type:varchar(32);index"`
	UserID string `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	ReqID  string `json:"req_id" gorm:"column:req_id;type:varchar(64);index"`
	Salt   string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*PayinValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type PayinValues struct {
	Status        *string          `json:"status" gorm:"column:status;type:varchar(16);index;default:'pending'"`
	Country       *string          `json:"country" gorm:"column:country;type:varchar(8)"`
	Ccy           *string          `json:"ccy" gorm:"column:ccy;type:varchar(16)"`
	Amount        *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Fee           *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(36,18);default:0"`
	ChannelCode   *string          `json:"channel_code" gorm:"column:channel_code;type:varchar(32)"`
	PaymentMethod *string          `json:"payment_method" gorm:"column:payment_method;type:varchar(32)"`
	NotifyURL     *string          `json:"notify_url" gorm:"column:notify_url;type:varchar(512)"`
	ReturnURL     *string          `json:"return_url" gorm:"column:return_url;type:varchar(512)"`
	OriTrxID      *string          `json:"ori_trx_id" gorm:"column:ori_trx_id;type:varchar(64)"` // 原交易ID(退款使用)
	Metadata      *string          `json:"metadata" gorm:"column:metadata;type:json"`
	Remark        *string          `json:"remark" gorm:"column:remark;type:varchar(512)"`
	UsdAmount     *decimal.Decimal `json:"usd_amount" gorm:"column:usd_amount;type:decimal(36,18)"`    // USD金额
	SettleID      *string          `json:"settle_id" gorm:"column:settle_id;type:varchar(64)"`         // 结算ID
	SettleStatus  *string          `json:"settle_status" gorm:"column:settle_status;type:varchar(16)"` // 结算状态
	SettledAt     *int64           `json:"settled_at" gorm:"column:settled_at"`                        // 结算时间
	ExpiredAt     *int64           `json:"expired_at" gorm:"column:expired_at"`
	ConfirmedAt   *int64           `json:"confirmed_at" gorm:"column:confirmed_at"`
	CanceledAt    *int64           `json:"canceled_at" gorm:"column:canceled_at"`
	UpdatedAt     int64            `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Payin) TableName() string {
	return "t_payins"
}

// GetStatus returns the Status value
func (pv *PayinValues) GetStatus() string {
	if pv.Status == nil {
		return ""
	}
	return *pv.Status
}

// GetCountry returns the Country value
func (pv *PayinValues) GetCountry() string {
	if pv.Country == nil {
		return ""
	}
	return *pv.Country
}

// GetCcy returns the Ccy value
func (pv *PayinValues) GetCcy() string {
	if pv.Ccy == nil {
		return ""
	}
	return *pv.Ccy
}

// GetAmount returns the Amount value
func (pv *PayinValues) GetAmount() decimal.Decimal {
	if pv.Amount == nil {
		return decimal.Zero
	}
	return *pv.Amount
}

// GetFee returns the Fee value
func (pv *PayinValues) GetFee() decimal.Decimal {
	if pv.Fee == nil {
		return decimal.Zero
	}
	return *pv.Fee
}

// GetChannelCode returns the ChannelCode value
func (pv *PayinValues) GetChannelCode() string {
	if pv.ChannelCode == nil {
		return ""
	}
	return *pv.ChannelCode
}

// GetPaymentMethod returns the PaymentMethod value
func (pv *PayinValues) GetPaymentMethod() string {
	if pv.PaymentMethod == nil {
		return ""
	}
	return *pv.PaymentMethod
}

// GetNotifyURL returns the NotifyURL value
func (pv *PayinValues) GetNotifyURL() string {
	if pv.NotifyURL == nil {
		return ""
	}
	return *pv.NotifyURL
}

// GetReturnURL returns the ReturnURL value
func (pv *PayinValues) GetReturnURL() string {
	if pv.ReturnURL == nil {
		return ""
	}
	return *pv.ReturnURL
}

// GetExpiredAt returns the ExpiredAt value
func (pv *PayinValues) GetExpiredAt() int64 {
	if pv.ExpiredAt == nil {
		return 0
	}
	return *pv.ExpiredAt
}

// GetConfirmedAt returns the ConfirmedAt value
func (pv *PayinValues) GetConfirmedAt() int64 {
	if pv.ConfirmedAt == nil {
		return 0
	}
	return *pv.ConfirmedAt
}

// GetCanceledAt returns the CanceledAt value
func (pv *PayinValues) GetCanceledAt() int64 {
	if pv.CanceledAt == nil {
		return 0
	}
	return *pv.CanceledAt
}

// GetUpdatedAt returns the UpdatedAt value
func (pv *PayinValues) GetUpdatedAt() int64 {
	return pv.UpdatedAt
}

// SetStatus sets the Status value
func (pv *PayinValues) SetStatus(value string) *PayinValues {
	pv.Status = &value
	return pv
}

// SetCountry sets the Country value
func (pv *PayinValues) SetCountry(value string) *PayinValues {
	pv.Country = &value
	return pv
}

// SetCcy sets the Ccy value
func (pv *PayinValues) SetCcy(value string) *PayinValues {
	pv.Ccy = &value
	return pv
}

// SetAmount sets the Amount value
func (pv *PayinValues) SetAmount(value decimal.Decimal) *PayinValues {
	pv.Amount = &value
	return pv
}

// SetFee sets the Fee value
func (pv *PayinValues) SetFee(value decimal.Decimal) *PayinValues {
	pv.Fee = &value
	return pv
}

// SetChannelCode sets the ChannelCode value
func (pv *PayinValues) SetChannelCode(value string) *PayinValues {
	pv.ChannelCode = &value
	return pv
}

// SetPaymentMethod sets the PaymentMethod value
func (pv *PayinValues) SetPaymentMethod(value string) *PayinValues {
	pv.PaymentMethod = &value
	return pv
}

// SetNotifyURL sets the NotifyURL value
func (pv *PayinValues) SetNotifyURL(value string) *PayinValues {
	pv.NotifyURL = &value
	return pv
}

// SetReturnURL sets the ReturnURL value
func (pv *PayinValues) SetReturnURL(value string) *PayinValues {
	pv.ReturnURL = &value
	return pv
}

// SetExpiredAt sets the ExpiredAt value
func (pv *PayinValues) SetExpiredAt(value int64) *PayinValues {
	pv.ExpiredAt = &value
	return pv
}

// SetConfirmedAt sets the ConfirmedAt value
func (pv *PayinValues) SetConfirmedAt(value int64) *PayinValues {
	pv.ConfirmedAt = &value
	return pv
}

// SetCanceledAt sets the CanceledAt value
func (pv *PayinValues) SetCanceledAt(value int64) *PayinValues {
	pv.CanceledAt = &value
	return pv
}

// SetUpdatedAt sets the UpdatedAt value
func (pv *PayinValues) SetUpdatedAt(value int64) *PayinValues {
	pv.UpdatedAt = value
	return pv
}

// SetValues sets multiple PayinValues fields at once
func (p *Payin) SetValues(values *PayinValues) *Payin {
	if values == nil {
		return p
	}

	if p.PayinValues == nil {
		p.PayinValues = &PayinValues{}
	}
	if values.Status != nil {
		p.PayinValues.SetStatus(*values.Status)
	}
	if values.Country != nil {
		p.PayinValues.SetCountry(*values.Country)
	}
	if values.Ccy != nil {
		p.PayinValues.SetCcy(*values.Ccy)
	}
	if values.Amount != nil {
		p.PayinValues.SetAmount(*values.Amount)
	}
	if values.Fee != nil {
		p.PayinValues.SetFee(*values.Fee)
	}
	if values.ChannelCode != nil {
		p.PayinValues.SetChannelCode(*values.ChannelCode)
	}
	if values.PaymentMethod != nil {
		p.PayinValues.SetPaymentMethod(*values.PaymentMethod)
	}
	if values.NotifyURL != nil {
		p.PayinValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.ReturnURL != nil {
		p.PayinValues.SetReturnURL(*values.ReturnURL)
	}
	if values.ExpiredAt != nil {
		p.PayinValues.SetExpiredAt(*values.ExpiredAt)
	}
	if values.ConfirmedAt != nil {
		p.PayinValues.SetConfirmedAt(*values.ConfirmedAt)
	}
	if values.CanceledAt != nil {
		p.PayinValues.SetCanceledAt(*values.CanceledAt)
	}
	// UpdatedAt is not a pointer, so we always set it
	p.PayinValues.SetUpdatedAt(values.UpdatedAt)

	return p
}

// ToTransaction converts Payin to Transaction
func (p *Payin) ToTransaction() *Transaction {
	if p == nil {
		return nil
	}

	transaction := &Transaction{
		ID:      p.ID,
		TrxID:   p.TrxID,
		Mid:     p.Mid,
		UserID:  p.UserID,
		ReqID:   p.ReqID,
		TrxType: protocol.TrxTypePayin, // Set transaction type to payin
		TransactionValues: &TransactionValues{
			Status:        p.PayinValues.Status,
			Amount:        p.PayinValues.Amount,
			Fee:           p.PayinValues.Fee,
			Ccy:           p.PayinValues.Ccy,
			ChannelCode:   p.PayinValues.ChannelCode,
			PaymentMethod: p.PayinValues.PaymentMethod,
			NotifyURL:     p.PayinValues.NotifyURL,
			ReturnURL:     p.PayinValues.ReturnURL,
			NotifyStatus:  nil, // Payin doesn't have NotifyStatus
			NotifyTimes:   nil, // Payin doesn't have NotifyTimes
			OriTrxID:      nil, // Payin doesn't have OriTrxID
			Metadata:      nil, // Payin doesn't have Metadata
			Remark:        nil, // Payin doesn't have Remark
			ExpiredAt:     p.PayinValues.ExpiredAt,
			ConfirmedAt:   p.PayinValues.ConfirmedAt,
			CanceledAt:    p.PayinValues.CanceledAt,
			UpdatedAt:     p.PayinValues.UpdatedAt,
		},
		CreatedAt: p.CreatedAt,
	}

	return transaction
}
