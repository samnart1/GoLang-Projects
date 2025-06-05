package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use: "server",
	Short: "A simple server",
	Long: "A simple server in golang",
}

func Execute() error {
	return rootCmd.Execute()
}