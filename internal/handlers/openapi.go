package handlers

import (
	"inpayos/internal/config"
	"inpayos/internal/middleware"
	"inpayos/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OpenApi 第一层：纯开放 API 接口处理器
// 面向第三方开发者集成，需要 API Key 认证
type OpenApi struct {
	*config.ServiceConfig
	Transaction *services.TransactionService
	Checkout    *services.CheckoutService
}

// NewOpenApi 创建 OpenApi 处理器
func NewOpenApi() *OpenApi {
	cfg := config.Get()
	if cfg == nil || cfg.Server.OpenAPI == nil {
		return nil
	}

	return &OpenApi{
		ServiceConfig: cfg.Server.OpenAPI,
		Transaction:   services.GetTransactionService(),
		Checkout:      services.GetCheckoutService(),
	}
}

// ToServer 创建 HTTP服务器
func (a *OpenApi) ToServer() *http.Server {
	return &http.Server{
		Addr: ":" + a.Port,
	}
}

// SetupRouter 设置路由
func (a *OpenApi) SetupRouter() *gin.Engine {
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
		c.JSON(200, gin.H{"status": "ok", "service": "openapi"})
	})

	// API路由组 - 需要API Key认证
	prefix := a.Prefix
	if prefix == "" {
		prefix = "/openapi"
	}
	apiGroup := router.Group(prefix)
	apiGroup.Use(middleware.APIKeyAuth())
	{
		// 代收接口
		apiGroup.POST("/payin", a.Payin)
		// 代付接口
		apiGroup.POST("/payout", a.Payout)
		apiGroup.POST("/cancel", a.Cancel)

		// 收银台接口
		checkout := apiGroup.Group("/checkout")
		{
			checkout.POST("", a.CreateCheckout)
			checkout.GET("/info", a.GetCheckout)
			checkout.POST("/confirm", a.ConfirmCheckout)
			checkout.POST("/cancel", a.CancelCheckout)
		}

		// 查询接口
		apiGroup.POST("/balance", a.Balance)
		apiGroup.POST("/query", a.Query)
	}

	return router
}
