/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap_test

import (
	"testing"

	"github.com/kordax/basic-utils/umap"
)

func BenchmarkUniqueMultiMap_Set(b *testing.B) {
	multiMap := umap.NewUniqueMultiMap[int, string]()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		key := generateTestKey(i)
		values := generateTestValues(i)
		b.StartTimer()

		multiMap.Set(key, values...)
	}
}

func BenchmarkUniqueMultiMap_Get(b *testing.B) {
	multiMap := umap.NewUniqueMultiMap[int, string]()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		key := generateTestKey(i)
		b.StartTimer()

		_, _ = multiMap.Get(key)
	}
}
