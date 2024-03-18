package cache

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"sync"
	"time"

	"github.com/dgryski/go-farm"
	arrayutils "github.com/kordax/basic-utils/array-utils"
	"github.com/kordax/basic-utils/opt"
)

type container[K CompositeKey, T Comparable] struct {
	pairs map[int][]arrayutils.Pair[K, T]
	node  map[int]any
}

type Cache[K CompositeKey, T Comparable] interface {
	Put(key K, values ...T) // Adds a new value to the key
	Set(key K, values ...T) // Overwrites key values
	Get(key K) []T
	AddSilently(key K, values ...T) // Adds the values but update cache state/doesn't add any changes to the cache
	Changes() []K
	Drop()
	DropKey(key K)
	Outdated(key opt.Opt[K]) bool
}

// InMemoryTreeCache provides an in-memory caching mechanism with support for compound keys.
// The cache leverages tree-like structures to store and organize data, allowing efficient
// operations even with composite keys. The cache supports optional TTL (time-to-live) for entries,
// ensuring that outdated entries can be identified and potentially purged. Concurrency-safe
// operations are ensured through the use of a mutex.
//
// Benchmark insights:
// - Put operation performance is fast for shallow depth keys but slows down as the depth increases.
// - Get operation is particularly efficient, especially for shallow depth keys.
// - Set operation's performance is consistent regardless of the depth of the key.
type InMemoryTreeCache[K CompositeKey, T Comparable] struct {
	values  map[int]any
	changes []K

	lastUpdatedKeys map[string]time.Time
	lastUpdated     time.Time
	ttl             *time.Duration

	vMtx sync.Mutex
}

// NewInMemoryTreeCache creates a new instance of the InMemoryTreeCache.
// It takes an optional TTL (time-to-live) parameter to set expiration time for cache entries.
// If the TTL is not provided, cache entries will not expire.
func NewInMemoryTreeCache[K CompositeKey, T Comparable](ttl opt.Opt[time.Duration]) *InMemoryTreeCache[K, T] {
	c := &InMemoryTreeCache[K, T]{
		values:          make(map[int]any),
		changes:         make([]K, 0),
		lastUpdatedKeys: make(map[string]time.Time),
	}
	ttl.IfPresent(func(t time.Duration) {
		c.ttl = &t
	})

	return c
}

// Put inserts a new value(s) into the cache associated with the given key.
// If the key already exists in the cache, it appends the new value(s) to the existing values.
// This operation is relatively fast for shallow depth keys, but becomes slower as the depth increases.
func (c *InMemoryTreeCache[K, T]) Put(key K, val ...T) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.put(key, val...)
	c.lastUpdatedKeys[key.String()] = time.Now()
	c.lastUpdated = time.Now()
}

// Set inserts a new value(s) into the cache associated with the given key.
// If the key already exists in the cache, this method will overwrite the existing values.
func (c *InMemoryTreeCache[K, T]) Set(key K, val ...T) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.dropKeyRecursively(key.Keys(), 0, c.values)
	c.put(key, val...)
	c.lastUpdatedKeys[key.String()] = time.Now()
	c.lastUpdated = time.Now()
}

// AddSilently behaves like the Put method but does not update the cache state or add any changes to the cache.
// This method is useful when you want to add values to the cache without triggering any side effects.
func (c *InMemoryTreeCache[K, T]) AddSilently(key K, val ...T) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.addTran(key, val...)
	c.lastUpdatedKeys[key.String()] = time.Now()
	c.lastUpdated = time.Now()
}

// Get retrieves the value(s) associated with the given key from the cache.
// If the key is not found, it returns an empty slice.
// Retrieval is fast, especially for shallow depth keys.
func (c *InMemoryTreeCache[K, T]) Get(key K) []T {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.changes = nil
	bucket := c.tryToGetBucket(key.Keys())
	result := make([]T, 0)
	for _, pairs := range bucket {
		for _, p := range pairs {
			result = append(result, p.Right)
		}
	}

	return result
}

// Changes returns a slice of keys that have been modified in the cache.
// This method provides a way to track changes made to the cache, useful for scenarios like cache syncing.
func (c *InMemoryTreeCache[K, T]) Changes() []K {
	return c.changes
}

// Drop removes all entries from the cache.
// This is a complete reset of the cache, useful when you want to clear the cache and start fresh.
func (c *InMemoryTreeCache[K, T]) Drop() {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.dropAll()
	c.lastUpdatedKeys = make(map[string]time.Time)
}

