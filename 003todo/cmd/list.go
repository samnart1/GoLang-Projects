package cmd

import (
	"fmt"

	"github.com/samnart1/GoLang-Projects/003todo/internal/config"
	"github.com/samnart1/GoLang-Projects/003todo/internal/storage"
	"github.com/samnart1/GoLang-Projects/003todo/internal/task"
	"github.com/samnart1/GoLang-Projects/003todo/internal/ui"
)

func ListTasks(cfg *config.Config, args []string) {
	storage := storage.NewJSONStorage(cfg)

	tasks, err := storage.LoadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	filter := task.NewFilter()
	tableFormat := false
	sortBy := "created"

	for i, arg := range args {
		switch arg { 
		case "--pending":
			filter.SetPendingOnly()
		case "--completed":
			filter.SetCompletedOnly()
		case "--priority":
			if i+1 < len(args) {
				if priority, err := task.ParsePriority(args[i+1]); err == nil {
					filter.Priority = &priority
				}
			}
		case "--table":
			tableFormat = true
		case "--sort":
			if i+1 < len(args) {
				sortBy = args[i+1]
			}
		}
	}

	filteredTasks := filter.Apply(tasks)

	manager := task.NewManager()
	manager.LoadTasks(filteredTasks)

	switch sortBy {
	case "priority":
		manager.SortByPriority()
	case "due":
		manager.SortByDueDate()
	default:
		manager.SortByCreated()
	}

	if tableFormat {
		tableFormatter := ui.NewTableFormatter()
		fmt.Println(tableFormatter.FormatTable(manager.GetTasks()))
	} else {
		formatter := ui.NewTaskFormatter()
		fmt.Println(formatter.FormatTaskList(manager.GetTasks()))
	}

	if len(filteredTasks) != len(tasks) {
		fmt.Printf("\nShowing %d of %d tasks\n", len(filteredTasks), len(tasks))
	}
}