/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uset

import (
	"sync"
)

// SynchronizedHashSet is a thread-safe wrapper around HashSet that ensures concurrent safety.
type SynchronizedHashSet[T comparable] struct {
	hs *HashSet[T]

	mtx sync.RWMutex
}

func NewSynchronizedHashSet[T comparable](values ...T) *SynchronizedHashSet[T] {
	return &SynchronizedHashSet[T]{hs: NewHashSet[T](values...)}
}

func NewSynchronizedHashSetFromSet[T comparable](hs *HashSet[T]) *SynchronizedHashSet[T] {
	return &SynchronizedHashSet[T]{hs: hs}
}

func (s *SynchronizedHashSet[T]) Add(value T) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.hs.Add(value)
}

func (s *SynchronizedHashSet[T]) Contains(value T) bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.hs.Contains(value)
}

func (s *SynchronizedHashSet[T]) Remove(value T) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.hs.Remove(value)
}

func (s *SynchronizedHashSet[T]) Size() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.hs.Size()
}

func (s *SynchronizedHashSet[T]) Clear() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.hs.Clear()
}

// Values retrieves all the values
func (s *SynchronizedHashSet[T]) Values() []T {
	return s.hs.Values()
}
