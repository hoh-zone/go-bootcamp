package contextdemo

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDoWithTimeoutSuccess(t *testing.T) {
	err := DoWithTimeout(context.Background(), 50*time.Millisecond, func(ctx context.Context) error {
		select {
		case <-time.After(10 * time.Millisecond):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoWithTimeoutExceeded(t *testing.T) {
	err := DoWithTimeout(context.Background(), 5*time.Millisecond, func(ctx context.Context) error {
		select {
		case <-time.After(50 * time.Millisecond):
			return errors.New("should not happen")
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
}
