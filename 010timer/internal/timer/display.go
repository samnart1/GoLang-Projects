package timer

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/samnart1/golang/010timer/pkg/types"
)

type Display struct {
	showProgress 	bool
	verbose			bool
	spinner			*spinner.Spinner
	lastLine		string
}

func NewDisplay(showProgress, verbose bool) *Display {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " "

	return &Display{
		showProgress: showProgress,
		verbose: verbose,
		spinner: s,
	}
}

func (d *Display) ShowStart(duration time.Duration) {
	if d.verbose {
		color.Green("Timer started for %v", d.FormatDuration(duration))
	}
	if d.showProgress {
		d.spinner.Start()
	}
}

func (d *Display) ShowProgress(remaining, total time.Duration) {
	if !d.showProgress {
		return
	}

	progress := float64(total-remaining) / float64(total)
	progressBar := d.createProgressBar(progress, 30)

	line := fmt.Sprintf("\r%s %s remaining: %s",
		progressBar,
		color.CyanString("Time"),
		color.YellowString(d.FormatDuration(remaining)))

		if line != d.lastLine {
			fmt.Print(line)
			d.lastLine = line
		}
}

func (d *Display) ShowComplete(message string) {
	if d.showProgress {
		d.spinner.Stop()
		fmt.Print("\r" + strings.Repeat(" ", len(d.lastLine)) + "\r")
	}

	color.Red("%s", message)
	fmt.Println()
}

func (d *Display) ShowStopwatch(elapsed time.Duration, state types.TimerState, laps []types.LapTime) {
	stateColor := color.GreenString
	stateIcon := "▶️"

	switch state {
	case types.StatePaused:
		stateColor = color.YellowString
		stateIcon = "⏸️"
	case types.StateStopped:
		stateColor = color.RedString
		stateIcon = "⏹️"
	}

	line := fmt.Sprintf("\r%s %s: %s",
		stateIcon,
		stateColor("Elapsed"),
		color.CyanString(d.FormatDurationPrecise(elapsed)))

	if len(laps) > 0 {
		lastLap := laps[len(laps)-1]
		line += fmt.Sprintf(" | Last lap: %s",
			color.MagentaString(d.FormatDuration(lastLap.Time)))
	}

	if line != d.lastLine {
		fmt.Print(line)
		d.lastLine = line
	}
}

func (d *Display) ShowAlarmStart(targetTime time.Time, duration time.Duration) {
	if d.verbose {
		color.Blue("Alarm set for %v", targetTime.Format("15:04:05"))
		color.Cyan("Time until alarm: %v", d.FormatDuration(duration))
	}
}

func (d *Display) ShowAlarmProgress(remaining time.Duration, targetTime time.Time) {
	line := fmt.Sprintf("\r Alarm at %s | Time remaining: %s",
		color.BlueString(targetTime.Format("15:04:05")),
		color.YellowString(d.FormatDuration(remaining)))

	if line != d.lastLine {
		fmt.Print(line)
		d.lastLine = line
	}
}

func (d *Display) ShowAlarmTriggered(message string) {
	fmt.Print("\r" + strings.Repeat(" ", len(d.lastLine)) + "\r")
	color.Red("ALARM: %s", message)
	fmt.Println()
}

func (d *Display) FormatDuration(duration time.Duration) string {
	if duration < 0 {
		return "00:00"
	}

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func (d *Display) FormatDurationPrecise(duration time.Duration) string {
	if duration < 0 {
		return "00:00.000"
	}

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	milliseconds := int(duration.Milliseconds()) % 1000

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
	}
	return fmt.Sprintf("%02d:%02d.%03d", minutes, seconds, milliseconds)
}

func (d *Display) createProgressBar(progress float64, width int) string {
	if progress > 1.0 {
		progress = 1.0
	}
	if progress < 0.0 {
		progress = 0.0
	}

	filled := int(progress * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

	percentage := int(progress * 100)
	return fmt.Sprintf("[%s] %3d%%", color.GreenString(bar), percentage) 
}