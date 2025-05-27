package errors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidExpression	= errors.New("invalid mathematical expresson")
	ErrDivisionByZero		= errors.New("division by zero is not allowed")
	ErrInvalidOperation		= errors.New("invalid operaton")
	ErrInvalidNumber		= errors.New("invalid number format")
)

type CalculatorError struct {
	Operation	string
	Input		string
	Err			error
}

func (e *CalculatorError) Error() string {
	return fmt.Sprintf("calculator error in %s with input '%s': %v", e.Operation, e.Input, e.Err)
}

func NewCalculatorError(operation, input string, err error) *CalculatorError {
	return &CalculatorError{
		Operation: 	operation,
		Input: 		input,
		Err: 		err,
	}
}