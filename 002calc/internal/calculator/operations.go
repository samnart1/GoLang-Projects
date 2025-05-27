package calculator

import (
	"math"

	calcErrors "github.com/samnart1/GoLang-Projects/tree/main/002calc/pkg/errors"
)

type Operation struct {
	Symbol 		string
	Name 		string
	Description	string
	Function	func(a, b float64) (float64, error)
}

func GetSupportedOperations() map[string]Operation {
	return map[string]Operation {
		"+": {
			Symbol: 		"+",
			Name: 			"Addition",
			Description: 	"Addtion of two numbers",
			Function: 		add,
		},
		"-": {
			Symbol: 		"-",
			Name: 			"Subtration",
			Description: 	"Subrate second number from first",
			Function: 		subtract,
		},
		"/": {
			Symbol: 		"/",
			Name: 			"Division",
			Description: 	"Divide first number by second number",
			Function: 		divide,
		},
		"*": {
			Symbol: 		"*",
			Name: 			"Multiplication",
			Description: 	"Multiply two numbers together",
			Function: 		multiply,
		},
		"^": {
			Symbol: 		"^",
			Name: 			"Power",
			Description: 	"Raise first number to the power of the second number",
			Function: 		power,
		},
		"%": {
			Symbol: 		"%",
			Name: 			"Modulos",
			Description: 	"Return remainder after first number is divided by second number",
			Function: 		mod,
		},
	}
}

func add(a, b float64) (float64, error) {
	return a + b, nil
}

func subtract(a, b float64) (float64, error) {
	return a - b, nil
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, calcErrors.ErrDivisionByZero
	}
	return a / b, nil
}

func multiply(a, b float64) (float64, error) {
	return a * b, nil
}

func power(a, b float64) (float64, error) {
	result := math.Pow(a, b)
	if math.IsInf(result, 0) || math.IsNaN(result) {
		return 0, calcErrors.NewCalculatorError("power", "", calcErrors.ErrInvalidOperation)
	}
	return result, nil
}

func mod(a, b float64) (float64, error) {
	if b == 0 {
		return 0, calcErrors.ErrDivisionByZero
	}
	return math.Mod(a, b), nil
}