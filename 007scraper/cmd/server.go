package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/samnart1/golang/007scrapper/internal/server"
	"github.com/spf13/cobra"
)

var (
	serverPort string
	serverHost string
)

var serverCmd = &cobra.Command{
	Use: "server",
	Short: "Start HTTP server for web scraping API",
	Long: `Start an HTTP server that provides a REST API for web scraping.
	The server provides the following endpoints:
		GET 	/health			- Health check
		POST 	/scrape			- Scrape a single URL
		POST	/scrape/batch	- scrape multiple URLs
		
	Examples:
		webscraper server
		webscraper server --port 8080 --host localhost`,

	RunE: runServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&serverPort, "port", "p", "", "Server port (default: from config)")
	serverCmd.Flags().StringVar(&serverHost, "host", "", "Server host (default: from config)")
}

func runServer(cmd *cobra.Command, args []string) error {
	host := cfg.ServerHost
	if serverHost != "" {
		host = serverHost
	}

	port := cfg.ServerPort
	if serverPort != "" {
		port = serverPort
	}

	address := host + ":" + port

	srv := server.New(cfg)

	fmt.Printf("Starting web scraper server on %s\n", address)
	fmt.Printf("Available endpoints:\n")
	fmt.Printf("	GET  /health		- Health check/n")
	fmt.Printf("	POST /scrape		- Scrape a single URL\n")
	fmt.Printf("	POST /scrape/batch	- Scrape multiple URLs\n")
	fmt.Printf("\nPress Ctrl+C to stop the server\n\n")

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Println("\nShutting down server...")
		os.Exit(0)
	}()

	if err := http.ListenAndServe(address, srv.Router()); err != nil {
		return fmt.Errorf("server failed to start: %v", err)
	}

	return nil
}