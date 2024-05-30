//go:build integration_test

/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ucache

import (
	"testing"
	"time"

	"github.com/kordax/basic-utils/uopt"
	"github.com/stretchr/testify/assert"
)

func TestManagedCache_TTL(t *testing.T) {
	ttl := 500 * time.Millisecond
	cache := NewInMemoryHashMapCache[IntKey, string, uint64](func(key int64) uint64 {
		return uint64(key)
	}, uopt.Of(ttl))
	managedCache := NewManagedCache(cache, uopt.Of(ttl))
	defer managedCache.Stop()

	key := IntKey(1)
	value := "TestValue"

	managedCache.Set(key, value)

	time.Sleep(2 * ttl)
	managedCache.performCleanup()

	v, ok := managedCache.Get(key)
	assert.False(t, ok)
	assert.Nil(t, v)

	managedCache.Set(key, value)
	v, ok = managedCache.Get(key)
	assert.True(t, ok)
	assert.NotNil(t, v)

	time.Sleep(2 * ttl)
	managedCache.performCleanup()

	// Key should not be removed now as it was requested before cleanupByGet and it's ttl was updated
	v, ok = managedCache.Get(key)
	assert.True(t, ok)
	assert.NotNil(t, v)
}
