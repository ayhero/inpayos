package protocol

// Transaction 交易响应
type Transaction struct {
	// 基础交易信息
	ID      int64  `json:"id,omitempty"`
	Tid     string `json:"tid,omitempty"`
	TrxID   string `json:"trx_id,omitempty"`
	TrxType string `json:"trx_type,omitempty"`
	Mid     string `json:"mid,omitempty"`
	ReqID   string `json:"req_id,omitempty"`
	UserID  string `json:"user_id,omitempty"`

	// 原始交易信息
	OriTrxID  string `json:"ori_trx_id,omitempty"`
	OriReqID  string `json:"ori_req_id,omitempty"`
	OriFlowNo string `json:"ori_flow_no,omitempty"`

	// 交易方式和模式
	TrxMethod string `json:"trx_method,omitempty"`
	TrxMode   string `json:"trx_mode,omitempty"`
	TrxApp    string `json:"trx_app,omitempty"`
	Pkg       string `json:"pkg,omitempty"`
	Did       string `json:"did,omitempty"`
	ProductID string `json:"product_id,omitempty"`

	// 用户信息
	UserIP  string `json:"user_ip,omitempty"`
	Email   string `json:"email,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Country string `json:"country,omitempty"`

	// 金额信息
	Ccy          string `json:"ccy,omitempty"`
	Amount       string `json:"amount,omitempty"`
	ActualAmount string `json:"actual_amount,omitempty"`
	UsdAmount    string `json:"usd_amount,omitempty"`

	// 费用信息
	FeeCcy       string `json:"fee_ccy,omitempty"`
	FeeAmount    string `json:"fee_amount,omitempty"`
	FeeUsdAmount string `json:"fee_usd_amount,omitempty"`
	FeeUsdRate   string `json:"fee_usd_rate,omitempty"`

	// 账户信息
	AccountNo   string `json:"account_no,omitempty"`
	AccountName string `json:"account_name,omitempty"`
	AccountType string `json:"account_type,omitempty"`
	BankCode    string `json:"bank_code,omitempty"`
	BankName    string `json:"bank_name,omitempty"`

	// 状态信息
	Status        string `json:"status,omitempty"`
	ChannelStatus string `json:"channel_status,omitempty"`
	ResCode       string `json:"res_code,omitempty"`
	ResMsg        string `json:"res_msg,omitempty"`
	Reason        string `json:"reason,omitempty"`
	FailureReason string `json:"failure_reason,omitempty"`

	// 渠道信息
	ChannelTrxID        string `json:"channel_trx_id,omitempty"`
	ChannelCode         string `json:"channel_code,omitempty"`
	ChannelAccount      string `json:"channel_account,omitempty"`
	ChannelGroup        string `json:"channel_group,omitempty"`
	ChannelFeeCcy       string `json:"channel_fee_ccy,omitempty"`
	ChannelFeeAmount    string `json:"channel_fee_amount,omitempty"`
	ChannelFeeUsdAmount string `json:"channel_fee_usd_amount,omitempty"`
	ChannelFeeUsdRate   string `json:"channel_fee_usd_rate,omitempty"`

	// 流程信息
	FlowNo    string `json:"flow_no,omitempty"`
	Link      string `json:"link,omitempty"`
	NotifyURL string `json:"notify_url,omitempty"`
	ReturnURL string `json:"return_url,omitempty"`
	Remark    string `json:"remark,omitempty"`

	// 退款信息
	RefundedCount     int    `json:"refunded_count,omitempty"`
	RefundedAmount    string `json:"refunded_amount,omitempty"`
	RefundedUsdAmount string `json:"refunded_usd_amount,omitempty"`
	LastRefundedAt    int64  `json:"last_refunded_at,omitempty"`

	// 结算信息
	SettleStatus string `json:"settle_status,omitempty"`
	SettleID     string `json:"settle_id,omitempty"`
	SettledAt    int64  `json:"settled_at,omitempty"`

	// 取消信息
	CanceledAt         int64  `json:"canceled_at,omitempty"`
	CancelReason       string `json:"cancel_reason,omitempty"`
	CancelFailedResult string `json:"cancel_failed_result,omitempty"`

	// 时间戳
	ConfirmedAt int64 `json:"confirmed_at,omitempty"`
	CompletedAt int64 `json:"completed_at,omitempty"`
	ExpiredAt   int64 `json:"expired_at,omitempty"`
	CreatedAt   int64 `json:"created_at,omitempty"`
	UpdatedAt   int64 `json:"updated_at,omitempty"`

	// 扩展信息
	Metadata MapData        `json:"metadata,omitempty"`
	Detail   map[string]any `json:"detail,omitempty"`
	Version  int64          `json:"version,omitempty"`

	// 收银员信息
	CashierID string `json:"cashier_id,omitempty"`
}
