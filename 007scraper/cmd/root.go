package cmd

import (
	"fmt"

	"github.com/samnart1/golang/007scrapper/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfg *config.Config
)

var rootCmd = &cobra.Command{
	Use: "webscraper",
	Short: "A powerful web scraper built with Go",
	Long: `A comprehensive web scraper that can fetch and extract information from web pages.
	Features:
		- Extract titles, descriptions, keywords, links and images
		- Batch processing of multiple URLs
		- HTTP server mode for API access
		- Configurable timeouts and user agents
		- Support for custom headers`,
	Version: "1.0.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.webscraper.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

var cfgFile string
var verbose bool

func initConfig() {
	cfg = config.Load()

	if verbose {
		fmt.Printf("Configuration loaded:\n")
		fmt.Printf("	Timeout: %v\n", cfg.Timeout)
		fmt.Printf("	User Agent: %s\n", cfg.UserAgent)
		fmt.Printf("	Server: %s\n", cfg.ServerAddress())
		fmt.Printf("	Max Concurrent: %d\n", cfg.MaxConcurrent)
	}
}