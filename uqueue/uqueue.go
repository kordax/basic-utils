/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uqueue

import (
	"context"
	"time"

	"github.com/kordax/basic-utils/v4/uopt"
)

type Queue[T any] interface {
	Queue(t T)
	Fetch() uopt.Opt[T]
	Poll(timeout time.Duration) uopt.Opt[T]
	PollContext(ctx context.Context) uopt.Opt[T]
	Peek() uopt.Opt[T]
	Drain(limit ...int) []T
	Clear()
	Empty() bool
	Len() uint64
}

type PriorityQueue[T any] interface {
	Queue(t T, priority int)
	Fetch() uopt.Opt[T]
	Poll(timeout time.Duration) uopt.Opt[T]
	PollContext(ctx context.Context) uopt.Opt[T]
	Peek() uopt.Opt[T]
	Drain(limit ...int) []T
	Clear()
	Empty() bool
	Len() uint64
}
