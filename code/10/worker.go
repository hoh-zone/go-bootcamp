package contextdemo

import "context"

// ProcessAll consumes jobs until channel关闭或 context取消。
// handle 应该尊重 ctx；当 ctx 取消时返回 ctx.Err()。
func ProcessAll(ctx context.Context, jobs <-chan int, handle func(context.Context, int) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case j, ok := <-jobs:
			if !ok {
				return nil
			}
			if err := handle(ctx, j); err != nil {
				return err
			}
		}
	}
}
