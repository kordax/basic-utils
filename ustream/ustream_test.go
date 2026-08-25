/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ustream_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kordax/basic-utils/v4/ustream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStream_NewStream(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	stream := ustream.Of(values)

	require.NotNil(t, stream, "Of returned nil")
	assert.Equal(t, values, stream.Collect(), "Of did not properly initialize with the values")
}

func TestStream_Filter(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	stream := ustream.Of(values)
	filtered := stream.Filter(func(v int) bool {
		return v%2 == 0
	}).Collect()

	expected := []int{2, 4}
	assert.Equal(t, expected, filtered, "Filter function failed")
}

func TestStream_FilterOut(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	stream := ustream.Of(values)
	filteredOut := stream.FilterOut(func(v int) bool {
		return v%2 == 0
	}).Collect()

	expected := []int{1, 3, 5}
	assert.Equal(t, expected, filteredOut, "FilterOut function failed")
}

func TestStream_Map(t *testing.T) {
	values := []int{1, 2, 3}
	stream := ustream.Of(values)
	mapped := stream.Map(func(v int) string {
		return fmt.Sprintf("Num: %d", v)
	}).Collect()

	expected := []string{"Num: 1", "Num: 2", "Num: 3"}
	assert.Equal(t, expected, mapped, "Map function failed")
}

func TestStream_FluentPipeline(t *testing.T) {
	var peeked atomic.Int32

	result := ustream.From(5, 1, 2, 2, 3, 4, 0).
		Filter(func(v int) bool {
			return v >= 2
		}).
		DistinctBy(func(v int) string {
			return fmt.Sprintf("%d", v)
		}).
		Transform(func(v int) int {
			return v * 10
		}).
		Peek(func(v int) {
			peeked.Add(int32(v))
		}).
		Sort(func(a, b int) bool {
			return a < b
		}).
		Skip(1).
		Limit(2).
		Reverse().
		Collect()

	assert.Equal(t, []int{40, 30}, result)
	assert.EqualValues(t, 140, peeked.Load())
}

func TestStream_TakeDropWhileAndMatches(t *testing.T) {
	stream := ustream.From(1, 2, 3, 4, 1)

	assert.Equal(t, []int{1, 2, 3}, stream.TakeWhile(func(v int) bool { return v < 4 }).Collect())
	assert.Equal(t, []int{4, 1}, stream.DropWhile(func(v int) bool { return v < 4 }).Collect())
	assert.True(t, stream.AnyMatch(func(v int) bool { return v == 3 }))
	assert.True(t, stream.AllMatch(func(v int) bool { return v > 0 }))
	assert.True(t, stream.NoneMatch(func(v int) bool { return v < 0 }))
	assert.Equal(t, 2, stream.FindIndex(func(v int) bool { return v == 3 }))
	assert.Equal(t, 5, stream.Count())
	require.NotNil(t, stream.Find(func(v int) bool { return v == 4 }))
	assert.Equal(t, 11, stream.Reduce(0, func(acc int, v int) int { return acc + v }))
}

func TestStream_FlatMapAndCompactFunc(t *testing.T) {
	result := ustream.From("a", "", "b").
		CompactFunc(func(v string) bool {
			return v == ""
		}).
		FlatMap(func(v string) []string {
			return []string{v, strings.ToUpper(v)}
		}).
		Collect()

	assert.Equal(t, []string{"a", "A", "b", "B"}, result)
}

func TestStream_GenericMapFlatMapReduceAndCollectors(t *testing.T) {
	stream := ustream.From(1, 2, 3)

	mapped := stream.Map(func(v int) string {
		return fmt.Sprintf("n-%d", v)
	})
	assert.Equal(t, []string{"n-1", "n-2", "n-3"}, mapped.Collect())

	flatMapped := stream.FlatMap(func(v int) []string {
		return []string{fmt.Sprintf("%d", v), fmt.Sprintf("%d", v*10)}
	})
	assert.Equal(t, []string{"1", "10", "2", "20", "3", "30"}, flatMapped.Collect())

	joined := stream.Reduce("", func(acc string, v int) string {
		return fmt.Sprintf("%s%d", acc, v)
	})
	assert.Equal(t, "123", joined)

	toMap := mapped.CollectWith(ustream.ToMap(func(v string) (int, string) {
		return len(v), v
	}))
	assert.Equal(t, map[int]string{3: "n-3"}, toMap)

	grouped := stream.CollectWith(ustream.GroupingBy(func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	}))
	assert.Equal(t, map[string][]int{"odd": {1, 3}, "even": {2}}, grouped)

	counted := stream.CollectWith(ustream.CountingBy(func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	}))
	assert.Equal(t, map[string]int{"odd": 2, "even": 1}, counted)

	distinct := stream.DistinctBy(func(v int) int {
		return v % 2
	})
	assert.Equal(t, []int{1, 2}, distinct.Collect())
}

