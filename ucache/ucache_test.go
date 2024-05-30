/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ucache_test

import (
	"sync"
	"testing"
	"time"

	"github.com/kordax/basic-utils/ucache"
	"github.com/kordax/basic-utils/uopt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashMapCache_CompositeKey(t *testing.T) {
	c := ucache.NewDefaultHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
	key := ucache.StringKey("category")
	key2 := ucache.StringKey("category2")
	val := 10
	val2 := 236261

	c.SetQuietly(key, val)
	c.SetQuietly(key2, val2)

	result, ok := c.Get(key)
	require.True(t, ok, "value was expected to be cached")
	result2, ok := c.Get(key2)
	require.True(t, ok, "value was expected to be cached")
	assert.EqualValues(t, val, *result)
	assert.EqualValues(t, val2, *result2)
}

func TestHashMapCache_PutQuietly(t *testing.T) {
	c := ucache.NewDefaultHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
	key := ucache.StringKey("kp_1")
	val := 10
	val2 := 15

	c.SetQuietly(key, val)
	c.SetQuietly(key, val)
	c.SetQuietly(key, val)

	result, ok := c.Get(key)
	require.True(t, ok, "value was expected to be cached")
	assert.Equal(t, *result, val)

	c.SetQuietly(key, val2)
	result, ok = c.Get(key)
	require.True(t, ok, "value was expected to be cached")
	assert.Equal(t, *result, val2)

	c.SetQuietly(key, val)
	result, ok = c.Get(key)
	require.True(t, ok, "value was expected to be cached")
	assert.Equal(t, *result, val)
}

func TestHashMapCache_TTLExpiry(t *testing.T) {
	ttl := 100 * time.Millisecond
	c := ucache.NewDefaultHashMapCache[ucache.StringKey, int](uopt.Of(ttl))
	key := ucache.StringKey("ttlKey")
	val := 42

	c.Set(key, val)
	time.Sleep(2 * ttl)
	outdated := c.Outdated(uopt.Of(key))
	assert.True(t, outdated, "key should be marked as outdated")
}

func TestHashMapCache_Concurrency(t *testing.T) {
	c := ucache.NewDefaultHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
	key := ucache.StringKey("concurrencyKey")
	val := 42

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Set(key, val+i)
			result, _ := c.Get(key)
			assert.NotNil(t, result)
		}(i)
	}
	wg.Wait()
}

func TestHashMapCache_EmptyCache(t *testing.T) {
	c := ucache.NewDefaultHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
	key := ucache.StringKey("emptyKey")

	_, ok := c.Get(key)
	assert.False(t, ok, "key should not be found in an empty cache")

	c.DropKey(key)
	_, ok = c.Get(key)
	assert.False(t, ok, "key should not be found in an empty cache")
}

func TestHashMapCache_DropAll(t *testing.T) {
	c := ucache.NewDefaultHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
	key := ucache.StringKey("key1")
	key2 := ucache.StringKey("key2")
	c.Set(key, 1)
	c.Set(key2, 2)

	c.Drop()
	_, ok1 := c.Get(key)
	_, ok2 := c.Get(key2)

	assert.False(t, ok1, "key1 should be dropped")
	assert.False(t, ok2, "key2 should be dropped")
}

