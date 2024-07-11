/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap

import (
	"slices"
)

// HashMultiMap is a generic data structure that implements a multi-value map,
// where each key K can be associated with multiple values of type V.
//
// This implementation allows duplicate values under the same key,
// as values are stored in a slice. This approach supports storing multiple
// identical items under the same key.
//
// Usage:
// This map is particularly useful in scenarios where it is necessary to maintain a
// collection of items (including duplicates) for each key. It is optimized for
// scenarios where quick lookup, insertion, and deletion of items are required.
//
// !IMPORTANT: This map is not safe for concurrent operations. Use ConcurrentMultiMap for concurrent operations.
// This map does not preserve the order of insertion and allows storing multiple
// identical items under the same key.
type HashMultiMap[K any, V any] struct {
	store map[any][]V
}

func NewHashMultiMap[K any, V any]() *HashMultiMap[K, V] {
	return &HashMultiMap[K, V]{
		store: make(map[any][]V),
	}
}

func (m *HashMultiMap[K, V]) Get(key K) ([]V, bool) {
	values, ok := m.store[key]
	return values, ok
}

func (m *HashMultiMap[K, V]) Set(key K, values ...V) int {
	m.store[key] = values
	return len(values)
}

func (m *HashMultiMap[K, V]) Append(key K, values ...V) int {
	m.store[key] = append(m.store[key], values...)
	return len(values)
}

func (m *HashMultiMap[K, V]) Remove(key K, predicate func(v V) bool) int {
	values, exists := m.store[key]
	if !exists {
		return 0
	}

	removalCount := 0
	indexesToRemove := make([]int, 0)
	for ind, value := range values {
		if predicate(value) {
			indexesToRemove = append(indexesToRemove, ind)
			removalCount++
		}
	}
	m.store[key] = withoutIndexes(values, indexesToRemove)

	return removalCount
}

func (m *HashMultiMap[K, V]) Clear(key K) bool {
	_, exists := m.store[key]
	if exists {
		delete(m.store, key)
		return true
	}

	return false
}

func withoutIndexes[T any](src []T, indexes []int) []T {
	indexMap := make(map[int]struct{})
	for _, index := range indexes {
		indexMap[index] = dummy
	}

	uniqueIndexes := make([]int, 0, len(indexMap))
	for index := range indexMap {
		uniqueIndexes = append(uniqueIndexes, index)
	}

	slices.Sort(uniqueIndexes)
	slices.Reverse(uniqueIndexes)

	for _, index := range uniqueIndexes {
		if index < len(src) {
			src = append(src[:index], src[index+1:]...)
		}
	}

	return src
}
