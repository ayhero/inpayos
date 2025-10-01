package protocol

import (
	"github.com/shopspring/decimal"
)

// TrxRouterRequest 交易路由请求
type TrxRouterRequest struct {
	RouterID string `json:"router_id" binding:"required"`
	Name     string `json:"name" binding:"required"`

	// 路由条件 - 三选一，互斥
	CashierTypes []string `json:"cashier_types,omitempty"` // Cashier类型数组
	CashierIDs   []string `json:"cashier_ids,omitempty"`   // 指定Cashier ID数组
	CashierGroup string   `json:"cashier_group,omitempty"` // Cashier分组

	// 匹配条件
	PaymentMethod   string           `json:"payment_method,omitempty"`   // 支付方式
	TransactionType string           `json:"transaction_type,omitempty"` // receipt, payment
	Currency        string           `json:"currency,omitempty"`         // 币种
	Country         string           `json:"country,omitempty"`          // 国家代码
	MinAmount       *decimal.Decimal `json:"min_amount,omitempty"`       // 最小金额
	MaxAmount       *decimal.Decimal `json:"max_amount,omitempty"`       // 最大金额
	Priority        int              `json:"priority"`                   // 优先级，默认100
	Status          string           `json:"status"`                     // 状态
	EffectiveTime   *int64           `json:"effective_time,omitempty"`   // 生效时间
	ExpireTime      *int64           `json:"expire_time,omitempty"`      // 过期时间
	Remark          string           `json:"remark,omitempty"`           // 备注
}

// TrxRouterResponse 交易路由响应
type TrxRouterResponse struct {
	ID              uint64           `json:"id"`
	RouterID        string           `json:"router_id"`
	Name            string           `json:"name"`
	CashierTypes    []string         `json:"cashier_types"`
	CashierIDs      []string         `json:"cashier_ids"`
	CashierGroup    string           `json:"cashier_group"`
	PaymentMethod   string           `json:"payment_method"`
	TransactionType string           `json:"transaction_type"`
	Currency        string           `json:"currency"`
	Country         string           `json:"country"`
	MinAmount       *decimal.Decimal `json:"min_amount"`
	MaxAmount       *decimal.Decimal `json:"max_amount"`
	Priority        int              `json:"priority"`
	Status          string           `json:"status"`
	EffectiveTime   *int64           `json:"effective_time"`
	ExpireTime      *int64           `json:"expire_time"`
	Remark          string           `json:"remark"`
	CreatedAt       int64            `json:"created_at"`
	UpdatedAt       int64            `json:"updated_at"`
}

// TrxRouterUpdateRequest 交易路由更新请求
type TrxRouterUpdateRequest struct {
	Name            *string          `json:"name,omitempty"`
	CashierTypes    []string         `json:"cashier_types,omitempty"`
	CashierIDs      []string         `json:"cashier_ids,omitempty"`
	CashierGroup    *string          `json:"cashier_group,omitempty"`
	PaymentMethod   *string          `json:"payment_method,omitempty"`
	TransactionType *string          `json:"transaction_type,omitempty"`
	Currency        *string          `json:"currency,omitempty"`
	Country         *string          `json:"country,omitempty"`
	MinAmount       *decimal.Decimal `json:"min_amount,omitempty"`
	MaxAmount       *decimal.Decimal `json:"max_amount,omitempty"`
	Priority        *int             `json:"priority,omitempty"`
	Status          *string          `json:"status,omitempty"`
	EffectiveTime   *int64           `json:"effective_time,omitempty"`
	ExpireTime      *int64           `json:"expire_time,omitempty"`
	Remark          *string          `json:"remark,omitempty"`
}

// TrxRouterListRequest 交易路由列表请求
type TrxRouterListRequest struct {
	Pagination
	Status          string `json:"status,omitempty" form:"status"`
	TransactionType string `json:"transaction_type,omitempty" form:"transaction_type"`
	PaymentMethod   string `json:"payment_method,omitempty" form:"payment_method"`
	Currency        string `json:"currency,omitempty" form:"currency"`
	Country         string `json:"country,omitempty" form:"country"`
	CashierType     string `json:"cashier_type,omitempty" form:"cashier_type"`
	CashierGroup    string `json:"cashier_group,omitempty" form:"cashier_group"`
}

// TrxRouterListResponse 交易路由列表响应
type TrxRouterListResponse struct {
	*PageResult
}

// RouteMatchRequest 路由匹配请求
type RouteMatchRequest struct {
	MerchantID      string          `json:"merchant_id" binding:"required"`
	TransactionType string          `json:"transaction_type" binding:"required"` // receipt, payment
	PaymentMethod   string          `json:"payment_method" binding:"required"`
	Amount          decimal.Decimal `json:"amount" binding:"required"`
	Currency        string          `json:"currency" binding:"required"`
	Country         string          `json:"country,omitempty"`
	BillID          string          `json:"bill_id,omitempty"`

	// 可选参数
	PreferredTypes []string `json:"preferred_types,omitempty"` // 偏好类型
	ForceRouterID  string   `json:"force_router_id,omitempty"` // 强制指定路由ID
}

// RouteMatchResponse 路由匹配响应
type RouteMatchResponse struct {
	Router      *TrxRouterResponse     `json:"router"`
	Cashier     *CashierInfoResponse   `json:"cashier"`
	MatchReason string                 `json:"match_reason"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// CashierInfoResponse Cashier信息响应
type CashierInfoResponse struct {
	CashierID   string                 `json:"cashier_id"`
	CashierType string                 `json:"cashier_type"` // business, personal
	Group       string                 `json:"group"`
	Status      string                 `json:"status"`
	Available   bool                   `json:"available"`
	Config      map[string]interface{} `json:"config"`
}

// RouterTestRequest 路由测试请求
type RouterTestRequest struct {
	RouterID        string          `json:"router_id" binding:"required"`
	TransactionType string          `json:"transaction_type" binding:"required"`
	PaymentMethod   string          `json:"payment_method" binding:"required"`
	Amount          decimal.Decimal `json:"amount" binding:"required"`
	Currency        string          `json:"currency" binding:"required"`
	Country         string          `json:"country,omitempty"`
}

// RouterTestResponse 路由测试响应
type RouterTestResponse struct {
	Matched          bool                   `json:"matched"`
	MatchDetails     map[string]interface{} `json:"match_details"`
	FailedReasons    []string               `json:"failed_reasons,omitempty"`
	AvailableCashier *CashierInfoResponse   `json:"available_cashier,omitempty"`
}

// RouterStatisticsRequest 路由统计请求
type RouterStatisticsRequest struct {
	RouterID  string `json:"router_id,omitempty" form:"router_id"`
	StartTime int64  `json:"start_time,omitempty" form:"start_time"`
	EndTime   int64  `json:"end_time,omitempty" form:"end_time"`
	DateRange string `json:"date_range,omitempty" form:"date_range"` // today, week, month
}

// RouterStatisticsResponse 路由统计响应
type RouterStatisticsResponse struct {
	RouterID        string                 `json:"router_id"`
	RouterName      string                 `json:"router_name"`
	TotalRequests   int64                  `json:"total_requests"`
	SuccessCount    int64                  `json:"success_count"`
	FailureCount    int64                  `json:"failure_count"`
	SuccessRate     float64                `json:"success_rate"`
	AvgResponseTime float64                `json:"avg_response_time"`
	Details         map[string]interface{} `json:"details"`
}
