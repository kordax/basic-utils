/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uevent

import (
	"context"
)

type watchFunc[T any] func(ctx context.Context, t T)

// Subscription controls the lifetime of a registered event listener.
type Subscription interface {
	// Unsubscribe removes the listener. It returns true only for the call that performs the removal.
	Unsubscribe() bool
}

// Watcher starts watching a channel and dispatches its messages to registered listeners.
type Watcher[T any] interface {
	Register(f watchFunc[T])        // Registers watching function depending on watcher implementation
	Watch(ctx context.Context) bool // Starts watching
}
