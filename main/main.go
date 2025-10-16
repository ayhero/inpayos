package main

import (
	"fmt"

	"inpayos/internal/channels"
	"inpayos/internal/config"
	"inpayos/internal/handlers"
	"inpayos/internal/i18n"
	"inpayos/internal/log"
	"inpayos/internal/models"
	"inpayos/internal/services"

	_ "inpayos/docs/admin"         // 导入Admin文档
	_ "inpayos/docs/cashier"       // 导入Cashier文档
	_ "inpayos/docs/cashier_admin" // 导入CashierAdmin文档
	_ "inpayos/docs/merchant"      // 导入Merchant文档
	_ "inpayos/docs/openapi"       // 导入OpenAPI文档

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func InitialConfig() error {
	cfg := config.Get()

	// 初始化数据库
	if err := models.InitDB(cfg); err != nil {
		log.Get().Fatalf("Failed to connect to database: %v", err)
		return err
	}
	// 运行数据库迁移
	if err := models.AutoMigrate(); err != nil {
		log.Get().Fatalf("Failed to run database migrations: %v", err)
		return err
	}
	if err := models.InitRedis(); err != nil {
		log.Get().Fatalf("Failed to connect to Redis: %v", err)
		return err
	}

	// 初始化国际化
	translator := i18n.NewFileTranslator()
	i18n.SetGlobalTranslator(translator)

	// 加载翻译资源，仅在配置了路径时加载
	if cfg.I18n != nil && cfg.I18n.LocalesDir != "" {
		localesDir := cfg.I18n.LocalesDir
		if err := translator.LoadTranslations(localesDir); err != nil {
			log.Get().Warnf("I18n: Failed to load translations from '%s': %v", localesDir, err)
		} else {
			log.Get().Infof("I18n: Translations loaded successfully from '%s'", localesDir)
		}
	} else {
		log.Get().Info("I18n: locales_dir not configured, using default English messages")
	}
	go channels.LoadChannelOpenApiService()
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
	err = services.InitializeMerchantServices()
	if err != nil {
		log.Get().Fatalf("Failed to initialize services: %v", err)
	}

	// 启动OpenAPI服务
	g.Go(func() error {
		app := handlers.NewOpenApi()
		if app == nil {
			log.Get().Fatal("Failed to create OpenAPI service - configuration may be invalid")
			return fmt.Errorf("failed to create OpenAPI service")
		}
		server := app.ToServer()
		server.Handler = app.SetupRouter()

		log.Get().Printf("Starting OpenAPI Service on port %s", app.Port)
		err := server.ListenAndServe()
		if err != nil {
			log.Get().Printf("OpenAPI Service error: %v", err)
		} else {
			log.Get().Println("OpenAPI Service started successfully")
		}
		return err
	})

	// 启动Merchant服务
	g.Go(func() error {
		app := handlers.NewMerchantAdmin()
		if app == nil {
			log.Get().Fatal("Failed to create Merchant service - configuration may be invalid")
			return fmt.Errorf("failed to create Merchant service")
		}
		server := app.ToServer()
		server.Handler = app.SetupRouter()

		log.Get().Printf("Starting Merchant Service on port %s", app.Port)
		err := server.ListenAndServe()
		if err != nil {
			log.Get().Printf("Merchant Service error: %v", err)
		} else {
			log.Get().Println("Merchant Service started successfully")
		}
		return err
	})

	// 启动Cashier Admin服务
	g.Go(func() error {
		app := handlers.NewCashierApi()
		if app == nil {
			log.Get().Fatal("Failed to create Cashier service - configuration may be invalid")
			return fmt.Errorf("failed to create Cashier service")
		}
		server := app.ToServer()
		server.Handler = app.SetupRouter()

		log.Get().Printf("Starting Cashier Service on port %s", app.Port)
		err := server.ListenAndServe()
		if err != nil {
			log.Get().Printf("Cashier Service error: %v", err)
		} else {
			log.Get().Println("Cashier Service started successfully")
		}
		return err
	})
	// 启动Cashier Admin服务
	g.Go(func() error {
		app := handlers.NewCashierAdmin()
		if app == nil {
			log.Get().Fatal("Failed to create Cashier Admin service - configuration may be invalid")
			return fmt.Errorf("failed to create Cashier Admin service")
		}
		server := app.ToServer()
		server.Handler = app.SetupRouter()

		log.Get().Printf("Starting Cashier Admin Service on port %s", app.Port)
		err := server.ListenAndServe()
		if err != nil {
			log.Get().Printf("Cashier Admin Service error: %v", err)
		} else {
			log.Get().Println("Cashier Admin Service started successfully")
		}
		return err
	})
	// 启动Admin服务
	g.Go(func() error {
		app := handlers.NewAdmin()
		if app == nil {
			log.Get().Fatal("Failed to create Admin service - configuration may be invalid")
			return fmt.Errorf("failed to create Admin service")
		}
		server := app.ToServer()
		server.Handler = app.SetupRouter()

		log.Get().Printf("Starting Admin Service on port %s", app.Port)
		err := server.ListenAndServe()
		if err != nil {
			log.Get().Printf("Admin Service error: %v", err)
		} else {
			log.Get().Println("Admin Service started successfully")
		}
		return err
	})

	// 等待所有服务
	if err := g.Wait(); err != nil {
		panic(err)
	}
}
