package protocol

import "github.com/shopspring/decimal"

// 出纳员相关请求/响应
type CreateCashierRequest struct {
	AccountID    string          `json:"account_id" binding:"required"`                   // 关联的账户ID
	Type         string          `json:"type" binding:"required,oneof=private corporate"` // private(私户), corporate(公户)
	BankCode     string          `json:"bank_code" binding:"required"`
	BankName     string          `json:"bank_name" binding:"required"`
	CardNumber   string          `json:"card_number" binding:"required"`
	HolderName   string          `json:"holder_name" binding:"required"`
	HolderPhone  string          `json:"holder_phone"`
	HolderEmail  string          `json:"holder_email"`
	Country      string          `json:"country"`
	CountryCode  string          `json:"country_code"`
	Province     string          `json:"province"`
	City         string          `json:"city"`
	Currency     string          `json:"currency"`
	Usage        int32           `json:"usage"` // 1-收款，2-付款，4-储存
	Status       string          `json:"status"`
	DailyLimit   decimal.Decimal `json:"daily_limit"`
	MonthlyLimit decimal.Decimal `json:"monthly_limit"`
	ExpireAt     int64           `json:"expire_at"`
	Logo         string          `json:"logo"`
	Remark       string          `json:"remark"`
}

type UpdateCashierRequest struct {
	BankName     *string          `json:"bank_name"`
	HolderName   *string          `json:"holder_name"`
	HolderPhone  *string          `json:"holder_phone"`
	HolderEmail  *string          `json:"holder_email"`
	Province     *string          `json:"province"`
	City         *string          `json:"city"`
	Currency     *string          `json:"currency"`
	Usage        *int32           `json:"usage"`
	Status       *string          `json:"status"`
	DailyLimit   *decimal.Decimal `json:"daily_limit"`
	MonthlyLimit *decimal.Decimal `json:"monthly_limit"`
	DailyUsed    *decimal.Decimal `json:"daily_used"`
	MonthlyUsed  *decimal.Decimal `json:"monthly_used"`
	ExpireAt     *int64           `json:"expire_at"`
	Logo         *string          `json:"logo"`
	Remark       *string          `json:"remark"`
}

type ListCashiersRequest struct {
	Pagination
	AccountID  string `json:"account_id" form:"account_id"`
	Type       string `json:"type" form:"type"`
	BankCode   string `json:"bank_code" form:"bank_code"`
	Status     string `json:"status" form:"status"`
	Currency   string `json:"currency" form:"currency"`
	HolderName string `json:"holder_name" form:"holder_name"`
}

type Cashier struct {
	ID           uint64          `json:"id"`
	CashierID    string          `json:"cashier_id"`
	AccountID    string          `json:"account_id"`
	Type         string          `json:"type"`
	BankCode     string          `json:"bank_code"`
	BankName     string          `json:"bank_name"`
	CardNumber   string          `json:"card_number"`
	HolderName   string          `json:"holder_name"`
	HolderPhone  string          `json:"holder_phone"`
	HolderEmail  string          `json:"holder_email"`
	Country      string          `json:"country"`
	CountryCode  string          `json:"country_code"`
	Province     string          `json:"province"`
	City         string          `json:"city"`
	Currency     string          `json:"currency"`
	Usage        int32           `json:"usage"`
	Status       string          `json:"status"`
	DailyLimit   decimal.Decimal `json:"daily_limit"`
	MonthlyLimit decimal.Decimal `json:"monthly_limit"`
	DailyUsed    decimal.Decimal `json:"daily_used"`
	MonthlyUsed  decimal.Decimal `json:"monthly_used"`
	ExpireAt     int64           `json:"expire_at"`
	Logo         string          `json:"logo"`
	Remark       string          `json:"remark"`
	CreatedAt    int64           `json:"created_at"`
	UpdatedAt    int64           `json:"updated_at"`
}

type CashierTeam struct {
	ID      int64  `json:"id"`               // 主键ID
	Tid     string `json:"tid"`              // 团队ID
	Name    string `json:"name"`             // 车队名称
	Type    string `json:"type"`             // 车队类型
	Email   string `json:"email"`            // 车队邮箱
	Phone   string `json:"phone"`            // 车队电话
	Status  string `json:"status"`           // 车队状态
	Region  string `json:"region,omitempty"` // 车队区域
	Avatar  string `json:"avatar,omitempty"` // 车队头像
	HasG2FA bool   `json:"has_g2fa"`         // 是否启用二次验证
}
