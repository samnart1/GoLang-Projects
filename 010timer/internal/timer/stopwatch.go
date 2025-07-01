package timer

import "github.com/samnart1/golang/010timer/pkg/types"

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