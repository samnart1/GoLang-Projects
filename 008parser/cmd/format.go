package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/samnart1/golang/008parser/internal/parser"
	"github.com/samnart1/golang/008parser/pkg/types"
	"github.com/spf13/cobra"
)

var formatCmd = &cobra.Command{
	Use: "format [file...]",
	Short: "Format and pretty-print JSON files",
	Long: `Format JSON files with customizable indentation and output options.
		
		This command read  JSON files, validates them, and outputs formatted versions`,
		
	Args: cobra.MinimumNArgs(1),
	RunE: runFormat,
}

var (
	formatIndent int
	useTabs 	bool
	sortKeys 	bool
	inPlace 	bool
	backup 		bool
)

func init() {
	rootCmd.AddCommand(formatCmd)

	formatCmd.Flags().IntVar(&formatIndent, "indent", 2, "Number of spaces for indentation")
	formatCmd.Flags().BoolVar(&useTabs, "use-tabs", false, "Use tabs instead of spaces")
	formatCmd.Flags().BoolVar(&sortKeys, "sort-keys", false, "Sort object keys alphabetically")
	formatCmd.Flags().BoolVar(&inPlace, "in-place", false, "Modify files in place")
	formatCmd.Flags().BoolVar(&backup, "backup", false, "Create backup files when using --in-place")
	formatCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (single file only)")
}

func runFormat(cmd *cobra.Command, args []string) error {
	if len(args) > 1 && outputFile != "" {
		return fmt.Errorf("cannot use --output with multiple input files")
	}

	if outputFile != "" && inPlace {
		return fmt.Errorf("cannot use --output with --in-place")
	}

	p := parser.New(cfg)
	formatter := parser.NewFormatter(cfg)

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
			if err := formatFile(p, formatter, match); err != nil {
				if !cfg.Quiet {
					fmt.Fprintf(os.Stderr, "Error formatting '%s': %v\n", match, err)
				}
				hasErrors = true
			}
		}
	}

	if hasErrors {
		os.Exit(1)
	}

	return nil
}

func formatFile(p *parser.Parser, formatter *parser.Formatter, filePath string) error {
	result, err := p.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse: %w", err)
	}

	options := types.OutputOptions{
		Format: "pretty",
		Indent: formatIndent,
		UseTabs: useTabs,
		SortKeys: sortKeys,
		Colors: false,
	}

	formatted, err := formatter.Format(result.Data, options)
	if err != nil {
		return fmt.Errorf("failed to format: %w", err)
	}

	var outputPath string
	if inPlace {
		outputPath = filePath

		if backup {
			backupPath := filePath + ".bak"
			if err := copyFile(filePath, backupPath); err != nil {
				return fmt.Errorf("failed to create backup: %w", err)
			}
		}
	} else if outputFile != "" {
		outputPath = outputFile
	}

	if outputPath != "" {
		err := os.WriteFile(outputPath, []byte(formatted), 0644)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		if !cfg.Quiet {
			if inPlace {
				fmt.Printf("Formatted: %s\n", filePath)
			} else {
				fmt.Printf("Formatted %s -> %s\n", filePath, outputPath)
			}
		}
	} else {
		fmt.Print(formatted)
	}

	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}