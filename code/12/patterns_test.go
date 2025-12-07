package patterns

import (
	"context"
	"errors"
	"sort"
	"testing"
	"time"
)

func TestPipelineDoubleThenAdd(t *testing.T) {
	in := make(chan int, 3)
	for _, v := range []int{1, 2, 3} {
		in <- v
	}
	close(in)

	var got []int
	for v := range PipelineDoubleThenAdd(in) {
		got = append(got, v)
	}
	want := []int{3, 5, 7}
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want %d", i, got[i], want[i])
		}
	}
}

func TestFanOutSquare(t *testing.T) {
	ctx := context.Background()
	in := make(chan int, 4)
	for _, v := range []int{1, 2, 3, 4} {
		in <- v
	}
	close(in)

	var got []int
	for v := range FanOutSquare(ctx, in, 2) {
		got = append(got, v)
	}
	sort.Ints(got)
	want := []int{1, 4, 9, 16}
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want %d", i, got[i], want[i])
		}
	}
}

func TestSendWithTimeout(t *testing.T) {
	ctx := context.Background()
	ch := make(chan int, 1)

	if err := SendWithTimeout(ctx, ch, 42, 10*time.Millisecond); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := <-ch; got != 42 {
		t.Fatalf("got %d want 42", got)
	}

	// now channel is full; expect timeout
	ch <- 1
	err := SendWithTimeout(ctx, ch, 2, 5*time.Millisecond)
	if err == nil {
		t.Fatalf("expected timeout error")
	}
	if !errors.Is(err, context.DeadlineExceeded) && err.Error() != "send timeout" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendWithTimeoutContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ch := make(chan int)
	if err := SendWithTimeout(ctx, ch, 1, 5*time.Millisecond); err == nil {
		t.Fatalf("expected ctx error")
	}
}
