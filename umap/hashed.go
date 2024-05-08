/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap

import (
	"sync"
)

type Hashed interface {
	Hash() int64
}

type HashedMultiMap[K comparable, V Hashed] struct {
	sync.Mutex
	store map[K]map[int64]V
}

func NewHashedMultiMap[K comparable, V Hashed]() *HashedMultiMap[K, V] {
	return &HashedMultiMap[K, V]{
		store: make(map[K]map[int64]V),
	}
}

func (m *HashedMultiMap[K, V]) Get(key K) ([]V, bool) {
	m.Lock()
	defer m.Unlock()

	hashMap, ok := m.store[key]
	if !ok {
		return nil, false
	}

	values := make([]V, 0, len(hashMap))
	for _, value := range hashMap {
		values = append(values, value)
	}

	return values, true
}

func (m *HashedMultiMap[K, V]) Set(key K, values ...V) int {
	m.Lock()
	defer m.Unlock()

	hashMap := make(map[int64]V)
	addedCount := 0
	for _, value := range values {
		hash := value.Hash()
		if _, exists := hashMap[hash]; !exists {
			hashMap[hash] = value
			addedCount++
		}
	}

	m.store[key] = hashMap
	return addedCount
}

func (m *HashedMultiMap[K, V]) Append(key K, values ...V) int {
	m.Lock()
	defer m.Unlock()

	hashMap, exists := m.store[key]
	if !exists {
		hashMap = make(map[int64]V)
		m.store[key] = hashMap
	}

	addedCount := 0
	for _, value := range values {
		hash := value.Hash()
		if _, found := hashMap[hash]; !found {
			hashMap[hash] = value
			addedCount++
		}
	}

	return addedCount
}

func (m *HashedMultiMap[K, V]) Remove(key K, predicate func(v V) bool) int {
	m.Lock()
	defer m.Unlock()

	hashMap, exists := m.store[key]
	if !exists {
		return 0
	}

	removalCount := 0
	for hash, value := range hashMap {
		if predicate(value) {
			delete(hashMap, hash)
			removalCount++
		}
	}

	return removalCount
}

func (m *HashedMultiMap[K, V]) Clear(key K) bool {
	m.Lock()
	defer m.Unlock()

	_, exists := m.store[key]
	if exists {
		delete(m.store, key)
		return true
	}

	return false
}
