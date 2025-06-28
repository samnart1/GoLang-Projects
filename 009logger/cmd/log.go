package cmd

import (
	"fmt"

	"github.com/samnart1/golang/009l9gger/internal/config"
	"github.com/samnart1/golang/009l9gger/internal/logger"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use: "log [message]",
	Short: "Log a message with timestamp",
	Long: `Log a message to the specified log file or console with timestamp.
	
		Example:
			go-simple-logger log "Hello, World!"
			go-simple-logger log --level error "Something went wrong"
			go-simple-logger log --file /var/log/app.log "Application started"`,

	Args: cobra.MinimumNArgs(1),
	RunE: runLog,
}

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().String("level", "info", "log level (debug, info, warn, error, fatal)")
	logCmd.Flags().String("source", "", "log source/component name")
	logCmd.Flags().Bool("console", false, "also output to console")
}

func runLog(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	level, _ := cmd.Flags().GetString("level")
	source, _ := cmd.Flags().GetString("source")
	console, _ := cmd.Flags().GetBool("console")

	if level != "info" {
		cfg.Log.Level = level
	}
	if console {
		cfg.Log.Output = "both"
	}

	l, err := logger.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer l.Close()

	message := ""
	for i, arg := range args {
		if i > 0 {
			message += " "
		}
		message += arg
	}

	entry := logger.Entry{
		Message: message,
		Level: level,
		Source: source,
	}

	if err := l.Log(entry); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	fmt.Printf("Message logged successfully\n")
	return nil
}