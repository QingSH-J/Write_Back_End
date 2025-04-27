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
	log.Printf("Connecting to database...")
	gormConfig := &gorm.Config{}
	if cfg.Environment != "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Warn)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Printf("Connected to database")

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Failed to get database instance: %v", err)
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	} else {
		sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBConnMaxLifetimeMinutes) * time.Minute)

		if err := sqlDB.Ping(); err != nil {
			log.Printf("Failed to ping database: %v", err)
			return nil, fmt.Errorf("failed to ping database: %v", err)
		} else {
			log.Printf("Database pinged successfully")
		}
	}
	log.Printf("Database connection successful")
	log.Printf("Database pinged successfully")
	log.Printf("Trying to auto migrate")

	if err := db.AutoMigrate(
		&models.User{},
		&models.Wallet{},
	); err != nil {
		log.Printf("Failed to auto migrate: %v", err)
		return nil, fmt.Errorf("failed to auto migrate: %v", err)
	} else {
		log.Printf("Auto migrate successful")
	}

	return db, nil

}
