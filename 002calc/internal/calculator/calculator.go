package calculator

import (
	"fmt"

	calcErrors "github.com/samnart1/GoLang-Projects/tree/main/002calc/pkg/errors"
)

type Calculator struct {
	operations map[string]Operation
	history []CalculationResult
}

type CalculationResult struct {
	Expression *Expression
	Result float64
	Error error
}

func New() *Calculator {
	return &Calculator{
		operations: GetSupportedOperations(),
		history: make([]CalculationResult, 0),
	}
}

func (c *Calculator) Calculate(expression *Expression) (float64, error) {
	operation, exists := c.operations[expression.Operator]
	if !exists {
		err := calcErrors.NewCalculatorError("calculate", expression.Operator, calcErrors.ErrInvalidOperation)
		c.addToHistory(expression, 0, err)
		return 0, err
	}

	result, err := operation.Function(expression.FirstNumber, expression.SecondNumber)
	c.addToHistory(expression, result, err)

	return result, err
}

func (c *Calculator) CalculateFromString(input string) (float64, error) {
	expression, err := ParseExpression(input)
	if err != nil {
		return 0, err
	}

	return c.Calculate(expression)
}

func (c *Calculator) GetHistory() []CalculationResult {
	return c.history
}

func (c *Calculator) ClearHistory() {
	c.history = make([]CalculationResult, 0)
}

func (c *Calculator) GetSupportedOperationsInfo() map[string]Operation {
	return c.operations
}

func (c *Calculator) addToHistory(expression *Expression, result float64, err error) {
	c.history = append(c.history, CalculationResult{
		Expression: expression,
		Result: result,
		Error: err,
	})
} 

func (c *Calculator) FormatResult(result float64, expression *Expression) string {
	if result == float64(int64(result)) {
		return fmt.Sprintf("%.0f", result)
	}

	return fmt.Sprintf("%.6g", result)
}