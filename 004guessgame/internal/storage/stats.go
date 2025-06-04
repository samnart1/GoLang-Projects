package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Stats struct {
	GamesPlayed 	int `json:"games_played"`
	GamesWon		int	`json:"games_won"`
	TotalGuesses	int	`json:"total_guesses"`
}

func LoadStats() (*Stats, error) {
	path, err := statsPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Stats{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var stats Stats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}


func SaveStats(stats *Stats) error {
	path, err := statsPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(stats, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}


func ResetAll() error {
	if err := saveScores(Scores{}); err != nil {
		return err
	}

	return SaveStats(&Stats{})
}


func statsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "guess-game", "stats.json"), nil
}