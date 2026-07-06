/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ucache

import (
	"sync"
	"sync/atomic"
	"time"

	"git.casinomodule.org/casino27/basic-utils/v3/uconst"
	"git.casinomodule.org/casino27/basic-utils/v3/uopt"
)

type BaseCache[K, T any] interface {
	// Set updates the cache value for the provided key. If the key already exists,
	// its previous value is removed before adding the new value. This method should be thread-safe.
	Set(key K, value T)

	// Get retrieves the value associated with the provided key from the cache.
	// It returns the value and a boolean indicating whether the key was found.
	// This method should be thread-safe. Get operation drops down change state of the item, meaning that item becomes
	// actual after Get operation.
	Get(key K) (*T, bool)

	// GetValue retrieves the value associated with the provided key without allocating a pointer result.
	GetValue(key K) (T, bool)

	// Changes returns a slice of keys that have been modified in the cache.
	// This method provides a way to track changes made to the cache, useful for scenarios like cache syncing.
	// Cache changes will be updated only on modifying operations, but not on Drop() call, meaning that in-fact, changes contain all the present keys.
	Changes() []K

	// Keys returns a snapshot of all currently stored keys.
	// This method is used by managed cleanup and should not depend on change tracking.
	Keys() []K

	// Drop completely clears the cache, removing all entries. This method should be thread-safe.
	Drop()

	// DropKey removes the value associated with the provided key from the cache. This method should be thread-safe.
	// This method also clears up changes associated with this key.
	DropKey(key K)

	// Outdated checks if the provided key or the entire cache (if no key is provided)
	// is outdated based on the set TTL (time-to-live). Returns true if outdated, false otherwise.
	// This method should be thread-safe.
	// If key was not found returns false.
	Outdated(key uopt.Opt[K]) bool

	// SetQuietly is an optimized method adds a value to the cache for the provided key but does so without
	// altering the change history. This method is useful when modifications should not trigger cache change diff.
	// This method should be thread-safe.
	// This operation is much faster and can be used to optimize cache performance in case you don't want to track changes.
	SetQuietly(key K, value T)
}

// The Cache interface defines a set of methods for a generic cache implementation.
// This interface supports setting, getting, and managing cache entries with composite keys.
// Unlike MultiCache, it is designed to handle only one value per key and does not support hierarchical composite keys.
type Cache[K uconst.Unique, T any] interface {
	BaseCache[K, T]
}

// The ComparableCache is the same as Cache, but is more generic and allows comparable keys.
type ComparableCache[K comparable, T any] interface {
	BaseCache[K, T]
}

type ttlConfigurable interface {
	SetTTL(ttl uopt.Opt[time.Duration])
}

type hashValueContainer[K uconst.Unique, T any] struct {
	key   K
	value T
}

// InMemoryHashMapCache provides an in-memory caching mechanism using hashmaps for single-value entries.
// Unlike InMemoryHashMapMultiCache, it stores only one value per key.
// This implementation supports linked-chain collision resolution, so at the worst it should be O(n) complexity.
// This structure translates composite keys into a hash value using a user-provided hashing function.
// Supports optional TTL for entries and ensures concurrency-safe operations using a mutex.
// TTL parameter in cache doesn't automatically clean up all the entries.
// Use ManagedCache wrapper to automatically manage outdated keys.
type InMemoryHashMapCache[K uconst.Unique, T any] struct {
	values  map[int64][]hashValueContainer[K, T]
	changes map[int64][]K

	lastUpdatedKeys map[int64][]keyContainer[K]
	lastUpdated     time.Time
	ttl             *time.Duration

	vMtx sync.RWMutex
}

// NewInMemoryHashMapCache creates a new instance of the InMemoryHashMapCache.
// It takes a hashing function to translate the composite keys to a desired hash type,
// and an optional time-to-live duration for the cache entries.
func NewInMemoryHashMapCache[K uconst.Unique, T any](ttl uopt.Opt[time.Duration]) *InMemoryHashMapCache[K, T] {
	c := &InMemoryHashMapCache[K, T]{
		values:          make(map[int64][]hashValueContainer[K, T]),
		changes:         make(map[int64][]K),
		lastUpdatedKeys: make(map[int64][]keyContainer[K]),
	}
	c.SetTTL(ttl)

	return c
}

