/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package queue

import (
	"time"

	"github.com/kordax/basic-utils/opt"
)

// FIFOQueueImpl represents a generic implementation of a First-In-First-Out (FIFO) data structure.
// This struct uses a slice to maintain elements in the order they are added, ensuring that
// the oldest item (the first one added) is the first to be fetched.
//
// Fields:
// - queue: Slice holding the actual elements. It grows dynamically as new elements are added.
//
//   - ch: A communication channel utilized in the Poll() method. The channel is used to help
//     in fetching elements with a specified timeout. When a new item is queued and the channel
//     is not full, the new item's pointer is sent into the channel.
//
// Note: This implementation isn't thread-safe. If concurrent access is required please use ConcurrentFIFOQueueImpl.
type FIFOQueueImpl[T any] struct {
	queue []T
	ch    chan *T
}

func NewFIFOQueue[T any](elements ...T) *FIFOQueueImpl[T] {
	return &FIFOQueueImpl[T]{
		queue: elements,
		ch:    make(chan *T),
	}
}

// Queue queues an item. This operation is not thread-safe, and a synchronization wrapper should be provided in case
// consistent results are required in an async environment.
func (q *FIFOQueueImpl[T]) Queue(t T) {
	q.queue = append(q.queue, t)

	select {
	case q.ch <- &t:
	default:
	}
}

func (q *FIFOQueueImpl[T]) Fetch() opt.Opt[T] {
	if len(q.queue) > 0 {
		first := q.queue[0]
		q.queue = q.queue[1:]

		return opt.Of(first)
	} else {
		return opt.Null[T]()
	}
}

func (q *FIFOQueueImpl[T]) Poll(timeout time.Duration) opt.Opt[T] {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for len(q.queue) == 0 {
		select {
		case <-timer.C:
			return opt.Null[T]()
		case t := <-q.ch:
			return opt.OfNullable(t)
		}
	}

	var result *T
	if len(q.queue) > 0 {
		result = &q.queue[0]
		q.queue = q.queue[1:]
	}

	return opt.OfNullable(result)
}

func (q *FIFOQueueImpl[T]) Len() uint64 {
	return uint64(len(q.queue))
}
