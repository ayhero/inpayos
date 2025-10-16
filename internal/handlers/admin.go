package handlers

import (
	admindocs "inpayos/docs/admin"
	"inpayos/internal/config"
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title InPayOS Admin API
// @version 1.0
// @description 平台管理API接口文档，提供管理员管理、系统配置等功能
// @termsOfService http://swagger.io/terms/
// @contact.name InPayOS Support
// @contact.url http://www.inpayos.com/support
// @contact.email support@inpayos.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:6082
// @BasePath /
// @securityDefinitions.bearer BearerAuth
// @in header
// @name Authorization
// @description JWT认证令牌，请在请求头中添加Authorization: Bearer <token>

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
		c.JSON(200, gin.H{"status": "ok", "service": "admin"})
	})

	// Initialize Swagger Info for Admin
	admindocs.SwaggerInfoadmin.Title = "InPayOS Admin API"
	admindocs.SwaggerInfoadmin.Description = "平台管理API接口文档，提供管理员管理、系统配置等功能"
	admindocs.SwaggerInfoadmin.Version = "1.0"
	admindocs.SwaggerInfoadmin.Host = "localhost:6082"
	admindocs.SwaggerInfoadmin.BasePath = "/"
	admindocs.SwaggerInfoadmin.Schemes = []string{"http", "https"}

	// 添加Swagger文档路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName("admin")))

	// 需要JWT认证的端点
	adminAPI := router.Group("/")
	adminAPI.Use(middleware.AdminJWTAuth())
	adminAPI.Use(middleware.PermissionCheck())
	return router
}
