package handlers

import (
	"fmt"
	"inpayos/internal/config"
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"

	apidocs "inpayos/docs/openapi"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title InPayOS OpenAPI
// @version 1.0
// @description 统一支付网关OpenAPI接口文档，提供代收、代付、收银台等核心支付功能
// @termsOfService http://swagger.io/terms/
// @contact.name InPayOS Support
// @contact.url http://www.inpayos.com/support
// @contact.email support@inpayos.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:6080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-api-key
// @description API密钥认证，请在请求头中添加x-api-key字段
type OpenApi struct {
	*config.ServiceConfig
	Transaction *services.MerchantTransactionService
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
		Transaction:   services.GetMerchantTransactionService(),
		Checkout:      services.GetCheckoutService(),
	}
}

// ToServer 创建 HTTP服务器
func (a *OpenApi) ToServer() *http.Server {
	return &http.Server{
		Addr: fmt.Sprintf(":%v", a.Port),
	}
}

// SetupRouter 设置路由
func (a *OpenApi) SetupRouter() *gin.Engine {
	// 设置Gin模式
	cfg := config.Get()
	if cfg.Env == protocol.EnvProduction {
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

	// Initialize Swagger Info
	apidocs.SwaggerInfoopenapi.Title = "OpenAPI"
	apidocs.SwaggerInfoopenapi.Description = "OpenAPI Documentation"
	apidocs.SwaggerInfoopenapi.Version = "1.0"
	apidocs.SwaggerInfoopenapi.Host = ""
	apidocs.SwaggerInfoopenapi.BasePath = "/"
	apidocs.SwaggerInfoopenapi.Schemes = []string{"http", "https"}

	// Swagger UI 路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName("openapi")))

	// API路由组 - 需要API Key认证
	apiGroup := router.Group("")
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
		}

		// 查询接口
		apiGroup.POST("/balance", a.Balance)
		apiGroup.POST("/query", a.Query)
	}

	return router
}
