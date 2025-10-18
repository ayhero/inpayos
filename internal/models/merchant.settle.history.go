package models

// MerchantSettleHistory 商户结算历史记录，用于记录商户的结算状态变更
type MerchantSettleHistory struct {
	ID         int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	HistoryID  string `json:"history_id" gorm:"column:history_id;uniqueIndex"` // 历史记录ID
	SettleID   string `json:"settle_id" gorm:"column:settle_id;index"`         // 结算ID
	Mid        string `json:"mid" gorm:"column:mid;type:varchar(64);index"`    // 商户ID
	FromStatus int    `json:"from_status" gorm:"column:from_status"`           // 前状态
	ToStatus   int    `json:"to_status" gorm:"column:to_status"`               // 后状态
	ChangedBy  string `json:"changed_by" gorm:"column:changed_by"`             // 修改人
	BeforeData string `json:"before_data" gorm:"column:before_data;type:text"` // 修改前JSON数据
	AfterData  string `json:"after_data" gorm:"column:after_data;type:text"`   // 修改后JSON数据
	Remark     string `json:"remark" gorm:"column:remark"`                     // 备注信息
	CreatedAt  int64  `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt  int64  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

func (t MerchantSettleHistory) TableName() string {
	return "t_merchant_settle_history"
}
