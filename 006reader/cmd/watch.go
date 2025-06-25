package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/samnart1/GoLang/006reader/internal/logger"
	"github.com/samnart1/GoLang/006reader/internal/reader"
	"github.com/spf13/cobra"
)

var (
	tailMode bool
	interval int
)

var watchCmd = &cobra.Command{
	Use: "watch [file]",
	Short: "Watch file for changes and display updates",
	Long: `Watch a file for changes and display updates in real-time
	
	Examples:
		go-file-reader watch file.txt
		go-file-reader watch --tail file.log
		go-file-reader watch --interval 500 file.txt`,
	Args: cobra.ExactArgs(1),
	RunE: runWatch,
}

func init() {
	rootCmd.AddCommand(watchCmd)

	watchCmd.Flags().BoolVarP(&tailMode, "tail", "t", false, "tail mode (show only new content)")
	watchCmd.Flags().IntVarP(&interval, "interval", "i", 1000, "polling interval in milliseconds")
}

func runWatch(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	if err := validateFile(filePath); err != nil {
		return err
	}

	watcher := reader.NewWatcher(&reader.WatchConfig{
		FilePath:		filePath,
		TailMode:		tailMode,
		PollInterval:	interval,
		BufferSize:		cfg.Reader.BufferSize,
	})

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Watching file: %s (Press Ctrl+C to stop)\n", filePath)
	fmt.Println("------------------------------------------")

	if err := watcher.Start(); err != nil {
		return fmt.Errorf("failed to start watcher: %w", err)
	} 
	defer watcher.Stop()

	for {
		select {
		case event := <-watcher.Events():
			fmt.Printf("[%s] %s\n", event.Time.Format("15:04:05"), event.Content)

		case err := <-watcher.Errors():
			log.Error("Watcher error", logger.Error(err))

		case <-signChan:
			fmt.Println("\nStopping watcher...")
			log.Info("Watcher stopped by user")
			return nil
		}
	}
}