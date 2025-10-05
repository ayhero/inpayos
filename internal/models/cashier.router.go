package models

import (
	"inpayos/internal/log"

	"github.com/shopspring/decimal"
)

type CashierRouter struct {
	ID  int64  `json:"id" gorm:"column:id;primaryKey;AUTO_INCREMENT"`
	Tid string `json:"tid" gorm:"column:tid"`
	*CashierRouterValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

type CashierRouterValues struct {
	Ccy          *string          `json:"ccy" gorm:"column:ccy"`
	Country      *string          `json:"country" gorm:"column:country"` // 国家代码
	MinAmount    *decimal.Decimal `json:"min_amount" gorm:"column:min_amount"`
	MaxAmount    *decimal.Decimal `json:"max_amount" gorm:"column:max_amount"`
	MinUsdAmount *decimal.Decimal `json:"min_usd_amount" gorm:"column:min_usd_amount"`
	MaxUsdAmount *decimal.Decimal `json:"max_usd_amount" gorm:"column:max_usd_amount"`
	CashierID    *string          `json:"cashier_id" gorm:"column:cashier_id"`
	CashierGroup *string          `json:"cashier_group" gorm:"column:cashier_group"`
	Priority     *int64           `json:"priority" gorm:"column:priority"`
	Status       *string          `json:"status" gorm:"column:status"`
	Version      *int64           `json:"version" gorm:"column:version"`
}

func (t *CashierRouter) TableName() string {
	return "t_cashier_routers"
}

type CashierRouters []*CashierRouter

func ListRouterByCashierByProirity(mid, trx_type string) (data CashierRouters) {
	err := ReadDB.Where("mid = ? and trx_type=?", mid, trx_type).
		Order("priority desc").
		Find(&data).Error
	if err != nil {
		log.Get().Errorf("ListRouterByCashier error: %v", err)
		return nil
	}
	return
}

// GetCcy returns the Ccy value
func (mrv *CashierRouterValues) GetCcy() string {
	if mrv.Ccy == nil {
		return ""
	}
	return *mrv.Ccy
}

// GetCountry returns the Country value
func (mrv *CashierRouterValues) GetCountry() string {
	if mrv.Country == nil {
		return ""
	}
	return *mrv.Country
}

// GetMinAmount returns the MinAmount value
func (mrv *CashierRouterValues) GetMinAmount() decimal.Decimal {
	if mrv.MinAmount == nil {
		return decimal.Zero
	}
	return *mrv.MinAmount
}

// GetMaxAmount returns the MaxAmount value
func (mrv *CashierRouterValues) GetMaxAmount() decimal.Decimal {
	if mrv.MaxAmount == nil {
		return decimal.Zero
	}
	return *mrv.MaxAmount
}

// GetMinUsdAmount returns the MinUsdAmount value
func (mrv *CashierRouterValues) GetMinUsdAmount() decimal.Decimal {
	if mrv.MinUsdAmount == nil {
		return decimal.Zero
	}
	return *mrv.MinUsdAmount
}

// GetMaxUsdAmount returns the MaxUsdAmount value
func (mrv *CashierRouterValues) GetMaxUsdAmount() decimal.Decimal {
	if mrv.MaxUsdAmount == nil {
		return decimal.Zero
	}
	return *mrv.MaxUsdAmount
}

// GetPriority returns the Priority value
func (mrv *CashierRouterValues) GetPriority() int64 {
	if mrv.Priority == nil {
		return 0
	}
	return *mrv.Priority
}

// GetStatus returns the Status value
func (mrv *CashierRouterValues) GetStatus() string {
	if mrv.Status == nil {
		return ""
	}
	return *mrv.Status
}

// GetVersion returns the Version value
func (mrv *CashierRouterValues) GetVersion() int64 {
	if mrv.Version == nil {
		return 0
	}
	return *mrv.Version
}

// SetCcy sets the Ccy value
func (mrv *CashierRouterValues) SetCcy(value string) *CashierRouterValues {
	mrv.Ccy = &value
	return mrv
}

// SetCountry sets the Country value
func (mrv *CashierRouterValues) SetCountry(value string) *CashierRouterValues {
	mrv.Country = &value
	return mrv
}

// SetMinAmount sets the MinAmount value
func (mrv *CashierRouterValues) SetMinAmount(value decimal.Decimal) *CashierRouterValues {
	mrv.MinAmount = &value
	return mrv
}

// SetMaxAmount sets the MaxAmount value
func (mrv *CashierRouterValues) SetMaxAmount(value decimal.Decimal) *CashierRouterValues {
	mrv.MaxAmount = &value
	return mrv
}

// SetMinUsdAmount sets the MinUsdAmount value
func (mrv *CashierRouterValues) SetMinUsdAmount(value decimal.Decimal) *CashierRouterValues {
	mrv.MinUsdAmount = &value
	return mrv
}

// SetMaxUsdAmount sets the MaxUsdAmount value
func (mrv *CashierRouterValues) SetMaxUsdAmount(value decimal.Decimal) *CashierRouterValues {
	mrv.MaxUsdAmount = &value
	return mrv
}

// SetPriority sets the Priority value
func (mrv *CashierRouterValues) SetPriority(value int64) *CashierRouterValues {
	mrv.Priority = &value
	return mrv
}

// SetStatus sets the Status value
func (mrv *CashierRouterValues) SetStatus(value string) *CashierRouterValues {
	mrv.Status = &value
	return mrv
}

// SetVersion sets the Version value
func (mrv *CashierRouterValues) SetVersion(value int64) *CashierRouterValues {
	mrv.Version = &value
	return mrv
}

// SetValues sets multiple CashierRouterValues fields at once
func (mr *CashierRouter) SetValues(values *CashierRouterValues) *CashierRouter {
	if values == nil {
		return mr
	}

	if mr.CashierRouterValues == nil {
		mr.CashierRouterValues = &CashierRouterValues{}
	}

	if values.Ccy != nil {
		mr.CashierRouterValues.SetCcy(*values.Ccy)
	}
	if values.Country != nil {
		mr.CashierRouterValues.SetCountry(*values.Country)
	}
	if values.MinAmount != nil {
		mr.CashierRouterValues.SetMinAmount(*values.MinAmount)
	}
	if values.MaxAmount != nil {
		mr.CashierRouterValues.SetMaxAmount(*values.MaxAmount)
	}
	if values.MinUsdAmount != nil {
		mr.CashierRouterValues.SetMinUsdAmount(*values.MinUsdAmount)
	}
	if values.MaxUsdAmount != nil {
		mr.CashierRouterValues.SetMaxUsdAmount(*values.MaxUsdAmount)
	}
	if values.Priority != nil {
		mr.CashierRouterValues.SetPriority(*values.Priority)
	}
	if values.Status != nil {
		mr.CashierRouterValues.SetStatus(*values.Status)
	}
	if values.Version != nil {
		mr.CashierRouterValues.SetVersion(*values.Version)
	}

	return mr
}
