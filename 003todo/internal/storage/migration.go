package storage

import (
	"encoding/json"
	"os"

	"github.com/samnart1/GoLang-Projects/003todo/internal/task"
)

type LegacyTask struct {
	ID			int		`json:"id"`
	Description	string	`json:"description"`
	Done		bool	`json:"done"`
}

func (s *JSONStorage) MigrateLegacyData() error {
	data, err := os.ReadFile(s.config.TasksFile)
	if err != nil {
		return err
	}

	var legacyTasks []LegacyTask
	if err := json.Unmarshal(data, &legacyTasks); err != nil {
		return nil		// not legacy format, not needed
	}

	var newTasks []*task.Task
	for _, legacy := range legacyTasks {
		newTask := task.NewTask(legacy.ID, legacy.Description)
		if legacy.Done {
			newTask.Complete()
		}
		newTasks = append(newTasks, newTask)
	}

	return s.SaveTasks(newTasks)
}