/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ustream

import (
	"cmp"
	"context"
	"iter"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/kordax/basic-utils/v4/uarray"
)

const unknownSizeHint = -1

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

// Stream describes a lazy, ordered and reusable sequence of values.
type Stream[T any] struct {
	seq         iter.Seq[T]
	sizeHint    int
	materialize func() []T
}

type panicRelay struct {
	once  sync.Once
	done  chan struct{}
	value any
}

func newPanicRelay() *panicRelay {
	return &panicRelay{done: make(chan struct{})}
}

func (r *panicRelay) fail(value any) {
	r.once.Do(func() {
		r.value = value
		close(r.done)
	})
}

func (r *panicRelay) stopped() bool {
	select {
	case <-r.done:
		return true
	default:
		return false
	}
}

func (r *panicRelay) repanic() {
	if r.stopped() {
		panic(r.value)
	}
}

func newStream[T any](seq iter.Seq[T], sizeHint int) *Stream[T] {
	return &Stream[T]{seq: seq, sizeHint: sizeHint}
}

func newMaterializingStream[T any](materialize func() []T, sizeHint int) *Stream[T] {
	return &Stream[T]{
		seq: func(yield func(T) bool) {
			for _, value := range materialize() {
				if !yield(value) {
					return
				}
			}
		},
		sizeHint:    sizeHint,
		materialize: materialize,
	}
}

func emptySeq[T any](yield func(T) bool) {}

// Empty creates a new empty Stream.
func Empty[T any]() *Stream[T] {
	return newStream[T](emptySeq[T], 0)
}

// From creates a new Stream from variadic values.
func From[T any](values ...T) *Stream[T] {
	return Of(values)
}

// Of creates a lazy Stream over the given slice.
func Of[T any](values []T) *Stream[T] {
	if values == nil {
		return Empty[T]()
	}

	return newStream(func(yield func(T) bool) {
		for _, value := range values {
			if !yield(value) {
				return
			}
		}
	}, len(values))
}

// FromSeq creates a Stream backed by seq. Its replayability is determined by seq.
func FromSeq[T any](seq iter.Seq[T]) *Stream[T] {
	if seq == nil {
		return Empty[T]()
	}

	return newStream(seq, unknownSizeHint)
}

// Seq returns the lazy sequence represented by the stream.
func (s *Stream[T]) Seq() iter.Seq[T] {
	if s == nil || s.seq == nil {
		return emptySeq[T]
	}

	return s.seq
}

func (s *Stream[T]) allocationHint() int {
	if s == nil || s.sizeHint < 0 {
		return 0
	}

	return s.sizeHint
}

// Generate creates a lazy Stream by calling supplier at most count times per traversal.
func Generate[T any](count int, supplier func(index int) T) *Stream[T] {
	if count <= 0 {
		return Empty[T]()
	}

	return newStream(func(yield func(T) bool) {
		for index := 0; index < count; index++ {
			if !yield(supplier(index)) {
				return
			}
		}
	}, count)
}

// Iterate creates a lazy Stream by repeatedly applying next to the previous value.
func Iterate[T any](seed T, count int, next func(T) T) *Stream[T] {
	if count <= 0 {
		return Empty[T]()
	}

	return newStream(func(yield func(T) bool) {
		value := seed
		for index := 0; index < count; index++ {
			if !yield(value) {
				return
			}
			if index+1 < count {
				value = next(value)
			}
		}
	}, count)
}

// Concat lazily appends values from other after this stream.
func (s *Stream[T]) Concat(other *Stream[T]) *Stream[T] {
	sizeHint := unknownSizeHint
	leftKnown := s == nil || s.sizeHint >= 0
	rightKnown := other == nil || other.sizeHint >= 0
	leftHint := s.allocationHint()
	rightHint := other.allocationHint()
	maxInt := int(^uint(0) >> 1)
	if leftKnown && rightKnown && leftHint <= maxInt-rightHint {
		sizeHint = leftHint + rightHint
	}

	return newStream(func(yield func(T) bool) {
		for value := range s.Seq() {
			if !yield(value) {
				return
			}
		}
		for value := range other.Seq() {
			if !yield(value) {
				return
			}
		}
	}, sizeHint)
}

