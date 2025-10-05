package services

import (
	"errors"
	"inpayos/internal/models"
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
)

type RouterRequest struct {
	Mid       string           `json:"mid"`
	TrxType   string           `json:"trx_type"`
	ReqID     string           `json:"req_id"`
	Ccy       string           `json:"ccy"`
	Amount    *decimal.Decimal `json:"amount"`
	TrxMethod string           `json:"trx_method"`
	TrxMode   string           `json:"trx_mode"`
	TrxApp    string           `json:"trx_app"`
	Pkg       string           `json:"pkg"`
	Did       string           `json:"did"`
	ProductID string           `json:"product_id"`
}

func GetChannelByMerchant(req *RouterRequest) (r *protocol.RouterInfo) {
	routers := models.ListRouterByMerchantByProirity(req.Mid, req.TrxType)
	if len(routers) == 0 {
		return
	}
	r = &protocol.RouterInfo{
		Mid:             req.Mid,
		ChannelAccounts: []string{},
		ChannelCodeLib:  map[string]string{},
	}
	for _, router := range routers {
		if err := ValidateRouter(router, req); err != nil {
			continue
		}
		if router.ChannelAccount != nil && *router.ChannelAccount != "" {
			account := models.GetChannelAccountsByAccountID(*router.ChannelAccount)
			if account != nil {
				r.ChannelAccounts = append(r.ChannelAccounts, account.GetAccountID())
				r.ChannelCodeLib[account.GetAccountID()] = account.ChannelCode
				return
			}
		}
		if router.ChannelCode != nil && *router.ChannelCode != "" {
			account := models.GetActiveChannelAccountByCode(req.Mid, *router.ChannelCode)
			if account != nil {
				r.ChannelAccounts = append(r.ChannelAccounts, account.GetAccountID())
				r.ChannelCodeLib[account.GetAccountID()] = account.ChannelCode
				return
			}
		}
		if router.ChannelGroup != nil && *router.ChannelGroup != "" {
			group := models.GetActiveChannelGroupByCode(*router.ChannelGroup)
			if group != nil {
				accounts := []string{}
				for _, member := range group.Members {
					account := models.GetChannelAccountsByAccountID(member.Member)
					if account != nil {
						accounts = append(accounts, account.GetAccountID())
						r.ChannelCodeLib[account.GetAccountID()] = account.ChannelCode
					}
				}
				if len(accounts) > 0 {
					r.ChannelAccounts = accounts
					r.Strategy = protocol.RouterStrategyAll
					if group.Setting != nil && group.Setting.Strategy != "" {
						r.Strategy = group.Setting.Strategy
					}
					return
				}
			}
		}
	}
	return nil
}

func ValidateRouter(router *models.MerchantRouter, req *RouterRequest) (err error) {
	if err = ValidateAmount(router, req.Ccy, req.Amount); err != nil {
		return err
	}
	if router.TrxMethod != nil && *router.TrxMethod != "" && *router.TrxMethod != req.TrxMethod {
		err = errors.New("trx_method not supported")
	}
	if router.TrxMode != nil && *router.TrxMode != "" && *router.TrxMode != req.TrxMode {
		err = errors.New("trx_mode not supported")
	}
	if router.TrxApp != nil && *router.TrxApp != "" && *router.TrxApp != req.TrxApp {
		err = errors.New("trx_app not supported")
	}
	if router.Did != nil && *router.Did != "" && *router.Did != req.Did {
		err = errors.New("did not supported")
	}
	if router.Pkg != nil && *router.Pkg != "" && *router.Pkg != req.Pkg {
		err = errors.New("pkg not supported")
	}

	return
}

func ValidateAmount(router *models.MerchantRouter, ccy string, amount *decimal.Decimal) (err error) {
	if router.Ccy != nil && *router.Ccy != "" && *router.Ccy != ccy {
		err = errors.New("currency not supported")
		return
	}
	if router.MinAmount == nil && router.MaxAmount == nil {
		return
	}
	if amount == nil {
		err = errors.New("amount not supported")
		return
	}
	if router.MinAmount != nil && router.MinAmount.GreaterThan(*amount) {
		err = errors.New("MinAmount not supported")
		return
	}
	if router.MaxAmount != nil && router.MaxAmount.LessThan(*amount) {
		err = errors.New("MaxAmount not supported")
		return
	}
	return
}
