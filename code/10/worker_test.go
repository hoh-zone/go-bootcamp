package contextdemo

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestProcessAllStopsOnCancel(t *testing.T) {
	jobs := make(chan int, 3)
	for i := 0; i < 3; i++ {
		jobs <- i
	}
	close(jobs)

	ctx, cancel := context.WithCancel(context.Background())
	var handled int32

	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	err := ProcessAll(ctx, jobs, func(ctx context.Context, v int) error {
		atomic.AddInt32(&handled, 1)
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	if err == nil {
		t.Fatalf("expected context cancellation error")
	}
	if handled == 0 {
		t.Fatalf("expected to handle some jobs before cancel")
	}
	if handled == 3 {
		t.Fatalf("expected cancel to stop early, handled=%d", handled)
	}
}
