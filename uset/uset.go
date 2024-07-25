/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uset

import "github.com/kordax/basic-utils/uconst"

// Set that contains unique elements.
type Set[T comparable] interface {
	Add(value T) bool
	Contains(value T) bool
	Remove(value T) bool
	Size() int
	Clear()
	Values() []T
}

// OrderedSet is an interface that represents a set of unique items that can be retrieved in the original insertion order.
// It is parameterized by T, which must implement uconst.UniqueKey[K], and K, which must be comparable.
type OrderedSet[T uconst.UniqueKey[K], K comparable] interface {
	Add(value T) bool
	Contains(value T) bool
	Get(key K) *T // Retrieve value by key
	Delete(key K) bool
	Remove(value T) bool
	Size() int
	Clear()
	OrderedList() []T // OrderedList returns the slice of ordered items while preserving their natural order
}
