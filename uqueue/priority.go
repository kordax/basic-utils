/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uqueue

import (
	"container/heap"
	"context"
	"sync"
	"time"

	"github.com/kordax/basic-utils/v3/uopt"
)

// container represents an individual item in the priority queue.
// It wraps an item of generic type T, associating a priority and an index with it.
//
// Fields:
// - t: Holds the actual data element of type T.
//
// - priority: Represents the priority associated with this item. Higher values denote higher priorities.
//
//   - index: Stores the current position of the item in the heap. It aids in maintaining
//     the priority order of items during swap operations.
type container[T any] struct {
	t        *T
	priority int
	index    int
}

// PriorityQueueImpl represents a generic implementation of a priority queue data structure.
// Items are organized based on their priorities, with higher-priority items being fetched before lower-priority ones.
// This struct uses a heap data structure (as implemented in the "container/heap" package) to efficiently manage
// the priorities and retrieval of items.
//
// Fields:
// - queue: The underlying heap structure (prioritizedQueue) that manages the prioritized items.
//
//   - ch: A communication channel utilized in the Poll() method. The channel is used to notify
//     waiting pollers when queue state changes.
type PriorityQueueImpl[T any] struct {
	mu    sync.Mutex
	queue *prioritizedQueue[T]
	ch    chan struct{}
}

func NewPriorityQueue[T any]() *PriorityQueueImpl[T] {
	pq := &prioritizedQueue[T]{}
	heap.Init(pq)
	return &PriorityQueueImpl[T]{
		queue: pq,
		ch:    make(chan struct{}, 1),
	}
}

func (q *PriorityQueueImpl[T]) Queue(t T, priority int) {
	q.mu.Lock()
	heap.Push(q.queue, &container[T]{
		t:        &t,
		priority: priority,
		index:    q.queue.Len(),
	})
	q.mu.Unlock()

	q.notify()
}

func (q *PriorityQueueImpl[T]) Fetch() uopt.Opt[T] {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.queue.Len() == 0 {
		return uopt.Null[T]()
	}

	r := heap.Pop(q.queue)
	if r == nil {
		return uopt.Null[T]()
	}
	if q.queue.Len() > 0 {
		defer q.notify()
	}

	return uopt.OfNullable[T](r.(*container[T]).t)
}

func (q *PriorityQueueImpl[T]) Poll(timeout time.Duration) uopt.Opt[T] {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return q.PollContext(ctx)
}

func (q *PriorityQueueImpl[T]) PollContext(ctx context.Context) uopt.Opt[T] {
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

func (q *PriorityQueueImpl[T]) Peek() uopt.Opt[T] {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.queue.Len() == 0 {
		return uopt.Null[T]()
	}

	return uopt.OfNullable(q.queue.e[0].t)
}

func (q *PriorityQueueImpl[T]) Drain(limit ...int) []T {
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

func (q *PriorityQueueImpl[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	for i := range q.queue.e {
		q.queue.e[i] = nil
	}
	q.queue.e = nil
}

func (q *PriorityQueueImpl[T]) Empty() bool {
	return q.Len() == 0
}

func (q *PriorityQueueImpl[T]) Len() uint64 {
	q.mu.Lock()
	defer q.mu.Unlock()

	return uint64(q.queue.Len())
}

func (q *PriorityQueueImpl[T]) notify() {
	select {
	case q.ch <- struct{}{}:
	default:
	}
}

type prioritizedQueue[T any] struct {
	e []*container[T]
}

func (pq *prioritizedQueue[T]) Len() int { return len(pq.e) }

func (pq *prioritizedQueue[T]) Less(i, j int) bool {
	return pq.e[i].priority > pq.e[j].priority
}

func (pq *prioritizedQueue[T]) Swap(i, j int) {
	pq.e[i], pq.e[j] = pq.e[j], pq.e[i]
	pq.e[i].index = i
	pq.e[j].index = j
}

func (pq *prioritizedQueue[T]) Push(x any) {
	if c, ok := x.(*container[T]); ok {
		c.index = pq.Len()
		pq.e = append(pq.e, c)
	}
}

func (pq *prioritizedQueue[T]) Pop() any {
	old := pq.e
	n := len(old)
	if n == 0 {
		return nil
	}

	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	pq.e = old[:n-1]

	return item
}
