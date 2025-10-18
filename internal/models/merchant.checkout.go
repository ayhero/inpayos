package models

import (
	"inpayos/internal/protocol"
	"slices"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// MerchantCheckout 收银台会话表
type MerchantCheckout struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	CheckoutID string `gorm:"column:checkout_id;type:varchar(64);uniqueIndex;not null" json:"checkout_id"`
	Mid        string `gorm:"column:mid;type:varchar(64);not null;index" json:"mid"`
	ReqID      string `gorm:"column:req_id;type:varchar(64);index" json:"req_id"`
	TrxType    string `gorm:"column:trx_type;type:varchar(32);index" json:"trx_type"` // 交易类型: payin-代收, payout-代付
	*MerchantCheckoutValues
	CreatedAt int64 `gorm:"column:created_at;type:bigint;autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64 `gorm:"column:updated_at;type:bigint;autoUpdateTime:milli" json:"updated_at"`
}

type MerchantCheckoutValues struct {
	Ccy          *string           `gorm:"column:ccy;type:varchar(10)" json:"ccy"`
	Amount       *decimal.Decimal  `gorm:"column:amount;type:decimal(20,8)" json:"amount"`
	Country      *string           `gorm:"column:country;type:varchar(3)" json:"country"`
	TrxID        *string           `gorm:"column:trx_id;type:varchar(64);index" json:"trx_id"`
	TrxApp       *string           `gorm:"column:trx_app;type:varchar(32)" json:"trx_app"`
	TrxMethod    *string           `gorm:"column:trx_method;type:varchar(32)" json:"trx_method"`
	ReturnURL    *string           `gorm:"column:return_url;type:varchar(1024)" json:"return_url"`
	NotifyURL    *string           `gorm:"column:notify_url;type:varchar(1024)" json:"notify_url"`
	Status       *string           `gorm:"column:status;type:varchar(32);default:'created'" json:"status"` // created, pending, completed, cancelled, expired
	ChannelCode  *string           `gorm:"column:channel_code;type:varchar(32)" json:"channel_code"`
	CheckoutURL  *string           `gorm:"column:checkout_url;type:varchar(1024)" json:"checkout_url"`
	Transactions []*Transaction    `gorm:"column:transactions;type:json;serializer:json" json:"transactions"` // 关联的交易记录，不存储在数据库中
	Metadata     *protocol.MapData `gorm:"column:meta_data;type:json;serializer:json" json:"meta_data"`       // JSON格式的元数据
	ErrorCode    *string           `gorm:"column:error_code;type:varchar(32)" json:"error_code"`
	ErrorMsg     *string           `gorm:"column:error_msg;type:varchar(512)" json:"error_msg"`
	ExpiredAt    *int64            `gorm:"column:expired_at;type:bigint" json:"expired_at"`
	SubmitedAt   *int64            `gorm:"column:submited_at;type:bigint" json:"submited_at"`
	ConfirmedAt  *int64            `gorm:"column:confirmed_at;type:bigint" json:"confirmed_at"`
	CanceledAt   *int64            `gorm:"column:canceled_at;type:bigint" json:"canceled_at"`
	CompletedAt  *int64            `gorm:"column:completed_at;type:bigint" json:"completed_at"`
}

// TableName 返回表名
func (MerchantCheckout) TableName() string {
	return "t_merchant_checkouts"
}

// CheckoutValues Getter Methods
// GetCcy returns the Ccy value
func (cv *MerchantCheckoutValues) GetCcy() string {
	if cv.Ccy == nil {
		return ""
	}
	return *cv.Ccy
}

// GetAmount returns the Amount value
func (cv *MerchantCheckoutValues) GetAmount() decimal.Decimal {
	if cv.Amount == nil {
		return decimal.Zero
	}
	return *cv.Amount
}

// GetCountry returns the Country value
func (cv *MerchantCheckoutValues) GetCountry() string {
	if cv.Country == nil {
		return ""
	}
	return *cv.Country
}

// GetTrxID returns the TrxID value
func (cv *MerchantCheckoutValues) GetTrxID() string {
	if cv.TrxID == nil {
		return ""
	}
	return *cv.TrxID
}

// GetTrxMethod returns the PaymentMethod value
func (cv *MerchantCheckoutValues) GetTrxMethod() string {
	if cv.TrxMethod == nil {
		return ""
	}
	return *cv.TrxMethod
}

// GetReturnURL returns the ReturnURL value
func (cv *MerchantCheckoutValues) GetReturnURL() string {
	if cv.ReturnURL == nil {
		return ""
	}
	return *cv.ReturnURL
}

// GetNotifyURL returns the NotifyURL value
func (cv *MerchantCheckoutValues) GetNotifyURL() string {
	if cv.NotifyURL == nil {
		return ""
	}
	return *cv.NotifyURL
}

// GetStatus returns the Status value
func (cv *MerchantCheckoutValues) GetStatus() string {
	if cv.Status == nil {
		return ""
	}
	return *cv.Status
}

// GetChannelCode returns the ChannelCode value
func (cv *MerchantCheckoutValues) GetChannelCode() string {
	if cv.ChannelCode == nil {
		return ""
	}
	return *cv.ChannelCode
}

// GetMetadata returns the Metadata value
func (cv *MerchantCheckoutValues) GetMetadata() protocol.MapData {
	if cv.Metadata == nil {
		return protocol.MapData{}
	}
	return *cv.Metadata
}

// GetErrorCode returns the ErrorCode value
func (cv *MerchantCheckoutValues) GetErrorCode() string {
	if cv.ErrorCode == nil {
		return ""
	}
	return *cv.ErrorCode
}

// GetErrorMsg returns the ErrorMsg value
func (cv *MerchantCheckoutValues) GetErrorMsg() string {
	if cv.ErrorMsg == nil {
		return ""
	}
	return *cv.ErrorMsg
}

// GetExpiredAt returns the ExpiredAt value
func (cv *MerchantCheckoutValues) GetExpiredAt() int64 {
	if cv.ExpiredAt == nil {
		return 0
	}
	return *cv.ExpiredAt
}

// GetCompletedAt returns the CompletedAt value
func (cv *MerchantCheckoutValues) GetCompletedAt() int64 {
	if cv.CompletedAt == nil {
		return 0
	}
	return *cv.CompletedAt
}

// GetTrxApp returns the TrxApp value
func (cv *MerchantCheckoutValues) GetTrxApp() string {
	if cv.TrxApp == nil {
		return ""
	}
	return *cv.TrxApp
}

// GetCheckoutURL returns the CheckoutURL value
func (cv *MerchantCheckoutValues) GetCheckoutURL() string {
	if cv.CheckoutURL == nil {
		return ""
	}
	return *cv.CheckoutURL
}

// GetSubmitedAt returns the SubmitedAt value
func (cv *MerchantCheckoutValues) GetSubmitedAt() int64 {
	if cv.SubmitedAt == nil {
		return 0
	}
	return *cv.SubmitedAt
}

// GetConfirmedAt returns the ConfirmedAt value
func (cv *MerchantCheckoutValues) GetConfirmedAt() int64 {
	if cv.ConfirmedAt == nil {
		return 0
	}
	return *cv.ConfirmedAt
}

// GetCanceledAt returns the CanceledAt value
func (cv *MerchantCheckoutValues) GetCanceledAt() int64 {
	if cv.CanceledAt == nil {
		return 0
	}
	return *cv.CanceledAt
}

// GetTransactions returns the Transactions value
func (cv *MerchantCheckoutValues) GetTransactions() []*Transaction {
	if cv.Transactions == nil {
		return []*Transaction{}
	}
	return cv.Transactions
}

// CheckoutValues Setter Methods (support method chaining)
// SetCcy sets the Ccy value
func (cv *MerchantCheckoutValues) SetCcy(value string) *MerchantCheckoutValues {
	cv.Ccy = &value
	return cv
}

// SetAmount sets the Amount value
func (cv *MerchantCheckoutValues) SetAmount(value decimal.Decimal) *MerchantCheckoutValues {
	cv.Amount = &value
	return cv
}

// SetCountry sets the Country value
func (cv *MerchantCheckoutValues) SetCountry(value string) *MerchantCheckoutValues {
	cv.Country = &value
	return cv
}

// SetTrxID sets the TrxID value
func (cv *MerchantCheckoutValues) SetTrxID(value string) *MerchantCheckoutValues {
	cv.TrxID = &value
	return cv
}

// SetTrxMethod sets the PaymentMethod value
func (cv *MerchantCheckoutValues) SetTrxMethod(value string) *MerchantCheckoutValues {
	cv.TrxMethod = &value
	return cv
}

// SetReturnURL sets the ReturnURL value
func (cv *MerchantCheckoutValues) SetReturnURL(value string) *MerchantCheckoutValues {
	cv.ReturnURL = &value
	return cv
}

// SetNotifyURL sets the NotifyURL value
func (cv *MerchantCheckoutValues) SetNotifyURL(value string) *MerchantCheckoutValues {
	cv.NotifyURL = &value
	return cv
}

// SetStatus sets the Status value
func (cv *MerchantCheckoutValues) SetStatus(value string) *MerchantCheckoutValues {
	cv.Status = &value
	return cv
}

// SetChannelCode sets the ChannelCode value
func (cv *MerchantCheckoutValues) SetChannelCode(value string) *MerchantCheckoutValues {
	cv.ChannelCode = &value
	return cv
}

// SetMetadata sets the Metadata value
func (cv *MerchantCheckoutValues) SetMetadata(value protocol.MapData) *MerchantCheckoutValues {
	cv.Metadata = &value
	return cv
}

// SetErrorCode sets the ErrorCode value
func (cv *MerchantCheckoutValues) SetErrorCode(value string) *MerchantCheckoutValues {
	cv.ErrorCode = &value
	return cv
}

// SetErrorMsg sets the ErrorMsg value
func (cv *MerchantCheckoutValues) SetErrorMsg(value string) *MerchantCheckoutValues {
	cv.ErrorMsg = &value
	return cv
}

// SetExpiredAt sets the ExpiredAt value
func (cv *MerchantCheckoutValues) SetExpiredAt(value int64) *MerchantCheckoutValues {
	cv.ExpiredAt = &value
	return cv
}

// SetCompletedAt sets the CompletedAt value
func (cv *MerchantCheckoutValues) SetCompletedAt(value int64) *MerchantCheckoutValues {
	cv.CompletedAt = &value
	return cv
}

// SetTrxApp sets the TrxApp value
func (cv *MerchantCheckoutValues) SetTrxApp(value string) *MerchantCheckoutValues {
	cv.TrxApp = &value
	return cv
}

// SetCheckoutURL sets the CheckoutURL value
func (cv *MerchantCheckoutValues) SetCheckoutURL(value string) *MerchantCheckoutValues {
	cv.CheckoutURL = &value
	return cv
}

// SetSubmitedAt sets the SubmitedAt value
func (cv *MerchantCheckoutValues) SetSubmitedAt(value int64) *MerchantCheckoutValues {
	cv.SubmitedAt = &value
	return cv
}

// SetConfirmedAt sets the ConfirmedAt value
func (cv *MerchantCheckoutValues) SetConfirmedAt(value int64) *MerchantCheckoutValues {
	cv.ConfirmedAt = &value
	return cv
}

// SetCanceledAt sets the CanceledAt value
func (cv *MerchantCheckoutValues) SetCanceledAt(value int64) *MerchantCheckoutValues {
	cv.CanceledAt = &value
	return cv
}

// SetValues sets multiple CheckoutValues fields at once
func (c *MerchantCheckout) SetValues(values *MerchantCheckoutValues) *MerchantCheckout {
	if values == nil {
		return c
	}

	if c.MerchantCheckoutValues == nil {
		c.MerchantCheckoutValues = &MerchantCheckoutValues{}
	}

	// Set all fields from the provided values
	if values.Ccy != nil {
		c.MerchantCheckoutValues.SetCcy(*values.Ccy)
	}
	if values.Amount != nil {
		c.MerchantCheckoutValues.SetAmount(*values.Amount)
	}
	if values.Country != nil {
		c.MerchantCheckoutValues.SetCountry(*values.Country)
	}
	if values.TrxID != nil {
		c.MerchantCheckoutValues.SetTrxID(*values.TrxID)
	}
	if values.TrxApp != nil {
		c.MerchantCheckoutValues.SetTrxApp(*values.TrxApp)
	}
	if values.TrxMethod != nil {
		c.MerchantCheckoutValues.SetTrxMethod(*values.TrxMethod)
	}
	if values.ReturnURL != nil {
		c.MerchantCheckoutValues.SetReturnURL(*values.ReturnURL)
	}
	if values.NotifyURL != nil {
		c.MerchantCheckoutValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.Status != nil {
		c.MerchantCheckoutValues.SetStatus(*values.Status)
	}
	if values.ChannelCode != nil {
		c.MerchantCheckoutValues.SetChannelCode(*values.ChannelCode)
	}
	if values.CheckoutURL != nil {
		c.MerchantCheckoutValues.SetCheckoutURL(*values.CheckoutURL)
	}
	// Metadata is not a pointer, so we always set it
	c.MerchantCheckoutValues.SetMetadata(values.GetMetadata())
	if values.ErrorCode != nil {
		c.MerchantCheckoutValues.SetErrorCode(*values.ErrorCode)
	}
	if values.ErrorMsg != nil {
		c.MerchantCheckoutValues.SetErrorMsg(*values.ErrorMsg)
	}
	if values.ExpiredAt != nil {
		c.MerchantCheckoutValues.SetExpiredAt(*values.ExpiredAt)
	}
	if values.SubmitedAt != nil {
		c.MerchantCheckoutValues.SetSubmitedAt(*values.SubmitedAt)
	}
	if values.ConfirmedAt != nil {
		c.MerchantCheckoutValues.SetConfirmedAt(*values.ConfirmedAt)
	}
	if values.CanceledAt != nil {
		c.MerchantCheckoutValues.SetCanceledAt(*values.CanceledAt)
	}
	if values.CompletedAt != nil {
		c.MerchantCheckoutValues.SetCompletedAt(*values.CompletedAt)
	}
	// Set Transactions field
	if values.Transactions != nil {
		c.MerchantCheckoutValues.SetTransactions(values.Transactions)
	}

	return c
}

func (t *MerchantCheckout) CanCancel() bool {
	return !slices.Contains(protocol.MerchantCheckoutStatusList, t.GetStatus())
}

// Database Operations

// GetMerchantCheckoutByReqID 根据商户ID和请求ID获取收银台记录
func GetMerchantCheckoutByReqID(mid, reqID string) *MerchantCheckout {
	var checkout MerchantCheckout
	if err := ReadDB.Where("mid = ? AND req_id = ?", mid, reqID).First(&checkout).Error; err == nil {
		return &checkout
	}
	return nil
}

// GetMerchantCheckoutByCheckoutID 根据收银台ID获取记录
func GetMerchantCheckoutByCheckoutID(checkoutID string) *MerchantCheckout {
	var checkout MerchantCheckout
	if err := ReadDB.Where("checkout_id = ?", checkoutID).First(&checkout).Error; err != nil {
		return nil
	}
	return &checkout
}

func SaveMerchantCheckout(tx *gorm.DB, checkout *MerchantCheckout, values *MerchantCheckoutValues) (err error) {
	defer func() {
		if err == nil {
			checkout.SetValues(values)
		}
	}()
	return tx.Model(checkout).UpdateColumns(values).Error
}

// Transactions字段操作方法

// AddTransaction 添加交易记录到Checkout.Transactions
func (c *MerchantCheckout) AddTransaction(trx *Transaction) {
	if c.MerchantCheckoutValues == nil {
		c.MerchantCheckoutValues = &MerchantCheckoutValues{}
	}
	if c.MerchantCheckoutValues.Transactions == nil {
		c.MerchantCheckoutValues.Transactions = []*Transaction{}
	}
	c.MerchantCheckoutValues.Transactions = append(c.MerchantCheckoutValues.Transactions, trx)
}

// FindTransactionByTrxMethod 根据支付方式查找交易记录
func (c *MerchantCheckout) FindTransactionByTrxMethod(trxMethod string) *Transaction {
	for _, trx := range c.GetTransactions() {
		if trx.TrxMethod == trxMethod {
			return trx
		}
	}
	return nil
}

// FindTransactionByTrxID 根据交易ID查找交易记录
func (c *MerchantCheckout) FindTransactionByTrxID(trxID string) *Transaction {
	for _, trx := range c.GetTransactions() {
		if trx.TrxID == trxID {
			return trx
		}
	}
	return nil
}

// GetTransactions 获取所有交易记录
func (c *MerchantCheckout) GetTransactions() []*Transaction {
	if c.MerchantCheckoutValues == nil || c.MerchantCheckoutValues.Transactions == nil {
		return []*Transaction{}
	}
	return c.Transactions
}

// SetTransactions sets the Transactions value
func (cv *MerchantCheckoutValues) SetTransactions(txs []*Transaction) *MerchantCheckoutValues {
	cv.Transactions = txs
	return cv
}

// Protocol conversion method
func (c *MerchantCheckout) Protocol() *protocol.Checkout {
	if c == nil {
		return nil
	}

	return &protocol.Checkout{
		CheckoutID:  c.CheckoutID,
		Mid:         c.Mid,
		ReqID:       c.ReqID,
		TrxID:       c.GetTrxID(),
		Amount:      c.GetAmount().String(),
		Ccy:         c.GetCcy(),
		Country:     c.GetCountry(),
		TrxMethod:   c.GetTrxMethod(),
		Status:      c.GetStatus(),
		NotifyURL:   c.GetNotifyURL(),
		ReturnURL:   c.GetReturnURL(),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		ExpiredAt:   c.GetExpiredAt(),
		ErrorCode:   c.GetErrorCode(),
		ErrorMsg:    c.GetErrorMsg(),
		CheckoutURL: c.GetCheckoutURL(),
		SubmitedAt:  c.GetSubmitedAt(),
		ConfirmedAt: c.GetConfirmedAt(),
		CanceledAt:  c.GetCanceledAt(),
		CompletedAt: c.GetCompletedAt(),
	}
}
