/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ustream

import (
	"cmp"
	"slices"
	"sync"
	"time"

	"github.com/kordax/basic-utils/v3/uarray"
)

// Collector defines the legacy interface for collecting elements from a stream.
type Collector[T any] interface {
	Collect() []T
	CollectToMap(func(T) (any, any)) map[any]any
}

// ResultCollector collects stream values into a custom result.
type ResultCollector[T, R any] interface {
	Collect(values []T) R
}

// ResultCollectorFunc adapts a function to ResultCollector.
type ResultCollectorFunc[T, R any] func(values []T) R

// Collect adapts a function to ResultCollector.
func (f ResultCollectorFunc[T, R]) Collect(values []T) R {
	return f(values)
}

// Stream wraps a slice and exposes fluent collection operations.
type Stream[T any] struct {
	values []T
}

// Empty creates a new empty Stream.
func Empty[T any]() *Stream[T] {
	return &Stream[T]{values: []T{}}
}

// From creates a new Stream from variadic values.
func From[T any](values ...T) *Stream[T] {
	return Of(values)
}

// Of creates a new Stream from the given slice.
func Of[T any](values []T) *Stream[T] {
	if values == nil {
		values = []T{}
	}

	return &Stream[T]{values: values}
}

// Generate creates a Stream by calling supplier count times.
func Generate[T any](count int, supplier func(index int) T) *Stream[T] {
	if count <= 0 {
		return Empty[T]()
	}

	values := make([]T, count)
	for i := 0; i < count; i++ {
		values[i] = supplier(i)
	}

	return Of(values)
}

// Iterate creates a Stream by repeatedly applying next to the previous value.
func Iterate[T any](seed T, count int, next func(T) T) *Stream[T] {
	if count <= 0 {
		return Empty[T]()
	}

	values := make([]T, count)
	values[0] = seed
	for i := 1; i < count; i++ {
		values[i] = next(values[i-1])
	}

	return Of(values)
}

// Filter keeps values matching predicate.
func (s *Stream[T]) Filter(predicate func(T) bool) *Stream[T] {
	if len(s.values) == 0 {
		return Empty[T]()
	}

	result := make([]T, 0, len(s.values))
	for _, value := range s.values {
		if predicate(value) {
			result = append(result, value)
		}
	}

	return Of(result)
}

// FilterOut removes values matching predicate.
func (s *Stream[T]) FilterOut(predicate func(T) bool) *Stream[T] {
	return s.Filter(func(value T) bool {
		return !predicate(value)
	})
}

// Map maps values into any and returns a terminal stream for legacy compatibility.
func (s *Stream[T]) Map(mapper func(T) any) *TerminalStream[any] {
	return NewTerminalStream(uarray.Map(s.values, mapper))
}

// Transform maps values without changing their type.
func (s *Stream[T]) Transform(mapper func(T) T) *Stream[T] {
	if len(s.values) == 0 {
		return Empty[T]()
	}

	result := make([]T, len(s.values))
	for i, value := range s.values {
		result[i] = mapper(value)
	}

	return Of(result)
}

// FlatMap expands each value into zero or more values of the same type.
func (s *Stream[T]) FlatMap(mapper func(T) []T) *Stream[T] {
	result := make([]T, 0, len(s.values))
	for _, value := range s.values {
		result = append(result, mapper(value)...)
	}

	return Of(result)
}

// Peek executes action for each value and returns the same stream values.
func (s *Stream[T]) Peek(action func(T)) *Stream[T] {
	for _, value := range s.values {
		action(value)
	}

	return s
}

// Limit keeps at most n first values.
func (s *Stream[T]) Limit(n int) *Stream[T] {
	if n <= 0 {
		return Empty[T]()
	}
	if n >= len(s.values) {
		return Of(s.values)
	}

	return Of(s.values[:n])
}

// Skip removes n first values.
func (s *Stream[T]) Skip(n int) *Stream[T] {
	if n <= 0 {
		return Of(s.values)
	}
	if n >= len(s.values) {
		return Empty[T]()
	}

	return Of(s.values[n:])
}

// TakeWhile keeps leading values while predicate returns true.
func (s *Stream[T]) TakeWhile(predicate func(T) bool) *Stream[T] {
	for i, value := range s.values {
		if !predicate(value) {
			return Of(s.values[:i])
		}
	}

	return Of(s.values)
}

// DropWhile skips leading values while predicate returns true.
func (s *Stream[T]) DropWhile(predicate func(T) bool) *Stream[T] {
	for i, value := range s.values {
		if !predicate(value) {
			return Of(s.values[i:])
		}
	}

	return Empty[T]()
}

