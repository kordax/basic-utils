package ustream_test

import (
	"iter"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kordax/basic-utils/v4/ustream"
)

func TestStream_LazyPipelineIsReplayable(t *testing.T) {
	reads := 0
	maps := 0
	peeks := 0
	source := iter.Seq[int](func(yield func(int) bool) {
		for _, value := range []int{1, 2, 3, 4} {
			reads++
			if !yield(value) {
				return
			}
		}
	})

	pipeline := ustream.FromSeq(source).
		Map(func(value int) int {
			maps++
			return value * 2
		}).
		Filter(func(value int) bool { return value%4 == 0 }).
		Peek(func(int) { peeks++ })

	assert.Zero(t, reads)
	assert.Zero(t, maps)
	assert.Zero(t, peeks)

	assert.Equal(t, []int{4, 8}, pipeline.Collect())
	assert.Equal(t, 4, reads)
	assert.Equal(t, 4, maps)
	assert.Equal(t, 2, peeks)

	assert.Equal(t, []int{4, 8}, pipeline.Collect())
	assert.Equal(t, 8, reads)
	assert.Equal(t, 8, maps)
	assert.Equal(t, 4, peeks)
}

func TestStream_LimitStopsSourceWithoutReadingAhead(t *testing.T) {
	reads := 0
	maps := 0
	source := iter.Seq[int](func(yield func(int) bool) {
		for value := 1; value <= 10; value++ {
			reads++
			if !yield(value) {
				return
			}
		}
	})

	actual := ustream.FromSeq(source).
		Map(func(value int) int {
			maps++
			return value * 10
		}).
		Limit(3).
		Collect()

	assert.Equal(t, []int{10, 20, 30}, actual)
	assert.Equal(t, 3, reads)
	assert.Equal(t, 3, maps)
}

func TestStream_TerminalsShortCircuitSource(t *testing.T) {
	newSource := func(reads *int) *ustream.Stream[int] {
		return ustream.FromSeq(iter.Seq[int](func(yield func(int) bool) {
			for value := 1; value <= 10; value++ {
				*reads++
				if !yield(value) {
					return
				}
			}
		}))
	}

	reads := 0
	assert.True(t, newSource(&reads).AnyMatch(func(value int) bool { return value == 4 }))
	assert.Equal(t, 4, reads)

	reads = 0
	found := newSource(&reads).Find(func(value int) bool { return value == 3 })
	require.NotNil(t, found)
	assert.Equal(t, 3, *found)
	assert.Equal(t, 3, reads)

	reads = 0
	assert.False(t, newSource(&reads).AllMatch(func(value int) bool { return value < 3 }))
	assert.Equal(t, 3, reads)
}

func TestStream_FlatMapPropagatesShortCircuit(t *testing.T) {
	reads := 0
	maps := 0
	source := ustream.FromSeq(iter.Seq[int](func(yield func(int) bool) {
		for value := 1; value <= 10; value++ {
			reads++
			if !yield(value) {
				return
			}
		}
	}))

	actual := source.
		FlatMap(func(value int) []int {
			maps++
			return []int{value, -value}
		}).
		Limit(3).
		Collect()

	assert.Equal(t, []int{1, -1, 2}, actual)
	assert.Equal(t, 2, reads)
	assert.Equal(t, 2, maps)
}

func TestStream_GenerateAndIterateAreLazy(t *testing.T) {
	generated := 0
	pipeline := ustream.Generate(10, func(int) int {
		generated++
		return generated
	}).Limit(3)

	assert.Zero(t, generated)
	assert.Equal(t, []int{1, 2, 3}, pipeline.Collect())
	assert.Equal(t, 3, generated)

	unusedSupplierCalls := 0
	emptyGenerated := ustream.Generate(0, func(int) int {
		unusedSupplierCalls++
		return 1
	})
	assert.Empty(t, emptyGenerated.Collect())
	assert.Zero(t, unusedSupplierCalls)

	nextCalls := 0
	iterated := ustream.Iterate(1, 4, func(value int) int {
		nextCalls++
		return value + 2
	})
	assert.Zero(t, nextCalls)
	assert.Equal(t, []int{1, 3, 5, 7}, iterated.Collect())
	assert.Equal(t, 3, nextCalls)

	unusedNextCalls := 0
	emptyIterated := ustream.Iterate(1, -1, func(value int) int {
		unusedNextCalls++
		return value + 1
	})
	assert.Empty(t, emptyIterated.Collect())
	assert.Zero(t, unusedNextCalls)
}

