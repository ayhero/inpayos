package models

import (
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Transaction 通用交易记录表（作为所有业务交易的抽象层）
// 每个具体业务表（Payin, Payout等）通过 ToTransaction() 方法转换为此通用模型
type Transaction struct {
	ID        uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TeamID    string `json:"team_id" gorm:"column:team_id;type:varchar(32);index"`
	CashierID string `json:"cashier_id" gorm:"column:cashier_id;type:varchar(32);index"`
	Mid       string `json:"mid" gorm:"column:mid;type:varchar(32);index"`
	UserID    string `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	TrxID     string `json:"transaction_id" gorm:"column:transaction_id;type:varchar(64);uniqueIndex"`
	ReqID     string `json:"req_id" gorm:"column:req_id;type:varchar(64);index"`
	TrxType   string `json:"trx_type" gorm:"column:trx_type;type:varchar(16);index"` // receipt, payment, refund, transfer
	*TransactionValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type TransactionValues struct {
	Status        *string          `json:"status" gorm:"column:status;type:varchar(16);index"` // pending, processing, success, failed
	Amount        *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Fee           *decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(36,18);default:0"`
	Ccy           *string          `json:"ccy" gorm:"column:ccy;type:varchar(16)"`
	ChannelCode   *string          `json:"channel_code" gorm:"column:channel_code;type:varchar(32)"`
	PaymentMethod *string          `json:"payment_method" gorm:"column:payment_method;type:varchar(32)"`
	NotifyURL     *string          `json:"notify_url" gorm:"column:notify_url;type:varchar(512)"`
	ReturnURL     *string          `json:"return_url" gorm:"column:return_url;type:varchar(512)"`
	NotifyStatus  *string          `json:"notify_status" gorm:"column:notify_status;type:varchar(16);default:'pending'"`
	NotifyTimes   *int             `json:"notify_times" gorm:"column:notify_times;type:int;default:0"`
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

func (Transaction) TableName() string {
	return "t_transactions"
}

// TrxTypeTableMap 定义交易类型和对应的表名映射关系
var TrxTypeTableMap = map[string]string{
	protocol.TrxTypePayin:  "t_payin",
	protocol.TrxTypePayout: "t_payouts",
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

// GetAmount returns the Amount value
func (tv *TransactionValues) GetAmount() decimal.Decimal {
	if tv.Amount == nil {
		return decimal.Zero
	}
	return *tv.Amount
}

// GetFee returns the Fee value
func (tv *TransactionValues) GetFee() decimal.Decimal {
	if tv.Fee == nil {
		return decimal.Zero
	}
	return *tv.Fee
}

// GetCcy returns the Ccy value
func (tv *TransactionValues) GetCcy() string {
	if tv.Ccy == nil {
		return ""
	}
	return *tv.Ccy
}

// GetChannelCode returns the ChannelCode value
func (tv *TransactionValues) GetChannelCode() string {
	if tv.ChannelCode == nil {
		return ""
	}
	return *tv.ChannelCode
}

// GetPaymentMethod returns the PaymentMethod value
func (tv *TransactionValues) GetPaymentMethod() string {
	if tv.PaymentMethod == nil {
		return ""
	}
	return *tv.PaymentMethod
}

// GetNotifyURL returns the NotifyURL value
func (tv *TransactionValues) GetNotifyURL() string {
	if tv.NotifyURL == nil {
		return ""
	}
	return *tv.NotifyURL
}

// GetReturnURL returns the ReturnURL value
func (tv *TransactionValues) GetReturnURL() string {
	if tv.ReturnURL == nil {
		return ""
	}
	return *tv.ReturnURL
}

// GetNotifyStatus returns the NotifyStatus value
func (tv *TransactionValues) GetNotifyStatus() string {
	if tv.NotifyStatus == nil {
		return ""
	}
	return *tv.NotifyStatus
}

// GetNotifyTimes returns the NotifyTimes value
func (tv *TransactionValues) GetNotifyTimes() int {
	if tv.NotifyTimes == nil {
		return 0
	}
	return *tv.NotifyTimes
}

// GetOriTrxID returns the OriTrxID value
func (tv *TransactionValues) GetOriTrxID() string {
	if tv.OriTrxID == nil {
		return ""
	}
	return *tv.OriTrxID
}

// GetMetadata returns the Metadata value
func (tv *TransactionValues) GetMetadata() string {
	if tv.Metadata == nil {
		return ""
	}
	return *tv.Metadata
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

// GetUpdatedAt returns the UpdatedAt value
func (tv *TransactionValues) GetUpdatedAt() int64 {
	return tv.UpdatedAt
}

// SetStatus sets the Status value
func (tv *TransactionValues) SetStatus(value string) *TransactionValues {
	tv.Status = &value
	return tv
}

// SetAmount sets the Amount value
func (tv *TransactionValues) SetAmount(value decimal.Decimal) *TransactionValues {
	tv.Amount = &value
	return tv
}

// SetFee sets the Fee value
func (tv *TransactionValues) SetFee(value decimal.Decimal) *TransactionValues {
	tv.Fee = &value
	return tv
}

// SetCcy sets the Ccy value
func (tv *TransactionValues) SetCcy(value string) *TransactionValues {
	tv.Ccy = &value
	return tv
}

// SetChannelCode sets the ChannelCode value
func (tv *TransactionValues) SetChannelCode(value string) *TransactionValues {
	tv.ChannelCode = &value
	return tv
}

// SetPaymentMethod sets the PaymentMethod value
func (tv *TransactionValues) SetPaymentMethod(value string) *TransactionValues {
	tv.PaymentMethod = &value
	return tv
}

// SetNotifyURL sets the NotifyURL value
func (tv *TransactionValues) SetNotifyURL(value string) *TransactionValues {
	tv.NotifyURL = &value
	return tv
}

// SetReturnURL sets the ReturnURL value
func (tv *TransactionValues) SetReturnURL(value string) *TransactionValues {
	tv.ReturnURL = &value
	return tv
}

// SetNotifyStatus sets the NotifyStatus value
func (tv *TransactionValues) SetNotifyStatus(value string) *TransactionValues {
	tv.NotifyStatus = &value
	return tv
}

// SetNotifyTimes sets the NotifyTimes value
func (tv *TransactionValues) SetNotifyTimes(value int) *TransactionValues {
	tv.NotifyTimes = &value
	return tv
}

// SetOriTrxID sets the OriTrxID value
func (tv *TransactionValues) SetOriTrxID(value string) *TransactionValues {
	tv.OriTrxID = &value
	return tv
}

// SetMetadata sets the Metadata value
func (tv *TransactionValues) SetMetadata(value string) *TransactionValues {
	tv.Metadata = &value
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

// SetUpdatedAt sets the UpdatedAt value
func (tv *TransactionValues) SetUpdatedAt(value int64) *TransactionValues {
	tv.UpdatedAt = value
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
	if values.Amount != nil {
		t.TransactionValues.SetAmount(*values.Amount)
	}
	if values.Fee != nil {
		t.TransactionValues.SetFee(*values.Fee)
	}
	if values.Ccy != nil {
		t.TransactionValues.SetCcy(*values.Ccy)
	}
	if values.ChannelCode != nil {
		t.TransactionValues.SetChannelCode(*values.ChannelCode)
	}
	if values.PaymentMethod != nil {
		t.TransactionValues.SetPaymentMethod(*values.PaymentMethod)
	}
	if values.NotifyURL != nil {
		t.TransactionValues.SetNotifyURL(*values.NotifyURL)
	}
	if values.ReturnURL != nil {
		t.TransactionValues.SetReturnURL(*values.ReturnURL)
	}
	if values.NotifyStatus != nil {
		t.TransactionValues.SetNotifyStatus(*values.NotifyStatus)
	}
	if values.NotifyTimes != nil {
		t.TransactionValues.SetNotifyTimes(*values.NotifyTimes)
	}
	if values.OriTrxID != nil {
		t.TransactionValues.SetOriTrxID(*values.OriTrxID)
	}
	if values.Metadata != nil {
		t.TransactionValues.SetMetadata(*values.Metadata)
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
	if values.UsdAmount != nil {
		t.TransactionValues.SetUsdAmount(*values.UsdAmount)
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
	// UpdatedAt is not a pointer, so we always set it
	t.TransactionValues.SetUpdatedAt(values.UpdatedAt)

	return t
}

// GetUsdAmount 获取USD金额
func (tv *TransactionValues) GetUsdAmount() decimal.Decimal {
	if tv.UsdAmount == nil {
		return decimal.Zero
	}
	return *tv.UsdAmount
}

// SetUsdAmount 设置USD金额
func (tv *TransactionValues) SetUsdAmount(value decimal.Decimal) *TransactionValues {
	tv.UsdAmount = &value
	return tv
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
	db := query.BuildQuery(ReadDB.Model(&Transaction{}))
	err := db.Count(&count).Error
	return count, err
}

// ListTransactionByQuery 根据查询条件获取交易列表
func ListTransactionByQuery(query *TrxQuery, offset, limit int) ([]*Transaction, error) {
	var transactions []*Transaction
	db := query.BuildQuery(ReadDB)
	err := db.Offset(offset).Limit(limit).Find(&transactions).Error
	return transactions, err
}

// NewTrxValues 创建新的TransactionValues用于更新
func NewTrxValues() *TransactionValues {
	return &TransactionValues{}
}

// SaveTransactionValues 保存交易值更新
func SaveTransactionValues(db *gorm.DB, trx *Transaction, values *TransactionValues) error {
	// 创建更新映射
	updates := make(map[string]interface{})

	if values.Status != nil {
		updates["status"] = *values.Status
	}
	if values.Amount != nil {
		updates["amount"] = *values.Amount
	}
	if values.Fee != nil {
		updates["fee"] = *values.Fee
	}
	if values.Ccy != nil {
		updates["ccy"] = *values.Ccy
	}
	if values.ChannelCode != nil {
		updates["channel_code"] = *values.ChannelCode
	}
	if values.PaymentMethod != nil {
		updates["payment_method"] = *values.PaymentMethod
	}
	if values.NotifyURL != nil {
		updates["notify_url"] = *values.NotifyURL
	}
	if values.ReturnURL != nil {
		updates["return_url"] = *values.ReturnURL
	}
	if values.NotifyStatus != nil {
		updates["notify_status"] = *values.NotifyStatus
	}
	if values.NotifyTimes != nil {
		updates["notify_times"] = *values.NotifyTimes
	}
	if values.OriTrxID != nil {
		updates["ori_trx_id"] = *values.OriTrxID
	}
	if values.Metadata != nil {
		updates["metadata"] = *values.Metadata
	}
	if values.Remark != nil {
		updates["remark"] = *values.Remark
	}
	if values.ExpiredAt != nil {
		updates["expired_at"] = *values.ExpiredAt
	}
	if values.ConfirmedAt != nil {
		updates["confirmed_at"] = *values.ConfirmedAt
	}
	if values.CanceledAt != nil {
		updates["canceled_at"] = *values.CanceledAt
	}
	if values.UsdAmount != nil {
		updates["usd_amount"] = *values.UsdAmount
	}
	if values.SettleID != nil {
		updates["settle_id"] = *values.SettleID
	}
	if values.SettleStatus != nil {
		updates["settle_status"] = *values.SettleStatus
	}
	if values.SettledAt != nil {
		updates["settled_at"] = *values.SettledAt
	}
	// UpdatedAt总是更新
	updates["updated_at"] = values.UpdatedAt

	// 执行更新
	return db.Model(trx).Updates(updates).Error
}
