package cmd

import (
	"fmt"

	"github.com/samnart1/GoLang-Projects/003todo/internal/config"
	"github.com/samnart1/GoLang-Projects/003todo/internal/storage"
	"github.com/samnart1/GoLang-Projects/003todo/internal/task"
	"github.com/samnart1/GoLang-Projects/003todo/internal/ui"
)

func SearchTasks(cfg *config.Config, searchTerm string) {
	storage := storage.NewJSONStorage(cfg)

	tasks, err := storage.LoadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	filter := task.NewFilter()
	filter.SearchTerm = searchTerm

	results := filter.Apply(tasks)

	formatter := ui.NewTaskFormatter()
	fmt.Printf("Search results for %s:\n\n", ui.Yellow(`"`+searchTerm+`"`))

	if len(results) == 0 {
		fmt.Println("No tasks found matching the search term")
		return
	}

	fmt.Println(formatter.FormatTaskList(results))
	fmt.Printf("\nFound %d task(s)\n", len(results))
}