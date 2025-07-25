/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uevent

import (
	"context"
	"sync"
	"sync/atomic"
)

// ParallelWatcher is a generic implementation of the Watcher interface. It ensures that only a single goroutine
// receives specific message from the provided channel. In case you need to listen to each messages in every listener,
// consider using BroadcastWatcher.
//
// The ParallelWatcher struct is parameterized with two types:
// - T: The type of the messages in the channel.
// - C: The type of the channel itself, which should be a channel of T.
//
// Fields:
// - ch: The channel to be watched for incoming messages.
// - f: The function to be called with each message read from the channel.
// - m: A mutex to ensure that Watch is thread-safe.
// - watching: An atomic boolean to track whether a watcher is currently active.
type ParallelWatcher[T any] struct {
	ch <-chan T

	m        sync.Mutex
	watching atomic.Bool

	f atomic.Pointer[watchFunc[T]]
}

func NewParallelWatcher[T any](ch <-chan T, f watchFunc[T]) *ParallelWatcher[T] {
	w := &ParallelWatcher[T]{ch: ch}
	w.Register(f)

	return w
}

// Register replaces the watching function in case it's required.
func (w *ParallelWatcher[T]) Register(f watchFunc[T]) {
	w.f.Store(&f)
}

// Watch starts a goroutine to watch the channel for incoming messages and call the provided function
// with each message. If the channel is nil, the function is nil, or a watcher is already active, it returns false.
// The context.Context parameter allows for cancelling the watching operation.
// Parameters:
// - ctx: The context to control cancellation of the watching operation.
// Returns:
// - A boolean indicating whether the watcher was successfully started.
func (w *ParallelWatcher[T]) Watch(ctx context.Context) bool {
	w.m.Lock()
	defer w.m.Unlock()

	if w.ch == nil || w.f.Load() == nil || !w.watching.CompareAndSwap(false, true) {
		return false
	}

	go func() {
		defer w.watching.Store(false)
		for {
			orig, ok := <-w.ch
			if !ok {
				return
			}

			if ctx.Err() != nil {
				continue
			}

			v := orig
			fptr := w.f.Load()
			go (*fptr)(ctx, &v)
		}
	}()

	return true
}
