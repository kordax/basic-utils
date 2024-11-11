/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uarray

import (
	"testing"
)

// Benchmark for Contains function
func BenchmarkContains(b *testing.B) {
	sampleSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	val := 5

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Contains(sampleSlice, val)
	}
}

// Benchmark for ContainsAny function
func BenchmarkContainsAny(b *testing.B) {
	sampleSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	values := []int{15, 20, 5, 30}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ContainsAny(sampleSlice, values...)
	}
}

// Benchmark for EqualValues function
func BenchmarkEqualValues(b *testing.B) {
	slice1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	slice2 := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EqualValues(slice1, slice2)
	}
}

func BenchmarkCopyWithoutIndexes(b *testing.B) {
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	indexes := []int{2, 5, 8}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CopyWithoutIndexes(src, indexes)
	}
}

func BenchmarkMapAndGroupToMapBy(b *testing.B) {
	sampleSlice := []string{"apple", "banana", "cherry", "avocado", "blueberry", "grape", "melon"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MapAndGroupToMapBy(sampleSlice, func(v *string) (int, *string) {
			return len(*v), v
		})
	}
}
