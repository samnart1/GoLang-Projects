package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samnart1/GoLang-Projects/005server/internal/config"
	"github.com/samnart1/GoLang-Projects/005server/internal/logger"
	"github.com/samnart1/GoLang-Projects/005server/internal/server"
	"github.com/spf13/cobra"
)

var (
	port		string
	environment	string
	logLevel	string
)

var serverCmd = &cobra.Command{
	Use: "server",
	Short: "Start the HTTP server",
	Long: `Start the Http server with specified config`,
	Run: runServer,
}

func init() {
	serverCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to run the server on")
	serverCmd.Flags().StringVarP(&environment, "env", "e", "development", "Environment (development, production)")
	serverCmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (debug, info, warn, error)")
}

func runServer(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	if cmd.Flags().Changed("port") {
		cfg.Port = port
	}
	if cmd.Flags().Changed("env") {
		cfg.Environment = environment
	}
	if cmd.Flags().Changed("log-level") {
		cfg.LogLevel = logLevel
	}

	log := logger.New(cfg.LogLevel, cfg.Environment)

	srv := server.New(cfg, log)

	go func() {
		log.Info("Starting http server", "port", cfg.Port, "env", cfg.Environment)
		if err := srv.Start(); err != nil {
			log.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	//shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	log.Info("Server exited gracefully")
}