/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ucache

import (
	"sync"
	"time"

	"github.com/kordax/basic-utils/v2/uconst"
	"github.com/kordax/basic-utils/v2/umap"
	"github.com/kordax/basic-utils/v2/uopt"
	"github.com/kordax/basic-utils/v2/uset"
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

	// Changes returns a slice of keys that have been modified in the cache.
	// This method provides a way to track changes made to the cache, useful for scenarios like cache syncing.
	// Cache changes will be updated only on modifying operations, but not on Drop() call, meaning that in-fact, changes contain all the present keys.
	Changes() []K

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
	changes map[int64]K

	lastUpdatedKeys map[int64]keyContainer[K]
	lastUpdated     time.Time
	ttl             *time.Duration

	vMtx sync.Mutex
}

// NewInMemoryHashMapCache creates a new instance of the InMemoryHashMapCache.
// It takes a hashing function to translate the composite keys to a desired hash type,
// and an optional time-to-live duration for the cache entries.
func NewInMemoryHashMapCache[K uconst.Unique, T any](ttl uopt.Opt[time.Duration]) *InMemoryHashMapCache[K, T] {
	c := &InMemoryHashMapCache[K, T]{
		values:          make(map[int64][]hashValueContainer[K, T]),
		changes:         make(map[int64]K),
		lastUpdatedKeys: make(map[int64]keyContainer[K]),
	}
	ttl.IfPresent(func(t time.Duration) {
		c.ttl = &t
	})

	return c
}

// Set updates the cache value for the provided key. If the key already exists,
// its previous value are removed before adding the new value. The operation is thread-safe.
func (c *InMemoryHashMapCache[K, T]) Set(key K, value T) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.put(key, value)
	n := time.Now()
	c.lastUpdatedKeys[key.Key()] = keyContainer[K]{
		key:       key,
		updatedAt: n,
	}
	c.lastUpdated = n
}

// SetQuietly is an optimized method that adds value to the cache for the provided key but does so without
// altering the change history. This operation can be used when modifications should not trigger cache change diff.
// This operation is much faster and can be used to optimize cache performance in case you don't want to track changes.
func (c *InMemoryHashMapCache[K, T]) SetQuietly(key K, value T) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.addTran(key, value)
	n := time.Now()
	c.lastUpdatedKeys[key.Key()] = keyContainer[K]{
		key:       key,
		updatedAt: n,
	}
	c.lastUpdated = n
}

// Get retrieves the value associated with the provided key from the cache.
// The operation is thread-safe and does not alter the change history.
func (c *InMemoryHashMapCache[K, T]) Get(key K) (*T, bool) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()

	values, ok := c.values[key.Key()]
	if !ok {
		return nil, false
	}

	if len(values) > 0 {
		for _, v := range values {
			if v.key.Equals(key) {
				return &v.value, true
			}
		}

		return nil, false
	}

	return &values[0].value, ok
}

// Changes returns a slice of keys that have been modified in the cache.
// This method provides a way to track changes made to the cache, useful for scenarios like cache syncing.
func (c *InMemoryHashMapCache[K, T]) Changes() []K {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	return umap.Values(c.changes)
}

// Drop completely clears the cache, removing all entries. The operation is thread-safe.
func (c *InMemoryHashMapCache[K, T]) Drop() {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.dropAll()
	c.changes = nil
	c.lastUpdatedKeys = make(map[int64]keyContainer[K])
}

// DropKey removes the value associated with the provided key from the cache. The operation is thread-safe.
func (c *InMemoryHashMapCache[K, T]) DropKey(key K) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	hash := key.Key()
	c.dropKey(key.Key())
	delete(c.changes, hash)
	delete(c.lastUpdatedKeys, key.Key())
}

// Outdated checks if the provided key or the entire cache (if no key is provided)
// is outdated based on the set TTL. Returns true if outdated, false otherwise.
func (c *InMemoryHashMapCache[K, T]) Outdated(key uopt.Opt[K]) bool {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()

	if c.ttl == nil {
		return false
	} else {
		if key.Present() {
			k := key.Get()
			if lu, ok := c.lastUpdatedKeys[(*k).Key()]; ok {
				return time.Since(lu.updatedAt) > *c.ttl
			} else {
				return true
			}
		} else {
			return false
		}
	}
}

func (c *InMemoryHashMapCache[K, T]) dropAll() {
	c.values = make(map[int64][]hashValueContainer[K, T])
}

