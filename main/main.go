package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"inpayos/internal/config"
	"inpayos/internal/handlers"
	"inpayos/internal/i18n"
	"inpayos/internal/models"
	"inpayos/internal/services"
)

func InitialConfig() error {
	cfg := config.Get()

	// 初始化数据库
	if err := models.InitDB(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return err
	}
	// 运行数据库迁移
	if err := models.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
		return err
	}

	// 初始化国际化系统
	i18n.InitI18n("internal/locales")

	// 初始化服务
	services.SetupService()

	return nil
}

func main() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		panic(err)
	}

	// 初始化应用配置
	err = InitialConfig()
	if err != nil {
		panic(err)
	}

	cfg := config.Get()

	// 设置Gin模式
	if cfg.Env == config.ProdEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	// 启动三层架构服务
	router := handlers.SetupRouter()

	log.Printf("Starting inpayos server on port :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