func TestStream_BufferingStagesStayLazy(t *testing.T) {
	t.Run("sort", func(t *testing.T) {
		reads := 0
		stream := countedStream([]int{3, 1, 2}, &reads).
			Sort(func(left, right int) bool { return left < right })
		assert.Zero(t, reads)
		assert.Equal(t, []int{1, 2, 3}, stream.Collect())
		assert.Equal(t, 3, reads)
	})

	t.Run("reverse", func(t *testing.T) {
		reads := 0
		stream := countedStream([]int{1, 2, 3}, &reads).Reverse()
		assert.Zero(t, reads)
		assert.Equal(t, []int{3, 2, 1}, stream.Collect())
		assert.Equal(t, 3, reads)
	})

	t.Run("parallel map", func(t *testing.T) {
		reads := 0
		var maps atomic.Int32
		stream := countedStream([]int{1, 2, 3}, &reads).
			ParallelMap(func(value int) int {
				maps.Add(1)
				return value * 10
			}, 2)
		assert.Zero(t, reads)
		assert.Zero(t, maps.Load())
		assert.Equal(t, []int{10, 20, 30}, stream.Collect())
		assert.Equal(t, 3, reads)
		assert.EqualValues(t, 3, maps.Load())
	})
}

func TestStream_ParallelMapPreservesOrderAndUsesWorkers(t *testing.T) {
	values := make([]int, 32)
	expected := make([]int, len(values))
	for index := range values {
		values[index] = index
		expected[index] = index * index
	}

	var active atomic.Int32
	var maximum atomic.Int32
	actual := ustream.Of(values).ParallelMap(func(value int) int {
		current := active.Add(1)
		for {
			observed := maximum.Load()
			if current <= observed || maximum.CompareAndSwap(observed, current) {
				break
			}
		}
		time.Sleep(time.Millisecond)
		active.Add(-1)
		return value * value
	}, 4).Collect()

	assert.Equal(t, expected, actual)
	assert.Greater(t, maximum.Load(), int32(1))
}

func TestStream_ParallelMapSequentialFallbackIsLazy(t *testing.T) {
	calls := 0
	stream := ustream.Of([]int{1, 2, 3}).ParallelMap(func(value int) int {
		calls++
		return value + 1
	}, 0)

	assert.Zero(t, calls)
	assert.Equal(t, []int{2, 3, 4}, stream.Collect())
	assert.Equal(t, 3, calls)
}

func TestStream_SeqSupportsConsumerShortCircuit(t *testing.T) {
	reads := 0
	stream := countedStream([]int{1, 2, 3, 4}, &reads)
	actual := make([]int, 0, 2)

	stream.Seq()(func(value int) bool {
		actual = append(actual, value)
		return len(actual) < 2
	})

	assert.Equal(t, []int{1, 2}, actual)
	assert.Equal(t, 2, reads)
}

func TestStream_EmptyAndBoundaryOperations(t *testing.T) {
	var nilStream *ustream.Stream[int]
	var zeroStream ustream.Stream[int]

	assert.Empty(t, nilStream.Collect())
	assert.Empty(t, zeroStream.Collect())
	assert.Empty(t, ustream.Empty[int]().Collect())
	assert.Empty(t, ustream.Of[int](nil).Collect())
	assert.Empty(t, ustream.FromSeq[int](nil).Collect())
	assert.Empty(t, ustream.Of([]int{1, 2, 3}).Limit(0).Collect())
	assert.Equal(t, []int{1, 2, 3}, ustream.Of([]int{1, 2, 3}).Limit(100).Collect())
	assert.Equal(t, []int{1, 2, 3}, ustream.Of([]int{1, 2, 3}).Skip(-1).Collect())
	assert.Empty(t, ustream.Of([]int{1, 2, 3}).Skip(100).Collect())
	assert.Empty(t, ustream.Of([]int{1, 2, 3}).TakeWhile(func(int) bool { return false }).Collect())
	assert.Equal(t, []int{1, 2, 3}, ustream.Of([]int{1, 2, 3}).TakeWhile(func(int) bool { return true }).Collect())
	assert.Equal(t, []int{1, 2, 3}, ustream.Of([]int{1, 2, 3}).DropWhile(func(int) bool { return false }).Collect())
	assert.Empty(t, ustream.Of([]int{1, 2, 3}).DropWhile(func(int) bool { return true }).Collect())
}

