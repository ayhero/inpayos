package handlers

import (
	"fmt"
	merchantdocs "inpayos/docs/merchant"
	"inpayos/internal/config"
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"inpayos/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title InPayOS Merchant API
// @version 1.0
// @description 商户管理API接口文档，提供商户注册、认证、交易管理等功能
// @termsOfService http://swagger.io/terms/
// @contact.name InPayOS Support
// @contact.url http://www.inpayos.com/support
// @contact.email support@inpayos.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:6081
// @BasePath /
// @securityDefinitions.bearer BearerAuth
// @in header
// @name Authorization
// @description JWT认证令牌，请在请求头中添加Authorization: Bearer <token>
type MerchantAdmin struct {
	*config.ServiceConfig
	Checkout *services.CheckoutService
}

func NewMerchantAdmin() *MerchantAdmin {
	cfg := config.Get()
	if cfg == nil || cfg.Server.Merchant == nil {
		return nil
	}

	return &MerchantAdmin{
		ServiceConfig: cfg.Server.Merchant,
		Checkout:      services.GetCheckoutService(),
	}
}

func (t *MerchantAdmin) SetupRouter() *gin.Engine {
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

	// Initialize Swagger Info for Merchant
	merchantdocs.SwaggerInfomerchant.Title = "InPayOS Merchant API"
	merchantdocs.SwaggerInfomerchant.Description = "商户管理API接口文档，提供商户注册、认证、交易管理等功能"
	merchantdocs.SwaggerInfomerchant.Version = "1.0"
	merchantdocs.SwaggerInfomerchant.Host = "localhost:6081"
	merchantdocs.SwaggerInfomerchant.BasePath = "/"
	merchantdocs.SwaggerInfomerchant.Schemes = []string{"http", "https"}

	// 添加Swagger文档路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName("merchant")))
	// API路由组
	api := router.Group(fmt.Sprintf("/%s", t.Prefix))
	api.POST("/auth", t.Auth)     // 注册授权路由
	api.POST("/logout", t.Logout) // 注册登出路由
	// 无需认证的路由
	api.POST("/verifycode/send", SendVerifyCode) // 发送验证码
	api.POST("/verifycode/verify", VerifyCode)   // 验证验证码
	api.POST("/register", t.Register)            // 注册商户

	// 注册JWT中间件
	api.Use(middleware.MerchantJWTAuth())
	api.POST("/info", t.Info)                      // 商户信息
	api.POST("/password/change", t.ChangePassword) // 修改密码
	api.POST("/password/reset", t.ResetPassword)   // 重置密码
	api.POST("/checkout/info", t.CheckoutInfo)

	// G2FA相关路由
	g2fa := api.Group("/g2fa")
	{
		g2fa.POST("/bind", t.BindG2FA) // G2FA绑定
		g2fa.POST("/new", t.NewG2FA)   // 生成新的G2FA密钥
	}

	// 交易相关路由
	transactions := api.Group("/transactions")
	{
		transactions.POST("/list", t.ListTransactions)                // 交易列表
		transactions.POST("/detail", t.TransactionDetail)             // 交易详情
		transactions.POST("/today-stats", t.GetTransactionTodayStats) // 今日统计
	}

	// Dashboard相关路由
	dashboard := api.Group("/dashboard")
	{
		dashboard.POST("/today-stats", t.GetTodayStats)         // 今日统计
		dashboard.POST("/overview", t.GetDashboardOverview)     // Dashboard概览
		dashboard.POST("/account-balance", t.GetAccountBalance) // 账户余额
	}

	// 账户相关路由
	account := api.Group("/account")
	{
		account.GET("/list", t.AccountList)           // 账户列表
		account.POST("/flow/list", t.AccountFlowList) // 账户流水列表
	}
	checkout := api.Group("/checkout")
	{
		checkout.POST("/submit", t.SubmitCheckout)
		checkout.POST("/services", t.CheckoutServices)
		checkout.POST("/confirm", t.ConfirmCheckout)
		checkout.POST("/cancel", t.CancelCheckout)
	}

	return router
}
