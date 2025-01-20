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

func BenchmarkFind(b *testing.B) {
	largeSlice := make([]int, 10000)
	for i := range largeSlice {
		largeSlice[i] = i
	}
	toFind := 9999

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Find(largeSlice, func(v *int) bool {
			return *v == toFind
		})
	}
}

func BenchmarkFindBinary(b *testing.B) {
	largeSlice := make([]int, 10000)
	for i := range largeSlice {
		largeSlice[i] = i
	}
	toFind := 9999

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindBinary(largeSlice, func(v *int) bool {
			return *v == toFind
		})
	}
}

func BenchmarkSortFind(b *testing.B) {
	largeSlice := make([]int, 10000)
	for i := range largeSlice {
		largeSlice[i] = i
	}
	toFind := 5846

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SortFind(largeSlice, func(a, b int) bool {
			return a < b
		}, func(v *int) bool {
			return *v == toFind
		})
	}
}
