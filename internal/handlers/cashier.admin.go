package handlers

import (
	"fmt"
	cashieradmindocs "inpayos/docs/cashier_admin"
	"inpayos/internal/config"
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title InPayOS CashierAdmin API
// @version 1.0
// @description 出纳员管理API接口文档，提供出纳员认证、交易管理等功能
// @termsOfService http://swagger.io/terms/
// @contact.name InPayOS Support
// @contact.url http://www.inpayos.com/support
// @contact.email support@inpayos.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:6084
// @BasePath /
// @securityDefinitions.bearer BearerAuth
// @in header
// @name Authorization
// @description JWT认证令牌，请在请求头中添加Authorization: Bearer <token>
type CashierAdmin struct {
	*config.ServiceConfig
}

func NewCashierAdmin() *CashierAdmin {
	cfg := config.Get()
	if cfg == nil || cfg.Server.CashierAdmin == nil {
		return nil
	}

	return &CashierAdmin{
		ServiceConfig: cfg.Server.CashierAdmin,
	}
}

func (t *CashierAdmin) SetupRouter() *gin.Engine {
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

	// Initialize Swagger Info for CashierAdmin
	cashieradmindocs.SwaggerInfocashieradmin.Title = "InPayOS CashierAdmin API"
	cashieradmindocs.SwaggerInfocashieradmin.Description = "出纳员管理API接口文档，提供出纳员认证、交易管理等功能"
	cashieradmindocs.SwaggerInfocashieradmin.Version = "1.0"
	cashieradmindocs.SwaggerInfocashieradmin.Host = "localhost:6084"
	cashieradmindocs.SwaggerInfocashieradmin.BasePath = "/"
	cashieradmindocs.SwaggerInfocashieradmin.Schemes = []string{"http", "https"}

	// 添加Swagger文档路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName("cashieradmin")))
	// API路由组
	api := router.Group(fmt.Sprintf("/%s", t.Prefix))
	api.POST("/auth", t.Auth)     // 注册授权路由
	api.POST("/logout", t.Logout) // 注册登出路由
	// 无需认证的路由
	api.POST("/verifycode/send", SendVerifyCode) // 发送验证码
	api.POST("/verifycode/verify", VerifyCode)   // 验证验证码
	api.POST("/register", t.Register)            // 注册商户
	api.POST("/password/reset", t.ResetPassword) // 重置密码
	// 注册JWT中间件
	api.Use(middleware.CashierTeamJWTAuth())
	api.POST("/info", t.Info)                      // 商户信息
	api.POST("/password/change", t.ChangePassword) // 修改密码

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

	// 出纳员相关路由
	cashiers := api.Group("/cashiers")
	{
		cashiers.POST("/list", t.ListCashiers)    // 出纳员列表
		cashiers.POST("/detail", t.CashierDetail) // 出纳员详情
	}

	return router
}
