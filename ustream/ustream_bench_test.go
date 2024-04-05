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

	"github.com/kordax/basic-utils/uarray"
	"github.com/kordax/basic-utils/ustream"
)

func BenchmarkParallelExecute(b *testing.B) {
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

func BenchmarkParallelExecute_HigherOrder(b *testing.B) {
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
