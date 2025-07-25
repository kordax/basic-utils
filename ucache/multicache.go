package ucache

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dgryski/go-farm"
	"github.com/kordax/basic-utils/v2/uarray"
	"github.com/kordax/basic-utils/v2/uconst"
	"github.com/kordax/basic-utils/v2/umap"
	"github.com/kordax/basic-utils/v2/uopt"
)

type container[K CompositeKey, T uconst.Comparable] struct {
	pairs map[int64][]uarray.Pair[K, T]
	node  map[int64]any
}

// The MultiCache interface defines a set of methods for a generic cache implementation.
// This interface supports setting, getting, and managing cache entries with composite keys.
// Unlike Cache, it is designed to handle multiple values per key and has a hierarchical key handling.
//
// Example:
// - Set([1, 2, 3], "Value1") => [1, 2, 3] is set with "Value1".
// - Set([1, 2, 3, 4], "Value2") => [1, 2, 3] is replaced, and [1, 2, 3, 4] is set with "Value2".
// - Get([1, 2, 3]) => returns nil, as it has been replaced by [1, 2, 3, 4].
//
// This hierarchical key handling is useful for scenarios where more specific keys should override
// the values of their parent keys, providing a clear and structured way to manage cache entries.
type MultiCache[K CompositeKey, T any] interface {
	// Put inserts a new value(s) into the cache associated with the given key.
	// If the key already exists in the cache, it appends the new value(s) to the existing values.
	// This operation is relatively fast for shallow depth keys, but becomes slower as the depth increases.
	Put(key K, values ...T)

	// Set inserts a new value(s) into the cache associated with the given key.
	// If the key already exists in the cache, this method will overwrite the existing values.
	Set(key K, values ...T)

	// Get retrieves the value(s) associated with the given key from the cache.
	// If the key is not found, it returns an empty slice.
	// Retrieval is fast, especially for shallow depth keys.
	// Supports retrieving a value using a broader key (e.g., [1, 2]) or a full/shallow key (e.g., [1, 2, 3, 4])
	Get(key K) []T

	// Changes returns a slice of keys that have been modified in the cache.
	// This method provides a way to track changes made to the cache, useful for scenarios like cache syncing.
	// Cache changes will be updated only on modifying operations, meaning that in-fact, changes contain all the present keys.
	Changes() []K

	// Drop removes all entries from the cache.
	// This is a complete reset of the cache, useful when you want to clear the cache and start fresh.
	Drop()

	// DropKey removes the value(s) associated with the given key from the cache.
	DropKey(key K)

	// Outdated checks if a given key or the entire cache is outdated based on the TTL.
	// If no key is provided it checks the last updated time of the entire cache.
	// If a key is provided and found, it checks the last updated time of that specific key.
	// If key was not found returns false.
	Outdated(key uopt.Opt[K]) bool

	// PutQuietly behaves like the Put method but does not update the cache state or add any changes to the cache, making it
	// much faster alternative to Put and Set.
	// This method is useful when you want to add values to the cache without triggering any side effects.
	PutQuietly(key K, values ...T)
}

// InMemoryTreeMultiCache provides an in-memory caching mechanism with support for compound keys.
// The cache leverages tree-like structures to store and organize data, allowing efficient
// operations even with composite keys. The cache supports optional TTL (time-to-live) for entries,
// ensuring that outdated entries can be identified and potentially purged. Concurrency-safe
// operations are ensured through the use of a mutex.
//
// Benchmark insights:
// - Put operation performance is fast for shallow depth keys but slows down as the depth increases.
// - Get operation is particularly efficient, especially for shallow depth keys.
// - Set operation's performance is consistent regardless of the depth of the key.
// TTL parameter in cache doesn't automatically clean up all the entries.
// Use ManagedMultiCache wrapper to automatically manage outdated keys.
type InMemoryTreeMultiCache[K CompositeKey, T uconst.Comparable] struct {
	values  map[int64]any
	changes []K

	lastUpdatedKeys map[string]time.Time
	lastUpdated     time.Time
	ttl             *time.Duration

	vMtx sync.Mutex
}

