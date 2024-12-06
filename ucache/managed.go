package ucache

import (
	"sync"
	"time"

	"github.com/kordax/basic-utils/uconst"
	"github.com/kordax/basic-utils/uopt"
)

// ManagedCache provides a wrapper around a Cache implementation to manage
// periodic cleanup of outdated cache entries. It uses a background goroutine to perform
// cleanup tasks based on the provided TTL (time-to-live) value.
// The Stop method must be called to clean up resources if you want to stop managing the cache.
type ManagedCache[K any, T any] struct {
	cache    BaseCache[K, T]
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewManagedCache[K any, T any](cache BaseCache[K, T], tick time.Duration) *ManagedCache[K, T] {
	b := &ManagedCache[K, T]{
		cache:    cache,
		stopChan: make(chan struct{}),
	}

	b.wg.Add(1)
	go b.cleanupRoutine(tick)

	return b
}

func (b *ManagedCache[K, T]) cleanupRoutine(tick time.Duration) {
	defer b.wg.Done()
	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			b.ForceCleanup()
		case <-b.stopChan:
			return
		}
	}
}

func (b *ManagedCache[K, T]) ForceCleanup() {
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

func (b *ManagedCache[K, T]) Changes() []K {
	return b.cache.Changes()
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

// ManagedMultiCache provides a wrapper around a MultiCache implementation to manage
// periodic cleanup of outdated cache entries. It uses a background goroutine to perform
// cleanup tasks based on the provided TTL (time-to-live) value.
// The Stop method must be called to clean up resources if you want to stop managing the cache.
type ManagedMultiCache[K CompositeKey, T uconst.Comparable] struct {
	cache    MultiCache[K, T]
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewManagedMultiCache[K CompositeKey, T uconst.Comparable](cache MultiCache[K, T], tick time.Duration) *ManagedMultiCache[K, T] {
	b := &ManagedMultiCache[K, T]{
		cache:    cache,
		stopChan: make(chan struct{}),
	}

	b.wg.Add(1)
	go b.cleanupRoutine(tick)

	return b
}

func (b *ManagedMultiCache[K, T]) cleanupRoutine(tick time.Duration) {
	defer b.wg.Done()
	ticker := time.NewTicker(tick)
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
