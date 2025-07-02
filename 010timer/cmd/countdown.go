package cmd

import (
	"fmt"
	"time"

	"github.com/samnart1/golang/010timer/internal/timer"
	"github.com/samnart1/golang/010timer/pkg/types"
	"github.com/spf13/cobra"
)

var (
	duration string
	message string
	sound bool
	desktop bool
)

var countdownCmd = &cobra.Command{
	Use: "countdown",
	Short: "Start a countdown timer",
	Long: `Start a countdown timer with the specified duration
	
		Duration can be specified in various formats:
			- 30s (30 seconds)
			- 5m (5 minutes)
			- 1h (1 hour)
			- 1h30m (1 hour 30 minutes)
			- 90s (90 seconds)
			
		Examples: 
			timer countdown -d 5m
			timer countdown -d 30s -m "Time's up!"
			timer countdown -d 1h30s --sound`,

	RunE: runCountdown,
}

func init() {
	rootCmd.AddCommand(countdownCmd)

	countdownCmd.Flags().StringVarP(&duration, "duration", "d", "", "timre duration (required)")
	countdownCmd.Flags().StringVarP(&message, "message", "m", "Time's up!", "message to display when time completes")
	countdownCmd.Flags().BoolVar(&sound, "sound", true, "play sound when timer completes")
	countdownCmd.Flags().BoolVar(&desktop, "desktop", true, "show desktop notification when timer completes")

	countdownCmd.MarkFlagRequired("duration")
}

func runCountdown(cmd *cobra.Command, args []string) error {
	d, err := time.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("invalid duration format '%s': %v\nValid examples: 30s, 5m, 1h, 1h30m", duration, err)
	}

	if d <= 0 {
		return fmt.Errorf("duration must be positive")
	}

	config := GetConfig()
	if cmd.Flags().Changed("sound") {
		config.Sound.Enabled = sound && !quiet
	}
	if cmd.Flags().Changed("desktop") {
		config.Notifications.Desktop = desktop && !quiet
	}

	timerConfig := &types.TimerConfig{
		Duration: d,
		Message: message,
		Sound: config.Sound.Enabled,
		Desktop: config.Notifications.Desktop,
		Terminal: config.Notifications.Terminal,
		ShowProgress: !quiet,
		Verbose: verbose,
	}

	countdown := timer.NewCountdown(timerConfig)
	return countdown.Start()
}