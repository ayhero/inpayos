package protocol

import "github.com/shopspring/decimal"

// CcyCode 支持的币种
const (
	CcyCNY = "CNY" // 人民币
	CcyUSD = "USD" // 美元
	CcyEUR = "EUR" // 欧元
	CcyGBP = "GBP" // 英镑
	CcyJPY = "JPY" // 日元
)

// 账户相关请求/响应
type CreateAccountRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	UserType string `json:"user_type" binding:"required,oneof=merchant cashier bank"`
	Ccy      string `json:"ccy" binding:"required"`
}

type UpdateBalanceRequest struct {
	UserID        string          `json:"user_id" binding:"required"`
	UserType      string          `json:"user_type" binding:"required"`
	Ccy           string          `json:"ccy" binding:"required"`
	Operation     string          `json:"operation" binding:"required,oneof=add subtract freeze unfreeze margin release_margin"`
	Amount        decimal.Decimal `json:"amount" binding:"required"`
	TransactionID string          `json:"trx_id"`
	BillID        string          `json:"bill_id"`
	BusinessType  string          `json:"business_type"`
	Description   string          `json:"description"`
}

type Balance struct {
	Balance          string `json:"balance"`
	AvailableBalance string `json:"available_balance"`
	FrozenBalance    string `json:"frozen_balance"`
	MarginBalance    string `json:"margin_balance"`
	ReserveBalance   string `json:"reserve_balance"`
	Currency         string `json:"currency"`
	UpdatedAt        int64  `json:"updated_at"`
}

// 账户信息响应
type Account struct {
	AccountID    string   `json:"account_id"`
	UserID       string   `json:"user_id"`
	UserType     string   `json:"user_type"`
	Currency     string   `json:"currency"`
	Balance      *Balance `json:"balance"`
	Status       int      `json:"status"`
	Version      int64    `json:"version"`
	LastActiveAt int64    `json:"last_active_at"`
	CreatedAt    int64    `json:"created_at"`
	UpdatedAt    int64    `json:"updated_at"`
}

// 账户流水查询请求
type AccountFlowListRequest struct {
	FlowNo         string `json:"flow_no"`                      // 流水号
	TrxID          string `json:"trx_id"`                       // 关联业务ID
	TrxType        string `json:"trx_type"`                     // 业务类型
	Direction      string `json:"direction"`                    // 流水方向：1-进账 2-出账
	Ccy            string `json:"ccy"`                          // 币种
	CreatedAtStart int64  `json:"created_at_start"`             // 开始时间
	CreatedAtEnd   int64  `json:"created_at_end"`               // 结束时间
	Page           int    `json:"page" binding:"min=1"`         // 页码
	Size           int    `json:"size" binding:"min=1,max=100"` // 每页记录数
}

// 流水相关响应
type FundFlow struct {
	ID             uint64          `json:"id"`
	FlowNo         string          `json:"flow_no"`     // 流水号
	OriFlowNo      string          `json:"ori_flow_no"` // 原始流水号
	UserID         string          `json:"user_id"`
	UserType       string          `json:"user_type"`
	AccountID      string          `json:"account_id"`
	AccountVersion int64           `json:"account_version"` // 账户版本号
	Direction      string          `json:"direction"`       // 流水方向
	TrxID          string          `json:"trx_id"`          // 关联业务ID
	TrxType        string          `json:"trx_type"`        // 业务类型
	Ccy            string          `json:"ccy"`             // 币种
	Amount         decimal.Decimal `json:"amount"`          // 交易金额
	BeforeBalance  decimal.Decimal `json:"before_balance"`  // 变动前余额
	AfterBalance   decimal.Decimal `json:"after_balance"`   // 变动后余额
	Remark         string          `json:"remark"`          // 备注
	OperatorId     string          `json:"operator_id"`     // 操作人ID
	CreatedAt      int64           `json:"created_at"`
	UpdatedAt      int64           `json:"updated_at"`
}
