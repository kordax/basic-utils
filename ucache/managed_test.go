/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ucache_test

import (
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/kordax/basic-utils/v2/ucache"
	"github.com/kordax/basic-utils/v2/uopt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManagedMultiCache_SetAndGet(t *testing.T) {
	cache := ucache.NewInMemoryTreeMultiCache[ucache.StrCompositeKey, DummyComparable](uopt.Null[time.Duration]())
	managedCache := ucache.NewManagedMultiCache(cache, time.Second)
	defer managedCache.Stop()

	key := ucache.NewStrCompositeKey("category", "key1")
	value := DummyComparable{Val: 42}

	managedCache.Set(key, value)
	results := managedCache.Get(key)
	assert.Len(t, results, 1)
	assert.Equal(t, value, results[0])
}

func TestManagedMultiCache_Drop(t *testing.T) {
	cache := ucache.NewInMemoryTreeMultiCache[ucache.StrCompositeKey, DummyComparable](uopt.Null[time.Duration]())
	managedCache := ucache.NewManagedMultiCache(cache, time.Second)
	defer managedCache.Stop()

	key := ucache.NewStrCompositeKey("category", "key1")
	value := DummyComparable{Val: 42}

	managedCache.Set(key, value)
	managedCache.Drop()
	results := managedCache.Get(key)
	assert.Empty(t, results)
}

func TestManagedMultiCache_DropKey(t *testing.T) {
	cache := ucache.NewInMemoryTreeMultiCache[ucache.StrCompositeKey, DummyComparable](uopt.Null[time.Duration]())
	managedCache := ucache.NewManagedMultiCache(cache, time.Second)
	defer managedCache.Stop()

	key := ucache.NewStrCompositeKey("category", "key1")
	value := DummyComparable{Val: 42}

	managedCache.Set(key, value)
	managedCache.DropKey(key)
	results := managedCache.Get(key)
	assert.Empty(t, results)
}

func TestManagedMultiCache_PutQuietly(t *testing.T) {
	cache := ucache.NewInMemoryTreeMultiCache[ucache.StrCompositeKey, DummyComparable](uopt.Null[time.Duration]())
	managedCache := ucache.NewManagedMultiCache(cache, time.Second)
	defer managedCache.Stop()

	key := ucache.NewStrCompositeKey("category", "key1")
	value := DummyComparable{Val: 42}

	managedCache.PutQuietly(key, value)
	results := managedCache.Get(key)
	assert.Len(t, results, 1)
	assert.Equal(t, value, results[0])
}

func TestManagedMultiCache_Outdated(t *testing.T) {
	ttl := 100 * time.Millisecond
	cache := ucache.NewInMemoryTreeMultiCache[ucache.StrCompositeKey, DummyComparable](uopt.Of(ttl))
	managedCache := ucache.NewManagedMultiCache(cache, time.Nanosecond)
	defer managedCache.Stop()

	key := ucache.NewStrCompositeKey("category", "key1")
	value := DummyComparable{Val: 42}

	managedCache.Set(key, value)
	time.Sleep(2 * ttl)
	values := managedCache.Get(key)
	assert.Empty(t, values)
}

func TestManagedCache_SetAndGet(t *testing.T) {
	cache := ucache.NewInMemoryHashMapCache[ucache.IntKey, string](uopt.Null[time.Duration]())
	managedCache := ucache.NewManagedCache(cache, time.Second)
	defer managedCache.Stop()

	key := ucache.IntKey(1)
	value := "TestValue"

	managedCache.Set(key, value)
	v, ok := managedCache.Get(key)
	assert.True(t, ok)
	assert.Equal(t, value, *v)
}

func TestManagedCache_Drop(t *testing.T) {
	cache := ucache.NewInMemoryHashMapCache[ucache.IntKey, string](uopt.Null[time.Duration]())
	managedCache := ucache.NewManagedCache(cache, time.Second)
	defer managedCache.Stop()

	key := ucache.IntKey(1)
	value := "TestValue"

	managedCache.Set(key, value)
	managedCache.Drop()
	v, ok := managedCache.Get(key)
	assert.False(t, ok)
	assert.Nil(t, v)
}

func TestManagedCache_DropKey(t *testing.T) {
	cache := ucache.NewInMemoryHashMapCache[ucache.IntKey, string](uopt.Null[time.Duration]())
	managedCache := ucache.NewManagedCache(cache, time.Second)
	defer managedCache.Stop()

	key := ucache.IntKey(1)
	value := "TestValue"

	managedCache.Set(key, value)
	managedCache.DropKey(key)
	v, ok := managedCache.Get(key)
	assert.False(t, ok)
	assert.Nil(t, v)
}

func TestManagedCache_SetQuietly(t *testing.T) {
	cache := ucache.NewInMemoryHashMapCache[ucache.IntKey, string](uopt.Null[time.Duration]())
	managedCache := ucache.NewManagedCache(cache, time.Second)
	defer managedCache.Stop()

	key := ucache.IntKey(1)
	value := "TestValue"

	managedCache.SetQuietly(key, value)
	v, ok := managedCache.Get(key)
	assert.True(t, ok)
	assert.Equal(t, value, *v)
}

func TestManagedCache_Outdated(t *testing.T) {
	ttl := 1 * time.Millisecond
	cache := ucache.NewInMemoryHashMapCache[ucache.IntKey, string](uopt.Of(ttl))
	managedCache := ucache.NewManagedCache(cache, time.Nanosecond)
	defer managedCache.Stop()

	key := ucache.IntKey(1)
	value := "TestValue"

	managedCache.Set(key, value)
	time.Sleep(10 * ttl)
	_, ok := managedCache.Get(key)
	assert.False(t, ok)
}

func TestManagedCache_MemoryLeaks(t *testing.T) {
	ttl := time.Nanosecond
	cache := ucache.NewInMemoryHashMapCache[ucache.IntKey, string](uopt.Of(ttl))
	managedCache := ucache.NewManagedCache(cache, time.Nanosecond)
	defer managedCache.Stop()

	iterations := 100000

	before := new(runtime.MemStats)
	runtime.ReadMemStats(before)
	for i := range iterations {
		managedCache.Set(ucache.IntKey(i), "TestValue"+strconv.Itoa(i))
	}

	for i := range iterations {
		_, ok := managedCache.Get(ucache.IntKey(i))
		require.False(t, ok)
	}

	runtime.GC()
	after := new(runtime.MemStats)
	runtime.ReadMemStats(after)

	assert.LessOrEqual(t, after.HeapAlloc, before.HeapAlloc*3)
}
