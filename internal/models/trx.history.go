package models

import (
	"encoding/json"
	"inpayos/internal/protocol"
)

// TrxHistory 交易历史记录，用于记录交易状态变更
type TrxHistory struct {
	ID         int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	HistoryID  string `json:"history_id" gorm:"column:history_id"`                      // 历史记录ID
	TrxID      string `json:"trx_id" gorm:"column:trx_id;index"`                        // 交易ID
	TrxType    string `json:"trx_type" gorm:"column:trx_type;index"`                    // 交易类型：payment, payout, refund等
	FromStatus string `json:"from_status" gorm:"column:from_status"`                    // 前状态
	ToStatus   string `json:"to_status" gorm:"column:to_status"`                        // 后状态
	ChangedBy  string `json:"changed_by" gorm:"column:changed_by"`                      // 修改人
	BeforeData string `json:"before_data" gorm:"column:before_data;type:text"`          // 修改前JSON数据
	AfterData  string `json:"after_data" gorm:"column:after_data;type:text"`            // 修改后JSON数据
	Remark     string `json:"remark" gorm:"column:remark"`                              // 备注信息
	CreatedAt  int64  `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"` // 创建时间
	UpdatedAt  int64  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

func NewTrxHistory() *TrxHistory {
	return &TrxHistory{}
}

func NewTrxHistoryByTransaction(trx *Transaction) *TrxHistory {
	history := NewTrxHistory()
	history.FillTransaction(trx)
	return history
}

func (*TrxHistory) TableName() string {
	return "t_trx_history"
}

func (trx *TrxHistory) FillTransaction(p *Transaction) {
	if p == nil {
		return
	}
	trx.TrxID = p.TrxID
	trx.TrxType = protocol.TrxTypePayin
	_json, _ := json.Marshal(p.TransactionValues)
	trx.FromStatus = p.GetStatus()
	trx.ToStatus = p.GetStatus()
	trx.BeforeData = string(_json)
	trx.AfterData = string(_json)
}
func (trx *TrxHistory) FillValues(v *TransactionValues) {
	if v == nil {
		return
	}
	_json, _ := json.Marshal(v)
	trx.AfterData = string(_json)
	trx.ToStatus = v.GetStatus()
}

// CreateHistory 保存支付交易历史
func CreateHistory(history *TrxHistory) error {
	return WriteDB.Create(&history).Error
}

// GetPaymentTrxHistory 获取交易历史记录
func GetPaymentTrxHistory(trxID string) ([]*TrxHistory, error) {
	var histories []*TrxHistory
	err := ReadDB.Where("trx_id = ? AND trx_type = ?", trxID, protocol.TrxTypePayin).Order("created_at desc").Find(&histories).Error
	return histories, err
}
