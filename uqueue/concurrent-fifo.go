/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uqueue

import (
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/kordax/basic-utils/v2/uopt"
)

type node[T any] struct {
	value T
	next  unsafe.Pointer
}

// ConcurrentFIFOQueueImpl implementation that queues/dequeues items according to https://www.cs.rochester.edu/~scott/papers/1996_PODC_queues.pdf
type ConcurrentFIFOQueueImpl[T any] struct {
	head unsafe.Pointer
	tail unsafe.Pointer
	ch   chan *T
	l    atomic.Uint64
}

func NewConcurrentFIFOQueueImpl[T any]() *ConcurrentFIFOQueueImpl[T] {
	n := &node[T]{}
	return &ConcurrentFIFOQueueImpl[T]{
		head: unsafe.Pointer(n),
		tail: unsafe.Pointer(n),
		ch:   make(chan *T),
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

	select {
	case q.ch <- &t:
	default:
	}
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
					q.l.Add(^uint64(0))
					return uopt.Of(value)
				}
			}
		}
	}
}

// Poll items fetches item in the finite time.
func (q *ConcurrentFIFOQueueImpl[T]) Poll(timeout time.Duration) uopt.Opt[T] {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	r := q.Fetch()
	for !r.Present() {
		select {
		case <-timer.C:
			return uopt.Null[T]()
		case <-q.ch:
			return q.Fetch()
		}
	}

	return r
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
