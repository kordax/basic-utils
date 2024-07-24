package uevent

import (
	"context"
	"sync"
	"sync/atomic"
)

// BroadcastWatcher is an implementation of the Watcher interface that allows multiple parallel listeners
// to receive the same messages from a single channel.
// It's pretty similar to ParallelWatcher, but allows multiple parallel listeners to receive the same messages.
//
// Fields:
// - input: The input channel to be watched for incoming messages.
// - listeners: A slice of listeners channels for registered listeners.
// - m: A mutex to ensure that Register and other operations are thread-safe.
// - started: An atomic boolean to ensure the Start method is only called once.
type BroadcastWatcher[T any] struct {
	input     <-chan T
	listeners []watchFunc[T]

	m       sync.Mutex
	started atomic.Bool
}

// NewBroadcastWatcher creates a new instance of BroadcastWatcher.
// Parameters:
// - input: The input channel to be watched for incoming messages.
// Returns:
// - A pointer to a newly created BroadcastWatcher instance.
func NewBroadcastWatcher[T any](input <-chan T) *BroadcastWatcher[T] {
	return &BroadcastWatcher[T]{input: input, listeners: make([]watchFunc[T], 0)}
}

// Register registers a new listener and returns a read-only channel for the listener to receive messages.
// This method is thread-safe and can be called concurrently by multiple goroutines. Note that using locks
// can introduce contention and affect performance in highly concurrent environments.
func (w *BroadcastWatcher[T]) Register(f watchFunc[T]) {
	w.m.Lock()
	defer w.m.Unlock()

	w.listeners = append(w.listeners, f)
}

// Watch starts the broadcasting process, sending each message from the input channel to all registered listeners.
// This method is thread-safe and ensures that the broadcasting process can only be started once. Multiple calls to
// this method will have no effect after the first call. Note that using locks can introduce contention and affect
// performance in highly concurrent environments.
// Parameters:
// - ctx: The context to control cancellation of the watching operation.
// Returns:
// - A boolean indicating whether the broadcasting process was successfully started.
func (w *BroadcastWatcher[T]) Watch(ctx context.Context) bool {
	if w.started.CompareAndSwap(false, true) {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-w.input:
					if !ok {
						return
					}
					w.broadcast(ctx, &msg)
				}
			}
		}()
		return true
	} else {
		return false
	}
}

func (w *BroadcastWatcher[T]) broadcast(ctx context.Context, msg *T) {
	w.m.Lock()
	defer w.m.Unlock()

	for _, listener := range w.listeners {
		go func() {
			select {
			case <-ctx.Done():
				return
			default:
				listener(ctx, msg)
			}
		}()
	}
}