func (c *InMemoryHashMapCache[K, T]) SetTTL(ttl uopt.Opt[time.Duration]) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()

	if ttl.Present() {
		t := ttl.OrElse(0)
		c.ttl = &t
		return
	}
	c.ttl = nil
}

// Set updates the cache value for the provided key. If the key already exists,
// its previous value are removed before adding the new value. The operation is thread-safe.
func (c *InMemoryHashMapCache[K, T]) Set(key K, value T) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.put(key, value)
	n := time.Now()
	c.setLastUpdatedKey(key, n)
	c.lastUpdated = n
}

func (c *InMemoryHashMapCache[K, T]) setLastUpdatedKey(key K, updatedAt time.Time) {
	hash := key.Key()
	keys := c.lastUpdatedKeys[hash]
	for i := range keys {
		if keys[i].key.Equals(key) {
			keys[i].updatedAt = updatedAt
			keys[i].key = key
			c.lastUpdatedKeys[hash] = keys
			return
		}
	}
	c.lastUpdatedKeys[hash] = append(keys, keyContainer[K]{
		key:       key,
		updatedAt: updatedAt,
	})
}

// SetQuietly is an optimized method that adds value to the cache for the provided key but does so without
// altering the change history. This operation can be used when modifications should not trigger cache change diff.
// This operation is much faster and can be used to optimize cache performance in case you don't want to track changes.
func (c *InMemoryHashMapCache[K, T]) SetQuietly(key K, value T) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.addTran(key, value)
	n := time.Now()
	c.setLastUpdatedKey(key, n)
	c.lastUpdated = n
}

// Get retrieves the value associated with the provided key from the cache.
// The operation is thread-safe and does not alter the change history.
func (c *InMemoryHashMapCache[K, T]) Get(key K) (*T, bool) {
	value, ok := c.GetValue(key)
	if !ok {
		return nil, false
	}
	return &value, true
}

func (c *InMemoryHashMapCache[K, T]) GetValue(key K) (T, bool) {
	c.vMtx.RLock()
	if c.ttl == nil {
		defer c.vMtx.RUnlock()
		return c.getUnsafe(key)
	}
	c.vMtx.RUnlock()

	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	if c.outdatedKeyUnsafe(key) {
		c.dropKey(key)
		c.removeChangedKey(key)
		c.removeLastUpdatedKey(key)
		c.compactIfEmpty()
		var zero T
		return zero, false
	}

	return c.getUnsafe(key)
}

func (c *InMemoryHashMapCache[K, T]) getUnsafe(key K) (T, bool) {
	values, ok := c.values[key.Key()]
	if !ok {
		var zero T
		return zero, false
	}

	for i := range values {
		if values[i].key.Equals(key) {
			return values[i].value, true
		}
	}

	var zero T
	return zero, false
}

func (c *InMemoryHashMapCache[K, T]) outdatedKeyUnsafe(key K) bool {
	if c.ttl == nil {
		return false
	}
	for _, lu := range c.lastUpdatedKeys[key.Key()] {
		if lu.key.Equals(key) {
			return time.Since(lu.updatedAt) > *c.ttl
		}
	}

	return true
}

// Changes returns a slice of keys that have been modified in the cache.
// This method provides a way to track changes made to the cache, useful for scenarios like cache syncing.
func (c *InMemoryHashMapCache[K, T]) Changes() []K {
	c.vMtx.RLock()
	defer c.vMtx.RUnlock()

	result := make([]K, 0, len(c.changes))
	for _, keys := range c.changes {
		result = append(result, keys...)
	}

	return result
}

func (c *InMemoryHashMapCache[K, T]) Keys() []K {
	c.vMtx.RLock()
	defer c.vMtx.RUnlock()

	result := make([]K, 0, len(c.values))
	for _, values := range c.values {
		for i := range values {
			result = append(result, values[i].key)
		}
	}

	return result
}

// Drop completely clears the cache, removing all entries. The operation is thread-safe.
func (c *InMemoryHashMapCache[K, T]) Drop() {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.dropAll()
	c.changes = make(map[int64][]K)
	c.lastUpdatedKeys = make(map[int64][]keyContainer[K])
}

// DropKey removes the value associated with the provided key from the cache. The operation is thread-safe.
func (c *InMemoryHashMapCache[K, T]) DropKey(key K) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.dropKey(key)
	c.removeChangedKey(key)
	c.removeLastUpdatedKey(key)
	c.compactIfEmpty()
}

