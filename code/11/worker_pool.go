package concurrency

import (
	"context"
	"errors"
	"sync"
)

// ProcessWithPool runs jobs concurrently with at most `workers` goroutines.
// It returns squared results in the same order as inputs.
func ProcessWithPool(ctx context.Context, inputs []int, workers int) ([]int, error) {
	if workers <= 0 {
		return nil, errors.New("workers must be positive")
	}
	type job struct {
		idx int
		val int
	}

	jobs := make(chan job)
	results := make([]int, len(inputs))
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for j := range jobs {
			select {
			case <-ctx.Done():
				return
			default:
			}
			results[j.idx] = j.val * j.val
		}
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}

	for i, v := range inputs {
		if err := ctx.Err(); err != nil {
			close(jobs)
			wg.Wait()
			return results, err
		}
		select {
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			return results, ctx.Err()
		case jobs <- job{idx: i, val: v}:
		}
	}
	close(jobs)
	wg.Wait()
	if err := ctx.Err(); err != nil {
		return results, err
	}
	return results, nil
}
