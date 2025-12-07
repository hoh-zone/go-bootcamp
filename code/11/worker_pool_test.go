package concurrency

import (
	"context"
	"testing"
	"time"
)

func TestProcessWithPool(t *testing.T) {
	ctx := context.Background()
	inputs := []int{1, 2, 3, 4, 5}
	workers := 2

	results, err := ProcessWithPool(ctx, inputs, workers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []int{1, 4, 9, 16, 25}
	for i, v := range expected {
		if results[i] != v {
			t.Fatalf("result[%d]=%d want %d", i, results[i], v)
		}
	}
}

func TestProcessWithPoolCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	_, err := ProcessWithPool(ctx, []int{1, 2, 3}, 2)
	if err == nil {
		t.Fatalf("expected timeout error")
	}
}

func TestProcessWithPoolInvalidWorkerCount(t *testing.T) {
	_, err := ProcessWithPool(context.Background(), []int{1}, 0)
	if err == nil {
		t.Fatalf("expected error for non-positive workers")
	}
}
