package concurrency

import (
	"sync"
	"testing"
)

func TestSafeCounter(t *testing.T) {
	var wg sync.WaitGroup
	c := &SafeCounter{}
	total := 1000

	wg.Add(total)
	for i := 0; i < total; i++ {
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()

	if got := c.Value(); got != total {
		t.Fatalf("counter = %d, want %d", got, total)
	}
}