// DropKey removes the value(s) associated with the given key from the cache.
func (c *InMemoryTreeCache[K, T]) DropKey(key K) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.dropKeyRecursively(key.Keys(), 0, c.values)
	c.lastUpdatedKeys[key.String()] = time.Now()
	c.lastUpdated = time.Now()
}

// Outdated checks if a given key or the entire cache is outdated based on the TTL.
// If no key is provided, it checks the last updated time of the entire cache.
// If a key is provided, it checks the last updated time of that specific key.
func (c *InMemoryTreeCache[K, T]) Outdated(key opt.Opt[K]) bool {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()

	if c.ttl == nil {
		return false
	} else {
		if key.Present() {
			k := key.Get()
			if lu, ok := c.lastUpdatedKeys[(*k).String()]; ok {
				return time.Since(lu) > *c.ttl
			} else {
				return true
			}
		} else {
			return time.Since(c.lastUpdated) > *c.ttl
		}
	}
}

func (c *InMemoryTreeCache[K, T]) dropAll() {
	c.values = make(map[int]any)
	c.changes = nil
}

func (c *InMemoryTreeCache[K, T]) put(key K, val ...T) {
	c.addTran(key, val...)
	changes := len(c.changes) == 0
	found := false
	for _, diff := range c.changes {
		if arrayutils.EqualsWithOrder(diff.Keys(), key.Keys()) {
			if !diff.Equals(key) {
				changes = true
				break
			}
			found = true
			continue
		}
	}
	if changes || !found {
		c.changes = append(c.changes, key)
	}
}

func (c *InMemoryTreeCache[K, T]) addTran(key K, values ...T) {
	hashes := key.Keys()
	if len(hashes) == 0 {
		return
	}

	bucket := c.tryToGetBucket(hashes)
	lHash := key.Keys()[len(hashes)-1]

	for _, value := range values {
		if ind, _ := arrayutils.ContainsPredicate(bucket[lHash], func(v *arrayutils.Pair[K, T]) bool {
			return v.Right.Equals(value)
		}); ind > -1 {
			bucket[lHash][ind] = *arrayutils.NewPair[K, T](key, value)
		} else {
			bucket[lHash] = append(bucket[lHash], *arrayutils.NewPair[K, T](key, value))
		}
	}
}

func (c *InMemoryTreeCache[K, T]) dropKeyRecursively(keys []int, n int, bucket map[int]any) {
	hash := keys[n]
	interBucket := bucket[hash]
	if interBucket != nil {
		switch b := interBucket.(type) {
		case container[K, T]:
			if n+1 == len(keys) {
				delete(bucket, hash)
			} else {
				c.dropKeyRecursively(keys, n+1, b.node)
			}
		default:
			delete(bucket, hash)
		}
	}
}

func (c *InMemoryTreeCache[K, T]) tryToGetBucket(keys []int) map[int][]arrayutils.Pair[K, T] {
	return c.getBucket(keys, 0, c.values)
}

func (c *InMemoryTreeCache[K, T]) getBucket(keys []int, n int, interBucket map[int]any) map[int][]arrayutils.Pair[K, T] {
	if keys == nil || n >= len(keys) {
		return nil
	}

	if bucket, ok := interBucket[keys[n]]; ok {
		switch b := bucket.(type) {
		case map[int][]arrayutils.Pair[K, T]:
			if n+1 < len(keys) {
				interBucket[keys[n]] = container[K, T]{
					node:  make(map[int]any),
					pairs: b,
				}
				return c.getBucket(keys, n+1, interBucket[keys[n]].(container[K, T]).node)
			} else {
				return b
			}
		case container[K, T]:
			if n+1 == len(keys) {
				result := make(map[int][]arrayutils.Pair[K, T])
				for k, e := range b.pairs {
					result[k] = append(result[k], e...)
				}
				if b.node != nil {
					result = c.getNodePairsFlat(b.node, result)
				}

				return result
			}

			return c.getBucket(keys, n+1, b.node)
		}
	} else {
		if n+1 == len(keys) {
			interBucket[keys[n]] = map[int][]arrayutils.Pair[K, T]{
				keys[n]: nil,
			}
			return interBucket[keys[n]].(map[int][]arrayutils.Pair[K, T])
		} else {
			if entry, ok := interBucket[keys[n]]; !ok {
				interBucket[keys[n]] = container[K, T]{
					node:  make(map[int]any),
					pairs: make(map[int][]arrayutils.Pair[K, T]),
				}
				return c.getBucket(keys, n+1, interBucket[keys[n]].(container[K, T]).node)
			} else {
				switch e := entry.(type) {
				case map[int][]arrayutils.Pair[K, T]:
					interBucket[keys[n]] = container[K, T]{
						node:  make(map[int]any),
						pairs: e,
					}
					return c.getBucket(keys, n+1, interBucket[keys[n]].(container[K, T]).node)
				case container[K, T]:
					interBucket[keys[n]] = container[K, T]{
						pairs: e.pairs,
					}
					return c.getBucket(keys, n+1, e.node)
				}
			}
		}
	}

	return nil
}

