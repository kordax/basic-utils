/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ucache_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/kordax/basic-utils/ucache"
	"github.com/kordax/basic-utils/uopt"
)

func BenchmarkInMemoryHashMapCachePut(b *testing.B) {
	cache := ucache.NewDefaultHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
	keys := make([]ucache.StringKey, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = ucache.StringKey(fmt.Sprintf("key%d", i))
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Set(keys[i], i)
	}
}

func BenchmarkInMemoryHashMapCachePutConcurrent(b *testing.B) {
	cache := ucache.NewDefaultHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
	keys := make([]ucache.StringKey, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = ucache.StringKey(fmt.Sprintf("key%d", i))
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := keys[rand.Intn(b.N)]
			cache.Set(key, rand.Int())
		}
	})
}

func BenchmarkInMemoryHashMapCacheGet(b *testing.B) {
	numItems := 10000
	cache := ucache.NewDefaultHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
	keys := make([]ucache.StringKey, numItems)
	for i := 0; i < numItems; i++ {
		keys[i] = ucache.StringKey(fmt.Sprintf("key%d", i))
		cache.Set(keys[i], i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Get(keys[i%numItems])
	}
}

func BenchmarkInMemoryHashMapCacheGetConcurrent(b *testing.B) {
	numItems := 10000
	cache := ucache.NewDefaultHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
	keys := make([]ucache.StringKey, numItems)
	for i := 0; i < numItems; i++ {
		keys[i] = ucache.StringKey(fmt.Sprintf("key%d", i))
		cache.Set(keys[i], i)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := keys[rand.Intn(numItems)]
			cache.Get(key)
		}
	})
}
