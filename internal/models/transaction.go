package models

import (
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Transaction 通用交易记录表（作为所有业务交易的抽象层）
// 每个具体业务表（Payin, Payout等）通过 ToTransaction() 方法转换为此通用模型
type Transaction struct {
	ID                 int64            `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Tid                string           `json:"tid" gorm:"column:tid;type:varchar(32);index"`
	CashierID          string           `json:"cashier_id" gorm:"column:cashier_id;type:varchar(32);index"`
	Mid                string           `json:"mid" gorm:"column:mid;type:varchar(32);index"`
	UserID             string           `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	TrxID              string           `json:"trx_id" gorm:"column:trx_id;type:varchar(64);uniqueIndex"`
	TrxType            string           `json:"trx_type" gorm:"column:trx_type;type:varchar(32);index"` // 交易类型：payin, payout
	ReqID              string           `json:"req_id" gorm:"column:req_id;type:varchar(64);index"`
	OriTrxID           string           `json:"ori_trx_id" gorm:"column:ori_trx_id;index;<-:create"`
	OriReqID           string           `json:"ori_req_id" gorm:"column:ori_req_id;index;<-:create"`
	OriFlowNo          string           `json:"ori_flow_no" gorm:"column:ori_flow_no"`
	TrxMethod          string           `json:"trx_method" gorm:"column:trx_method;<-:create"`
	TrxMode            string           `json:"trx_mode" gorm:"column:trx_mode;<-:create"`
	TrxApp             string           `json:"trx_app" gorm:"column:trx_app;<-:create"`
	Pkg                string           `json:"pkg" gorm:"column:pkg;<-:create"`
	Did                string           `json:"did" gorm:"column:did;<-:create"`
	ProductID          string           `json:"product_id" gorm:"column:product_id;<-:create"`
	UserIP             string           `json:"user_ip" gorm:"column:user_ip;<-:create"`
	Email              string           `json:"email" gorm:"column:email;<-:create"`
	Phone              string           `json:"phone" gorm:"column:phone;<-:create"`
	Ccy                string           `json:"ccy" gorm:"column:ccy;<-:create"`
	Amount             *decimal.Decimal `json:"amount" gorm:"column:amount;<-:create"`
	UsdAmount          *decimal.Decimal `json:"usd_amount" gorm:"column:usd_amount;<-:create"`
	AccountNo          string           `json:"account_no" gorm:"column:account_no;<-:create"`
	AccountName        string           `json:"account_name" gorm:"column:account_name;<-:create"`
	AccountType        string           `json:"account_type" gorm:"column:account_type;<-:create"`
	BankCode           string           `json:"bank_code" gorm:"column:bank_code;<-:create"`
	BankName           string           `json:"bank_name" gorm:"column:bank_name;<-:create"`
	ReturnURL          string           `json:"return_url" gorm:"column:return_url;<-:create"`
	*TransactionValues `gorm:"embedded"`
	CreatedAt          int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt          int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

type TransactionValues struct {
	MetaData *protocol.MapData `json:"metadata" gorm:"column:metadata;serializer:json;type:json"`

	// Refund related fields
	RefundedCount     *int             `json:"refunded_count" gorm:"column:refunded_count"`
	RefundedAmount    *decimal.Decimal `json:"refunded_amount" gorm:"column:refunded_amount"`
	RefundedUsdAmount *decimal.Decimal `json:"refunded_usd_amount" gorm:"column:refunded_usd_amount"`
	LastRefundedAt    *int64           `json:"last_refunded_at" gorm:"column:last_refunded_at"`

	// Settlement related fields
	SettleStatus *string `json:"settle_status" gorm:"column:settle_status;index"` // SettleStatus 结算状态
	SettleID     *string `json:"settle_id" gorm:"column:settle_id;index"`
	SettledAt    *int64  `json:"settled_at" gorm:"column:settled_at"`

	// Basic transaction fields
	Country       *string        `json:"country" gorm:"column:country"`
	Remark        *string        `json:"remark" gorm:"column:remark"`
	FlowNo        *string        `json:"flow_no" gorm:"column:flow_no;index"`
	Status        *string        `json:"status" gorm:"column:status;index"`
	ChannelStatus *string        `json:"channel_status" gorm:"column:channel_status"`
	ResCode       *string        `json:"res_code" gorm:"column:res_code"`
	ResMsg        *string        `json:"res_msg" gorm:"column:res_msg"`
	Reason        *string        `json:"reason" gorm:"column:reason"`
	Link          *string        `json:"link" gorm:"column:link"`
	Detail        map[string]any `json:"detail" gorm:"column:detail;serializer:json;type:json"`
	NotifyURL     *string        `json:"notify_url" gorm:"column:notify_url"`

	// Fee related fields
	FeeCcy       *string          `json:"fee_ccy" gorm:"column:fee_ccy"`
	FeeAmount    *decimal.Decimal `json:"fee_amount" gorm:"column:fee_amount"`
	FeeUsdAmount *decimal.Decimal `json:"fee_usd_amount" gorm:"column:fee_usd_amount"`
	FeeUsdRate   *decimal.Decimal `json:"fee_usd_rate" gorm:"column:fee_usd_rate"`

	// Channel related fields
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

// TrxTypeTableMap 定义交易类型和对应的表名映射关系
var TrxTypeTableMap = map[string]string{
	protocol.TrxTypePayin:  "t_merchant_payins",
	protocol.TrxTypePayout: "t_merchant_payouts",
}

func GetTransactionQueryByType(trxType string) *gorm.DB {
	if _v, ok := TrxTypeTableMap[trxType]; ok {
		return ReadDB.Table(_v)
	}
	return ReadDB
}

// TrxQuery 交易查询参数
type TrxQuery struct {
	Mid            string `json:"mid"`             // 商户ID
	TrxType        string `json:"trx_type"`        // 交易类型
	TrxID          string `json:"trx_id"`          // 交易ID
	ReqID          string `json:"req_id"`          // 商户订单号
	TrxMethod      string `json:"trx_method"`      // 交易方式
	TrxMode        string `json:"trx_mode"`        // 交易模式
	Status         string `json:"status"`          // 交易状态
	FlowNo         string `json:"flow_no"`         // 流水号
	ChannelCode    string `json:"channel_code"`    // 渠道代码
	ChannelAccount string `json:"channel_account"` // 渠道账号
	ChannelGroup   string `json:"channel_group"`   // 渠道组
	ChannelTrxID   string `json:"channel_trx_id"`  // 渠道交易ID

	StatusList         []string `json:"status_list"`          // 交易状态列表
	MidList            []string `json:"mid_list"`             // 商户ID列表
	TrxIDList          []string `json:"trx_id_list"`          // 交易ID列表
	ReqIDList          []string `json:"req_id_list"`          // 商户订单号列表
	FlowNoList         []string `json:"flow_no_list"`         // 流水号列表
	ChannelTrxIDList   []string `json:"channel_trx_id_list"`  // 渠道交易ID列表
	ChannelAccountList []string `json:"channel_account_list"` // 渠道账号列表
	ChannelGroupList   []string `json:"channel_group_list"`   // 渠道组列表
	ChannelCodeList    []string `json:"channel_code_list"`    // 渠道代码列表
	TrxMethodList      []string `json:"trx_method_list"`      // 交易方式列表
	TrxModeList        []string `json:"trx_mode_list"`        // 交易模式列表

	SettleStatus     string `json:"settle_status"`      // 结算状态
	SettleStatusList []int  `json:"settle_status_list"` // 结算状态列表
	SettledAtStart   int64  `json:"settled_at_start"`   // 结算开始时间
	SettledAtEnd     int64  `json:"settled_at_end"`     // 结算结束时间

	CompletedAtStart int64 `json:"completed_at_start"` // 交易完成开始时间
	CompletedAtEnd   int64 `json:"completed_at_end"`   // 交易完成结束时间

	CreatedAtStart int64 `json:"created_at_start"` // 开始时间
	CreatedAtEnd   int64 `json:"created_at_end"`   // 结束时间
	Page           int   `json:"page"`             // 页码
	Size           int   `json:"size"`             // 每页记录数
}

// GetOffset 获取数据库查询的偏移量
func (q *TrxQuery) GetOffset() int {
	return (q.Page - 1) * q.Size
}

// GetLimit 获取数据库查询的限制数
func (q *TrxQuery) GetLimit() int {
	return q.Size
}

// BuildQuery 构建查询条件
func (q *TrxQuery) BuildQuery(db *gorm.DB) *gorm.DB {
	db = db.Where("mid = ?", q.Mid)
	if q.TrxType != "" {
		db = db.Where("trx_type = ?", q.TrxType)
	}
	if q.CreatedAtStart > 0 {
		db = db.Where("created_at >= ?", q.CreatedAtStart)
	}
	if q.CreatedAtEnd > 0 {
		db = db.Where("created_at <= ?", q.CreatedAtEnd)
	}
	if q.TrxID != "" {
		db = db.Where("trx_id = ?", q.TrxID)
	}
	if q.ReqID != "" {
		db = db.Where("req_id = ?", q.ReqID)
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}
	if len(q.StatusList) > 0 {
		db = db.Where("status IN ?", q.StatusList)
	}
	if q.SettleStatus != "" {
		db = db.Where("settle_status = ?", q.SettleStatus)
	}
	if len(q.SettleStatusList) > 0 {
		db = db.Where("settle_status IN ?", q.SettleStatusList)
	}
	if q.SettledAtStart > 0 {
		db = db.Where("settled_at >= ?", q.SettledAtStart)
	}
	if q.SettledAtEnd > 0 {
		db = db.Where("settled_at <= ?", q.SettledAtEnd)
	}
	if q.CompletedAtStart > 0 {
		db = db.Where("completed_at >= ?", q.CompletedAtStart)
	}
	if q.CompletedAtEnd > 0 {
		db = db.Where("completed_at <= ?", q.CompletedAtEnd)
	}
	return db
}

// GetStatus returns the Status value
func (tv *TransactionValues) GetStatus() string {
	if tv.Status == nil {
		return ""
	}
	return *tv.Status
}

// GetChannelCode returns the ChannelCode value
func (tv *TransactionValues) GetChannelCode() string {
	if tv.ChannelCode == nil {
		return ""
	}
	return *tv.ChannelCode
}

// GetNotifyURL returns the NotifyURL value
func (tv *TransactionValues) GetNotifyURL() string {
	if tv.NotifyURL == nil {
		return ""
	}
	return *tv.NotifyURL
}

// GetRemark returns the Remark value
func (tv *TransactionValues) GetRemark() string {
	if tv.Remark == nil {
		return ""
	}
	return *tv.Remark
}

// GetExpiredAt returns the ExpiredAt value
func (tv *TransactionValues) GetExpiredAt() int64 {
	if tv.ExpiredAt == nil {
		return 0
	}
	return *tv.ExpiredAt
}

// GetConfirmedAt returns the ConfirmedAt value
func (tv *TransactionValues) GetConfirmedAt() int64 {
	if tv.ConfirmedAt == nil {
		return 0
	}
	return *tv.ConfirmedAt
}

// GetCanceledAt returns the CanceledAt value
func (tv *TransactionValues) GetCanceledAt() int64 {
	if tv.CanceledAt == nil {
		return 0
	}
	return *tv.CanceledAt
}

// GetCompletedAt returns the CompletedAt value
func (tv *TransactionValues) GetCompletedAt() int64 {
	if tv.CompletedAt == nil {
		return 0
	}
	return *tv.CompletedAt
}

// GetRefundedCount returns the RefundedCount value
func (tv *TransactionValues) GetRefundedCount() int {
	if tv.RefundedCount == nil {
		return 0
	}
	return *tv.RefundedCount
}

// GetRefundedAmount returns the RefundedAmount value
func (tv *TransactionValues) GetRefundedAmount() decimal.Decimal {
	if tv.RefundedAmount == nil {
		return decimal.Zero
	}
	return *tv.RefundedAmount
}

// GetRefundedUsdAmount returns the RefundedUsdAmount value
func (tv *TransactionValues) GetRefundedUsdAmount() decimal.Decimal {
	if tv.RefundedUsdAmount == nil {
		return decimal.Zero
	}
	return *tv.RefundedUsdAmount
}

// GetLastRefundedAt returns the LastRefundedAt value
func (tv *TransactionValues) GetLastRefundedAt() int64 {
	if tv.LastRefundedAt == nil {
		return 0
	}
	return *tv.LastRefundedAt
}

// GetCountry returns the Country value
func (tv *TransactionValues) GetCountry() string {
	if tv.Country == nil {
		return ""
	}
	return *tv.Country
}

// GetFlowNo returns the FlowNo value
func (tv *TransactionValues) GetFlowNo() string {
	if tv.FlowNo == nil {
		return ""
	}
	return *tv.FlowNo
}

// GetReason returns the Reason value
func (tv *TransactionValues) GetReason() string {
	if tv.Reason == nil {
		return ""
	}
	return *tv.Reason
}

// GetLink returns the Link value
func (tv *TransactionValues) GetLink() string {
	if tv.Link == nil {
		return ""
	}
	return *tv.Link
}

// GetFeeCcy returns the FeeCcy value
func (tv *TransactionValues) GetFeeCcy() string {
	if tv.FeeCcy == nil {
		return ""
	}
	return *tv.FeeCcy
}

// GetFeeAmount returns the FeeAmount value
func (tv *TransactionValues) GetFeeAmount() decimal.Decimal {
	if tv.FeeAmount == nil {
		return decimal.Zero
	}
	return *tv.FeeAmount
}

// GetFeeUsdAmount returns the FeeUsdAmount value
func (tv *TransactionValues) GetFeeUsdAmount() decimal.Decimal {
	if tv.FeeUsdAmount == nil {
		return decimal.Zero
	}
	return *tv.FeeUsdAmount
}

// GetFeeUsdRate returns the FeeUsdRate value
func (tv *TransactionValues) GetFeeUsdRate() decimal.Decimal {
	if tv.FeeUsdRate == nil {
		return decimal.Zero
	}
	return *tv.FeeUsdRate
}

// GetChannelStatus returns the ChannelStatus value
func (tv *TransactionValues) GetChannelStatus() string {
	if tv.ChannelStatus == nil {
		return ""
	}
	return *tv.ChannelStatus
}

// GetResCode returns the ResCode value
func (tv *TransactionValues) GetResCode() string {
	if tv.ResCode == nil {
		return ""
	}
	return *tv.ResCode
}

// GetResMsg returns the ResMsg value
func (tv *TransactionValues) GetResMsg() string {
	if tv.ResMsg == nil {
		return ""
	}
	return *tv.ResMsg
}

// GetChannelTrxID returns the ChannelTrxID value
func (tv *TransactionValues) GetChannelTrxID() string {
	if tv.ChannelTrxID == nil {
		return ""
	}
	return *tv.ChannelTrxID
}

// GetChannelAccount returns the ChannelAccount value
func (tv *TransactionValues) GetChannelAccount() string {
	if tv.ChannelAccount == nil {
		return ""
	}
	return *tv.ChannelAccount
}

// GetChannelGroup returns the ChannelGroup value
func (tv *TransactionValues) GetChannelGroup() string {
	if tv.ChannelGroup == nil {
		return ""
	}
	return *tv.ChannelGroup
}

// GetChannelFeeCcy returns the ChannelFeeCcy value
func (tv *TransactionValues) GetChannelFeeCcy() string {
	if tv.ChannelFeeCcy == nil {
		return ""
	}
	return *tv.ChannelFeeCcy
}

// GetChannelFeeAmount returns the ChannelFeeAmount value
func (tv *TransactionValues) GetChannelFeeAmount() decimal.Decimal {
	if tv.ChannelFeeAmount == nil {
		return decimal.Zero
	}
	return *tv.ChannelFeeAmount
}

// GetChannelFeeUsdAmount returns the ChannelFeeUsdAmount value
func (tv *TransactionValues) GetChannelFeeUsdAmount() decimal.Decimal {
	if tv.ChannelFeeUsdAmount == nil {
		return decimal.Zero
	}
	return *tv.ChannelFeeUsdAmount
}

// GetChannelFeeUsdRate returns the ChannelFeeUsdRate value
func (tv *TransactionValues) GetChannelFeeUsdRate() decimal.Decimal {
	if tv.ChannelFeeUsdRate == nil {
		return decimal.Zero
	}
	return *tv.ChannelFeeUsdRate
}

// GetMetaData returns the MetaData value
func (tv *TransactionValues) GetMetaData() *protocol.MapData {
	return tv.MetaData
}

// GetDetail returns the Detail value
func (tv *TransactionValues) GetDetail() map[string]any {
	return tv.Detail
}

// GetCancelReason returns the CancelReason value
func (tv *TransactionValues) GetCancelReason() string {
	if tv.CancelReason == nil {
		return ""
	}
	return *tv.CancelReason
}

// GetCancelFailedResult returns the CancelFailedResult value
func (tv *TransactionValues) GetCancelFailedResult() string {
	if tv.CancelFailedResult == nil {
		return ""
	}
	return *tv.CancelFailedResult
}

// GetVersion returns the Version value
func (tv *TransactionValues) GetVersion() int64 {
	if tv.Version == nil {
		return 0
	}
	return *tv.Version
}

// SetStatus sets the Status value
func (tv *TransactionValues) SetStatus(value string) *TransactionValues {
	tv.Status = &value
	return tv
}

// SetChannelCode sets the ChannelCode value
func (tv *TransactionValues) SetChannelCode(value string) *TransactionValues {
	tv.ChannelCode = &value
	return tv
}

// SetNotifyURL sets the NotifyURL value
func (tv *TransactionValues) SetNotifyURL(value string) *TransactionValues {
	tv.NotifyURL = &value
	return tv
}

// SetRemark sets the Remark value
func (tv *TransactionValues) SetRemark(value string) *TransactionValues {
	tv.Remark = &value
	return tv
}

// SetExpiredAt sets the ExpiredAt value
func (tv *TransactionValues) SetExpiredAt(value int64) *TransactionValues {
	tv.ExpiredAt = &value
	return tv
}

// SetConfirmedAt sets the ConfirmedAt value
func (tv *TransactionValues) SetConfirmedAt(value int64) *TransactionValues {
	tv.ConfirmedAt = &value
	return tv
}

// SetCanceledAt sets the CanceledAt value
func (tv *TransactionValues) SetCanceledAt(value int64) *TransactionValues {
	tv.CanceledAt = &value
	return tv
}

// SetCompletedAt sets the CompletedAt value
func (tv *TransactionValues) SetCompletedAt(value int64) *TransactionValues {
	tv.CompletedAt = &value
	return tv
}

// SetRefundedCount sets the RefundedCount value
func (tv *TransactionValues) SetRefundedCount(value int) *TransactionValues {
	tv.RefundedCount = &value
	return tv
}

// SetRefundedAmount sets the RefundedAmount value
func (tv *TransactionValues) SetRefundedAmount(value decimal.Decimal) *TransactionValues {
	tv.RefundedAmount = &value
	return tv
}

// SetRefundedUsdAmount sets the RefundedUsdAmount value
func (tv *TransactionValues) SetRefundedUsdAmount(value decimal.Decimal) *TransactionValues {
	tv.RefundedUsdAmount = &value
	return tv
}

// SetLastRefundedAt sets the LastRefundedAt value
func (tv *TransactionValues) SetLastRefundedAt(value int64) *TransactionValues {
	tv.LastRefundedAt = &value
	return tv
}

// SetCountry sets the Country value
func (tv *TransactionValues) SetCountry(value string) *TransactionValues {
	tv.Country = &value
	return tv
}

// SetFlowNo sets the FlowNo value
func (tv *TransactionValues) SetFlowNo(value string) *TransactionValues {
	tv.FlowNo = &value
	return tv
}

// SetReason sets the Reason value
func (tv *TransactionValues) SetReason(value string) *TransactionValues {
	tv.Reason = &value
	return tv
}

// SetLink sets the Link value
func (tv *TransactionValues) SetLink(value string) *TransactionValues {
	tv.Link = &value
	return tv
}

// SetFeeCcy sets the FeeCcy value
func (tv *TransactionValues) SetFeeCcy(value string) *TransactionValues {
	tv.FeeCcy = &value
	return tv
}

// SetFeeAmount sets the FeeAmount value
func (tv *TransactionValues) SetFeeAmount(value decimal.Decimal) *TransactionValues {
	tv.FeeAmount = &value
	return tv
}

// SetFeeUsdAmount sets the FeeUsdAmount value
func (tv *TransactionValues) SetFeeUsdAmount(value decimal.Decimal) *TransactionValues {
	tv.FeeUsdAmount = &value
	return tv
}

// SetFeeUsdRate sets the FeeUsdRate value
func (tv *TransactionValues) SetFeeUsdRate(value decimal.Decimal) *TransactionValues {
	tv.FeeUsdRate = &value
	return tv
}

// SetChannelStatus sets the ChannelStatus value
func (tv *TransactionValues) SetChannelStatus(value string) *TransactionValues {
	tv.ChannelStatus = &value
	return tv
}

// SetResCode sets the ResCode value
func (tv *TransactionValues) SetResCode(value string) *TransactionValues {
	tv.ResCode = &value
	return tv
}

// SetResMsg sets the ResMsg value
func (tv *TransactionValues) SetResMsg(value string) *TransactionValues {
	tv.ResMsg = &value
	return tv
}

// SetChannelTrxID sets the ChannelTrxID value
func (tv *TransactionValues) SetChannelTrxID(value string) *TransactionValues {
	tv.ChannelTrxID = &value
	return tv
}

// SetChannelAccount sets the ChannelAccount value
func (tv *TransactionValues) SetChannelAccount(value string) *TransactionValues {
	tv.ChannelAccount = &value
	return tv
}

// SetChannelGroup sets the ChannelGroup value
func (tv *TransactionValues) SetChannelGroup(value string) *TransactionValues {
	tv.ChannelGroup = &value
	return tv
}

// SetChannelFeeCcy sets the ChannelFeeCcy value
func (tv *TransactionValues) SetChannelFeeCcy(value string) *TransactionValues {
	tv.ChannelFeeCcy = &value
	return tv
}

// SetChannelFeeAmount sets the ChannelFeeAmount value
func (tv *TransactionValues) SetChannelFeeAmount(value decimal.Decimal) *TransactionValues {
	tv.ChannelFeeAmount = &value
	return tv
}

// SetChannelFeeUsdAmount sets the ChannelFeeUsdAmount value
func (tv *TransactionValues) SetChannelFeeUsdAmount(value decimal.Decimal) *TransactionValues {
	tv.ChannelFeeUsdAmount = &value
	return tv
}

// SetChannelFeeUsdRate sets the ChannelFeeUsdRate value
func (tv *TransactionValues) SetChannelFeeUsdRate(value decimal.Decimal) *TransactionValues {
	tv.ChannelFeeUsdRate = &value
	return tv
}

// SetMetaData sets the MetaData value
func (tv *TransactionValues) SetMetaData(value *protocol.MapData) *TransactionValues {
	tv.MetaData = value
	return tv
}

// SetDetail sets the Detail value
func (tv *TransactionValues) SetDetail(value map[string]any) *TransactionValues {
	tv.Detail = value
	return tv
}

// SetCancelReason sets the CancelReason value
func (tv *TransactionValues) SetCancelReason(value string) *TransactionValues {
	tv.CancelReason = &value
	return tv
}

// SetCancelFailedResult sets the CancelFailedResult value
func (tv *TransactionValues) SetCancelFailedResult(value string) *TransactionValues {
	tv.CancelFailedResult = &value
	return tv
}

// SetVersion sets the Version value
func (tv *TransactionValues) SetVersion(value int64) *TransactionValues {
	tv.Version = &value
	return tv
}

// SetValues sets multiple TransactionValues fields at once
func (t *Transaction) SetValues(values *TransactionValues) *Transaction {
	if values == nil {
		return t
	}

	if t.TransactionValues == nil {
		t.TransactionValues = &TransactionValues{}
	}

	if values.Status != nil {
		t.TransactionValues.SetStatus(*values.Status)
	}
	if values.ChannelCode != nil {
		t.TransactionValues.SetChannelCode(*values.ChannelCode)
	}
	if values.NotifyURL != nil {
		t.TransactionValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.Remark != nil {
		t.TransactionValues.SetRemark(*values.Remark)
	}
	if values.ExpiredAt != nil {
		t.TransactionValues.SetExpiredAt(*values.ExpiredAt)
	}
	if values.ConfirmedAt != nil {
		t.TransactionValues.SetConfirmedAt(*values.ConfirmedAt)
	}
	if values.CanceledAt != nil {
		t.TransactionValues.SetCanceledAt(*values.CanceledAt)
	}

	if values.SettleID != nil {
		t.TransactionValues.SetSettleID(*values.SettleID)
	}
	if values.SettleStatus != nil {
		t.TransactionValues.SetSettleStatus(*values.SettleStatus)
	}
	if values.SettledAt != nil {
		t.TransactionValues.SetSettledAt(*values.SettledAt)
	}
	if values.CompletedAt != nil {
		t.TransactionValues.SetCompletedAt(*values.CompletedAt)
	}
	if values.RefundedCount != nil {
		t.TransactionValues.SetRefundedCount(*values.RefundedCount)
	}
	if values.RefundedAmount != nil {
		t.TransactionValues.SetRefundedAmount(*values.RefundedAmount)
	}
	if values.RefundedUsdAmount != nil {
		t.TransactionValues.SetRefundedUsdAmount(*values.RefundedUsdAmount)
	}
	if values.LastRefundedAt != nil {
		t.TransactionValues.SetLastRefundedAt(*values.LastRefundedAt)
	}
	if values.Country != nil {
		t.TransactionValues.SetCountry(*values.Country)
	}
	if values.FlowNo != nil {
		t.TransactionValues.SetFlowNo(*values.FlowNo)
	}
	if values.Reason != nil {
		t.TransactionValues.SetReason(*values.Reason)
	}
	if values.Link != nil {
		t.TransactionValues.SetLink(*values.Link)
	}
	if values.FeeCcy != nil {
		t.TransactionValues.SetFeeCcy(*values.FeeCcy)
	}
	if values.FeeAmount != nil {
		t.TransactionValues.SetFeeAmount(*values.FeeAmount)
	}
	if values.FeeUsdAmount != nil {
		t.TransactionValues.SetFeeUsdAmount(*values.FeeUsdAmount)
	}
	if values.FeeUsdRate != nil {
		t.TransactionValues.SetFeeUsdRate(*values.FeeUsdRate)
	}
	if values.ChannelStatus != nil {
		t.TransactionValues.SetChannelStatus(*values.ChannelStatus)
	}
	if values.ResCode != nil {
		t.TransactionValues.SetResCode(*values.ResCode)
	}
	if values.ResMsg != nil {
		t.TransactionValues.SetResMsg(*values.ResMsg)
	}
	if values.ChannelTrxID != nil {
		t.TransactionValues.SetChannelTrxID(*values.ChannelTrxID)
	}
	if values.ChannelAccount != nil {
		t.TransactionValues.SetChannelAccount(*values.ChannelAccount)
	}
	if values.ChannelGroup != nil {
		t.TransactionValues.SetChannelGroup(*values.ChannelGroup)
	}
	if values.ChannelFeeCcy != nil {
		t.TransactionValues.SetChannelFeeCcy(*values.ChannelFeeCcy)
	}
	if values.ChannelFeeAmount != nil {
		t.TransactionValues.SetChannelFeeAmount(*values.ChannelFeeAmount)
	}
	if values.ChannelFeeUsdAmount != nil {
		t.TransactionValues.SetChannelFeeUsdAmount(*values.ChannelFeeUsdAmount)
	}
	if values.ChannelFeeUsdRate != nil {
		t.TransactionValues.SetChannelFeeUsdRate(*values.ChannelFeeUsdRate)
	}
	if values.MetaData != nil {
		t.TransactionValues.SetMetaData(values.MetaData)
	}
	if values.Detail != nil {
		t.TransactionValues.SetDetail(values.Detail)
	}
	if values.CancelReason != nil {
		t.TransactionValues.SetCancelReason(*values.CancelReason)
	}
	if values.CancelFailedResult != nil {
		t.TransactionValues.SetCancelFailedResult(*values.CancelFailedResult)
	}
	if values.Version != nil {
		t.TransactionValues.SetVersion(*values.Version)
	}
	return t
}

// GetSettleID 获取结算ID
func (tv *TransactionValues) GetSettleID() string {
	if tv.SettleID == nil {
		return ""
	}
	return *tv.SettleID
}

// SetSettleID 设置结算ID
func (tv *TransactionValues) SetSettleID(value string) *TransactionValues {
	tv.SettleID = &value
	return tv
}

// GetSettledAt 获取结算时间
func (tv *TransactionValues) GetSettledAt() int64 {
	if tv.SettledAt == nil {
		return 0
	}
	return *tv.SettledAt
}

// SetSettledAt 设置结算时间
func (tv *TransactionValues) SetSettledAt(value int64) *TransactionValues {
	tv.SettledAt = &value
	return tv
}

// GetSettleStatus 获取结算状态
func (tv *TransactionValues) GetSettleStatus() string {
	if tv.SettleStatus == nil {
		return ""
	}
	return *tv.SettleStatus
}

// SetSettleStatus 设置结算状态
func (tv *TransactionValues) SetSettleStatus(value string) *TransactionValues {
	tv.SettleStatus = &value
	return tv
}

// CountTransactionByQuery 根据查询条件统计交易数量
func CountTransactionByQuery(query *TrxQuery) (int64, error) {
	var count int64
	db := query.BuildQuery(GetTransactionQueryByType(query.TrxType))
	err := db.Count(&count).Error
	return count, err
}

// ListTransactionByQuery 根据查询条件获取交易列表
func ListTransactionByQuery(query *TrxQuery, offset, limit int) ([]*Transaction, error) {
	var transactions []*Transaction
	db := query.BuildQuery(GetTransactionQueryByType(query.TrxType))
	err := db.Offset(offset).Limit(limit).Find(&transactions).Error
	return transactions, err
}

// NewTrxValues 创建新的TransactionValues用于更新
func NewTrxValues() *TransactionValues {
	return &TransactionValues{}
}

// SaveTransactionValues 保存交易值更新
func SaveTransactionValues(db *gorm.DB, trx *Transaction, values *TransactionValues) (err error) {
	defer func() {
		if err == nil {
			trx.SetValues(values)
		}
	}()
	// 执行更新
	err = GetTransactionQueryByType(trx.TrxType).Where("trx_id=?", trx.TrxID).UpdateColumns(values).Error
	return
}

func (t *Transaction) Protocol() *protocol.Transaction {
	if t == nil {
		return nil
	}

	info := &protocol.Transaction{
		// 基础交易信息
		ID:        t.ID,
		Tid:       t.Tid,
		TrxID:     t.TrxID,
		TrxType:   t.TrxType,
		Mid:       t.Mid,
		ReqID:     t.ReqID,
		UserID:    t.UserID,
		CashierID: t.CashierID,

		// 原始交易信息
		OriTrxID:  t.OriTrxID,
		OriReqID:  t.OriReqID,
		OriFlowNo: t.OriFlowNo,

		// 交易方式和模式
		TrxMethod: t.TrxMethod,
		TrxMode:   t.TrxMode,
		TrxApp:    t.TrxApp,
		Pkg:       t.Pkg,
		Did:       t.Did,
		ProductID: t.ProductID,

		// 用户信息
		UserIP: t.UserIP,
		Email:  t.Email,
		Phone:  t.Phone,

		// 币种和账户信息
		Ccy:         t.Ccy,
		AccountNo:   t.AccountNo,
		AccountName: t.AccountName,
		AccountType: t.AccountType,
		BankCode:    t.BankCode,
		BankName:    t.BankName,

		// URL信息
		ReturnURL: t.ReturnURL,

		// 时间戳
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}

	// 处理金额字段 - 转换为字符串
	if t.Amount != nil {
		info.Amount = t.Amount.String()

		// 处理实际金额（从Amount减去费用后的金额）
		actualAmount := *t.Amount
		if t.TransactionValues != nil && t.TransactionValues.FeeAmount != nil {
			actualAmount = actualAmount.Sub(*t.TransactionValues.FeeAmount)
		}
		info.ActualAmount = actualAmount.String()
	}

	if t.UsdAmount != nil {
		info.UsdAmount = t.UsdAmount.String()
	}

	// 安全处理 TransactionValues 字段
	if t.TransactionValues != nil {
		// 基础状态信息
		info.Country = t.TransactionValues.GetCountry()
		info.Status = t.TransactionValues.GetStatus()
		info.ChannelStatus = t.TransactionValues.GetChannelStatus()
		info.ResCode = t.TransactionValues.GetResCode()
		info.ResMsg = t.TransactionValues.GetResMsg()
		info.Reason = t.TransactionValues.GetReason()
		info.NotifyURL = t.TransactionValues.GetNotifyURL()
		info.Remark = t.TransactionValues.GetRemark()

		// 流程信息
		info.FlowNo = t.TransactionValues.GetFlowNo()
		info.Link = t.TransactionValues.GetLink()

		// 费用信息
		info.FeeCcy = t.TransactionValues.GetFeeCcy()
		if t.TransactionValues.FeeAmount != nil {
			info.FeeAmount = t.TransactionValues.FeeAmount.String()
		}
		if t.TransactionValues.FeeUsdAmount != nil {
			info.FeeUsdAmount = t.TransactionValues.FeeUsdAmount.String()
		}
		if t.TransactionValues.FeeUsdRate != nil {
			info.FeeUsdRate = t.TransactionValues.FeeUsdRate.String()
		}

		// 渠道信息
		info.ChannelTrxID = t.TransactionValues.GetChannelTrxID()
		info.ChannelCode = t.TransactionValues.GetChannelCode()
		info.ChannelAccount = t.TransactionValues.GetChannelAccount()
		info.ChannelGroup = t.TransactionValues.GetChannelGroup()
		info.ChannelFeeCcy = t.TransactionValues.GetChannelFeeCcy()
		if t.TransactionValues.ChannelFeeAmount != nil {
			info.ChannelFeeAmount = t.TransactionValues.ChannelFeeAmount.String()
		}
		if t.TransactionValues.ChannelFeeUsdAmount != nil {
			info.ChannelFeeUsdAmount = t.TransactionValues.ChannelFeeUsdAmount.String()
		}
		if t.TransactionValues.ChannelFeeUsdRate != nil {
			info.ChannelFeeUsdRate = t.TransactionValues.ChannelFeeUsdRate.String()
		}

		// 退款信息
		info.RefundedCount = t.TransactionValues.GetRefundedCount()
		if t.TransactionValues.RefundedAmount != nil {
			info.RefundedAmount = t.TransactionValues.RefundedAmount.String()
		}
		if t.TransactionValues.RefundedUsdAmount != nil {
			info.RefundedUsdAmount = t.TransactionValues.RefundedUsdAmount.String()
		}
		info.LastRefundedAt = t.TransactionValues.GetLastRefundedAt()

		// 结算信息
		info.SettleStatus = t.TransactionValues.GetSettleStatus()
		info.SettleID = t.TransactionValues.GetSettleID()
		info.SettledAt = t.TransactionValues.GetSettledAt()

		// 时间字段
		info.ConfirmedAt = t.TransactionValues.GetConfirmedAt()
		info.CompletedAt = t.TransactionValues.GetCompletedAt()
		info.ExpiredAt = t.TransactionValues.GetExpiredAt()
		info.CanceledAt = t.TransactionValues.GetCanceledAt()
		info.CancelReason = t.TransactionValues.GetCancelReason()
		info.CancelFailedResult = t.TransactionValues.GetCancelFailedResult()

		// 扩展信息
		if t.TransactionValues.GetMetaData() != nil {
			info.Metadata = *t.TransactionValues.GetMetaData()
		}
		info.Detail = t.TransactionValues.GetDetail()
		info.Version = t.TransactionValues.GetVersion()

		// 处理失败原因
		if t.TransactionValues.GetStatus() == protocol.StatusFailed {
			info.FailureReason = t.TransactionValues.GetReason()
		}
	}

	return info
}
