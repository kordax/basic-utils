/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap

import (
	"hash"
	"sync"
)

type ReflectiveMultiMap[K comparable, V any] struct {
	h hash.Hash
	sync.Mutex
	store map[K]map[int64][]V
}

func NewReflectiveMultiMap[K comparable, V any](hashingMethod hash.Hash) *ReflectiveMultiMap[K, V] {
	return &ReflectiveMultiMap[K, V]{h: hashingMethod, store: make(map[K]map[int64][]V)}
}

func (m *ReflectiveMultiMap[K, V]) Get(key K) ([]V, bool) {
	m.Lock()
	defer m.Unlock()

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
	m.Lock()
	defer m.Unlock()

	oldStore, exists := m.store[key]
	matchCount := 0
	if exists {
		for _, v := range values {
			hashCalc := computeHash(m.h, v)
			if existingValues, found := oldStore[hashCalc]; found {
				matchCount += len(existingValues)
			}
		}
	}

	newStore := make(map[int64][]V)
	for _, v := range values {
		hashCalc := computeHash(m.h, v)
		if collisions, collisionsFound := newStore[hashCalc]; collisionsFound {
			newStore[hashCalc] = append(collisions, v)
		} else {
			newStore[hashCalc] = []V{v}
		}
	}

	m.store[key] = newStore
	return matchCount
}

func (m *ReflectiveMultiMap[K, V]) Append(key K, values ...V) int {
	m.Lock()
	defer m.Unlock()

	hashMap, exists := m.store[key]
	if !exists {
		hashMap = make(map[int64][]V)
		m.store[key] = hashMap
	}

	duplicateCount := 0
	for _, v := range values {
		hashCalc := computeHash(m.h, v)
		if vs, found := hashMap[hashCalc]; found {
			duplicateCount += len(vs)
		}
		hashMap[hashCalc] = append(hashMap[hashCalc], v)
	}

	return duplicateCount
}

func (m *ReflectiveMultiMap[K, V]) Remove(key K, predicate func(v V) bool) int {
	m.Lock()
	defer m.Unlock()

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
	m.Lock()
	defer m.Unlock()

	_, exists := m.store[key]
	if exists {
		delete(m.store, key)
	}

	return exists
}
