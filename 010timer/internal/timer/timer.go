package timer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samnart1/golang/010timer/pkg/types"
)

type Timer interface {
	Start() error
	Stop() error
	Pause() error
	Resume() error
}

type BaseTimer struct {
	state types.TimerState
	startTime time.Time
	pauseTime time.Time
	elapsed time.Duration
	ctx context.Context
	cancel context.CancelFunc
}

func NewBaseTimer() *BaseTimer {
	ctx, cancel := context.WithCancel(context.Background())
	return &BaseTimer{
		state: types.StateStopped,
		ctx: ctx,
		cancel: cancel,
	}
}

func (bt *BaseTimer) GetState() types.TimerState {
	return bt.state
}

func (bt *BaseTimer) GetElapsed() time.Duration {
	switch bt.state {
	case types.StateRunning:
		return bt.elapsed + time.Since(bt.startTime)
	case types.StatePaused:
		return bt.elapsed + bt.pauseTime.Sub(bt.startTime)
	default:
		return bt.elapsed
	}
}

func (bt *BaseTimer) start() {
	bt.state = types.StateRunning
	bt.startTime = time.Now()
}

func (bt *BaseTimer) stop() {
	bt.state = types.StateStopped
	bt.elapsed = 0
	bt.cancel()
}

func (bt *BaseTimer) pause() {
	if bt.state == types.StateRunning {
		bt.state = types.StatePaused
		bt.pauseTime = time.Now()
		bt.elapsed += bt.pauseTime.Sub(bt.startTime)
	}
}

func (bt *BaseTimer) resume() {
	if bt.state == types.StatePaused {
		bt.state = types.StateRunning
		bt.startTime = time.Now()
	}
}

func setupSignalHandler() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nReceived interrupt signal, shutting down...")
		cancel()
	}()

	return ctx
}