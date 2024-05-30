/*
* @kordax (Dmitry Morozov)
* dmorozov@valoru-software.com
* Copyright (c) 2024.
 */

package ucache

//import (
//	"bytes"
//	"sync"
//	"time"
//
//	"github.com/dgryski/go-farm"
//	"github.com/kordax/basic-utils/uopt"
//)
//
//// The CompositeCache is the same as Cache, but allows CompositeKeys
//type CompositeCache[K CompositeKey, T any] interface {
//	Cache[K, T]
//}
//
//// InMemoryHashMapCompositeCache composite version of InMemoryHashMapCompositeCache
//type InMemoryHashMapCompositeCache[K Unique, T any, H comparable] struct {
//	values  map[H]T
//	changes []K
//
//	lastUpdatedKeys map[int64]time.Time
//	lastUpdated     time.Time
//	ttl             *time.Duration
//
//	toHash func(key int64) H
//	vMtx   sync.Mutex
//}
//
//// NewInMemoryHashMapCompositeCache creates a new instance of the InMemoryHashMapCompositeCache.
//func NewInMemoryHashMapCompositeCache[K Unique, T any, H comparable](toHash func(key int64) H, ttl uopt.Opt[time.Duration]) Cache[K, T] {
//	c := &InMemoryHashMapCompositeCache[K, T, H]{
//		values:          make(map[H]T),
//		changes:         make([]K, 0),
//		lastUpdatedKeys: make(map[int64]time.Time),
//		toHash:          toHash,
//	}
//	ttl.IfPresent(func(t time.Duration) {
//		c.ttl = &t
//	})
//
//	return c
//}
//
//// NewDefaultHashMapCache creates a new instance of the InMemoryHashMapCache using SHA256 as the hashing algorithm.
//func NewDefaultHashMapCache[K Unique, T any](ttl uopt.Opt[time.Duration]) Cache[K, T] {
//	return NewFarmHashMapCache[K, T](ttl)
//}
//
//func NewFarmHashMapCache[K Unique, T any](ttl uopt.Opt[time.Duration]) Cache[K, T] {
//	return NewInMemoryHashMapCache[K, T, uint64](func(key int64) uint64 {
//		buffer := new(bytes.Buffer)
//		return farm.Hash64(intToBytes(buffer, key))
//	}, ttl)
//}
//
//// Set updates the cache value for the provided key. If the key already exists,
//// its previous value are removed before adding the new value. The operation is thread-safe.
//func (c *InMemoryHashMapCompositeCache[K, T, H]) Set(key K, value T) {
//	c.vMtx.Lock()
//	defer c.vMtx.Unlock()
//	c.put(key, value)
//	c.lastUpdatedKeys[key.Key()] = time.Now()
//	c.lastUpdated = time.Now()
//}
//
//// SetQuietly is an optimized method that adds value to the cache for the provided key but does so without
//// altering the change history. This operation can be used when modifications should not trigger cache change diff.
//// This operation is much faster and can be used to optimize cache performance in case you don't want to track changes.
//func (c *InMemoryHashMapCompositeCache[K, T, H]) SetQuietly(key K, value T) {
//	c.vMtx.Lock()
//	defer c.vMtx.Unlock()
//	c.addTran(key, value)
//	c.lastUpdatedKeys[key.Key()] = time.Now()
//	c.lastUpdated = time.Now()
//}
//
//// Get retrieves the value associated with the provided key from the cache.
//// The operation is thread-safe and does not alter the change history.
//func (c *InMemoryHashMapCompositeCache[K, T, H]) Get(key K) (*T, bool) {
//	c.vMtx.Lock()
//	defer c.vMtx.Unlock()
//	c.changes = nil
//
//	value, ok := c.values[c.toHash(key.Key())]
//	if !ok {
//		return nil, false
//	}
//	return &value, ok
//}
//
//// Drop completely clears the cache, removing all entries. The operation is thread-safe.
//func (c *InMemoryHashMapCompositeCache[K, T, H]) Drop() {
//	c.vMtx.Lock()
//	defer c.vMtx.Unlock()
//	c.dropAll()
//	c.lastUpdatedKeys = make(map[int64]time.Time)
//}
//
//// DropKey removes the value associated with the provided key from the cache. The operation is thread-safe.
//func (c *InMemoryHashMapCompositeCache[K, T, H]) DropKey(key K) {
//	c.vMtx.Lock()
//	defer c.vMtx.Unlock()
//	c.dropKey(key.Key())
//	c.lastUpdatedKeys[key.Key()] = time.Now()
//	c.lastUpdated = time.Now()
//}
//
//// Outdated checks if the provided key or the entire cache (if no key is provided)
//// is outdated based on the set TTL. Returns true if outdated, false otherwise.
//func (c *InMemoryHashMapCompositeCache[K, T, H]) Outdated(key uopt.Opt[K]) bool {
//	c.vMtx.Lock()
//	defer c.vMtx.Unlock()
//
//	if c.ttl == nil {
//		return false
//	} else {
//		if key.Present() {
//			k := key.Get()
//			if lu, ok := c.lastUpdatedKeys[(*k).Key()]; ok {
//				return time.Since(lu) > *c.ttl
//			} else {
//				return true
//			}
//		} else {
//			return time.Since(c.lastUpdated) > *c.ttl
//		}
//	}
//}
//
//func (c *InMemoryHashMapCompositeCache[K, T, H]) dropAll() {
//	c.values = make(map[H]T)
//	c.changes = nil
//}
//
//func (c *InMemoryHashMapCompositeCache[K, T, H]) put(key K, value T) {
//	c.addTran(key, value)
//	changes := len(c.changes) == 0
//	found := false
//	for _, diff := range c.changes {
//		if diff.Key() == key.Key() {
//			if !diff.Equals(key) {
//				changes = true
//				break
//			}
//			found = true
//			continue
//		}
//	}
//	if changes || !found {
//		c.changes = append(c.changes, key)
//	}
//}
//
//func (c *InMemoryHashMapCompositeCache[K, T, H]) addTran(key K, value T) {
//	c.values[c.toHash(key.Key())] = value
//}
//
//func (c *InMemoryHashMapCompositeCache[K, T, H]) dropKey(key int64) {
//	delete(c.values, c.toHash(key))
//}
