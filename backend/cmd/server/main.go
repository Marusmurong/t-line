package main

import (
	"log"

	"github.com/t-line/backend/internal/config"
	"github.com/t-line/backend/internal/pkg/database"
	"github.com/t-line/backend/internal/pkg/logger"
	"github.com/t-line/backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger.Init(cfg.Log)

	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	rdb, err := database.NewRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	srv := server.New(cfg, db, rdb)
	if err := srv.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
