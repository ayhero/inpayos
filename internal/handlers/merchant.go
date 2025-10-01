package handlers

import (
	"inpayos/internal/config"
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Merchant 第二层：商户后台接口处理器
// 面向商户用户，需要 JWT + Session 认证
type Merchant struct {
	*config.ServiceConfig
}

// NewMerchant 创建 Merchant 处理器
func NewMerchant() *Merchant {
	cfg := config.Get()
	if cfg == nil || cfg.Server.Merchant == nil {
		return nil
	}

	return &Merchant{
		ServiceConfig: cfg.Server.Merchant,
	}
}

// ToServer 创建 HTTP服务器
func (m *Merchant) ToServer() *http.Server {
	return &http.Server{
		Addr: ":" + m.Port,
	}
}

// SetupRouter 设置路由
func (m *Merchant) SetupRouter() *gin.Engine {
	// 设置Gin模式
	cfg := config.Get()
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// 添加中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.LanguageMiddleware())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "merchant"})
	})

	// 商户认证相关（无前缀）
	router.POST("/login", m.Login)

	// 需要JWT认证的端点
	merchantAPI := router.Group("/")
	merchantAPI.Use(middleware.JWTAuth())
	{
		merchantAPI.POST("/logout", m.Logout)
		merchantAPI.POST("/register", m.Register)
		merchantAPI.POST("/verify", m.Verify)
		merchantAPI.POST("/reset-password", m.ResetPassword)
		merchantAPI.POST("/change-password", m.ChangePassword)

		// 账户信息管理
		account := merchantAPI.Group("/account")
		{
			account.GET("/info", m.GetAccountInfo)
			account.PUT("/info", m.UpdateAccountInfo)
			account.GET("/balance", m.GetAccountBalance)
			account.GET("/history", m.GetAccountHistory)
		}

		// API配置管理
		api := merchantAPI.Group("/api")
		{
			api.POST("/public-key", m.SubmitPublicKey)
			api.POST("/webhook", m.SubmitWebhook)
			api.POST("/whitelist", m.SubmitWhitelist)
		}

		// 交易管理
		transactions := merchantAPI.Group("/transactions")
		{
			transactions.GET("", m.ListTransactions)
			transactions.GET("/:id", m.GetTransactionDetail)
			transactions.GET("/export", m.ExportTransactions)
			transactions.POST("/refund", m.CreateRefund)
		}

		// 收银台管理
		checkouts := merchantAPI.Group("/checkouts")
		{
			checkouts.GET("", m.ListCheckouts)
			checkouts.GET("/services", m.GetCheckoutServices)
		}
	}

	return router
}

// =============================================================================
// 商户认证管理
// =============================================================================

// Login 商户登录
func (m *Merchant) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Invalid login parameters: "+err.Error()))
		return
	}

	// TODO: 实现商户登录逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(gin.H{
		"token":      "jwt_token_here",
		"expires_in": 3600,
	}))
}

// Logout 商户登出
func (m *Merchant) Logout(c *gin.Context) {
	// TODO: 实现商户登出逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(nil))
}

// Register 商户注册
func (m *Merchant) Register(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=8"`
		CompanyName string `json:"company_name" binding:"required"`
		ContactName string `json:"contact_name" binding:"required"`
		Phone       string `json:"phone" binding:"required"`
		Country     string `json:"country" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Invalid registration parameters: "+err.Error()))
		return
	}

	// TODO: 实现商户注册逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(gin.H{
		"merchant_id": "merchant_id_here",
		"status":      "pending_approval",
	}))
}

// Verify 账户验证
func (m *Merchant) Verify(c *gin.Context) {
	var req struct {
		VerifyCode string `json:"verify_code" binding:"required"`
		Type       string `json:"type" binding:"required,oneof=email sms g2fa"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Invalid verify parameters: "+err.Error()))
		return
	}

	// TODO: 实现账户验证逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(nil))
}

// ResetPassword 重置密码
func (m *Merchant) ResetPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Invalid email: "+err.Error()))
		return
	}

	// TODO: 实现重置密码逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(nil))
}

