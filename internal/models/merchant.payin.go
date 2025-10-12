package models

import (
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
)

// MerchantPayin 代收记录表
type MerchantPayin struct {
	ID                   uint64           `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TrxID                string           `json:"transaction_id" gorm:"column:transaction_id;type:varchar(64);uniqueIndex"`
	TrxType              string           `json:"trx_type" gorm:"column:trx_type;type:varchar(16);index;default:'payin'"`
	Mid                  string           `json:"mid" gorm:"column:mid;type:varchar(32);index"`
	UserID               string           `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	ReqID                string           `json:"req_id" gorm:"column:req_id;type:varchar(64);index"`
	OriTrxID             string           `json:"ori_trx_id" gorm:"column:ori_trx_id;index;<-:create"`
	OriReqID             string           `json:"ori_req_id" gorm:"column:ori_req_id;index;<-:create"`
	OriFlowNo            string           `json:"ori_flow_no" gorm:"column:ori_flow_no"`
	TrxMethod            string           `json:"trx_method" gorm:"column:trx_method;<-:create"`
	TrxMode              string           `json:"trx_mode" gorm:"column:trx_mode;<-:create"`
	TrxApp               string           `json:"trx_app" gorm:"column:trx_app;<-:create"`
	Pkg                  string           `json:"pkg" gorm:"column:pkg;<-:create"`
	Did                  string           `json:"did" gorm:"column:did;<-:create"`
	ProductID            string           `json:"product_id" gorm:"column:product_id;<-:create"`
	UserIP               string           `json:"user_ip" gorm:"column:user_ip;<-:create"`
	Email                string           `json:"email" gorm:"column:email;<-:create"`
	Phone                string           `json:"phone" gorm:"column:phone;<-:create"`
	Ccy                  string           `json:"ccy" gorm:"column:ccy;<-:create"`
	Amount               *decimal.Decimal `json:"amount" gorm:"column:amount;<-:create"`
	UsdAmount            *decimal.Decimal `json:"usd_amount" gorm:"column:usd_amount;<-:create"`
	AccountNo            string           `json:"account_no" gorm:"column:account_no;<-:create"`
	AccountName          string           `json:"account_name" gorm:"column:account_name;<-:create"`
	AccountType          string           `json:"account_type" gorm:"column:account_type;<-:create"`
	BankCode             string           `json:"bank_code" gorm:"column:bank_code;<-:create"`
	BankName             string           `json:"bank_name" gorm:"column:bank_name;<-:create"`
	ReturnURL            string           `json:"return_url" gorm:"column:return_url;<-:create"`
	*MerchantPayinValues `gorm:"embedded"`
	CreatedAt            int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt            int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

type MerchantPayinValues struct {
	MetaData *protocol.MapData `json:"metadata" gorm:"column:metadata;serializer:json;type:json"`

	// Refund related fields
	RefundedCount     *int             `json:"refunded_count" gorm:"column:refunded_count"`
	RefundedAmount    *decimal.Decimal `json:"refunded_amount" gorm:"column:refunded_amount"`
	RefundedUsdAmount *decimal.Decimal `json:"refunded_usd_amount" gorm:"column:refunded_usd_amount"`
	LastRefundedAt    *int64           `json:"last_refunded_at" gorm:"column:last_refunded_at"`

	// Settlement related fields
	FlowNo       *string `json:"flow_no" gorm:"column:flow_no"`
	SettleStatus *string `json:"settle_status" gorm:"column:settle_status;index"` // SettleStatus 结算状态
	SettleID     *string `json:"settle_id" gorm:"column:settle_id;index"`
	SettledAt    *int64  `json:"settled_at" gorm:"column:settled_at"`

	// Basic transaction fields
	Country   *string        `json:"country" gorm:"column:country"`
	Remark    *string        `json:"remark" gorm:"column:remark"`
	Status    *string        `json:"status" gorm:"column:status;index"`
	Reason    *string        `json:"reason" gorm:"column:reason"`
	Link      *string        `json:"link" gorm:"column:link"`
	Detail    map[string]any `json:"detail" gorm:"column:detail;serializer:json;type:json"`
	NotifyURL *string        `json:"notify_url" gorm:"column:notify_url"`

	// Fee related fields
	FeeCcy       *string          `json:"fee_ccy" gorm:"column:fee_ccy"`
	FeeAmount    *decimal.Decimal `json:"fee_amount" gorm:"column:fee_amount"`
	FeeUsdAmount *decimal.Decimal `json:"fee_usd_amount" gorm:"column:fee_usd_amount"`
	FeeUsdRate   *decimal.Decimal `json:"fee_usd_rate" gorm:"column:fee_usd_rate"`

	// Channel related fields
	ChannelStatus       *string          `json:"channel_status" gorm:"column:channel_status"`
	ResCode             *string          `json:"res_code" gorm:"column:res_code"`
	ResMsg              *string          `json:"res_msg" gorm:"column:res_msg"`
	ChannelTrxID        *string          `json:"channel_trx_id" gorm:"column:channel_trx_id;index"`
	ChannelCode         *string          `json:"channel_code" gorm:"column:channel_code;index"`
	ChannelAccount      *string          `json:"channel_account" gorm:"column:channel_account"`
	ChannelGroup        *string          `json:"channel_group" gorm:"column:channel_group"`
	ChannelFeeCcy       *string          `json:"channel_fee_ccy" gorm:"column:channel_fee_ccy"`
	ChannelFeeAmount    *decimal.Decimal `json:"channel_fee_amount" gorm:"column:channel_fee_amount"`
	ChannelFeeUsdAmount *decimal.Decimal `json:"channel_fee_usd_amount" gorm:"column:channel_fee_usd_amount"`
	ChannelFeeUsdRate   *decimal.Decimal `json:"channel_fee_usd_rate" gorm:"column:channel_fee_usd_rate"`

	// Timing fields
	ConfirmedAt        *int64  `json:"confirmed_at" gorm:"column:confirmed_at"`
	CompletedAt        *int64  `json:"completed_at" gorm:"column:completed_at"`
	ExpiredAt          *int64  `json:"expired_at" gorm:"column:expired_at"`
	CanceledAt         *int64  `json:"canceled_at" gorm:"column:canceled_at"`
	CancelReason       *string `json:"cancel_reason" gorm:"column:cancel_reason"`
	CancelFailedResult *string `json:"cancel_failed_result" gorm:"column:cancel_failed_result"`
	Version            *int64  `json:"version" gorm:"column:version"`
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

// GetChannelCode returns the ChannelCode value
func (pv *MerchantPayinValues) GetChannelCode() string {
	if pv.ChannelCode == nil {
		return ""
	}
	return *pv.ChannelCode
}

// GetNotifyURL returns the NotifyURL value
func (pv *MerchantPayinValues) GetNotifyURL() string {
	if pv.NotifyURL == nil {
		return ""
	}
	return *pv.NotifyURL
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

// GetRemark returns the Remark value
func (pv *MerchantPayinValues) GetRemark() string {
	if pv.Remark == nil {
		return ""
	}
	return *pv.Remark
}

// GetSettleID returns the SettleID value
func (pv *MerchantPayinValues) GetSettleID() string {
	if pv.SettleID == nil {
		return ""
	}
	return *pv.SettleID
}

// GetSettleStatus returns the SettleStatus value
func (pv *MerchantPayinValues) GetSettleStatus() string {
	if pv.SettleStatus == nil {
		return ""
	}
	return *pv.SettleStatus
}

// GetSettledAt returns the SettledAt value
func (pv *MerchantPayinValues) GetSettledAt() int64 {
	if pv.SettledAt == nil {
		return 0
	}
	return *pv.SettledAt
}

// GetRefundedCount returns the RefundedCount value
func (pv *MerchantPayinValues) GetRefundedCount() int {
	if pv.RefundedCount == nil {
		return 0
	}
	return *pv.RefundedCount
}

// GetRefundedAmount returns the RefundedAmount value
func (pv *MerchantPayinValues) GetRefundedAmount() decimal.Decimal {
	if pv.RefundedAmount == nil {
		return decimal.Zero
	}
	return *pv.RefundedAmount
}

// GetRefundedUsdAmount returns the RefundedUsdAmount value
func (pv *MerchantPayinValues) GetRefundedUsdAmount() decimal.Decimal {
	if pv.RefundedUsdAmount == nil {
		return decimal.Zero
	}
	return *pv.RefundedUsdAmount
}

// GetLastRefundedAt returns the LastRefundedAt value
func (pv *MerchantPayinValues) GetLastRefundedAt() int64 {
	if pv.LastRefundedAt == nil {
		return 0
	}
	return *pv.LastRefundedAt
}

// GetFlowNo returns the FlowNo value
func (pv *MerchantPayinValues) GetFlowNo() string {
	if pv.FlowNo == nil {
		return ""
	}
	return *pv.FlowNo
}

// GetChannelStatus returns the ChannelStatus value
func (pv *MerchantPayinValues) GetChannelStatus() string {
	if pv.ChannelStatus == nil {
		return ""
	}
	return *pv.ChannelStatus
}

// GetResCode returns the ResCode value
func (pv *MerchantPayinValues) GetResCode() string {
	if pv.ResCode == nil {
		return ""
	}
	return *pv.ResCode
}

// GetResMsg returns the ResMsg value
func (pv *MerchantPayinValues) GetResMsg() string {
	if pv.ResMsg == nil {
		return ""
	}
	return *pv.ResMsg
}

// GetReason returns the Reason value
func (pv *MerchantPayinValues) GetReason() string {
	if pv.Reason == nil {
		return ""
	}
	return *pv.Reason
}

// GetLink returns the Link value
func (pv *MerchantPayinValues) GetLink() string {
	if pv.Link == nil {
		return ""
	}
	return *pv.Link
}

// GetFeeCcy returns the FeeCcy value
func (pv *MerchantPayinValues) GetFeeCcy() string {
	if pv.FeeCcy == nil {
		return ""
	}
	return *pv.FeeCcy
}

// GetFeeAmount returns the FeeAmount value
func (pv *MerchantPayinValues) GetFeeAmount() decimal.Decimal {
	if pv.FeeAmount == nil {
		return decimal.Zero
	}
	return *pv.FeeAmount
}

// GetFeeUsdAmount returns the FeeUsdAmount value
func (pv *MerchantPayinValues) GetFeeUsdAmount() decimal.Decimal {
	if pv.FeeUsdAmount == nil {
		return decimal.Zero
	}
	return *pv.FeeUsdAmount
}

// GetFeeUsdRate returns the FeeUsdRate value
func (pv *MerchantPayinValues) GetFeeUsdRate() decimal.Decimal {
	if pv.FeeUsdRate == nil {
		return decimal.Zero
	}
	return *pv.FeeUsdRate
}

// GetChannelTrxID returns the ChannelTrxID value
func (pv *MerchantPayinValues) GetChannelTrxID() string {
	if pv.ChannelTrxID == nil {
		return ""
	}
	return *pv.ChannelTrxID
}

// GetChannelAccount returns the ChannelAccount value
func (pv *MerchantPayinValues) GetChannelAccount() string {
	if pv.ChannelAccount == nil {
		return ""
	}
	return *pv.ChannelAccount
}

// GetChannelGroup returns the ChannelGroup value
func (pv *MerchantPayinValues) GetChannelGroup() string {
	if pv.ChannelGroup == nil {
		return ""
	}
	return *pv.ChannelGroup
}

// GetChannelFeeCcy returns the ChannelFeeCcy value
func (pv *MerchantPayinValues) GetChannelFeeCcy() string {
	if pv.ChannelFeeCcy == nil {
		return ""
	}
	return *pv.ChannelFeeCcy
}

// GetChannelFeeAmount returns the ChannelFeeAmount value
func (pv *MerchantPayinValues) GetChannelFeeAmount() decimal.Decimal {
	if pv.ChannelFeeAmount == nil {
		return decimal.Zero
	}
	return *pv.ChannelFeeAmount
}

// GetChannelFeeUsdAmount returns the ChannelFeeUsdAmount value
func (pv *MerchantPayinValues) GetChannelFeeUsdAmount() decimal.Decimal {
	if pv.ChannelFeeUsdAmount == nil {
		return decimal.Zero
	}
	return *pv.ChannelFeeUsdAmount
}

// GetChannelFeeUsdRate returns the ChannelFeeUsdRate value
func (pv *MerchantPayinValues) GetChannelFeeUsdRate() decimal.Decimal {
	if pv.ChannelFeeUsdRate == nil {
		return decimal.Zero
	}
	return *pv.ChannelFeeUsdRate
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

// SetChannelCode sets the ChannelCode value
func (pv *MerchantPayinValues) SetChannelCode(value string) *MerchantPayinValues {
	pv.ChannelCode = &value
	return pv
}

// SetNotifyURL sets the NotifyURL value
func (pv *MerchantPayinValues) SetNotifyURL(value string) *MerchantPayinValues {
	pv.NotifyURL = &value
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

// SetRemark sets the Remark value
func (pv *MerchantPayinValues) SetRemark(value string) *MerchantPayinValues {
	pv.Remark = &value
	return pv
}

// SetSettleID sets the SettleID value
func (pv *MerchantPayinValues) SetSettleID(value string) *MerchantPayinValues {
	pv.SettleID = &value
	return pv
}

// SetSettleStatus sets the SettleStatus value
func (pv *MerchantPayinValues) SetSettleStatus(value string) *MerchantPayinValues {
	pv.SettleStatus = &value
	return pv
}

// SetSettledAt sets the SettledAt value
func (pv *MerchantPayinValues) SetSettledAt(value int64) *MerchantPayinValues {
	pv.SettledAt = &value
	return pv
}

// SetRefundedCount sets the RefundedCount value
func (pv *MerchantPayinValues) SetRefundedCount(value int) *MerchantPayinValues {
	pv.RefundedCount = &value
	return pv
}

// SetRefundedAmount sets the RefundedAmount value
func (pv *MerchantPayinValues) SetRefundedAmount(value decimal.Decimal) *MerchantPayinValues {
	pv.RefundedAmount = &value
	return pv
}

// SetRefundedUsdAmount sets the RefundedUsdAmount value
func (pv *MerchantPayinValues) SetRefundedUsdAmount(value decimal.Decimal) *MerchantPayinValues {
	pv.RefundedUsdAmount = &value
	return pv
}

// SetLastRefundedAt sets the LastRefundedAt value
func (pv *MerchantPayinValues) SetLastRefundedAt(value int64) *MerchantPayinValues {
	pv.LastRefundedAt = &value
	return pv
}

// SetFlowNo sets the FlowNo value
func (pv *MerchantPayinValues) SetFlowNo(value string) *MerchantPayinValues {
	pv.FlowNo = &value
	return pv
}

// SetChannelStatus sets the ChannelStatus value
func (pv *MerchantPayinValues) SetChannelStatus(value string) *MerchantPayinValues {
	pv.ChannelStatus = &value
	return pv
}

// SetResCode sets the ResCode value
func (pv *MerchantPayinValues) SetResCode(value string) *MerchantPayinValues {
	pv.ResCode = &value
	return pv
}

// SetResMsg sets the ResMsg value
func (pv *MerchantPayinValues) SetResMsg(value string) *MerchantPayinValues {
	pv.ResMsg = &value
	return pv
}

// SetReason sets the Reason value
func (pv *MerchantPayinValues) SetReason(value string) *MerchantPayinValues {
	pv.Reason = &value
	return pv
}

// SetLink sets the Link value
func (pv *MerchantPayinValues) SetLink(value string) *MerchantPayinValues {
	pv.Link = &value
	return pv
}

// SetFeeCcy sets the FeeCcy value
func (pv *MerchantPayinValues) SetFeeCcy(value string) *MerchantPayinValues {
	pv.FeeCcy = &value
	return pv
}

// SetFeeAmount sets the FeeAmount value
func (pv *MerchantPayinValues) SetFeeAmount(value decimal.Decimal) *MerchantPayinValues {
	pv.FeeAmount = &value
	return pv
}

// SetFeeUsdAmount sets the FeeUsdAmount value
func (pv *MerchantPayinValues) SetFeeUsdAmount(value decimal.Decimal) *MerchantPayinValues {
	pv.FeeUsdAmount = &value
	return pv
}

// SetFeeUsdRate sets the FeeUsdRate value
func (pv *MerchantPayinValues) SetFeeUsdRate(value decimal.Decimal) *MerchantPayinValues {
	pv.FeeUsdRate = &value
	return pv
}

// SetChannelTrxID sets the ChannelTrxID value
func (pv *MerchantPayinValues) SetChannelTrxID(value string) *MerchantPayinValues {
	pv.ChannelTrxID = &value
	return pv
}

// SetChannelAccount sets the ChannelAccount value
func (pv *MerchantPayinValues) SetChannelAccount(value string) *MerchantPayinValues {
	pv.ChannelAccount = &value
	return pv
}

// SetChannelGroup sets the ChannelGroup value
func (pv *MerchantPayinValues) SetChannelGroup(value string) *MerchantPayinValues {
	pv.ChannelGroup = &value
	return pv
}

// SetChannelFeeCcy sets the ChannelFeeCcy value
func (pv *MerchantPayinValues) SetChannelFeeCcy(value string) *MerchantPayinValues {
	pv.ChannelFeeCcy = &value
	return pv
}

// SetChannelFeeAmount sets the ChannelFeeAmount value
func (pv *MerchantPayinValues) SetChannelFeeAmount(value decimal.Decimal) *MerchantPayinValues {
	pv.ChannelFeeAmount = &value
	return pv
}

// SetChannelFeeUsdAmount sets the ChannelFeeUsdAmount value
func (pv *MerchantPayinValues) SetChannelFeeUsdAmount(value decimal.Decimal) *MerchantPayinValues {
	pv.ChannelFeeUsdAmount = &value
	return pv
}

// SetChannelFeeUsdRate sets the ChannelFeeUsdRate value
func (pv *MerchantPayinValues) SetChannelFeeUsdRate(value decimal.Decimal) *MerchantPayinValues {
	pv.ChannelFeeUsdRate = &value
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
	if values.ChannelCode != nil {
		p.MerchantPayinValues.SetChannelCode(*values.ChannelCode)
	}
	if values.NotifyURL != nil {
		p.MerchantPayinValues.SetNotifyURL(*values.NotifyURL)
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

	if values.Remark != nil {
		p.MerchantPayinValues.SetRemark(*values.Remark)
	}

	if values.SettleID != nil {
		p.MerchantPayinValues.SetSettleID(*values.SettleID)
	}
	if values.SettleStatus != nil {
		p.MerchantPayinValues.SetSettleStatus(*values.SettleStatus)
	}
	if values.SettledAt != nil {
		p.MerchantPayinValues.SetSettledAt(*values.SettledAt)
	}
	if values.RefundedCount != nil {
		p.MerchantPayinValues.SetRefundedCount(*values.RefundedCount)
	}
	if values.RefundedAmount != nil {
		p.MerchantPayinValues.SetRefundedAmount(*values.RefundedAmount)
	}
	if values.RefundedUsdAmount != nil {
		p.MerchantPayinValues.SetRefundedUsdAmount(*values.RefundedUsdAmount)
	}
	if values.LastRefundedAt != nil {
		p.MerchantPayinValues.SetLastRefundedAt(*values.LastRefundedAt)
	}
	if values.FlowNo != nil {
		p.MerchantPayinValues.SetFlowNo(*values.FlowNo)
	}
	if values.ChannelStatus != nil {
		p.MerchantPayinValues.SetChannelStatus(*values.ChannelStatus)
	}
	if values.ResCode != nil {
		p.MerchantPayinValues.SetResCode(*values.ResCode)
	}
	if values.ResMsg != nil {
		p.MerchantPayinValues.SetResMsg(*values.ResMsg)
	}
	if values.Reason != nil {
		p.MerchantPayinValues.SetReason(*values.Reason)
	}
	if values.Link != nil {
		p.MerchantPayinValues.SetLink(*values.Link)
	}
	if values.FeeCcy != nil {
		p.MerchantPayinValues.SetFeeCcy(*values.FeeCcy)
	}
	if values.FeeAmount != nil {
		p.MerchantPayinValues.SetFeeAmount(*values.FeeAmount)
	}
	if values.FeeUsdAmount != nil {
		p.MerchantPayinValues.SetFeeUsdAmount(*values.FeeUsdAmount)
	}
	if values.FeeUsdRate != nil {
		p.MerchantPayinValues.SetFeeUsdRate(*values.FeeUsdRate)
	}
	if values.ChannelTrxID != nil {
		p.MerchantPayinValues.SetChannelTrxID(*values.ChannelTrxID)
	}
	if values.ChannelAccount != nil {
		p.MerchantPayinValues.SetChannelAccount(*values.ChannelAccount)
	}
	if values.ChannelGroup != nil {
		p.MerchantPayinValues.SetChannelGroup(*values.ChannelGroup)
	}
	if values.ChannelFeeCcy != nil {
		p.MerchantPayinValues.SetChannelFeeCcy(*values.ChannelFeeCcy)
	}
	if values.ChannelFeeAmount != nil {
		p.MerchantPayinValues.SetChannelFeeAmount(*values.ChannelFeeAmount)
	}
	if values.ChannelFeeUsdAmount != nil {
		p.MerchantPayinValues.SetChannelFeeUsdAmount(*values.ChannelFeeUsdAmount)
	}
	if values.ChannelFeeUsdRate != nil {
		p.MerchantPayinValues.SetChannelFeeUsdRate(*values.ChannelFeeUsdRate)
	}

	return p
}

// ToTransaction converts Payin to Transaction
func (p *MerchantPayin) ToTransaction() *Transaction {
	if p == nil {
		return nil
	}

	transaction := &Transaction{
		ID:          int64(p.ID), // Convert uint64 to int64
		Mid:         p.Mid,
		TrxType:     protocol.TrxTypePayin,
		UserID:      p.UserID,
		TrxID:       p.TrxID,
		ReqID:       p.ReqID,
		OriTrxID:    p.OriTrxID,
		OriReqID:    p.OriReqID,
		OriFlowNo:   p.OriFlowNo,
		TrxMethod:   p.TrxMethod,
		TrxMode:     p.TrxMode,
		TrxApp:      p.TrxApp,
		Pkg:         p.Pkg,
		Did:         p.Did,
		ProductID:   p.ProductID,
		UserIP:      p.UserIP,
		Email:       p.Email,
		Phone:       p.Phone,
		Ccy:         p.Ccy,
		Amount:      p.Amount,
		UsdAmount:   p.UsdAmount,
		AccountNo:   p.AccountNo,
		AccountName: p.AccountName,
		AccountType: p.AccountType,
		BankCode:    p.BankCode,
		BankName:    p.BankName,
		ReturnURL:   p.ReturnURL,
		TransactionValues: &TransactionValues{
			MetaData:            p.MerchantPayinValues.MetaData,
			RefundedCount:       p.MerchantPayinValues.RefundedCount,
			RefundedAmount:      p.MerchantPayinValues.RefundedAmount,
			RefundedUsdAmount:   p.MerchantPayinValues.RefundedUsdAmount,
			LastRefundedAt:      p.MerchantPayinValues.LastRefundedAt,
			SettleStatus:        p.MerchantPayinValues.SettleStatus,
			SettleID:            p.MerchantPayinValues.SettleID,
			SettledAt:           p.MerchantPayinValues.SettledAt,
			Country:             p.MerchantPayinValues.Country,
			Remark:              p.MerchantPayinValues.Remark,
			FlowNo:              p.MerchantPayinValues.FlowNo,
			Status:              p.MerchantPayinValues.Status,
			Reason:              p.MerchantPayinValues.Reason,
			Link:                p.MerchantPayinValues.Link,
			Detail:              p.MerchantPayinValues.Detail,
			NotifyURL:           p.MerchantPayinValues.NotifyURL,
			FeeCcy:              p.MerchantPayinValues.FeeCcy,
			FeeAmount:           p.MerchantPayinValues.FeeAmount,
			FeeUsdAmount:        p.MerchantPayinValues.FeeUsdAmount,
			FeeUsdRate:          p.MerchantPayinValues.FeeUsdRate,
			ChannelStatus:       p.MerchantPayinValues.ChannelStatus,
			ResCode:             p.MerchantPayinValues.ResCode,
			ResMsg:              p.MerchantPayinValues.ResMsg,
			ChannelTrxID:        p.MerchantPayinValues.ChannelTrxID,
			ChannelCode:         p.MerchantPayinValues.ChannelCode,
			ChannelAccount:      p.MerchantPayinValues.ChannelAccount,
			ChannelGroup:        p.MerchantPayinValues.ChannelGroup,
			ChannelFeeCcy:       p.MerchantPayinValues.ChannelFeeCcy,
			ChannelFeeAmount:    p.MerchantPayinValues.ChannelFeeAmount,
			ChannelFeeUsdAmount: p.MerchantPayinValues.ChannelFeeUsdAmount,
			ChannelFeeUsdRate:   p.MerchantPayinValues.ChannelFeeUsdRate,
			ConfirmedAt:         p.MerchantPayinValues.ConfirmedAt,
			CompletedAt:         p.MerchantPayinValues.CompletedAt,
			ExpiredAt:           p.MerchantPayinValues.ExpiredAt,
			CanceledAt:          p.MerchantPayinValues.CanceledAt,
			CancelReason:        p.MerchantPayinValues.CancelReason,
			CancelFailedResult:  p.MerchantPayinValues.CancelFailedResult,
			Version:             p.MerchantPayinValues.Version,
		},
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	return transaction
}

func GetMerchantPayinByReqID(mid, reqID string) *MerchantPayin {
	var existingTrx MerchantPayin
	if err := ReadDB.Where("mid = ? AND req_id = ?", mid, reqID).First(&existingTrx).Error; err == nil {
		return &existingTrx
	}
	return nil
}

func GetMerchantPayinByTrxID(mid, trxID string) *Transaction {
	var trx Transaction
	if err := ReadDB.Model(&MerchantPayin{}).Where("trx_id = ? AND mid = ?", trxID, mid).First(&trx).Error; err == nil {
		return &trx
	}
	return nil
}
