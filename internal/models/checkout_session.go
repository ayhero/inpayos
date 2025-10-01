package models

// CheckoutSession 收银台会话表
type CheckoutSession struct {
	ID            uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	SessionID     string `gorm:"column:session_id;type:varchar(64);uniqueIndex;not null" json:"session_id"`
	MerchantID    string `gorm:"column:merchant_id;type:varchar(64);not null;index" json:"merchant_id"`
	BillID        string `gorm:"column:bill_id;type:varchar(64);index" json:"bill_id"`
	Amount        string `gorm:"column:amount;type:decimal(20,8);not null" json:"amount"`
	Currency      string `gorm:"column:currency;type:varchar(10);not null" json:"currency"`
	Country       string `gorm:"column:country;type:varchar(3)" json:"country"`
	PaymentMethod string `gorm:"column:payment_method;type:varchar(32)" json:"payment_method"`
	ReturnURL     string `gorm:"column:return_url;type:varchar(1024)" json:"return_url"`
	CancelURL     string `gorm:"column:cancel_url;type:varchar(1024)" json:"cancel_url"`
	NotifyURL     string `gorm:"column:notify_url;type:varchar(1024)" json:"notify_url"`
	Status        string `gorm:"column:status;type:varchar(32);not null;default:'created'" json:"status"` // created, pending, completed, cancelled, expired
	TransactionID string `gorm:"column:transaction_id;type:varchar(64);index" json:"transaction_id"`      // 关联的交易ID
	ChannelCode   string `gorm:"column:channel_code;type:varchar(32)" json:"channel_code"`
	CustomerInfo  string `gorm:"column:customer_info;type:text" json:"customer_info"` // JSON格式的客户信息
	PaymentInfo   string `gorm:"column:payment_info;type:text" json:"payment_info"`   // JSON格式的支付信息
	Metadata      string `gorm:"column:metadata;type:text" json:"metadata"`           // JSON格式的元数据
	ErrorCode     string `gorm:"column:error_code;type:varchar(32)" json:"error_code"`
	ErrorMsg      string `gorm:"column:error_msg;type:varchar(512)" json:"error_msg"`
	ExpiredAt     int64  `gorm:"column:expired_at;type:bigint" json:"expired_at"`
	CompletedAt   int64  `gorm:"column:completed_at;type:bigint" json:"completed_at"`
	CreatedAt     int64  `gorm:"column:created_at;type:bigint;autoCreateTime:milli" json:"created_at"`
	UpdatedAt     int64  `gorm:"column:updated_at;type:bigint;autoUpdateTime:milli" json:"updated_at"`
	DeletedAt     int64  `gorm:"column:deleted_at;type:bigint;index" json:"deleted_at,omitempty"`
}

// TableName 返回表名
func (CheckoutSession) TableName() string {
	return "t_checkout_sessions"
}

// IsExpired 检查会话是否过期
func (cs *CheckoutSession) IsExpired() bool {
	if cs.ExpiredAt == 0 {
		return false
	}
	return cs.ExpiredAt < getCurrentTimeMillis()
}

// IsCompleted 检查会话是否已完成
func (cs *CheckoutSession) IsCompleted() bool {
	return cs.Status == "completed"
}

// IsCancelled 检查会话是否已取消
func (cs *CheckoutSession) IsCancelled() bool {
	return cs.Status == "cancelled"
}

// CanCancel 检查是否可以取消
func (cs *CheckoutSession) CanCancel() bool {
	return cs.Status == "created" || cs.Status == "pending"
}

// CanComplete 检查是否可以完成
func (cs *CheckoutSession) CanComplete() bool {
	return cs.Status == "pending" && !cs.IsExpired()
}

// CheckoutSessionResponse 收银台会话响应结构
type CheckoutSessionResponse struct {
	ID            uint64 `json:"id"`
	SessionID     string `json:"session_id"`
	MerchantID    string `json:"merchant_id"`
	BillID        string `json:"bill_id"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	Country       string `json:"country"`
	PaymentMethod string `json:"payment_method"`
	ReturnURL     string `json:"return_url,omitempty"`
	CancelURL     string `json:"cancel_url,omitempty"`
	NotifyURL     string `json:"notify_url,omitempty"`
	Status        string `json:"status"`
	TransactionID string `json:"transaction_id,omitempty"`
	ChannelCode   string `json:"channel_code,omitempty"`
	CustomerInfo  string `json:"customer_info,omitempty"`
	PaymentInfo   string `json:"payment_info,omitempty"`
	Metadata      string `json:"metadata,omitempty"`
	ErrorCode     string `json:"error_code,omitempty"`
	ErrorMsg      string `json:"error_msg,omitempty"`
	ExpiredAt     int64  `json:"expired_at"`
	CompletedAt   int64  `json:"completed_at"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}

// ToResponse 转换为响应结构
func (cs *CheckoutSession) ToResponse() *CheckoutSessionResponse {
	return &CheckoutSessionResponse{
		ID:            cs.ID,
		SessionID:     cs.SessionID,
		MerchantID:    cs.MerchantID,
		BillID:        cs.BillID,
		Amount:        cs.Amount,
		Currency:      cs.Currency,
		Country:       cs.Country,
		PaymentMethod: cs.PaymentMethod,
		ReturnURL:     cs.ReturnURL,
		CancelURL:     cs.CancelURL,
		NotifyURL:     cs.NotifyURL,
		Status:        cs.Status,
		TransactionID: cs.TransactionID,
		ChannelCode:   cs.ChannelCode,
		CustomerInfo:  cs.CustomerInfo,
		PaymentInfo:   cs.PaymentInfo,
		Metadata:      cs.Metadata,
		ErrorCode:     cs.ErrorCode,
		ErrorMsg:      cs.ErrorMsg,
		ExpiredAt:     cs.ExpiredAt,
		CompletedAt:   cs.CompletedAt,
		CreatedAt:     cs.CreatedAt,
		UpdatedAt:     cs.UpdatedAt,
	}
}

// GetCheckoutSessionByID 根据SessionID获取收银台会话
func GetCheckoutSessionByID(sessionID string) (*CheckoutSession, error) {
	var session CheckoutSession
	err := DB.Where("session_id = ?", sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetCheckoutSessionByBillID 根据商户订单号获取收银台会话
func GetCheckoutSessionByBillID(merchantID, billID string) (*CheckoutSession, error) {
	var session CheckoutSession
	err := DB.Where("merchant_id = ? AND bill_id = ?", merchantID, billID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// UpdateCheckoutSessionStatus 更新收银台会话状态
func (cs *CheckoutSession) UpdateStatus(status string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": getCurrentTimeMillis(),
	}

	// 根据状态设置时间戳
	if status == "completed" {
		updates["completed_at"] = getCurrentTimeMillis()
	}

	return DB.Model(cs).Updates(updates).Error
}

// SetTransactionID 设置关联的交易ID
func (cs *CheckoutSession) SetTransactionID(transactionID string) error {
	return DB.Model(cs).Updates(map[string]interface{}{
		"transaction_id": transactionID,
		"updated_at":     getCurrentTimeMillis(),
	}).Error
}

// SetError 设置错误信息
func (cs *CheckoutSession) SetError(errorCode, errorMsg string) error {
	return DB.Model(cs).Updates(map[string]interface{}{
		"error_code": errorCode,
		"error_msg":  errorMsg,
		"updated_at": getCurrentTimeMillis(),
	}).Error
}

// 移除重复的getCurrentTimeMillis函数
