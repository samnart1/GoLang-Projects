package cmd

import (
	"fmt"

	"github.com/samnart1/GoLang-Projects/003todo/internal/config"
	"github.com/samnart1/GoLang-Projects/003todo/internal/storage"
	"github.com/samnart1/GoLang-Projects/003todo/internal/task"
	"github.com/samnart1/GoLang-Projects/003todo/internal/ui"
)

func RemoveTask(cfg *config.Config, id int) {
	storage := storage.NewJSONStorage(cfg)

	tasks, err := storage.LoadTasks()
	if err != nil {
		fmt.Printf("Error loadng tasks: %v\n", err)
		return
	}

	manager := task.NewManager()
	manager.LoadTasks(tasks)

	taskToRemove, err := manager.GetTaskByID(id)
	if err != nil {
		fmt.Printf("Error finding task: %v\n", err)
		return
	}

	if err := manager.RemoveTask(id); err != nil {
		fmt.Printf("Error removing task: %v\n", err)
		return
	}

	if err := storage.SaveTasks(manager.GetTasks()); err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
		return
	}

	formatter := ui.NewTaskFormatter()
	fmt.Printf("Removed: %s\n", ui.Red("✓"))
	fmt.Println(formatter.FormatTask(taskToRemove))
}