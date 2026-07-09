package uasync

import (
	"context"
	"testing"
	"time"
)

func BenchmarkFutureCompleteAndWait(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		f := NewFuture[int](ctx)
		value := i
		f.Complete(&value)
		_, _ = f.Wait()
	}
}

func BenchmarkFutureWaitAlreadyCompleted(b *testing.B) {
	f := NewFuture[int](context.Background())
	value := 42
	f.Complete(&value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = f.Wait()
	}
}

func BenchmarkExecute(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		_, _ = Execute(ctx, func(ctx context.Context) (*int, error) {
			value := i
			return &value, nil
		})
	}
}

func BenchmarkAsyncTaskExecuteAndWait(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		task := NewAsyncTask(ctx, func(ctx context.Context) (*int, error) {
			value := i
			return &value, nil
		}, 0)
		task.ExecuteAsync()
		_, _ = task.Wait()
	}
}

func BenchmarkScheduledTaskImmediate(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		task := NewScheduledTask(ctx, time.Now(), func(ctx context.Context) (*int, error) {
			value := i
			return &value, nil
		}, 0)
		f := task.Schedule()
		_, _ = f.Wait()
	}
}
