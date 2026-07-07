/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func BenchmarkComputeHash_MapField(b *testing.B) {
	metadata := make(map[string]string, 100)
	for i := 0; i < 100; i++ {
		metadata[fmt.Sprintf("key-%03d", i)] = fmt.Sprintf("value-%03d", i)
	}

	value := hashTestStruct{
		ID:       1,
		Name:     "John Doe",
		IsActive: true,
		Ratings:  []int{5, 4, 3, 2, 1},
		Metadata: metadata,
		Profile: &profile{
			Title:   "Manager",
			Company: "Tech Inc",
		},
	}

	h := sha256.New()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = computeHash(h, value)
	}
}
