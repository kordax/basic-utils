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
		Contains(val, sampleSlice)
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

// Benchmark for MapKeys function
func BenchmarkMapKeys(b *testing.B) {
	sampleMap := map[int]string{1: "one", 2: "two", 3: "three", 4: "four", 5: "five"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapKeys(sampleMap)
	}
}
