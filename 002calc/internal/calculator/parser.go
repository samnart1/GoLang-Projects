package calculator

import (
	"regexp"
	"strconv"
	"strings"

	calcErrors "github.com/samnart1/GoLang-Projects/002calc/pkg/errors"
)

type Expression struct {
	FirstNumber float64
	Operator	string
	SecondNumber	float64
	Raw string
}

func ParseExpression(input string) (*Expression, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, calcErrors.NewCalculatorError("parse", input, calcErrors.ErrInvalidExpression)
	}

	re := regexp.MustCompile(`^(-?\d+\.?\d*)\s*([+\-*/^%])\s*(-?\d+\.?\d*)$`)
	matches := re.FindStringSubmatch(input)

	if len(matches) != 4 {
		return nil, calcErrors.NewCalculatorError("parse", input, calcErrors.ErrInvalidExpression)
	}

	firstNumber, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return nil, calcErrors.NewCalculatorError("parse", matches[1], calcErrors.ErrInvalidNumber)
	}

	operator := matches[2]

	secondNumber, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return nil, calcErrors.NewCalculatorError("parse", matches[3], calcErrors.ErrInvalidNumber)
	}

	return &Expression{
		FirstNumber: firstNumber,
		Operator: operator,
		SecondNumber: secondNumber,
		Raw: input,
	}, nil
}