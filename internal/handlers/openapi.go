package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OpenApi 第一层：纯开放 API 接口处理器
// 面向第三方开发者集成，需要 API Key 认证
type OpenApi struct {
}

// NewOpenApi 创建 OpenApi 处理器
func NewOpenApi() *OpenApi {
	return &OpenApi{}
}

// =============================================================================
// 代收（Receipt）接口
// =============================================================================

// CreateReceipt 创建代收订单
func (a *OpenApi) CreateReceipt(c *gin.Context) {
	var req protocol.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	req.Type = "receipt" // 设置交易类型为代收
	transactionService := services.GetTransactionService()
	code, receipt := transactionService.CreateReceipt(&req)

	lang := middleware.GetLanguage(c)
	result := protocol.HandleServiceResult(code, receipt, lang)
	c.JSON(http.StatusOK, result)
}

// CancelReceipt 取消代收订单
func (a *OpenApi) CancelReceipt(c *gin.Context) {
	var req struct {
		BillID string `json:"bill_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// TODO: 实现取消代收逻辑
	lang := middleware.GetLanguage(c)
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(gin.H{"message": "Receipt cancelled"}, lang))
}

// =============================================================================
// 代付（Payment）接口
// =============================================================================

// CreatePayment 创建代付订单
func (a *OpenApi) CreatePayment(c *gin.Context) {
	var req protocol.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	req.Type = "payment" // 设置交易类型为代付
	transactionService := services.GetTransactionService()
	payment, err := transactionService.CreatePayment(&req)

	lang := middleware.GetLanguage(c)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.DatabaseError, lang))
		return
	}
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(payment, lang))
}

// CancelPayment 取消代付订单
func (a *OpenApi) CancelPayment(c *gin.Context) {
	var req struct {
		BillID string `json:"bill_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// TODO: 实现取消代付逻辑
	lang := middleware.GetLanguage(c)
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(gin.H{"message": "Payment cancelled"}, lang))
}

// =============================================================================
// 收银台（Checkout）接口
// =============================================================================

// CreateCheckout 创建收银台会话
func (a *OpenApi) CreateCheckout(c *gin.Context) {
	var req struct {
		BillID        string `json:"bill_id" binding:"required"`
		Amount        string `json:"amount" binding:"required"`
		Currency      string `json:"currency" binding:"required"`
		Country       string `json:"country" binding:"required"`
		PaymentMethod string `json:"payment_method" binding:"required"`
		ReturnURL     string `json:"return_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// TODO: 实现创建收银台逻辑
	lang := middleware.GetLanguage(c)
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(gin.H{
		"checkout_id":  "checkout_" + req.BillID,
		"checkout_url": "https://checkout.inpayos.com/" + req.BillID,
		"status":       "created",
	}, lang))
}

// GetCheckout 获取收银台信息
func (a *OpenApi) GetCheckout(c *gin.Context) {
	checkoutID := c.Param("id")
	if checkoutID == "" {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.MissingParams, lang))
		return
	}

	// TODO: 实现获取收银台信息逻辑
	lang := middleware.GetLanguage(c)
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(gin.H{
		"checkout_id": checkoutID,
		"status":      "pending",
	}, lang))
}

// GetCheckoutConfigs 获取收银台服务配置
func (a *OpenApi) GetCheckoutConfigs(c *gin.Context) {
	var req struct {
		Country  string `json:"country"`
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// TODO: 实现获取收银台配置逻辑
	lang := middleware.GetLanguage(c)
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(gin.H{
		"configs": []gin.H{
			{"method": "bank_transfer", "fee": "0.5%"},
			{"method": "upi", "fee": "0.3%"},
		},
	}, lang))
}

// ConfirmCheckout 确认收银台支付
func (a *OpenApi) ConfirmCheckout(c *gin.Context) {
	var req struct {
		CheckoutID string `json:"checkout_id" binding:"required"`
		PayProofID string `json:"pay_proof_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// TODO: 实现确认收银台支付逻辑
	lang := middleware.GetLanguage(c)
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(gin.H{"message": "Checkout confirmed"}, lang))
}

// CancelCheckout 取消收银台会话
func (a *OpenApi) CancelCheckout(c *gin.Context) {
	var req struct {
		CheckoutID string `json:"checkout_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// TODO: 实现取消收银台逻辑
	lang := middleware.GetLanguage(c)
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(gin.H{"message": "Checkout cancelled"}, lang))
}

// =============================================================================
// 统一查询接口
// =============================================================================

// QueryTransaction 查询交易状态/详情
func (a *OpenApi) QueryTransaction(c *gin.Context) {
	var req struct {
		BillID   string `json:"bill_id"`
		RecordID string `json:"record_id"`
		TxType   string `json:"tx_type"` // receipt, payment, checkout
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	if req.BillID == "" && req.RecordID == "" {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.MissingParams, lang))
		return
	}

	// TODO: 实现查询交易逻辑
	lang := middleware.GetLanguage(c)
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(gin.H{
		"bill_id":   req.BillID,
		"record_id": req.RecordID,
		"status":    "completed",
		"type":      req.TxType,
	}, lang))
}

// ListTransactions 交易列表查询
func (a *OpenApi) ListTransactions(c *gin.Context) {
	var req struct {
		UserID    string `json:"user_id"`
		Type      string `json:"type"`
		Status    string `json:"status"`
		StartTime int64  `json:"start_time"`
		EndTime   int64  `json:"end_time"`
		Page      int    `json:"page"`
		Size      int    `json:"size"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}

	// TODO: 实现交易列表查询逻辑
	lang := middleware.GetLanguage(c)
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(gin.H{
		"records": []gin.H{},
		"total":   0,
		"page":    req.Page,
		"size":    req.Size,
	}, lang))
}

// =============================================================================
// 余额查询
// =============================================================================

// QueryBalance 查询账户余额
func (a *OpenApi) QueryBalance(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id" binding:"required"`
		UserType string `json:"user_type" binding:"required"`
		Currency string `json:"currency" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	accountService := services.GetAccountService()
	balance, err := accountService.GetBalance(req.UserID, req.UserType, req.Currency)
	lang := middleware.GetLanguage(c)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.BalanceNotFound, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(balance, lang))
}
