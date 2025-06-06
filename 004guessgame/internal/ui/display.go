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
	fmt.Println("Let's see how many guesses it takes you!")
	fmt.Println()
}

func ShowHint(hint string) {
	fmt.Printf("%s\n", hint)
	fmt.Println()
}

func ShowWin(guesses int, duration time.Duration) {
	fmt.Printf("Congratulations! You guessed it in %d atempts!\n", guesses)
	fmt.Printf("Time taken %v\n", duration.Round(time.Second))

	switch {
	case guesses == 1:
		fmt.Println("Incredible! You got it on the first try!")
	case guesses <= 3:
		fmt.Println("Excellent guessing!")
	case guesses <= 5:
		fmt.Println("Good job!")
	case guesses <= 10:
		fmt.Println("Not bad!")
	default:
		fmt.Println("Better luck next time!")
	}
	fmt.Println()
}

func ShowTimeout() {
	fmt.Println("Time's up! Better luck next time!")
	fmt.Println()
}

func ShowAnswer(answer int) {
	fmt.Printf("The number was: %d\n", answer)
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

	fmt.Printf("Current Streak: %d\n", stats.CurrentStreak)
	fmt.Printf("Best Streak: %d\n", stats.BestStreak)

	fmt.Println()

	if len(scores) > 0 {
		fmt.Println("Recent High Scores")
		fmt.Println("==================")

		displayCount := 5
		if len(scores) < displayCount {
			displayCount = len(scores)
		}

		for i := 0; i < displayCount; i++ {
			score := scores[i]
			fmt.Printf("%s: %d guesses in %v (%s)\n",
				score.Difficulty,
				score.Guesses,
				score.Time.Round(time.Second),
				score.Date.Format("20006-01-02"),
			)
		}
		fmt.Println()
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