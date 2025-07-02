package timer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/samnart1/golang/010timer/pkg/types"
)

type Stopwatch struct {
	*BaseTimer
	config *types.StopWatchConfig
	laps []types.LapTime
	display *Display
	lapCount int
}

func NewStopwatch(config *types.StopWatchConfig) *Stopwatch {
	return &Stopwatch{
		BaseTimer: NewBaseTimer(),
		config: config,
		laps: make([]types.LapTime, 0),
		display: NewDisplay(true, config.Verbose),
		lapCount: 0,
	}
}

func (s *Stopwatch) Start() error {
	fmt.Println("Stopwatch started! Controls:")
	fmt.Println("	SPACE: Start/Stop")
	if s.config.ShowLaps {
		fmt.Println("	L: Record lap")
	}
	fmt.Println("	R: Reset")
	fmt.Println("	Q: Quit")
	fmt.Println()

	s.start()

	go s.updateDisplay()

	return s.handleInput()
}

func (s *Stopwatch) Stop() error {
	s.stop()
	return nil
}

func (s *Stopwatch) Pause() error {
	s.pause()
	return nil
}

func (s *Stopwatch) Resume() error {
	s.resume()
	return nil
}

func (s *Stopwatch) RecordLap() {
	if s.state != types.StateRunning {
		return
	}

	s.lapCount++
	elapsed := s.GetElapsed()

	var lapTime time.Duration
	if len(s.laps) > 0 {
		lapTime = elapsed - s.laps[len(s.laps)-1].Total
	} else {
		lapTime = elapsed
	}

	lap := types.LapTime{
		Number: s.lapCount,
		Time: lapTime,
		Total: elapsed,
		Created: time.Now(),
	}

	s.laps = append(s.laps, lap)

	if s.config.Verbose {
		fmt.Printf("Lap %d: %v (Total: %v)\n", lap.Number,
			s.display.FormatDuration(lapTime),
			s.display.FormatDuration(lap.Total))
	}
}

func (s *Stopwatch) Reset() {
	s.elapsed = 0
	s.laps = make([]types.LapTime, 0)
	s.lapCount = 0
	s.state = types.StateStopped

	if s.config.Verbose {
		fmt.Println("Stopwatch reset")
	}
}

func (s *Stopwatch) updateDisplay() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			elapsed := s.GetElapsed()
			s.display.ShowStopwatch(elapsed, s.state, s.laps)
		}
	}
}

func (s *Stopwatch) handleInput() error {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))

		switch input {
		case " ", "space":
			if s.state == types.StateRunning {}
		}
	}
}