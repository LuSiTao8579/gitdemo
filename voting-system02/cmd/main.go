package main

import (
	"log"
	"net/http"

	"voting-system/internal/config"
	"voting-system/internal/router"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化路由
	r := router.SetupRouter()

	// 启动服务器
	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