// Outdated checks if the provided key or the entire cache (if no key is provided)
// is outdated based on the set TTL. Returns true if outdated, false otherwise.
func (c *InMemoryHashMapCache[K, T]) Outdated(key uopt.Opt[K]) bool {
	c.vMtx.RLock()
	defer c.vMtx.RUnlock()

	if c.ttl == nil {
		return false
	}

	if key.Present() {
		k := key.Get()
		return c.outdatedKeyUnsafe(*k)
	}

	return false
}

func (c *InMemoryHashMapCache[K, T]) RemoveOutdated() int {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()

	if c.ttl == nil {
		return 0
	}

	var expired []K
	for _, keys := range c.lastUpdatedKeys {
		for _, lu := range keys {
			if time.Since(lu.updatedAt) > *c.ttl {
				expired = append(expired, lu.key)
			}
		}
	}
	for _, key := range expired {
		c.dropKey(key)
		c.removeChangedKey(key)
		c.removeLastUpdatedKey(key)
	}
	c.compactIfEmpty()

	return len(expired)
}

func (c *InMemoryHashMapCache[K, T]) dropAll() {
	c.values = make(map[int64][]hashValueContainer[K, T])
}

func (c *InMemoryHashMapCache[K, T]) compactIfEmpty() {
	if len(c.values) != 0 {
		return
	}
	c.values = make(map[int64][]hashValueContainer[K, T])
	c.changes = make(map[int64][]K)
	c.lastUpdatedKeys = make(map[int64][]keyContainer[K])
}

func (c *InMemoryHashMapCache[K, T]) put(key K, value T) {
	hash := c.addTran(key, value)
	c.addChangedKey(hash, key)
}

func (c *InMemoryHashMapCache[K, T]) addChangedKey(hash int64, key K) {
	keys := c.changes[hash]
	for _, changed := range keys {
		if changed.Equals(key) {
			return
		}
	}
	c.changes[hash] = append(keys, key)
}

func (c *InMemoryHashMapCache[K, T]) removeChangedKey(key K) {
	hash := key.Key()
	keys := c.changes[hash]
	for i := range keys {
		if keys[i].Equals(key) {
			var zero K
			keys[i] = zero
			keys = append(keys[:i], keys[i+1:]...)
			break
		}
	}
	if len(keys) == 0 {
		delete(c.changes, hash)
		return
	}
	c.changes[hash] = keys
}

func (c *InMemoryHashMapCache[K, T]) removeLastUpdatedKey(key K) {
	hash := key.Key()
	keys := c.lastUpdatedKeys[hash]
	for i := range keys {
		if keys[i].key.Equals(key) {
			var zero keyContainer[K]
			keys[i] = zero
			keys = append(keys[:i], keys[i+1:]...)
			break
		}
	}
	if len(keys) == 0 {
		delete(c.lastUpdatedKeys, hash)
		return
	}
	c.lastUpdatedKeys[hash] = keys
}

func (c *InMemoryHashMapCache[K, T]) addTran(key K, value T) int64 {
	keyHash := key.Key()
	values := c.values[keyHash]
	if len(values) == 0 {
		values = make([]hashValueContainer[K, T], 0)
		values = append(values, hashValueContainer[K, T]{
			key:   key,
			value: value,
		})
		c.values[keyHash] = values
	} else {
		for i, v := range values {
			if v.key.Equals(key) {
				values[i] = hashValueContainer[K, T]{
					key:   key,
					value: value,
				}
				c.values[keyHash] = values
				return keyHash
			}
		}
		values = append(values, hashValueContainer[K, T]{
			key:   key,
			value: value,
		})
		c.values[keyHash] = values
	}

	return keyHash
}

func (c *InMemoryHashMapCache[K, T]) dropKey(key K) {
	hash := key.Key()
	values := c.values[hash]
	for i := range values {
		if values[i].key.Equals(key) {
			var zero hashValueContainer[K, T]
			values[i] = zero
			values = append(values[:i], values[i+1:]...)
			break
		}
	}
	if len(values) == 0 {
		delete(c.values, hash)
		return
	}
	c.values[hash] = values
}

