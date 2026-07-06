/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uqueue

import (
	"context"
	"sync/atomic"
	"time"
	"unsafe"

	"git.casinomodule.org/casino27/basic-utils/v2/uopt"
)

type node[T any] struct {
	value T
	next  unsafe.Pointer
}

// ConcurrentFIFOQueueImpl implementation that queues/dequeues items according to https://www.cs.rochester.edu/~scott/papers/1996_PODC_queues.pdf
type ConcurrentFIFOQueueImpl[T any] struct {
	head unsafe.Pointer
	tail unsafe.Pointer
	ch   chan struct{}
	l    atomic.Uint64
}

func NewConcurrentFIFOQueueImpl[T any]() *ConcurrentFIFOQueueImpl[T] {
	n := &node[T]{}
	return &ConcurrentFIFOQueueImpl[T]{
		head: unsafe.Pointer(n),
		tail: unsafe.Pointer(n),
		ch:   make(chan struct{}, 1),
	}
}

// Queue queues an item in the finite time. This operation is thread-safe yet is not "synchronized" by its nature.
// Implementation uses Michael-Scott CAS implementation.
func (q *ConcurrentFIFOQueueImpl[T]) Queue(t T) {
	newNode := &node[T]{value: t}
	for {
		tail := load[T](&q.tail)
		next := load[T](&tail.next)
		if tail == load[T](&q.tail) {
			if next == nil {
				if cas[T](&tail.next, next, newNode) {
					break
				}
			} else {
				cas[T](&q.tail, tail, next)
			}
		}
	}

	q.l.Add(1)

	q.notify()
}

// Fetch fetches item in the finite time.
func (q *ConcurrentFIFOQueueImpl[T]) Fetch() uopt.Opt[T] {
	for {
		head := load[T](&q.head)
		tail := load[T](&q.tail)
		next := load[T](&head.next)
		if head == load[T](&q.head) {
			if head == tail {
				if next == nil {
					return uopt.Null[T]()
				}
				cas[T](&q.tail, tail, next)
			} else {
				value := next.value
				if cas[T](&q.head, head, next) {
					if q.l.Add(^uint64(0)) > 0 {
						q.notify()
					}
					return uopt.Of(value)
				}
			}
		}
	}
}

// Poll items fetches item in the finite time.
func (q *ConcurrentFIFOQueueImpl[T]) Poll(timeout time.Duration) uopt.Opt[T] {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return q.PollContext(ctx)
}

func (q *ConcurrentFIFOQueueImpl[T]) PollContext(ctx context.Context) uopt.Opt[T] {
	if ctx == nil {
		ctx = context.Background()
	}

	for {
		if result := q.Fetch(); result.Present() {
			return result
		}

		select {
		case <-ctx.Done():
			return uopt.Null[T]()
		case <-q.ch:
		}
	}
}

func (q *ConcurrentFIFOQueueImpl[T]) Peek() uopt.Opt[T] {
	head := load[T](&q.head)
	next := load[T](&head.next)
	if next == nil {
		return uopt.Null[T]()
	}

	return uopt.Of(next.value)
}

func (q *ConcurrentFIFOQueueImpl[T]) Drain(limit ...int) []T {
	n := int(q.Len())
	if len(limit) > 0 && limit[0] >= 0 && limit[0] < n {
		n = limit[0]
	}
	if n == 0 {
		return []T{}
	}

	result := make([]T, 0, n)
	for len(result) < n {
		value := q.Fetch()
		if !value.Present() {
			break
		}
		result = append(result, value.OrElse(*new(T)))
	}

	return result
}

func (q *ConcurrentFIFOQueueImpl[T]) Clear() {
	n := &node[T]{}
	atomic.StorePointer(&q.head, unsafe.Pointer(n))
	atomic.StorePointer(&q.tail, unsafe.Pointer(n))
	q.l.Store(0)
}

func (q *ConcurrentFIFOQueueImpl[T]) Empty() bool {
	return q.Len() == 0
}

func (q *ConcurrentFIFOQueueImpl[T]) Len() uint64 {
	return q.l.Load()
}

func load[T any](ptr *unsafe.Pointer) *node[T] {
	return (*node[T])(atomic.LoadPointer(ptr))
}

func cas[T any](ptr *unsafe.Pointer, old, new *node[T]) bool {
	return atomic.CompareAndSwapPointer(ptr, unsafe.Pointer(old), unsafe.Pointer(new))
}

func (q *ConcurrentFIFOQueueImpl[T]) notify() {
	select {
	case q.ch <- struct{}{}:
	default:
	}
}
