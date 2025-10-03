package main

import (
	"fmt"
	"log"

	"inpayos/internal/config"
	"inpayos/internal/handlers"
	"inpayos/internal/i18n"
	"inpayos/internal/models"
	"inpayos/internal/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
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

	return nil
}

var (
	g errgroup.Group
)

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

	// 初始化服务单例
	err = services.InitializeServices()
	if err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	// 启动OpenAPI服务
	g.Go(func() error {
		openApiService := handlers.NewOpenApi()
		if openApiService == nil {
			log.Fatal("Failed to create OpenAPI service - configuration may be invalid")
			return fmt.Errorf("failed to create OpenAPI service")
		}
		server := openApiService.ToServer()
		server.Handler = openApiService.SetupRouter()

		log.Printf("Starting OpenAPI Service on port %s", openApiService.Port)
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("OpenAPI Service error: %v", err)
		} else {
			log.Println("OpenAPI Service started successfully")
		}
		return err
	})

	// 启动Merchant服务
	g.Go(func() error {
		merchantService := handlers.NewMerchantAdmin()
		if merchantService == nil {
			log.Fatal("Failed to create Merchant service - configuration may be invalid")
			return fmt.Errorf("failed to create Merchant service")
		}
		server := merchantService.ToServer()
		server.Handler = merchantService.SetupRouter()

		log.Printf("Starting Merchant Service on port %s", merchantService.Port)
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Merchant Service error: %v", err)
		} else {
			log.Println("Merchant Service started successfully")
		}
		return err
	})

	// 启动Admin服务
	g.Go(func() error {
		adminService := handlers.NewAdmin()
		if adminService == nil {
			log.Fatal("Failed to create Admin service - configuration may be invalid")
			return fmt.Errorf("failed to create Admin service")
		}
		server := adminService.ToServer()
		server.Handler = adminService.SetupRouter()

		log.Printf("Starting Admin Service on port %s", adminService.Port)
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Admin Service error: %v", err)
		} else {
			log.Println("Admin Service started successfully")
		}
		return err
	})

	// 等待所有服务
	if err := g.Wait(); err != nil {
		panic(err)
	}
}
