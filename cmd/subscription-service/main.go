package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chixxx1/subscription-service/internal/config"
	postgres_pool "github.com/chixxx1/subscription-service/internal/db/postgres"
	"github.com/chixxx1/subscription-service/internal/logger"
	sub_posgres_repo "github.com/chixxx1/subscription-service/internal/repository/postgres"
	sub_service "github.com/chixxx1/subscription-service/internal/service/subscription"
	transport_http "github.com/chixxx1/subscription-service/internal/transport/http"
	"go.uber.org/zap"
)

// @title        Subscription Service API
// @version      1.0
// @description  API for managing user subscriptions and calculating costs.
// @host         localhost:8080
// @BasePath     /api/v1
func main() {
	ctx := context.Background()

	logger, err := logger.NewLogger("dev")
	if err != nil {
		log.Fatalf("failed to initialized logger: %v", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("Starting application...")

	config, err := config.NewDBConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}
	logger.Info("Config load successfully")

	connPool, err := postgres_pool.NewConnectionPool(ctx, config)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer connPool.Close()
	logger.Info("Connected to PostgreSQL successfully")

	subRepo := sub_posgres_repo.NewSubscriptionRepo(connPool.Pool)

	subService := sub_service.NewSubscriptionService(subRepo, logger)

	router := transport_http.InitRoutes(subService, logger)

	logger.Info("Server is ready to accept connections")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal("HTTP server ListenAndServe", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
