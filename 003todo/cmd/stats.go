package cmd

import (
	"fmt"
	"strings"

	"github.com/samnart1/GoLang-Projects/003todo/internal/config"
	"github.com/samnart1/GoLang-Projects/003todo/internal/storage"
	"github.com/samnart1/GoLang-Projects/003todo/internal/task"
	"github.com/samnart1/GoLang-Projects/003todo/internal/ui"
)

func ShowStat(cfg *config.Config) {
	storage := storage.NewJSONStorage(cfg)

	tasks, err := storage.LoadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	manager := task.NewManager()
	manager.LoadTasks(tasks)
	stats := manager.GetStats()

	completionRate := float64(0)
	if stats["total"] > 0 {
		completionRate = float64(stats["completed"])
	}

	fmt.Printf("%s\n", ui.Bold("Task Statistics"))
	fmt.Println(strings.Repeat("=", 20))
	fmt.Printf("Total tasks:		%s\n", ui.Blue(fmt.Sprintf("%d", stats["total"])))
	fmt.Printf("Completed:			%s\n", ui.Green(fmt.Sprintf("%d", stats["completed"])))
	fmt.Printf("Pending:			%s\n", ui.Yellow(fmt.Sprintf("%d", stats["pending"])))
	fmt.Printf("Overdue:			%s\n", ui.Red(fmt.Sprintf("%d", stats["overdue"])))
	fmt.Printf("Completion rate:	%s\n", ui.Cyan(fmt.Sprintf("%.1f%%", completionRate)))

	fmt.Printf("\n%s\n", ui.Bold("Priority Breakdown"))
	fmt.Println(strings.Repeat("-", 20))

	priorityStats := make(map[task.Priority]int)
	for _, t := range tasks {
		if !t.Completed {
			priorityStats[t.Priority]++
		}
	}

	fmt.Printf("High priority:		%s\n", ui.Red(fmt.Sprintf("%d", priorityStats[task.High])))
	fmt.Printf("Medium priority:	%s\n", ui.Yellow(fmt.Sprintf("%d", priorityStats[task.Medium])))
	fmt.Printf("Low prriority:		%s\n", ui.Green(fmt.Sprintf("%d", priorityStats[task.Low])))
}