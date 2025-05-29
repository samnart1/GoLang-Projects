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

	dataDir := filepath.Join(homeDir, ".todo")

	return &Config{
		DataDir: dataDir,
		BackupDir: filepath.Join(homeDir, "backups"),
		TasksFile: filepath.Join(homeDir, "tasks.json"),
		MaxBackups: 10,
	}
}

func (c *Config) EnsureDirectories() error {
	dirs := []string{c.DataDir, c.BackupDir}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 07555); err != nil {
			return err
		}
	}
	
	return nil
}