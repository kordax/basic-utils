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

type Queue[T any] interface {
	Queue(t T)
	Fetch() opt.Opt[T]
	Poll(timeout time.Duration) opt.Opt[T]
	Len() uint64
}

type PriorityQueue[T any] interface {
	Queue(t T, priority int)
	Fetch() opt.Opt[T]
	Poll(timeout time.Duration) opt.Opt[T]
	Len() uint64
}
