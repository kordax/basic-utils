/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uqueue

import (
	"slices"
	"sync"
	"time"

	"git.casinomodule.org/casino27/basic-utils/v2/uopt"
)

// FIFOQueueImpl represents a generic implementation of a First-In-First-Out (FIFO) data structure.
// This struct uses a slice to maintain elements in the order they are added, ensuring that
// the oldest item (the first one added) is the first to be fetched.
//
// Fields:
// - queue: Slice holding the actual elements. It grows dynamically as new elements are added.
//
//   - ch: A communication channel utilized in the Poll() method. The channel is used to notify
//     waiting pollers when queue state changes.
type FIFOQueueImpl[T any] struct {
	mu    sync.Mutex
	queue []T
	ch    chan struct{}
}

func NewFIFOQueue[T any](elements ...T) *FIFOQueueImpl[T] {
	return &FIFOQueueImpl[T]{
		queue: slices.Clone(elements),
		ch:    make(chan struct{}, 1),
	}
}

// Queue queues an item.
func (q *FIFOQueueImpl[T]) Queue(t T) {
	q.mu.Lock()
	q.queue = append(q.queue, t)
	q.mu.Unlock()

	q.notify()
}

func (q *FIFOQueueImpl[T]) Fetch() uopt.Opt[T] {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.queue) == 0 {
		return uopt.Null[T]()
	}

	first := q.queue[0]
	var zero T
	q.queue[0] = zero
	q.queue = q.queue[1:]
	if len(q.queue) > 0 {
		defer q.notify()
	}

	return uopt.Of(first)
}

func (q *FIFOQueueImpl[T]) Poll(timeout time.Duration) uopt.Opt[T] {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		if result := q.Fetch(); result.Present() {
			return result
		}

		select {
		case <-timer.C:
			return uopt.Null[T]()
		case <-q.ch:
		}
	}
}

func (q *FIFOQueueImpl[T]) Len() uint64 {
	q.mu.Lock()
	defer q.mu.Unlock()

	return uint64(len(q.queue))
}

func (q *FIFOQueueImpl[T]) notify() {
	select {
	case q.ch <- struct{}{}:
	default:
	}
}
