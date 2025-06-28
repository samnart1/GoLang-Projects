package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var tailCmd = &cobra.Command{
	Use: "tail [log-file]",
	Short: "Monitor log file in real-time",
	Long: `Monitor a log file in real-time, similar to 'tail -f'
	
		Example:
			go-simple-logger tail /var/log/app.log
			go-simple-logger tail --lines 50 /var/log/app.log
			go-simple-logger tail --level error /var/log/app.log`,

	Args: cobra.ExactArgs(1),
	RunE: runTail,
}

func init() {
	rootCmd.AddCommand(tailCmd)

	tailCmd.Flags().Int("lines", 10, "number of lines to show initially")
	tailCmd.Flags().String("level", "", "filter by log level")
	tailCmd.Flags().Bool("follow", true, "follow the log file for new entries")
}

func runTail(cmd *cobra.Command, args []string) error {
	filename := args[0]
	lines, _ := cmd.Flags().GetInt("lines")
	levelFilter, _ := cmd.Flags().GetString("level")
	follow, _ := cmd.Flags().GetBool("follow")

	if _, err := os.Stat(filename); err != nil {
		return fmt.Errorf("log file does not exist: %s", filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	if err := showLastLines(file, lines, levelFilter); err != nil {
		return fmt.Errorf("failed to read log file: %w", err)
	}

	if !follow {
		return nil
	}

	fmt.Println("--- following log file (Ctrl+C to stop) ---")
	return followFile(filename, levelFilter)
}

func showLastLines(file *os.File, n int, levelFilter string) error {
	scanner := bufio.NewScanner(file)
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		if levelFilter == "" || strings.Contains(strings.ToLower(line), strings.ToLower(levelFilter)) {
			lines = append(lines, line)
		}
	}

	start := len(lines) - n
	if start < 0 {
		start = 0
	}

	for i := start; i < len(lines); i++ {
		fmt.Println(lines[i])
	}

	return scanner.Err()
}

func followFile(filename, levelFilter string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Seek(0, 2)
	scanner := bufio.NewScanner(file)

	for {
		if scanner.Scan() {
			line := scanner.Text()
			if levelFilter == "" || strings.Contains(strings.ToLower(line), strings.ToLower(levelFilter)) {
				fmt.Println(line)
			}
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}