// NewInMemoryTreeMultiCache creates a new instance of the InMemoryTreeMultiCache.
// It takes an optional TTL (time-to-live) parameter to set expiration time for cache entries.
// If the TTL is not provided, cache entries will not expire.
// Note on hierarchical key handling:
//   - If a composite key (e.g., [1, 2, 3]) is already set, any broader keys that share the same prefix
//     (e.g., [1, 2]) are considered "busy" as part of the hierarchy.
//   - Setting a more specific key (e.g., [1, 2, 3, 4]) will replace the broader key's value (e.g., [1, 2, 3]).
//   - This design ensures that more specific keys take precedence and can replace the values of their parent keys.
//   - Additionally, retrieving a value using a broader key (e.g., [1, 2]) will return the values of the most specific key
//     that shares the prefix (e.g., [1, 2, 3, 4]).
func NewInMemoryTreeMultiCache[K CompositeKey, T uconst.Comparable](ttl uopt.Opt[time.Duration]) MultiCache[K, T] {
	c := &InMemoryTreeMultiCache[K, T]{
		values:          make(map[int64]any),
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
func (c *InMemoryTreeMultiCache[K, T]) Put(key K, val ...T) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.put(key, val...)
	c.lastUpdatedKeys[keysAsString(key.Keys())] = time.Now()
	c.lastUpdated = time.Now()
}

// Set inserts a new value(s) into the cache associated with the given key.
// If the key already exists in the cache, this method will overwrite the existing values.
func (c *InMemoryTreeMultiCache[K, T]) Set(key K, val ...T) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.dropKeyRecursively(key.Keys(), 0, c.values)
	c.put(key, val...)
	c.lastUpdatedKeys[keysAsString(key.Keys())] = time.Now()
	c.lastUpdated = time.Now()
}

// PutQuietly behaves like the Put method but does not update the cache state or add any changes to the cache, making it
// much faster alternative to Put and Set.
// This method is useful when you want to add values to the cache without triggering any side effects.
func (c *InMemoryTreeMultiCache[K, T]) PutQuietly(key K, val ...T) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.addTran(key, val...)
	c.lastUpdatedKeys[keysAsString(key.Keys())] = time.Now()
	c.lastUpdated = time.Now()
}

// Get retrieves the value(s) associated with the given key from the cache.
// If the key is not found, it returns an empty slice.
// Retrieval is fast, especially for shallow depth keys.
func (c *InMemoryTreeMultiCache[K, T]) Get(key K) []T {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()

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
func (c *InMemoryTreeMultiCache[K, T]) Changes() []K {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	return c.changes
}

// Drop removes all entries from the cache.
// This is a complete reset of the cache, useful when you want to clear the cache and start fresh.
func (c *InMemoryTreeMultiCache[K, T]) Drop() {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.dropAll()
	c.lastUpdatedKeys = make(map[string]time.Time)
}

// DropKey removes the value(s) associated with the given key from the cache.
func (c *InMemoryTreeMultiCache[K, T]) DropKey(key K) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.dropKeyRecursively(key.Keys(), 0, c.values)
	delete(c.lastUpdatedKeys, keysAsString(key.Keys()))
	ind, _ := uarray.ContainsPredicate(c.changes, func(v K) bool {
		return v.Equals(key)
	})
	if ind > -1 {
		c.changes = uarray.CopyWithoutIndex(c.changes, ind)
	}
}

// Outdated checks if a given key or the entire cache is outdated based on the TTL.
// If no key is provided or key was not found, it checks the last updated time of the entire cache.
// If a key is provided and found, it checks the last updated time of that specific key.
func (c *InMemoryTreeMultiCache[K, T]) Outdated(key uopt.Opt[K]) bool {
	if !key.Present() {
		return time.Since(c.lastUpdated) > *c.ttl
	}

	c.vMtx.Lock()
	defer c.vMtx.Unlock()

	if c.ttl == nil {
		return false
	} else {
		if key.Present() {
			k := key.Get()
			if lu, ok := c.lastUpdatedKeys[keysAsString((*k).Keys())]; ok {
				return time.Since(lu) > *c.ttl
			} else {
				return true
			}
		} else {
			return false
		}
	}
}

func (c *InMemoryTreeMultiCache[K, T]) dropAll() {
	c.values = make(map[int64]any)
	c.changes = nil
}