func (c *InMemoryTreeCache[K, T]) getNodePairsFlat(node map[int]any, result map[int][]arrayutils.Pair[K, T]) map[int][]arrayutils.Pair[K, T] {
	for _, entry := range node {
		switch e := entry.(type) {
		case map[int][]arrayutils.Pair[K, T]:
			for hash, pair := range e {
				result[hash] = append(result[hash], pair...)
			}
		case container[K, T]:
			for hash, pair := range e.pairs {
				result[hash] = append(result[hash], pair...)
			}
			result = c.getNodePairsFlat(e.node, result)
		}
	}

	return result
}

// InMemoryHashMapCache provides an in-memory caching mechanism using hashmaps.
// This cache structure translates composite keys into a hash value using a user-provided
// hashing function. The cache supports optional TTL (time-to-live) for entries.
// Concurrency-safe operations are ensured through the use of a mutex.
//
// Performance Comparison with InMemoryTreeCache:
// - Insertions: InMemoryTreeCache is slightly faster for single-depth insertions and significantly faster for deeper depths.
// - Retrievals: InMemoryHashMapCache (especially with FarmHash) is faster for both single-depth and deeper retrievals.
// - Setting Values: Performance is relatively close between the two, with minor variations.
//
// The choice between InMemoryTreeCache and InMemoryHashMapCache would depend on the specific use case,
// especially the depth of the keys and the frequency of retrieval operations.
type InMemoryHashMapCache[K CompositeKey, T Comparable, H comparable] struct {
	values  map[H][]T
	changes []K

	lastUpdatedKeys map[string]time.Time
	lastUpdated     time.Time
	ttl             *time.Duration

	toHash func(keys []int) H
	vMtx   sync.Mutex
}

// NewInMemoryHashMapCache creates a new instance of the InMemoryHashMapCache.
// It takes a hashing function to translate the composite keys to a desired hash type,
// and an optional time-to-live duration for the cache entries.
func NewInMemoryHashMapCache[K CompositeKey, T Comparable, H comparable](toHash func(keys []int) H, ttl opt.Opt[time.Duration]) *InMemoryHashMapCache[K, T, H] {
	c := &InMemoryHashMapCache[K, T, H]{
		values:          make(map[H][]T),
		changes:         make([]K, 0),
		lastUpdatedKeys: make(map[string]time.Time),
		toHash:          toHash,
	}
	ttl.IfPresent(func(t time.Duration) {
		c.ttl = &t
	})

	return c
}

// NewDefaultHashMapCache creates a new instance of the InMemoryHashMapCache using SHA256 as the hashing algorithm.
func NewDefaultHashMapCache[K CompositeKey, T Comparable](ttl opt.Opt[time.Duration]) *InMemoryHashMapCache[K, T, uint64] {
	return NewFarmHashMapCache[K, T](ttl)
}

func NewFarmHashMapCache[K CompositeKey, T Comparable](ttl opt.Opt[time.Duration]) *InMemoryHashMapCache[K, T, uint64] {
	buffer := new(bytes.Buffer)
	return NewInMemoryHashMapCache[K, T, uint64](func(keys []int) uint64 {
		arr := make([]byte, 0)
		for _, hash := range keys {
			arr = append(arr, intToBytes(buffer, hash)...)
		}

		return farm.Hash64(arr)
	}, ttl)
}

func NewSha256HashMapCache[K CompositeKey, T Comparable](ttl opt.Opt[time.Duration]) *InMemoryHashMapCache[K, T, string] {
	buffer := new(bytes.Buffer)
	return NewInMemoryHashMapCache[K, T, string](func(keys []int) string {
		arr := make([]byte, 0)
		for _, hash := range keys {
			arr = append(arr, intToBytes(buffer, hash)...)
		}

		h := sha256.New()
		h.Write(arr)

		return string(h.Sum(nil))
	}, ttl)
}

