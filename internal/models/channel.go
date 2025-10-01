package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

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
	CashierTypes    StringArray      `json:"cashier_types" gorm:"column:cashier_types;type:text"`              // Cashier类型数组
	CashierIDs      StringArray      `json:"cashier_ids" gorm:"column:cashier_ids;type:text"`                  // 指定Cashier ID数组
	CashierGroup    *string          `json:"cashier_group" gorm:"column:cashier_group;type:varchar(64)"`       // Cashier分组
	PaymentMethod   *string          `json:"payment_method" gorm:"column:payment_method;type:varchar(32)"`     // 支付方式
	TransactionType *string          `json:"transaction_type" gorm:"column:transaction_type;type:varchar(16)"` // receipt, payment
	Currency        *string          `json:"currency" gorm:"column:currency;type:varchar(10)"`                 // 币种
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

// Validate 验证路由配置
func (v *TrxRouterValues) Validate() error {
	hasTypes := len(v.CashierTypes) > 0
	hasIDs := len(v.CashierIDs) > 0
	hasGroup := v.CashierGroup != nil && *v.CashierGroup != ""

	// 计算有值的字段数量
	count := 0
	if hasTypes {
		count++
	}
	if hasIDs {
		count++
	}
	if hasGroup {
		count++
	}

	// 确保只有一个字段有值
	if count != 1 {
		return fmt.Errorf("CashierTypes, CashierIDs, and CashierGroup cannot have values simultaneously, exactly one must be set")
	}

	return nil
}

// =============================================================================
// TrxRouter Chainable setters
// =============================================================================

func (v *TrxRouterValues) SetName(name string) *TrxRouterValues {
	v.Name = &name
	return v
}

func (v *TrxRouterValues) SetCashierTypes(types []string) *TrxRouterValues {
	v.CashierTypes = StringArray(types)
	// 清空其他字段确保互斥
	v.CashierIDs = StringArray{}
	v.CashierGroup = nil
	return v
}

func (v *TrxRouterValues) SetCashierIDs(ids []string) *TrxRouterValues {
	v.CashierIDs = StringArray(ids)
	// 清空其他字段确保互斥
	v.CashierTypes = StringArray{}
	v.CashierGroup = nil
	return v
}

func (v *TrxRouterValues) SetCashierGroup(group string) *TrxRouterValues {
	v.CashierGroup = &group
	// 清空其他字段确保互斥
	v.CashierTypes = StringArray{}
	v.CashierIDs = StringArray{}
	return v
}

func (v *TrxRouterValues) SetPaymentMethod(method string) *TrxRouterValues {
	v.PaymentMethod = &method
	return v
}

func (v *TrxRouterValues) SetTransactionType(txType string) *TrxRouterValues {
	v.TransactionType = &txType
	return v
}

func (v *TrxRouterValues) SetCurrency(currency string) *TrxRouterValues {
	v.Currency = &currency
	return v
}

func (v *TrxRouterValues) SetCountry(country string) *TrxRouterValues {
	v.Country = &country
	return v
}

func (v *TrxRouterValues) SetAmountRange(min, max decimal.Decimal) *TrxRouterValues {
	v.MinAmount = &min
	v.MaxAmount = &max
	return v
}

func (v *TrxRouterValues) SetMinAmount(amount decimal.Decimal) *TrxRouterValues {
	v.MinAmount = &amount
	return v
}

func (v *TrxRouterValues) SetMaxAmount(amount decimal.Decimal) *TrxRouterValues {
	v.MaxAmount = &amount
	return v
}

func (v *TrxRouterValues) SetPriority(priority int) *TrxRouterValues {
	v.Priority = &priority
	return v
}

func (v *TrxRouterValues) SetStatus(status string) *TrxRouterValues {
	v.Status = &status
	return v
}

func (v *TrxRouterValues) SetTimeRange(effective, expire int64) *TrxRouterValues {
	v.EffectiveTime = &effective
	v.ExpireTime = &expire
	return v
}

