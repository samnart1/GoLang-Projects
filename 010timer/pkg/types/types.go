package types

import (
	"fmt"
	"time"
)

type TimerConfig struct {
	Duration		time.Duration
	Message			string
	Sound			bool
	Desktop			bool
	Terminal		bool
	ShowProgress	bool
	Verbose			bool
}

type StopWatchConfig struct {
	ShowLaps	bool
	Verbose		bool
	Color		bool
}

type AlarmConfig struct {
	TargetTime		time.Time
	Duration		time.Duration
	Message			string
	Sound			bool
	Desktop			bool
	Terminal		bool
	ShowProgress	bool
	Verbose			bool
}

type LapTime struct {
	Number 	int
	Time	time.Duration
	Total	time.Duration
	Created	time.Time
}

type TimerState int

const (
	StateStopped TimerState = iota
	StateRunning
	StatePaused
	StateCompleted
)

func (s TimerState) String() string {
	switch s {
	case StateStopped:
		return "stopped"
	case StateRunning:
		return "running"
	case StatePaused:
		return "paused"
	case StateCompleted:
		return "completed"
	default:
		return "unknown"
	}
}

type NotificationType int

const (
	NotificationSound NotificationType = iota
	NotificationDesktop 
	NotificationTerminal
)

func (n NotificationType) String() string {
	switch n {
	case NotificationSound:
		return "sound"
	case NotificationDesktop:
		return "desktop"
	case NotificationTerminal:
		return "terminal"
	default:
		return "unknown"
	}
}

type TimerStats struct {
	StartTime 	time.Time
	EndTime		time.Time
	Duration	time.Duration
	State		TimerState
	Pauses		int
	TotalPaused	time.Duration
}

type NotificationPreferences struct {
	Sound struct {
		Enabled		bool
		File		string
		Volume		int
		RepeatCount int
	}
	Desktop struct {
		Enabled	bool
		Style	string
		Timeout	int
	}
	Terminal struct {
		Enabled 	bool
		UseColor	bool
		Style		string
	}
}

type TimerMetadata struct {
	ID			string
	CreatedAt	time.Time
	UpdatedAt	time.Time
	CreatedBy	string
	Description	string
	Tags		[]string
}

type ValidationError struct {
	Field 	string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type TimerHistory struct {
	Timer		*TimerMetadata
	Stats		*TimerStats
	StartedAt	time.Time
	CompletedAt	time.Time
	Success		bool
	Error		string
}

func NewTimerConfig() *TimerConfig {
	return &TimerConfig{
		Sound: true,
		Desktop: true,
		Terminal: true,
		ShowProgress: true,
		Verbose: false,
	}
}

func NewStopwatchConfig() *StopWatchConfig {
	return &StopWatchConfig{
		ShowLaps: false,
		Verbose: false,
		Color: true,
	}
}

func NewAlarmConfig() *AlarmConfig {
	return &AlarmConfig{
		Sound: true,
		Desktop: true,
		Terminal: true,
		ShowProgress: true,
		Verbose: false,
	}
}

func (c *TimerConfig) Validate() error {
	if c.Duration <= 0 {
		return &ValidationError{
			Field: "Duration",
			Message: "must be greater than 0",
		}
	}
	return nil
}

func (c *AlarmConfig) Validate() error {
	if c.TargetTime.Before(time.Now()) {
		return &ValidationError{
			Field: "TargetTime",
			Message: "must be in the future",
		}
	}
	return nil
}