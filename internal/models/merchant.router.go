package models

import (
	"inpayos/internal/log"

	"github.com/shopspring/decimal"
)

type MerchantRouter struct {
	ID  int64  `json:"id" gorm:"column:id;primaryKey;AUTO_INCREMENT"`
	MID string `json:"mid" gorm:"column:mid"`
	*MerchantRouterValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

type MerchantRouterValues struct {
	Pkg            *string          `json:"pkg" gorm:"column:pkg"`
	Did            *string          `json:"did" gorm:"column:did"`
	TrxType        *string          `json:"trx_type" gorm:"column:trx_type"`
	TrxSubType     *string          `json:"trx_sub_type" gorm:"column:trx_sub_type"`
	TrxMethod      *string          `json:"trx_method" gorm:"column:trx_method"`
	TrxMode        *string          `json:"trx_mode" gorm:"column:trx_mode"`
	TrxApp         *string          `json:"trx_app" gorm:"column:trx_app"`
	Ccy            *string          `json:"ccy" gorm:"column:ccy"`
	Country        *string          `json:"country" gorm:"column:country"` // 国家代码
	MinAmount      *decimal.Decimal `json:"min_amount" gorm:"column:min_amount"`
	MaxAmount      *decimal.Decimal `json:"max_amount" gorm:"column:max_amount"`
	MinUsdAmount   *decimal.Decimal `json:"min_usd_amount" gorm:"column:min_usd_amount"`
	MaxUsdAmount   *decimal.Decimal `json:"max_usd_amount" gorm:"column:max_usd_amount"`
	ChannelCode    *string          `json:"channel_code" gorm:"column:channel_code"`
	ChannelAccount *string          `json:"channel_account" gorm:"column:channel_account"`
	ChannelGroup   *string          `json:"channel_group" gorm:"column:channel_group"`
	Priority       *int64           `json:"priority" gorm:"column:priority"`
	Status         *string          `json:"status" gorm:"column:status"`
	Version        *int64           `json:"version" gorm:"column:version"`
}

func (t *MerchantRouter) TableName() string {
	return "t_merchant_routers"
}

type MerchantRouters []*MerchantRouter

func ListRouterByMerchantByProirity(mid, trx_type string) (data MerchantRouters) {
	err := ReadDB.Where("mid = ? and trx_type=?", mid, trx_type).
		Order("priority desc").
		Find(&data).Error
	if err != nil {
		log.Get().Errorf("ListRouterByMerchant error: %v", err)
		return nil
	}
	return
}

// MerchantRouterValues Getter Methods
// GetPkg returns the Pkg value
func (mrv *MerchantRouterValues) GetPkg() string {
	if mrv.Pkg == nil {
		return ""
	}
	return *mrv.Pkg
}

// GetDid returns the Did value
func (mrv *MerchantRouterValues) GetDid() string {
	if mrv.Did == nil {
		return ""
	}
	return *mrv.Did
}

// GetTrxType returns the TrxType value
func (mrv *MerchantRouterValues) GetTrxType() string {
	if mrv.TrxType == nil {
		return ""
	}
	return *mrv.TrxType
}

// GetTrxSubType returns the TrxSubType value
func (mrv *MerchantRouterValues) GetTrxSubType() string {
	if mrv.TrxSubType == nil {
		return ""
	}
	return *mrv.TrxSubType
}

// GetTrxMethod returns the TrxMethod value
func (mrv *MerchantRouterValues) GetTrxMethod() string {
	if mrv.TrxMethod == nil {
		return ""
	}
	return *mrv.TrxMethod
}

// GetTrxMode returns the TrxMode value
func (mrv *MerchantRouterValues) GetTrxMode() string {
	if mrv.TrxMode == nil {
		return ""
	}
	return *mrv.TrxMode
}

// GetTrxApp returns the TrxApp value
func (mrv *MerchantRouterValues) GetTrxApp() string {
	if mrv.TrxApp == nil {
		return ""
	}
	return *mrv.TrxApp
}

// GetCcy returns the Ccy value
func (mrv *MerchantRouterValues) GetCcy() string {
	if mrv.Ccy == nil {
		return ""
	}
	return *mrv.Ccy
}

// GetCountry returns the Country value
func (mrv *MerchantRouterValues) GetCountry() string {
	if mrv.Country == nil {
		return ""
	}
	return *mrv.Country
}

// GetMinAmount returns the MinAmount value
func (mrv *MerchantRouterValues) GetMinAmount() decimal.Decimal {
	if mrv.MinAmount == nil {
		return decimal.Zero
	}
	return *mrv.MinAmount
}

// GetMaxAmount returns the MaxAmount value
func (mrv *MerchantRouterValues) GetMaxAmount() decimal.Decimal {
	if mrv.MaxAmount == nil {
		return decimal.Zero
	}
	return *mrv.MaxAmount
}

// GetMinUsdAmount returns the MinUsdAmount value
func (mrv *MerchantRouterValues) GetMinUsdAmount() decimal.Decimal {
	if mrv.MinUsdAmount == nil {
		return decimal.Zero
	}
	return *mrv.MinUsdAmount
}

// GetMaxUsdAmount returns the MaxUsdAmount value
func (mrv *MerchantRouterValues) GetMaxUsdAmount() decimal.Decimal {
	if mrv.MaxUsdAmount == nil {
		return decimal.Zero
	}
	return *mrv.MaxUsdAmount
}

// GetChannelCode returns the ChannelCode value
func (mrv *MerchantRouterValues) GetChannelCode() string {
	if mrv.ChannelCode == nil {
		return ""
	}
	return *mrv.ChannelCode
}

// GetChannelAccount returns the ChannelAccount value
func (mrv *MerchantRouterValues) GetChannelAccount() string {
	if mrv.ChannelAccount == nil {
		return ""
	}
	return *mrv.ChannelAccount
}

// GetChannelGroup returns the ChannelGroup value
func (mrv *MerchantRouterValues) GetChannelGroup() string {
	if mrv.ChannelGroup == nil {
		return ""
	}
	return *mrv.ChannelGroup
}

// GetPriority returns the Priority value
func (mrv *MerchantRouterValues) GetPriority() int64 {
	if mrv.Priority == nil {
		return 0
	}
	return *mrv.Priority
}

// GetStatus returns the Status value
func (mrv *MerchantRouterValues) GetStatus() string {
	if mrv.Status == nil {
		return ""
	}
	return *mrv.Status
}

// GetVersion returns the Version value
func (mrv *MerchantRouterValues) GetVersion() int64 {
	if mrv.Version == nil {
		return 0
	}
	return *mrv.Version
}

// MerchantRouterValues Setter Methods (support method chaining)
// SetPkg sets the Pkg value
func (mrv *MerchantRouterValues) SetPkg(value string) *MerchantRouterValues {
	mrv.Pkg = &value
	return mrv
}

// SetDid sets the Did value
func (mrv *MerchantRouterValues) SetDid(value string) *MerchantRouterValues {
	mrv.Did = &value
	return mrv
}

// SetTrxType sets the TrxType value
func (mrv *MerchantRouterValues) SetTrxType(value string) *MerchantRouterValues {
	mrv.TrxType = &value
	return mrv
}

// SetTrxSubType sets the TrxSubType value
func (mrv *MerchantRouterValues) SetTrxSubType(value string) *MerchantRouterValues {
	mrv.TrxSubType = &value
	return mrv
}

// SetTrxMethod sets the TrxMethod value
func (mrv *MerchantRouterValues) SetTrxMethod(value string) *MerchantRouterValues {
	mrv.TrxMethod = &value
	return mrv
}

// SetTrxMode sets the TrxMode value
func (mrv *MerchantRouterValues) SetTrxMode(value string) *MerchantRouterValues {
	mrv.TrxMode = &value
	return mrv
}

// SetTrxApp sets the TrxApp value
func (mrv *MerchantRouterValues) SetTrxApp(value string) *MerchantRouterValues {
	mrv.TrxApp = &value
	return mrv
}

// SetCcy sets the Ccy value
func (mrv *MerchantRouterValues) SetCcy(value string) *MerchantRouterValues {
	mrv.Ccy = &value
	return mrv
}

// SetCountry sets the Country value
func (mrv *MerchantRouterValues) SetCountry(value string) *MerchantRouterValues {
	mrv.Country = &value
	return mrv
}

// SetMinAmount sets the MinAmount value
func (mrv *MerchantRouterValues) SetMinAmount(value decimal.Decimal) *MerchantRouterValues {
	mrv.MinAmount = &value
	return mrv
}

// SetMaxAmount sets the MaxAmount value
func (mrv *MerchantRouterValues) SetMaxAmount(value decimal.Decimal) *MerchantRouterValues {
	mrv.MaxAmount = &value
	return mrv
}

// SetMinUsdAmount sets the MinUsdAmount value
func (mrv *MerchantRouterValues) SetMinUsdAmount(value decimal.Decimal) *MerchantRouterValues {
	mrv.MinUsdAmount = &value
	return mrv
}

// SetMaxUsdAmount sets the MaxUsdAmount value
func (mrv *MerchantRouterValues) SetMaxUsdAmount(value decimal.Decimal) *MerchantRouterValues {
	mrv.MaxUsdAmount = &value
	return mrv
}

// SetChannelCode sets the ChannelCode value
func (mrv *MerchantRouterValues) SetChannelCode(value string) *MerchantRouterValues {
	mrv.ChannelCode = &value
	return mrv
}

// SetChannelAccount sets the ChannelAccount value
func (mrv *MerchantRouterValues) SetChannelAccount(value string) *MerchantRouterValues {
	mrv.ChannelAccount = &value
	return mrv
}

// SetChannelGroup sets the ChannelGroup value
func (mrv *MerchantRouterValues) SetChannelGroup(value string) *MerchantRouterValues {
	mrv.ChannelGroup = &value
	return mrv
}

// SetPriority sets the Priority value
func (mrv *MerchantRouterValues) SetPriority(value int64) *MerchantRouterValues {
	mrv.Priority = &value
	return mrv
}

// SetStatus sets the Status value
func (mrv *MerchantRouterValues) SetStatus(value string) *MerchantRouterValues {
	mrv.Status = &value
	return mrv
}

// SetVersion sets the Version value
func (mrv *MerchantRouterValues) SetVersion(value int64) *MerchantRouterValues {
	mrv.Version = &value
	return mrv
}

// SetValues sets multiple MerchantRouterValues fields at once
func (mr *MerchantRouter) SetValues(values *MerchantRouterValues) *MerchantRouter {
	if values == nil {
		return mr
	}

	if mr.MerchantRouterValues == nil {
		mr.MerchantRouterValues = &MerchantRouterValues{}
	}

	// Set all fields from the provided values
	if values.Pkg != nil {
		mr.MerchantRouterValues.SetPkg(*values.Pkg)
	}
	if values.Did != nil {
		mr.MerchantRouterValues.SetDid(*values.Did)
	}
	if values.TrxType != nil {
		mr.MerchantRouterValues.SetTrxType(*values.TrxType)
	}
	if values.TrxSubType != nil {
		mr.MerchantRouterValues.SetTrxSubType(*values.TrxSubType)
	}
	if values.TrxMethod != nil {
		mr.MerchantRouterValues.SetTrxMethod(*values.TrxMethod)
	}
	if values.TrxMode != nil {
		mr.MerchantRouterValues.SetTrxMode(*values.TrxMode)
	}
	if values.TrxApp != nil {
		mr.MerchantRouterValues.SetTrxApp(*values.TrxApp)
	}
	if values.Ccy != nil {
		mr.MerchantRouterValues.SetCcy(*values.Ccy)
	}
	if values.Country != nil {
		mr.MerchantRouterValues.SetCountry(*values.Country)
	}
	if values.MinAmount != nil {
		mr.MerchantRouterValues.SetMinAmount(*values.MinAmount)
	}
	if values.MaxAmount != nil {
		mr.MerchantRouterValues.SetMaxAmount(*values.MaxAmount)
	}
	if values.MinUsdAmount != nil {
		mr.MerchantRouterValues.SetMinUsdAmount(*values.MinUsdAmount)
	}
	if values.MaxUsdAmount != nil {
		mr.MerchantRouterValues.SetMaxUsdAmount(*values.MaxUsdAmount)
	}
	if values.ChannelCode != nil {
		mr.MerchantRouterValues.SetChannelCode(*values.ChannelCode)
	}
	if values.ChannelAccount != nil {
		mr.MerchantRouterValues.SetChannelAccount(*values.ChannelAccount)
	}
	if values.ChannelGroup != nil {
		mr.MerchantRouterValues.SetChannelGroup(*values.ChannelGroup)
	}
	if values.Priority != nil {
		mr.MerchantRouterValues.SetPriority(*values.Priority)
	}
	if values.Status != nil {
		mr.MerchantRouterValues.SetStatus(*values.Status)
	}
	if values.Version != nil {
		mr.MerchantRouterValues.SetVersion(*values.Version)
	}

	return mr
}
