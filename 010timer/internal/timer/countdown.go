package timer

import (
	"fmt"
	"time"

	"github.com/samnart1/golang/010timer/internal/notifications"
	"github.com/samnart1/golang/010timer/pkg/types"
)

type Countdaown struct {
	*BaseTimer
	config *types.TimerConfig
	remaining time.Duration
	display *Display
	notifier *notifications.Notifier
}

func NewCountdown(config *types.TimerConfig) *Countdaown {
	return &Countdaown{
		BaseTimer: NewBaseTimer(),
		config: config,
		remaining: config.Duration,
		display: NewDisplay(config.ShowProgress, config.Verbose),
		notifier: notifications.NewNotifier(),
	}
}

func (c *Countdaown) Start() error {
	ctx := setupSignalHandler()

	if c.config.Verbose {
		fmt.Printf("Starting countdown timer for %v\n", c.config.Duration)
	}

	c.start()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	c.display.ShowStart(c.config.Duration)

	for {
		select {
		case <-ctx.Done():
			c.stop()
			return nil
		case <-ticker.C:
			elapsed := c.GetElapsed()
			c.remaining = c.config.Duration - elapsed

			if c.remaining <= 0 {
				return c.complete()
			}

			if c.config.ShowProgress {
				c.display.ShowProgress(c.remaining, c.config.Duration)
			}
		}
	}
}

func (c *Countdaown) Stop() error {
	c.stop()
	return nil
}

func (c *Countdaown) Pause() error {
	c.pause()
	return nil
}

func (c *Countdaown) Resume() error {
	c.resume()
	return nil
}

func (c *Countdaown) complete() error {
	c.state = types.StateCompleted

	c.display.ShowComplete(c.config.Message)

	if c.config.Sound {
		c.notifier.PlaySound()
	}

	if c.config.Desktop {
		c.notifier.ShowDesktop("Timer Complete", c.config.Message)
	}

	if c.config.Terminal {
		c.notifier.ShowTerminal(c.config.Message)
	}

	return nil
}