package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/samnart1/golang/008parser/internal/parser"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use: "validate [file...]",
	Short: "Validate JSON files",
	Long: `Validate one or more JSON files for syntax errors and structural issues
		Returns exit code 0 if all files are valid, 1 if any file is invalid`,

	Args: cobra.MinimumNArgs(1),
	RunE: runValidate,
}

var (
	strict 		bool
	showValid	bool
)

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().BoolVar(&strict, "strict", false, "Enable strict validation (no duplicate keys, etc.)")
	validateCmd.Flags().BoolVar(&showValid, "show-valid", false, "Show valid files in output")
}

func runValidate(cmd *cobra.Command, args []string) error {
	validator := parser.NewValidator(cfg)

	var totalFiles, validFiles, invalidFiles int
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
			totalFiles++

			result := validator.ValidateFile(match, strict)

			if result.IsValid {
				validFiles++
				if showValid && !cfg.Quiet {
					if cfg.EnableColors && !cfg.NoColor {
						fmt.Printf("\033[32m✓\033[0m %s\n", match)
					} else {
						fmt.Printf("✓ %s\n", match)
					}
				}
			} else {
				invalidFiles++
				hasErrors = true

				if !cfg.Quiet {
					if cfg.EnableColors && !cfg.NoColor {
						fmt.Printf("\033[31mx\033[0m %s\n", match)
					} else {
						fmt.Printf("x %s\n", match)
					}

					for _, err := range result.Errors {
						fmt.Printf("	Error: %s\n", err)
					}

					if len(result.Warnings) > 0 {
						for _, warning := range result.Warnings {
							fmt.Printf("	Warning: %s\n", warning)
						}
					}
				}
			}
		}
	}

	if !cfg.Quiet {
		fmt.Printf("\nValidation Summary:\n")
		fmt.Printf("	Total files: %d\n", totalFiles)

		if cfg.EnableColors && !cfg.NoColor {
			fmt.Printf("	Valid: \033[32m%d\033[0m\n", invalidFiles)
		} else {
			fmt.Printf("	Valid: %d\n", validFiles)
			fmt.Printf("	Invalid: %d\n", invalidFiles)
		}
	}

	if hasErrors {
		os.Exit(1)
	}

	return nil
}