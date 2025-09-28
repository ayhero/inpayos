package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Cashier 出纳员/收银员表（区分公户和私户）
type Cashier struct {
	ID        uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	CashierID string `json:"cashier_id" gorm:"column:cashier_id;type:varchar(64);uniqueIndex"`
	AccountID string `json:"account_id" gorm:"column:account_id;type:varchar(64);index"` // 关联的账户ID
	*CashierValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type CashierValues struct {
	// 基础信息
	Type        *string `json:"type" gorm:"column:type;type:varchar(16);index;default:'private'"` // private(私户), corporate(公户)
	BankCode    *string `json:"bank_code" gorm:"column:bank_code;type:varchar(32);index"`         // 银行代码
	BankName    *string `json:"bank_name" gorm:"column:bank_name;type:varchar(128)"`              // 银行名称
	CardNumber  *string `json:"card_number" gorm:"column:card_number;type:varchar(32);index"`     // 卡号
	HolderName  *string `json:"holder_name" gorm:"column:holder_name;type:varchar(128)"`          // 持卡人姓名
	HolderPhone *string `json:"holder_phone" gorm:"column:holder_phone;type:varchar(32)"`         // 持卡人手机
	HolderEmail *string `json:"holder_email" gorm:"column:holder_email;type:varchar(128)"`        // 持卡人邮箱

	// 地域信息
	Country     *string `json:"country" gorm:"column:country;type:varchar(8);index"`     // 国家
	CountryCode *string `json:"country_code" gorm:"column:country_code;type:varchar(8)"` // 国家代码
	Province    *string `json:"province" gorm:"column:province;type:varchar(64)"`        // 省/州
	City        *string `json:"city" gorm:"column:city;type:varchar(64)"`                // 城市

	// 业务配置
	Currency     *string          `json:"currency" gorm:"column:currency;type:varchar(8);index;default:'CNY'"`    // 币种
	Usage        *int32           `json:"usage" gorm:"column:usage;default:0"`                                    // 用途权限位授权：1-收款，2-付款，4-储存
	Status       *string          `json:"status" gorm:"column:status;type:varchar(16);default:'active'"`          // active, inactive, frozen, suspended
	DailyLimit   *decimal.Decimal `json:"daily_limit" gorm:"column:daily_limit;type:decimal(20,8);default:0"`     // 每日限额
	MonthlyLimit *decimal.Decimal `json:"monthly_limit" gorm:"column:monthly_limit;type:decimal(20,8);default:0"` // 每月限额
	DailyUsed    *decimal.Decimal `json:"daily_used" gorm:"column:daily_used;type:decimal(20,8);default:0"`       // 每日已使用
	MonthlyUsed  *decimal.Decimal `json:"monthly_used" gorm:"column:monthly_used;type:decimal(20,8);default:0"`   // 每月已使用

	// 其他信息
	ExpireAt *int64  `json:"expire_at" gorm:"column:expire_at"`             // 过期时间
	Logo     *string `json:"logo" gorm:"column:logo;type:varchar(512)"`     // 头像/标志
	Remark   *string `json:"remark" gorm:"column:remark;type:varchar(512)"` // 备注
}

// 表名
func (Cashier) TableName() string {
	return "t_cashiers"
}

// Chainable setters
func (v *CashierValues) SetType(cashierType string) *CashierValues {
	v.Type = &cashierType
	return v
}

func (v *CashierValues) SetBankCode(code string) *CashierValues {
	v.BankCode = &code
	return v
}

func (v *CashierValues) SetBankName(name string) *CashierValues {
	v.BankName = &name
	return v
}

func (v *CashierValues) SetCardNumber(number string) *CashierValues {
	v.CardNumber = &number
	return v
}

func (v *CashierValues) SetHolderName(name string) *CashierValues {
	v.HolderName = &name
	return v
}

func (v *CashierValues) SetHolderPhone(phone string) *CashierValues {
	v.HolderPhone = &phone
	return v
}

func (v *CashierValues) SetHolderEmail(email string) *CashierValues {
	v.HolderEmail = &email
	return v
}

func (v *CashierValues) SetCountry(country string) *CashierValues {
	v.Country = &country
	return v
}

func (v *CashierValues) SetCountryCode(code string) *CashierValues {
	v.CountryCode = &code
	return v
}

func (v *CashierValues) SetProvince(province string) *CashierValues {
	v.Province = &province
	return v
}

func (v *CashierValues) SetCity(city string) *CashierValues {
	v.City = &city
	return v
}

func (v *CashierValues) SetCurrency(currency string) *CashierValues {
	v.Currency = &currency
	return v
}

func (v *CashierValues) SetUsage(usage int32) *CashierValues {
	v.Usage = &usage
	return v
}

func (v *CashierValues) SetStatus(status string) *CashierValues {
	v.Status = &status
	return v
}

func (v *CashierValues) SetDailyLimit(limit decimal.Decimal) *CashierValues {
	v.DailyLimit = &limit
	return v
}

func (v *CashierValues) SetMonthlyLimit(limit decimal.Decimal) *CashierValues {
	v.MonthlyLimit = &limit
	return v
}

func (v *CashierValues) SetDailyUsed(used decimal.Decimal) *CashierValues {
	v.DailyUsed = &used
	return v
}

func (v *CashierValues) SetMonthlyUsed(used decimal.Decimal) *CashierValues {
	v.MonthlyUsed = &used
	return v
}

func (v *CashierValues) SetExpireAt(time int64) *CashierValues {
	v.ExpireAt = &time
	return v
}

func (v *CashierValues) SetLogo(logo string) *CashierValues {
	v.Logo = &logo
	return v
}

func (v *CashierValues) SetRemark(remark string) *CashierValues {
	v.Remark = &remark
	return v
}

// Chainable getters
func (v *CashierValues) GetType() string {
	if v.Type == nil {
		return "private"
	}
	return *v.Type
}

func (v *CashierValues) GetBankCode() string {
	if v.BankCode == nil {
		return ""
	}
	return *v.BankCode
}

func (v *CashierValues) GetBankName() string {
	if v.BankName == nil {
		return ""
	}
	return *v.BankName
}

func (v *CashierValues) GetCardNumber() string {
	if v.CardNumber == nil {
		return ""
	}
	return *v.CardNumber
}

func (v *CashierValues) GetHolderName() string {
	if v.HolderName == nil {
		return ""
	}
	return *v.HolderName
}

func (v *CashierValues) GetHolderPhone() string {
	if v.HolderPhone == nil {
		return ""
	}
	return *v.HolderPhone
}

func (v *CashierValues) GetHolderEmail() string {
	if v.HolderEmail == nil {
		return ""
	}
	return *v.HolderEmail
}

func (v *CashierValues) GetCountry() string {
	if v.Country == nil {
		return ""
	}
	return *v.Country
}

func (v *CashierValues) GetCountryCode() string {
	if v.CountryCode == nil {
		return ""
	}
	return *v.CountryCode
}

func (v *CashierValues) GetProvince() string {
	if v.Province == nil {
		return ""
	}
	return *v.Province
}

func (v *CashierValues) GetCity() string {
	if v.City == nil {
		return ""
	}
	return *v.City
}

func (v *CashierValues) GetCurrency() string {
	if v.Currency == nil {
		return "CNY"
	}
	return *v.Currency
}

func (v *CashierValues) GetUsage() int32 {
	if v.Usage == nil {
		return 0
	}
	return *v.Usage
}

func (v *CashierValues) GetStatus() string {
	if v.Status == nil {
		return "active"
	}
	return *v.Status
}

func (v *CashierValues) GetDailyLimit() decimal.Decimal {
	if v.DailyLimit == nil {
		return decimal.Zero
	}
	return *v.DailyLimit
}

func (v *CashierValues) GetMonthlyLimit() decimal.Decimal {
	if v.MonthlyLimit == nil {
		return decimal.Zero
	}
	return *v.MonthlyLimit
}

func (v *CashierValues) GetDailyUsed() decimal.Decimal {
	if v.DailyUsed == nil {
		return decimal.Zero
	}
	return *v.DailyUsed
}

func (v *CashierValues) GetMonthlyUsed() decimal.Decimal {
	if v.MonthlyUsed == nil {
		return decimal.Zero
	}
	return *v.MonthlyUsed
}

func (v *CashierValues) GetExpireAt() int64 {
	if v.ExpireAt == nil {
		return 0
	}
	return *v.ExpireAt
}

func (v *CashierValues) GetLogo() string {
	if v.Logo == nil {
		return ""
	}
	return *v.Logo
}

func (v *CashierValues) GetRemark() string {
	if v.Remark == nil {
		return ""
	}
	return *v.Remark
}

// 业务方法

// IsPrivate 检查是否为私户
func (v *CashierValues) IsPrivate() bool {
	return v.GetType() == "private"
}

// IsCorporate 检查是否为公户
func (v *CashierValues) IsCorporate() bool {
	return v.GetType() == "corporate"
}

// IsActive 检查是否为活跃状态
func (v *CashierValues) IsActive() bool {
	return v.GetStatus() == "active"
}

// CanReceive 检查是否可以收款
func (v *CashierValues) CanReceive() bool {
	return (v.GetUsage() & 1) != 0
}

// CanPay 检查是否可以付款
func (v *CashierValues) CanPay() bool {
	return (v.GetUsage() & 2) != 0
}

// CanStorage 检查是否可以储存
func (v *CashierValues) CanStorage() bool {
	return (v.GetUsage() & 4) != 0
}

// IsExpired 检查是否已过期
func (v *CashierValues) IsExpired() bool {
	if v.GetExpireAt() == 0 {
		return false
	}
	now := time.Now().UnixMilli()
	return now > v.GetExpireAt()
}

// CanUseToday 检查今日是否还可以使用（未超过每日限额）
func (v *CashierValues) CanUseToday(amount decimal.Decimal) bool {
	dailyLimit := v.GetDailyLimit()
	if dailyLimit.IsZero() {
		return true // 无限额
	}

	dailyUsed := v.GetDailyUsed()
	return dailyUsed.Add(amount).LessThanOrEqual(dailyLimit)
}

// CanUseThisMonth 检查本月是否还可以使用（未超过每月限额）
func (v *CashierValues) CanUseThisMonth(amount decimal.Decimal) bool {
	monthlyLimit := v.GetMonthlyLimit()
	if monthlyLimit.IsZero() {
		return true // 无限额
	}

	monthlyUsed := v.GetMonthlyUsed()
	return monthlyUsed.Add(amount).LessThanOrEqual(monthlyLimit)
}

// AddDailyUsage 增加每日使用量
func (v *CashierValues) AddDailyUsage(amount decimal.Decimal) *CashierValues {
	current := v.GetDailyUsed()
	v.SetDailyUsed(current.Add(amount))
	return v
}

// AddMonthlyUsage 增加每月使用量
func (v *CashierValues) AddMonthlyUsage(amount decimal.Decimal) *CashierValues {
	current := v.GetMonthlyUsed()
	v.SetMonthlyUsed(current.Add(amount))
	return v
}

// ResetDailyUsage 重置每日使用量
func (v *CashierValues) ResetDailyUsage() *CashierValues {
	v.SetDailyUsed(decimal.Zero)
	return v
}

// ResetMonthlyUsage 重置每月使用量
func (v *CashierValues) ResetMonthlyUsage() *CashierValues {
	v.SetMonthlyUsed(decimal.Zero)
	return v
}
