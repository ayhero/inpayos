package handlers

import (
	"inpayos/internal/config"
	"inpayos/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Admin 第三层：平台管理接口处理器
// 面向运营人员，需要 JWT + 权限控制认证
type Admin struct {
	*config.ServiceConfig
}

// NewAdmin 创建 Admin 处理器
func NewAdmin() *Admin {
	cfg := config.Get()
	if cfg == nil || cfg.Server.Admin == nil {
		return nil
	}

	return &Admin{
		ServiceConfig: cfg.Server.Admin,
	}
}

// ToServer 创建 HTTP服务器
func (a *Admin) ToServer() *http.Server {
	return &http.Server{
		Addr: ":" + a.Port,
	}
}

// SetupRouter 设置路由
func (a *Admin) SetupRouter() *gin.Engine {
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
		c.JSON(200, gin.H{"status": "ok", "service": "admin"})
	})

	// 需要JWT认证的端点
	adminAPI := router.Group("/")
	adminAPI.Use(middleware.JWTAuth())
	adminAPI.Use(middleware.PermissionCheck())
	return router
}
