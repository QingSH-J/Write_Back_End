package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/jinxinyu/go_backend/internal/config"
	"github.com/jinxinyu/go_backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
		cfg.DBTimeZone,
	)
	log.Printf("尝试连接数据库: %s@%s:%s/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	gormConfig := &gorm.Config{}
	if cfg.Environment != "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Warn)
	}

	// 添加重试逻辑
	var db *gorm.DB
	var err error
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("数据库连接尝试 %d/%d", attempt+1, maxRetries)
			time.Sleep(time.Second * time.Duration(attempt*2)) // 递增等待时间
		}

		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err == nil {
			break
		}

		log.Printf("数据库连接尝试 %d 失败: %v", attempt+1, err)
	}

	if err != nil {
		log.Printf("经过多次尝试后仍然无法连接数据库: %v", err)
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Printf("数据库连接成功")

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("无法获取数据库实例: %v", err)
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	} else {
		// 设置适当的连接池配置
		sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBConnMaxLifetimeMinutes) * time.Minute)

		// 设置最短连接保持时间
		sqlDB.SetConnMaxIdleTime(time.Minute * 5)

		pingStart := time.Now()
		if err := sqlDB.Ping(); err != nil {
			log.Printf("数据库Ping失败: %v", err)
			return nil, fmt.Errorf("failed to ping database: %v", err)
		} else {
			log.Printf("数据库Ping成功，耗时: %v", time.Since(pingStart))
		}
	}

	log.Printf("开始自动迁移表结构")
	migrateStart := time.Now()

	if err := db.AutoMigrate(
		&models.User{},
		&models.WriteLog{},
	); err != nil {
		log.Printf("自动迁移失败: %v", err)
		return nil, fmt.Errorf("failed to auto migrate: %v", err)
	} else {
		log.Printf("自动迁移成功，耗时: %v", time.Since(migrateStart))
	}

	return db, nil
}
