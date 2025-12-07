package concurrency

import (
	"context"
	"time"
)

// DoWithTimeout runs fn with a derived context that times out.
func DoWithTimeout(parent context.Context, timeout time.Duration, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	done := make(chan error, 1)

	go func() {
		done <- fn(ctx)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
