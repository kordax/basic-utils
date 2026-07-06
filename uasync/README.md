# uasync

Async execution helpers built around `context.Context`.

Includes futures, async tasks, scheduled tasks, retry helpers, and grouped task execution with concurrency limits.

```go
result, err := uasync.Retry(ctx, 3, time.Millisecond, func(ctx context.Context) (*int, error) {
	v := 42
	return &v, nil
})
```
