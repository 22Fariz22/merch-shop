package main

import (
	"log"

	"github.com/22Fariz22/merch-shop/config"
	"github.com/22Fariz22/merch-shop/internal/server"
	"github.com/22Fariz22/merch-shop/pkg/db/postgres"
	"github.com/22Fariz22/merch-shop/pkg/db/redis"
	"github.com/22Fariz22/merch-shop/pkg/logger"
)

func main() {
	log.Println("Starting api server")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis connected")

	s := server.NewServer(cfg, psqlDB, redisClient, appLogger)
	if err = s.Run(); err != nil {
		appLogger.Fatalf("Error in main NewServer(): ", err)
	}
}
