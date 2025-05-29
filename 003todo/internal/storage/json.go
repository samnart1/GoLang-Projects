package storage

import (
	"encoding/json"
	"os"

	"github.com/samnart1/GoLang-Projects/003todo/internal/config"
	"github.com/samnart1/GoLang-Projects/003todo/internal/task"
	"github.com/samnart1/GoLang-Projects/003todo/pkg/errors"
)

type JSONStorage struct {
	config *config.Config
}

func NewJSONStorage(cfg *config.Config) *JSONStorage {
	return &JSONStorage{
		config: cfg,
	}
}


func (s *JSONStorage) LoadTasks() ([]*task.Task, error) {
	if err := s.config.EnsureDirectories(); err != nil {
		return nil, errors.NewStorageError(s.config.DataDir, "create directories", err)
	}

	if _, err := os.Stat(s.config.TasksFile); os.IsNotExist(err) {
		return []*task.Task{}, nil
	}

	data, err := os.ReadFile(s.config.TasksFile)
	if err != nil {
		return nil, errors.NewStorageError(s.config.TasksFile, "read", err)
	}

	var tasks []*task.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, errors.NewStorageError(s.config.TasksFile, "unmarshal", err)
	}

	return tasks, nil
}


func (s *JSONStorage) SaveTasks(tasks []*task.Task) error {
	if err := s.config.EnsureDirectories(); err != nil {
		return errors.NewStorageError(s.config.DataDir, "create directories", err)
	}

	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return errors.NewStorageError(s.config.TasksFile, "marshal", err)
	}

	tempFile := s.config.GetTempPath()
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return errors.NewStorageError(tempFile, "write temp", err)
	}

	if err := os.Rename(tempFile, s.config.TasksFile); err != nil {
		os.Remove(tempFile)		// clean tempfile
		return errors.NewStorageError(s.config.TasksFile, "rename", err)
	}

	return nil
}