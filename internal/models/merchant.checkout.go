package models

import "inpayos/internal/protocol"

// MerchantCheckout 收银台会话表
type MerchantCheckout struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	CheckoutID string `gorm:"column:checkout_id;type:varchar(64);uniqueIndex;not null" json:"checkout_id"`
	Mid        string `gorm:"column:mid;type:varchar(64);not null;index" json:"mid"`
	ReqID      string `gorm:"column:req_id;type:varchar(64);index" json:"req_id"`
	TrxID      string `gorm:"column:trx_id;type:varchar(64);index" json:"trx_id"`
	TrxType    string `gorm:"column:trx_type;type:varchar(32);index" json:"trx_type"` // 交易类型: payin-代收, payout-代付
	*CheckoutValues
	CreatedAt int64 `gorm:"column:created_at;type:bigint;autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64 `gorm:"column:updated_at;type:bigint;autoUpdateTime:milli" json:"updated_at"`
}

type CheckoutValues struct {
	Ccy           *string          `gorm:"column:ccy;type:varchar(10)" json:"ccy"`
	Amount        *string          `gorm:"column:amount;type:decimal(20,8)" json:"amount"`
	Country       *string          `gorm:"column:country;type:varchar(3)" json:"country"`
	PaymentMethod *string          `gorm:"column:payment_method;type:varchar(32)" json:"payment_method"`
	ReturnURL     *string          `gorm:"column:return_url;type:varchar(1024)" json:"return_url"`
	CancelURL     *string          `gorm:"column:cancel_url;type:varchar(1024)" json:"cancel_url"`
	NotifyURL     *string          `gorm:"column:notify_url;type:varchar(1024)" json:"notify_url"`
	Status        *string          `gorm:"column:status;type:varchar(32);default:'created'" json:"status"` // created, pending, completed, cancelled, expired
	TrxID         *string          `gorm:"column:trx_id;type:varchar(64);index" json:"trx_id"`             // 关联的交易ID
	ChannelCode   *string          `gorm:"column:channel_code;type:varchar(32)" json:"channel_code"`
	Metadata      protocol.MapData `gorm:"column:metadata;type:text" json:"metadata"` // JSON格式的元数据
	ErrorCode     *string          `gorm:"column:error_code;type:varchar(32)" json:"error_code"`
	ErrorMsg      *string          `gorm:"column:error_msg;type:varchar(512)" json:"error_msg"`
	ExpiredAt     *int64           `gorm:"column:expired_at;type:bigint" json:"expired_at"`
	CompletedAt   *int64           `gorm:"column:completed_at;type:bigint" json:"completed_at"`
}

// TableName 返回表名
func (MerchantCheckout) TableName() string {
	return "t_merchant_checkouts"
}

// CheckoutValues Getter Methods
// GetCcy returns the Ccy value
func (cv *CheckoutValues) GetCcy() string {
	if cv.Ccy == nil {
		return ""
	}
	return *cv.Ccy
}

// GetAmount returns the Amount value
func (cv *CheckoutValues) GetAmount() string {
	if cv.Amount == nil {
		return ""
	}
	return *cv.Amount
}

// GetCountry returns the Country value
func (cv *CheckoutValues) GetCountry() string {
	if cv.Country == nil {
		return ""
	}
	return *cv.Country
}

// GetPaymentMethod returns the PaymentMethod value
func (cv *CheckoutValues) GetPaymentMethod() string {
	if cv.PaymentMethod == nil {
		return ""
	}
	return *cv.PaymentMethod
}

// GetReturnURL returns the ReturnURL value
func (cv *CheckoutValues) GetReturnURL() string {
	if cv.ReturnURL == nil {
		return ""
	}
	return *cv.ReturnURL
}

// GetCancelURL returns the CancelURL value
func (cv *CheckoutValues) GetCancelURL() string {
	if cv.CancelURL == nil {
		return ""
	}
	return *cv.CancelURL
}

// GetNotifyURL returns the NotifyURL value
func (cv *CheckoutValues) GetNotifyURL() string {
	if cv.NotifyURL == nil {
		return ""
	}
	return *cv.NotifyURL
}

// GetStatus returns the Status value
func (cv *CheckoutValues) GetStatus() string {
	if cv.Status == nil {
		return ""
	}
	return *cv.Status
}

// GetTrxID returns the TransactionID value
func (cv *CheckoutValues) GetTrxID() string {
	if cv.TrxID == nil {
		return ""
	}
	return *cv.TrxID
}

// GetChannelCode returns the ChannelCode value
func (cv *CheckoutValues) GetChannelCode() string {
	if cv.ChannelCode == nil {
		return ""
	}
	return *cv.ChannelCode
}

// GetMetadata returns the Metadata value
func (cv *CheckoutValues) GetMetadata() protocol.MapData {
	return cv.Metadata
}

// GetErrorCode returns the ErrorCode value
func (cv *CheckoutValues) GetErrorCode() string {
	if cv.ErrorCode == nil {
		return ""
	}
	return *cv.ErrorCode
}

// GetErrorMsg returns the ErrorMsg value
func (cv *CheckoutValues) GetErrorMsg() string {
	if cv.ErrorMsg == nil {
		return ""
	}
	return *cv.ErrorMsg
}

// GetExpiredAt returns the ExpiredAt value
func (cv *CheckoutValues) GetExpiredAt() int64 {
	if cv.ExpiredAt == nil {
		return 0
	}
	return *cv.ExpiredAt
}

// GetCompletedAt returns the CompletedAt value
func (cv *CheckoutValues) GetCompletedAt() int64 {
	if cv.CompletedAt == nil {
		return 0
	}
	return *cv.CompletedAt
}

// CheckoutValues Setter Methods (support method chaining)
// SetCcy sets the Ccy value
func (cv *CheckoutValues) SetCcy(value string) *CheckoutValues {
	cv.Ccy = &value
	return cv
}

// SetAmount sets the Amount value
func (cv *CheckoutValues) SetAmount(value string) *CheckoutValues {
	cv.Amount = &value
	return cv
}

// SetCountry sets the Country value
func (cv *CheckoutValues) SetCountry(value string) *CheckoutValues {
	cv.Country = &value
	return cv
}

// SetPaymentMethod sets the PaymentMethod value
func (cv *CheckoutValues) SetPaymentMethod(value string) *CheckoutValues {
	cv.PaymentMethod = &value
	return cv
}

// SetReturnURL sets the ReturnURL value
func (cv *CheckoutValues) SetReturnURL(value string) *CheckoutValues {
	cv.ReturnURL = &value
	return cv
}

// SetCancelURL sets the CancelURL value
func (cv *CheckoutValues) SetCancelURL(value string) *CheckoutValues {
	cv.CancelURL = &value
	return cv
}

// SetNotifyURL sets the NotifyURL value
func (cv *CheckoutValues) SetNotifyURL(value string) *CheckoutValues {
	cv.NotifyURL = &value
	return cv
}

// SetStatus sets the Status value
func (cv *CheckoutValues) SetStatus(value string) *CheckoutValues {
	cv.Status = &value
	return cv
}

// SetTrxID sets the TransactionID value
func (cv *CheckoutValues) SetTrxID(value string) *CheckoutValues {
	cv.TrxID = &value
	return cv
}

// SetChannelCode sets the ChannelCode value
func (cv *CheckoutValues) SetChannelCode(value string) *CheckoutValues {
	cv.ChannelCode = &value
	return cv
}

// SetMetadata sets the Metadata value
func (cv *CheckoutValues) SetMetadata(value protocol.MapData) *CheckoutValues {
	cv.Metadata = value
	return cv
}

// SetErrorCode sets the ErrorCode value
func (cv *CheckoutValues) SetErrorCode(value string) *CheckoutValues {
	cv.ErrorCode = &value
	return cv
}

// SetErrorMsg sets the ErrorMsg value
func (cv *CheckoutValues) SetErrorMsg(value string) *CheckoutValues {
	cv.ErrorMsg = &value
	return cv
}

// SetExpiredAt sets the ExpiredAt value
func (cv *CheckoutValues) SetExpiredAt(value int64) *CheckoutValues {
	cv.ExpiredAt = &value
	return cv
}

// SetCompletedAt sets the CompletedAt value
func (cv *CheckoutValues) SetCompletedAt(value int64) *CheckoutValues {
	cv.CompletedAt = &value
	return cv
}

// SetValues sets multiple CheckoutValues fields at once
func (c *MerchantCheckout) SetValues(values *CheckoutValues) *MerchantCheckout {
	if values == nil {
		return c
	}

	if c.CheckoutValues == nil {
		c.CheckoutValues = &CheckoutValues{}
	}

	// Set all fields from the provided values
	if values.Ccy != nil {
		c.CheckoutValues.SetCcy(*values.Ccy)
	}
	if values.Amount != nil {
		c.CheckoutValues.SetAmount(*values.Amount)
	}
	if values.Country != nil {
		c.CheckoutValues.SetCountry(*values.Country)
	}
	if values.PaymentMethod != nil {
		c.CheckoutValues.SetPaymentMethod(*values.PaymentMethod)
	}
	if values.ReturnURL != nil {
		c.CheckoutValues.SetReturnURL(*values.ReturnURL)
	}
	if values.CancelURL != nil {
		c.CheckoutValues.SetCancelURL(*values.CancelURL)
	}
	if values.NotifyURL != nil {
		c.CheckoutValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.Status != nil {
		c.CheckoutValues.SetStatus(*values.Status)
	}
	if values.TrxID != nil {
		c.CheckoutValues.SetTrxID(*values.TrxID)
	}
	if values.ChannelCode != nil {
		c.CheckoutValues.SetChannelCode(*values.ChannelCode)
	}
	// Metadata is not a pointer, so we always set it
	c.CheckoutValues.SetMetadata(values.Metadata)
	if values.ErrorCode != nil {
		c.CheckoutValues.SetErrorCode(*values.ErrorCode)
	}
	if values.ErrorMsg != nil {
		c.CheckoutValues.SetErrorMsg(*values.ErrorMsg)
	}
	if values.ExpiredAt != nil {
		c.CheckoutValues.SetExpiredAt(*values.ExpiredAt)
	}
	if values.CompletedAt != nil {
		c.CheckoutValues.SetCompletedAt(*values.CompletedAt)
	}

	return c
}