// Reverse returns a reversed copy of stream values.
func (s *Stream[T]) Reverse() *Stream[T] {
	result := slices.Clone(s.values)
	slices.Reverse(result)

	return Of(result)
}

// Sort returns a sorted copy using less.
func (s *Stream[T]) Sort(less func(a, b T) bool) *Stream[T] {
	result := slices.Clone(s.values)
	slices.SortFunc(result, func(a, b T) int {
		switch {
		case less(a, b):
			return -1
		case less(b, a):
			return 1
		default:
			return 0
		}
	})

	return Of(result)
}

// DistinctBy keeps the first value for each string key.
func (s *Stream[T]) DistinctBy(key func(T) string) *Stream[T] {
	if len(s.values) == 0 {
		return Empty[T]()
	}

	seen := make(map[string]struct{}, len(s.values))
	result := make([]T, 0, len(s.values))
	for _, value := range s.values {
		k := key(value)
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		result = append(result, value)
	}

	return Of(result)
}

// CompactFunc removes values matching empty.
func (s *Stream[T]) CompactFunc(empty func(T) bool) *Stream[T] {
	return s.FilterOut(empty)
}

// ForEach executes action for each value.
func (s *Stream[T]) ForEach(action func(T)) {
	for _, value := range s.values {
		action(value)
	}
}

// Find returns the first value matching predicate.
func (s *Stream[T]) Find(predicate func(T) bool) *T {
	return uarray.Find(s.values, predicate)
}

// FindIndex returns the first index matching predicate, or -1.
func (s *Stream[T]) FindIndex(predicate func(T) bool) int {
	return uarray.FindIndex(s.values, predicate)
}

// AnyMatch checks whether any value matches predicate.
func (s *Stream[T]) AnyMatch(predicate func(T) bool) bool {
	return uarray.AnyMatch(s.values, predicate)
}

// AllMatch checks whether all values match predicate.
func (s *Stream[T]) AllMatch(predicate func(T) bool) bool {
	return uarray.AllMatch(s.values, predicate)
}

// NoneMatch checks whether no values match predicate.
func (s *Stream[T]) NoneMatch(predicate func(T) bool) bool {
	return !s.AnyMatch(predicate)
}

// Count returns stream size.
func (s *Stream[T]) Count() int {
	return len(s.values)
}

// Reduce folds values into one value.
func (s *Stream[T]) Reduce(initial T, reducer func(acc T, value T) T) T {
	return uarray.Reduce(s.values, initial, reducer)
}

// Collect returns stream values.
func (s *Stream[T]) Collect() []T {
	return s.values
}

// CollectCopy returns a copy of stream values.
func (s *Stream[T]) CollectCopy() []T {
	return slices.Clone(s.values)
}

// CollectToMap collects stream values into a map. Duplicate keys are overwritten.
func (s *Stream[T]) CollectToMap(mapper func(T) (any, any)) map[any]any {
	return uarray.ToMap(s.values, mapper)
}

// CollectToMultiMap collects stream values into a grouped map.
func (s *Stream[T]) CollectToMultiMap(mapper func(T) (any, any)) map[any][]any {
	return uarray.ToMultiMap(s.values, mapper)
}

// CollectWith collects stream values using a custom collector.
func (s *Stream[T]) CollectWith(collector func([]T) any) any {
	return collector(s.values)
}

// ToTerminal converts a stream to a terminal stream.
func (s *Stream[T]) ToTerminal() *TerminalStream[T] {
	return NewTerminalStream(s.values)
}

// TerminalStream represents a stream that can only be collected or executed.
type TerminalStream[T any] struct {
	values []T
}

// NewTerminalStream creates a terminal stream from values.
func NewTerminalStream[T any](values []T) *TerminalStream[T] {
	if values == nil {
		values = []T{}
	}

	return &TerminalStream[T]{values: values}
}

// ParallelExecute executes fn concurrently on each stream value.
func (s *TerminalStream[T]) ParallelExecute(fn func(int, *T), parallelism int) {
	if parallelism <= 0 {
		parallelism = 1
	}

	var wg sync.WaitGroup
	in := make(chan int)

	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range in {
				fn(index, &s.values[index])
			}
		}()
	}

	for i := range s.values {
		in <- i
	}

	close(in)
	wg.Wait()
}