func TestStream_ParallelMapPreservesOrderAndBoundsConcurrency(t *testing.T) {
	const size = 32
	const parallelism = 4

	values := make([]int, size)
	expected := make([]string, size)
	for i := range values {
		values[i] = i
		expected[i] = fmt.Sprintf("n-%d", i)
	}

	calls := make([]int, size)
	active := 0
	maxActive := 0
	var mu sync.Mutex

	mapped := ustream.Of(values).ParallelMap(func(value int) string {
		mu.Lock()
		calls[value]++
		active++
		maxActive = max(maxActive, active)
		mu.Unlock()

		time.Sleep(time.Duration((size-value)%5+1) * time.Millisecond)

		mu.Lock()
		active--
		mu.Unlock()

		return fmt.Sprintf("n-%d", value)
	}, parallelism)

	assert.Equal(t, expected, mapped.Collect())
	for value, count := range calls {
		assert.Equalf(t, 1, count, "mapper calls for value %d", value)
	}
	assert.Greater(t, maxActive, 1)
	assert.LessOrEqual(t, maxActive, parallelism)
}

func TestStream_ParallelMapEmpty(t *testing.T) {
	mapped := ustream.Empty[int]().ParallelMap(func(value int) string {
		t.Fatalf("mapper called for empty stream with value %d", value)
		return ""
	}, 4)

	assert.Empty(t, mapped.Collect())
}

func TestStream_ParallelMapRepanicsWithoutEmittingPartialResults(t *testing.T) {
	const panicValue = "parallel map panic"
	const parallelism = 4

	values := make([]int, 128)
	for index := range values {
		values[index] = index
	}

	started := make(chan struct{})
	var processed atomic.Int32
	var emitted atomic.Int32

	mapped := ustream.Of(values).
		ParallelMap(func(value int) int {
			if value == 0 {
				close(started)
				panic(panicValue)
			}

			<-started
			time.Sleep(time.Millisecond)
			processed.Add(1)
			return value
		}, parallelism).
		Peek(func(int) {
			emitted.Add(1)
		})

	require.PanicsWithValue(t, panicValue, func() {
		mapped.Collect()
	})
	assert.Zero(t, emitted.Load(), "partial mapped values reached downstream")
	assert.Less(t, int(processed.Load()), len(values)-1, "mapping continued after the panic")
}

func TestStream_TopLevelDistinctSortedCompact(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, ustream.Sorted(ustream.From(3, 1, 2)).Collect())
	assert.Equal(t, []int{1, 2, 3}, ustream.Distinct(ustream.From(1, 2, 1, 3)).Collect())
	assert.Equal(t, []int{1, 2}, ustream.Compact(ustream.From(0, 1, 0, 2)).Collect())

	type item struct {
		id   int
		name string
	}
	distinct := ustream.DistinctBy(ustream.From(
		item{id: 1, name: "first"},
		item{id: 1, name: "duplicate"},
		item{id: 2, name: "second"},
	), func(v item) int {
		return v.id
	}).Collect()

	assert.Equal(t, []item{{id: 1, name: "first"}, {id: 2, name: "second"}}, distinct)
}

func TestStream_CollectCopyDoesNotExposeBackingArray(t *testing.T) {
	stream := ustream.From(1, 2, 3)
	copied := stream.CollectCopy()
	copied[0] = 100

	assert.Equal(t, []int{1, 2, 3}, stream.Collect())
}

func TestStream_CollectToMap(t *testing.T) {
	stream := ustream.Of([]int{1, 2, 3, 4})

	legacy := stream.CollectToMap(func(value int) (any, any) {
		return value % 2, value
	})
	assert.Equal(t, map[any]any{0: 4, 1: 3}, legacy)

	typed := stream.ToMap(func(value int) (string, int) {
		return fmt.Sprintf("value-%d", value), value
	})
	assert.Equal(t, map[string]int{
		"value-1": 1,
		"value-2": 2,
		"value-3": 3,
		"value-4": 4,
	}, typed)

	grouped := stream.ToMultiMap(func(value int) (int, string) {
		return value % 2, fmt.Sprintf("value-%d", value)
	})
	assert.Equal(t, map[int][]string{
		0: {"value-2", "value-4"},
		1: {"value-1", "value-3"},
	}, grouped)
}

func TestTerminalStream_ParallelExecute(t *testing.T) {
	fn := func(index int, value *int) {}

	stream := ustream.Of([]int{1, 2, 3, 4, 5})
	stream.ToTerminal().ParallelExecute(fn, 4)
}

func TestTerminalStream_ParallelExecuteUsesSnapshot(t *testing.T) {
	stream := ustream.Of([]int{1, 2, 3})
	terminal := stream.ToTerminal()

	terminal.ParallelExecute(func(index int, value *int) {
		*value *= 10
	}, 2)

	assert.Equal(t, []int{10, 20, 30}, terminal.Collect())
	assert.Equal(t, []int{1, 2, 3}, stream.Collect())
}

