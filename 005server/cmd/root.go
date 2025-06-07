package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "go-http-server",
	Short: "A simple HTTP server built with go",
	Long: `A simple HTTP server built with Go that demonstrates web development concepts including routing, middleware, template rendering, and API endpoints`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(versionCmd)
}