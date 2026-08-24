package umap_test

import (
	"testing"

	"git.casinomodule.org/casino27/basic-utils/v4/umap"
)

func BenchmarkUniqueHashMultiMap_Set(b *testing.B) {
	m := umap.NewUniqueHashMultiMap[int, string](benchmarkHasher)

	for i := 0; i < b.N; i++ {
		m.Set(generateTestKey(i), generateTestValues(i)...)
	}
}

func BenchmarkUniqueHashMultiMap_Get(b *testing.B) {
	m := umap.NewUniqueHashMultiMap[int, string](benchmarkHasher)

	b.StopTimer()
	for i := 0; i < b.N; i++ {
		m.Set(generateTestKey(i), generateTestValues(i)...)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, _ = m.Get(generateTestKey(i))
	}
}