func (v *TrxRouterValues) SetRemark(remark string) *TrxRouterValues {
	v.Remark = &remark
	return v
}

// =============================================================================
// TrxRouter Chainable getters
// =============================================================================

func (v *TrxRouterValues) GetName() string {
	if v.Name == nil {
		return ""
	}
	return *v.Name
}

func (v *TrxRouterValues) GetCashierTypes() []string {
	return []string(v.CashierTypes)
}

func (v *TrxRouterValues) GetCashierIDs() []string {
	return []string(v.CashierIDs)
}

func (v *TrxRouterValues) GetCashierGroup() string {
	if v.CashierGroup == nil {
		return ""
	}
	return *v.CashierGroup
}

func (v *TrxRouterValues) GetPaymentMethod() string {
	if v.PaymentMethod == nil {
		return ""
	}
	return *v.PaymentMethod
}

func (v *TrxRouterValues) GetTransactionType() string {
	if v.TransactionType == nil {
		return ""
	}
	return *v.TransactionType
}

func (v *TrxRouterValues) GetCurrency() string {
	if v.Currency == nil {
		return ""
	}
	return *v.Currency
}

func (v *TrxRouterValues) GetCountry() string {
	if v.Country == nil {
		return ""
	}
	return *v.Country
}

func (v *TrxRouterValues) GetMinAmount() decimal.Decimal {
	if v.MinAmount == nil {
		return decimal.Zero
	}
	return *v.MinAmount
}

func (v *TrxRouterValues) GetMaxAmount() decimal.Decimal {
	if v.MaxAmount == nil {
		return decimal.Zero
	}
	return *v.MaxAmount
}

func (v *TrxRouterValues) GetPriority() int {
	if v.Priority == nil {
		return 100
	}
	return *v.Priority
}

func (v *TrxRouterValues) GetStatus() string {
	if v.Status == nil {
		return "active"
	}
	return *v.Status
}

func (v *TrxRouterValues) GetEffectiveTime() int64 {
	if v.EffectiveTime == nil {
		return 0
	}
	return *v.EffectiveTime
}

func (v *TrxRouterValues) GetExpireTime() int64 {
	if v.ExpireTime == nil {
		return 0
	}
	return *v.ExpireTime
}

func (v *TrxRouterValues) GetRemark() string {
	if v.Remark == nil {
		return ""
	}
	return *v.Remark
}

// =============================================================================
// TrxRouter 业务逻辑方法
// =============================================================================

// IsActive 检查路由是否激活
func (v *TrxRouterValues) IsActive() bool {
	return v.GetStatus() == "active"
}

