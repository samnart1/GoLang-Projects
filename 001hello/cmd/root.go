package cmd

import (
	"github.com/samnart/GoLang-projects/001hello/internal/app"
	"github.com/samnart/GoLang-projects/001hello/pkg/version"
	"github.com/spf13/cobra"
)

var (
	verbose	bool
	name	string
)

var rootCmd = &cobra.Command{
	Use: "001hello",
	Short: "A production-ready Hello World CLI application",
	Long: `A well-structured Go CLI application that demonstrates production best practices`,
	Version: version.Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := app.Config{
			Name: name,
			Verbose: verbose,
		}

		helloApp := app.New(config)
		return helloApp.Run()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().StringVarP(&name, "name", "n", "World", "Name to greet")
}