package cmd

import (
	"github.com/samnart1/GoLang-Projects/tree/main/001hello/internal/app"
	"github.com/samnart1/GoLang-Projects/tree/main/001hello/pkg/version"
	"github.com/spf13/cobra"
)

var (
	name 	string
	verbose bool
)

var rootCmd = &cobra.Command {
	Use: 		"001hello",
	Short: 		"This is a simple but well structured GoLang Project!",
	Long: 		"This project showcases how go projects are structured and follows a production ready norm or something",
	Version: 	version.Version,
	RunE:		func(cmd *cobra.Command, args []string) error {
					config := app.Config {
						Name: 		name,
						Verbose: 	verbose,
					}

					helloApp := app.New(config)
					return helloApp.Run()
				},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.Flags().StringVarP(&name, "name", "n", "World", "name for greeting")
}