/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap

import "hash"

type wrapper[T any] struct {
	h hash.Hash
	v T
}

func newWrapper[T any](h hash.Hash, v T) wrapper[T] {
	return wrapper[T]{h: h, v: v}
}

func (w wrapper[T]) Hash() int64 {
	w.h.Reset()
	return computeHash(w.h, w.v)
}

// UniqueHashMultiMap is the same as UniqueMultiMap, but allows to specify custom hash function.
type UniqueHashMultiMap[K comparable, V any] struct {
	h  hash.Hash
	mm *UniqueMultiMap[K, wrapper[V]]
}

func NewUniqueHashMultiMap[K comparable, V any](hashMethod hash.Hash) *UniqueHashMultiMap[K, V] {
	return &UniqueHashMultiMap[K, V]{
		h:  hashMethod,
		mm: NewUniqueMultiMap[K, wrapper[V]](),
	}
}

func (h *UniqueHashMultiMap[K, V]) Get(key K) (value []V, ok bool) {
	wrappedResult, ok := h.mm.Get(key)
	if !ok {
		return nil, false
	}

	result := make([]V, len(wrappedResult))
	for i, v := range wrappedResult {
		result[i] = v.v
	}

	return result, true
}

func (h *UniqueHashMultiMap[K, V]) Set(key K, values ...V) int {
	wrappedValues := make([]wrapper[V], len(values))
	for i, v := range values {
		wrappedValues[i] = newWrapper(h.h, v)
	}
	return h.mm.Set(key, wrappedValues...)
}

func (h *UniqueHashMultiMap[K, V]) Append(key K, values ...V) int {
	wrappedValues := make([]wrapper[V], len(values))
	for i, v := range values {
		wrappedValues[i] = newWrapper(h.h, v)
	}
	return h.mm.Append(key, wrappedValues...)
}

func (h *UniqueHashMultiMap[K, V]) Remove(key K, predicate func(v V) bool) int {
	wrappedPredicate := func(wv wrapper[V]) bool {
		return predicate(wv.v)
	}
	return h.mm.Remove(key, wrappedPredicate)
}

func (h *UniqueHashMultiMap[K, V]) Clear(key K) bool {
	return h.mm.Clear(key)
}