// InMemoryComparableMapCache provides an in-memory caching mechanism for single-value entries.
// It supports optional TTL for entries and ensures concurrency-safe operations.
// It is very similiar to InMemoryHashMapCache by behaviour, and the only difference is a key type constraint.
type InMemoryComparableMapCache[K comparable, T any] struct {
	values          sync.Map
	changes         sync.Map
	lastUpdatedKeys sync.Map
	ttlNanos        atomic.Int64

	bufferOnce      sync.Once
	bufferStarted   atomic.Bool
	bufferClosed    atomic.Bool
	bufferQueue     chan bufferedComparableMapEntry[K, T]
	bufferQueueSize int
	bufferWorkers   int
	bufferedMaxKeys int64
	bufferedKeys    sync.Map
	bufferedKeyLen  atomic.Int64
	bufferedSeq     atomic.Uint64
	bufferedLatest  sync.Map
	bufferWG        sync.WaitGroup
	bufferWorkersWG sync.WaitGroup
	bufferAddMtx    sync.Mutex
}

const (
	defaultComparableMapCacheBufferedWorkers   = 4
	defaultComparableMapCacheBufferedQueueSize = 65536
)

type InMemoryComparableMapCacheOptions struct {
	// TTL configures cache entry time-to-live.
	TTL uopt.Opt[time.Duration]
	// BufferedWorkers configures how many workers apply buffered writes.
	BufferedWorkers int
	// BufferedQueueSize configures the buffered write queue capacity.
	BufferedQueueSize int
	// BufferedMaxKeys limits keys accepted through buffered writes. Zero means unlimited.
	// Synchronous Set and SetQuietly keep exact write semantics and are not rejected by this limit.
	BufferedMaxKeys int64
}

type bufferedComparableMapEntry[K comparable, T any] struct {
	key          K
	value        T
	seq          uint64
	trackChanges bool
}

// NewInMemoryComparableMapCache creates a new instance of InMemoryComparableMapCache.
// It accepts an optional TTL (time-to-live) duration for cache entries.
func NewInMemoryComparableMapCache[K comparable, T any](ttl uopt.Opt[time.Duration]) *InMemoryComparableMapCache[K, T] {
	return NewInMemoryComparableMapCacheWithOptions[K, T](InMemoryComparableMapCacheOptions{
		TTL: ttl,
	})
}

// NewInMemoryComparableMapCacheWithOptions creates a new InMemoryComparableMapCache with TTL and buffered write options.
func NewInMemoryComparableMapCacheWithOptions[K comparable, T any](
	options InMemoryComparableMapCacheOptions,
) *InMemoryComparableMapCache[K, T] {
	c := &InMemoryComparableMapCache[K, T]{
		bufferQueueSize: normalizeComparableMapCacheBufferedQueueSize(options.BufferedQueueSize),
		bufferWorkers:   normalizeComparableMapCacheBufferedWorkers(options.BufferedWorkers),
		bufferedMaxKeys: normalizeComparableMapCacheBufferedMaxKeys(options.BufferedMaxKeys),
	}
	c.SetTTL(options.TTL)
	return c
}

func normalizeComparableMapCacheBufferedWorkers(workers int) int {
	if workers > 0 {
		return workers
	}

	return defaultComparableMapCacheBufferedWorkers
}

func normalizeComparableMapCacheBufferedQueueSize(queueSize int) int {
	if queueSize > 0 {
		return queueSize
	}

	return defaultComparableMapCacheBufferedQueueSize
}

func normalizeComparableMapCacheBufferedMaxKeys(maxKeys int64) int64 {
	if maxKeys > 0 {
		return maxKeys
	}

	return 0
}

func (c *InMemoryComparableMapCache[K, T]) SetTTL(ttl uopt.Opt[time.Duration]) {
	if ttl.Present() {
		t := ttl.OrElse(0)
		if t > 0 {
			c.ttlNanos.Store(int64(t))
			c.touchExisting(time.Now().UnixNano())
			return
		}
	}
	c.ttlNanos.Store(0)
	c.lastUpdatedKeys.Range(func(key, _ any) bool {
		c.lastUpdatedKeys.Delete(key)
		return true
	})
}

func (c *InMemoryComparableMapCache[K, T]) ttl() time.Duration {
	nanos := c.ttlNanos.Load()
	if nanos <= 0 {
		return 0
	}

	return time.Duration(nanos)
}

func (c *InMemoryComparableMapCache[K, T]) touch(key K) {
	if c.ttlNanos.Load() > 0 {
		c.lastUpdatedKeys.Store(key, time.Now().UnixNano())
	}
}

