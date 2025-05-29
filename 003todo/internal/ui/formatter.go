package ui

import (
	"fmt"
	"strings"

	"github.com/samnart1/GoLang-Projects/003todo/internal/task"
)

type TaskFormatter struct {
	showID			bool
	showStatus		bool
	showPriority	bool
	showDueDate		bool
	showTags		bool
	colorOutput		bool
}

func NewTaskFormatter() *TaskFormatter {
	return &TaskFormatter{
		showID: true,
		showStatus: true,
		showPriority: true,
		showDueDate: true,
		showTags: true,
		colorOutput: true,
	}
}

func (f *TaskFormatter) FormatTask(t *task.Task) string {
	var parts []string

	if f.showID {
		if f.colorOutput {
			parts = append(parts, Dim(fmt.Sprintf("[%d]", t.ID)))
		} else {
			parts = append(parts, fmt.Sprintf("[%d]", t.ID))
		}
	}

	if f.showStatus {
		status := "○"
		if t.Completed {
			status = "●"
		}
		if f.colorOutput {
			parts = append(parts, StatusColor(t.Completed)(status))
		} else {
			parts = append(parts, status)
		}
	}

	if f.showPriority {
		priority := fmt.Sprintf("(%s)", t.Priority.String())
		if f.colorOutput {
			parts = append(parts, PriorityColor(t.Priority.String())(priority))
		} else {
			parts = append(parts, priority)
		}
	}

	desc := t.Description
	if f.colorOutput && t.Completed {
		desc = Dim(desc)
	}
	parts = append(parts, desc)

	if f.showDueDate && t.DueDate != nil {
		dueStr := fmt.Sprintf("(due: %s)", t.DueDate.Format("2006-01-02"))
		if f.colorOutput {
			if t.IsOverdue() {
				dueStr = Red(dueStr)
			} else {
				dueStr = Cyan(dueStr)
			}
		}
		parts = append(parts, dueStr)
	}

	if f.showTags && len(t.Tags) > 0 {
		tagStr := "#" + strings.Join(t.Tags, " #")
		if f.colorOutput {
			tagStr = Magenta(tagStr)
		}
		parts = append(parts, tagStr)
	}

	return strings.Join(parts, " ")
}

func (f *TaskFormatter) FormatTaskList(tasks []*task.Task) string {
	if len(tasks) == 0 {
		return "No tasks found"
	}

	var lines []string
	for _, task := range tasks {
		lines = append(lines, f.FormatTask(task))
	}

	return strings.Join(lines, "\n")
}

func (f *TaskFormatter) SetOptions(opts map[string]bool) {
	if showID, ok := opts["showID"]; ok {
		f.showID = showID
	}
	if showStatus, ok := opts["showStatus"]; ok {
		f.showStatus = showStatus
	}
	if showPriority, ok := opts["showPriority"]; ok {
		f.showPriority = showPriority
	}
	if showDueDate, ok := opts["showDueDate"]; ok {
		f.showDueDate = showDueDate
	}
	if showTags, ok := opts["showTags"]; ok {
		f.showTags = showTags
	}
	if colorOutput, ok := opts["colorOutput"]; ok {
		f.colorOutput = colorOutput
	}
}