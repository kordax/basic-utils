package umap_test

import (
	"testing"

	"github.com/kordax/basic-utils/umap"
)

func BenchmarkUniqueHashMultiMap_Set(b *testing.B) {
	multiMap := umap.NewUniqueHashMultiMap[int, string](benchmarkHasher)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		key := generateTestKey(i)
		values := generateTestValues(i)
		b.StartTimer()

		multiMap.Set(key, values...)
	}
}

func BenchmarkUniqueHashMultiMap_Get(b *testing.B) {
	multiMap := umap.NewUniqueHashMultiMap[int, string](benchmarkHasher)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		key := generateTestKey(i)
		b.StartTimer()

		_, _ = multiMap.Get(key)
	}
}
