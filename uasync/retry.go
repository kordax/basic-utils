package uasync

import (
	"context"
	"fmt"
	"time"
)

// Retry executes fn until it succeeds, context is canceled, or attempts are exhausted.
// attempts <= 0 means retry forever.
func Retry[R any](ctx context.Context, attempts int, delay time.Duration, fn func(context.Context) (*R, error)) (*R, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	var lastErr error
	for attempt := 1; attempts <= 0 || attempt <= attempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		result, err := fn(ctx)
		if err == nil {
			return result, nil
		}
		lastErr = err

		if attempts > 0 && attempt >= attempts {
			break
		}
		if delay > 0 {
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			case <-timer.C:
			}
		}
	}

	return nil, fmt.Errorf("retry attempts exhausted: %w", lastErr)
}
