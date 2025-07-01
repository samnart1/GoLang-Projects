package timer

import (
	"fmt"
	"time"

	"github.com/samnart1/golang/010timer/internal/notifications"
	"github.com/samnart1/golang/010timer/pkg/types"
)

type Alarm struct {
	*BaseTimer
	config 		*types.AlarmConfig
	targetTime 	time.Time
	display 	*Display
	notifier 	*notifications.Notifier
}

func NewAlarm(config *types.AlarmConfig) *Alarm {
	return &Alarm{
		BaseTimer: NewBaseTimer(),
		config: config,
		targetTime: config.TargetTime,
		display: NewDisplay(config.ShowProgress, config.Verbose),
		notifier: notifications.NewNotifier(),
	}
}

func (a *Alarm) Start() error {
	ctx := setupSignalHandler()

	if a.config.Verbose {
		fmt.Printf("Alarm set for %v\n", a.targetTime.Format("15:04:05"))
		fmt.Printf("Time until alarm: %v\n", a.config.Duration)
	}

	a.Start()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	a.display.ShowAlarmStart(a.targetTime, a.config.Duration)

	for {
		select {
		case  <-ctx.Done():
			a.stop()
			return nil
		case <-ticker.C:
			now := time.Now()

			if now.After(a.targetTime) || now.Equal(a.targetTime) {
				return a.trigger()
			}

			remaining := a.targetTime.Sub(now)
			if a.config.ShowProgress {
				a.display.ShowAlarmProgress(remaining, a.targetTime)
			}
		}
	}
}

func (a *Alarm) Stop() error {
	a.stop()
	return nil
}

func (a *Alarm) Pause() error {
	return fmt.Errorf("alarms cannot be paused")
}

func (a *Alarm) Resume() error {
	return fmt.Errorf("allarms cannot be resumed")
}

func (a *Alarm) trigger() error {
	a.state = types.StateCompleted

	a.display.ShowAlarmTriggered(a.config.Message)

	if a.config.Sound {
		a.notifier.PlaySound()
	}

	if a.config.Desktop {
		a.notifier.ShowDesktop("Alarm", a.config.Message)
	}

	if a.config.Terminal {
		a.notifier.ShowTerminal(a.config.Message)
	}

	return nil
}