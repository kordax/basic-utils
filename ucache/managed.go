package ucache

import (
	"sync"
	"time"

	"github.com/kordax/basic-utils/uopt"
)

// ManagedMultiCache provides a wrapper around a MultiCache implementation to manage
// periodic cleanup of outdated cache entries. It uses a background goroutine to perform
// cleanup tasks based on the provided TTL (time-to-live) value.
// The Stop method must be called to clean up resources if you want to stop managing the cache.
type ManagedMultiCache[K CompositeKey, T Comparable] struct {
	cache    MultiCache[K, T]
	ttl      time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewManagedMultiCache[K CompositeKey, T Comparable](cache MultiCache[K, T], ttl time.Duration) *ManagedMultiCache[K, T] {
	b := &ManagedMultiCache[K, T]{
		cache:    cache,
		stopChan: make(chan struct{}),
		ttl:      ttl,
	}

	b.wg.Add(1)
	go b.cleanupRoutine()

	return b
}

func (b *ManagedMultiCache[K, T]) cleanupRoutine() {
	defer b.wg.Done()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			b.performCleanup()
		case <-b.stopChan:
			return
		}
	}
}

func (b *ManagedMultiCache[K, T]) performCleanup() {
	for _, key := range b.cache.Changes() {
		if b.cache.Outdated(uopt.Of(key)) {
			b.cache.DropKey(key)
		}
	}
}

func (b *ManagedMultiCache[K, T]) Stop() {
	close(b.stopChan)
	b.wg.Wait()
}

func (b *ManagedMultiCache[K, T]) Put(key K, values ...T) {
	b.cache.Put(key, values...)
}

func (b *ManagedMultiCache[K, T]) Set(key K, values ...T) {
	b.cache.Set(key, values...)
}

func (b *ManagedMultiCache[K, T]) Get(key K) []T {
	return b.cache.Get(key)
}

func (b *ManagedMultiCache[K, T]) Changes() []K {
	return b.cache.Changes()
}

func (b *ManagedMultiCache[K, T]) Drop() {
	b.cache.Drop()
}

func (b *ManagedMultiCache[K, T]) DropKey(key K) {
	b.cache.DropKey(key)
}

func (b *ManagedMultiCache[K, T]) Outdated(key uopt.Opt[K]) bool {
	return b.cache.Outdated(key)
}

func (b *ManagedMultiCache[K, T]) PutQuietly(key K, values ...T) {
	b.cache.PutQuietly(key, values...)
}

// ManagedCache provides a wrapper around a Cache implementation to manage
// periodic cleanup of outdated cache entries. It uses a background goroutine to perform
// cleanup tasks based on the provided TTL (time-to-live) value.
// The Stop method must be called to clean up resources if you want to stop managing the cache.
type ManagedCache[K Unique, T any] struct {
	cache    Cache[K, T]
	ttl      *time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewManagedCache[K Unique, T any](cache Cache[K, T], ttl uopt.Opt[time.Duration]) *ManagedCache[K, T] {
	b := &ManagedCache[K, T]{
		cache:    cache,
		stopChan: make(chan struct{}),
	}

	ttl.IfPresent(func(t time.Duration) {
		b.ttl = &t
	})

	b.wg.Add(1)
	go b.cleanupRoutine()

	return b
}

func (b *ManagedCache[K, T]) cleanupRoutine() {
	defer b.wg.Done()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			b.performCleanup()
		case <-b.stopChan:
			return
		}
	}
}

func (b *ManagedCache[K, T]) performCleanup() {
	for _, key := range b.cache.Changes() {
		if b.cache.Outdated(uopt.Of(key)) {
			b.cache.DropKey(key)
		}
	}
}

func (b *ManagedCache[K, T]) Stop() {
	close(b.stopChan)
	b.wg.Wait()
}

func (b *ManagedCache[K, T]) Set(key K, value T) {
	b.cache.Set(key, value)
}

func (b *ManagedCache[K, T]) Get(key K) (*T, bool) {
	return b.cache.Get(key)
}

func (b *ManagedCache[K, T]) Drop() {
	b.cache.Drop()
}

func (b *ManagedCache[K, T]) DropKey(key K) {
	b.cache.DropKey(key)
}

func (b *ManagedCache[K, T]) Outdated(key uopt.Opt[K]) bool {
	return b.cache.Outdated(key)
}

func (b *ManagedCache[K, T]) SetQuietly(key K, value T) {
	b.cache.SetQuietly(key, value)
}
