package cmd

import (
	"github.com/samnart1/golang/010timer/internal/timer"
	"github.com/samnart1/golang/010timer/pkg/types"
	"github.com/spf13/cobra"
)

var (
	showLaps bool
)

var stopwatchCmd = &cobra.Command{
	Use: "stopwatch",
	Short: "Start a stopwatch",
	Long: `Start a stopwatch to measure elapsed time.
	
		Controls:
			- SPACE: Start/Stop the stopwatch
			- L: Record a lap time (when --laps is enabled)
			- R: Reset the stopwatch
			- Q: Quit
			
		Example:
			timer stopwatch
			timer stopwatch --laps`,

	RunE: runStopwatch,
}

func init() {
	rootCmd.AddCommand(stopwatchCmd)

	stopwatchCmd.Flags().BoolVar(&showLaps, "laps", false, "enable lap time recording")
}

func runStopwatch(cmd *cobra.Command, args []string) error {
	config := GetConfig()

	stopwatchConfig := &types.StopWatchConfig{
		ShowLaps: showLaps,
		Verbose: verbose,
		Color: config.Display.Color,
	}
	sw := timer.NewStopwatch(stopwatchConfig)
	return sw.Start()
}