func (c *InMemoryComparableMapCache[K, T]) touchExisting(now int64) {
	c.values.Range(func(key, _ any) bool {
		c.lastUpdatedKeys.Store(key, now)
		return true
	})
}

func loadSyncMapValue[K comparable, T any](m *sync.Map, key K) (T, bool) {
	raw, ok := m.Load(key)
	if !ok {
		var zero T
		return zero, false
	}

	value, ok := raw.(T)
	if !ok {
		var zero T
		return zero, false
	}

	return value, true
}

func rangeSyncMapKeys[K comparable](m *sync.Map) []K {
	result := make([]K, 0)
	m.Range(func(key, _ any) bool {
		typedKey, ok := key.(K)
		if ok {
			result = append(result, typedKey)
		}
		return true
	})

	return result
}

// Set updates the cache value for the provided key. If the key already exists,
// its previous value is replaced with the new value. The operation is thread-safe.
func (c *InMemoryComparableMapCache[K, T]) Set(key K, value T) {
	c.values.Store(key, value)
	c.changes.Store(key, struct{}{})
	c.touch(key)
}

func (c *InMemoryComparableMapCache[K, T]) setBuffered(key K, value T) bool {
	return c.enqueueBuffered(key, value, true)
}

// SetQuietly adds a value to the cache for the provided key without altering the change history.
// This method is thread-safe and optimized for performance when change tracking is unnecessary.
func (c *InMemoryComparableMapCache[K, T]) SetQuietly(key K, value T) {
	c.values.Store(key, value)
	c.touch(key)
}

func (c *InMemoryComparableMapCache[K, T]) setQuietlyBuffered(key K, value T) bool {
	return c.enqueueBuffered(key, value, false)
}

func (c *InMemoryComparableMapCache[K, T]) enqueueBuffered(key K, value T, trackChanges bool) bool {
	c.bufferAddMtx.Lock()
	defer c.bufferAddMtx.Unlock()

	if c.bufferClosed.Load() {
		return false
	}

	if !c.admitBufferedKey(key) {
		return false
	}

	c.startBuffer()

	seq := c.bufferedSeq.Add(1)
	c.bufferedLatest.Store(key, seq)
	c.bufferWG.Add(1)
	c.bufferQueue <- bufferedComparableMapEntry[K, T]{
		key:          key,
		value:        value,
		seq:          seq,
		trackChanges: trackChanges,
	}

	return true
}

func (c *InMemoryComparableMapCache[K, T]) admitBufferedKey(key K) bool {
	if c.bufferedMaxKeys <= 0 {
		return true
	}

	if _, exists := c.bufferedKeys.Load(key); exists {
		return true
	}

	for {
		current := c.bufferedKeyLen.Load()
		if current >= c.bufferedMaxKeys {
			return false
		}
		if c.bufferedKeyLen.CompareAndSwap(current, current+1) {
			if _, loaded := c.bufferedKeys.LoadOrStore(key, struct{}{}); loaded {
				c.bufferedKeyLen.Add(-1)
			}
			return true
		}
	}
}

func (c *InMemoryComparableMapCache[K, T]) releaseBufferedKey(key K) {
	if c.bufferedMaxKeys <= 0 {
		return
	}
	if _, loaded := c.bufferedKeys.LoadAndDelete(key); loaded {
		c.bufferedKeyLen.Add(-1)
	}
}

func (c *InMemoryComparableMapCache[K, T]) clearBufferedKeys() {
	c.bufferedLatest.Range(func(key, _ any) bool {
		c.bufferedLatest.Delete(key)
		return true
	})

	if c.bufferedMaxKeys <= 0 {
		return
	}
	c.bufferedKeys.Range(func(key, _ any) bool {
		c.bufferedKeys.Delete(key)
		return true
	})
	c.bufferedKeyLen.Store(0)
}

func (c *InMemoryComparableMapCache[K, T]) startBuffer() {
	c.bufferOnce.Do(func() {
		c.bufferQueueSize = normalizeComparableMapCacheBufferedQueueSize(c.bufferQueueSize)
		c.bufferWorkers = normalizeComparableMapCacheBufferedWorkers(c.bufferWorkers)
		c.bufferQueue = make(chan bufferedComparableMapEntry[K, T], c.bufferQueueSize)

		for i := 0; i < c.bufferWorkers; i++ {
			c.bufferWorkersWG.Add(1)
			go c.bufferWorker()
		}

		c.bufferStarted.Store(true)
	})
}

