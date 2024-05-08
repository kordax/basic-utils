/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap_test

import (
	"crypto/sha512"
	"strconv"
)

func generateTestKey(i int) int {
	return i // Just return the benchmark iteration as the test key
}

func generateTestValues(i int) []string {
	return []string{strconv.Itoa(i)}
}

var benchmarkHasher = sha512.New()
