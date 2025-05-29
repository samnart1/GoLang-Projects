package config

import (
	"path/filepath"
	"time"
)

func (c *Config) GetBackupPath() string {
	timestamp := time.Now().Format("20060102_150405")
	filename := "tasks_" + timestamp + ".json"
	return filepath.Join(c.BackupDir, filename)
}

func (c *Config) GetTempPath() string {
	return c.TasksFile + ".temp"
}