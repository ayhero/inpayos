package handlers

import (
	"fmt"
	"inpayos/internal/config"
	"inpayos/internal/middleware"

	"github.com/gin-gonic/gin"
)

type MerchantAdmin struct {
	*config.ServiceConfig
}

func NewMerchantAdmin() *MerchantAdmin {
	cfg := config.Get()
	if cfg == nil || cfg.Server.Merchant == nil {
		return nil
	}

	return &MerchantAdmin{
		ServiceConfig: cfg.Server.Merchant,
	}
}

func (t *MerchantAdmin) SetupRouter() *gin.Engine {
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
		dashboard.POST("/today-stats", t.GetTodayStats)             // 今日统计
		dashboard.POST("/transaction-trend", t.GetTransactionTrend) // 交易趋势
		dashboard.POST("/settlement-trend", t.GetSettlementTrend)   // 结算趋势
		dashboard.POST("/overview", t.GetDashboardOverview)         // Dashboard概览
		dashboard.POST("/account-balance", t.GetAccountBalance)     // 账户余额
	}

	return router
}