// Put adds the given values to the cache associated with the provided key.
// If the key already exists, the values are updated. The insertion is thread-safe.
func (c *InMemoryHashMapCache[K, T, H]) Put(key K, values ...T) {
	if len(values) == 0 {
		return
	}
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.put(key, values...)
	c.lastUpdatedKeys[key.String()] = time.Now()
	c.lastUpdated = time.Now()
}

// Set updates the cache values for the provided key. If the key already exists,
// its previous values are removed before adding the new values. The operation is thread-safe.
func (c *InMemoryHashMapCache[K, T, H]) Set(key K, values ...T) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.dropKey(key.Keys())
	c.put(key, values...)
	c.lastUpdatedKeys[key.String()] = time.Now()
	c.lastUpdated = time.Now()
}

// AddSilently adds values to the cache for the provided key but does so without
// altering the change history. This operation can be used when modifications should not trigger cache change logs.
func (c *InMemoryHashMapCache[K, T, H]) AddSilently(key K, values ...T) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.addTran(key, values...)
	c.lastUpdatedKeys[key.String()] = time.Now()
	c.lastUpdated = time.Now()
}

// Get retrieves the values associated with the provided key from the cache.
// The operation is thread-safe and does not alter the change history.
func (c *InMemoryHashMapCache[K, T, H]) Get(key K) []T {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.changes = nil

	return c.values[c.toHash(key.Keys())]
}

// Changes returns a list of keys that have experienced changes in the cache since the last reset.
func (c *InMemoryHashMapCache[K, T, H]) Changes() []K {
	return c.changes
}

// Drop completely clears the cache, removing all entries. The operation is thread-safe.
func (c *InMemoryHashMapCache[K, T, H]) Drop() {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.dropAll()
	c.lastUpdatedKeys = make(map[string]time.Time)
}

// DropKey removes the values associated with the provided key from the cache. The operation is thread-safe.
func (c *InMemoryHashMapCache[K, T, H]) DropKey(key K) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.dropKey(key.Keys())
	c.lastUpdatedKeys[key.String()] = time.Now()
	c.lastUpdated = time.Now()
}

// Outdated checks if the provided key or the entire cache (if no key is provided)
// is outdated based on the set TTL. Returns true if outdated, false otherwise.
func (c *InMemoryHashMapCache[K, T, H]) Outdated(key opt.Opt[K]) bool {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()

	if c.ttl == nil {
		return false
	} else {
		if key.Present() {
			k := key.Get()
			if lu, ok := c.lastUpdatedKeys[(*k).String()]; ok {
				return time.Since(lu) > *c.ttl
			} else {
				return true
			}
		} else {
			return time.Since(c.lastUpdated) > *c.ttl
		}
	}
}

func (c *InMemoryHashMapCache[K, T, H]) dropAll() {
	c.values = make(map[H][]T)
	c.changes = nil
}

func (c *InMemoryHashMapCache[K, T, H]) put(key K, values ...T) {
	c.addTran(key, values...)
	changes := len(c.changes) == 0
	found := false
	for _, diff := range c.changes {
		if arrayutils.EqualsWithOrder(diff.Keys(), key.Keys()) {
			if !diff.Equals(key) {
				changes = true
				break
			}
			found = true
			continue
		}
	}
	if changes || !found {
		c.changes = append(c.changes, key)
	}
}

func (c *InMemoryHashMapCache[K, T, H]) addTran(key K, values ...T) {
	keys := key.Keys()
	if len(values) == 0 {
		return
	}

	for i := 0; i < len(keys); i++ {
		hash := c.toHash(keys[:i+1])
		for _, value := range values {
			if existing, found := c.values[hash]; found {
				if ind, entry := arrayutils.ContainsPredicate[T](existing, func(v *T) bool {
					return (*v).Equals(value)
				}); entry == nil {
					// Collision detected
					c.values[hash] = append(existing, value)
				} else {
					// Else replace, this value is already hashed
					c.values[hash][ind] = value
				}
			} else {
				c.values[hash] = []T{value}
			}
		}
	}
}

func (c *InMemoryHashMapCache[K, T, H]) dropKey(keys []int) {
	delete(c.values, c.toHash(keys))
}

func intToBytes(buffer *bytes.Buffer, num int) []byte {
	buffer.Reset()
	_ = binary.Write(buffer, binary.LittleEndian, int64(num))

	return buffer.Bytes()
}
