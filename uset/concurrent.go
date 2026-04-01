/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uset

import (
	"sync"

	"git.casinomodule.org/kordax/basic-utils/usrlz"
	"github.com/dgryski/go-farm"
)

type concurrentShard[T comparable] struct {
	mu sync.RWMutex
	m  map[T]struct{}
}

// ConcurrentHashSet is a concurrent-safe hash set implementation that uses sharding.
// Each shard is protected by its own RWMutex which allows high concurrency for
// operations on different keys.
type ConcurrentHashSet[T comparable] struct {
	hash   func(T) uint64
	shards []concurrentShard[T]
}

const defaultConcurrentHashSetShards = 32

func newConcurrentHashSetWithHash[T comparable](hash func(T) uint64, shardCount int) *ConcurrentHashSet[T] {
	if shardCount <= 0 {
		shardCount = defaultConcurrentHashSetShards
	}

	shards := make([]concurrentShard[T], shardCount)
	for i := range shards {
		shards[i].m = make(map[T]struct{})
	}

	return &ConcurrentHashSet[T]{
		hash:   hash,
		shards: shards,
	}
}

// NewConcurrentHashSet creates a ConcurrentHashSet with a default Farm64-based hash
// and a default number of shards.
func NewConcurrentHashSet[T comparable]() *ConcurrentHashSet[T] {
	defaultHash := func(value T) uint64 {
		return farm.Hash64(usrlz.ToBytes(&value))
	}

	return newConcurrentHashSetWithHash[T](defaultHash, defaultConcurrentHashSetShards)
}

// NewConcurrentHashSetWithHash creates a ConcurrentHashSet with a custom hash
// function and a default number of shards.
func NewConcurrentHashSetWithHash[T comparable](hash func(T) uint64) *ConcurrentHashSet[T] {
	return newConcurrentHashSetWithHash[T](hash, defaultConcurrentHashSetShards)
}

// NewConcurrentHashSetWithHashAndShards creates a ConcurrentHashSet with a custom
// hash function and a custom number of shards.
func NewConcurrentHashSetWithHashAndShards[T comparable](hash func(T) uint64, shardCount int) *ConcurrentHashSet[T] {
	return newConcurrentHashSetWithHash[T](hash, shardCount)
}

// NewCustomConcurrentHashSet is kept for backward compatibility.
// It is equivalent to NewConcurrentHashSetWithHash with the default shard count.
func NewCustomConcurrentHashSet[T comparable](hash func(T) uint64) *ConcurrentHashSet[T] {
	return NewConcurrentHashSetWithHash(hash)
}

func (s *ConcurrentHashSet[T]) shardFor(value T) *concurrentShard[T] {
	if len(s.shards) == 0 {
		return nil
	}

	idx := s.hash(value) % uint64(len(s.shards))
	return &s.shards[idx]
}

// Add inserts a value into the set and returns true if the value was not already present.
func (s *ConcurrentHashSet[T]) Add(value T) bool {
	sh := s.shardFor(value)
	if sh == nil {
		return false
	}

	sh.mu.Lock()
	defer sh.mu.Unlock()

	if sh.m == nil {
		sh.m = make(map[T]struct{})
	}

	if _, exists := sh.m[value]; exists {
		return false
	}

	sh.m[value] = struct{}{}
	return true
}

// Contains checks if a value is present in the set.
func (s *ConcurrentHashSet[T]) Contains(value T) bool {
	sh := s.shardFor(value)
	if sh == nil {
		return false
	}

	sh.mu.RLock()
	defer sh.mu.RUnlock()

	if sh.m == nil {
		return false
	}

	_, exists := sh.m[value]
	return exists
}

// Remove deletes a value from the set and returns true if the value was present.
func (s *ConcurrentHashSet[T]) Remove(value T) bool {
	sh := s.shardFor(value)
	if sh == nil {
		return false
	}

	sh.mu.Lock()
	defer sh.mu.Unlock()

	if sh.m == nil {
		return false
	}

	if _, exists := sh.m[value]; exists {
		delete(sh.m, value)
		return true
	}

	return false
}

// Size returns the number of elements in the set.
func (s *ConcurrentHashSet[T]) Size() int {
	total := 0
	for i := range s.shards {
		sh := &s.shards[i]
		sh.mu.RLock()
		total += len(sh.m)
		sh.mu.RUnlock()
	}

	return total
}

// Clear removes all elements from the set.
func (s *ConcurrentHashSet[T]) Clear() {
	for i := range s.shards {
		sh := &s.shards[i]
		sh.mu.Lock()
		sh.m = make(map[T]struct{})
		sh.mu.Unlock()
	}
}

// Values retrieves all values from the set.
// The snapshot is not atomic across shards but is safe to call concurrently.
func (s *ConcurrentHashSet[T]) Values() []T {
	// First, estimate the total size.
	total := 0
	for i := range s.shards {
		sh := &s.shards[i]
		sh.mu.RLock()
		total += len(sh.m)
		sh.mu.RUnlock()
	}

	result := make([]T, 0, total)
	for i := range s.shards {
		sh := &s.shards[i]
		sh.mu.RLock()
		for v := range sh.m {
			result = append(result, v)
		}
		sh.mu.RUnlock()
	}

	return result
}
