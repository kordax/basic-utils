/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uset

import (
	"git.casinomodule.org/kordax/basic-utils/v2/uconst"
	"git.casinomodule.org/kordax/basic-utils/v2/umap"
)

// ComparableHashSet is the same as HashSet, but allows for custom UniqueKey definition.
type ComparableHashSet[T uconst.UniqueKey[K], K comparable] struct {
	m map[K]T
}

// NewComparableHashSet creates a new instance of HashSet with initial capacity equal to values length.
func NewComparableHashSet[T uconst.UniqueKey[K], K comparable](values ...T) *ComparableHashSet[T, K] {
	m := make(map[K]T, len(values))
	for _, v := range values {
		m[v.Key()] = v
	}

	return &ComparableHashSet[T, K]{m: m}
}

// NewComparableHashSetWithSize creates a new instance of HashSet with a specified initial size.
func NewComparableHashSetWithSize[T uconst.UniqueKey[K], K comparable](size int) *ComparableHashSet[T, K] {
	return &ComparableHashSet[T, K]{m: make(map[K]T, size)}
}

// Add inserts a value into the set and returns true if the value was not already present.
func (s *ComparableHashSet[T, K]) Add(value T) bool {
	if s.m == nil {
		s.m = make(map[K]T)
	}

	key := value.Key()
	if _, exists := s.m[key]; exists {
		return false
	}

	s.m[key] = value
	return true
}

// Contains checks if a value is present in the set.
func (s *ComparableHashSet[T, K]) Contains(value T) bool {
	if s.m == nil {
		return false
	}

	_, exists := s.m[value.Key()]
	return exists
}

// Remove deletes a value from the set and returns true if the value was present.
func (s *ComparableHashSet[T, K]) Remove(value T) bool {
	if s.m == nil {
		return false
	}

	key := value.Key()
	if _, exists := s.m[key]; exists {
		delete(s.m, key)
		return true
	}

	return false
}

// Size returns the number of elements in the set.
func (s *ComparableHashSet[T, K]) Size() int {
	return len(s.m)
}

// Clear removes all elements from the set.
func (s *ComparableHashSet[T, K]) Clear() {
	s.m = nil
}

func (s *ComparableHashSet[T, K]) Values() []T {
	return umap.Values(s.m)
}
