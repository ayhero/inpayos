package models

import (
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
)

// MerchantPayout 代付记录表
type MerchantPayout struct {
	ID     uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TrxID  string `json:"transaction_id" gorm:"column:transaction_id;type:varchar(64);uniqueIndex"`
	Mid    string `json:"mid" gorm:"column:mid;type:varchar(32);index"`
	UserID string `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	ReqID  string `json:"req_id" gorm:"column:req_id;type:varchar(64);index"`
	*MerchantPayoutValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type MerchantPayoutValues struct {
	Status        *string          `json:"status" gorm:"column:status;type:varchar(16);index;default:'pending'"`
	Amount        *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Fee           *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(36,18);default:0"`
	Ccy           *string          `json:"ccy" gorm:"column:ccy;type:varchar(16)"`
	Country       *string          `json:"country" gorm:"column:country;type:varchar(8)"`
	ChannelCode   *string          `json:"channel_code" gorm:"column:channel_code;type:varchar(32)"`
	PaymentMethod *string          `json:"payment_method" gorm:"column:payment_method;type:varchar(32)"`
	RecipientInfo *string          `json:"recipient_info" gorm:"column:recipient_info;type:json"` // 收款方信息
	NotifyURL     *string          `json:"notify_url" gorm:"column:notify_url;type:varchar(512)"`
	OriTrxID      *string          `json:"ori_trx_id" gorm:"column:ori_trx_id;type:varchar(64)"` // 原交易ID(退款使用)
	Metadata      *string          `json:"metadata" gorm:"column:metadata;type:json"`
	Remark        *string          `json:"remark" gorm:"column:remark;type:varchar(512)"`
	UsdAmount     *decimal.Decimal `json:"usd_amount" gorm:"column:usd_amount;type:decimal(36,18)"`    // USD金额
	SettleID      *string          `json:"settle_id" gorm:"column:settle_id;type:varchar(64)"`         // 结算ID
	SettleStatus  *string          `json:"settle_status" gorm:"column:settle_status;type:varchar(16)"` // 结算状态
	SettledAt     *int64           `json:"settled_at" gorm:"column:settled_at"`                        // 结算时间
	ExpiredAt     *int64           `json:"expired_at" gorm:"column:expired_at"`
	CanceledAt    *int64           `json:"canceled_at" gorm:"column:canceled_at"`
}

func (MerchantPayout) TableName() string {
	return "t_merchant_payouts"
}

// GetStatus returns the Status value
func (pov *MerchantPayoutValues) GetStatus() string {
	if pov.Status == nil {
		return ""
	}
	return *pov.Status
}

// GetAmount returns the Amount value
func (pov *MerchantPayoutValues) GetAmount() decimal.Decimal {
	if pov.Amount == nil {
		return decimal.Zero
	}
	return *pov.Amount
}

// GetFee returns the Fee value
func (pov *MerchantPayoutValues) GetFee() decimal.Decimal {
	if pov.Fee == nil {
		return decimal.Zero
	}
	return *pov.Fee
}

// GetCcy returns the Ccy value
func (pov *MerchantPayoutValues) GetCcy() string {
	if pov.Ccy == nil {
		return ""
	}
	return *pov.Ccy
}

// GetChannelCode returns the ChannelCode value
func (pov *MerchantPayoutValues) GetChannelCode() string {
	if pov.ChannelCode == nil {
		return ""
	}
	return *pov.ChannelCode
}

// GetPaymentMethod returns the PaymentMethod value
func (pov *MerchantPayoutValues) GetPaymentMethod() string {
	if pov.PaymentMethod == nil {
		return ""
	}
	return *pov.PaymentMethod
}

// GetRecipientInfo returns the RecipientInfo value
func (pov *MerchantPayoutValues) GetRecipientInfo() string {
	if pov.RecipientInfo == nil {
		return ""
	}
	return *pov.RecipientInfo
}

// GetNotifyURL returns the NotifyURL value
func (pov *MerchantPayoutValues) GetNotifyURL() string {
	if pov.NotifyURL == nil {
		return ""
	}
	return *pov.NotifyURL
}

// GetCountry returns the Country value
func (pov *MerchantPayoutValues) GetCountry() string {
	if pov.Country == nil {
		return ""
	}
	return *pov.Country
}

// GetExpiredAt returns the ExpiredAt value
func (pov *MerchantPayoutValues) GetExpiredAt() int64 {
	if pov.ExpiredAt == nil {
		return 0
	}
	return *pov.ExpiredAt
}

// GetCanceledAt returns the CanceledAt value
func (pov *MerchantPayoutValues) GetCanceledAt() int64 {
	if pov.CanceledAt == nil {
		return 0
	}
	return *pov.CanceledAt
}

// SetStatus sets the Status value
func (pov *MerchantPayoutValues) SetStatus(value string) *MerchantPayoutValues {
	pov.Status = &value
	return pov
}

// SetAmount sets the Amount value
func (pov *MerchantPayoutValues) SetAmount(value decimal.Decimal) *MerchantPayoutValues {
	pov.Amount = &value
	return pov
}

