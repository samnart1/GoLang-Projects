package cmd

import (
	"fmt"

	"github.com/samnart1/GoLang-Projects/004guessgame/internal/storage"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use: "reset",
	Short: "Reset game data",
	Long: "Reset high scores and statistics",
	RunE: runReset,
}

func runReset(cmd *cobra.Command, args []string) error {
	if err := storage.ResetAll(); err != nil {
		return err
	}
	fmt.Println("Game data reset successfully!")
	return nil
}