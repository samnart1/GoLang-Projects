package cmd

import (
	"github.com/go-delve/delve/pkg/dwarf/reader"
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
		filePath:	filePath,
		maxLines:	maxLines
	})

}