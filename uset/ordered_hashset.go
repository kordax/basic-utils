/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uset

import "github.com/kordax/basic-utils/uconst"

// OrderedHashSet is a set implementation that preserves the order of elements as they were added.
// This implementation is slower than a traditional HashSet or ComparableHashset, therefore it's recommended to use this
// only in cases you need to retrieve the items in their original order.
type OrderedHashSet[T uconst.UniqueKey[K], K comparable] struct {
	m    map[K]T
	list []T
}

// NewOrderedHashSet creates a new instance of OrderedHashSet.
func NewOrderedHashSet[T uconst.UniqueKey[K], K comparable](values ...T) *OrderedHashSet[T, K] {
	m := make(map[K]T)
	list := make([]T, 0, len(values))

	for _, v := range values {
		key := v.Key()
		if _, exists := m[key]; !exists {
			m[key] = v
			list = append(list, v)
		}
	}

	return &OrderedHashSet[T, K]{m: m, list: list}
}

// Add inserts a value into the set and returns true if the value was not already present.
func (s *OrderedHashSet[T, K]) Add(value T) bool {
	if s.m == nil {
		s.m = make(map[K]T)
	}

	key := value.Key()
	if _, exists := s.m[key]; exists {
		return false
	}

	s.m[key] = value
	s.list = append(s.list, value)
	return true
}

// Contains checks if a value is present in the set.
func (s *OrderedHashSet[T, K]) Contains(value T) bool {
	if s.m == nil {
		s.m = make(map[K]T)
	}

	_, exists := s.m[value.Key()]
	return exists
}

// Get retrieves the element by unique key
func (s *OrderedHashSet[T, K]) Get(key K) *T {
	if s.m == nil {
		s.m = make(map[K]T)
	}

	v, exists := s.m[key]
	if !exists {
		return nil
	}

	return &v
}

// Delete deletes a value from the set by unique key and returns true if the value was present.
func (s *OrderedHashSet[T, K]) Delete(key K) bool {
	if s.m == nil {
		s.m = make(map[K]T)
	}

	if _, exists := s.m[key]; exists {
		delete(s.m, key)
		for i, v := range s.list {
			if v.Key() == key {
				s.list = append(s.list[:i], s.list[i+1:]...)
				break
			}
		}
		return true
	}
	return false
}

// Remove deletes a value from the set and returns true if the value was present.
func (s *OrderedHashSet[T, K]) Remove(value T) bool {
	return s.Delete(value.Key())
}

// Size returns the number of elements in the set.
func (s *OrderedHashSet[T, K]) Size() int {
	return len(s.m)
}

// Clear removes all elements from the set.
func (s *OrderedHashSet[T, K]) Clear() {
	s.m = make(map[K]T)
	s.list = []T{}
}

// OrderedList returns a slice of all elements in the set, in the order they were added.
func (s *OrderedHashSet[T, K]) OrderedList() []T {
	return s.list
}