func (c *InMemoryComparableMapCache[K, T]) bufferWorker() {
	defer c.bufferWorkersWG.Done()

	for entry := range c.bufferQueue {
		if !c.latestBufferedEntry(entry) {
			c.bufferWG.Done()
			continue
		}

		c.values.Store(entry.key, entry.value)
		if entry.trackChanges {
			c.changes.Store(entry.key, struct{}{})
		}
		c.touch(entry.key)
		c.bufferWG.Done()
	}
}

func (c *InMemoryComparableMapCache[K, T]) latestBufferedEntry(entry bufferedComparableMapEntry[K, T]) bool {
	raw, exists := c.bufferedLatest.Load(entry.key)
	if !exists {
		return false
	}
	seq, ok := raw.(uint64)
	return ok && seq == entry.seq
}

func (c *InMemoryComparableMapCache[K, T]) waitBuffered() {
	if !c.bufferStarted.Load() {
		return
	}

	c.bufferAddMtx.Lock()
	c.bufferWG.Wait()
	c.bufferAddMtx.Unlock()
}

func (c *InMemoryComparableMapCache[K, T]) closeBuffered() {
	c.bufferAddMtx.Lock()
	if !c.bufferClosed.Swap(true) {
		if c.bufferStarted.Load() {
			c.bufferWG.Wait()
			close(c.bufferQueue)
		}
	}
	c.bufferAddMtx.Unlock()

	c.bufferWorkersWG.Wait()
}

// Get retrieves the value associated with the provided key from the cache.
// It returns a pointer to the value and a boolean indicating whether the key was found.
// The operation is thread-safe.
func (c *InMemoryComparableMapCache[K, T]) Get(key K) (*T, bool) {
	value, ok := c.GetValue(key)
	if !ok {
		return nil, false
	}
	return &value, true
}

func (c *InMemoryComparableMapCache[K, T]) GetValue(key K) (T, bool) {
	if c.ttlNanos.Load() > 0 && c.outdatedKey(key) {
		c.deleteKey(key)
		var zero T
		return zero, false
	}

	return loadSyncMapValue[K, T](&c.values, key)
}

func (c *InMemoryComparableMapCache[K, T]) outdatedKey(key K) bool {
	ttl := c.ttl()
	if ttl <= 0 {
		return false
	}
	raw, exists := c.lastUpdatedKeys.Load(key)
	if !exists {
		return true
	}
	lastUpdated, ok := raw.(int64)
	if !ok {
		return true
	}

	return time.Since(time.Unix(0, lastUpdated)) > ttl
}

// Changes returns a slice of keys that have been modified in the cache since the last call to Changes.
// This method is thread-safe.
func (c *InMemoryComparableMapCache[K, T]) Changes() []K {
	return rangeSyncMapKeys[K](&c.changes)
}

func (c *InMemoryComparableMapCache[K, T]) Keys() []K {
	return rangeSyncMapKeys[K](&c.values)
}

// Drop completely clears the cache, removing all entries. The operation is thread-safe.
func (c *InMemoryComparableMapCache[K, T]) Drop() {
	c.waitBuffered()
	c.clearBufferedKeys()

	c.values.Range(func(key, _ any) bool {
		c.values.Delete(key)
		return true
	})
	c.changes.Range(func(key, _ any) bool {
		c.changes.Delete(key)
		return true
	})
	c.lastUpdatedKeys.Range(func(key, _ any) bool {
		c.lastUpdatedKeys.Delete(key)
		return true
	})
}

// DropKey removes the value associated with the provided key from the cache.
// The operation is thread-safe.
func (c *InMemoryComparableMapCache[K, T]) DropKey(key K) {
	c.waitBuffered()
	c.deleteKey(key)
}

func (c *InMemoryComparableMapCache[K, T]) deleteKey(key K) {
	c.values.Delete(key)
	c.changes.Delete(key)
	c.lastUpdatedKeys.Delete(key)
	c.bufferedLatest.Delete(key)
	c.releaseBufferedKey(key)
}

