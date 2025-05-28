package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DataDir		string
	BackupDir	string
	TasksFile	string
	MaxBackups	int
}

func New() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	dataDir := filepath.Join(homeDir, "todo")

	return &Config{
		DataDir: dataDir,
		BackupDir: filepath.Join(dataDir, "backups"),
		TasksFile: filepath.Join(dataDir, "tasks.json"),
		MaxBackups: 10,
	}
}