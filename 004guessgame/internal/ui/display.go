package ui

import (
	"fmt"
	"time"

	"github.com/samnart1/GoLang-Projects/004guessgame/internal/config"
	"github.com/samnart1/GoLang-Projects/004guessgame/internal/storage"
)

func ShowWelcome() {
	fmt.Println("Welcome to the Number Guessing Game!")
	fmt.Println("====================================")
	fmt.Println()
}

func ShowGameStart(min, max int, difficulty string, timeLimit time.Duration) {
	fmt.Printf("Starting %s game (Range: %d-%d)\n", difficulty, min, max)
	if timeLimit > 0 {
		fmt.Printf("Time limit: %v\n", timeLimit)
	}
	fmt.Printf("I'm thinking of a number between %d and %d\n", min, max)
	fmt.Println()
}

func ShowHint(hint string) {
	fmt.Printf("%s\n", hint)
	fmt.Println()
}

func ShowWin(guesses int, duration time.Duration) {
	fmt.Printf("Congratulations! You guessed it in %d atempts!\n", guesses)
	fmt.Printf("Time taken %v\n", duration)
	fmt.Println()
}

func ShowTimeout() {
	fmt.Println("Time's up! Better luck next time!")
	fmt.Println()
}

func ShowError(err error) {
	fmt.Printf("Error: %v\n", err)
	fmt.Println()
}

func ShowStats(stats *storage.Stats, scores storage.Scores) {
	fmt.Println("Game Statistics")
	fmt.Println("================")
	fmt.Printf("Games Played: %d\n", stats.GamesPlayed)
	fmt.Printf("Games Won: %d\n", stats.GamesWon)

	if stats.GamesPlayed > 0 {
		winRate := float64(stats.GamesWon) / float64(stats.GamesPlayed) * 100
		fmt.Printf("Win Rate: %.1f%%\n", winRate)
	}

	if stats.GamesWon > 0 {
		avgGuesses := float64(stats.TotalGuesses) / float64(stats.GamesWon)
		fmt.Printf("Average Guesses: %.1f\n", avgGuesses)
	}

	fmt.Println()

	if len(scores) > 0 {
		fmt.Println("Recent High Scores")
		fmt.Println("==================")

		start := len(scores) - 5
		if start < 0 {
			start = 0
		}

		for i := len(scores) - 1; i >= start; i-- {
			score := scores[i]
			fmt.Printf("%s: %d guesses in %v (%s)\n",
				score.Difficulty,
				score.Guesses,
				score.Time.Round(time.Second),
				score.Date.Format("20006-01-02"),
			)
		}
	}
}

func ShowConfig(cfg *config.Config) {
	fmt.Println("Configuration")
	fmt.Println("=============")
	fmt.Printf("Default Difficulty: %s\n", cfg.DefaultDifficulty)
	fmt.Printf("Colors Enabled: %t\n", cfg.EnableColors)
	fmt.Printf("Sound Enabled: %t\n", cfg.EnableSound)
	fmt.Printf("Default Time Limit: %d seconds\n", cfg.DefaultTimeLimit)
}