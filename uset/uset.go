/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uset

import "github.com/kordax/basic-utils/uconst"

type Set[T comparable] interface {
	Add(value T) bool
	Contains(value T) bool
	Remove(value T) bool
	Size() int
	Clear()
}

type ComparableSet[T uconst.UniqueKey[K], K comparable] interface {
	Add(value T) bool
	Contains(value T) bool
	Remove(value T) bool
	Size() int
	Clear()
	Compare(lhv, rhv T) bool
}
