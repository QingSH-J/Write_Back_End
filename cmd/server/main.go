package main

import (
	"log"

	"github.com/jinxinyu/go_backend/internal/config"
	"github.com/jinxinyu/go_backend/internal/router"
	"github.com/jinxinyu/go_backend/internal/storage"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. 初始化数据库连接
	db, err := storage.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 3. 设置路由
	r := router.SetupRouter(db)

	// 4. 启动服务器
	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
