package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samnart1/GoLang-Projects/005server/cmd"
	"github.com/samnart1/GoLang-Projects/005server/internal/config"
	"github.com/samnart1/GoLang-Projects/005server/internal/logger"
	"github.com/samnart1/GoLang-Projects/005server/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logger.New(cfg.LogLevel, cfg.Environment)

	// Check if we're running a CLI command
	if len(os.Args) > 1 {
		cmd.Execute()
		return
	}

	// Create and start HTTP server
	srv := server.New(cfg, logger)

	// Start server in a goroutine
	go func() {
		logger.Info("Starting HTTP server", "port", cfg.Port, "env", cfg.Environment)
		if err := srv.Start(); err != nil {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Server exited gracefully")
}