// Filter lazily keeps values matching predicate.
func (s *Stream[T]) Filter(predicate func(T) bool) *Stream[T] {
	return newStream(func(yield func(T) bool) {
		for value := range s.Seq() {
			if predicate(value) && !yield(value) {
				return
			}
		}
	}, s.allocationHint())
}

// FilterOut lazily removes values matching predicate.
func (s *Stream[T]) FilterOut(predicate func(T) bool) *Stream[T] {
	return s.Filter(func(value T) bool {
		return !predicate(value)
	})
}

// Map lazily maps values into another type.
func (s *Stream[T]) Map[R any](mapper func(T) R) *Stream[R] {
	return newStream(func(yield func(R) bool) {
		for value := range s.Seq() {
			if !yield(mapper(value)) {
				return
			}
		}
	}, s.allocationHint())
}

// Transform lazily maps values without changing their type.
func (s *Stream[T]) Transform(mapper func(T) T) *Stream[T] {
	return s.Map(mapper)
}

// FlatMap lazily maps each value into zero or more values of another type.
func (s *Stream[T]) FlatMap[R any](mapper func(T) []R) *Stream[R] {
	return newStream(func(yield func(R) bool) {
		for value := range s.Seq() {
			for _, mapped := range mapper(value) {
				if !yield(mapped) {
					return
				}
			}
		}
	}, unknownSizeHint)
}

// MapMulti lazily emits zero or more mapped values without allocating a slice per input.
// The mapper must emit synchronously, must not retain emit, and should return when emit returns false.
func (s *Stream[T]) MapMulti[R any](mapper func(T, func(R) bool)) *Stream[R] {
	return newStream(func(yield func(R) bool) {
		stopped := false
		s.Seq()(func(value T) bool {
			mapper(value, func(mapped R) bool {
				if stopped {
					return false
				}
				if !yield(mapped) {
					stopped = true
					return false
				}
				return true
			})
			return !stopped
		})
	}, s.allocationHint())
}

// ParallelMap lazily schedules an ordered parallel mapping barrier.
// A mapper panic stops scheduling, waits for the workers, and is re-panicked by the caller.
func (s *Stream[T]) ParallelMap[R any](mapper func(T) R, parallelism int) *Stream[R] {
	if parallelism <= 1 {
		return s.Map(mapper)
	}

	return newStream(func(yield func(R) bool) {
		values := s.Collect()
		if len(values) == 0 {
			return
		}
		if len(values) == 1 {
			yield(mapper(values[0]))
			return
		}

		workers := min(parallelism, len(values))
		const chunksPerWorker = 8

		result := make([]R, len(values))
		chunkSize := max(1, len(values)/workers/chunksPerWorker)
		relay := newPanicRelay()
		var next atomic.Int64
		var wg sync.WaitGroup

		wg.Add(workers)
		for range workers {
			go func() {
				defer wg.Done()
				defer func() {
					if recovered := recover(); recovered != nil {
						relay.fail(recovered)
					}
				}()

				for {
					if relay.stopped() {
						return
					}

					start := int(next.Add(int64(chunkSize))) - chunkSize
					if start >= len(values) {
						return
					}

					end := min(start+chunkSize, len(values))
					for index := start; index < end; index++ {
						if relay.stopped() {
							return
						}
						result[index] = mapper(values[index])
					}
				}
			}()
		}

		wg.Wait()
		relay.repanic()

		for _, value := range result {
			if !yield(value) {
				return
			}
		}
	}, s.allocationHint())
}

// Peek lazily executes action as values are consumed.
func (s *Stream[T]) Peek(action func(T)) *Stream[T] {
	return newStream(func(yield func(T) bool) {
		for value := range s.Seq() {
			action(value)
			if !yield(value) {
				return
			}
		}
	}, s.allocationHint())
}

