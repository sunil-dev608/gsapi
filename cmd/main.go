package main

import (
	"context"
	"gsapi/config"
	"gsapi/internal/server"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	err = godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file:", zap.Error(err))
	}
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Fatal("Error loading .env file:", zap.Error(err))
	}

	shuttingDown := false
	server := server.New(cfg)
	go func() {
		if err := server.Start(); err != nil {
			if !shuttingDown {
				logger.Fatal("failed to start server", zap.Error(err))
			}
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	shuttingDown = true
	logger.Warn("Receiving signal to Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("server shutdown failed", zap.Error(err))
	}

}
