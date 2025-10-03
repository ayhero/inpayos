package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

// Channel 支付渠道表
type Channel struct {
	ID        uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ChannelID string `json:"channel_id" gorm:"column:channel_id;type:varchar(64);uniqueIndex"`
	*ChannelValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type ChannelValues struct {
	Code         *string          `json:"code" gorm:"column:code;type:varchar(32);index"`
	Name         *string          `json:"name" gorm:"column:name;type:varchar(128)"`
	Type         *string          `json:"type" gorm:"column:type;type:varchar(16);index"`                      // receipt, payment, deposit, withdraw
	PayMethods   *string          `json:"pay_methods" gorm:"column:pay_methods;type:text"`                     // JSON数组格式的支付方式
	Currencies   *string          `json:"currencies" gorm:"column:currencies;type:text"`                       // JSON数组格式的货币
	Status       *string          `json:"status" gorm:"column:status;type:varchar(16);index;default:'active'"` // active, inactive, maintain
	FeeType      *string          `json:"fee_type" gorm:"column:fee_type;type:varchar(16)"`                    // fixed, percent
	FeeValue     *decimal.Decimal `json:"fee_value" gorm:"column:fee_value;type:decimal(20,8);default:0"`
	MinAmount    *decimal.Decimal `json:"min_amount" gorm:"column:min_amount;type:decimal(20,8);default:0"`
	MaxAmount    *decimal.Decimal `json:"max_amount" gorm:"column:max_amount;type:decimal(20,8);default:0"`
	DailyLimit   *decimal.Decimal `json:"daily_limit" gorm:"column:daily_limit;type:decimal(20,8);default:0"`
	MonthlyLimit *decimal.Decimal `json:"monthly_limit" gorm:"column:monthly_limit;type:decimal(20,8);default:0"`
	DailyUsed    *decimal.Decimal `json:"daily_used" gorm:"column:daily_used;type:decimal(20,8);default:0"`
	MonthlyUsed  *decimal.Decimal `json:"monthly_used" gorm:"column:monthly_used;type:decimal(20,8);default:0"`
	Config       *string          `json:"config" gorm:"column:config;type:text"` // JSON格式的配置信息
	Remark       *string          `json:"remark" gorm:"column:remark;type:varchar(512)"`
}

// 表名
func (Channel) TableName() string {
	return "t_channels"
}

// Chainable setters
func (v *ChannelValues) SetCode(code string) *ChannelValues {
	v.Code = &code
	return v
}

func (v *ChannelValues) SetName(name string) *ChannelValues {
	v.Name = &name
	return v
}

func (v *ChannelValues) SetType(txType string) *ChannelValues {
	v.Type = &txType
	return v
}

func (v *ChannelValues) SetPayMethods(methods string) *ChannelValues {
	v.PayMethods = &methods
	return v
}

func (v *ChannelValues) SetCurrencies(currencies string) *ChannelValues {
	v.Currencies = &currencies
	return v
}

func (v *ChannelValues) SetStatus(status string) *ChannelValues {
	v.Status = &status
	return v
}

func (v *ChannelValues) SetFeeType(feeType string) *ChannelValues {
	v.FeeType = &feeType
	return v
}

func (v *ChannelValues) SetFeeValue(value decimal.Decimal) *ChannelValues {
	v.FeeValue = &value
	return v
}

func (v *ChannelValues) SetFeeFixed(fee decimal.Decimal) *ChannelValues {
	v.FeeType = StringPtr("fixed")
	v.FeeValue = &fee
	return v
}

func (v *ChannelValues) SetFeePercent(percent decimal.Decimal) *ChannelValues {
	v.FeeType = StringPtr("percent")
	v.FeeValue = &percent
	return v
}

func (v *ChannelValues) SetMinAmount(amount decimal.Decimal) *ChannelValues {
	v.MinAmount = &amount
	return v
}

func (v *ChannelValues) SetMaxAmount(amount decimal.Decimal) *ChannelValues {
	v.MaxAmount = &amount
	return v
}

func (v *ChannelValues) SetDailyLimit(limit decimal.Decimal) *ChannelValues {
	v.DailyLimit = &limit
	return v
}

func (v *ChannelValues) SetMonthlyLimit(limit decimal.Decimal) *ChannelValues {
	v.MonthlyLimit = &limit
	return v
}

func (v *ChannelValues) SetDailyUsed(used decimal.Decimal) *ChannelValues {
	v.DailyUsed = &used
	return v
}

func (v *ChannelValues) SetMonthlyUsed(used decimal.Decimal) *ChannelValues {
	v.MonthlyUsed = &used
	return v
}

func (v *ChannelValues) SetConfig(config string) *ChannelValues {
	v.Config = &config
	return v
}

func (v *ChannelValues) SetRemark(remark string) *ChannelValues {
	v.Remark = &remark
	return v
}

// Helper function
func StringPtr(s string) *string {
	return &s
}

// Chainable getters
func (v *ChannelValues) GetCode() string {
	if v.Code == nil {
		return ""
	}
	return *v.Code
}

func (v *ChannelValues) GetName() string {
	if v.Name == nil {
		return ""
	}
	return *v.Name
}

func (v *ChannelValues) GetType() string {
	if v.Type == nil {
		return ""
	}
	return *v.Type
}

func (v *ChannelValues) GetPayMethods() string {
	if v.PayMethods == nil {
		return "[]"
	}
	return *v.PayMethods
}

func (v *ChannelValues) GetCurrencies() string {
	if v.Currencies == nil {
		return "[]"
	}
	return *v.Currencies
}

func (v *ChannelValues) GetStatus() string {
	if v.Status == nil {
		return "active"
	}
	return *v.Status
}

func (v *ChannelValues) GetFeeType() string {
	if v.FeeType == nil {
		return "fixed"
	}
	return *v.FeeType
}

func (v *ChannelValues) GetFeeValue() decimal.Decimal {
	if v.FeeValue == nil {
		return decimal.Zero
	}
	return *v.FeeValue
}

func (v *ChannelValues) GetFeeFixed() decimal.Decimal {
	if v.FeeType != nil && *v.FeeType == "fixed" && v.FeeValue != nil {
		return *v.FeeValue
	}
	return decimal.Zero
}

func (v *ChannelValues) GetFeePercent() decimal.Decimal {
	if v.FeeType != nil && *v.FeeType == "percent" && v.FeeValue != nil {
		return *v.FeeValue
	}
	return decimal.Zero
}

func (v *ChannelValues) GetMinAmount() decimal.Decimal {
	if v.MinAmount == nil {
		return decimal.Zero
	}
	return *v.MinAmount
}

func (v *ChannelValues) GetMaxAmount() decimal.Decimal {
	if v.MaxAmount == nil {
		return decimal.Zero
	}
	return *v.MaxAmount
}

func (v *ChannelValues) GetDailyLimit() decimal.Decimal {
	if v.DailyLimit == nil {
		return decimal.Zero
	}
	return *v.DailyLimit
}

func (v *ChannelValues) GetMonthlyLimit() decimal.Decimal {
	if v.MonthlyLimit == nil {
		return decimal.Zero
	}
	return *v.MonthlyLimit
}

func (v *ChannelValues) GetDailyUsed() decimal.Decimal {
	if v.DailyUsed == nil {
		return decimal.Zero
	}
	return *v.DailyUsed
}

func (v *ChannelValues) GetMonthlyUsed() decimal.Decimal {
	if v.MonthlyUsed == nil {
		return decimal.Zero
	}
	return *v.MonthlyUsed
}

func (v *ChannelValues) GetConfig() string {
	if v.Config == nil {
		return ""
	}
	return *v.Config
}

func (v *ChannelValues) GetRemark() string {
	if v.Remark == nil {
		return ""
	}
	return *v.Remark
}

// CalculateFee 计算手续费
func (v *ChannelValues) CalculateFee(amount decimal.Decimal) decimal.Decimal {
	feeType := v.GetFeeType()
	feeValue := v.GetFeeValue()

	if feeType == "fixed" {
		return feeValue
	} else if feeType == "percent" {
		return amount.Mul(feeValue).Div(decimal.NewFromInt(100))
	}

	return decimal.Zero
}

// IsAmountValid 验证金额是否在允许范围内
func (v *ChannelValues) IsAmountValid(amount decimal.Decimal) bool {
	minAmount := v.GetMinAmount()
	maxAmount := v.GetMaxAmount()

	if minAmount.GreaterThan(decimal.Zero) && amount.LessThan(minAmount) {
		return false
	}

	if maxAmount.GreaterThan(decimal.Zero) && amount.GreaterThan(maxAmount) {
		return false
	}

	return true
}

// IsActive 检查渠道是否可用
func (v *ChannelValues) IsActive() bool {
	return v.GetStatus() == "active"
}

// =============================================================================
// TrxRouter 交易路由相关定义
// =============================================================================

// StringArray 自定义字符串数组类型，用于存储JSON数组
type StringArray []string

// Scan 实现 sql.Scanner 接口
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = StringArray{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return fmt.Errorf("cannot scan %T into StringArray", value)
	}
}

