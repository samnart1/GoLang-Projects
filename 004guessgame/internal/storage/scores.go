package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Score struct {
	Difficulty 	string			`json:"difficulty"`
	Guesses		int				`json:"guesses"`
	Time		time.Duration	`json:"time"`
	Date		time.Time		`json:"date"`
}

type Scores []Score

func LoadScores() (Scores, error) {
	path, err := scoresPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Scores{}, nil
	}

	data, err := os.ReadFile(path);
	if err != nil {
		return nil, err
	}

	var scores Scores
	if err := json.Unmarshal(data, &scores); err != nil {
		return nil, err
	}

	return scores, nil
}


func SaveScores(score *Score) error {
	scores, err := LoadScores()
	if err != nil {
		return err
	}

	scores = append(scores, *score)
	return saveScores(scores)
}


func saveScores(scores Scores) error {
	path, err := scoresPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(scores, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}


func scoresPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "guess-game", "score.json"), nil
}