func TestStream_TerminalEdgeCases(t *testing.T) {
	stream := ustream.Of([]int{1, 2, 3})
	sum := 0
	stream.ForEach(func(value int) { sum += value })
	assert.Equal(t, 6, sum)
	assert.Nil(t, stream.Find(func(value int) bool { return value > 10 }))
	assert.Equal(t, -1, stream.FindIndex(func(value int) bool { return value > 10 }))
	assert.False(t, stream.AnyMatch(func(value int) bool { return value > 10 }))
	assert.True(t, ustream.Empty[int]().AllMatch(func(int) bool { return false }))
	assert.True(t, ustream.Empty[int]().NoneMatch(func(int) bool { return true }))

	multiMap := stream.CollectToMultiMap(func(value int) (any, any) {
		return value % 2, value * 10
	})
	assert.Equal(t, map[any][]any{0: {20}, 1: {10, 30}}, multiMap)
}

func TestStream_DeprecatedCompatibilityFunctions(t *testing.T) {
	stream := ustream.Of([]int{1, 2, 3})

	assert.Equal(t, []string{"1", "2", "3"}, ustream.Map(stream, func(value int) string {
		return string(rune('0' + value))
	}).Collect())
	assert.Equal(t, []int{2, 4, 6}, ustream.ParallelMap(stream, func(value int) int { return value * 2 }, 2).Collect())
	assert.Equal(t, []int{1, -1, 2, -2, 3, -3}, ustream.FlatMap(stream, func(value int) []int {
		return []int{value, -value}
	}).Collect())
	assert.Equal(t, 6, ustream.Reduce(stream, 0, func(result, value int) int { return result + value }))
	assert.Equal(t, []int{1, 2, 3}, ustream.Collect(stream, ustream.ToSlice[int]()))
	assert.Equal(t, []int{1, 2, 3}, ustream.Collect(stream, ustream.ToSliceCopy[int]()))
}

func TestTerminalStream_CollectionOperations(t *testing.T) {
	empty := ustream.NewTerminalStream[int](nil)
	assert.Empty(t, empty.Collect())

	terminal := ustream.NewTerminalStream([]int{1, 2, 3})
	copyValues := terminal.CollectCopy()
	copyValues[0] = 100
	assert.Equal(t, []int{1, 2, 3}, terminal.Collect())

	actual := terminal.CollectToMap(func(value int) (any, int) {
		return value % 2, value * 10
	})
	assert.Equal(t, map[any][]int{0: {20}, 1: {10, 30}}, actual)
}

func TestStream_MapMultiIsLazyOrderedAndReplayable(t *testing.T) {
	mapperCalls := 0
	stream := ustream.Of([]int{1, 2, 3, 4}).MapMulti(func(value int, emit func(int64) bool) {
		mapperCalls++
		if value%2 == 0 {
			return
		}
		if !emit(int64(value)) {
			return
		}
		emit(-int64(value))
	})

	assert.Zero(t, mapperCalls)
	assert.Equal(t, []int64{1, -1, 3, -3}, stream.Collect())
	assert.Equal(t, 4, mapperCalls)
	assert.Equal(t, []int64{1, -1, 3, -3}, stream.Collect())
	assert.Equal(t, 8, mapperCalls)
}

func TestStream_MapMultiPropagatesShortCircuit(t *testing.T) {
	reads := 0
	mapperCalls := 0
	emitCalls := 0
	stream := countedStream([]int{1, 2, 3, 4}, &reads).
		MapMulti(func(value int, emit func(int) bool) {
			mapperCalls++
			emitCalls++
			if !emit(value) {
				return
			}
			emitCalls++
			emit(-value)
		}).
		Limit(3)

	assert.Equal(t, []int{1, -1, 2}, stream.Collect())
	assert.Equal(t, 2, reads)
	assert.Equal(t, 2, mapperCalls)
	assert.Equal(t, 3, emitCalls)
}

