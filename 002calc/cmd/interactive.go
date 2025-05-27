package cmd

import (
	"github.com/samnart1/GoLang-Projects/tree/main/002calc/internal/ui"
	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command {
	Use: "interactive",
	Aliases: []string{"i", "repl"},
	Short: "Start interactive calculator mode",
	Long: `Start the calculator in interactive mode where you can perform multiple calculations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ui.InteractiveMode(verbose)
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}