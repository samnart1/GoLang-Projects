package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/samnart1/golang/008parser/internal/parser"
	"github.com/samnart1/golang/008parser/pkg/types"
	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use: "parse [file...]",
	Short: "Parse JSON files and display their contents",
	Long: `Parse one or more JSON files an display their contents in various formats
		
		Supported output formats:
			- compact: Single line, no extra whitespace
			- pretty: Formatted with indentation
			- tree: Tree-like structure showing hierarchy
			- table: Tabular format for arrays of objects`,

	Args: cobra.MinimumNArgs(1),
	RunE: runParse,
}

var (
	outputFormat	string
	outputFile		string
	indent			int
	showKeys		bool
	showTypes		bool
	maxDepth		int
)

func init() {
	rootCmd.AddCommand(parseCmd)

	parseCmd.Flags().StringVarP(&outputFormat, "format", "f", "pretty", "Output format (compact, pretty, tree, table)")
	parseCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	parseCmd.Flags().IntVar(&indent, "indent", 2, "Indentation spaces for pretty format")
	parseCmd.Flags().BoolVar(&showKeys, "show-keys", false, "Show object keys in tree format")
	parseCmd.Flags().BoolVar(&showTypes, "show-types", false, "Show value types")
	parseCmd.Flags().IntVar(&maxDepth, "max-depth", 0, "Maximum depth to display (0 = unlimited)")
}

func runParse(cmd *cobra.Command, args []string) error {
	p := parser.New(cfg)

	var allResults []types.ParseResult
	var hasErrors bool

	for _, filePath := range args {
		matches, err := filepath.Glob(filePath)
		if err != nil {
			return fmt.Errorf("invalid file pattern '%s': %w", filePath, err)
		}

		if len(matches) == 0 {
			if !cfg.Quiet {
				fmt.Fprintf(os.Stderr, "Warning: no files match pattern '%s'\n", filePath)
			}
			continue
		}

		for _, match := range matches {
			result, err := p.ParseFile(match)
			if err != nil {
				if !cfg.Quiet {
					fmt.Fprintf(os.Stderr, "Error parsing '%s': %v\n", match, err)
				}
				hasErrors = true
				continue
			}

			allResults = append(allResults, *result)
		}
	}

	if len(allResults) == 0 {
		return fmt.Errorf("no files were successfully parsed")
	}

	options := types.OutputOptions{
		Format: outputFormat,
		Indent: indent,
		ShowKeys: showKeys,
		ShowTypes: showTypes,
		MaxDepth: maxDepth,
		Colors: cfg.EnableColors && !cfg.NoColor,
	}

	formatter := parser.NewFormatter(cfg)

	var output strings.Builder
	for i, result := range allResults {
		if len(allResults) > 1 {
			if i > 0 {
				output.WriteString("\n")
			}
			output.WriteString(fmt.Sprintf("=== %s ===\n", result.FilePath))
		}

		formatted, err := formatter.Format(result.Data, options)
		if err != nil {
			return fmt.Errorf("error formatting output: %w", err)
		}

		output.WriteString(formatted)
		if i < len(allResults)-1 {
			output.WriteString("\n")
		}
	}

	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(output.String()), 0644)
		if err != nil {
			return fmt.Errorf("error writing to file '%s': %w", outputFile, err)
		}
		if !cfg.Quiet {
			fmt.Printf("Ooutput written to: %s\n", outputFile)
		}
	} else {
		fmt.Print(output.String())
	}

	if hasErrors {
		os.Exit(1)
	}

	return nil
}