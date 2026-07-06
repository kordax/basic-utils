package ristretto_benchmark

import (
	"strconv"
	"testing"
	"time"

	"git.casinomodule.org/casino27/basic-utils/v2/ucache"
	"git.casinomodule.org/casino27/basic-utils/v2/uopt"
	"github.com/dgraph-io/ristretto/v2"
)

const numItems = 10000

func newRistrettoCache(b testing.TB) *ristretto.Cache[string, int] {
	b.Helper()

	cache, err := ristretto.NewCache(&ristretto.Config[string, int]{
		NumCounters: numItems * 10,
		MaxCost:     numItems,
		BufferItems: 64,
	})
	if err != nil {
		b.Fatal(err)
	}

	return cache
}

func seedUCache() (*ucache.InMemoryComparableMapCache[string, int], []string) {
	cache := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())
	keys := make([]string, numItems)
	for i := 0; i < numItems; i++ {
		key := "key-" + strconv.Itoa(i)
		keys[i] = key
		cache.SetQuietly(key, i)
	}

	return cache, keys
}

func seedRistretto(b testing.TB) (*ristretto.Cache[string, int], []string) {
	cache := newRistrettoCache(b)
	keys := make([]string, numItems)
	for i := 0; i < numItems; i++ {
		key := "key-" + strconv.Itoa(i)
		keys[i] = key
		cache.Set(key, i, 1)
	}
	cache.Wait()

	return cache, keys
}

func BenchmarkUCacheComparableSetQuietly(b *testing.B) {
	cache := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.SetQuietly("key-"+strconv.Itoa(i), i)
	}
}

func BenchmarkUCacheBufferedComparableSetBounded(b *testing.B) {
	cache := ucache.NewInMemoryBufferedComparableMapCacheWithOptions[string, int](ucache.InMemoryComparableMapCacheOptions{
		BufferedMaxKeys: int64(numItems),
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key-"+strconv.Itoa(i), i)
	}
	cache.Wait()
	cache.CloseBuffered()
}

func BenchmarkUCacheBufferedComparableSetQuietlyBounded(b *testing.B) {
	cache := ucache.NewInMemoryBufferedComparableMapCacheWithOptions[string, int](ucache.InMemoryComparableMapCacheOptions{
		BufferedMaxKeys: int64(numItems),
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.SetQuietly("key-"+strconv.Itoa(i), i)
	}
	cache.Wait()
	cache.CloseBuffered()
}

func BenchmarkRistrettoSet(b *testing.B) {
	cache := newRistrettoCache(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key-"+strconv.Itoa(i), i, 1)
	}
	cache.Wait()
}

func BenchmarkUCacheComparableGetValue(b *testing.B) {
	cache, keys := seedUCache()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetValue(keys[i%numItems])
	}
}

func BenchmarkRistrettoGet(b *testing.B) {
	cache, keys := seedRistretto(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(keys[i%numItems])
	}
}

func BenchmarkUCacheComparableGetValueParallel(b *testing.B) {
	cache, keys := seedUCache()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.GetValue(keys[i%numItems])
			i++
		}
	})
}

func BenchmarkRistrettoGetParallel(b *testing.B) {
	cache, keys := seedRistretto(b)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Get(keys[i%numItems])
			i++
		}
	})
}
