package errors

import "fmt"

type TaskError struct {
	Op		string
	Err		error
}

type StorageError struct {
	Op		string
	Path	string
	Err		error
}

type ValidationError struct {
	Field	string
	Message	string
}

func (e *TaskError) Error() string {
	return fmt.Sprintf("Task %s: %v\n", e.Op, e.Err)
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("Storage %s with path %s: %v\n", e.Op, e.Path, e.Err)
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error: %s - %s", e.Field, e.Message)
}

func NewTaskError(op string, err error) *TaskError {
	return &TaskError{
		Op: op,
		Err: err,
	}
}

func NewStorageError(op, path string, err error) *StorageError {
	return &StorageError{
		Op: op,
		Path: path,
		Err: err,
	}
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field: field,
		Message: message,
	}
}
