package ucache

import (
	"sync"
	"time"

	"git.casinomodule.org/casino27/basic-utils/v2/uconst"
	"git.casinomodule.org/casino27/basic-utils/v2/uopt"
)

// ManagedCache provides a wrapper around a Cache implementation to manage
// periodic cleanup of outdated cache entries. It uses a background goroutine to perform
// cleanup tasks based on the provided TTL (time-to-live) value.
// The Stop method must be called to clean up resources if you want to stop managing the cache.
type ManagedCache[K any, T any] struct {
	cache    BaseCache[K, T]
	stopChan chan struct{}
	wg       sync.WaitGroup
	stopOnce sync.Once
}

type ManagedCacheOptions struct {
	CleanupInterval time.Duration
	TTL             uopt.Opt[time.Duration]
}

type outdatedCleaner interface {
	RemoveOutdated() int
}

func NewManagedCache[K any, T any](cache BaseCache[K, T], tick time.Duration) *ManagedCache[K, T] {
	return NewManagedCacheWithOptions[K, T](cache, ManagedCacheOptions{CleanupInterval: tick})
}

func NewManagedCacheWithOptions[K any, T any](cache BaseCache[K, T], options ManagedCacheOptions) *ManagedCache[K, T] {
	if options.TTL.Present() {
		if configurable, ok := any(cache).(ttlConfigurable); ok {
			configurable.SetTTL(options.TTL)
		}
	}
	tick := options.CleanupInterval
	if tick <= 0 {
		tick = time.Minute
	}

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
	if cleaner, ok := any(b.cache).(outdatedCleaner); ok {
		cleaner.RemoveOutdated()
		return
	}
	for _, key := range b.cache.Keys() {
		if b.cache.Outdated(uopt.Of(key)) {
			b.cache.DropKey(key)
		}
	}
}

func (b *ManagedCache[K, T]) Stop() {
	b.stopOnce.Do(func() {
		close(b.stopChan)
	})
	b.wg.Wait()
}

func (b *ManagedCache[K, T]) Set(key K, value T) {
	b.cache.Set(key, value)
}

func (b *ManagedCache[K, T]) Get(key K) (*T, bool) {
	return b.cache.Get(key)
}

func (b *ManagedCache[K, T]) GetValue(key K) (T, bool) {
	return b.cache.GetValue(key)
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
	stopOnce sync.Once
}

func NewManagedMultiCache[K CompositeKey, T uconst.Comparable](cache MultiCache[K, T], tick time.Duration) *ManagedMultiCache[K, T] {
	return NewManagedMultiCacheWithOptions[K, T](cache, ManagedCacheOptions{CleanupInterval: tick})
}

func NewManagedMultiCacheWithOptions[K CompositeKey, T uconst.Comparable](cache MultiCache[K, T], options ManagedCacheOptions) *ManagedMultiCache[K, T] {
	if options.TTL.Present() {
		if configurable, ok := any(cache).(ttlConfigurable); ok {
			configurable.SetTTL(options.TTL)
		}
	}
	tick := options.CleanupInterval
	if tick <= 0 {
		tick = time.Minute
	}

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
			b.ForceCleanup()
		case <-b.stopChan:
			return
		}
	}
}

func (b *ManagedMultiCache[K, T]) ForceCleanup() {
	if cleaner, ok := any(b.cache).(outdatedCleaner); ok {
		cleaner.RemoveOutdated()
		return
	}
	for _, key := range b.cache.Keys() {
		if b.cache.Outdated(uopt.Of(key)) {
			b.cache.DropKey(key)
		}
	}
}

func (b *ManagedMultiCache[K, T]) Stop() {
	b.stopOnce.Do(func() {
		close(b.stopChan)
	})
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
