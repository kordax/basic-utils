/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ustream

import (
	"fmt"
	"sync"

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

func (s *TerminalStream[T]) Collect() []T {
	return s.values
}

// CollectToMap uses the uarray.ToMultiMap function to collect elements of the Stream into a map.
func (s *TerminalStream[T]) CollectToMap(mapper func(*T) (any, T)) map[any][]T {
	return uarray.ToMultiMap(s.values, mapper)
}
