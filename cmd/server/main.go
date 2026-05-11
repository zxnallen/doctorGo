package main

import (
	"log"

	"doctor-go/internal/app"
	"doctor-go/internal/config"
)

// @title Doctor Go API
// @version 1.0
// @description 肿瘤医生项目后端接口文档
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	server, err := app.NewServer(cfg)
	if err != nil {
		log.Fatalf("create server: %v", err)
	}

	if err := server.Run(); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
