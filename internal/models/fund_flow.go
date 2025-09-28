package models

import (
	"github.com/shopspring/decimal"
)

// FundFlow 资金流水记录表
type FundFlow struct {
	ID     uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	FlowID string `json:"flow_id" gorm:"column:flow_id;type:varchar(64);uniqueIndex"`
	Salt   string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*FundFlowValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type FundFlowValues struct {
	UserID        *string          `json:"user_id" gorm:"column:user_id;type:varchar(32);index"`
	UserType      *string          `json:"user_type" gorm:"column:user_type;type:varchar(16);index"`
	AccountID     *string          `json:"account_id" gorm:"column:account_id;type:varchar(32);index"`
	TransactionID *string          `json:"transaction_id" gorm:"column:transaction_id;type:varchar(64);index"`
	BillID        *string          `json:"bill_id" gorm:"column:bill_id;type:varchar(64);index"`
	FlowType      *string          `json:"flow_type" gorm:"column:flow_type;type:varchar(16);index"` // income, expense, freeze, unfreeze, margin
	Amount        *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Currency      *string          `json:"currency" gorm:"column:currency;type:varchar(16)"`
	BeforeBalance *decimal.Decimal `json:"before_balance" gorm:"column:before_balance;type:decimal(36,18)"`
	AfterBalance  *decimal.Decimal `json:"after_balance" gorm:"column:after_balance;type:decimal(36,18)"`
	BusinessType  *string          `json:"business_type" gorm:"column:business_type;type:varchar(32)"` // payment, receipt, refund, settlement
	Description   *string          `json:"description" gorm:"column:description;type:varchar(512)"`
	FlowAt        *int64           `json:"flow_at" gorm:"column:flow_at;index"`

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (FundFlow) TableName() string {
	return "t_fund_flows"
}

// NewFundFlow 创建新资金流水
func NewFundFlow() *FundFlow {
	return &FundFlow{
		FundFlowValues: &FundFlowValues{},
	}
}

// SetValues 设置FundFlowValues
func (f *FundFlowValues) SetValues(values *FundFlowValues) {
	if values == nil {
		return
	}
	if values.UserID != nil {
		f.UserID = values.UserID
	}
	if values.UserType != nil {
		f.UserType = values.UserType
	}
	if values.AccountID != nil {
		f.AccountID = values.AccountID
	}
	if values.TransactionID != nil {
		f.TransactionID = values.TransactionID
	}
	if values.BillID != nil {
		f.BillID = values.BillID
	}
	if values.FlowType != nil {
		f.FlowType = values.FlowType
	}
	if values.Amount != nil {
		f.Amount = values.Amount
	}
	if values.Currency != nil {
		f.Currency = values.Currency
	}
	if values.BeforeBalance != nil {
		f.BeforeBalance = values.BeforeBalance
	}
	if values.AfterBalance != nil {
		f.AfterBalance = values.AfterBalance
	}
	if values.BusinessType != nil {
		f.BusinessType = values.BusinessType
	}
	if values.Description != nil {
		f.Description = values.Description
	}
	if values.FlowAt != nil {
		f.FlowAt = values.FlowAt
	}
}

// Getter方法
func (f *FundFlowValues) GetUserID() string {
	if f.UserID == nil {
		return ""
	}
	return *f.UserID
}

func (f *FundFlowValues) GetFlowType() string {
	if f.FlowType == nil {
		return ""
	}
	return *f.FlowType
}

func (f *FundFlowValues) GetAmount() decimal.Decimal {
	if f.Amount == nil {
		return decimal.Zero
	}
	return *f.Amount
}

func (f *FundFlowValues) GetCurrency() string {
	if f.Currency == nil {
		return ""
	}
	return *f.Currency
}

func (f *FundFlowValues) GetBeforeBalance() decimal.Decimal {
	if f.BeforeBalance == nil {
		return decimal.Zero
	}
	return *f.BeforeBalance
}

func (f *FundFlowValues) GetAfterBalance() decimal.Decimal {
	if f.AfterBalance == nil {
		return decimal.Zero
	}
	return *f.AfterBalance
}

// Setter方法(支持链式调用)
func (f *FundFlowValues) SetUserID(userID string) *FundFlowValues {
	f.UserID = &userID
	return f
}

func (f *FundFlowValues) SetUserType(userType string) *FundFlowValues {
	f.UserType = &userType
	return f
}

func (f *FundFlowValues) SetAccountID(accountID string) *FundFlowValues {
	f.AccountID = &accountID
	return f
}

func (f *FundFlowValues) SetTransactionID(transactionID string) *FundFlowValues {
	f.TransactionID = &transactionID
	return f
}

func (f *FundFlowValues) SetBillID(billID string) *FundFlowValues {
	f.BillID = &billID
	return f
}

func (f *FundFlowValues) SetFlowType(flowType string) *FundFlowValues {
	f.FlowType = &flowType
	return f
}

func (f *FundFlowValues) SetAmount(amount decimal.Decimal) *FundFlowValues {
	f.Amount = &amount
	return f
}

func (f *FundFlowValues) SetCurrency(currency string) *FundFlowValues {
	f.Currency = &currency
	return f
}

func (f *FundFlowValues) SetBeforeBalance(beforeBalance decimal.Decimal) *FundFlowValues {
	f.BeforeBalance = &beforeBalance
	return f
}

func (f *FundFlowValues) SetAfterBalance(afterBalance decimal.Decimal) *FundFlowValues {
	f.AfterBalance = &afterBalance
	return f
}

func (f *FundFlowValues) SetBusinessType(businessType string) *FundFlowValues {
	f.BusinessType = &businessType
	return f
}

func (f *FundFlowValues) SetDescription(description string) *FundFlowValues {
	f.Description = &description
	return f
}

func (f *FundFlowValues) SetFlowAt(flowAt int64) *FundFlowValues {
	f.FlowAt = &flowAt
	return f
}

// 查询方法
func GetFundFlowByID(id uint64) (*FundFlow, error) {
	var fundFlow FundFlow
	err := DB.Where("id = ?", id).First(&fundFlow).Error
	if err != nil {
		return nil, err
	}
	return &fundFlow, nil
}

func GetFundFlowByFlowID(flowID string) (*FundFlow, error) {
	var fundFlow FundFlow
	err := DB.Where("flow_id = ?", flowID).First(&fundFlow).Error
	if err != nil {
		return nil, err
	}
	return &fundFlow, nil
}

func GetFundFlowsByUserID(userID, userType string, limit, offset int) ([]*FundFlow, error) {
	var fundFlows []*FundFlow
	err := DB.Where("user_id = ? AND user_type = ?", userID, userType).
		Order("flow_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&fundFlows).Error
	if err != nil {
		return nil, err
	}
	return fundFlows, nil
}

func GetFundFlowsByTransactionID(transactionID string) ([]*FundFlow, error) {
	var fundFlows []*FundFlow
	err := DB.Where("transaction_id = ?", transactionID).
		Order("created_at ASC").
		Find(&fundFlows).Error
	if err != nil {
		return nil, err
	}
	return fundFlows, nil
}