// IsInAmountRange 检查金额是否在路由范围内
func (v *TrxRouterValues) IsInAmountRange(amount decimal.Decimal) bool {
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

// IsInTimeRange 检查是否在有效时间范围内
func (v *TrxRouterValues) IsInTimeRange(timestamp int64) bool {
	effectiveTime := v.GetEffectiveTime()
	expireTime := v.GetExpireTime()

	// 如果没有设置时间范围，则认为总是有效
	if effectiveTime == 0 && expireTime == 0 {
		return true
	}

	// 检查生效时间
	if effectiveTime > 0 && timestamp < effectiveTime {
		return false
	}

	// 检查过期时间
	if expireTime > 0 && timestamp > expireTime {
		return false
	}

	return true
}

// IsEffective 检查路由是否当前有效
func (v *TrxRouterValues) IsEffective() bool {
	return v.IsActive() && v.IsInTimeRange(getCurrentTimeMillis())
}

// MatchesRequest 检查路由是否匹配请求
func (v *TrxRouterValues) MatchesRequest(txType, paymentMethod, currency, country string, amount decimal.Decimal) bool {
	// 检查交易类型
	if v.GetTransactionType() != "" && v.GetTransactionType() != "*" && v.GetTransactionType() != txType {
		return false
	}

	// 检查支付方式
	if v.GetPaymentMethod() != "" && v.GetPaymentMethod() != "*" && v.GetPaymentMethod() != paymentMethod {
		return false
	}

	// 检查币种
	if v.GetCurrency() != "" && v.GetCurrency() != "*" && v.GetCurrency() != currency {
		return false
	}

	// 检查国家
	if v.GetCountry() != "" && v.GetCountry() != "*" && v.GetCountry() != country {
		return false
	}

	// 检查金额范围
	if !v.IsInAmountRange(amount) {
		return false
	}

	// 检查路由是否有效
	if !v.IsEffective() {
		return false
	}

	return true
}

// GetCashierSelector 获取Cashier选择器信息
func (v *TrxRouterValues) GetCashierSelector() map[string]interface{} {
	result := make(map[string]interface{})

	if len(v.CashierTypes) > 0 {
		result["type"] = "cashier_types"
		result["values"] = v.GetCashierTypes()
	} else if len(v.CashierIDs) > 0 {
		result["type"] = "cashier_ids"
		result["values"] = v.GetCashierIDs()
	} else if v.GetCashierGroup() != "" {
		result["type"] = "cashier_group"
		result["values"] = v.GetCashierGroup()
	}

	return result
}

// HasCashierType 检查是否包含指定的Cashier类型
func (v *TrxRouterValues) HasCashierType(cashierType string) bool {
	for _, t := range v.CashierTypes {
		if t == cashierType {
			return true
		}
	}
	return false
}

// HasCashierID 检查是否包含指定的Cashier ID
func (v *TrxRouterValues) HasCashierID(cashierID string) bool {
	for _, id := range v.CashierIDs {
		if id == cashierID {
			return true
		}
	}
	return false
}

// String 实现 Stringer 接口，便于调试
func (v *TrxRouterValues) String() string {
	var parts []string

	if v.GetName() != "" {
		parts = append(parts, fmt.Sprintf("Name:%s", v.GetName()))
	}

	if len(v.CashierTypes) > 0 {
		parts = append(parts, fmt.Sprintf("Types:%v", v.GetCashierTypes()))
	}

	if len(v.CashierIDs) > 0 {
		parts = append(parts, fmt.Sprintf("IDs:%v", v.GetCashierIDs()))
	}

	if v.GetCashierGroup() != "" {
		parts = append(parts, fmt.Sprintf("Group:%s", v.GetCashierGroup()))
	}

	if v.GetTransactionType() != "" {
		parts = append(parts, fmt.Sprintf("TxType:%s", v.GetTransactionType()))
	}

	return strings.Join(parts, ", ")
}

// =============================================================================
// TrxRouter 数据库操作方法
// =============================================================================

// GetTrxRoutersByCondition 根据条件获取交易路由
func GetTrxRoutersByCondition(txType, paymentMethod, currency, country string, amount decimal.Decimal) ([]*TrxRouter, error) {
	var routers []*TrxRouter

	query := DB.Where("status = ?", "active")

	// 添加匹配条件
	if txType != "" {
		query = query.Where("(transaction_type = ? OR transaction_type = '*' OR transaction_type IS NULL)", txType)
	}

	if paymentMethod != "" {
		query = query.Where("(payment_method = ? OR payment_method = '*' OR payment_method IS NULL)", paymentMethod)
	}

	if currency != "" {
		query = query.Where("(currency = ? OR currency = '*' OR currency IS NULL)", currency)
	}

	if country != "" {
		query = query.Where("(country = ? OR country = '*' OR country IS NULL)", country)
	}

	// 金额范围过滤
	if amount.GreaterThan(decimal.Zero) {
		query = query.Where("(min_amount IS NULL OR min_amount <= ?)", amount).
			Where("(max_amount IS NULL OR max_amount >= ?)", amount)
	}

	// 按优先级排序
	err := query.Order("priority DESC, created_at ASC").Find(&routers).Error
	if err != nil {
		return nil, err
	}

	// 进一步过滤匹配的路由
	var matchedRouters []*TrxRouter
	for _, router := range routers {
		if router.TrxRouterValues.MatchesRequest(txType, paymentMethod, currency, country, amount) {
			matchedRouters = append(matchedRouters, router)
		}
	}

	return matchedRouters, nil
}

// GetTrxRouterByID 根据路由ID获取交易路由
func GetTrxRouterByID(routerID string) (*TrxRouter, error) {
	var router TrxRouter
	err := DB.Where("router_id = ?", routerID).First(&router).Error
	if err != nil {
		return nil, err
	}
	return &router, nil
}

// CreateTrxRouter 创建交易路由
func CreateTrxRouter(router *TrxRouter) error {
	// 验证路由配置
	if err := router.TrxRouterValues.Validate(); err != nil {
		return err
	}

	return DB.Create(router).Error
}

// UpdateTrxRouter 更新交易路由
func UpdateTrxRouter(router *TrxRouter) error {
	// 验证路由配置
	if err := router.TrxRouterValues.Validate(); err != nil {
		return err
	}

	return DB.Save(router).Error
}

// NewTrxRouter 创建新的交易路由实例
func NewTrxRouter() *TrxRouter {
	status := "active"
	priority := 100

	return &TrxRouter{
		RouterID: generateRouterID(),
		TrxRouterValues: &TrxRouterValues{
			Status:       &status,
			Priority:     &priority,
			CashierTypes: StringArray{},
			CashierIDs:   StringArray{},
		},
	}
}

// generateRouterID 生成路由ID
func generateRouterID() string {
	return fmt.Sprintf("TXR%d", getCurrentTimeMillis())
}

// =============================================================================
// 路由匹配引擎相关类型
// =============================================================================

// TransactionRequest 交易请求
type TransactionRequest struct {
	MerchantID      string          `json:"merchant_id"`
	TransactionType string          `json:"transaction_type"` // receipt, payment
	PaymentMethod   string          `json:"payment_method"`
	Amount          decimal.Decimal `json:"amount"`
	Currency        string          `json:"currency"`
	Country         string          `json:"country"`
	BillID          string          `json:"bill_id"`

	// 业务偏好设置
	PreferredTypes []string `json:"preferred_types"` // 可选，指定偏好类型
}

// RouterResult 路由结果
type RouterResult struct {
	Router   *TrxRouter             `json:"router"`
	Cashier  *CashierInfo           `json:"cashier"`
	Metadata map[string]interface{} `json:"metadata"`
}

// CashierInfo Cashier信息
type CashierInfo struct {
	CashierID   string                 `json:"cashier_id"`
	CashierType string                 `json:"cashier_type"` // business, personal
	Group       string                 `json:"group"`
	Status      string                 `json:"status"`
	Config      map[string]interface{} `json:"config"`
}

// RouterMatcher 路由匹配器
type RouterMatcher struct {
	routers []*TrxRouter
}

// NewRouterMatcher 创建路由匹配器
func NewRouterMatcher() *RouterMatcher {
	return &RouterMatcher{}
}

// LoadRouters 加载路由配置
func (rm *RouterMatcher) LoadRouters() error {
	var routers []*TrxRouter
	err := DB.Where("status = ?", "active").Find(&routers).Error
	if err != nil {
		return err
	}
	rm.routers = routers
	return nil
}

// FindBestRouter 查找最佳路由
func (rm *RouterMatcher) FindBestRouter(req *TransactionRequest) (*RouterResult, error) {
	// 1. 基础条件过滤
	matchedRouters := rm.filterByBasicConditions(req)

	// 2. 按优先级排序
	sort.Slice(matchedRouters, func(i, j int) bool {
		return matchedRouters[i].GetPriority() > matchedRouters[j].GetPriority()
	})

	// 3. 逐一验证可用性
	for _, router := range matchedRouters {
		cashierInfo, err := rm.selectAvailableCashier(router, req)
		if err != nil {
			continue
		}

		if cashierInfo != nil {
			return &RouterResult{
				Router:  router,
				Cashier: cashierInfo,
				Metadata: map[string]interface{}{
					"selected_at":      time.Now().Unix(),
					"selection_reason": "best_match",
				},
			}, nil
		}
	}

	return nil, errors.New("no available route found")
}

// filterByBasicConditions 基础条件过滤
func (rm *RouterMatcher) filterByBasicConditions(req *TransactionRequest) []*TrxRouter {
	var matchedRouters []*TrxRouter

	for _, router := range rm.routers {
		if router.TrxRouterValues.MatchesRequest(
			req.TransactionType,
			req.PaymentMethod,
			req.Currency,
			req.Country,
			req.Amount,
		) {
			matchedRouters = append(matchedRouters, router)
		}
	}

	return matchedRouters
}

// selectAvailableCashier 选择可用的Cashier
func (rm *RouterMatcher) selectAvailableCashier(router *TrxRouter, req *TransactionRequest) (*CashierInfo, error) {
	selector := router.GetCashierSelector()

	switch selector["type"] {
	case "cashier_types":
		types := selector["values"].([]string)
		return rm.selectCashierByTypes(types, req)

	case "cashier_ids":
		ids := selector["values"].([]string)
		return rm.selectCashierByIDs(ids, req)

	case "cashier_group":
		group := selector["values"].(string)
		return rm.selectCashierByGroup(group, req)

	default:
		return nil, errors.New("invalid cashier selector")
	}
}

// selectCashierByTypes 按类型选择Cashier
func (rm *RouterMatcher) selectCashierByTypes(types []string, req *TransactionRequest) (*CashierInfo, error) {
	// 这里应该调用实际的Cashier服务获取可用Cashier
	// 暂时返回模拟数据
	for _, cashierType := range types {
		if rm.isCashierTypeAvailable(cashierType, req) {
			return &CashierInfo{
				CashierID:   fmt.Sprintf("CASH_%s_%d", strings.ToUpper(cashierType), time.Now().Unix()%1000),
				CashierType: cashierType,
				Status:      "active",
				Config: map[string]interface{}{
					"selected_by": "type",
					"available":   true,
				},
			}, nil
		}
	}
	return nil, errors.New("no available cashier of specified types")
}

// selectCashierByIDs 按ID选择Cashier
func (rm *RouterMatcher) selectCashierByIDs(ids []string, req *TransactionRequest) (*CashierInfo, error) {
	// 这里应该调用实际的Cashier服务检查ID可用性
	// 暂时返回模拟数据
	for _, cashierID := range ids {
		if rm.isCashierIDAvailable(cashierID, req) {
			return &CashierInfo{
				CashierID:   cashierID,
				CashierType: "specific",
				Status:      "active",
				Config: map[string]interface{}{
					"selected_by": "id",
					"available":   true,
				},
			}, nil
		}
	}
	return nil, errors.New("no available cashier of specified IDs")
}

// selectCashierByGroup 按分组选择Cashier
func (rm *RouterMatcher) selectCashierByGroup(group string, req *TransactionRequest) (*CashierInfo, error) {
	// 这里应该调用实际的Cashier服务获取分组内可用Cashier
	// 暂时返回模拟数据
	if rm.isCashierGroupAvailable(group, req) {
		return &CashierInfo{
			CashierID:   fmt.Sprintf("CASH_%s_%d", strings.ToUpper(group), time.Now().Unix()%1000),
			CashierType: "group",
			Group:       group,
			Status:      "active",
			Config: map[string]interface{}{
				"selected_by": "group",
				"available":   true,
			},
		}, nil
	}
	return nil, errors.New("no available cashier in specified group")
}

// 可用性检查方法（模拟实现）
func (rm *RouterMatcher) isCashierTypeAvailable(cashierType string, req *TransactionRequest) bool {
	// 这里应该实现真实的可用性检查逻辑
	return true
}

func (rm *RouterMatcher) isCashierIDAvailable(cashierID string, req *TransactionRequest) bool {
	// 这里应该实现真实的ID可用性检查逻辑
	return true
}

func (rm *RouterMatcher) isCashierGroupAvailable(group string, req *TransactionRequest) bool {
	// 这里应该实现真实的分组可用性检查逻辑
	return true
}
