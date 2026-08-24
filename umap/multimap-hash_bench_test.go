/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap_test

import (
	"testing"

	"git.casinomodule.org/casino27/basic-utils/v4/umap"
)

func BenchmarkHashMultiMap_Set(b *testing.B) {
	m := umap.NewHashMultiMap[int, string]()

	for i := 0; i < b.N; i++ {
		m.Set(generateTestKey(i), generateTestValues(i)...)
	}
}

func BenchmarkHashMultiMap_Get(b *testing.B) {
	m := umap.NewHashMultiMap[int, string]()

	b.StopTimer()
	for i := 0; i < b.N; i++ {
		m.Set(generateTestKey(i), generateTestValues(i)...)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, _ = m.Get(generateTestKey(i))
	}
}

func BenchmarkHashMultiMap_Remove(b *testing.B) {
	values := make([]int, 10000)
	for i := range values {
		values[i] = i
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := umap.NewHashMultiMap[int, int]()
		m.Set(1, values...)
		_ = m.Remove(1, func(v int) bool {
			return v%3 == 0
		})
	}
}
