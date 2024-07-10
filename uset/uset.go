/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uset

type Set[T comparable] interface {
	Add(value T) bool
	Contains(value T) bool
	Remove(value T) bool
	Size() int
	Clear()
}
