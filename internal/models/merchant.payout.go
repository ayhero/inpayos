package models

import (
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
)

// MerchantPayout 代付记录表
type MerchantPayout struct {
	ID                    uint64           `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TrxID                 string           `json:"trx_id" gorm:"column:trx_id;type:varchar(64);uniqueIndex"`
	TrxType               string           `json:"trx_type" gorm:"column:trx_type;type:varchar(16);index;default:'payin'"`
	Mid                   string           `json:"mid" gorm:"column:mid;type:varchar(32);index"`
	UserID                string           `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	ReqID                 string           `json:"req_id" gorm:"column:req_id;type:varchar(64);index"`
	OriTrxID              string           `json:"ori_trx_id" gorm:"column:ori_trx_id;index;<-:create"`
	OriReqID              string           `json:"ori_req_id" gorm:"column:ori_req_id;index;<-:create"`
	OriFlowNo             string           `json:"ori_flow_no" gorm:"column:ori_flow_no"`
	TrxMethod             string           `json:"trx_method" gorm:"column:trx_method;<-:create"`
	TrxMode               string           `json:"trx_mode" gorm:"column:trx_mode;<-:create"`
	TrxApp                string           `json:"trx_app" gorm:"column:trx_app;<-:create"`
	Pkg                   string           `json:"pkg" gorm:"column:pkg;<-:create"`
	Did                   string           `json:"did" gorm:"column:did;<-:create"`
	ProductID             string           `json:"product_id" gorm:"column:product_id;<-:create"`
	UserIP                string           `json:"user_ip" gorm:"column:user_ip;<-:create"`
	Email                 string           `json:"email" gorm:"column:email;<-:create"`
	Phone                 string           `json:"phone" gorm:"column:phone;<-:create"`
	Ccy                   string           `json:"ccy" gorm:"column:ccy;<-:create"`
	Amount                *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(19,4);<-:create"`
	UsdAmount             *decimal.Decimal `json:"usd_amount" gorm:"column:usd_amount;type:decimal(19,4);<-:create"`
	AccountNo             string           `json:"account_no" gorm:"column:account_no;<-:create"`
	AccountName           string           `json:"account_name" gorm:"column:account_name;<-:create"`
	AccountType           string           `json:"account_type" gorm:"column:account_type;<-:create"`
	BankCode              string           `json:"bank_code" gorm:"column:bank_code;<-:create"`
	BankName              string           `json:"bank_name" gorm:"column:bank_name;<-:create"`
	ReturnURL             string           `json:"return_url" gorm:"column:return_url;<-:create"`
	*MerchantPayoutValues `gorm:"embedded"`
	CreatedAt             int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt             int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

type MerchantPayoutValues struct {
	MetaData *protocol.MapData `json:"metadata" gorm:"column:metadata;serializer:json;type:json"`

	// Refund related fields
	RefundedCount     *int             `json:"refunded_count" gorm:"column:refunded_count"`
	RefundedAmount    *decimal.Decimal `json:"refunded_amount" gorm:"column:refunded_amount;type:decimal(19,4)"`
	RefundedUsdAmount *decimal.Decimal `json:"refunded_usd_amount" gorm:"column:refunded_usd_amount;type:decimal(19,4)"`
	LastRefundedAt    *int64           `json:"last_refunded_at" gorm:"column:last_refunded_at"`

	// Settlement related fields
	FlowNo       *string `json:"flow_no" gorm:"column:flow_no;index"`
	SettleStatus *string `json:"settle_status" gorm:"column:settle_status;index"` // SettleStatus 结算状态
	SettleID     *string `json:"settle_id" gorm:"column:settle_id;index"`
	SettledAt    *int64  `json:"settled_at" gorm:"column:settled_at"`

	// Basic transaction fields
	Country   *string        `json:"country" gorm:"column:country"`
	Remark    *string        `json:"remark" gorm:"column:remark"`
	Status    *string        `json:"status" gorm:"column:status;index"`
	Link      *string        `json:"link" gorm:"column:link"`
	Detail    map[string]any `json:"detail" gorm:"column:detail;serializer:json;type:json"`
	NotifyURL *string        `json:"notify_url" gorm:"column:notify_url"`

	// Fee related fields
	FeeCcy       *string          `json:"fee_ccy" gorm:"column:fee_ccy"`
	FeeAmount    *decimal.Decimal `json:"fee_amount" gorm:"column:fee_amount;type:decimal(19,4)"`
	FeeUsdAmount *decimal.Decimal `json:"fee_usd_amount" gorm:"column:fee_usd_amount;type:decimal(19,4)"`
	FeeUsdRate   *decimal.Decimal `json:"fee_usd_rate" gorm:"column:fee_usd_rate;type:decimal(19,4)"`

	// Channel related fields
	ChannelStatus       *string          `json:"channel_status" gorm:"column:channel_status"`
	ResCode             *string          `json:"res_code" gorm:"column:res_code"`
	ResMsg              *string          `json:"res_msg" gorm:"column:res_msg"`
	Reason              *string          `json:"reason" gorm:"column:reason"`
	ChannelTrxID        *string          `json:"channel_trx_id" gorm:"column:channel_trx_id;index"`
	ChannelCode         *string          `json:"channel_code" gorm:"column:channel_code;index"`
	ChannelAccount      *string          `json:"channel_account" gorm:"column:channel_account"`
	ChannelGroup        *string          `json:"channel_group" gorm:"column:channel_group"`
	ChannelFeeCcy       *string          `json:"channel_fee_ccy" gorm:"column:channel_fee_ccy"`
	ChannelFeeAmount    *decimal.Decimal `json:"channel_fee_amount" gorm:"column:channel_fee_amount;type:decimal(19,4)"`
	ChannelFeeUsdAmount *decimal.Decimal `json:"channel_fee_usd_amount" gorm:"column:channel_fee_usd_amount;type:decimal(19,4)"`
	ChannelFeeUsdRate   *decimal.Decimal `json:"channel_fee_usd_rate" gorm:"column:channel_fee_usd_rate;type:decimal(19,4)"`

	// Timing fields
	ConfirmedAt        *int64  `json:"confirmed_at" gorm:"column:confirmed_at"`
	CompletedAt        *int64  `json:"completed_at" gorm:"column:completed_at"`
	ExpiredAt          *int64  `json:"expired_at" gorm:"column:expired_at"`
	CanceledAt         *int64  `json:"canceled_at" gorm:"column:canceled_at"`
	CancelReason       *string `json:"cancel_reason" gorm:"column:cancel_reason"`
	CancelFailedResult *string `json:"cancel_failed_result" gorm:"column:cancel_failed_result"`
	Version            *int64  `json:"version" gorm:"column:version"`
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

// GetChannelCode returns the ChannelCode value
func (pov *MerchantPayoutValues) GetChannelCode() string {
	if pov.ChannelCode == nil {
		return ""
	}
	return *pov.ChannelCode
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

// GetRemark returns the Remark value
func (pov *MerchantPayoutValues) GetRemark() string {
	if pov.Remark == nil {
		return ""
	}
	return *pov.Remark
}

// GetSettleID returns the SettleID value
func (pov *MerchantPayoutValues) GetSettleID() string {
	if pov.SettleID == nil {
		return ""
	}
	return *pov.SettleID
}

// GetSettleStatus returns the SettleStatus value
func (pov *MerchantPayoutValues) GetSettleStatus() string {
	if pov.SettleStatus == nil {
		return ""
	}
	return *pov.SettleStatus
}

// GetSettledAt returns the SettledAt value
func (pov *MerchantPayoutValues) GetSettledAt() int64 {
	if pov.SettledAt == nil {
		return 0
	}
	return *pov.SettledAt
}

// GetFlowNo returns the FlowNo value
func (pov *MerchantPayoutValues) GetFlowNo() string {
	if pov.FlowNo == nil {
		return ""
	}
	return *pov.FlowNo
}

// GetChannelStatus returns the ChannelStatus value
func (pov *MerchantPayoutValues) GetChannelStatus() string {
	if pov.ChannelStatus == nil {
		return ""
	}
	return *pov.ChannelStatus
}

// GetResCode returns the ResCode value
func (pov *MerchantPayoutValues) GetResCode() string {
	if pov.ResCode == nil {
		return ""
	}
	return *pov.ResCode
}

// GetResMsg returns the ResMsg value
func (pov *MerchantPayoutValues) GetResMsg() string {
	if pov.ResMsg == nil {
		return ""
	}
	return *pov.ResMsg
}

// GetReason returns the Reason value
func (pov *MerchantPayoutValues) GetReason() string {
	if pov.Reason == nil {
		return ""
	}
	return *pov.Reason
}

// GetLink returns the Link value
func (pov *MerchantPayoutValues) GetLink() string {
	if pov.Link == nil {
		return ""
	}
	return *pov.Link
}

// GetFeeCcy returns the FeeCcy value
func (pov *MerchantPayoutValues) GetFeeCcy() string {
	if pov.FeeCcy == nil {
		return ""
	}
	return *pov.FeeCcy
}

// GetFeeAmount returns the FeeAmount value
func (pov *MerchantPayoutValues) GetFeeAmount() decimal.Decimal {
	if pov.FeeAmount == nil {
		return decimal.Zero
	}
	return *pov.FeeAmount
}

// GetFeeUsdAmount returns the FeeUsdAmount value
func (pov *MerchantPayoutValues) GetFeeUsdAmount() decimal.Decimal {
	if pov.FeeUsdAmount == nil {
		return decimal.Zero
	}
	return *pov.FeeUsdAmount
}

// GetFeeUsdRate returns the FeeUsdRate value
func (pov *MerchantPayoutValues) GetFeeUsdRate() decimal.Decimal {
	if pov.FeeUsdRate == nil {
		return decimal.Zero
	}
	return *pov.FeeUsdRate
}

// GetChannelTrxID returns the ChannelTrxID value
func (pov *MerchantPayoutValues) GetChannelTrxID() string {
	if pov.ChannelTrxID == nil {
		return ""
	}
	return *pov.ChannelTrxID
}

// GetChannelAccount returns the ChannelAccount value
func (pov *MerchantPayoutValues) GetChannelAccount() string {
	if pov.ChannelAccount == nil {
		return ""
	}
	return *pov.ChannelAccount
}

// GetChannelGroup returns the ChannelGroup value
func (pov *MerchantPayoutValues) GetChannelGroup() string {
	if pov.ChannelGroup == nil {
		return ""
	}
	return *pov.ChannelGroup
}

// GetChannelFeeCcy returns the ChannelFeeCcy value
func (pov *MerchantPayoutValues) GetChannelFeeCcy() string {
	if pov.ChannelFeeCcy == nil {
		return ""
	}
	return *pov.ChannelFeeCcy
}

// GetChannelFeeAmount returns the ChannelFeeAmount value
func (pov *MerchantPayoutValues) GetChannelFeeAmount() decimal.Decimal {
	if pov.ChannelFeeAmount == nil {
		return decimal.Zero
	}
	return *pov.ChannelFeeAmount
}

// GetChannelFeeUsdAmount returns the ChannelFeeUsdAmount value
func (pov *MerchantPayoutValues) GetChannelFeeUsdAmount() decimal.Decimal {
	if pov.ChannelFeeUsdAmount == nil {
		return decimal.Zero
	}
	return *pov.ChannelFeeUsdAmount
}

// GetChannelFeeUsdRate returns the ChannelFeeUsdRate value
func (pov *MerchantPayoutValues) GetChannelFeeUsdRate() decimal.Decimal {
	if pov.ChannelFeeUsdRate == nil {
		return decimal.Zero
	}
	return *pov.ChannelFeeUsdRate
}

// GetMetaData returns the MetaData value
func (pov *MerchantPayoutValues) GetMetaData() *protocol.MapData {
	return pov.MetaData
}

// GetRefundedCount returns the RefundedCount value
func (pov *MerchantPayoutValues) GetRefundedCount() int {
	if pov.RefundedCount == nil {
		return 0
	}
	return *pov.RefundedCount
}

// GetRefundedAmount returns the RefundedAmount value
func (pov *MerchantPayoutValues) GetRefundedAmount() decimal.Decimal {
	if pov.RefundedAmount == nil {
		return decimal.Zero
	}
	return *pov.RefundedAmount
}

// GetRefundedUsdAmount returns the RefundedUsdAmount value
func (pov *MerchantPayoutValues) GetRefundedUsdAmount() decimal.Decimal {
	if pov.RefundedUsdAmount == nil {
		return decimal.Zero
	}
	return *pov.RefundedUsdAmount
}

// GetLastRefundedAt returns the LastRefundedAt value
func (pov *MerchantPayoutValues) GetLastRefundedAt() int64 {
	if pov.LastRefundedAt == nil {
		return 0
	}
	return *pov.LastRefundedAt
}

// GetDetail returns the Detail value
func (pov *MerchantPayoutValues) GetDetail() map[string]any {
	return pov.Detail
}

// GetConfirmedAt returns the ConfirmedAt value
func (pov *MerchantPayoutValues) GetConfirmedAt() int64 {
	if pov.ConfirmedAt == nil {
		return 0
	}
	return *pov.ConfirmedAt
}

// GetCompletedAt returns the CompletedAt value
func (pov *MerchantPayoutValues) GetCompletedAt() int64 {
	if pov.CompletedAt == nil {
		return 0
	}
	return *pov.CompletedAt
}

// GetCancelReason returns the CancelReason value
func (pov *MerchantPayoutValues) GetCancelReason() string {
	if pov.CancelReason == nil {
		return ""
	}
	return *pov.CancelReason
}

// GetCancelFailedResult returns the CancelFailedResult value
func (pov *MerchantPayoutValues) GetCancelFailedResult() string {
	if pov.CancelFailedResult == nil {
		return ""
	}
	return *pov.CancelFailedResult
}

// GetVersion returns the Version value
func (pov *MerchantPayoutValues) GetVersion() int64 {
	if pov.Version == nil {
		return 0
	}
	return *pov.Version
}

// SetStatus sets the Status value
func (pov *MerchantPayoutValues) SetStatus(value string) *MerchantPayoutValues {
	pov.Status = &value
	return pov
}

// SetChannelCode sets the ChannelCode value
func (pov *MerchantPayoutValues) SetChannelCode(value string) *MerchantPayoutValues {
	pov.ChannelCode = &value
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

// SetRemark sets the Remark value
func (pov *MerchantPayoutValues) SetRemark(value string) *MerchantPayoutValues {
	pov.Remark = &value
	return pov
}

// SetSettleID sets the SettleID value
func (pov *MerchantPayoutValues) SetSettleID(value string) *MerchantPayoutValues {
	pov.SettleID = &value
	return pov
}

// SetSettleStatus sets the SettleStatus value
func (pov *MerchantPayoutValues) SetSettleStatus(value string) *MerchantPayoutValues {
	pov.SettleStatus = &value
	return pov
}

// SetSettledAt sets the SettledAt value
func (pov *MerchantPayoutValues) SetSettledAt(value int64) *MerchantPayoutValues {
	pov.SettledAt = &value
	return pov
}

// SetFlowNo sets the FlowNo value
func (pov *MerchantPayoutValues) SetFlowNo(value string) *MerchantPayoutValues {
	pov.FlowNo = &value
	return pov
}

// SetChannelStatus sets the ChannelStatus value
func (pov *MerchantPayoutValues) SetChannelStatus(value string) *MerchantPayoutValues {
	pov.ChannelStatus = &value
	return pov
}

// SetResCode sets the ResCode value
func (pov *MerchantPayoutValues) SetResCode(value string) *MerchantPayoutValues {
	pov.ResCode = &value
	return pov
}

// SetResMsg sets the ResMsg value
func (pov *MerchantPayoutValues) SetResMsg(value string) *MerchantPayoutValues {
	pov.ResMsg = &value
	return pov
}

// SetReason sets the Reason value
func (pov *MerchantPayoutValues) SetReason(value string) *MerchantPayoutValues {
	pov.Reason = &value
	return pov
}

// SetLink sets the Link value
func (pov *MerchantPayoutValues) SetLink(value string) *MerchantPayoutValues {
	pov.Link = &value
	return pov
}

// SetFeeCcy sets the FeeCcy value
func (pov *MerchantPayoutValues) SetFeeCcy(value string) *MerchantPayoutValues {
	pov.FeeCcy = &value
	return pov
}

// SetFeeAmount sets the FeeAmount value
func (pov *MerchantPayoutValues) SetFeeAmount(value decimal.Decimal) *MerchantPayoutValues {
	pov.FeeAmount = &value
	return pov
}

// SetFeeUsdAmount sets the FeeUsdAmount value
func (pov *MerchantPayoutValues) SetFeeUsdAmount(value decimal.Decimal) *MerchantPayoutValues {
	pov.FeeUsdAmount = &value
	return pov
}

// SetFeeUsdRate sets the FeeUsdRate value
func (pov *MerchantPayoutValues) SetFeeUsdRate(value decimal.Decimal) *MerchantPayoutValues {
	pov.FeeUsdRate = &value
	return pov
}

// SetChannelTrxID sets the ChannelTrxID value
func (pov *MerchantPayoutValues) SetChannelTrxID(value string) *MerchantPayoutValues {
	pov.ChannelTrxID = &value
	return pov
}

// SetChannelAccount sets the ChannelAccount value
func (pov *MerchantPayoutValues) SetChannelAccount(value string) *MerchantPayoutValues {
	pov.ChannelAccount = &value
	return pov
}

// SetChannelGroup sets the ChannelGroup value
func (pov *MerchantPayoutValues) SetChannelGroup(value string) *MerchantPayoutValues {
	pov.ChannelGroup = &value
	return pov
}

// SetChannelFeeCcy sets the ChannelFeeCcy value
func (pov *MerchantPayoutValues) SetChannelFeeCcy(value string) *MerchantPayoutValues {
	pov.ChannelFeeCcy = &value
	return pov
}

// SetChannelFeeAmount sets the ChannelFeeAmount value
func (pov *MerchantPayoutValues) SetChannelFeeAmount(value decimal.Decimal) *MerchantPayoutValues {
	pov.ChannelFeeAmount = &value
	return pov
}

// SetChannelFeeUsdAmount sets the ChannelFeeUsdAmount value
func (pov *MerchantPayoutValues) SetChannelFeeUsdAmount(value decimal.Decimal) *MerchantPayoutValues {
	pov.ChannelFeeUsdAmount = &value
	return pov
}

// SetChannelFeeUsdRate sets the ChannelFeeUsdRate value
func (pov *MerchantPayoutValues) SetChannelFeeUsdRate(value decimal.Decimal) *MerchantPayoutValues {
	pov.ChannelFeeUsdRate = &value
	return pov
}

// SetMetaData sets the MetaData value
func (pov *MerchantPayoutValues) SetMetaData(value *protocol.MapData) *MerchantPayoutValues {
	pov.MetaData = value
	return pov
}

// SetRefundedCount sets the RefundedCount value
func (pov *MerchantPayoutValues) SetRefundedCount(value int) *MerchantPayoutValues {
	pov.RefundedCount = &value
	return pov
}

// SetRefundedAmount sets the RefundedAmount value
func (pov *MerchantPayoutValues) SetRefundedAmount(value decimal.Decimal) *MerchantPayoutValues {
	pov.RefundedAmount = &value
	return pov
}

// SetRefundedUsdAmount sets the RefundedUsdAmount value
func (pov *MerchantPayoutValues) SetRefundedUsdAmount(value decimal.Decimal) *MerchantPayoutValues {
	pov.RefundedUsdAmount = &value
	return pov
}

// SetLastRefundedAt sets the LastRefundedAt value
func (pov *MerchantPayoutValues) SetLastRefundedAt(value int64) *MerchantPayoutValues {
	pov.LastRefundedAt = &value
	return pov
}

// SetDetail sets the Detail value
func (pov *MerchantPayoutValues) SetDetail(value map[string]any) *MerchantPayoutValues {
	pov.Detail = value
	return pov
}

// SetConfirmedAt sets the ConfirmedAt value
func (pov *MerchantPayoutValues) SetConfirmedAt(value int64) *MerchantPayoutValues {
	pov.ConfirmedAt = &value
	return pov
}

// SetCompletedAt sets the CompletedAt value
func (pov *MerchantPayoutValues) SetCompletedAt(value int64) *MerchantPayoutValues {
	pov.CompletedAt = &value
	return pov
}

// SetCancelReason sets the CancelReason value
func (pov *MerchantPayoutValues) SetCancelReason(value string) *MerchantPayoutValues {
	pov.CancelReason = &value
	return pov
}

// SetCancelFailedResult sets the CancelFailedResult value
func (pov *MerchantPayoutValues) SetCancelFailedResult(value string) *MerchantPayoutValues {
	pov.CancelFailedResult = &value
	return pov
}

// SetVersion sets the Version value
func (pov *MerchantPayoutValues) SetVersion(value int64) *MerchantPayoutValues {
	pov.Version = &value
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
	if values.ChannelCode != nil {
		p.MerchantPayoutValues.SetChannelCode(*values.ChannelCode)
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

	if values.Remark != nil {
		p.MerchantPayoutValues.SetRemark(*values.Remark)
	}

	if values.SettleID != nil {
		p.MerchantPayoutValues.SetSettleID(*values.SettleID)
	}
	if values.SettleStatus != nil {
		p.MerchantPayoutValues.SetSettleStatus(*values.SettleStatus)
	}
	if values.SettledAt != nil {
		p.MerchantPayoutValues.SetSettledAt(*values.SettledAt)
	}
	if values.FlowNo != nil {
		p.MerchantPayoutValues.SetFlowNo(*values.FlowNo)
	}
	if values.ChannelStatus != nil {
		p.MerchantPayoutValues.SetChannelStatus(*values.ChannelStatus)
	}
	if values.ResCode != nil {
		p.MerchantPayoutValues.SetResCode(*values.ResCode)
	}
	if values.ResMsg != nil {
		p.MerchantPayoutValues.SetResMsg(*values.ResMsg)
	}
	if values.Reason != nil {
		p.MerchantPayoutValues.SetReason(*values.Reason)
	}
	if values.Link != nil {
		p.MerchantPayoutValues.SetLink(*values.Link)
	}
	if values.FeeCcy != nil {
		p.MerchantPayoutValues.SetFeeCcy(*values.FeeCcy)
	}
	if values.FeeAmount != nil {
		p.MerchantPayoutValues.SetFeeAmount(*values.FeeAmount)
	}
	if values.FeeUsdAmount != nil {
		p.MerchantPayoutValues.SetFeeUsdAmount(*values.FeeUsdAmount)
	}
	if values.FeeUsdRate != nil {
		p.MerchantPayoutValues.SetFeeUsdRate(*values.FeeUsdRate)
	}
	if values.ChannelTrxID != nil {
		p.MerchantPayoutValues.SetChannelTrxID(*values.ChannelTrxID)
	}
	if values.ChannelAccount != nil {
		p.MerchantPayoutValues.SetChannelAccount(*values.ChannelAccount)
	}
	if values.ChannelGroup != nil {
		p.MerchantPayoutValues.SetChannelGroup(*values.ChannelGroup)
	}
	if values.ChannelFeeCcy != nil {
		p.MerchantPayoutValues.SetChannelFeeCcy(*values.ChannelFeeCcy)
	}
	if values.ChannelFeeAmount != nil {
		p.MerchantPayoutValues.SetChannelFeeAmount(*values.ChannelFeeAmount)
	}
	if values.ChannelFeeUsdAmount != nil {
		p.MerchantPayoutValues.SetChannelFeeUsdAmount(*values.ChannelFeeUsdAmount)
	}
	if values.ChannelFeeUsdRate != nil {
		p.MerchantPayoutValues.SetChannelFeeUsdRate(*values.ChannelFeeUsdRate)
	}
	if values.MetaData != nil {
		p.MerchantPayoutValues.SetMetaData(values.MetaData)
	}
	if values.RefundedCount != nil {
		p.MerchantPayoutValues.SetRefundedCount(*values.RefundedCount)
	}
	if values.RefundedAmount != nil {
		p.MerchantPayoutValues.SetRefundedAmount(*values.RefundedAmount)
	}
	if values.RefundedUsdAmount != nil {
		p.MerchantPayoutValues.SetRefundedUsdAmount(*values.RefundedUsdAmount)
	}
	if values.LastRefundedAt != nil {
		p.MerchantPayoutValues.SetLastRefundedAt(*values.LastRefundedAt)
	}
	if values.Detail != nil {
		p.MerchantPayoutValues.SetDetail(values.Detail)
	}
	if values.ConfirmedAt != nil {
		p.MerchantPayoutValues.SetConfirmedAt(*values.ConfirmedAt)
	}
	if values.CompletedAt != nil {
		p.MerchantPayoutValues.SetCompletedAt(*values.CompletedAt)
	}
	if values.CancelReason != nil {
		p.MerchantPayoutValues.SetCancelReason(*values.CancelReason)
	}
	if values.CancelFailedResult != nil {
		p.MerchantPayoutValues.SetCancelFailedResult(*values.CancelFailedResult)
	}
	if values.Version != nil {
		p.MerchantPayoutValues.SetVersion(*values.Version)
	}

	return p
}

func GetMerchantPayoutByReqID(mid, reqID string) *MerchantPayout {
	var existingTrx MerchantPayout
	if err := ReadDB.Where("mid = ? AND req_id = ?", mid, reqID).First(&existingTrx).Error; err == nil {
		return &existingTrx
	}
	return nil
}
func GetMerchantPayoutByTrxID(mid, trxID string) *Transaction {
	var trx Transaction
	if err := ReadDB.Model(&MerchantPayout{}).Where("trx_id = ? AND mid = ?", trxID, mid).First(&trx).Error; err == nil {
		return &trx
	}
	return nil
}

// ToTransaction converts Payout to Transaction
func (p *MerchantPayout) ToTransaction() *Transaction {
	if p == nil {
		return nil
	}

	transaction := &Transaction{
		ID:          int64(p.ID), // Convert uint64 to int64
		Mid:         p.Mid,
		TrxType:     protocol.TrxTypePayout,
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
			MetaData:            p.MerchantPayoutValues.MetaData,
			RefundedCount:       p.MerchantPayoutValues.RefundedCount,
			RefundedAmount:      p.MerchantPayoutValues.RefundedAmount,
			RefundedUsdAmount:   p.MerchantPayoutValues.RefundedUsdAmount,
			LastRefundedAt:      p.MerchantPayoutValues.LastRefundedAt,
			SettleStatus:        p.MerchantPayoutValues.SettleStatus,
			SettleID:            p.MerchantPayoutValues.SettleID,
			SettledAt:           p.MerchantPayoutValues.SettledAt,
			Country:             p.MerchantPayoutValues.Country,
			Remark:              p.MerchantPayoutValues.Remark,
			FlowNo:              p.MerchantPayoutValues.FlowNo,
			Status:              p.MerchantPayoutValues.Status,
			Reason:              p.MerchantPayoutValues.Reason,
			Link:                p.MerchantPayoutValues.Link,
			Detail:              p.MerchantPayoutValues.Detail,
			NotifyURL:           p.MerchantPayoutValues.NotifyURL,
			FeeCcy:              p.MerchantPayoutValues.FeeCcy,
			FeeAmount:           p.MerchantPayoutValues.FeeAmount,
			FeeUsdAmount:        p.MerchantPayoutValues.FeeUsdAmount,
			FeeUsdRate:          p.MerchantPayoutValues.FeeUsdRate,
			ChannelStatus:       p.MerchantPayoutValues.ChannelStatus,
			ResCode:             p.MerchantPayoutValues.ResCode,
			ResMsg:              p.MerchantPayoutValues.ResMsg,
			ChannelTrxID:        p.MerchantPayoutValues.ChannelTrxID,
			ChannelCode:         p.MerchantPayoutValues.ChannelCode,
			ChannelAccount:      p.MerchantPayoutValues.ChannelAccount,
			ChannelGroup:        p.MerchantPayoutValues.ChannelGroup,
			ChannelFeeCcy:       p.MerchantPayoutValues.ChannelFeeCcy,
			ChannelFeeAmount:    p.MerchantPayoutValues.ChannelFeeAmount,
			ChannelFeeUsdAmount: p.MerchantPayoutValues.ChannelFeeUsdAmount,
			ChannelFeeUsdRate:   p.MerchantPayoutValues.ChannelFeeUsdRate,
			ConfirmedAt:         p.MerchantPayoutValues.ConfirmedAt,
			CompletedAt:         p.MerchantPayoutValues.CompletedAt,
			ExpiredAt:           p.MerchantPayoutValues.ExpiredAt,
			CanceledAt:          p.MerchantPayoutValues.CanceledAt,
			CancelReason:        p.MerchantPayoutValues.CancelReason,
			CancelFailedResult:  p.MerchantPayoutValues.CancelFailedResult,
			Version:             p.MerchantPayoutValues.Version,
		},
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	return transaction
}
