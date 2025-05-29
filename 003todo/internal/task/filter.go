package task

import (
	"strings"
	"time"
)

type Filter struct {
	ShowCompleted	bool
	ShowPending		bool
	Priority		*Priority
	Tag				string
	DueBefore			*time.Time
	DueAfter		*time.Time
	SearchTerm		string
}

func NewFilter() *Filter {
	return &Filter{
		ShowCompleted: true,
		ShowPending: true,
	}
}

func (f *Filter) Apply(tasks []*Task) []*Task {
	var filtered []*Task

	for _, task := range tasks {
		if f.matches(task) {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

func (f *Filter) matches(task *Task) bool {
	if task.Completed && !f.ShowCompleted {
		return false
	}
	if !task.Completed && !f.ShowPending {
		return false
	}
	if f.Priority != nil && task.Priority != *f.Priority {
		return false
	}
	if f.Tag != "" && !task.HashTag(f.Tag) {
		return false
	}
	if f.DueBefore != nil && task.DueDate != nil && !task.DueDate.Before(*f.DueBefore) {
		return false
	}
	if f.DueAfter != nil && task.DueDate != nil && !task.DueDate.After(*f.DueAfter) {
		return false
	}

	if f.SearchTerm != "" {
		searchLower := strings.ToLower(f.SearchTerm)
		descLower := strings.ToLower(task.Description)
		if !strings.Contains(descLower, searchLower) {
			return false
		}
	}
	
	return true
}

func (f *Filter) SetPendingOnly() {
	f.ShowCompleted = false
	f.ShowPending = true
}

func (f *Filter) SetCompletedOnly() {
	f.ShowCompleted = true
	f.ShowPending = false
}