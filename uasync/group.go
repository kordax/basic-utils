package uasync

import (
	"context"
	"sync"
)

// RunGroup runs tasks with a concurrency limit and cancels the group on the first error.
func RunGroup(ctx context.Context, concurrency int, tasks ...func(context.Context) error) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if concurrency <= 0 || concurrency > len(tasks) {
		concurrency = len(tasks)
	}
	if concurrency == 0 {
		return nil
	}

	groupCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	taskCh := make(chan func(context.Context) error)
	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if groupCtx.Err() != nil {
					return
				}
				if err := task(groupCtx); err != nil {
					select {
					case errCh <- err:
						cancel()
					default:
					}
					return
				}
			}
		}()
	}

	for _, task := range tasks {
		select {
		case <-groupCtx.Done():
			close(taskCh)
			wg.Wait()
			select {
			case err := <-errCh:
				return err
			default:
				return groupCtx.Err()
			}
		case taskCh <- task:
		}
	}
	close(taskCh)
	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return ctx.Err()
	}
}