func TestStream_MapMultiSuppressesEmissionsAfterStop(t *testing.T) {
	mapperCalls := 0
	emitCalls := 0
	actual := ustream.Of([]int{1, 2}).
		MapMulti(func(value int, emit func(int) bool) {
			mapperCalls++
			emitCalls++
			emit(value)
			emitCalls++
			emit(-value)
		}).
		Limit(1).
		Collect()

	assert.Equal(t, []int{1}, actual)
	assert.Equal(t, 1, mapperCalls)
	assert.Equal(t, 2, emitCalls)
}

func TestStream_ConcatIsLazyOrderedAndShortCircuits(t *testing.T) {
	leftReads := 0
	rightReads := 0
	stream := countedStream([]int{1, 2}, &leftReads).
		Concat(countedStream([]int{3, 4}, &rightReads))

	assert.Zero(t, leftReads)
	assert.Zero(t, rightReads)
	assert.Equal(t, []int{1, 2, 3, 4}, stream.Collect())
	assert.Equal(t, 2, leftReads)
	assert.Equal(t, 2, rightReads)

	leftReads = 0
	rightReads = 0
	assert.Equal(t, []int{1}, stream.Limit(1).Collect())
	assert.Equal(t, 1, leftReads)
	assert.Zero(t, rightReads)
}

func TestStream_ConcatTreatsNilStreamsAsEmpty(t *testing.T) {
	var nilStream *ustream.Stream[int]

	assert.Equal(t, []int{1, 2}, nilStream.Concat(ustream.Of([]int{1, 2})).Collect())
	assert.Equal(t, []int{1, 2}, ustream.Of([]int{1, 2}).Concat(nil).Collect())
}

func TestStream_FindFirstShortCircuitsAndHandlesEmpty(t *testing.T) {
	reads := 0
	first := countedStream([]int{7, 8, 9}, &reads).FindFirst()

	require.NotNil(t, first)
	assert.Equal(t, 7, *first)
	assert.Equal(t, 1, reads)
	assert.Nil(t, ustream.Empty[int]().FindFirst())
}

func TestStream_MinMaxUseFirstValueOnTies(t *testing.T) {
	type scored struct {
		name  string
		score int
	}
	values := []scored{
		{name: "middle", score: 2},
		{name: "minimum-first", score: 1},
		{name: "minimum-second", score: 1},
		{name: "maximum-first", score: 3},
		{name: "maximum-second", score: 3},
	}
	less := func(left, right scored) bool { return left.score < right.score }
	stream := ustream.Of(values)

	minimum := stream.Min(less)
	maximum := stream.Max(less)
	require.NotNil(t, minimum)
	require.NotNil(t, maximum)
	assert.Equal(t, "minimum-first", minimum.name)
	assert.Equal(t, "maximum-first", maximum.name)
	assert.Nil(t, ustream.Empty[scored]().Min(less))
	assert.Nil(t, ustream.Empty[scored]().Max(less))
}

func TestStream_BufferingCollectIsLazyReplayableAndFresh(t *testing.T) {
	tests := []struct {
		name     string
		pipeline func(*ustream.Stream[int]) *ustream.Stream[int]
		expected []int
	}{
		{name: "reverse", pipeline: func(stream *ustream.Stream[int]) *ustream.Stream[int] {
			return stream.Reverse()
		}, expected: []int{2, 1, 3}},
		{name: "sort", pipeline: func(stream *ustream.Stream[int]) *ustream.Stream[int] {
			return stream.Sort(func(left, right int) bool { return left < right })
		}, expected: []int{1, 2, 3}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reads := 0
			stream := test.pipeline(countedStream([]int{3, 1, 2}, &reads))
			assert.Zero(t, reads)

			first := stream.Collect()
			assert.Equal(t, test.expected, first)
			assert.Equal(t, 3, reads)
			first[0] = 100

			assert.Equal(t, test.expected, stream.Collect())
			assert.Equal(t, 6, reads)
		})
	}
}

func TestStream_BufferingStageKeepsLazyDownstreamSemantics(t *testing.T) {
	reads := 0
	actual := countedStream([]int{1, 2, 3}, &reads).
		Reverse().
		Map(func(value int) int { return value * 10 }).
		Limit(2).
		Collect()

	assert.Equal(t, []int{30, 20}, actual)
	assert.Equal(t, 3, reads)
}

func countedStream(values []int, reads *int) *ustream.Stream[int] {
	return ustream.FromSeq(iter.Seq[int](func(yield func(int) bool) {
		for _, value := range values {
			*reads++
			if !yield(value) {
				return
			}
		}
	}))
}
