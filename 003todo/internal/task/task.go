package task

import (
	"strings"
	"time"

	"github.com/samnart1/GoLang-Projects/003todo/pkg/errors"
)

type Priority int

const (
	Low Priority = iota
	Medium
	High
)

func (p Priority) String() string {
	switch p {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	default:
		return "Unknown"
	}
}

type Task struct {
	ID 			int			`json:"id"`
	Description string		`json:"description"`
	Completed 	bool		`json:"completed"`
	CreatedAt 	time.Time	`json:"created_at"`
	CompletedAt	*time.Time	`json:"completed_at,omitempty"`
	DueDate 	*time.Time	`json:"due_date,omitempty"`
	Priority 	Priority	`json:"priority"`
	Tags 		[]string	`json:"tags,omitempty"`
}

func NewTask(id int, description string) *Task {
	return &Task{
		ID: id,
		Description: description,
		Completed: false,
		CreatedAt: time.Now(),
		Priority: Medium,
		Tags: []string{},
	}
}

func (t *Task) Complete() {
	t.Completed = true
	now := time.Now()
	t.CompletedAt = &now
} 

func (t *Task) Uncomplete() {
	t.Completed = false
	t.CompletedAt = nil
}

func (t *Task) SetPriority(priority Priority) {
	t.Priority = priority
}

func (t *Task) AddTag(tag string) {
	tag = strings.TrimSpace(strings.ToLower(tag))
	if tag == "" {
		return
	}

	for _, existingTag := range t.Tags {
		if existingTag == tag {
			return
		}
	}

	t.Tags = append(t.Tags, tag)
}

func (t *Task) RemoveTag(tag string) {
	tag = strings.TrimSpace(strings.ToLower(tag))
	for i, existingTag := range t.Tags {
		if existingTag == tag {
			t.Tags = append(t.Tags[:i], t.Tags[i + 1:]...)
			return
		}
	}
}

func (t *Task) HashTag(tag string) bool {
	tag = strings.TrimSpace(strings.ToLower(tag))
	for _, existingTag := range t.Tags {
		if existingTag == tag {
			return true
		}
	}
	return false
}

func (t *Task) IsOverdue() bool {
	if t.DueDate == nil || t.Completed {
		return false
	}
	return t.DueDate.Before(time.Now())
}

func (t *Task) Validate() error {
	if strings.TrimSpace(t.Description) == "" {
		return errors.NewValidationError("description", "cannot be empty")
	}

	if t.ID <= 0 {
		return errors.NewValidationError("id", "must be positive")
	}

	return nil
}