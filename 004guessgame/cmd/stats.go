package cmd

import (
	"github.com/samnart1/GoLang-Projects/004guessgame/internal/storage"
	"github.com/samnart1/GoLang-Projects/004guessgame/internal/ui"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use: "stats",
	Short: "Show game statistics",
	Long: "Display your game statistics including wins, average guesses, and high scores",
	RunE: runStats,
}

func runStats(cmd *cobra.Command, args []string) error {
	stats, err := storage.LoadStats()
	if err != nil {
		return err
	}

	scores, err := storage.LoadScores()
	if err != nil {
		return err
	}

	ui.ShowStats(stats, scores)
	return nil
}