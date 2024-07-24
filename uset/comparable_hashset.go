/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uset

import (
	"github.com/kordax/basic-utils/uconst"
	"github.com/kordax/basic-utils/umap"
)

// ComparableHashSet is the same as HashSet
type ComparableHashSet[T uconst.UniqueKey[K], K comparable] struct {
	m map[K]T
}

// NewComparableHashSet creates a new instance of HashSet with the default size
func NewComparableHashSet[T uconst.UniqueKey[K], K comparable](values ...T) *ComparableHashSet[T, K] {
	m := make(map[K]T)
	for _, v := range values {
		m[v.Key()] = v
	}

	return &ComparableHashSet[T, K]{m: m}
}

// NewComparableHashSetWithSize creates a new instance of HashSet with a specified initial size
func NewComparableHashSetWithSize[T uconst.UniqueKey[K], K comparable](size int) *ComparableHashSet[T, K] {
	return &ComparableHashSet[T, K]{m: make(map[K]T, size)}
}

// Add inserts a value into the set and returns true if the value was not already present
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

// Contains checks if a value is present in the set
func (s *ComparableHashSet[T, K]) Contains(value T) bool {
	if s.m == nil {
		s.m = make(map[K]T)
	}

	_, exists := s.m[value.Key()]
	return exists
}

// Remove deletes a value from the set and returns true if the value was present
func (s *ComparableHashSet[T, K]) Remove(value T) bool {
	if s.m == nil {
		s.m = make(map[K]T)
	}

	key := value.Key()
	if _, exists := s.m[key]; exists {
		delete(s.m, key)
		return true
	}
	return false
}

// Size returns the number of elements in the set
func (s *ComparableHashSet[T, K]) Size() int {
	return len(s.m)
}

// Clear returns the number of elements in the set
func (s *ComparableHashSet[T, K]) Clear() {
	s.m = make(map[K]T)
}

// Compare compares two values of type T based on their keys and returns an integer indicating their relative order.
// The keys are obtained by calling the Key() method on the values of type T.
//
// Parameters:
// - lhv: The left-hand value to be compared.
// - rhv: The right-hand value to be compared.
//
// Returns:
// - true if keys are equals
// - false otherwise
//
// The method assumes that T implements a Key() method that returns an integer value which can be used for comparison.
func (s *ComparableHashSet[T, K]) Compare(lhv, rhv T) bool {
	return lhv.Key() == rhv.Key()
}

func (s *ComparableHashSet[T, K]) Values() []T {
	return umap.Values(s.m)
}
