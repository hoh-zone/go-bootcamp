package concurrency

import "sync"

// SafeCounter protects an integer with a mutex to avoid data races.
type SafeCounter struct {
	mu sync.Mutex
	n  int
}

// Inc increments the counter safely.
func (c *SafeCounter) Inc() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}

// Value returns the current count.
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}
