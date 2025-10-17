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
		// 基础实体
		&Account{},
		&Admin{},
		&Merchant{},
		&MerchantSecret{},
		&Cashier{},
		&CashierTeam{},

		// 交易相关
		&MerchantPayin{},
		&MerchantPayout{},
		&MerchantCheckout{},
		&Deposit{},
		&Withdraw{},
		&CashierPayin{},
		&CashierPayout{},

		//渠道相关
		&ChannelAccount{},
		&ChannelGroup{},

		// 资金和流水
		&FundFlow{},
		&TrxHistory{},

		// 配置相关
		&MerchantConfig{},
		&MerchantFeeConfig{},
		&APIConfig{},
		&Contract{},

		// 路由和渠道
		&MerchantRouter{},
		&CashierRouter{},
		&ChannelGroup{},

		// 结算相关
		&MerchantSettleLog{},
		&MerchantSettleTransaction{},
		&MerchantSettleHistory{},
		&SettleRule{},

		// 通知和消息
		&Webhook{},
		&MessageTemplate{},
		&FCMToken{},

		// 统计和系统
		&SummaryStats{},
		&Task{},
	)
}
