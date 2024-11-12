/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uasync

import (
	"context"
	"testing"
	"time"
)

// Benchmark for ExecuteAsync method of AsyncTask
func BenchmarkAsyncTaskExecuteAsync(b *testing.B) {
	sampleFn := func(ctx context.Context) (*int, error) {
		val := 5
		return &val, nil
	}
	ctx := context.Background()
	task := NewAsyncTask(ctx, sampleFn, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task.ExecuteAsync()
	}
}

// Benchmark for Wait method of AsyncTask
func BenchmarkAsyncTaskWait(b *testing.B) {
	sampleFn := func(ctx context.Context) (*int, error) {
		val := 5
		return &val, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	task := NewAsyncTask(ctx, sampleFn, 0)
	task.ExecuteAsync()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = task.Wait()
	}
}
