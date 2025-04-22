package main

import (
	"context"
	"os/signal"
	_ "plata/docs"
	cU "plata/internal/app/cron"
	"plata/internal/app/postgres"
	"plata/internal/clients/exchange"
	"plata/internal/common/log"
	"plata/internal/config"
	qr "plata/internal/repository/quote"
	qs "plata/internal/services/quote"
	"plata/internal/transport/api"
	"syscall"
	"time"
)

// @title Quotes API
// @version 1.0
// @description Currency quote service: async updates and retrieval
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
// @schemes http

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := log.NewZapLogger()
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Errorf("Failed to load config: %v", err)
		return
	}
	db, err := postgres.NewPostgresDB(cfg.Postgres, logger)
	if err != nil {
		logger.Errorf("Failed to connect to DB: %v", err)
		return
	}
	defer db.Stop()

	repoQuote := qr.New(db.Primary(), db.Replica())
	exchClient := exchange.New(cfg.Exchange, logger)
	service := qs.New(repoQuote, exchClient, logger)
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	server := api.NewServer(logger)
	defer server.Stop(shutdownCtx)
	handler := api.NewHandler(service)
	if err = server.InitServer(handler); err != nil {
		logger.Errorf("Failed to initialize server: %v", err)
		return
	}
	quoteUpdater := cU.New(cfg.Cron, *repoQuote, exchClient, logger)
	defer quoteUpdater.Stop()
	if err = quoteUpdater.Run(); err != nil {
		logger.Errorf("Failed to start quote updater: %v", err)
		return
	}

	<-ctx.Done()
	logger.Info("Shutdown signal received")
	defer cancel()
}
