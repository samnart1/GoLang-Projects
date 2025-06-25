package cmd

import (
	"fmt"
	"os"

	"github.com/samnart1/GoLang/006reader/internal/formatter"
	"github.com/samnart1/GoLang/006reader/internal/logger"
	"github.com/samnart1/GoLang/006reader/internal/reader"
	"github.com/spf13/cobra"
)

var (
	format		string
	lineNumber	bool
	maxLines	int
	encoding	string
)

var readCmd = &cobra.Command{
	Use: "read [file]",
	Short: "Read and display file contents",
	Long: `Read and display the contents of a text file with version formatiing options
	
	Examples:
		go-file-reader read file.txt
		go-file-reader read --format json file.txt
		go-file-reader read --lines --max-lines 50 file.txt`,
	Args: cobra.ExactArgs(1),
	RunE: runRead,
}

func init() {
	rootCmd.AddCommand(readCmd)

	readCmd.Flags().StringVarP(&format, "format", "f", "plain", "output format (plain, json, table)")
	readCmd.Flags().BoolVarP(&lineNumber, "lines", "l", false, "show line numbers")
	readCmd.Flags().IntVarP(&maxLines, "max-lines", "n", 0, "maximum number of lines to read (0 for all)")
	readCmd.Flags().StringVarP(&encoding, "encoding", "e", "utf-8", "file encoding")
}

func runRead(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	if err := validateFile(filePath); err != nil {
		return err
	}

	fileReader := reader.New(&reader.Config{
		FilePath:	filePath,
		MaxLines:	maxLines,
		ShowLines:	lineNumber,
		Encoding:	encoding,
		BufferSize:	cfg.Reader.BufferSize,
	})

	content, err := fileReader.Read()
	if err != nil {
		log.Error("Failed to read file",
			logger.String("file", filePath),
			logger.Error(err))
		return fmt.Errorf("failed to read file: %w", err)
	}

	fmt_handler, err := formatter.New(format, &formatter.Config{
		ShowLineNumbers:	lineNumber,
		MaxWidth:			cfg.Formatter.MaxWidth,
		Theme:				cfg.Formatter.Theme,
	})
	if err != nil {
		return fmt.Errorf("failed to create formatter: %w", err)
	}

	output, err := fmt_handler.Format(content)
	if err != nil {
		return fmt.Errorf("failed to format content: %w", err)
	}

	fmt.Println(output)

	log.Info("File read successfully",
 		logger.String("file", filePath),
		logger.Int("lines", len(content.Lines)),
		logger.String("format", format))

	return nil
}

func validateFile(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filePath)
		}
		return fmt.Errorf("cannot access file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", filePath)
	}
	if info.Size() > int64(cfg.Reader.MaxFileSize) {
		return fmt.Errorf("file too large: %d bytes (max: %d)", info.Size(), cfg.Reader.MaxFileSize)
	}

	return nil
}