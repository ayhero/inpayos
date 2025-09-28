package handlers

import (
	"inpayos/internal/config"
	"inpayos/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置三层架构路由
func SetupRouter() *gin.Engine {
	// 设置Gin模式
	cfg := config.Get()
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// 添加中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LanguageMiddleware()) // 语言处理中间件
	router.Use(middleware.RequestLogger())
	router.Use(middleware.ErrorHandler())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "inpayos"})
	})

	// =============================================================================
	// 第一层：纯开放 API 接口 (需要 API Key 认证)
	// =============================================================================
	openapi := NewOpenApi()
	apiV1 := router.Group("/api/v1")
	// 添加API Key认证中间件
	apiV1.Use(middleware.APIKeyMiddleware())
	{
		// 代收接口 (需要 receipt 权限)
		receipt := apiV1.Group("/receipt")
		receipt.Use(middleware.RequirePermission("receipt"))
		{
			receipt.POST("", openapi.CreateReceipt)
			receipt.POST("/cancel", openapi.CancelReceipt)
		}

		// 代付接口 (需要 payment 权限)
		payment := apiV1.Group("/payment")
		payment.Use(middleware.RequirePermission("payment"))
		{
			payment.POST("", openapi.CreatePayment)
			payment.POST("/cancel", openapi.CancelPayment)
		}

		// 收银台接口 (需要 checkout 权限)
		checkout := apiV1.Group("/checkout")
		checkout.Use(middleware.RequirePermission("checkout"))
		{
			checkout.POST("", openapi.CreateCheckout)
			checkout.GET("/:id", openapi.GetCheckout)
			checkout.POST("/configs", openapi.GetCheckoutConfigs)
			checkout.POST("/confirm", openapi.ConfirmCheckout)
			checkout.POST("/cancel", openapi.CancelCheckout)
		}

		// 统一查询接口 (需要 query 权限)
		transaction := apiV1.Group("/transaction")
		transaction.Use(middleware.RequirePermission("query"))
		{
			transaction.POST("/query", openapi.QueryTransaction)
			transaction.POST("/list", openapi.ListTransactions)
		}

		// 余额查询 (需要 query 权限)
		balance := apiV1.Group("/balance")
		balance.Use(middleware.RequirePermission("query"))
		{
			balance.POST("/query", openapi.QueryBalance)
		}
	}

	// =============================================================================
	// 第二层：Merchant 商户后台接口 (需要 JWT + Session 认证)
	// =============================================================================
	merchant := NewMerchant()
	merchantV1 := router.Group("/merchant/v1")

	// 商户认证管理（无需JWT验证）
	auth := merchantV1.Group("/auth")
	{
		auth.POST("/login", merchant.Login)
		auth.POST("/logout", merchant.Logout)
		auth.POST("/register", merchant.Register)
		auth.POST("/verify", merchant.Verify)
		auth.POST("/reset-password", merchant.ResetPassword)
	}

	// 需要JWT认证的商户接口
	merchantAuth := merchantV1.Group("")
	merchantAuth.Use(middleware.JWTMiddleware(cfg.JWT.Secret), middleware.MerchantPermissionMiddleware())
	{
		// 修改密码需要认证
		auth.POST("/change-password", merchant.ChangePassword)

		// 账户信息管理
		account := merchantAuth.Group("/account")
		{
			account.POST("/info", merchant.GetAccountInfo)
			account.POST("/update", merchant.UpdateAccountInfo)
			account.POST("/balance", merchant.GetAccountBalance)
			account.POST("/history", merchant.GetAccountHistory)
		}

		// 开发者配置管理
		dev := merchantAuth.Group("/dev")
		{
			dev.POST("/submit-public-key", merchant.SubmitPublicKey)
			dev.POST("/submit-webhook", merchant.SubmitWebhook)
			dev.POST("/submit-whitelist", merchant.SubmitWhitelist)
		}

		// 交易状态查询类
		transactions := merchantAuth.Group("/transactions")
		{
			transactions.POST("/list", merchant.ListTransactions)
			transactions.POST("/detail", merchant.GetTransactionDetail)
			transactions.POST("/export", merchant.ExportTransactions)
			transactions.POST("/refund", merchant.CreateRefund)
		}

		// 收银台管理
		checkout := merchantAuth.Group("/checkout")
		{
			checkout.POST("/list", merchant.ListCheckouts)
			checkout.POST("/services", merchant.GetCheckoutServices)
		}

		// Webhook 管理
		webhooks := merchantAuth.Group("/webhooks")
		{
			webhooks.POST("/list", merchant.ListWebhooks)
			webhooks.POST("/retry", merchant.RetryWebhook)
			webhooks.POST("/test", merchant.TestWebhook)
		}
	}

	// =============================================================================
	// 第三层：Admin 平台管理接口 (需要 JWT + 权限控制认证)
	// =============================================================================
	admin := NewAdmin()
	adminV1 := router.Group("/admin/v1")
	adminV1.Use(middleware.JWTMiddleware(cfg.JWT.Secret), middleware.AdminPermissionMiddleware())
	{
		// 商户管理
		merchants := adminV1.Group("/merchants")
		{
			merchants.POST("/list", admin.ListMerchants)
			merchants.POST("/detail", admin.GetMerchantDetail)
			merchants.POST("/create", admin.CreateMerchant)
			merchants.POST("/update", admin.UpdateMerchant)
			merchants.POST("/status", admin.UpdateMerchantStatus)
			merchants.POST("/approve", admin.ApproveMerchant)
		}

		// Cashier 管理
		cashiers := adminV1.Group("/cashiers")
		{
			cashiers.POST("/list", admin.ListCashiers)
			cashiers.POST("/detail", admin.GetCashierDetail)
			cashiers.POST("/create", admin.CreateCashier)
			cashiers.POST("/update", admin.UpdateCashier)
			cashiers.POST("/status", admin.UpdateCashierStatus)
			cashiers.POST("/delete", admin.DeleteCashier)
			cashiers.POST("/assign", admin.AssignCashier)
			cashiers.POST("/stats", admin.GetCashierStats)
		}

		// 渠道管理
		channels := adminV1.Group("/channels")
		{
			channels.POST("/list", admin.ListChannels)
			channels.POST("/create", admin.CreateChannel)
			channels.POST("/update", admin.UpdateChannel)
			channels.POST("/status", admin.UpdateChannelStatus)
			channels.POST("/stats", admin.GetChannelStats)
		}

		// 系统配置管理
		configs := adminV1.Group("/configs")
		{
			configs.POST("/api", admin.GetApiConfigs)
			configs.POST("/api/save", admin.SaveApiConfig)
			configs.POST("/fees", admin.GetFeeConfigs)
			configs.POST("/fees/save", admin.SaveFeeConfig)
		}

		// 交易管理
		transactions := adminV1.Group("/transactions")
		{
			transactions.POST("/list", admin.ListAllTransactions)
			transactions.POST("/detail", admin.GetTransactionDetail)
			transactions.POST("/audit", admin.AuditTransaction)
		}
	}

	return router
}