func (c *InMemoryTreeMultiCache[K, T]) put(key K, val ...T) {
	c.addTran(key, val...)
	changes := len(c.changes) == 0
	found := false
	for _, diff := range c.changes {
		if uarray.EqualsWithOrder(diff.Keys(), key.Keys()) {
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

func (c *InMemoryTreeMultiCache[K, T]) addTran(key K, values ...T) {
	keys := key.Keys()
	if len(keys) == 0 {
		return
	}

	bucket := c.tryToGetBucket(keys)
	lowKey := key.Keys()[len(keys)-1].Key()

	for _, value := range values {
		if ind, _ := uarray.ContainsPredicate(bucket[lowKey], func(v uarray.Pair[K, T]) bool {
			return v.Right.Equals(value)
		}); ind > -1 {
			bucket[lowKey][ind] = *uarray.NewPair[K, T](key, value)
		} else {
			bucket[lowKey] = append(bucket[lowKey], *uarray.NewPair[K, T](key, value))
		}
	}
}

func (c *InMemoryTreeMultiCache[K, T]) dropKeyRecursively(keys []uconst.Unique, n int, bucket map[int64]any) {
	key := keys[n].Key()
	interBucket := bucket[key]
	if interBucket != nil {
		switch b := interBucket.(type) {
		case container[K, T]:
			if n+1 == len(keys) {
				delete(bucket, key)
			} else {
				c.dropKeyRecursively(keys, n+1, b.node)
			}
		default:
			delete(bucket, key)
		}
	}
}

func (c *InMemoryTreeMultiCache[K, T]) tryToGetBucket(keys []uconst.Unique) map[int64][]uarray.Pair[K, T] {
	return c.getBucket(keys, 0, c.values)
}

func (c *InMemoryTreeMultiCache[K, T]) getBucket(keys []uconst.Unique, n int, interBucket map[int64]any) map[int64][]uarray.Pair[K, T] {
	if keys == nil || n >= len(keys) {
		return nil
	}

	hash := keys[n].Key()
	if bucket, ok := interBucket[hash]; ok {
		switch b := bucket.(type) {
		case map[int64][]uarray.Pair[K, T]:
			if n+1 < len(keys) {
				interBucket[hash] = container[K, T]{
					node:  make(map[int64]any),
					pairs: b,
				}
				return c.getBucket(keys, n+1, interBucket[hash].(container[K, T]).node)
			} else {
				return b
			}
		case container[K, T]:
			if n+1 == len(keys) {
				result := make(map[int64][]uarray.Pair[K, T])
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
			interBucket[hash] = map[int64][]uarray.Pair[K, T]{
				hash: nil,
			}
			return interBucket[hash].(map[int64][]uarray.Pair[K, T])
		} else {
			if entry, ok := interBucket[hash]; !ok {
				interBucket[hash] = container[K, T]{
					node:  make(map[int64]any),
					pairs: make(map[int64][]uarray.Pair[K, T]),
				}
				return c.getBucket(keys, n+1, interBucket[hash].(container[K, T]).node)
			} else {
				switch e := entry.(type) {
				case map[int64][]uarray.Pair[K, T]:
					interBucket[hash] = container[K, T]{
						node:  make(map[int64]any),
						pairs: e,
					}
					return c.getBucket(keys, n+1, interBucket[hash].(container[K, T]).node)
				case container[K, T]:
					interBucket[hash] = container[K, T]{
						pairs: e.pairs,
					}
					return c.getBucket(keys, n+1, e.node)
				}
			}
		}
	}

	return nil
}

func (c *InMemoryTreeMultiCache[K, T]) getNodePairsFlat(node map[int64]any, result map[int64][]uarray.Pair[K, T]) map[int64][]uarray.Pair[K, T] {
	for _, entry := range node {
		switch e := entry.(type) {
		case map[int64][]uarray.Pair[K, T]:
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

// InMemoryHashMapMultiCache provides an in-memory caching mechanism using hashmaps.
// Unlike InMemoryHashMapCache, it stores multiple values per key.
// This cache structure translates composite keys into a hash value using a user-provided
// hashing function. The cache supports optional TTL (time-to-live) for entries (Please read below regarding clean up.)
// Concurrency-safe operations are ensured through the use of a mutex.
// Unlike InMemoryTreeMultiCache, it doesn't support a hierarchy for keys, so each key is unique.
// - Setting a more specific key (e.g., [1, 2, 3, 4]) WILL NOT replace the value of a broader key (e.g., [1, 2, 3]).
// - Deleting a specific key or a parent key (e.g., [1, 2]) will not affect the other one.
//
// Performance Comparison with InMemoryTreeMultiCache:
// - Insertions: InMemoryTreeMultiCache is slightly faster for single-depth insertions and significantly faster for deeper depths.
// - Retrievals: InMemoryHashMapMultiCache (especially with FarmHash) is faster for both single-depth and deeper retrievals.
// - Setting Values: Performance is relatively close between the two, with minor variations.
//
// The choice between InMemoryTreeMultiCache and InMemoryHashMapMultiCache would depend on the specific use case,
// especially the depth of the keys and the frequency of retrieval operations.
// TTL parameter in cache doesn't automatically clean up all the entries.
// Use ManagedMultiCache wrapper to automatically manage outdated keys.
type InMemoryHashMapMultiCache[K CompositeKey, T any, H comparable] struct {
	values  map[H][]T
	changes map[H]K

	lastUpdatedKeys map[string]keyContainer[K]
	lastUpdated     time.Time
	ttl             *time.Duration

	toHash func(keys []uconst.Unique) H
	vMtx   sync.Mutex
}

// NewInMemoryHashMapMultiCache creates a new instance of the InMemoryHashMapMultiCache.
// It takes a hashing function to translate the composite keys to a desired hash type,
// and an optional time-to-live duration for the cache entries.
func NewInMemoryHashMapMultiCache[K CompositeKey, T any, H comparable](toHash func(keys []uconst.Unique) H, ttl uopt.Opt[time.Duration]) *InMemoryHashMapMultiCache[K, T, H] {
	c := &InMemoryHashMapMultiCache[K, T, H]{
		values:          make(map[H][]T),
		changes:         make(map[H]K, 0),
		lastUpdatedKeys: make(map[string]keyContainer[K]),
		toHash:          toHash,
	}
	ttl.IfPresent(func(t time.Duration) {
		c.ttl = &t
	})

	return c
}

// NewDefaultHashMapMultiCache creates a new instance of the InMemoryHashMapMultiCache using SHA256 as the hashing algorithm.
func NewDefaultHashMapMultiCache[K CompositeKey, T uconst.Comparable](ttl uopt.Opt[time.Duration]) *InMemoryHashMapMultiCache[K, T, uint64] {
	return NewFarmHashMapMultiCache[K, T](ttl)
}

func NewFarmHashMapMultiCache[K CompositeKey, T uconst.Comparable](ttl uopt.Opt[time.Duration]) *InMemoryHashMapMultiCache[K, T, uint64] {
	return NewInMemoryHashMapMultiCache[K, T, uint64](func(keys []uconst.Unique) uint64 {
		buffer := new(bytes.Buffer)
		arr := make([]byte, 0)
		for _, hash := range keys {
			arr = append(arr, intToBytes(buffer, hash.Key())...)
		}

		return farm.Hash64(arr)
	}, ttl)
}

func NewSha256HashMapMultiCache[K CompositeKey, T uconst.Comparable](ttl uopt.Opt[time.Duration]) MultiCache[K, T] {
	return NewInMemoryHashMapMultiCache[K, T, string](func(keys []uconst.Unique) string {
		buffer := new(bytes.Buffer)
		arr := make([]byte, 0)
		for _, hash := range keys {
			arr = append(arr, intToBytes(buffer, hash.Key())...)
		}

		h := sha256.New()
		h.Write(arr)

		return string(h.Sum(nil))
	}, ttl)
}

// Put adds the given values to the cache associated with the provided key.
// If the key already exists, the values are updated. The insertion is thread-safe.
func (c *InMemoryHashMapMultiCache[K, T, H]) Put(key K, values ...T) {
	if len(values) == 0 {
		return
	}
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.put(key, values...)
	n := time.Now()
	c.lastUpdatedKeys[keysAsString(key.Keys())] = keyContainer[K]{
		key:       key,
		updatedAt: n,
	}
	c.lastUpdated = n
}

// Set updates the cache values for the provided key. If the key already exists,
// its previous values are removed before adding the new values. The operation is thread-safe.
func (c *InMemoryHashMapMultiCache[K, T, H]) Set(key K, values ...T) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.dropKey(key.Keys())
	c.put(key, values...)
	n := time.Now()
	c.lastUpdatedKeys[keysAsString(key.Keys())] = keyContainer[K]{
		key:       key,
		updatedAt: n,
	}
	c.lastUpdated = n
}

// PutQuietly adds values to the cache for the provided key but does so without
// altering the change history. This operation can be used when modifications should not trigger cache change diff.
func (c *InMemoryHashMapMultiCache[K, T, H]) PutQuietly(key K, values ...T) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.addTran(key, values...)
	n := time.Now()
	c.lastUpdatedKeys[keysAsString(key.Keys())] = keyContainer[K]{
		key:       key,
		updatedAt: n,
	}
	c.lastUpdated = n
}

// Get retrieves the values associated with the provided key from the cache.
// The operation is thread-safe and does not alter the change history.
func (c *InMemoryHashMapMultiCache[K, T, H]) Get(key K) []T {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	return c.values[c.toHash(key.Keys())]
}

// Changes returns a list of keys that have experienced changes in the cache since the last reset.
func (c *InMemoryHashMapMultiCache[K, T, H]) Changes() []K {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()

	return umap.Values(c.changes)
}

// Drop completely clears the cache, removing all entries. The operation is thread-safe.
func (c *InMemoryHashMapMultiCache[K, T, H]) Drop() {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	c.dropAll()
	c.lastUpdatedKeys = make(map[string]keyContainer[K])
}

// DropKey removes the values associated with the provided key from the cache. The operation is thread-safe.
func (c *InMemoryHashMapMultiCache[K, T, H]) DropKey(key K) {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()
	hash := c.dropKey(key.Keys())
	delete(c.lastUpdatedKeys, keysAsString(key.Keys()))
	delete(c.changes, hash)
}

// Outdated checks if the provided key or the entire cache (if no key is provided)
// is outdated based on the set TTL. Returns true if outdated, false otherwise.
func (c *InMemoryHashMapMultiCache[K, T, H]) Outdated(key uopt.Opt[K]) bool {
	c.vMtx.Lock()
	defer c.vMtx.Unlock()

	if c.ttl == nil {
		return false
	} else {
		if key.Present() {
			k := key.Get()
			if lu, ok := c.lastUpdatedKeys[keysAsString((*k).Keys())]; ok {
				return time.Since(lu.updatedAt) > *c.ttl
			} else {
				return true
			}
		} else {
			return time.Since(c.lastUpdated) > *c.ttl
		}
	}
}

func (c *InMemoryHashMapMultiCache[K, T, H]) dropAll() {
	c.values = make(map[H][]T)
	c.changes = nil
}

func (c *InMemoryHashMapMultiCache[K, T, H]) put(key K, values ...T) {
	hash := c.addTran(key, values...)
	changes := len(c.changes) == 0
	found := false
	for _, diff := range c.changes {
		if uarray.EqualsWithOrder(diff.Keys(), key.Keys()) {
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

func (c *InMemoryHashMapMultiCache[K, T, H]) addTran(key K, values ...T) H {
	hash := c.toHash(key.Keys())
	c.values[hash] = append(c.values[hash], values...)

	return hash
}

func (c *InMemoryHashMapMultiCache[K, T, H]) dropKey(keys []uconst.Unique) H {
	hash := c.toHash(keys)
	delete(c.values, c.toHash(keys))
	return hash
}

func intToBytes(buffer *bytes.Buffer, num int64) []byte {
	buffer.Reset()
	_ = binary.Write(buffer, binary.LittleEndian, num)

	return buffer.Bytes()
}

func keysAsString(keys []uconst.Unique) string {
	var sb strings.Builder
	for _, key := range keys {
		sb.WriteString(strconv.FormatInt(key.Key(), 10))
	}
	return sb.String()
}
