package protocol

import "github.com/shopspring/decimal"

// 渠道相关请求/响应
type CreateChannelRequest struct {
	Code         string                 `json:"code" binding:"required"`
	Name         string                 `json:"name" binding:"required"`
	Type         string                 `json:"type" binding:"required"`
	PayMethods   []string               `json:"pay_methods" binding:"required"`
	Currencies   []string               `json:"currencies" binding:"required"`
	Status       string                 `json:"status" binding:"required,oneof=active inactive maintain"`
	FeeType      string                 `json:"fee_type" binding:"required,oneof=fixed percent"`
	FeeValue     decimal.Decimal        `json:"fee_value" binding:"required"`
	MinAmount    decimal.Decimal        `json:"min_amount"`
	MaxAmount    decimal.Decimal        `json:"max_amount"`
	DailyLimit   decimal.Decimal        `json:"daily_limit"`
	MonthlyLimit decimal.Decimal        `json:"monthly_limit"`
	Config       map[string]interface{} `json:"config"`
	Remark       string                 `json:"remark"`
}

type UpdateChannelRequest struct {
	Name         *string                `json:"name"`
	Type         *string                `json:"type"`
	PayMethods   []string               `json:"pay_methods"`
	Currencies   []string               `json:"currencies"`
	Status       *string                `json:"status"`
	FeeType      *string                `json:"fee_type"`
	FeeValue     *decimal.Decimal       `json:"fee_value"`
	MinAmount    *decimal.Decimal       `json:"min_amount"`
	MaxAmount    *decimal.Decimal       `json:"max_amount"`
	DailyLimit   *decimal.Decimal       `json:"daily_limit"`
	MonthlyLimit *decimal.Decimal       `json:"monthly_limit"`
	Config       map[string]interface{} `json:"config"`
	Remark       *string                `json:"remark"`
}

type ListChannelsRequest struct {
	Pagination
	Code      string `json:"code" form:"code"`
	Name      string `json:"name" form:"name"`
	Type      string `json:"type" form:"type"`
	Status    string `json:"status" form:"status"`
	PayMethod string `json:"pay_method" form:"pay_method"`
	Currency  string `json:"currency" form:"currency"`
}

type ChannelResponse struct {
	ID           uint64                 `json:"id"`
	Code         string                 `json:"code"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	PayMethods   []string               `json:"pay_methods"`
	Currencies   []string               `json:"currencies"`
	Status       string                 `json:"status"`
	FeeType      string                 `json:"fee_type"`
	FeeValue     decimal.Decimal        `json:"fee_value"`
	MinAmount    decimal.Decimal        `json:"min_amount"`
	MaxAmount    decimal.Decimal        `json:"max_amount"`
	DailyLimit   decimal.Decimal        `json:"daily_limit"`
	MonthlyLimit decimal.Decimal        `json:"monthly_limit"`
	DailyUsed    decimal.Decimal        `json:"daily_used"`
	MonthlyUsed  decimal.Decimal        `json:"monthly_used"`
	Config       map[string]interface{} `json:"config"`
	Remark       string                 `json:"remark"`
	CreatedAt    int64                  `json:"created_at"`
	UpdatedAt    int64                  `json:"updated_at"`
}

type ChannelStatsResponse struct {
	TotalChannels    int64 `json:"total_channels"`
	ActiveChannels   int64 `json:"active_channels"`
	DisabledChannels int64 `json:"disabled_channels"`
	MaintainChannels int64 `json:"maintain_channels"`
}
