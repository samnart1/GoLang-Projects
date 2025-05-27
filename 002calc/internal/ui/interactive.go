package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/samnart1/GoLang-Projects/tree/main/002calc/internal/calculator"
)

func InteractiveMode(verbose bool) error {
	calc := calculator.New()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Interactive Calculator")
	fmt.Println("Enter mathematical expressions (e.g., 5 + 3, 10 * 5)")
	fmt.Println("Specail commands:")
	fmt.Println("	help		- Show supported operations")
	fmt.Println("	history		- Show calculation history")
	fmt.Println("	clear		- Clear history")
	fmt.Println("	exit		- Exit calculator")
	fmt.Println()

	for {
		fmt.Print("calc> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		switch input {
			case "exit", "quit":
				fmt.Println("Goodbye!")
				return nil
			case "help":
				showHelp(calc)
				continue
			case "history":
				showHistory(calc, verbose)
				continue
			case "clear":
				calc.ClearHistory()
				fmt.Println("History cleared.")
				continue
			case "":
				continue
		}

		result, err := calc.CalculateFromString(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		expr, _ := calculator.ParseExpression(input)
		formatedResult := calc.FormatResult(result, expr)

		fmt.Printf("=%s\n", formatedResult)

		if verbose {
			fmt.Printf("	(Input: %s, Result: %f)\n", input, result)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil
}

func showHelp(calc *calculator.Calculator) {
	fmt.Println("\nSupported Operations:")
	operations := calc.GetSupportedOperationsInfo()

	for symbol, op := range operations {
		fmt.Printf("	%s	%s - %s\n", symbol, op.Name, op.Description)
	}

	fmt.Println("\nExamples")
	fmt.Println("	5 + 3		-> Addition")
	fmt.Println("	10 - 4		-> Subtraction")
	fmt.Println("	7 * 6		-> Multiplication")
	fmt.Println("	15 / 3		-> Division")
	fmt.Println("	2 ^ 8		-> Power (2 to the power of 8)")
	fmt.Println("	17 % 5		-> Modulus (remainder)")
	fmt.Println(" 	-5 + 10		-> Works wih negative numbers")
	fmt.Println("	3.14 * 2	-> Works with decimals")
	fmt.Println()
}

func showHistory(calc *calculator.Calculator, verbose bool) {}