package cmd

import (
	"fmt"

	"github.com/samnart1/GoLang-Projects/004guessgame/internal/game"
	"github.com/samnart1/GoLang-Projects/004guessgame/internal/ui"
	"github.com/spf13/cobra"
)

var (
	difficulty	string
	timeLimit	int
	hints		bool
)

var playCmd = &cobra.Command{
	Use: "play",
	Short: "start a new guessing game",
	Long: "start a new guessing game with the specified difficulty levels",
	RunE: runPlay,
}

func init() {
	playCmd.Flags().StringVarP(&difficulty, "difficulty", "d", "medium", "Game difficulty (easy, medium, hard, custom)")
	playCmd.Flags().IntVarP(&timeLimit, "time", "t", 0, "Time limit in seconds (0 for no limit)")
	playCmd.Flags().BoolVarP(&hints, "hints", "h", true, "Enable hints")
}

func runPlay(cmd *cobra.Command, args []string) error {
	ui.ShowWelcome()

	g, err := game.New(difficulty, timeLimit, hints)
	if err != nil {
		return fmt.Errorf("failed to created game: %w", err)
	}

	return g.Play()
}