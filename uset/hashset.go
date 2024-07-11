/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uset

// HashSet is a generic set data structure that ensures all elements are unique.
// It uses a map to provide efficient operations for adding, removing, and checking elements.
type HashSet[T comparable] struct {
	m map[T]dummy
}

// NewHashSet creates a new instance of HashSet with the default size
func NewHashSet[T comparable](values ...T) *HashSet[T] {
	m := make(map[T]dummy)
	for _, v := range values {
		m[v] = dummy{}
	}

	return &HashSet[T]{m: m}
}

// NewHashSetWithSize creates a new instance of HashSet with a specified initial size
func NewHashSetWithSize[T comparable](size int) *HashSet[T] {
	return &HashSet[T]{m: make(map[T]dummy, size)}
}

// Add inserts a value into the set and returns true if the value was not already present
func (s *HashSet[T]) Add(value T) bool {
	if s.m == nil {
		s.m = make(map[T]dummy)
	}

	if _, exists := s.m[value]; exists {
		return false
	}

	s.m[value] = def
	return true
}

// Contains checks if a value is present in the set
func (s *HashSet[T]) Contains(value T) bool {
	if s.m == nil {
		s.m = make(map[T]dummy)
	}

	_, exists := s.m[value]
	return exists
}

// Remove deletes a value from the set and returns true if the value was present
func (s *HashSet[T]) Remove(value T) bool {
	if s.m == nil {
		s.m = make(map[T]dummy)
	}

	if _, exists := s.m[value]; exists {
		delete(s.m, value)
		return true
	}
	return false
}

// Size returns the number of elements in the set
func (s *HashSet[T]) Size() int {
	return len(s.m)
}

// Clear returns the number of elements in the set
func (s *HashSet[T]) Clear() {
	s.m = make(map[T]dummy)
}
