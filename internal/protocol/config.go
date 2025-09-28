package protocol

// Config Protocol Definitions
type GetConfigRequest struct {
	MerchantID string `json:"merchant_id" form:"merchant_id"`
	Type       string `json:"type" form:"type" binding:"required,oneof=receipt payment deposit withdraw refund"`
}

type SaveConfigRequest struct {
	MerchantID    string            `json:"merchant_id" binding:"required"`
	Type          string            `json:"type" binding:"required,oneof=receipt payment deposit withdraw refund"`
	MinAmount     map[string]string `json:"min_amount,omitempty"`
	MaxAmount     map[string]string `json:"max_amount,omitempty"`
	DailyLimit    map[string]string `json:"daily_limit,omitempty"`
	MonthlyLimit  map[string]string `json:"monthly_limit,omitempty"`
	FeeRate       map[string]string `json:"fee_rate,omitempty"`
	FeeFixed      map[string]string `json:"fee_fixed,omitempty"`
	Status        string            `json:"status,omitempty" binding:"omitempty,oneof=on off"`
	AutoConfirm   string            `json:"auto_confirm,omitempty" binding:"omitempty,oneof=on off"`
	NotifyURL     string            `json:"notify_url,omitempty"`
	TimeoutMinute int               `json:"timeout_minute,omitempty"`
}

type ConfigResponse struct {
	ID            uint64            `json:"id"`
	MerchantID    string            `json:"merchant_id"`
	Type          string            `json:"type"`
	MinAmount     map[string]string `json:"min_amount,omitempty"`
	MaxAmount     map[string]string `json:"max_amount,omitempty"`
	DailyLimit    map[string]string `json:"daily_limit,omitempty"`
	MonthlyLimit  map[string]string `json:"monthly_limit,omitempty"`
	FeeRate       map[string]string `json:"fee_rate,omitempty"`
	FeeFixed      map[string]string `json:"fee_fixed,omitempty"`
	Status        string            `json:"status"`
	AutoConfirm   string            `json:"auto_confirm"`
	NotifyURL     string            `json:"notify_url"`
	TimeoutMinute int               `json:"timeout_minute"`
	CreatedAt     int64             `json:"created_at"`
	UpdatedAt     int64             `json:"updated_at"`
}

type ListConfigsRequest struct {
	MerchantID string `json:"merchant_id" form:"merchant_id"`
}
