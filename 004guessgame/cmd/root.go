package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use: "guess-game",
	Short: "A number guessing game",
	Long: "A CLI number guessing game with multple difficulty levels",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(configCmd)
	// rootCmd.AddCommand(playCmd)
	// rootCmd.AddCommand(resetCmd)
	// rootCmd.AddCommand(statsCmd)
}