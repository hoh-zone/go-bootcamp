package patterns

import (
	"context"
	"errors"
	"sync"
	"time"
)

// PipelineDoubleThenAdd builds a two-stage pipeline:
// stage1: x*2, stage2: x+1. It returns the output channel.
func PipelineDoubleThenAdd(in <-chan int) <-chan int {
	stage1 := func(input <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for v := range input {
				out <- v * 2
			}
		}()
		return out
	}
	stage2 := func(input <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for v := range input {
				out <- v + 1
			}
		}()
		return out
	}
	return stage2(stage1(in))
}

// FanOutSquare starts workerCount goroutines to square numbers from in and merges them into one output channel.
// It closes the output when all workers finish or ctx is done.
func FanOutSquare(ctx context.Context, in <-chan int, workerCount int) <-chan int {
	out := make(chan int)
	if workerCount <= 0 {
		close(out)
		return out
	}

	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		for v := range in {
			select {
			case <-ctx.Done():
				return
			case out <- v * v:
			}
		}
	}

	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// SendWithTimeout tries to send v into ch, respecting ctx and timeout.
// It returns ctx.Err(), deadline exceeded, or nil on success.
func SendWithTimeout(ctx context.Context, ch chan<- int, v int, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = time.Nanosecond
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case ch <- v:
		return nil
	case <-timer.C:
		return errors.New("send timeout")
	}
}