// ParallelExecuteWithTimeout executes fn concurrently and calls cancel when an item exceeds timeout.
func (s *TerminalStream[T]) ParallelExecuteWithTimeout(fn func(int, T), cancel func(int, T), timeout time.Duration, parallelism int) {
	if parallelism <= 0 {
		parallelism = 1
	}

	var wg sync.WaitGroup
	in := make(chan int)

	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range in {
				value := s.values[index]
				if timeout <= 0 {
					if cancel != nil {
						cancel(index, value)
					}
					continue
				}

				done := make(chan struct{}, 1)
				go func() {
					fn(index, value)
					done <- struct{}{}
				}()

				timer := time.NewTimer(timeout)
				select {
				case <-done:
					if !timer.Stop() {
						<-timer.C
					}
				case <-timer.C:
					if cancel != nil {
						cancel(index, value)
					}
				}
			}
		}()
	}

	for i := range s.values {
		in <- i
	}
	close(in)
	wg.Wait()
}

// Collect returns terminal stream values.
func (s *TerminalStream[T]) Collect() []T {
	return s.values
}

// CollectCopy returns a copy of terminal stream values.
func (s *TerminalStream[T]) CollectCopy() []T {
	return slices.Clone(s.values)
}

// CollectToMap collects terminal stream values into a multimap.
func (s *TerminalStream[T]) CollectToMap(mapper func(T) (any, T)) map[any][]T {
	return uarray.ToMultiMap(s.values, mapper)
}

// Map maps stream values into another type.
func Map[T, R any](stream *Stream[T], mapper func(T) R) *Stream[R] {
	values := stream.Collect()
	result := make([]R, len(values))
	for i, value := range values {
		result[i] = mapper(value)
	}

	return Of(result)
}

// FlatMap maps stream values into another type and flattens the result.
func FlatMap[T, R any](stream *Stream[T], mapper func(T) []R) *Stream[R] {
	values := stream.Collect()
	result := make([]R, 0, len(values))
	for _, value := range values {
		result = append(result, mapper(value)...)
	}

	return Of(result)
}

// Reduce folds stream values into any result type.
func Reduce[T, R any](stream *Stream[T], initial R, reducer func(acc R, value T) R) R {
	return uarray.Reduce(stream.Collect(), initial, reducer)
}

// Collect applies a typed collector to a stream.
func Collect[T, R any](stream *Stream[T], collector ResultCollector[T, R]) R {
	return collector.Collect(stream.Collect())
}

// ToSlice creates a collector that returns stream values.
func ToSlice[T any]() ResultCollector[T, []T] {
	return ResultCollectorFunc[T, []T](func(values []T) []T {
		return values
	})
}

// ToSliceCopy creates a collector that returns a copy of stream values.
func ToSliceCopy[T any]() ResultCollector[T, []T] {
	return ResultCollectorFunc[T, []T](func(values []T) []T {
		return slices.Clone(values)
	})
}

// ToMap creates a collector that maps stream values into a map.
func ToMap[T any, K comparable, R any](mapper func(T) (K, R)) ResultCollector[T, map[K]R] {
	return ResultCollectorFunc[T, map[K]R](func(values []T) map[K]R {
		return uarray.ToMap(values, mapper)
	})
}

// ToMultiMap creates a collector that maps stream values into grouped map values.
func ToMultiMap[T any, K comparable, R any](mapper func(T) (K, R)) ResultCollector[T, map[K][]R] {
	return ResultCollectorFunc[T, map[K][]R](func(values []T) map[K][]R {
		return uarray.ToMultiMap(values, mapper)
	})
}

// GroupingBy creates a collector that groups stream values by key.
func GroupingBy[T any, K comparable](key func(T) K) ResultCollector[T, map[K][]T] {
	return ResultCollectorFunc[T, map[K][]T](func(values []T) map[K][]T {
		return uarray.GroupToMapBy(values, key)
	})
}

// CountingBy creates a collector that counts stream values by key.
func CountingBy[T any, K comparable](key func(T) K) ResultCollector[T, map[K]int] {
	return ResultCollectorFunc[T, map[K]int](func(values []T) map[K]int {
		return uarray.CountBy(values, key)
	})
}

// Sorted sorts ordered stream values.
func Sorted[T cmp.Ordered](stream *Stream[T]) *Stream[T] {
	values := slices.Clone(stream.Collect())
	slices.Sort(values)

	return Of(values)
}

// Distinct keeps the first occurrence of each comparable value.
func Distinct[T comparable](stream *Stream[T]) *Stream[T] {
	return Of(uarray.Unique(stream.Collect()))
}

// DistinctBy keeps the first value for each comparable key.
func DistinctBy[T any, K comparable](stream *Stream[T], key func(T) K) *Stream[T] {
	return Of(uarray.Uniq(stream.Collect(), key))
}

// Compact removes zero values from a comparable stream.
func Compact[T comparable](stream *Stream[T]) *Stream[T] {
	return Of(uarray.Compact(stream.Collect()))
}
