package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samnart1/golang/009l9gger/internal/config"
	"github.com/samnart1/golang/009l9gger/internal/logger"
	"github.com/samnart1/golang/009l9gger/internal/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use: "server",
	Short: "Start HTTP server for remote logging",
	Long: `Start an Http server that accepts log messages via REST API.
	
		The server provides endpoints for:
		- POST /log - Log a message
		- GET /health - Health check
		- GET /logs - View recent logs (if enabled)
		
		Example:
			go-simple-logger server --port 8080
			curl -X POST http://localhost:8080/log -d '{"message": "Hello", "level": "info"}'`,

	RunE: runServer,		
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().String("host", "localhost", "server host")
	serverCmd.Flags().Int("port", 8080, "server port")
	serverCmd.Flags().Duration("timeout", 30*time.Second, "server timeout")
}

func runServer(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetInt("port")
	timeout, _ := cmd.Flags().GetDuration("timeout")

	cfg.Server.Host = host
	cfg.Server.Port = port

	l, err := logger.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer l.Close()

	srv := server.New(cfg, l)

	httpServer := &http.Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
		Handler: srv.Router(),
		ReadTimeout: timeout,
		WriteTimeout: timeout,
	}

	go func(){
		fmt.Printf("Starting server on %s:%d", host, port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<- quit

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("serer forced to shutdown: %w", err)
	}

	fmt.Println("Server stopped")
	return nil
}