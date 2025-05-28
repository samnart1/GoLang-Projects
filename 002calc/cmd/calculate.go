package cmd

import (
	"fmt"
	"strings"

	"github.com/samnart1/GoLang-Projects/002calc/internal/calculator"
	"github.com/spf13/cobra"
)

var calculateCmd = &cobra.Command{
	Use: "calc [expression]",
	Short: "Calculate a mathematical expression",
	Long: `Calculate a mathematical expression directly. Example: calc "5 + 3"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		expression := strings.Join(args, " ")

		calc := calculator.New()
		result, err := calc.CalculateFromString(expression)
		if err != nil {
			return fmt.Errorf("calculation failed: %w", err)
		}

		expr, _ := calculator.ParseExpression(expression)
		formattedResult := calc.FormatResult(result, expr)

		if verbose {
			fmt.Printf("Expression: %s\n", expression)
			fmt.Printf("Result: %s\n", formattedResult)
			fmt.Printf("Raw result: %f\n", result)
		} else {
			fmt.Printf("%s = %s\n", expression, formattedResult)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(calculateCmd)
}