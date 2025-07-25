/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uqueue

import (
	"time"

	"github.com/kordax/basic-utils/v2/uopt"
)

type Queue[T any] interface {
	Queue(t T)
	Fetch() uopt.Opt[T]
	Poll(timeout time.Duration) uopt.Opt[T]
	Len() uint64
}

type PriorityQueue[T any] interface {
	Queue(t T, priority int)
	Fetch() uopt.Opt[T]
	Poll(timeout time.Duration) uopt.Opt[T]
	Len() uint64
}
