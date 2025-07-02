package cmd

import (
	"fmt"
	"time"

	"github.com/samnart1/golang/010timer/internal/timer"
	"github.com/samnart1/golang/010timer/pkg/types"
	"github.com/spf13/cobra"
)

var (
	alarmTime string
)

var alarmCmd = &cobra.Command{
	Use: "alarm",
	Short: "Set an alarm for a specific time",
	Long: `Set an alarm that will trigger at a specific time.
	
		Time format: HH:MM (24-hour format)
		
		Examples:
			timer alarm -t 14:30
			timer alarm -t 09:00 -m "Meeting!"`,

	RunE: runAlarm,
}

func init() {
	rootCmd.AddCommand(alarmCmd)

	alarmCmd.Flags().StringVarP(&alarmTime, "time", "t", "", "alarm time in HH:MM format (required)")
	alarmCmd.Flags().StringVarP(&message, "message", "m", "Alarm", "message to display wheh alarm triggers")
	alarmCmd.Flags().BoolVar(&sound, "sound", true, "play sound when alarm triggers")
	alarmCmd.Flags().BoolVar(&desktop, "desktop", true, "show desktop notification when alarm triggers")

	alarmCmd.MarkFlagRequired("time")
}

func runAlarm(cmd *cobra.Command, args []string) error {
	targetTime, err := time.Parse("15:04", alarmTime)
	if err != nil {
		return fmt.Errorf("invalid time format '%s': use HH:MM (24-hour format)", alarmTime)
	}

	now := time.Now()
	target := time.Date(now.Year(), now.Month(), now.Day(), targetTime.Hour(), targetTime.Minute(), 0, 0, now.Location())

	if target.Before(now) {
		target = target.Add(24 * time.Hour)
	}

	duration := target.Sub(now)

	config := GetConfig()
	if cmd.Flags().Changed("sound") {
		config.Notifications.Desktop = desktop && !quiet
	}

	alarmConfig := &types.AlarmConfig{
		TargetTime: target,
		Duration: duration,
		Message: message,
		Sound: config.Sound.Enabled,
		Desktop: config.Notifications.Desktop,
		Terminal: config.Notifications.Terminal,
		ShowProgress: !quiet,
		Verbose: verbose,
	}

	alarm := timer.NewAlarm(alarmConfig)
	return alarm.Start()
}