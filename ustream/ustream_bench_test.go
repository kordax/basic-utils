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

	"git.casinomodule.org/casino27/basic-utils/v4/uarray"
	"git.casinomodule.org/casino27/basic-utils/v4/ustream"
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

var benchmarkMappedValues []int

func BenchmarkStream_MapVsParallelMap(b *testing.B) {
	cheapValues := make([]int, 10_000)
	for i := range cheapValues {
		cheapValues[i] = i
	}
	cheapStream := ustream.Of(cheapValues)
	cheapMapper := func(value int) int {
		return value * 2
	}

	b.Run("cheap/Map", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkMappedValues = ustream.Map(cheapStream, cheapMapper).Collect()
		}
	})
	b.Run("cheap/ParallelMap-4", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkMappedValues = ustream.ParallelMap(cheapStream, cheapMapper, 4).Collect()
		}
	})

	cpuValues := make([]int, 1_024)
	for i := range cpuValues {
		cpuValues[i] = i
	}
	cpuStream := ustream.Of(cpuValues)
	cpuMapper := func(value int) int {
		for range 2_048 {
			value = value*1_664_525 + 1_013_904_223
		}
		return value
	}

	b.Run("cpu/Map", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkMappedValues = ustream.Map(cpuStream, cpuMapper).Collect()
		}
	})
	b.Run("cpu/ParallelMap-4", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkMappedValues = ustream.ParallelMap(cpuStream, cpuMapper, 4).Collect()
		}
	})
}
