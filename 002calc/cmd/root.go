package cmd

import (
	"github.com/samnart1/GoLang-Projects/002calc/pkg/version"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

var rootCmd = &cobra.Command {
	Use:	"calc",
	Short:	"A production-ready calculator CLI application",
	Long: 	`A powerful calculator that supports basic arithmetic operations with both direct calculation and interactive modes.`,
	Version:	version.Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}