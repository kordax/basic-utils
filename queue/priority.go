/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package queue

import (
	"container/heap"
	"time"

	"github.com/kordax/basic-utils/opt"
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
//   - ch: A communication channel utilized in the Poll() method. The channel is used to assist
//     in fetching elements with a specified timeout. When a new item is queued and the channel
//     is not full, the new item's pointer is sent into the channel.
//
// Note: This implementation isn't inherently thread-safe. If concurrent access is anticipated,
//
//	external synchronization mechanisms should be used, or you can use ConcurrentFIFOQueueImpl.
type PriorityQueueImpl[T any] struct {
	queue *prioritizedQueue[T]
	ch    chan *T
}

func NewPriorityQueue[T any]() *PriorityQueueImpl[T] {
	pq := &prioritizedQueue[T]{}
	heap.Init(pq)
	return &PriorityQueueImpl[T]{
		queue: pq,
		ch:    make(chan *T),
	}
}

func (q *PriorityQueueImpl[T]) Queue(t T, priority int) {
	heap.Push(q.queue, &container[T]{
		t:        &t,
		priority: priority,
		index:    q.queue.Len(),
	})

	select {
	case q.ch <- &t:
	default:
	}
}

func (q *PriorityQueueImpl[T]) Fetch() opt.Opt[T] {
	r := q.queue.Pop()
	if r == nil {
		return opt.Null[T]()
	} else {
		return opt.OfNullable[T](r.(*container[T]).t)
	}
}

func (q *PriorityQueueImpl[T]) Poll(timeout time.Duration) opt.Opt[T] {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	r := q.Fetch()
	for !r.Present() {
		select {
		case <-timer.C:
			return opt.Null[T]()
		case t := <-q.ch:
			return opt.OfNullable(t)
		}
	}

	return r
}

func (q *PriorityQueueImpl[T]) Len() uint64 {
	return uint64(q.queue.Len())
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
	old := make([]*container[T], len(pq.e))
	copy(old, pq.e)
	n := len(old)
	if n == 0 {
		return nil
	}

	item := old[0]
	old[0] = nil    // avoid memory leak
	item.index = -1 // for safety
	pq.e = pq.e[1:]

	return item
}
