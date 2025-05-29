package task

import (
	"sort"

	"github.com/samnart1/GoLang-Projects/003todo/pkg/errors"
)

type Manager struct {
	tasks 	[]*Task
	nextID 	int
}

func NewManager() *Manager {
	return &Manager{
		tasks: 	make([]*Task, 0),
		nextID: 1,
	}
}

func (m *Manager) LoadTasks(tasks []*Task) {
	m.tasks = tasks
	m.updateNextID()
}

func (m *Manager) GetTasks() []*Task {
	return m.tasks
}

func (m *Manager) AddTask(description string) (*Task, error) {
	task := NewTask(m.nextID, description)
	if err := task.Validate(); err != nil {
		return nil, err
	}

	m.tasks = append(m.tasks, task)
	m.nextID++

	return task, nil
} 

func (m *Manager) GetTaskByID(id int) (*Task, error) {
	for _, task := range m.tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return nil, errors.NewTaskError("get", errors.NewValidationError("id", "task not found!"))
}

func (m *Manager) RemoveTask(id int) error {
	for i, task := range m.tasks {
		if task.ID == id {
			m.tasks = append(m.tasks[:i], m.tasks[i + 1:]...)
			return nil
		}
	}
	return errors.NewTaskError("remove", errors.NewValidationError("id", "task not found"))
}

func (m *Manager) CompleteTask(id int) error {
	task, err := m.GetTaskByID(id)
	if err != nil {
		return err
	}

	task.Complete()
	return nil
}

func (m *Manager) EditTask(id int, description string) error {
	task, err := m.GetTaskByID(id)
	if err != nil {
		return err
	}

	task.Description = description
	return task.Validate()
}

func (m *Manager) GetStats() map[string]int {
	stats := map[string]int {
		"total": len(m.tasks),
		"completed": 0,
		"pending": 0,
		"overdue": 0,
	}

	for _, task := range m.tasks {
		if task.Completed {
			stats["completed"]++
		} else {
			stats["pending"]++
			if task.IsOverdue() {
				stats["overdue"]++
			}
		}
	}
	return stats
}

func (m *Manager) SortByPriority() {
	sort.Slice(m.tasks, func(i, j int) bool {
		return m.tasks[i].Priority > m.tasks[j].Priority
	})
}

func (m *Manager) SortByCreated() {
	sort.Slice(m.tasks, func(i, j int) bool {
		return m.tasks[i].CreatedAt.Before(*m.tasks[j].CompletedAt)
	})
}

func (m *Manager) SortByDueDate() {
	sort.Slice(m.tasks, func(i, j int) bool {
		if m.tasks[i].DueDate == nil {
			return false
		}
		if m.tasks[j].DueDate == nil {
			return true
		}
		return m.tasks[i].DueDate.Before(*m.tasks[j].DueDate)
	})
}

func (m *Manager) updateNextID() {
	maxID := 0
	for _, task := range m.tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	m.nextID = maxID + 1
	
}