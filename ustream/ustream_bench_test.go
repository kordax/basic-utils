/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ustream_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kordax/basic-utils/v3/uarray"
	"github.com/kordax/basic-utils/v3/ustream"
)

func BenchmarkTerminalStream_ParallelExecute(b *testing.B) {
	parallelisms := uarray.RangeWithStep(1, 40, 4)
	fn := func(index int, value *int) { time.Sleep(time.Nanosecond * 10000) } // Emulates the load

	for _, parallelism := range parallelisms {
		sliceSize := parallelism * 10
		values := make([]int, sliceSize)
		for i := 0; i < sliceSize; i++ {
			values[i] = i + 1
		}
		stream := ustream.NewTerminalStream(values)

		b.Run(fmt.Sprintf("Parallelism-%d", parallelism), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stream.ParallelExecute(fn, parallelism)
			}
		})
	}
}

func BenchmarkTerminalStream_ParallelExecute_HigherOrder(b *testing.B) {
	parallelisms := uarray.RangeWithStep(1, 500, 100)
	fn := func(index int, value *int) { time.Sleep(time.Nanosecond * 10000) } // Emulates the load

	for _, parallelism := range parallelisms {
		sliceSize := parallelism * 10
		values := make([]int, sliceSize)
		for i := 0; i < sliceSize; i++ {
			values[i] = i + 1
		}
		stream := ustream.NewTerminalStream(values)

		b.Run(fmt.Sprintf("Parallelism-%d", parallelism), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stream.ParallelExecute(fn, parallelism)
			}
		})
	}
}

func BenchmarkTerminalStream_ParallelExecuteWithTimeout(b *testing.B) {
	timeout := time.Second
	parallelisms := uarray.RangeWithStep(1, 40, 4)
	fn := func(index int, value int) { time.Sleep(time.Nanosecond * 10000) } // Emulates the load

	for _, parallelism := range parallelisms {
		sliceSize := parallelism * 10
		values := make([]int, sliceSize)
		for i := 0; i < sliceSize; i++ {
			values[i] = i + 1
		}
		stream := ustream.NewTerminalStream(values)

		b.Run(fmt.Sprintf("Parallelism-%d", parallelism), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stream.ParallelExecuteWithTimeout(fn, nil, timeout, parallelism)
			}
		})
	}
}

func BenchmarkTerminalStream_ParallelExecuteWithTimeout_HigherOrder(b *testing.B) {
	timeout := time.Second
	parallelisms := uarray.RangeWithStep(1, 500, 100)
	fn := func(index int, value int) { time.Sleep(time.Nanosecond * 10000) } // Emulates the load

	for _, parallelism := range parallelisms {
		sliceSize := parallelism * 10
		values := make([]int, sliceSize)
		for i := 0; i < sliceSize; i++ {
			values[i] = i + 1
		}
		stream := ustream.NewTerminalStream(values)

		b.Run(fmt.Sprintf("Parallelism-%d", parallelism), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stream.ParallelExecuteWithTimeout(fn, nil, timeout, parallelism)
			}
		})
	}
}

func BenchmarkStream_FluentPipeline(b *testing.B) {
	values := make([]int, 10_000)
	for i := range values {
		values[i] = i
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ustream.Of(values).
			Filter(func(v int) bool {
				return v%2 == 0
			}).
			Transform(func(v int) int {
				return v * 2
			}).
			Skip(10).
			Limit(1000).
			Reverse().
			Reduce(0, func(acc int, v int) int {
				return acc + v
			})
	}
}

func BenchmarkStream_GenericMapCollect(b *testing.B) {
	values := make([]int, 10_000)
	for i := range values {
		values[i] = i
	}
	stream := ustream.Of(values)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapped := ustream.Map(stream, func(v int) int {
			return v * 2
		})
		_ = ustream.Collect(mapped, ustream.ToMap(func(v int) (int, int) {
			return v, v
		}))
	}
}
