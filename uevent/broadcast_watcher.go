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
// - listeners: Registered callbacks indexed by subscription ID.
// - m: A mutex to ensure that Register and other operations are thread-safe.
// - started: An atomic boolean to ensure the Start method is only called once.
type BroadcastWatcher[T any] struct {
	input     <-chan T
	listeners map[uint64]watchFunc[T]

	m              sync.RWMutex
	nextListenerID atomic.Uint64
	started        atomic.Bool
}

// NewBroadcastWatcher creates a new instance of BroadcastWatcher.
// Parameters:
// - input: The input channel to be watched for incoming messages.
// Returns:
// - A pointer to a newly created BroadcastWatcher instance.
func NewBroadcastWatcher[T any](input <-chan T) *BroadcastWatcher[T] {
	return &BroadcastWatcher[T]{input: input, listeners: make(map[uint64]watchFunc[T])}
}

type broadcastSubscription struct {
	once        sync.Once
	unsubscribe func() bool
}

func (s *broadcastSubscription) Unsubscribe() bool {
	removed := false
	s.once.Do(func() {
		removed = s.unsubscribe()
	})

	return removed
}

// Subscribe registers a listener and returns an idempotent handle that can remove it.
// A callback already dispatched when Unsubscribe is called may still finish.
func (w *BroadcastWatcher[T]) Subscribe(f watchFunc[T]) Subscription {
	id := w.nextListenerID.Add(1)

	w.m.Lock()
	if w.listeners == nil {
		w.listeners = make(map[uint64]watchFunc[T])
	}
	w.listeners[id] = f
	w.m.Unlock()

	return &broadcastSubscription{unsubscribe: func() bool {
		w.m.Lock()
		defer w.m.Unlock()

		if _, ok := w.listeners[id]; !ok {
			return false
		}

		delete(w.listeners, id)
		return true
	}}
}

// Register registers a listener and keeps it subscribed for the lifetime of the watcher.
func (w *BroadcastWatcher[T]) Register(f watchFunc[T]) {
	w.Subscribe(f)
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
					w.broadcast(ctx, msg)
				}
			}
		}()
		return true
	} else {
		return false
	}
}

func (w *BroadcastWatcher[T]) broadcast(ctx context.Context, msg T) {
	w.m.RLock()
	listeners := make([]watchFunc[T], 0, len(w.listeners))
	for _, listener := range w.listeners {
		listeners = append(listeners, listener)
	}
	w.m.RUnlock()

	for _, listener := range listeners {
		go func(listener watchFunc[T]) {
			select {
			case <-ctx.Done():
				return
			default:
				listener(ctx, msg)
			}
		}(listener)
	}
}
