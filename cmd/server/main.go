package main

import (
	"log"

	"github.com/alexedwards/argon2id"
	"github.com/jinxinyu/go_backend/internal/auth"
	"github.com/jinxinyu/go_backend/internal/config"
	"github.com/jinxinyu/go_backend/internal/router"
	"github.com/jinxinyu/go_backend/internal/storage"
	"github.com/jinxinyu/go_backend/internal/utils"
)

func main() {
	// load config
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	// initialize database
	db, err := storage.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// initialize hash utils
	hashutils := utils.NewHashedPassword(argon2id.DefaultParams)
	tokenmaker, err := utils.NewTokenGenerator(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize token maker: %v", err)
	}

	//initialize repo
	userRepo := storage.NewUserRepository(db)

	//initialize service
	authService := auth.NewService(userRepo, tokenmaker, hashutils)

	//initialize router
	router := router.SetupRouter(authService)

	//start server
	router.Run(":" + cfg.ServerPort)
}