func TestTerminalStream_ParallelExecuteWithInvalidParallelism(t *testing.T) {
	var processed atomic.Int32
	stream := ustream.Of([]int{1, 2, 3})

	assert.NotPanics(t, func() {
		stream.ToTerminal().ParallelExecute(func(index int, value *int) {
			processed.Add(1)
		}, 0)
	})
	assert.EqualValues(t, 3, processed.Load())
}

func TestTerminalStream_ParallelExecuteRepanicsAndStopsScheduling(t *testing.T) {
	const panicValue = "parallel execute panic"
	const parallelism = 4

	values := make([]int, 128)
	started := make(chan struct{})
	var processed atomic.Int32

	require.PanicsWithValue(t, panicValue, func() {
		ustream.NewTerminalStream(values).ParallelExecute(func(index int, value *int) {
			if index == 0 {
				close(started)
				panic(panicValue)
			}

			<-started
			time.Sleep(time.Millisecond)
			processed.Add(1)
		}, parallelism)
	})
	assert.Less(t, int(processed.Load()), len(values)-1, "execution continued after the panic")
}

func TestTerminalStream_ParallelExecuteContext(t *testing.T) {
	type contextKey struct{}
	key := contextKey{}
	ctx := context.WithValue(context.Background(), key, "expected")
	stream := ustream.NewTerminalStream([]int{1, 2, 3})
	var invalidContext atomic.Bool

	err := stream.ParallelExecuteContext(ctx, func(callbackCtx context.Context, index int, value *int) {
		if callbackCtx.Value(key) != "expected" {
			invalidContext.Store(true)
		}
		*value *= 10
	}, 2)

	require.NoError(t, err)
	assert.False(t, invalidContext.Load())
	assert.Equal(t, []int{10, 20, 30}, stream.Collect())
}

func TestTerminalStream_ParallelExecuteContextEmpty(t *testing.T) {
	var calls atomic.Int32

	err := ustream.NewTerminalStream[int](nil).ParallelExecuteContext(context.Background(), func(context.Context, int, *int) {
		calls.Add(1)
	}, 4)

	require.NoError(t, err)
	assert.Zero(t, calls.Load())
}

func TestTerminalStream_ParallelExecuteContextReturnsExistingCancellation(t *testing.T) {
	cause := errors.New("execution canceled")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(cause)
	var calls atomic.Int32

	err := ustream.NewTerminalStream([]int{1, 2, 3}).ParallelExecuteContext(ctx, func(context.Context, int, *int) {
		calls.Add(1)
	}, 2)

	require.ErrorIs(t, err, cause)
	assert.Zero(t, calls.Load())
}

func TestTerminalStream_ParallelExecuteContextStopsScheduling(t *testing.T) {
	cause := errors.New("stop execution")
	ctx, cancel := context.WithCancelCause(context.Background())
	values := make([]int, 128)
	var calls atomic.Int32

	err := ustream.NewTerminalStream(values).ParallelExecuteContext(ctx, func(callbackCtx context.Context, index int, value *int) {
		calls.Add(1)
		if index == 0 {
			cancel(cause)
			return
		}
		<-callbackCtx.Done()
	}, 4)

	require.ErrorIs(t, err, cause)
	assert.Less(t, int(calls.Load()), len(values), "execution continued scheduling after cancellation")
}

func TestTerminalStream_ParallelExecuteContextReturnsDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	var calls atomic.Int32

	err := ustream.NewTerminalStream([]int{1, 2, 3, 4}).ParallelExecuteContext(ctx, func(callbackCtx context.Context, index int, value *int) {
		calls.Add(1)
		<-callbackCtx.Done()
	}, 2)

	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.LessOrEqual(t, calls.Load(), int32(2))
}

func TestTerminalStream_ParallelExecuteContextPanicTakesPrecedence(t *testing.T) {
	const panicValue = "context execution panic"
	ctx, cancel := context.WithCancel(context.Background())

	require.PanicsWithValue(t, panicValue, func() {
		_ = ustream.NewTerminalStream([]int{1, 2, 3}).ParallelExecuteContext(ctx, func(callbackCtx context.Context, index int, value *int) {
			if index == 0 {
				cancel()
				panic(panicValue)
			}
			<-callbackCtx.Done()
		}, 2)
	})
}

func TestStreamToMultiMap(t *testing.T) {
	grouped := ustream.From("alpha", "atom", "beta").CollectWith(ustream.ToMultiMap(func(value string) (string, string) {
		return string(value[0]), value
	}))

	assert.Equal(t, map[string][]string{
		"a": {"alpha", "atom"},
		"b": {"beta"},
	}, grouped)
}
