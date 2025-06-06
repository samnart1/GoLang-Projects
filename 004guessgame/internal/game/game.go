package game

import (
	"context"
	"fmt"
	"time"

	"github.com/samnart1/GoLang-Projects/004guessgame/internal/storage"
	"github.com/samnart1/GoLang-Projects/004guessgame/internal/timer"
	"github.com/samnart1/GoLang-Projects/004guessgame/internal/ui"
	"github.com/samnart1/GoLang-Projects/004guessgame/pkg/random"
)

type Game struct {
	target		int
	min			int
	max			int
	difficulty	Difficulty
	timeLimit	time.Duration
	timer 		*timer.Timer
	hints		bool
	guesses		int
	startTime	time.Time
	won			bool
}

func New(difficultyStr string, timeLimit int, hints bool) (*Game, error) {
	difficulty, err := ParseDifficulty(difficultyStr)
	if err != nil {
		return nil, err
	}

	min, max := difficulty.Range()
	target, err := random.IntRange(min, max)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random number: %w", err)
	}

	var timeLimitDuration time.Duration
	if timeLimit > 0 {
		timeLimitDuration = time.Duration(timeLimit) * time.Second
	}

	return &Game{
		target: 	target,
		min: 		min,
		max: 		max,
		difficulty: difficulty,
		timeLimit: 	timeLimitDuration,
		hints: 		hints,
		startTime: 	time.Now(),
	}, nil
}

func (g *Game) Play() error {
	ctx := context.Background()
	var cancel context.CancelFunc

	if g.timeLimit > 0 {
		ctx, cancel = context.WithTimeout(ctx, g.timeLimit)
		defer cancel()

		g.timer = timer.New(g.timeLimit)
		go g.timer.Start()
		defer func() {
			if g.timer != nil {
				g.timer.Stop()
			}
		}()
	}

	ui.ShowGameStart(g.min, g.max, g.difficulty.String(), g.timeLimit)

	for {
		select {
		case <-ctx.Done():
			ui.ShowTimeout()
			ui.ShowAnswer(g.target)
			return g.endGame()
		default:
			guess, err := ui.GetGuess(g.min, g.max)
			if err != nil {
				ui.ShowError(err)
				continue
			}
			
			g.guesses++

			if guess == g.target {
				g.won = true
				ui.ShowWin(g.guesses, time.Since(g.startTime))
				return g.endGame()
			}

			if g.hints {
				hint := g.generateHint(guess)
				ui.ShowHint(hint)
			} else {
				// still showing basic higher/lower without detailed hints
				if guess < g.target {
					ui.ShowHint("Too low!")
				} else {
					ui.ShowHint("Too high!")
				}
			}
		}
	}
}

func (g *Game) generateHint(guess int) string {
	diff := abs(g.target - guess)

	if guess < g.target {
		if diff <= 2 {
			return "Very close! Go higher."
		} else if diff <= 5 {
			return "Close! Go higher."
		} else if diff <= 10 {
			return "Go higher!!"
		} else {
			return "much higher"
		}
	} else {
		if diff <= 2 {
			return "Very close! Go lower."
		} else if diff <= 5 {
			return "Close! Go lower."
		} else if diff <= 10 {
			return "Goo lower!"
		} else {
			return "Much lower"
		}
	}
}

func (g *Game) endGame() error {
	duration := time.Since(g.startTime)

	stats, err := storage.LoadStats()
	if err != nil {
		stats = &storage.Stats{}
	}

	stats.GamesPlayed++
	if g.won {
		stats.GamesWon++
		stats.TotalGuesses += g.guesses

		score := &storage.Score{
			Difficulty: g.difficulty.String(),
			Guesses: 	g.guesses,
			Time: 		duration,
			Date: 		time.Now(),
		}

		if err := storage.SaveScores(score); err != nil {
			return err
		}
	}

	return storage.SaveStats(stats)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}