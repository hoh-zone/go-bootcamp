package contextdemo

import (
	"context"
	"time"
)

// DoWithTimeout runs fn with a derived timeout context.
// fn should honor ctx.Done(); returns fn error or ctx error.
func DoWithTimeout(parent context.Context, timeout time.Duration, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	return fn(ctx)
}
