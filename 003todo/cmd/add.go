package cmd

import (
	"fmt"
	"strings"

	"github.com/samnart1/GoLang-Projects/003todo/internal/config"
	"github.com/samnart1/GoLang-Projects/003todo/internal/storage"
	"github.com/samnart1/GoLang-Projects/003todo/internal/task"
	"github.com/samnart1/GoLang-Projects/003todo/internal/ui"
)

func AddTask(cfg *config.Config, args []string) {
	storage := storage.NewJSONStorage(cfg)

	tasks, err := storage.LoadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	manager := task.NewManager()
	manager.LoadTasks(tasks)

	description := strings.Join(args, " ")

	newTask, err := manager.AddTask(description)
	if err != nil {
		fmt.Printf("Error adding task: %v\n", err)
		return
	}

	if err := storage.SaveTasks(manager.GetTasks()); err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
		return
	}

	formatter := ui.NewTaskFormatter()
	fmt.Printf("Added: %s\n", ui.Green("✓"))
	fmt.Println(formatter.FormatTask(newTask))
}