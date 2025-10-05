package models

import (
	"fmt"
	"inpayos/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	WriteDB *gorm.DB
	ReadDB  *gorm.DB
)

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return WriteDB
}

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) error {
	if db, err := initDBByConfig(cfg.WriteDB); err != nil {
		return err
	} else {
		WriteDB = db
	}
	if db, err := initDBByConfig(cfg.ReadDB); err != nil {
		return err
	} else {
		ReadDB = db
	}
	return nil
}

func initDBByConfig(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	if cfg.DSN == "" {
		return nil, fmt.Errorf("database DSN is empty")
	}

	var logLevel logger.LogLevel
	if cfg.Debug {
		logLevel = logger.Info
	} else {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return db, nil
}

// AutoMigrate runs database migrations
func AutoMigrate() error {
	return WriteDB.AutoMigrate(
		&Account{},
		&Asset{},
		&FundFlow{},
		&MerchantPayin{},
		&MerchantPayout{},
		&Webhook{},
		&Cashier{},
		&MerchantConfig{},
		&MerchantSecret{},
		// 新增的核心业务表
		&MerchantFeeConfig{},
		&MerchantCheckout{},
		&APIConfig{},
	)
}
