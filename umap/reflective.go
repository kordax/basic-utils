/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap

import "iter"

// ReflectiveMultiMap is a multi-value map: key K -> []V.
// !IMPORTANT: This map is NOT safe for concurrent use. Wrap it with your own sync if needed.
type ReflectiveMultiMap[K comparable, V comparable] struct {
	store map[K]map[V][]V
}

func NewReflectiveMultiMap[K comparable, V comparable]() *ReflectiveMultiMap[K, V] {
	return &ReflectiveMultiMap[K, V]{store: make(map[K]map[V][]V)}
}

func (m *ReflectiveMultiMap[K, V]) Get(key K) ([]V, bool) {
	hashMap, ok := m.store[key]
	if !ok {
		return nil, false
	}

	values := make([]V, 0)
	for _, v := range hashMap {
		values = append(values, v...)
	}

	return values, true
}

// Set replaces all values for a key. It returns the number of previously stored
// values that are equal to any of the provided values.
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

	newStore := make(map[V][]V, len(values))
	for _, v := range values {
		newStore[v] = append(newStore[v], v)
	}

	m.store[key] = newStore
	return matchCount
}

// Append adds values to an existing key. It returns the number of duplicates
// encountered (counted against already stored values under that key).
func (m *ReflectiveMultiMap[K, V]) Append(key K, values ...V) int {
	hashMap, exists := m.store[key]
	if !exists {
		hashMap = make(map[V][]V, len(values))
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

// Remove deletes values under key that match predicate. It returns the number
// of removed values.
func (m *ReflectiveMultiMap[K, V]) Remove(key K, predicate func(v V) bool) int {
	hashMap, exists := m.store[key]
	if !exists {
		return 0
	}

	removalCount := 0
	for hashKey, values := range hashMap {
		newValues := make([]V, 0, len(values))
		for _, v := range values {
			if predicate(v) {
				removalCount++
			} else {
				newValues = append(newValues, v)
			}
		}
		if len(newValues) == 0 {
			delete(hashMap, hashKey)
		} else {
			hashMap[hashKey] = newValues
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

// Iterator returns an iterator over all (key, []V-bucket) pairs.
// A single logical key K may be yielded multiple times if it has
// multiple buckets (i.e. multiple distinct V "hash keys").
func (m *ReflectiveMultiMap[K, V]) Iterator() iter.Seq2[K, []V] {
	return func(yield func(K, []V) bool) {
		for k, bucketMap := range m.store {
			for _, v := range bucketMap {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}
