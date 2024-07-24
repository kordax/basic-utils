/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uevent

import (
	"context"
)

type watchFunc[T any] func(ctx context.Context, t *T)

// Watcher is an interface that defines a method to start watching a channel for incoming messages.
// The Watch method takes a context.Context to allow for cancellation and returns a boolean indicating
// whether the watcher was successfully started.
type Watcher[T any] interface {
	Register(f watchFunc[T])        // Registers watching function depending on watcher implementation
	Watch(ctx context.Context) bool // Starts watching
}
