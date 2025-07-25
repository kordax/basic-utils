/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap_test

import (
	"testing"

	"github.com/kordax/basic-utils/v2/umap"
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