// Outdated checks if the provided key is outdated based on the set TTL (time-to-live).
// Returns true if outdated, false otherwise.
// If no TTL is set, it returns false. If TTL is set and the key does not exist, it returns true.
func (c *InMemoryComparableMapCache[K, T]) Outdated(key uopt.Opt[K]) bool {
	if c.ttlNanos.Load() <= 0 {
		return false
	}

	if k := key.Get(); k != nil {
		return c.outdatedKey(*k)
	}

	return false
}

func (c *InMemoryComparableMapCache[K, T]) RemoveOutdated() int {
	ttl := c.ttl()
	if ttl <= 0 {
		return 0
	}

	removed := 0
	c.lastUpdatedKeys.Range(func(key, value any) bool {
		lastUpdated, ok := value.(int64)
		if !ok || time.Since(time.Unix(0, lastUpdated)) > ttl {
			c.values.Delete(key)
			c.changes.Delete(key)
			c.lastUpdatedKeys.Delete(key)
			c.bufferedLatest.Delete(key)
			if typedKey, typed := key.(K); typed {
				c.releaseBufferedKey(typedKey)
			}
			removed++
		}
		return true
	})

	return removed
}

// InMemoryBufferedComparableMapCache wraps InMemoryComparableMapCache and uses buffered writes by default.
// Set and SetQuietly enqueue writes and return immediately. Call Wait to make accepted writes visible deterministically.
// If BufferedMaxKeys is configured, new keys above the limit can be rejected silently by Set and SetQuietly.
type InMemoryBufferedComparableMapCache[K comparable, T any] struct {
	cache *InMemoryComparableMapCache[K, T]
}

// NewInMemoryBufferedComparableMapCache creates a buffered comparable-key cache with default buffered write settings.
func NewInMemoryBufferedComparableMapCache[K comparable, T any](
	ttl uopt.Opt[time.Duration],
) *InMemoryBufferedComparableMapCache[K, T] {
	return NewInMemoryBufferedComparableMapCacheWithOptions[K, T](InMemoryComparableMapCacheOptions{
		TTL: ttl,
	})
}

// NewInMemoryBufferedComparableMapCacheWithOptions creates a buffered comparable-key cache with custom settings.
func NewInMemoryBufferedComparableMapCacheWithOptions[K comparable, T any](
	options InMemoryComparableMapCacheOptions,
) *InMemoryBufferedComparableMapCache[K, T] {
	return &InMemoryBufferedComparableMapCache[K, T]{
		cache: NewInMemoryComparableMapCacheWithOptions[K, T](options),
	}
}

func (c *InMemoryBufferedComparableMapCache[K, T]) SetTTL(ttl uopt.Opt[time.Duration]) {
	c.cache.SetTTL(ttl)
}

func (c *InMemoryBufferedComparableMapCache[K, T]) Set(key K, value T) {
	c.cache.setBuffered(key, value)
}

func (c *InMemoryBufferedComparableMapCache[K, T]) SetQuietly(key K, value T) {
	c.cache.setQuietlyBuffered(key, value)
}

// Wait blocks until all currently accepted buffered updates are applied.
func (c *InMemoryBufferedComparableMapCache[K, T]) Wait() {
	c.cache.waitBuffered()
}

// CloseBuffered waits for buffered updates and stops background buffer workers.
// After CloseBuffered, Set and SetQuietly stop accepting buffered writes.
func (c *InMemoryBufferedComparableMapCache[K, T]) CloseBuffered() {
	c.cache.closeBuffered()
}

func (c *InMemoryBufferedComparableMapCache[K, T]) Get(key K) (*T, bool) {
	return c.cache.Get(key)
}

func (c *InMemoryBufferedComparableMapCache[K, T]) GetValue(key K) (T, bool) {
	return c.cache.GetValue(key)
}

func (c *InMemoryBufferedComparableMapCache[K, T]) Changes() []K {
	return c.cache.Changes()
}

func (c *InMemoryBufferedComparableMapCache[K, T]) Keys() []K {
	return c.cache.Keys()
}

func (c *InMemoryBufferedComparableMapCache[K, T]) Drop() {
	c.cache.Drop()
}

func (c *InMemoryBufferedComparableMapCache[K, T]) DropKey(key K) {
	c.cache.DropKey(key)
}

func (c *InMemoryBufferedComparableMapCache[K, T]) Outdated(key uopt.Opt[K]) bool {
	return c.cache.Outdated(key)
}

func (c *InMemoryBufferedComparableMapCache[K, T]) RemoveOutdated() int {
	return c.cache.RemoveOutdated()
}
