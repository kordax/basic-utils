/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uset

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/dgryski/go-farm"
	"github.com/kordax/basic-utils/v2/uref"
	"github.com/kordax/basic-utils/v2/usrlz"
)

type node[T comparable] struct {
	value T
	next  unsafe.Pointer
}

// ConcurrentHashSet is a work-in-progress experimental thread-safe set implementation that uses atomic operations and sync.Map.
// !!! This implementation is not yet optimized and may contain bugs. For stable usage, prefer SynchronizedHashSet.
// Deprecated, use SynchronizedHashSet instead.
type ConcurrentHashSet[T comparable] struct {
	buckets sync.Map

	hash func(value T) uint64
}

// NewCustomConcurrentHashSet creates a new instance of ConcurrentHashSet with specific hashing implementation
func NewCustomConcurrentHashSet[T comparable](hash func(value T) uint64) *ConcurrentHashSet[T] {
	return &ConcurrentHashSet[T]{
		hash: hash,
	}
}

// NewConcurrentHashSet creates a new instance of ConcurrentHashSet with default Farm64 hash implementation
func NewConcurrentHashSet[T comparable]() *ConcurrentHashSet[T] {
	return NewCustomConcurrentHashSet[T](func(value T) uint64 {
		return farm.Hash64(usrlz.ToBytes(&value))
	})
}

// Add inserts a value into the set
func (s *ConcurrentHashSet[T]) Add(value T) bool {
	index := s.hash(value)
	bval, ok := s.buckets.Load(index)

	var head *node[T]
	newNode := &node[T]{value: value}
	if ok {
		ptr := bval.(*unsafe.Pointer)
		head = (*node[T])(atomic.LoadPointer(ptr))
		for n := head; n != nil; n = (*node[T])(atomic.LoadPointer(&n.next)) {
			if n.value == value {
				return false
			}
		}

		for {
			newNode.next = unsafe.Pointer(head)
			if atomic.CompareAndSwapPointer(ptr, unsafe.Pointer(head), unsafe.Pointer(newNode)) {
				return true
			}
			head = (*node[T])(atomic.LoadPointer(ptr))
		}
	} else {
		s.buckets.Store(index, uref.Ref(unsafe.Pointer(newNode)))
		return true
	}
}

// Contains checks if a value is present in the set
func (s *ConcurrentHashSet[T]) Contains(value T) bool {
	index := s.hash(value)
	bval, ok := s.buckets.Load(index)

	if ok {
		head := (*node[T])(atomic.LoadPointer(bval.(*unsafe.Pointer)))

		for n := head; n != nil; n = (*node[T])(atomic.LoadPointer(&n.next)) {
			if n.value == value {
				return true
			}
		}
	}

	return false
}

// Remove deletes a value from the set
func (s *ConcurrentHashSet[T]) Remove(value T) bool {
	index := s.hash(value)
	bval, _ := s.buckets.Load(index)
	ptr := bval.(*unsafe.Pointer)
	head := (*node[T])(atomic.LoadPointer(ptr))

	var prev *node[T]
	for n := head; n != nil; n = (*node[T])(atomic.LoadPointer(&n.next)) {
		if n.value == value {
			if prev == nil {
				return atomic.CompareAndSwapPointer(ptr, unsafe.Pointer(n), n.next)
			}
			return atomic.CompareAndSwapPointer(&prev.next, unsafe.Pointer(n), n.next)
		}
		prev = n
	}

	return false
}

// Size returns the number of elements in the set
func (s *ConcurrentHashSet[T]) Size() int {
	size := 0
	s.buckets.Range(func(k, v interface{}) bool {
		head := (*node[T])(atomic.LoadPointer(v.(*unsafe.Pointer)))
		for n := head; n != nil; n = (*node[T])(atomic.LoadPointer(&n.next)) {
			size++
		}

		return true
	})

	return size
}

// Clear removes all elements from the set
func (s *ConcurrentHashSet[T]) Clear() {
	s.buckets.Range(func(k, v interface{}) bool {
		atomic.StorePointer(v.(*unsafe.Pointer), nil)
		return true
	})
}

// Values retrieves all the values
func (s *ConcurrentHashSet[T]) Values() []T {
	values := make([]T, 0, s.Size())

	s.buckets.Range(func(_, v interface{}) bool {
		values = append(values, v.(T))
		return true
	})

	return values
}
