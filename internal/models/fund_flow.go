package models

import (
	"inpayos/internal/protocol"

	"github.com/shopspring/decimal"
)

// FundFlow 资金流水表
type FundFlow struct {
	ID             uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	FlowNo         string           `gorm:"type:varchar(64);uniqueIndex;not null" json:"flow_no"`              // 流水号
	OriFlowNo      string           `gorm:"type:varchar(64);index" json:"ori_flow_no"`                         // 原始流水号(关联流水记录)
	UserID         string           `gorm:"type:varchar(32);index;not null" json:"user_id"`                    // 用户ID
	UserType       string           `gorm:"type:varchar(16);index;not null" json:"user_type"`                  // 用户类型: merchant-商户, cashier-收款员, system-系统，cashier-team-收款员团队
	AccountID      string           `gorm:"type:varchar(32);index;not null" json:"account_id"`                 // 账户ID
	AccountVersion int64            `gorm:"not null" json:"account_version"`                                   // 账户版本号
	Direction      string           `gorm:"type:varchar(10);not null" json:"direction"`                        // 流水方向：1-进账 2-出账
	TrxID          string           `gorm:"type:varchar(64);index" json:"trx_id"`                              // 关联业务ID
	TrxType        string           `gorm:"type:varchar(20);not null" json:"trx_type"`                         // 业务类型：payin-代收,payout-代付，margin-保证金，withdraw-提现，fee-手续费，adjust-调账
	Ccy            string           `gorm:"type:varchar(16);not null" json:"ccy"`                              // 交易币种
	Amount         *decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"amount"`                         // 交易金额
	BeforeAsset    *Asset           `gorm:"column:before_asset;serializer:json;type:json" json:"before_asset"` // 资金信息(包含币种)
	AfterAsset     *Asset           `gorm:"column:after_asset;serializer:json;type:json" json:"after_asset"`   // 资金信息(包含币种)
	Remark         string           `gorm:"type:varchar(255)" json:"remark"`                                   // 备注
	OperatorId     string           `gorm:"type:varchar(32)" json:"operator_id"`                               // 操作人ID
	CreatedAt      int64            `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt      int64            `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

func (t FundFlow) TableName() string {
	return "t_fund_flows"
}

// Protocol 转换为协议格式
func (f *FundFlow) Protocol() *protocol.FundFlow {
	// 获取变动前后余额
	var beforeBalance, afterBalance decimal.Decimal
	if f.BeforeAsset != nil {
		beforeBalance = f.BeforeAsset.Balance
	}
	if f.AfterAsset != nil {
		afterBalance = f.AfterAsset.Balance
	}

	return &protocol.FundFlow{
		ID:             f.ID,
		FlowNo:         f.FlowNo,
		OriFlowNo:      f.OriFlowNo,
		UserID:         f.UserID,
		UserType:       f.UserType,
		AccountID:      f.AccountID,
		AccountVersion: f.AccountVersion,
		Direction:      f.Direction,
		TrxID:          f.TrxID,
		TrxType:        f.TrxType,
		Ccy:            f.Ccy,
		Amount:         *f.Amount,
		BeforeBalance:  beforeBalance,
		AfterBalance:   afterBalance,
		Remark:         f.Remark,
		OperatorId:     f.OperatorId,
		CreatedAt:      f.CreatedAt,
		UpdatedAt:      f.UpdatedAt,
	}
}
