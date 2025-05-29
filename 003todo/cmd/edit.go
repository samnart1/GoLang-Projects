package cmd

import (
	"fmt"
	"strings"

	"github.com/samnart1/GoLang-Projects/003todo/internal/config"
	"github.com/samnart1/GoLang-Projects/003todo/internal/storage"
	"github.com/samnart1/GoLang-Projects/003todo/internal/task"
	"github.com/samnart1/GoLang-Projects/003todo/internal/ui"
)

func EditTask(cfg *config.Config, id int, args []string) {
	storage := storage.NewJSONStorage(cfg)

	tasks, err := storage.LoadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	manager := task.NewManager()
	manager.LoadTasks(tasks)

	newDescription := strings.Join(args, " ")
	if err := manager.EditTask(id, newDescription); err != nil {
		fmt.Printf("Error editing task: %v\n", err)
		return
	}

	if err := storage.SaveTasks(manager.GetTasks()); err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
		return
	}

	editedTask, _ := manager.GetTaskByID(id)
	
	formatter := ui.NewTaskFormatter()
	fmt.Printf("Edited: %s\n", ui.Blue("✎"))
	fmt.Println(formatter.FormatTask(editedTask))

}