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

const (
	DirectionIn  = "in"  // 进账
	DirectionOut = "out" // 出账
)

var (
	AccountDirectionMap = map[string]string{
		TrxTypePayin:         DirectionIn,
		TrxTypePayout:        DirectionOut,
		TrxTypeRefund:        DirectionIn,
		TrxTypeDeposit:       DirectionIn,
		TrxTypeWithdraw:      DirectionOut,
		TrxTypeFreeze:        DirectionOut,
		TrxTypeUnfreeze:      DirectionIn,
		TrxTypeMarginDeposit: DirectionIn,
		TrxTypeMarginRelease: DirectionOut,
		TrxTypeFee:           DirectionOut,
		TrxTypeAdjustment:    DirectionIn,
		TrxTypeChargeback:    DirectionOut,
		TrxTypeSettle:        DirectionOut,
		TrxTypeTransfer:      DirectionOut,
		TrxTypeDividend:      DirectionIn,
		TrxTypeRfRecover:     DirectionOut,
		TrxTypeWdRecover:     DirectionIn,
	}
)

// 账户相关请求/响应
type CreateAccountRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	UserType string `json:"user_type" binding:"required,oneof=merchant cashier bank"`
	Ccy      string `json:"ccy" binding:"required"`
}

type UpdateBalanceRequest struct {
	UserID      string          `json:"user_id" binding:"required"`
	UserType    string          `json:"user_type" binding:"required"`
	Ccy         string          `json:"ccy" binding:"required"`
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	TrxID       string          `json:"trx_id"`
	TrxType     string          `json:"trx_type"`
	ReqID       string          `json:"req_id"`
	Description string          `json:"description"`
}

type Assert struct {
	Balance                string `json:"balance"`
	AvailableBalance       string `json:"available_balance"`
	FrozenBalance          string `json:"frozen_balance"`
	MarginBalance          string `json:"margin_balance"`
	AvailableMarginBalance string `json:"available_margin_balance"`
	FrozenMarginBalance    string `json:"frozen_margin_balance"`
	Ccy                    string `json:"ccy"`
	UpdatedAt              int64  `json:"updated_at"`
}

// 账户信息响应
type Account struct {
	AccountID    string  `json:"account_id"`
	UserID       string  `json:"user_id"`
	UserType     string  `json:"user_type"`
	Ccy          string  `json:"ccy"`
	Balance      *Assert `json:"balance"`
	Status       string  `json:"status"`
	Version      int64   `json:"version"`
	LastActiveAt int64   `json:"last_active_at"`
	CreatedAt    int64   `json:"created_at"`
	UpdatedAt    int64   `json:"updated_at"`
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
