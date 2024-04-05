/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ustream

import "github.com/kordax/basic-utils/uarray"

// Collector defines the interface for collecting elements from a stream.
type Collector[T any] interface {
	Collect() []T
	CollectToMap(func(*T) (any, any)) map[any]any
}

type Stream[T any] struct {
	values []T
}

// NewStream creates a new Stream from the given slice.
func NewStream[T any](values []T) *Stream[T] {
	return &Stream[T]{values: values}
}

// Filter wraps the uarray.Filter function in the stream API.
func (s *Stream[T]) Filter(predicate func(*T) bool) *Stream[T] {
	return NewStream(uarray.Filter(s.values, predicate))
}

// FilterOut wraps the uarray.FilterOut function in the stream API.
func (s *Stream[T]) FilterOut(predicate func(*T) bool) *Stream[T] {
	return NewStream(uarray.FilterOut(s.values, predicate))
}

// Map wraps the uarray.Map function in the stream API.
// This operation can only be the last in a pipeline.
func (s *Stream[T]) Map(mapper func(*T) any) *TerminalStream[any] {
	return NewTerminalStream(uarray.Map(s.values, mapper))
}

func (s *Stream[T]) Collect() []T {
	return s.values
}

// CollectToMap uses the uarray.ToMap function to collect elements of the Stream into a map.
func (s *Stream[T]) CollectToMap(mapper func(*T) (any, any)) map[any]any {
	return uarray.ToMap(s.values, mapper)
}

// TerminalStream represents a stream that can only be collected.
type TerminalStream[T any] struct {
	values []T
}

func NewTerminalStream[T any](values []T) *TerminalStream[T] {
	return &TerminalStream[T]{values: values}
}

func (s *TerminalStream[T]) Collect() []T {
	return s.values
}

func (s *TerminalStream[T]) CollectToMap(mapper func(*T) (any, T)) map[any]T {
	return uarray.ToMap(s.values, mapper)
}
