package timer

import (
	"fmt"
	"time"
)

type Timer struct {
	duration 	time.Duration
	done 		chan bool
}

func New(duration time.Duration) *Timer {
	return &Timer{
		duration: 	duration,
		done: 		make(chan bool),
	}
}

func (t *Timer) Start() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	deadline := time.Now().Add(t.duration)

	for {
		select {
		case <- t.done:
			return
		case <- ticker.C:
			remaining := time.Until(deadline)
			if remaining <= 0 {
				return
			}

			// show countdown for last 30 seconds
			if remaining <= 30*time.Second && remaining > 20*time.Second {
				fmt.Printf("%v remaining!\n", remaining.Round(time.Second))
			} else if remaining <= 10*time.Second {
				fmt.Printf("Only %v left!\n", remaining.Round(time.Second))
			}
		}
	}
}

func (t *Timer) Stop() {
	if t.done != nil {
		select {
		case <-t.done:
			// Channel already closed
		default:
			close(t.done)
		}
	}
}