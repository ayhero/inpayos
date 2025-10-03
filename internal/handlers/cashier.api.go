package handlers

import (
	"inpayos/internal/config"
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CashierApi 第一层：纯开放 API 接口处理器
// 面向第三方开发者集成，需要 API Key 认证
type CashierApi struct {
	*config.ServiceConfig
	Transaction *services.TransactionService
}

// NewCashierApi 创建 CashierApi 处理器
func NewCashierApi() *CashierApi {
	cfg := config.Get()
	if cfg == nil || cfg.Server.CashierAPI == nil {
		return nil
	}

	return &CashierApi{
		ServiceConfig: cfg.Server.CashierAPI,
		Transaction:   services.GetTransactionService(),
	}
}

// ToServer 创建 HTTP服务器
func (a *CashierApi) ToServer() *http.Server {
	return &http.Server{
		Addr: ":" + a.Port,
	}
}

// SetupRouter 设置路由
func (a *CashierApi) SetupRouter() *gin.Engine {
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
		c.JSON(200, gin.H{"status": "ok", "service": "cashier"})
	})

	// API路由组 - 需要API Key认证
	prefix := a.Prefix
	if prefix == "" {
		prefix = "/cashier"
	}
	apiGroup := router.Group(prefix)
	apiGroup.Use(middleware.APIKeyAuth())
	{
		// 代收接口
		apiGroup.POST("/payin", a.Payin)
		// 代付接口
		apiGroup.POST("/payout", a.Payout)
		apiGroup.POST("/cancel", a.Cancel)
		apiGroup.POST("/query", a.Query)
	}

	return router
}

// =============================================================================
// 代收（Payin）接口
// =============================================================================

// Payin 创建代收订单
func (a *CashierApi) Payin(c *gin.Context) {
	var req protocol.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 执行业务逻辑（代收类型）
	response, code := services.GetTransactionService().CreatePayin(&req)
	lang := middleware.GetLanguage(c)
	result := protocol.HandleServiceResult(code, response, lang)
	c.JSON(http.StatusOK, result)
}

// Cancel 取消订单
func (a *CashierApi) Cancel(c *gin.Context) {
	var req protocol.CancelTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 执行取消逻辑
	response, code := services.GetTransactionService().Cancel(&req)
	lang := middleware.GetLanguage(c)
	result := protocol.HandleServiceResult(code, response, lang)
	c.JSON(http.StatusOK, result)
}

// Payout 创建代付订单
func (a *CashierApi) Payout(c *gin.Context) {
	var req protocol.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 执行业务逻辑（代付类型）
	response, code := services.GetTransactionService().CreatePayout(&req)
	lang := middleware.GetLanguage(c)
	result := protocol.HandleServiceResult(code, response, lang)
	c.JSON(http.StatusOK, result)
}

// =============================================================================
// 统一查询接口
// =============================================================================

// Query 查询交易状态/详情
func (a *CashierApi) Query(c *gin.Context) {
	var req protocol.QueryTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	if req.ReqID == "" && req.TrxID == "" {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.MissingParams, lang))
		return
	}

	// 从上下文获取交易服务
	transactionService := services.GetTransactionService()
	if transactionService == nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.SystemError, lang))
		return
	}
	response, code := transactionService.Query(&req)
	lang := middleware.GetLanguage(c)
	result := protocol.HandleServiceResult(code, response, lang)
	c.JSON(http.StatusOK, result)
}