// Limit lazily keeps at most n first values.
func (s *Stream[T]) Limit(n int) *Stream[T] {
	if n <= 0 {
		return Empty[T]()
	}

	sizeHint := n
	if hint := s.allocationHint(); hint > 0 {
		sizeHint = min(hint, n)
	}

	return newStream(func(yield func(T) bool) {
		remaining := n
		for value := range s.Seq() {
			if !yield(value) {
				return
			}
			remaining--
			if remaining == 0 {
				return
			}
		}
	}, sizeHint)
}

// Skip lazily removes n first values.
func (s *Stream[T]) Skip(n int) *Stream[T] {
	if n <= 0 {
		return s
	}

	sizeHint := unknownSizeHint
	if s != nil && s.sizeHint >= 0 {
		sizeHint = max(0, s.sizeHint-n)
	}

	return newStream(func(yield func(T) bool) {
		skipped := 0
		for value := range s.Seq() {
			if skipped < n {
				skipped++
				continue
			}
			if !yield(value) {
				return
			}
		}
	}, sizeHint)
}

// TakeWhile lazily keeps leading values while predicate returns true.
func (s *Stream[T]) TakeWhile(predicate func(T) bool) *Stream[T] {
	return newStream(func(yield func(T) bool) {
		for value := range s.Seq() {
			if !predicate(value) || !yield(value) {
				return
			}
		}
	}, s.allocationHint())
}

// DropWhile lazily skips leading values while predicate returns true.
func (s *Stream[T]) DropWhile(predicate func(T) bool) *Stream[T] {
	return newStream(func(yield func(T) bool) {
		dropping := true
		for value := range s.Seq() {
			if dropping && predicate(value) {
				continue
			}
			dropping = false
			if !yield(value) {
				return
			}
		}
	}, s.allocationHint())
}

// Reverse returns a lazy stream that buffers and reverses values when traversed.
func (s *Stream[T]) Reverse() *Stream[T] {
	return newMaterializingStream(func() []T {
		values := s.Collect()
		slices.Reverse(values)
		return values
	}, s.allocationHint())
}

// Sort returns a lazy stream that buffers and sorts values when traversed.
func (s *Stream[T]) Sort(less func(a, b T) bool) *Stream[T] {
	return newMaterializingStream(func() []T {
		values := s.Collect()
		slices.SortFunc(values, func(a, b T) int {
			switch {
			case less(a, b):
				return -1
			case less(b, a):
				return 1
			default:
				return 0
			}
		})
		return values
	}, s.allocationHint())
}

// DistinctBy lazily keeps the first value for each comparable key.
func (s *Stream[T]) DistinctBy[K comparable](key func(T) K) *Stream[T] {
	return newStream(func(yield func(T) bool) {
		seen := make(map[K]struct{}, s.allocationHint())
		for value := range s.Seq() {
			k := key(value)
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			if !yield(value) {
				return
			}
		}
	}, s.allocationHint())
}

// CompactFunc lazily removes values matching empty.
func (s *Stream[T]) CompactFunc(empty func(T) bool) *Stream[T] {
	return s.FilterOut(empty)
}

// ForEach executes action for each value.
func (s *Stream[T]) ForEach(action func(T)) {
	for value := range s.Seq() {
		action(value)
	}
}

// Find returns the first value matching predicate.
func (s *Stream[T]) Find(predicate func(T) bool) *T {
	var result *T
	for value := range s.Seq() {
		if predicate(value) {
			found := value
			result = &found
			break
		}
	}

	return result
}

// FindFirst returns the first value, or nil when the stream is empty.
func (s *Stream[T]) FindFirst() *T {
	for value := range s.Seq() {
		first := value
		return &first
	}

	return nil
}

// Min returns the first minimum value according to less, or nil when the stream is empty.
func (s *Stream[T]) Min(less func(a, b T) bool) *T {
	var result T
	found := false
	for value := range s.Seq() {
		if !found || less(value, result) {
			result = value
			found = true
		}
	}
	if !found {
		return nil
	}

	return &result
}

