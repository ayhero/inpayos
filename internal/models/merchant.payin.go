package models

import (
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
)

// MerchantPayin 代收记录表
type MerchantPayin struct {
	ID     uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TrxID  string `json:"transaction_id" gorm:"column:transaction_id;type:varchar(64);uniqueIndex"`
	Mid    string `json:"mid" gorm:"column:mid;type:varchar(32);index"`
	UserID string `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	ReqID  string `json:"req_id" gorm:"column:req_id;type:varchar(64);index"`
	*MerchantPayinValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type MerchantPayinValues struct {
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

func (MerchantPayin) TableName() string {
	return "t_merchant_payins"
}

// GetStatus returns the Status value
func (pv *MerchantPayinValues) GetStatus() string {
	if pv.Status == nil {
		return ""
	}
	return *pv.Status
}

// GetCountry returns the Country value
func (pv *MerchantPayinValues) GetCountry() string {
	if pv.Country == nil {
		return ""
	}
	return *pv.Country
}

// GetCcy returns the Ccy value
func (pv *MerchantPayinValues) GetCcy() string {
	if pv.Ccy == nil {
		return ""
	}
	return *pv.Ccy
}

// GetAmount returns the Amount value
func (pv *MerchantPayinValues) GetAmount() decimal.Decimal {
	if pv.Amount == nil {
		return decimal.Zero
	}
	return *pv.Amount
}

// GetFee returns the Fee value
func (pv *MerchantPayinValues) GetFee() decimal.Decimal {
	if pv.Fee == nil {
		return decimal.Zero
	}
	return *pv.Fee
}

// GetChannelCode returns the ChannelCode value
func (pv *MerchantPayinValues) GetChannelCode() string {
	if pv.ChannelCode == nil {
		return ""
	}
	return *pv.ChannelCode
}

// GetPaymentMethod returns the PaymentMethod value
func (pv *MerchantPayinValues) GetPaymentMethod() string {
	if pv.PaymentMethod == nil {
		return ""
	}
	return *pv.PaymentMethod
}

// GetNotifyURL returns the NotifyURL value
func (pv *MerchantPayinValues) GetNotifyURL() string {
	if pv.NotifyURL == nil {
		return ""
	}
	return *pv.NotifyURL
}

// GetReturnURL returns the ReturnURL value
func (pv *MerchantPayinValues) GetReturnURL() string {
	if pv.ReturnURL == nil {
		return ""
	}
	return *pv.ReturnURL
}

// GetExpiredAt returns the ExpiredAt value
func (pv *MerchantPayinValues) GetExpiredAt() int64 {
	if pv.ExpiredAt == nil {
		return 0
	}
	return *pv.ExpiredAt
}

// GetConfirmedAt returns the ConfirmedAt value
func (pv *MerchantPayinValues) GetConfirmedAt() int64 {
	if pv.ConfirmedAt == nil {
		return 0
	}
	return *pv.ConfirmedAt
}

// GetCanceledAt returns the CanceledAt value
func (pv *MerchantPayinValues) GetCanceledAt() int64 {
	if pv.CanceledAt == nil {
		return 0
	}
	return *pv.CanceledAt
}

// GetUpdatedAt returns the UpdatedAt value
func (pv *MerchantPayinValues) GetUpdatedAt() int64 {
	return pv.UpdatedAt
}

// SetStatus sets the Status value
func (pv *MerchantPayinValues) SetStatus(value string) *MerchantPayinValues {
	pv.Status = &value
	return pv
}

// SetCountry sets the Country value
func (pv *MerchantPayinValues) SetCountry(value string) *MerchantPayinValues {
	pv.Country = &value
	return pv
}

// SetCcy sets the Ccy value
func (pv *MerchantPayinValues) SetCcy(value string) *MerchantPayinValues {
	pv.Ccy = &value
	return pv
}

// SetAmount sets the Amount value
func (pv *MerchantPayinValues) SetAmount(value decimal.Decimal) *MerchantPayinValues {
	pv.Amount = &value
	return pv
}

// SetFee sets the Fee value
func (pv *MerchantPayinValues) SetFee(value decimal.Decimal) *MerchantPayinValues {
	pv.Fee = &value
	return pv
}

// SetChannelCode sets the ChannelCode value
func (pv *MerchantPayinValues) SetChannelCode(value string) *MerchantPayinValues {
	pv.ChannelCode = &value
	return pv
}

// SetPaymentMethod sets the PaymentMethod value
func (pv *MerchantPayinValues) SetPaymentMethod(value string) *MerchantPayinValues {
	pv.PaymentMethod = &value
	return pv
}

// SetNotifyURL sets the NotifyURL value
func (pv *MerchantPayinValues) SetNotifyURL(value string) *MerchantPayinValues {
	pv.NotifyURL = &value
	return pv
}

// SetReturnURL sets the ReturnURL value
func (pv *MerchantPayinValues) SetReturnURL(value string) *MerchantPayinValues {
	pv.ReturnURL = &value
	return pv
}

// SetExpiredAt sets the ExpiredAt value
func (pv *MerchantPayinValues) SetExpiredAt(value int64) *MerchantPayinValues {
	pv.ExpiredAt = &value
	return pv
}

// SetConfirmedAt sets the ConfirmedAt value
func (pv *MerchantPayinValues) SetConfirmedAt(value int64) *MerchantPayinValues {
	pv.ConfirmedAt = &value
	return pv
}

// SetCanceledAt sets the CanceledAt value
func (pv *MerchantPayinValues) SetCanceledAt(value int64) *MerchantPayinValues {
	pv.CanceledAt = &value
	return pv
}

// SetUpdatedAt sets the UpdatedAt value
func (pv *MerchantPayinValues) SetUpdatedAt(value int64) *MerchantPayinValues {
	pv.UpdatedAt = value
	return pv
}

// SetValues sets multiple PayinValues fields at once
func (p *MerchantPayin) SetValues(values *MerchantPayinValues) *MerchantPayin {
	if values == nil {
		return p
	}

	if p.MerchantPayinValues == nil {
		p.MerchantPayinValues = &MerchantPayinValues{}
	}
	if values.Status != nil {
		p.MerchantPayinValues.SetStatus(*values.Status)
	}
	if values.Country != nil {
		p.MerchantPayinValues.SetCountry(*values.Country)
	}
	if values.Ccy != nil {
		p.MerchantPayinValues.SetCcy(*values.Ccy)
	}
	if values.Amount != nil {
		p.MerchantPayinValues.SetAmount(*values.Amount)
	}
	if values.Fee != nil {
		p.MerchantPayinValues.SetFee(*values.Fee)
	}
	if values.ChannelCode != nil {
		p.MerchantPayinValues.SetChannelCode(*values.ChannelCode)
	}
	if values.PaymentMethod != nil {
		p.MerchantPayinValues.SetPaymentMethod(*values.PaymentMethod)
	}
	if values.NotifyURL != nil {
		p.MerchantPayinValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.ReturnURL != nil {
		p.MerchantPayinValues.SetReturnURL(*values.ReturnURL)
	}
	if values.ExpiredAt != nil {
		p.MerchantPayinValues.SetExpiredAt(*values.ExpiredAt)
	}
	if values.ConfirmedAt != nil {
		p.MerchantPayinValues.SetConfirmedAt(*values.ConfirmedAt)
	}
	if values.CanceledAt != nil {
		p.MerchantPayinValues.SetCanceledAt(*values.CanceledAt)
	}
	// UpdatedAt is not a pointer, so we always set it
	p.MerchantPayinValues.SetUpdatedAt(values.UpdatedAt)

	return p
}

// ToTransaction converts Payin to Transaction
func (p *MerchantPayin) ToTransaction() *Transaction {
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
			Status:        p.MerchantPayinValues.Status,
			Amount:        p.MerchantPayinValues.Amount,
			Fee:           p.MerchantPayinValues.Fee,
			Ccy:           p.MerchantPayinValues.Ccy,
			ChannelCode:   p.MerchantPayinValues.ChannelCode,
			PaymentMethod: p.MerchantPayinValues.PaymentMethod,
			NotifyURL:     p.MerchantPayinValues.NotifyURL,
			ReturnURL:     p.MerchantPayinValues.ReturnURL,
			NotifyStatus:  nil, // Payin doesn't have NotifyStatus
			NotifyTimes:   nil, // Payin doesn't have NotifyTimes
			OriTrxID:      nil, // Payin doesn't have OriTrxID
			Metadata:      nil, // Payin doesn't have Metadata
			Remark:        nil, // Payin doesn't have Remark
			ExpiredAt:     p.MerchantPayinValues.ExpiredAt,
			ConfirmedAt:   p.MerchantPayinValues.ConfirmedAt,
			CanceledAt:    p.MerchantPayinValues.CanceledAt,
			UpdatedAt:     p.MerchantPayinValues.UpdatedAt,
		},
		CreatedAt: p.CreatedAt,
	}

	return transaction
}
