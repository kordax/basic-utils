/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap

type void struct{}

var dummy void

// UniqueMultiMap is a generic data structure that implements a multi-value map,
// where each key K can be associated with multiple unique values of type V.
//
// This implementation ensures that there are no duplicate values under the same key,
// as values are stored in a secondary map where the hash code of each value serves
// as the key, and the value itself as the map value. This approach effectively
// prevents storing duplicate values under the same key.
//
// Usage:
// This map is particularly useful in scenarios where it is necessary to maintain a
// collection of unique items (like a set) for each key, but unlike a traditional set,
// it organizes these items by their hash codes to speed up checks for duplicates,
// insertion, and deletion.
//
// !IMPORTANT: This map is not safe for concurrent operations. Use ConcurrentMultiMap for concurrent operations.
// This map does not preserve the order of insertion and cannot store multiple
// identical items under the same key. It is optimized for scenarios where quick
// lookup, insertion, and deletion of items are required, and where each item can be
// uniquely identified by a hash code.
type UniqueMultiMap[K any, V any] struct {
	store map[any]map[any]void
}

func NewUniqueMultiMap[K any, V any]() *UniqueMultiMap[K, V] {
	return &UniqueMultiMap[K, V]{
		store: make(map[any]map[any]void),
	}
}

func (m *UniqueMultiMap[K, V]) Get(key K) ([]V, bool) {
	hashMap, ok := m.store[key]
	if !ok {
		return nil, false
	}

	values := make([]V, 0, len(hashMap))
	for value := range hashMap {
		values = append(values, value.(V))
	}

	return values, true
}

func (m *UniqueMultiMap[K, V]) Set(key K, values ...V) int {
	addedCount := 0
	m.store[key] = make(map[any]void)
	ref := m.store[key]
	for _, value := range values {
		if _, exists := ref[value]; !exists {
			ref[value] = dummy
			addedCount++
		}
	}

	return addedCount
}

func (m *UniqueMultiMap[K, V]) Append(key K, values ...V) int {
	hashMap, exists := m.store[key]
	if !exists {
		hashMap = make(map[any]void)
		m.store[key] = hashMap
	}

	addedCount := 0
	for _, value := range values {
		if _, found := hashMap[value]; !found {
			hashMap[value] = dummy
			addedCount++
		}
	}

	return addedCount
}

func (m *UniqueMultiMap[K, V]) Remove(key K, predicate func(v V) bool) int {
	hashMap, exists := m.store[key]
	if !exists {
		return 0
	}

	removalCount := 0
	for value := range hashMap {
		if predicate(value.(V)) {
			delete(hashMap, value)
			removalCount++
		}
	}

	return removalCount
}

func (m *UniqueMultiMap[K, V]) Clear(key K) bool {
	_, exists := m.store[key]
	if exists {
		delete(m.store, key)
		return true
	}

	return false
}
