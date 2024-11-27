/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap

import (
	"hash"
	"iter"
)

type wrapper[T comparable] struct {
	h hash.Hash
	v T
}

func newWrapper[T comparable](h hash.Hash, v T) wrapper[T] {
	return wrapper[T]{h: h, v: v}
}

func (w wrapper[T]) Hash() int64 {
	w.h.Reset()
	return computeHash(w.h, w.v)
}

// UniqueHashMultiMap is the same as UniqueMultiMap, but allows to specify custom hash function.
type UniqueHashMultiMap[K comparable, V comparable] struct {
	h  hash.Hash
	mm *UniqueMultiMap[K, wrapper[V]]
}

func NewUniqueHashMultiMap[K comparable, V comparable](hashMethod hash.Hash) *UniqueHashMultiMap[K, V] {
	return &UniqueHashMultiMap[K, V]{
		h:  hashMethod,
		mm: NewUniqueMultiMap[K, wrapper[V]](),
	}
}

func (m *UniqueHashMultiMap[K, V]) Get(key K) (value []V, ok bool) {
	wrappedResult, ok := m.mm.Get(key)
	if !ok {
		return nil, false
	}

	result := make([]V, len(wrappedResult))
	for i, v := range wrappedResult {
		result[i] = v.v
	}

	return result, true
}

func (m *UniqueHashMultiMap[K, V]) Set(key K, values ...V) int {
	wrappedValues := make([]wrapper[V], len(values))
	for i, v := range values {
		wrappedValues[i] = newWrapper(m.h, v)
	}
	return m.mm.Set(key, wrappedValues...)
}

func (m *UniqueHashMultiMap[K, V]) Append(key K, values ...V) int {
	wrappedValues := make([]wrapper[V], len(values))
	for i, v := range values {
		wrappedValues[i] = newWrapper(m.h, v)
	}
	return m.mm.Append(key, wrappedValues...)
}

func (m *UniqueHashMultiMap[K, V]) Remove(key K, predicate func(v V) bool) int {
	wrappedPredicate := func(wv wrapper[V]) bool {
		return predicate(wv.v)
	}
	return m.mm.Remove(key, wrappedPredicate)
}

func (m *UniqueHashMultiMap[K, V]) Clear(key K) bool {
	return m.mm.Clear(key)
}

func (m *UniqueHashMultiMap[K, V]) Iterator() iter.Seq2[K, []V] {
	return func(yield func(K, []V) bool) {
		for i, wrappers := range m.mm.Iterator() {
			values := make([]V, len(wrappers))
			for n, w := range wrappers {
				values[n] = w.v
			}
			if !yield(i, values) {
				return
			}
		}
	}
}
