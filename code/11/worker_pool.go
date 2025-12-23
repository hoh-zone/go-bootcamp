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

/*
results 没有数据竞争的原因：

- 每个元素只会被写一次：job{idx, val} 在生产时就绑定了唯一的 idx，jobs 通道的元素只能被某一个 worker 取到，因此同一 idx 不会被两个 goroutine 同时写。
- 主 goroutine 不读 results，直到 wg.Wait() 完成所有 worker；写入阶段只有 worker 在写，读取阶段只有主 goroutine 在读，读写不重叠。
- 取消场景下，worker 直接返回，不会和其他 worker 同时写同一个下标，只是可能留下默认值 0。
*/