// SetFee sets the Fee value
func (pov *MerchantPayoutValues) SetFee(value decimal.Decimal) *MerchantPayoutValues {
	pov.Fee = &value
	return pov
}

// SetCcy sets the Ccy value
func (pov *MerchantPayoutValues) SetCcy(value string) *MerchantPayoutValues {
	pov.Ccy = &value
	return pov
}

// SetChannelCode sets the ChannelCode value
func (pov *MerchantPayoutValues) SetChannelCode(value string) *MerchantPayoutValues {
	pov.ChannelCode = &value
	return pov
}

// SetPaymentMethod sets the PaymentMethod value
func (pov *MerchantPayoutValues) SetPaymentMethod(value string) *MerchantPayoutValues {
	pov.PaymentMethod = &value
	return pov
}

// SetRecipientInfo sets the RecipientInfo value
func (pov *MerchantPayoutValues) SetRecipientInfo(value string) *MerchantPayoutValues {
	pov.RecipientInfo = &value
	return pov
}

// SetNotifyURL sets the NotifyURL value
func (pov *MerchantPayoutValues) SetNotifyURL(value string) *MerchantPayoutValues {
	pov.NotifyURL = &value
	return pov
}

// SetCountry sets the Country value
func (pov *MerchantPayoutValues) SetCountry(value string) *MerchantPayoutValues {
	pov.Country = &value
	return pov
}

// SetExpiredAt sets the ExpiredAt value
func (pov *MerchantPayoutValues) SetExpiredAt(value int64) *MerchantPayoutValues {
	pov.ExpiredAt = &value
	return pov
}

// SetCanceledAt sets the CanceledAt value
func (pov *MerchantPayoutValues) SetCanceledAt(value int64) *MerchantPayoutValues {
	pov.CanceledAt = &value
	return pov
}

// SetValues sets multiple PayoutValues fields at once
func (p *MerchantPayout) SetValues(values *MerchantPayoutValues) *MerchantPayout {
	if values == nil {
		return p
	}

	if p.MerchantPayoutValues == nil {
		p.MerchantPayoutValues = &MerchantPayoutValues{}
	}

	if values.Status != nil {
		p.MerchantPayoutValues.SetStatus(*values.Status)
	}
	if values.Amount != nil {
		p.MerchantPayoutValues.SetAmount(*values.Amount)
	}
	if values.Fee != nil {
		p.MerchantPayoutValues.SetFee(*values.Fee)
	}
	if values.Ccy != nil {
		p.MerchantPayoutValues.SetCcy(*values.Ccy)
	}
	if values.ChannelCode != nil {
		p.MerchantPayoutValues.SetChannelCode(*values.ChannelCode)
	}
	if values.PaymentMethod != nil {
		p.MerchantPayoutValues.SetPaymentMethod(*values.PaymentMethod)
	}
	if values.RecipientInfo != nil {
		p.MerchantPayoutValues.SetRecipientInfo(*values.RecipientInfo)
	}
	if values.NotifyURL != nil {
		p.MerchantPayoutValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.Country != nil {
		p.MerchantPayoutValues.SetCountry(*values.Country)
	}
	if values.ExpiredAt != nil {
		p.MerchantPayoutValues.SetExpiredAt(*values.ExpiredAt)
	}
	if values.CanceledAt != nil {
		p.MerchantPayoutValues.SetCanceledAt(*values.CanceledAt)
	}

	return p
}

// ToTransaction converts Payout to Transaction
func (p *MerchantPayout) ToTransaction() *Transaction {
	if p == nil {
		return nil
	}

	transaction := &Transaction{
		ID:      p.ID,
		TrxID:   p.TrxID,
		Mid:     p.Mid,
		UserID:  p.UserID,
		ReqID:   p.ReqID,
		TrxType: protocol.TrxTypePayout, // Set transaction type to payout
		TransactionValues: &TransactionValues{
			Status:        p.MerchantPayoutValues.Status,
			Amount:        p.MerchantPayoutValues.Amount,
			Fee:           p.MerchantPayoutValues.Fee,
			Ccy:           p.MerchantPayoutValues.Ccy,
			ChannelCode:   p.MerchantPayoutValues.ChannelCode,
			PaymentMethod: p.MerchantPayoutValues.PaymentMethod,
			NotifyURL:     p.MerchantPayoutValues.NotifyURL,
			ReturnURL:     nil,                                  // Payout doesn't have ReturnURL
			NotifyStatus:  nil,                                  // Payout doesn't have NotifyStatus
			NotifyTimes:   nil,                                  // Payout doesn't have NotifyTimes
			OriTrxID:      nil,                                  // Payout doesn't have OriTrxID
			Metadata:      p.MerchantPayoutValues.RecipientInfo, // Map RecipientInfo to Metadata
			Remark:        nil,                                  // Payout doesn't have Remark
			ExpiredAt:     p.MerchantPayoutValues.ExpiredAt,
			ConfirmedAt:   nil, // Payout doesn't have ConfirmedAt
			CanceledAt:    p.MerchantPayoutValues.CanceledAt,
			UpdatedAt:     p.UpdatedAt, // Payout has UpdatedAt in main struct
		},
		CreatedAt: p.CreatedAt,
	}

	return transaction
}
