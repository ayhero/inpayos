package models

import (
	"fmt"
	"inpayos/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) error {
	dsn := cfg.Database.DSN
	if dsn == "" {
		return fmt.Errorf("database DSN is empty")
	}

	var logLevel logger.LogLevel
	if cfg.Debug {
		logLevel = logger.Info
	} else {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.Database.ConnMaxIdleTime)

	DB = db
	return nil
}

// AutoMigrate runs database migrations
func AutoMigrate() error {
	return DB.AutoMigrate(
		&Account{},
		&Asset{},
		&FundFlow{},
		&Receipt{},
		&Payment{},
		&Deposit{},
		&Withdraw{},
		&Refund{},
		&Webhook{},
		&Channel{},
		&Cashier{},
		&MerchantConfig{},
		&MerchantSecret{},
		// 新增的核心业务表
		&FeeConfig{},
		&CheckoutSession{},
		&APIConfig{},
		&TrxRouter{},
	)
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
