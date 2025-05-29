package cmd

import (
	"fmt"

	"github.com/samnart1/GoLang-Projects/003todo/internal/config"
	"github.com/samnart1/GoLang-Projects/003todo/internal/storage"
	"github.com/samnart1/GoLang-Projects/003todo/internal/task"
	"github.com/samnart1/GoLang-Projects/003todo/internal/ui"
)

func CompleteTask(cfg *config.Config, id int) {
	storage := storage.NewJSONStorage(cfg)

	tasks, err := storage.LoadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	manager := task.NewManager()
	manager.LoadTasks(tasks)

	if err := manager.CompleteTask(id); err != nil {
		fmt.Printf("Error completing task: %v\n", err)
		return
	}

	if err := storage.SaveTasks(manager.GetTasks()); err != nil {
		fmt.Printf("Error saving tasks: %v", err)
		return
	}

	completedTask, _ := manager.GetTaskByID(id)

	formatter := ui.NewTaskFormatter()
	fmt.Printf("Completed: %s\n", ui.Green("✓"))
	fmt.Println(formatter.FormatTask(completedTask))
}