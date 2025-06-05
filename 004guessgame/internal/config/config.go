package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultDifficulty	string	`json:"default_difficulty"`
	EnableColors		bool	`json:"enable_colors"`
	EnableSound			bool	`json:"enable_sound"`
	DefaultTimeLimit	int		`json:"default_time_limit"`
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return defaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	} 

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	// create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func defaultConfig() *Config {
	return &Config{
		DefaultDifficulty: "medium",
		EnableColors: true,
		EnableSound: false,
		DefaultTimeLimit: 0,
	}
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "guess-game", "config.json"), nil
}