// ChangePassword 修改密码
func (m *Merchant) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Invalid password parameters: "+err.Error()))
		return
	}

	// TODO: 实现修改密码逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(nil))
}

// =============================================================================
// 账户信息管理
// =============================================================================

// GetAccountInfo 获取账户信息
func (m *Merchant) GetAccountInfo(c *gin.Context) {
	// TODO: 从JWT token中获取用户ID
	userID := "current_user_id"

	// TODO: 实现获取账户信息逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(gin.H{
		"user_id":      userID,
		"user_type":    "merchant",
		"company_name": "Demo Company",
		"status":       "active",
	}))
}

// UpdateAccountInfo 更新账户信息
func (m *Merchant) UpdateAccountInfo(c *gin.Context) {
	var req struct {
		CompanyName string `json:"company_name"`
		ContactName string `json:"contact_name"`
		Phone       string `json:"phone"`
		Address     string `json:"address"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Invalid update parameters: "+err.Error()))
		return
	}

	// TODO: 实现更新账户信息逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(nil))
}

// GetAccountBalance 账户余额详情
func (m *Merchant) GetAccountBalance(c *gin.Context) {
	var req struct {
		Currency string `json:"currency" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Currency is required: "+err.Error()))
		return
	}

	// TODO: 从JWT token中获取用户ID
	userID := "current_user_id"

	accountService := services.GetAccountService()
	balance, err := accountService.GetBalance(userID, "merchant", req.Currency)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewBusinessErrorResult("Balance not found: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(balance))
}

// GetAccountHistory 账户变动历史
func (m *Merchant) GetAccountHistory(c *gin.Context) {
	var req struct {
		Currency  string `json:"currency" binding:"required"`
		StartTime int64  `json:"start_time"`
		EndTime   int64  `json:"end_time"`
		Page      int    `json:"page"`
		Size      int    `json:"size"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Invalid parameters: "+err.Error()))
		return
	}

	// 设置默认分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}

	// TODO: 实现账户历史查询逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(gin.H{
		"records": []interface{}{},
		"total":   0,
		"page":    req.Page,
		"size":    req.Size,
	}))
}

// =============================================================================
// 开发者配置管理
// =============================================================================

// SubmitPublicKey 提交公钥
func (m *Merchant) SubmitPublicKey(c *gin.Context) {
	var req struct {
		PublicKey string `json:"public_key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Public key is required: "+err.Error()))
		return
	}

	// TODO: 实现提交公钥逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(nil))
}

// SubmitWebhook 设置 Webhook URL
func (m *Merchant) SubmitWebhook(c *gin.Context) {
	var req struct {
		WebhookURL string `json:"webhook_url" binding:"required,url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Valid webhook URL is required: "+err.Error()))
		return
	}

	// TODO: 实现设置Webhook逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(nil))
}

// SubmitWhitelist 设置 IP 白名单
func (m *Merchant) SubmitWhitelist(c *gin.Context) {
	var req struct {
		IpWhitelist []string `json:"ip_whitelist" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("IP whitelist is required: "+err.Error()))
		return
	}

	// TODO: 实现设置IP白名单逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(nil))
}

// =============================================================================
// 交易状态查询类（统一处理所有业务）
// =============================================================================

// ListTransactions 所有交易列表
func (m *Merchant) ListTransactions(c *gin.Context) {
	var req struct {
		Type      string `json:"type"` // receipt, payment, checkout, refund
		Status    string `json:"status"`
		Currency  string `json:"currency"`
		StartTime int64  `json:"start_time"`
		EndTime   int64  `json:"end_time"`
		Page      int    `json:"page"`
		Size      int    `json:"size"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Invalid parameters: "+err.Error()))
		return
	}

	// 设置默认分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}

	// TODO: 实现交易列表查询逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(gin.H{
		"records": []interface{}{},
		"total":   0,
		"page":    req.Page,
		"size":    req.Size,
	}))
}

// GetTransactionDetail 交易详情查询
func (m *Merchant) GetTransactionDetail(c *gin.Context) {
	var req struct {
		TransactionID string `json:"transaction_id"`
		BillID        string `json:"bill_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Invalid parameters: "+err.Error()))
		return
	}

	if req.TransactionID == "" && req.BillID == "" {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Either transaction_id or bill_id is required"))
		return
	}

	// TODO: 实现交易详情查询逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(gin.H{
		"transaction_id": req.TransactionID,
		"bill_id":        req.BillID,
		"status":         "completed",
	}))
}

