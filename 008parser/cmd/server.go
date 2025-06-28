package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/samnart1/golang/008parser/internal/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use: "server",
	Short: "Start HTTP server for JSON parsing",
	Long: `Start an HTTP server that provides JSON parsing, validation, and formatting endpoints
		
		Available endpoints:
		- POST /parse - Parse JSON data
		- POST /validate - Validate JSON data
		- POST /format - Format JSON data
		- GET /health - Health check`,

	RunE: runServer,
}

var (
	serverHost string
	serverPort int
	enableCORS bool
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVar(&serverHost, "host", "localhost", "Server host")
	serverCmd.Flags().IntVar(&serverPort, "port", 8080, "Server port")
	serverCmd.Flags().BoolVar(&enableCORS, "cors", false, "Enable CORS headers")
}

func runServer(cmd *cobra.Command, args []string) error {
	srv := server.New(cfg)

	httpServer := &http.Server{
		Addr: serverHost + ":" + strconv.Itoa(serverPort),
		Handler: 		srv.Handler(enableCORS),
		ReadTimeout: 	30 * time.Second,
		WriteTimeout: 	30 * time.Second,
		IdleTimeout: 	60 * time.Second,
	}

	go func() {
		fmt.Printf("Starting JSON Parser server on http://%s:%d\n", serverHost, serverPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<- quit

	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	fmt.Println("Server stopped")
	return nil
}