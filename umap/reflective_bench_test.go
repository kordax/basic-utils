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

func BenchmarkReflectiveMultiMap_Set(b *testing.B) {
	multiMap := umap.NewReflectiveMultiMap[int, string]()

	for i := 0; i < b.N; i++ {
		key := generateTestKey(i)
		values := generateTestValues(i)

		multiMap.Set(key, values...)
	}
}

func BenchmarkReflectiveMultiMap_Get(b *testing.B) {
	multiMap := umap.NewReflectiveMultiMap[int, string]()

	for i := 0; i < b.N; i++ {
		key := generateTestKey(i)

		_, _ = multiMap.Get(key)
	}
}