// Max returns the first maximum value according to less, or nil when the stream is empty.
func (s *Stream[T]) Max(less func(a, b T) bool) *T {
	var result T
	found := false
	for value := range s.Seq() {
		if !found || less(result, value) {
			result = value
			found = true
		}
	}
	if !found {
		return nil
	}

	return &result
}

// FindIndex returns the first index matching predicate, or -1.
func (s *Stream[T]) FindIndex(predicate func(T) bool) int {
	index := 0
	for value := range s.Seq() {
		if predicate(value) {
			return index
		}
		index++
	}

	return -1
}

// AnyMatch checks whether any value matches predicate.
func (s *Stream[T]) AnyMatch(predicate func(T) bool) bool {
	for value := range s.Seq() {
		if predicate(value) {
			return true
		}
	}

	return false
}

// AllMatch checks whether all values match predicate.
func (s *Stream[T]) AllMatch(predicate func(T) bool) bool {
	for value := range s.Seq() {
		if !predicate(value) {
			return false
		}
	}

	return true
}

// NoneMatch checks whether no values match predicate.
func (s *Stream[T]) NoneMatch(predicate func(T) bool) bool {
	return !s.AnyMatch(predicate)
}

// Count returns the number of values produced by the pipeline.
func (s *Stream[T]) Count() int {
	count := 0
	for range s.Seq() {
		count++
	}

	return count
}

// Reduce folds values into any result type.
func (s *Stream[T]) Reduce[R any](initial R, reducer func(acc R, value T) R) R {
	result := initial
	for value := range s.Seq() {
		result = reducer(result, value)
	}

	return result
}

// Collect materializes stream values into a fresh slice.
func (s *Stream[T]) Collect() []T {
	if s != nil && s.materialize != nil {
		return s.materialize()
	}

	result := make([]T, 0, s.allocationHint())
	for value := range s.Seq() {
		result = append(result, value)
	}

	return result
}

// CollectCopy materializes stream values into a fresh slice.
func (s *Stream[T]) CollectCopy() []T {
	return s.Collect()
}

// CollectToMap collects stream values into a map. Duplicate keys are overwritten.
func (s *Stream[T]) CollectToMap(mapper func(T) (any, any)) map[any]any {
	result := make(map[any]any, s.allocationHint())
	for value := range s.Seq() {
		key, mapped := mapper(value)
		result[key] = mapped
	}

	return result
}

// CollectToMultiMap collects stream values into a grouped map.
func (s *Stream[T]) CollectToMultiMap(mapper func(T) (any, any)) map[any][]any {
	result := make(map[any][]any)
	for value := range s.Seq() {
		key, mapped := mapper(value)
		result[key] = append(result[key], mapped)
	}

	return result
}

// CollectWith collects materialized stream values into any result type.
func (s *Stream[T]) CollectWith[R any](collector ResultCollector[T, R]) R {
	return collector.Collect(s.Collect())
}

// ToMap collects stream values into a typed map. Duplicate keys are overwritten.
func (s *Stream[T]) ToMap[K comparable, R any](mapper func(T) (K, R)) map[K]R {
	result := make(map[K]R, s.allocationHint())
	for value := range s.Seq() {
		key, mapped := mapper(value)
		result[key] = mapped
	}

	return result
}

// ToMultiMap collects stream values into a typed grouped map.
func (s *Stream[T]) ToMultiMap[K comparable, R any](mapper func(T) (K, R)) map[K][]R {
	result := make(map[K][]R)
	for value := range s.Seq() {
		key, mapped := mapper(value)
		result[key] = append(result[key], mapped)
	}

	return result
}

// ToTerminal materializes the pipeline into an execution-only snapshot.
func (s *Stream[T]) ToTerminal() *TerminalStream[T] {
	return NewTerminalStream(s.Collect())
}

