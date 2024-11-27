/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap

import "iter"

// ReflectiveMultiMap is a concurrent-safe data structure that implements a multi-value map.
// It associates keys of type K with multiple values of type V. and values under the same key
// // can have duplicates if they produce the same hash code. The structure employs a
// hash.Hash instance specified during creation to handle hash calculations for efficient storage
// and retrieval. This implementation ensures thread-safety using a sync.Mutex for concurrent access.
//
// Usage:
// ReflectiveMultiMap is beneficial in scenarios where keys need to map to multiple values,
// and concurrent access safety is essential. It provides methods to get, set, append, remove,
// and clear values associated with keys. The implementation handles potential hash collisions
// and allows efficient storage and retrieval of values under keys.
//
// !IMPORTANT: This map is not safe for concurrent operations. Use ConcurrentMultiMap for concurrent operations.
// This implementation does not preserve the order of insertion, and values under the same key
// can have duplicates if they produce the same hash code. It is optimized for concurrent access
// and efficient storage/retrieval operations, suitable for scenarios requiring a multi-value map
// with concurrency support.
type ReflectiveMultiMap[K comparable, V any] struct {
	store map[K]map[any][]V
}

func NewReflectiveMultiMap[K comparable, V any]() *ReflectiveMultiMap[K, V] {
	return &ReflectiveMultiMap[K, V]{store: make(map[K]map[any][]V)}
}

func (m *ReflectiveMultiMap[K, V]) Get(key K) ([]V, bool) {
	hashMap, ok := m.store[key]
	if !ok {
		return nil, false
	}

	var values []V
	for _, v := range hashMap {
		values = append(values, v...)
	}

	return values, true
}

func (m *ReflectiveMultiMap[K, V]) Set(key K, values ...V) int {
	oldStore, exists := m.store[key]
	matchCount := 0
	if exists {
		for _, v := range values {
			if existingValues, found := oldStore[v]; found {
				matchCount += len(existingValues)
			}
		}
	}

	newStore := make(map[any][]V)
	for _, v := range values {
		if collisions, collisionsFound := newStore[v]; collisionsFound {
			newStore[v] = append(collisions, v)
		} else {
			newStore[v] = []V{v}
		}
	}

	m.store[key] = newStore
	return matchCount
}

func (m *ReflectiveMultiMap[K, V]) Append(key K, values ...V) int {
	hashMap, exists := m.store[key]
	if !exists {
		hashMap = make(map[any][]V)
		m.store[key] = hashMap
	}

	duplicateCount := 0
	for _, v := range values {
		if vs, found := hashMap[v]; found {
			duplicateCount += len(vs)
		}
		hashMap[v] = append(hashMap[v], v)
	}

	return duplicateCount
}

func (m *ReflectiveMultiMap[K, V]) Remove(key K, predicate func(v V) bool) int {
	hashMap, exists := m.store[key]
	if !exists {
		return 0
	}

	removalCount := 0
	for hashCalc, values := range hashMap {
		newValues := make([]V, 0)
		for _, v := range values {
			if predicate(v) {
				removalCount++
			} else {
				newValues = append(newValues, v)
			}
		}
		if len(newValues) == 0 {
			delete(hashMap, hashCalc)
		} else {
			hashMap[hashCalc] = newValues
		}
	}

	return removalCount
}

func (m *ReflectiveMultiMap[K, V]) Clear(key K) bool {
	_, exists := m.store[key]
	if exists {
		delete(m.store, key)
	}

	return exists
}

func (m *ReflectiveMultiMap[K, V]) Iterator() iter.Seq2[K, []V] {
	return func(yield func(K, []V) bool) {
		for i, inm := range m.store {
			for _, v := range inm {
				if !yield(i, v) {
					return
				}
			}
		}
	}
}