// ExportTransactions 导出交易记录
func (m *Merchant) ExportTransactions(c *gin.Context) {
	var req struct {
		Type      string `json:"type"`
		Status    string `json:"status"`
		StartTime int64  `json:"start_time" binding:"required"`
		EndTime   int64  `json:"end_time" binding:"required"`
		Format    string `json:"format" binding:"required,oneof=csv xlsx"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Invalid parameters: "+err.Error()))
		return
	}

	// TODO: 实现导出交易记录逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(gin.H{
		"export_id": "export_id_here",
		"status":    "processing",
	}))
}

// CreateRefund 发起退款
func (m *Merchant) CreateRefund(c *gin.Context) {
	var req struct {
		OriginalTransactionID string `json:"original_transaction_id" binding:"required"`
		RefundAmount          string `json:"refund_amount" binding:"required"`
		Reason                string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Invalid refund parameters: "+err.Error()))
		return
	}

	// TODO: 实现发起退款逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(gin.H{
		"refund_id": "refund_id_here",
		"status":    "pending",
	}))
}

// =============================================================================
// 收银台管理
// =============================================================================

// ListCheckouts 收银台会话列表
func (m *Merchant) ListCheckouts(c *gin.Context) {
	var req struct {
		Status    string `json:"status"`
		StartTime int64  `json:"start_time"`
		EndTime   int64  `json:"end_time"`
		Page      int    `json:"page"`
		Size      int    `json:"size"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Invalid parameters: "+err.Error()))
		return
	}

	// 设置默认分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}

	// TODO: 实现收银台列表查询逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(gin.H{
		"records": []interface{}{},
		"total":   0,
		"page":    req.Page,
		"size":    req.Size,
	}))
}

// GetCheckoutServices 收银台服务配置
func (m *Merchant) GetCheckoutServices(c *gin.Context) {
	var req struct {
		Country string `json:"country"`
		Amount  string `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Invalid parameters: "+err.Error()))
		return
	}

	// TODO: 实现获取收银台服务配置逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(gin.H{
		"services": []interface{}{},
	}))
}

// =============================================================================
// Webhook 管理
// =============================================================================

// ListWebhooks Webhook 记录列表
func (m *Merchant) ListWebhooks(c *gin.Context) {
	var req struct {
		Status    string `json:"status"`
		StartTime int64  `json:"start_time"`
		EndTime   int64  `json:"end_time"`
		Page      int    `json:"page"`
		Size      int    `json:"size"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Invalid parameters: "+err.Error()))
		return
	}

	// 设置默认分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}

	// TODO: 实现Webhook列表查询逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(gin.H{
		"records": []gin.H{},
		"total":   0,
		"page":    req.Page,
		"size":    req.Size,
	}))
}

// RetryWebhook 重试 Webhook
func (m *Merchant) RetryWebhook(c *gin.Context) {
	var req struct {
		WebhookID string `json:"webhook_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Webhook ID is required: "+err.Error()))
		return
	}

	// TODO: 实现重试Webhook逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(gin.H{"message": "Webhook retried"}))
}

// TestWebhook 测试 Webhook
func (m *Merchant) TestWebhook(c *gin.Context) {
	var req struct {
		WebhookURL string `json:"webhook_url" binding:"required,url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult("Valid webhook URL is required: "+err.Error()))
		return
	}

	// TODO: 实现测试Webhook逻辑
	c.JSON(http.StatusOK, protocol.NewSuccessResult(gin.H{
		"test_result":   "success",
		"response_time": "150ms",
	}))
}
