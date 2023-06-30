/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package async_utils

import (
	"context"
	"testing"
	"time"
)

// Benchmark for NewAsyncTask function
func BenchmarkNewAsyncTask(b *testing.B) {
	sampleFn := func() (*int, error) {
		val := 5
		return &val, nil
	}
	cancelFunc := context.CancelFunc(func() {})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewAsyncTask(sampleFn, cancelFunc, 0)
	}
}

// Benchmark for ExecuteAsync method of AsyncTask
func BenchmarkAsyncTaskExecuteAsync(b *testing.B) {
	sampleFn := func() (*int, error) {
		val := 5
		return &val, nil
	}
	cancelFunc := context.CancelFunc(func() {})
	task := NewAsyncTask(sampleFn, cancelFunc, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task.ExecuteAsync()
	}
}

// Benchmark for Wait method of AsyncTask
func BenchmarkAsyncTaskWait(b *testing.B) {
	sampleFn := func() (*int, error) {
		val := 5
		return &val, nil
	}
	cancelFunc := context.CancelFunc(func() {})
	task := NewAsyncTask(sampleFn, cancelFunc, 0)
	task.ExecuteAsync()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = task.Wait(100 * time.Millisecond)
	}
}

// Benchmark for Cancel method of AsyncTask
func BenchmarkAsyncTaskCancel(b *testing.B) {
	sampleFn := func() (*int, error) {
		val := 5
		return &val, nil
	}
	cancelFunc := context.CancelFunc(func() {})
	task := NewAsyncTask(sampleFn, cancelFunc, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task.Cancel()
	}
}
