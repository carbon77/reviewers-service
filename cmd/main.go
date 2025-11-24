package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"reviewers/internal/config"
	"reviewers/internal/db"
	"reviewers/internal/handler"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := slog.Default()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	conn, err := db.Connect(cfg)
	if err != nil {
		logger.Error("Failed to connect to database")
		os.Exit(1)
	}
	defer func() {
		sqlDB, _ := conn.DB()
		sqlDB.Close()
		logger.Info("Databse connection closed")
	}()

	router := gin.Default()
	handler.InitHandlers(logger, conn, router)

	addr := fmt.Sprintf(":%d", cfg.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("Starting HTTP server", "address", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.Error("Server error", "error", err)

	case <-quit:
		logger.Info("Shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Shutdown failed", "error", err)
			server.Close()
		} else {
			logger.Info("Shutdown completed")
		}
	}
}
