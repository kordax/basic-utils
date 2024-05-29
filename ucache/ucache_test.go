/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ucache_test

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/kordax/basic-utils/ucache"
	"github.com/kordax/basic-utils/uopt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashMapCache(t *testing.T) {
	c := ucache.NewDefaultHashMapCache[SimpleCompositeKey[ucache.StringKey], int](uopt.Null[time.Duration]())
	key := NewSimpleCompositeKey[ucache.StringKey]("atest")
	key2 := NewSimpleCompositeKey[ucache.StringKey]("bSeCond-keyQ@!%!%#")
	val := 326
	c.Set(key, val)

	cached, ok := c.Get(key)
	require.True(t, ok, "value was expected to be cached")
	assert.Equal(t, val, *cached)

	for i := 0; i < 10; i++ {
		c.Set(key, i)
	}
	c.Set(key2, 65535)

	result, ok := c.Get(key2)
	require.True(t, ok, "value was expected to be cached")
	assert.Equal(t, *result, 65535)

	complexKeyBase := []ucache.StringKey{"p1", "p2", "p3"}
	partialComplexKey := NewSimpleCompositeKey[ucache.StringKey](complexKeyBase...)

	for i := 0; i < 10; i++ {
		complexKey := NewSimpleCompositeKey[ucache.StringKey](append(complexKeyBase, ucache.StringKey("number:"+strconv.Itoa(i)))...)
		c.Set(complexKey, i)
	}

	result, ok = c.Get(partialComplexKey)
	require.True(t, ok, "value was expected to be cached")
	assert.NotEmpty(t, result)
	assert.Equal(t, 9, *result)
}

func TestHashMapCache_CompositeKey(t *testing.T) {
	c := ucache.NewDefaultHashMapCache[ucache.StrCompositeKey, int](uopt.Null[time.Duration]())
	key := ucache.NewStrCompositeKey("category", "kp_2")
	key2 := ucache.NewStrCompositeKey("category2", "kp_2")
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

func TestHashMapCache_DropKey(t *testing.T) {
	c := ucache.NewDefaultHashMapCache[ucache.StrCompositeKey, int](uopt.Null[time.Duration]())
	categoryKey := ucache.NewStrCompositeKey("category")
	overlappingKey := ucache.NewStrCompositeKey("category", "kp_232626")
	key2 := ucache.NewStrCompositeKey("category2", "kp_232626")
	catVal := rand.Int()
	val := rand.Int()
	val2 := rand.Int()

	c.Set(categoryKey, catVal)
	c.Set(overlappingKey, val)
	c.Set(key2, val2)

	catRes, _ := c.Get(categoryKey)

	res, ok := c.Get(overlappingKey)
	require.True(t, ok, "value was expected to be cached")
	res2, ok := c.Get(key2)
	require.True(t, ok, "value was expected to be cached")
	assert.Equal(t, val, *catRes)
	assert.Equal(t, val, *res)
	assert.Equal(t, val2, *res2)

	c.DropKey(overlappingKey)
	catRes, ok = c.Get(categoryKey)
	require.True(t, ok, "value was expected to remain")
	assert.Equal(t, val, *catRes)
	_, ok = c.Get(overlappingKey)
	require.False(t, ok, "value was expected to be cleared out of the cache")
	res2, ok = c.Get(key2)
	require.True(t, ok, "value was expected to remain")
	assert.Equal(t, val2, *res2)
}

func TestHashMapCache_PutQuietly(t *testing.T) {
	c := ucache.NewDefaultHashMapCache[SimpleCompositeKey[ucache.StringKey], int](uopt.Null[time.Duration]())
	key := NewSimpleCompositeKey[ucache.StringKey]("kp_1", "kp_2")
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
	ttl := 1 * time.Second
	c := ucache.NewDefaultHashMapCache[SimpleCompositeKey[ucache.StringKey], int](uopt.Of(ttl))
	key := NewSimpleCompositeKey[ucache.StringKey]("ttlKey")
	val := 42

	c.Set(key, val)
	time.Sleep(2 * time.Second)
	outdated := c.Outdated(uopt.Of(key))
	assert.True(t, outdated, "key should be marked as outdated")
}

func TestHashMapCache_Concurrency(t *testing.T) {
	c := ucache.NewDefaultHashMapCache[SimpleCompositeKey[ucache.StringKey], int](uopt.Null[time.Duration]())
	key := NewSimpleCompositeKey[ucache.StringKey]("concurrencyKey")
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
	c := ucache.NewDefaultHashMapCache[SimpleCompositeKey[ucache.StringKey], int](uopt.Null[time.Duration]())
	key := NewSimpleCompositeKey[ucache.StringKey]("emptyKey")

	_, ok := c.Get(key)
	assert.False(t, ok, "key should not be found in an empty cache")

	c.DropKey(key)
	_, ok = c.Get(key)
	assert.False(t, ok, "key should not be found in an empty cache")
}

func TestHashMapCache_DropAll(t *testing.T) {
	c := ucache.NewDefaultHashMapCache[SimpleCompositeKey[ucache.StringKey], int](uopt.Null[time.Duration]())
	key := NewSimpleCompositeKey[ucache.StringKey]("key1")
	key2 := NewSimpleCompositeKey[ucache.StringKey]("key2")
	c.Set(key, 1)
	c.Set(key2, 2)

	c.Drop()
	_, ok1 := c.Get(key)
	_, ok2 := c.Get(key2)

	assert.False(t, ok1, "key1 should be dropped")
	assert.False(t, ok2, "key2 should be dropped")
}

func TestHashMapCache_PartialKeyMatch(t *testing.T) {
	c := ucache.NewDefaultHashMapCache[SimpleCompositeKey[ucache.StringKey], int](uopt.Null[time.Duration]())
	fullKey := NewSimpleCompositeKey[ucache.StringKey]("part1", "part2")
	partialKey := NewSimpleCompositeKey[ucache.StringKey]("part1")
	val := 123

	c.Set(fullKey, val)
	result, ok := c.Get(partialKey)
	assert.True(t, ok, "partial key should match")
	assert.Equal(t, val, *result)
}