func TestInMemoryHashMapCache(t *testing.T) {
	cache := ucache.NewInMemoryHashMapCache[ucache.IntKey, string, uint64](func(key int64) uint64 {
		return uint64(key)
	}, uopt.Null[time.Duration]())

	// Define multiple keys and values
	key1 := ucache.IntKey(1)
	value1 := "MyValue"

	key2 := ucache.IntKey(2)
	value2 := "AnotherValue"

	key3 := ucache.IntKey(3)
	value3 := "ThirdValue"

	key4 := ucache.IntKey(4)
	value4 := "FourthValue"

	key5 := ucache.IntKey(5)
	value5 := "FifthValue"

	// Test setting and getting multiple keys
	cache.Set(key1, value1)
	cache.Set(key2, value2)
	cache.Set(key3, value3)
	cache.Set(key4, value4)
	cache.Set(key5, value5)

	// Verify all keys return correct values
	retrievedValue, ok := cache.Get(key1)
	require.True(t, ok, "Expected to retrieve value for key1")
	assert.Equal(t, value1, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key2)
	require.True(t, ok, "Expected to retrieve value for key2")
	assert.Equal(t, value2, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key3)
	require.True(t, ok, "Expected to retrieve value for key3")
	assert.Equal(t, value3, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key4)
	require.True(t, ok, "Expected to retrieve value for key4")
	assert.Equal(t, value4, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key5)
	require.True(t, ok, "Expected to retrieve value for key5")
	assert.Equal(t, value5, *retrievedValue, "Retrieved value should match the set value")

	// Test updating values for existing keys
	updatedValue1 := "UpdatedMyValue"
	cache.Set(key1, updatedValue1)
	retrievedValue, ok = cache.Get(key1)
	require.True(t, ok, "Expected to retrieve updated value for key1")
	assert.Equal(t, updatedValue1, *retrievedValue, "Retrieved value should match the updated value")

	// Test removing keys
	cache.DropKey(key1)
	retrievedValue, ok = cache.Get(key1)
	assert.False(t, ok, "Expected key1 to be removed from cache")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key1 should be nil")

	// Ensure other keys are still retrievable and correct after removing key1
	retrievedValue, ok = cache.Get(key2)
	require.True(t, ok, "Expected to retrieve value for key2")
	assert.Equal(t, value2, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key3)
	require.True(t, ok, "Expected to retrieve value for key3")
	assert.Equal(t, value3, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key4)
	require.True(t, ok, "Expected to retrieve value for key4")
	assert.Equal(t, value4, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key5)
	require.True(t, ok, "Expected to retrieve value for key5")
	assert.Equal(t, value5, *retrievedValue, "Retrieved value should match the set value")

	// Test SetQuietly
	cache.SetQuietly(key1, updatedValue1)
	retrievedValue, ok = cache.Get(key1)
	require.True(t, ok, "Expected to retrieve value for key1 after SetQuietly")
	assert.Equal(t, updatedValue1, *retrievedValue, "Retrieved value should match the set value")

	// Test Drop (clearing the entire cache)
	cache.Drop()
	retrievedValue, ok = cache.Get(key1)
	assert.False(t, ok, "Expected key1 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key1 should be nil after Drop")

	retrievedValue, ok = cache.Get(key2)
	assert.False(t, ok, "Expected key2 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key2 should be nil after Drop")

	retrievedValue, ok = cache.Get(key3)
	assert.False(t, ok, "Expected key3 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key3 should be nil after Drop")

	retrievedValue, ok = cache.Get(key4)
	assert.False(t, ok, "Expected key4 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key4 should be nil after Drop")

	retrievedValue, ok = cache.Get(key5)
	assert.False(t, ok, "Expected key5 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key5 should be nil after Drop")
}

func TestHashMapCacheHighCollisionProbability(t *testing.T) {
	c := ucache.NewFarmHashMapCache[CollisionTestKey, ucache.Int64Value](uopt.Null[time.Duration]())

	// Define a set of keys that all produce the same hash code
	keys := []CollisionTestKey{
		{id: 1, hash: []int64{1, 2, 3}},
		{id: 2, hash: []int64{1, 2, 3}},
		{id: 3, hash: []int64{1, 2, 3}},
	}

	// Add values to the c for each key
	for i, key := range keys {
		c.Set(key, ucache.NewInt64Value(int64(i)))
	}

	// Ensure that all values can be retrieved despite the high collision probability
	for i, key := range keys {
		value, ok := c.Get(key)
		require.True(t, ok, "Expected to retrieve value for key")
		assert.EqualValues(t, ucache.NewInt64Value(int64(i)), *value, "Expected value for key %s", key.String())
	}
}