// TerminalStream represents a materialized stream that can only be collected or executed.
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
// A callback panic stops scheduling, waits for the workers, and is re-panicked by the caller.
// Mutations completed before a panic are not rolled back.
func (s *TerminalStream[T]) ParallelExecute(fn func(int, *T), parallelism int) {
	_ = s.ParallelExecuteContext(context.Background(), func(_ context.Context, index int, value *T) {
		fn(index, value)
	}, parallelism)
}

// ParallelExecuteContext executes fn concurrently and cooperatively stops on context cancellation.
// Running callbacks must observe ctx for prompt cancellation. Completed mutations are not rolled back.
// A callback panic takes precedence over cancellation and is re-panicked by the caller.
func (s *TerminalStream[T]) ParallelExecuteContext(ctx context.Context, fn func(context.Context, int, *T), parallelism int) error {
	if cause := context.Cause(ctx); cause != nil {
		return cause
	}
	if len(s.values) == 0 {
		return nil
	}
	if parallelism <= 0 {
		parallelism = 1
	}
	parallelism = min(parallelism, len(s.values))

	relay := newPanicRelay()
	var wg sync.WaitGroup
	in := make(chan int)

	wg.Add(parallelism)
	for range parallelism {
		go func() {
			defer wg.Done()
			defer func() {
				if recovered := recover(); recovered != nil {
					relay.fail(recovered)
				}
			}()

			for {
				select {
				case <-ctx.Done():
					return
				case <-relay.done:
					return
				case index, ok := <-in:
					if !ok {
						return
					}
					if relay.stopped() || context.Cause(ctx) != nil {
						return
					}
					fn(ctx, index, &s.values[index])
				}
			}
		}()
	}

	stopped := false
	for index := range s.values {
		if stopped {
			break
		}
		select {
		case <-ctx.Done():
			stopped = true
		case <-relay.done:
			stopped = true
		case in <- index:
		}
	}

	close(in)
	wg.Wait()
	relay.repanic()

	return context.Cause(ctx)
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
// Deprecated: use stream.Map(mapper).
func Map[T, R any](stream *Stream[T], mapper func(T) R) *Stream[R] {
	return stream.Map(mapper)
}

// ParallelMap maps stream values concurrently while preserving their original order.
// Deprecated: use stream.ParallelMap(mapper, parallelism).
func ParallelMap[T, R any](stream *Stream[T], mapper func(T) R, parallelism int) *Stream[R] {
	return stream.ParallelMap(mapper, parallelism)
}

// FlatMap maps stream values into another type and flattens the result.
// Deprecated: use stream.FlatMap(mapper).
func FlatMap[T, R any](stream *Stream[T], mapper func(T) []R) *Stream[R] {
	return stream.FlatMap(mapper)
}

// Reduce folds stream values into any result type.
// Deprecated: use stream.Reduce(initial, reducer).
func Reduce[T, R any](stream *Stream[T], initial R, reducer func(acc R, value T) R) R {
	return stream.Reduce(initial, reducer)
}

// Collect applies a typed collector to a stream.
// Deprecated: use stream.CollectWith(collector).
func Collect[T, R any](stream *Stream[T], collector ResultCollector[T, R]) R {
	return stream.CollectWith(collector)
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

// Sorted lazily sorts ordered stream values.
func Sorted[T cmp.Ordered](stream *Stream[T]) *Stream[T] {
	return stream.Sort(func(a, b T) bool {
		return a < b
	})
}

// Distinct lazily keeps the first occurrence of each comparable value.
func Distinct[T comparable](stream *Stream[T]) *Stream[T] {
	return stream.DistinctBy(func(value T) T {
		return value
	})
}

// DistinctBy lazily keeps the first value for each comparable key.
// Deprecated: use stream.DistinctBy(key).
func DistinctBy[T any, K comparable](stream *Stream[T], key func(T) K) *Stream[T] {
	return stream.DistinctBy(key)
}

// Compact lazily removes zero values from a comparable stream.
func Compact[T comparable](stream *Stream[T]) *Stream[T] {
	var zero T
	return stream.CompactFunc(func(value T) bool {
		return value == zero
	})
}
