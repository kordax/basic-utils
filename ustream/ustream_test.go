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
	"sync/atomic"
	"testing"
	"time"

	"github.com/kordax/basic-utils/v2/ustream"
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
