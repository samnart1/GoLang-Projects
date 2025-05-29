package storage

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/samnart1/GoLang-Projects/003todo/pkg/errors"
)

func (s *JSONStorage) CreateBackup() error {
	if _, err := os.Stat(s.config.TasksFile); os.IsNotExist(err) {
		return nil
	}

	if err := s.config.EnsureDirectories(); err != nil {
		return errors.NewStorageError(s.config.BackupDir, "create backup dir", err)
	}

	backPath := s.config.GetBackupPath()

	if err := copyFile(s.config.TasksFile, backPath); err != nil {
		return errors.NewStorageError(backPath, "create backup ", err)
	}

	if err := s.cleanupOldBackups(); err != nil {
		return nil
	}

	return nil
}


func (s *JSONStorage) cleanupOldBackups() error {
	files, err := filepath.Glob(filepath.Join(s.config.BackupDir, "tasks_*.json"))
	if err != nil {
		return err
	}

	if len(files) <= s.config.MaxBackups {
		return nil
	}

	filesToRemove := len(files) - s.config.MaxBackups
	for i := 0; i < filesToRemove; i++ {
		if err := os.Remove(files[i]); err != nil {
			return err
		}
	}
	return nil
}


func (s *JSONStorage) ListBackups() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(s.config.BackupDir, "tasks_*.json"))
	if err != nil {
		return nil, err
	}

	var backups []string
	for _, file := range files {
		basename := filepath.Base(file)
		if strings.HasPrefix(basename, "tasks_") && strings.HasSuffix(basename, ".json") {
			timestamp := strings.TrimPrefix(basename, "tasks_")
			timestamp = strings.TrimSuffix(timestamp, ".json")
			backups = append(backups, timestamp)
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(backups)))
	return backups, nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}