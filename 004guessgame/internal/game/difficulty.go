package game

import (
	"fmt"
	"strings"
)

type Difficulty int


const (
	Easy Difficulty = iota
	Medium
	Hard
	Custom
)

func (d Difficulty) String() string {
	switch d {
	case Easy:
		return "easy"
	case Medium:
		return "medium"
	case Hard:
		return "hard"
	case Custom:
		return "custom"
	default:
		return "unknown"
	}
}

func (d Difficulty) Range() (int, int) {
	switch d {
	case Easy:
		return 1, 10
	case Medium:
		return 1, 50
	case Hard:
		return 1, 100
	case Custom:
		return 1, 1000
	default:
		return 1, 10
	}

}

func ParseDifficulty(s string) (Difficulty, error) {
	switch strings.ToLower(s) {
	case "easy", "e":
		return Easy, nil
	case "medium", "m":
		return Medium, nil
	case "hard", "h":
		return Hard, nil
	case "custom", "c":
		return Custom, nil
	default:
		return Easy, fmt.Errorf("unknown difficulty: %s", s)
	}
}