func (c *InMemoryHashMapCache[K, T]) put(key K, value T) {
	hash := c.addTran(key, value)
	changes := len(c.changes) == 0
	found := false
	for _, diff := range c.changes {
		if diff.Key() == key.Key() {
			if !diff.Equals(key) {
				changes = true
				break
			}
			found = true
			continue
		}
	}
	if changes || !found {
		c.changes[hash] = key
	}
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
		ind := -1
		for i, v := range values {
			if v.key.Equals(key) {
				ind = i
			}
		}
		if ind != -1 {
			values[ind] = hashValueContainer[K, T]{
				key:   key,
				value: value,
			}
		} else {
			values = append(values, hashValueContainer[K, T]{
				key:   key,
				value: value,
			})
			c.values[keyHash] = values
		}
	}

	return keyHash
}

func (c *InMemoryHashMapCache[K, T]) dropKey(key int64) {
	delete(c.values, key)
}

// InMemoryComparableMapCache provides an in-memory caching mechanism using Go's native maps for single-value entries.
// It supports optional TTL for entries and ensures concurrency-safe operations using a mutex.
// It is very similiar to InMemoryHashMapCache by behaviour, and the only difference is a key type constraint.
type InMemoryComparableMapCache[K comparable, T any] struct {
	values  map[K]T
	changes uset.Set[K]

	lastUpdatedKeys map[K]time.Time
	lastUpdated     time.Time

	ttl  *time.Duration
	vMtx sync.Mutex
}

// NewInMemoryComparableMapCache creates a new instance of InMemoryComparableMapCache.
// It accepts an optional TTL (time-to-live) duration for cache entries.
func NewInMemoryComparableMapCache[K comparable, T any](ttl uopt.Opt[time.Duration]) *InMemoryComparableMapCache[K, T] {
	c := &InMemoryComparableMapCache[K, T]{
		values:          make(map[K]T),
		changes:         uset.NewHashSet[K](),
		lastUpdatedKeys: make(map[K]time.Time),
	}
	ttl.IfPresent(func(t time.Duration) {
		c.ttl = &t
	})
	return c
}

// Set updates the cache value for the provided key. If the key already exists,
// its previous value is replaced with the new value. The operation is thread-safe.
func (c *InMemoryComparableMapCache[K, T]) Set(key K, value T) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.values[key] = value
	c.changes.Add(key)
	now := time.Now()
	c.lastUpdatedKeys[key] = now
	c.lastUpdated = now
}

// SetQuietly adds a value to the cache for the provided key without altering the change history.
// This method is thread-safe and optimized for performance when change tracking is unnecessary.
func (c *InMemoryComparableMapCache[K, T]) SetQuietly(key K, value T) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.values[key] = value
	now := time.Now()
	c.lastUpdatedKeys[key] = now
	c.lastUpdated = now
}

// Get retrieves the value associated with the provided key from the cache.
// It returns a pointer to the value and a boolean indicating whether the key was found.
// The operation is thread-safe.
func (c *InMemoryComparableMapCache[K, T]) Get(key K) (*T, bool) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()

	value, ok := c.values[key]
	if !ok {
		return nil, false
	}
	return &value, true
}

// Changes returns a slice of keys that have been modified in the cache since the last call to Changes.
// This method is thread-safe.
func (c *InMemoryComparableMapCache[K, T]) Changes() []K {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()

	return c.changes.Values()
}

// Drop completely clears the cache, removing all entries. The operation is thread-safe.
func (c *InMemoryComparableMapCache[K, T]) Drop() {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.values = make(map[K]T)
	c.changes.Clear()
	c.lastUpdatedKeys = make(map[K]time.Time)
	c.lastUpdated = time.Time{}
}

// DropKey removes the value associated with the provided key from the cache.
// The operation is thread-safe.
func (c *InMemoryComparableMapCache[K, T]) DropKey(key K) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	delete(c.values, key)
	c.changes.Remove(key)
	delete(c.lastUpdatedKeys, key)
}

// Outdated checks if the provided key is outdated based on the set TTL (time-to-live).
// Returns true if outdated, false otherwise.
// If no TTL is set or the key does not exist, it returns false.
func (c *InMemoryComparableMapCache[K, T]) Outdated(key uopt.Opt[K]) bool {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()

	if c.ttl == nil {
		return false
	}

	if k := key.Get(); k != nil {
		lastUpdated, exists := c.lastUpdatedKeys[*k]
		if !exists {
			return true
		}
		return time.Since(lastUpdated) > *c.ttl
	}

	return false
}