// Value 实现 driver.Valuer 接口
func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

// TrxRouter 交易路由表
type TrxRouter struct {
	ID       uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	RouterID string `json:"router_id" gorm:"column:router_id;type:varchar(64);uniqueIndex"`
	*TrxRouterValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type TrxRouterValues struct {
	Name            *string          `json:"name" gorm:"column:name;type:varchar(128)"`                        // 路由名称
	CashierGroups   *[]string        `json:"cashier_group" gorm:"column:cashier_group;type:varchar(64)"`       // Cashier分组
	PaymentMethod   *string          `json:"payment_method" gorm:"column:payment_method;type:varchar(32)"`     // 支付方式
	TransactionType *string          `json:"transaction_type" gorm:"column:transaction_type;type:varchar(16)"` // receipt, payment
	Ccy             *string          `json:"ccy" gorm:"column:ccy;type:varchar(10)"`                           // 币种
	Country         *string          `json:"country" gorm:"column:country;type:varchar(3)"`                    // 国家代码
	MinAmount       *decimal.Decimal `json:"min_amount" gorm:"column:min_amount;type:decimal(20,8)"`           // 最小金额
	MaxAmount       *decimal.Decimal `json:"max_amount" gorm:"column:max_amount;type:decimal(20,8)"`           // 最大金额
	Priority        *int             `json:"priority" gorm:"column:priority;default:100"`                      // 优先级
	Status          *string          `json:"status" gorm:"column:status;type:varchar(16);default:'active'"`    // 状态
	EffectiveTime   *int64           `json:"effective_time" gorm:"column:effective_time"`                      // 生效时间
	ExpireTime      *int64           `json:"expire_time" gorm:"column:expire_time"`                            // 过期时间
	Remark          *string          `json:"remark" gorm:"column:remark;type:varchar(512)"`                    // 备注
}

// 表名
func (TrxRouter) TableName() string {
	return "t_trx_routers"
}
