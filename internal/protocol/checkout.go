package protocol

type CheckoutInfoRequest struct {
	CheckoutID string `json:"checkout_id" binding:"required"`
}

type CreateCheckoutRequest struct {
	ReqID         string `json:"req_id" binding:"required"`
	Ccy           string `json:"ccy" binding:"required"`
	Amount        string `json:"amount" binding:"required"`
	Country       string `json:"country"`
	PaymentMethod string `json:"payment_method" binding:"required"`
	NotifyURL     string `json:"notify_url"`
	ReturnURL     string `json:"return_url"`
}
type ConfirmCheckoutRequest struct {
	CheckoutID string `json:"checkout_id" binding:"required"`
	Ccy        string `json:"ccy" binding:"required"`
	Amount     string `json:"amount" binding:"required"`
}

type CancelCheckoutRequest struct {
	CheckoutID string `json:"checkout_id" binding:"required"`
}

type Checkout struct {
	CheckoutID    string `json:"checkout_id"`
	Mid           string `json:"mid"`
	ReqID         string `json:"req_id"`
	Amount        string `json:"amount"`
	Ccy           string `json:"ccy"`
	Country       string `json:"country"`
	PaymentMethod string `json:"payment_method"`
	Status        string `json:"status"`
	NotifyURL     string `json:"notify_url,omitempty"`
	ReturnURL     string `json:"return_url,omitempty"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
	ExpiredAt     int64  `json:"expired_at"`
	FailureReason string `json:"failure_reason,omitempty"`
	Metadata      string `json:"metadata,omitempty"`
}
