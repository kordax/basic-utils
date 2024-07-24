/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uset

import "github.com/kordax/basic-utils/uconst"

// OrderedHashSet is a set implementation that preserves the order of elements as they were added.
type OrderedHashSet[T uconst.UniqueKey[K], K comparable] struct {
	m    map[K]dummy
	list []T
}

// NewOrderedHashSet creates a new instance of OrderedHashSet.
func NewOrderedHashSet[T uconst.UniqueKey[K], K comparable](values ...T) *OrderedHashSet[T, K] {
	m := make(map[K]dummy)
	list := make([]T, 0, len(values))

	for _, v := range values {
		key := v.Key()
		if _, exists := m[key]; !exists {
			m[key] = def
			list = append(list, v)
		}
	}

	return &OrderedHashSet[T, K]{m: m, list: list}
}

// Add inserts a value into the set and returns true if the value was not already present.
func (s *OrderedHashSet[T, K]) Add(value T) bool {
	if s.m == nil {
		s.m = make(map[K]dummy)
	}

	key := value.Key()
	if _, exists := s.m[key]; exists {
		return false
	}

	s.m[key] = def
	s.list = append(s.list, value)
	return true
}

// Contains checks if a value is present in the set.
func (s *OrderedHashSet[T, K]) Contains(value T) bool {
	if s.m == nil {
		s.m = make(map[K]dummy)
	}

	_, exists := s.m[value.Key()]
	return exists
}

// Remove deletes a value from the set and returns true if the value was present.
func (s *OrderedHashSet[T, K]) Remove(value T) bool {
	if s.m == nil {
		s.m = make(map[K]dummy)
	}

	key := value.Key()
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

// Size returns the number of elements in the set.
func (s *OrderedHashSet[T, K]) Size() int {
	return len(s.m)
}

// Clear removes all elements from the set.
func (s *OrderedHashSet[T, K]) Clear() {
	s.m = make(map[K]dummy)
	s.list = []T{}
}

// AsSlice returns a slice of all elements in the set, in the order they were added.
func (s *OrderedHashSet[T, K]) AsSlice() []T {
	return s.list
}
