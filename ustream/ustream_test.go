/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ustream_test

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"git.casinomodule.org/casino27/basic-utils/v4/ustream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummy = struct{}

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
	mapped := stream.Map(func(v int) any {
		return fmt.Sprintf("Num: %d", v)
	}).Collect()

	expected := []interface{}{"Num: 1", "Num: 2", "Num: 3"}
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

	mapped := ustream.Map(stream, func(v int) string {
		return fmt.Sprintf("n-%d", v)
	})
	assert.Equal(t, []string{"n-1", "n-2", "n-3"}, mapped.Collect())

	flatMapped := ustream.FlatMap(stream, func(v int) []string {
		return []string{fmt.Sprintf("%d", v), fmt.Sprintf("%d", v*10)}
	})
	assert.Equal(t, []string{"1", "10", "2", "20", "3", "30"}, flatMapped.Collect())

	sum := ustream.Reduce(stream, 0, func(acc int, v int) int {
		return acc + v
	})
	assert.Equal(t, 6, sum)

	toMap := ustream.Collect(mapped, ustream.ToMap(func(v string) (int, string) {
		return len(v), v
	}))
	assert.Equal(t, map[int]string{3: "n-3"}, toMap)

	grouped := ustream.Collect(stream, ustream.GroupingBy(func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	}))
	assert.Equal(t, map[string][]int{"odd": {1, 3}, "even": {2}}, grouped)

	counted := ustream.Collect(stream, ustream.CountingBy(func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	}))
	assert.Equal(t, map[string]int{"odd": 2, "even": 1}, counted)
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

	mapped := ustream.ParallelMap(ustream.Of(values), func(value int) string {
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
	mapped := ustream.ParallelMap(ustream.Empty[int](), func(value int) string {
		t.Fatalf("mapper called for empty stream with value %d", value)
		return ""
	}, 4)

	assert.Empty(t, mapped.Collect())
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
	rand.New(rand.NewSource(time.Now().UnixNano()))

	var values []int
	for i := 0; i < 1000; i++ {
		values = append(values, rand.Intn(1000))
	}

	stream := ustream.Of(values)

	// Perform operations on the stream
	// Since Map returns a TerminalStream, we perform all transformations before mapping
	filteredStream := stream.Filter(func(v int) bool {
		return v%2 == 0 // Keep only even numbers
	}).FilterOut(func(v int) bool {
		return strings.Contains(fmt.Sprintf("%d", v), "100") // Remove numbers containing '100'
	})

	resultStream := filteredStream.Map(func(v int) any {
		return fmt.Sprintf("Even-%d", v) // Convert to string with a prefix
	})

	collectedMap := resultStream.CollectToMap(func(v any) (any, any) {
		return len((v).(string)), v
	})

	assert.NotEmpty(t, collectedMap, "CollectToMap should produce a non-empty map")
	for key, valueSlice := range collectedMap {
		assert.IsType(t, 0, key, "Keys in the map should be of int type")
		assert.NotEmpty(t, valueSlice, "Value slices in the map should be non-empty")
		for _, value := range valueSlice {
			assert.IsType(t, "", value, "Values in the map should be of string type")
		}
	}
}

func TestTerminalStream_ParallelExecute(t *testing.T) {
	fn := func(index int, value *int) {}

	stream := ustream.Of([]int{1, 2, 3, 4, 5})
	stream.ToTerminal().ParallelExecute(fn, 4)
}

func TestTerminalStream_ParallelExecuteMutatesOriginalValues(t *testing.T) {
	stream := ustream.Of([]int{1, 2, 3})

	stream.ToTerminal().ParallelExecute(func(index int, value *int) {
		*value *= 10
	}, 2)

	assert.Equal(t, []int{10, 20, 30}, stream.Collect())
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

func TestTerminalStream_ParallelExecuteWithTimeout(t *testing.T) {
	data := make([]dummy, 100)
	stream := ustream.NewTerminalStream(data)

	var counter atomic.Int32

	mockFn := func(index int, item dummy) {
		time.Sleep(1 * time.Millisecond) // Simulate some processing time
		counter.Add(1)
	}
	cancel := func(i int, d dummy) {
		assert.Fail(t, fmt.Sprintf("cancel for timeout called on item index %d", i))
	}

	stream.ParallelExecuteWithTimeout(mockFn, cancel, 3*time.Second, 10)

	require.EqualValues(t, len(data), counter.Load(), "Not all items were processed as expected")
}

func TestTerminalStream_ParallelExecuteWithTimeoutUsesPerItemTimeout(t *testing.T) {
	values := []int{1, 2, 3}
	stream := ustream.NewTerminalStream(values)

	var processed atomic.Int32
	var canceled atomic.Int32

	stream.ParallelExecuteWithTimeout(func(index int, item int) {
		if item == 1 {
			time.Sleep(30 * time.Millisecond)
			return
		}
		processed.Add(1)
	}, func(index int, item int) {
		canceled.Add(1)
	}, 5*time.Millisecond, 1)

	require.EqualValues(t, 2, processed.Load(), "fast items should still run after one item times out")
	require.EqualValues(t, 1, canceled.Load(), "only the slow item should be canceled")
}

func TestTerminalStream_ParallelExecuteWithTimeoutInvalidParallelism(t *testing.T) {
	stream := ustream.NewTerminalStream([]int{1, 2, 3})
	var processed atomic.Int32

	assert.NotPanics(t, func() {
		stream.ParallelExecuteWithTimeout(func(index int, item int) {
			processed.Add(1)
		}, func(index int, item int) {
			assert.Fail(t, "unexpected cancellation")
		}, time.Second, 0)
	})
	assert.EqualValues(t, 3, processed.Load())
}

func TestTerminalStream_ParallelExecuteWithTimeout_Timeout(t *testing.T) {
	data := make([]dummy, 100)
	stream := ustream.NewTerminalStream(data)

	var counter atomic.Int32

	mockFn := func(index int, item dummy) {
		time.Sleep(10 * time.Millisecond) // Simulate some processing time
	}
	cancel := func(i int, d dummy) {
		counter.Add(1)
	}

	stream.ParallelExecuteWithTimeout(mockFn, cancel, 0, 10)

	require.EqualValues(t, len(data), counter.Load(), "Not all items were processed as expected")
}
