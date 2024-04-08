/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ustream

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kordax/basic-utils/uarray"
)

type parallelTask[T any] struct {
	v     *T
	index int
}

// Collector defines the interface for collecting elements from a stream.
type Collector[T any] interface {
	Collect() []T
	CollectToMap(func(*T) (any, any)) map[any]any
}

type Stream[T any] struct {
	values []T
}

// NewStream creates a new Stream from the given slice.
func NewStream[T any](values []T) *Stream[T] {
	return &Stream[T]{values: values}
}

// Filter wraps the uarray.Filter function in the stream API.
func (s *Stream[T]) Filter(predicate func(*T) bool) *Stream[T] {
	return NewStream(uarray.Filter(s.values, predicate))
}

// FilterOut wraps the uarray.FilterOut function in the stream API.
func (s *Stream[T]) FilterOut(predicate func(*T) bool) *Stream[T] {
	return NewStream(uarray.FilterOut(s.values, predicate))
}

// Map wraps the uarray.Map function in the stream API.
// This operation can only be the last in a pipeline.
func (s *Stream[T]) Map(mapper func(*T) any) *TerminalStream[any] {
	return NewTerminalStream(uarray.Map(s.values, mapper))
}

func (s *Stream[T]) Collect() []T {
	return s.values
}

// CollectToMap uses the uarray.ToMultiMap function to collect elements of the Stream into a map.
func (s *Stream[T]) CollectToMap(mapper func(*T) (any, any)) map[any]any {
	return uarray.ToMap(s.values, mapper)
}

func (s *Stream[T]) ToTerminal() *TerminalStream[T] {
	return NewTerminalStream(s.values)
}

// TerminalStream represents a stream that can only be collected.
type TerminalStream[T any] struct {
	values []T
}

func NewTerminalStream[T any](values []T) *TerminalStream[T] {
	return &TerminalStream[T]{values: values}
}

// ParallelExecute executes the given function concurrently on each element of the stream's values
// using the specified level of parallelism.
func (s *TerminalStream[T]) ParallelExecute(fn func(int, *T), parallelism int) {
	if parallelism == 0 {
		panic(fmt.Errorf("parallelism cannot be zero"))
	}

	var wg sync.WaitGroup
	in := make(chan parallelTask[T])

	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for v := range in {
				fn(v.index, v.v)
			}
		}(i)
	}

	subsetSize := (len(s.values) + parallelism - 1) / parallelism
	startIndex := 0
	for i := 0; i < parallelism; i++ {
		endIndex := startIndex + subsetSize
		if endIndex > len(s.values) {
			endIndex = len(s.values)
		}
		subset := s.values[startIndex:endIndex]
		for j, value := range subset {
			in <- parallelTask[T]{v: &value, index: startIndex + j}
		}
		startIndex = endIndex
	}

	close(in)
	wg.Wait()
}

// ParallelExecuteWithTimeout executes a batch of tasks in parallel with a specified level of parallelism and timeout.
// This method divides the work into batches according to the parallelism parameter and executes each task asynchronously.
// Each task will be cancelled if it does not complete within the specified timeout duration.
//
// Parameters:
// - fn: The function to execute for each item in the TerminalStream. This function receives an index and a pointer to the item.
// - cancel: The cancel function to call for each item if it exceeds the timeout. This function receives an index and a pointer to the item.
// - timeout: The maximum duration to wait for each task to complete before cancelling it. If a task exceeds this duration, it is considered failed.
// - parallelism: The maximum number of tasks to execute concurrently. This controls the level of parallelism and helps manage resource utilization.
//
// Panics:
// This method panics if the parallelism parameter is less than or equal to zero, as it indicates an invalid configuration.
//
// Usage Example:
// Assuming a TerminalStream of some data type, you can process each item in parallel, with a specific timeout and level of parallelism:
//
//	s := NewTerminalStream[MyType](...your data...)
//	s.ParallelExecuteWithTimeout(func(i int, item *MyType) {
//	    // Process item
//	}, func(i int, item *MyType) {
//	    // Cancel item processing
//	}, 5*time.Second, 10) // Timeout of 5 seconds, with a parallelism of 10
//
// Note: The actual processing function (fn) does not return a value. If you need to collect results or errors from each task,
// you might need to use a different approach or modify the method accordingly.
func (s *TerminalStream[T]) ParallelExecuteWithTimeout(fn func(int, *T), cancel func(int, *T), timeout time.Duration, parallelism int) {
	if parallelism == 0 {
		panic(fmt.Errorf("parallelism cannot be zero"))
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	var wg sync.WaitGroup
	in := make(chan parallelTask[T], parallelism)

	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
		cycle:
			v, ok := <-in
			if !ok {
				return
			}

			select {
			case <-ctx.Done():
				cancel(v.index, v.v)
				goto cycle
			default:
				fn(v.index, v.v)
				goto cycle
			}
		}(i)
	}

	subsetSize := (len(s.values) + parallelism - 1) / parallelism
	startIndex := 0
	for i := 0; i < parallelism; i++ {
		endIndex := startIndex + subsetSize
		if endIndex > len(s.values) {
			endIndex = len(s.values)
		}
		subset := s.values[startIndex:endIndex]
		for j, value := range subset {
			in <- parallelTask[T]{v: &value, index: startIndex + j}
		}
		startIndex = endIndex
	}
	close(in)
	wg.Wait()
}

func (s *TerminalStream[T]) Collect() []T {
	return s.values
}

// CollectToMap uses the uarray.ToMultiMap function to collect elements of the Stream into a map.
func (s *TerminalStream[T]) CollectToMap(mapper func(*T) (any, T)) map[any][]T {
	return uarray.ToMultiMap(s.values, mapper)
}
