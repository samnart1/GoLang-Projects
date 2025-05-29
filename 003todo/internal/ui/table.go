package ui

import (
	"fmt"
	"strings"

	"github.com/samnart1/GoLang-Projects/003todo/internal/task"
)

type TableFormatter struct {
	colorOutput bool
}

func NewTableFormatter() *TableFormatter {
	return &TableFormatter{
		colorOutput: true,
	}
}

func (f *TableFormatter) FormatTable(tasks []*task.Task) string {
	if len(tasks) == 0 {
		return "No tasks found!"
	}

	idWidth := 3
	statusWidth := 6
	priorityWidth := 8
	descWidth := 30
	dueDateWidth := 12
	tagsWidth := 15

	for _, t := range tasks {
		if len(t.Description) > descWidth {
			descWidth = len(t.Description)
			if descWidth > 60 {
				descWidth = 60
			}
		}
	}

	var lines []string

	header := f.formatRow(
		"ID", "Status", "Priority", "Description", "Due Date", "Tags",
		idWidth, statusWidth, priorityWidth, descWidth, dueDateWidth, tagsWidth,
	)
	if f.colorOutput {
		header = Bold(header)
	}
	lines = append(lines, header)

	separator := strings.Repeat("-", idWidth+statusWidth+priorityWidth+descWidth+dueDateWidth+tagsWidth+15)
	lines = append(lines, separator)

	for _, t := range tasks {
		status := "Pending"
		if t.Completed {
			status = "Done"
		}

		desc := t.Description
		if len(desc) > descWidth {
			desc = desc[:descWidth-3] + "..."
		}

		dueDate := ""
		if t.DueDate != nil {
			dueDate = t.DueDate.Format("2006-01-02")
		}

		tags := ""
		if len(t.Tags) > 0 {
			tags = "#" + strings.Join(t.Tags, " #")
			if len(tags) > tagsWidth {
				tags = tags[:tagsWidth-3] + "..."
			}
		}

		row := f.formatRow(
			fmt.Sprintf("%d", t.ID),
			status,
			t.Priority.String(),
			desc,
			dueDate,
			tags,
			idWidth, statusWidth, priorityWidth, descWidth, dueDateWidth, tagsWidth,
		)

		if f.colorOutput {
			if t.Completed {
				row = Dim(row)
			}
			if t.IsOverdue() {
				row = Red(row)
			}
		}

		lines = append(lines, row)
	}

	return strings.Join(lines, "\n")
}

func (f *TableFormatter) formatRow(id, status, priority, desc, dueDate, tags string, idW, statusW, priorityW, descW, dueDateW, tagsW int) string {
	return fmt.Sprintf("%-*s %-*s %-*s %-*s %-*s %-*s",
		idW, id,
		statusW, status,
		priorityW, priority,
		descW, desc,
		dueDateW, dueDate,
		tagsW, tags,
	)
}
