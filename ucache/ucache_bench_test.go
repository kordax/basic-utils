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

	"git.casinomodule.org/casino27/basic-utils/v2/ucache"
	"git.casinomodule.org/casino27/basic-utils/v2/uopt"
)

func BenchmarkInMemoryHashMapCachePut(b *testing.B) {
	cache := ucache.NewInMemoryHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
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
	cache := ucache.NewInMemoryHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
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
	cache := ucache.NewInMemoryHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
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

func BenchmarkInMemoryHashMapCacheGetValue(b *testing.B) {
	numItems := 10000
	cache := ucache.NewInMemoryHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
	keys := make([]ucache.StringKey, numItems)
	for i := 0; i < numItems; i++ {
		keys[i] = ucache.StringKey(fmt.Sprintf("key%d", i))
		cache.Set(keys[i], i)
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.GetValue(keys[i%numItems])
	}
}

func BenchmarkInMemoryHashMapCacheGetConcurrent(b *testing.B) {
	numItems := 10000
	cache := ucache.NewInMemoryHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
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

func BenchmarkInMemoryComparableMapCacheSet(b *testing.B) {
	cache := ucache.NewInMemoryComparableMapCache[int, int](uopt.Null[time.Duration]())

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(i, i)
	}
}

func BenchmarkInMemoryComparableMapCacheSetQuietly(b *testing.B) {
	cache := ucache.NewInMemoryComparableMapCache[int, int](uopt.Null[time.Duration]())

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.SetQuietly(i, i)
	}
}

func BenchmarkInMemoryComparableMapCacheGet(b *testing.B) {
	const numItems = 10000
	cache := ucache.NewInMemoryComparableMapCache[int, int](uopt.Null[time.Duration]())
	for i := 0; i < numItems; i++ {
		cache.SetQuietly(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(i % numItems)
	}
}

func BenchmarkInMemoryComparableMapCacheGetValue(b *testing.B) {
	const numItems = 10000
	cache := ucache.NewInMemoryComparableMapCache[int, int](uopt.Null[time.Duration]())
	for i := 0; i < numItems; i++ {
		cache.SetQuietly(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetValue(i % numItems)
	}
}

func BenchmarkInMemoryComparableMapCacheGetValueConcurrent(b *testing.B) {
	const numItems = 10000
	cache := ucache.NewInMemoryComparableMapCache[int, int](uopt.Null[time.Duration]())
	for i := 0; i < numItems; i++ {
		cache.SetQuietly(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.GetValue(i % numItems)
			i++
		}
	})
}

func BenchmarkInMemoryComparableMapCacheGetConcurrent(b *testing.B) {
	const numItems = 10000
	cache := ucache.NewInMemoryComparableMapCache[int, int](uopt.Null[time.Duration]())
	for i := 0; i < numItems; i++ {
		cache.SetQuietly(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Get(i % numItems)
			i++
		}
	})
}
