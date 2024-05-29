/*
 * @Author: kordax, 10/5/23, 6:29 PM
 */

package ucache

import (
	"runtime"
	"testing"
	"time"

	"github.com/kordax/basic-utils/uopt"
)

const (
	numItems = 1000
	stdDepth = 3
)

// prepareCacheIntKey populates the cache with IntKey items.
func prepareCacheIntKey(c MultiCache[IntCompositeKey, Comparable], num int64) []IntCompositeKey {
	keys := make([]IntCompositeKey, num)
	for i := int64(0); i < num; i++ {
		key := NewIntCompositeKey(i)
		keys[i] = key
		c.Put(key, NewInt64Value(i))
	}

	return keys
}

// prepareCacheIntKey populates the cache with IntKey items.
func prepareCacheIntKeyWithDepth(c MultiCache[IntCompositeKey, Comparable], num, maxDepth int64) []IntCompositeKey {
	keys := make([]IntCompositeKey, num)
	for i := int64(0); i < num; i++ {
		var hashes []int64
		for h := int64(0); h < maxDepth; h++ {
			hashes = append(hashes, h)
		}
		key := NewIntCompositeKey(hashes...)
		keys[i] = key
		c.Put(key, NewInt64Value(i))
	}

	return keys
}

func BenchmarkSha256HashMapMultiCachePutIntKeySingle(b *testing.B) {
	c := NewSha256HashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := make([]IntCompositeKey, int64(b.N))
	for i := int64(0); i < int64(b.N); i++ {
		var hashes []int64
		for h := int64(0); h < stdDepth; h++ {
			hashes = append(hashes, h)
		}
		key := NewIntCompositeKey(hashes...)
		keys[i] = key
	}
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Put(keys[i], NewInt64Value(i))
	}
}

func BenchmarkSha256HashMapMultiCachePutIntKeySingleDepth100(b *testing.B) {
	c := NewSha256HashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := make([]IntCompositeKey, b.N)
	for i := int64(0); i < int64(b.N); i++ {
		var hashes []int64
		for h := int64(0); h < 100; h++ {
			hashes = append(hashes, h)
		}
		key := NewIntCompositeKey(hashes...)
		keys[i] = key
	}
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Put(keys[i], NewInt64Value(i))
	}
}

func BenchmarkSha256HashMapMultiCachePutIntKeyConcurrent(b *testing.B) {
	c := NewSha256HashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := make([]IntCompositeKey, b.N)
	for i := int64(0); i < int64(b.N); i++ {
		var hashes []int64
		for h := int64(0); h < stdDepth; h++ {
			hashes = append(hashes, h)
		}
		key := NewIntCompositeKey(hashes...)
		keys[i] = key
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Put(keys[int64(b.N-1)], NewInt64Value(int64(b.N)))
		}
	})
}

func BenchmarkSha256HashMapMultiCachePutIntKeyConcurrentDepth100(b *testing.B) {
	c := NewSha256HashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := make([]IntCompositeKey, b.N+1)
	for i := int64(0); i < int64(b.N); i++ {
		var hashes []int64
		for h := int64(0); h < 100; h++ {
			hashes = append(hashes, h)
		}
		key := NewIntCompositeKey(hashes...)
		keys[i] = key
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Put(keys[int64(b.N-1)], NewInt64Value(int64(b.N)))
		}
	})
}

func BenchmarkSha256HashMapMultiCacheGetIntKeySingle(b *testing.B) {
	c := NewSha256HashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%numItems])
	}
}

func BenchmarkSha256HashMapMultiCacheGetIntKeySingleDeepDepth100(b *testing.B) {
	c := NewSha256HashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKeyWithDepth(c, numItems, 100)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%numItems])
	}
}

func BenchmarkSha256HashMapMultiCacheGetIntKeyConcurrent(b *testing.B) {
	c := NewSha256HashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get(keys[int64(b.N)%numItems])
		}
	})
}

func BenchmarkSha256HashMapMultiCacheSetIntKeySingle(b *testing.B) {
	c := NewSha256HashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Set(keys[i%numItems], NewInt64Value(i))
	}
}

func BenchmarkSha256HashMapMultiCacheSetIntKeyConcurrent(b *testing.B) {
	c := NewSha256HashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set(keys[int64(b.N)%numItems], NewInt64Value(int64(b.N)))
		}
	})
}

func BenchmarkSha256HashMapMultiCacheGetIntKeySingle10xItems(b *testing.B) {
	num := int64(numItems * 10)
	c := NewSha256HashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, num)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%num])
	}
}

func BenchmarkFarmHashMapMultiCachePutIntKeySingle(b *testing.B) {
	c := NewFarmHashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := make([]IntCompositeKey, int64(b.N))
	for i := int64(0); i < int64(b.N); i++ {
		var hashes []int64
		for h := int64(0); h < stdDepth; h++ {
			hashes = append(hashes, h)
		}
		key := NewIntCompositeKey(hashes...)
		keys[i] = key
	}
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Put(keys[i], NewInt64Value(i))
	}
}

func BenchmarkFarmHashMapMultiCachePutIntKeySingleDepth100(b *testing.B) {
	c := NewFarmHashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := make([]IntCompositeKey, b.N)
	for i := int64(0); i < int64(b.N); i++ {
		var hashes []int64
		for h := int64(0); h < 100; h++ {
			hashes = append(hashes, h)
		}
		key := NewIntCompositeKey(hashes...)
		keys[i] = key
	}
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Put(keys[i], NewInt64Value(i))
	}
}

func BenchmarkFarmHashMapMultiCachePutIntKeyConcurrent(b *testing.B) {
	c := NewFarmHashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := make([]IntCompositeKey, b.N)
	for i := int64(0); i < int64(b.N); i++ {
		var hashes []int64
		for h := int64(0); h < stdDepth; h++ {
			hashes = append(hashes, h)
		}
		key := NewIntCompositeKey(hashes...)
		keys[i] = key
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Put(keys[int64(b.N-1)], NewInt64Value(int64(b.N)))
		}
	})
}

func BenchmarkFarmHashMapMultiCachePutIntKeyConcurrentDepth100(b *testing.B) {
	c := NewFarmHashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := make([]IntCompositeKey, b.N+1)
	for i := int64(0); i < int64(b.N); i++ {
		var hashes []int64
		for h := int64(0); h < 100; h++ {
			hashes = append(hashes, h)
		}
		key := NewIntCompositeKey(hashes...)
		keys[i] = key
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Put(keys[int64(b.N-1)], NewInt64Value(int64(b.N)))
		}
	})
}

func BenchmarkFarmHashMapMultiCacheGetIntKeySingle(b *testing.B) {
	c := NewFarmHashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%numItems])
	}
}

func BenchmarkFarmHashMapMultiCacheGetIntKeySingleDeepDepth100(b *testing.B) {
	c := NewFarmHashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKeyWithDepth(c, numItems, 100)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%numItems])
	}
}

func BenchmarkFarmHashMapMultiCacheGetIntKeyConcurrent(b *testing.B) {
	c := NewFarmHashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get(keys[int64(b.N)%numItems])
		}
	})
}

func BenchmarkFarmHashMapMultiCacheSetIntKeySingle(b *testing.B) {
	c := NewFarmHashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Set(keys[i%numItems], NewInt64Value(i))
	}
}

func BenchmarkFarmHashMapMultiCacheSetIntKeyConcurrent(b *testing.B) {
	c := NewFarmHashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set(keys[int64(b.N)%numItems], NewInt64Value(int64(b.N)))
		}
	})
}

func BenchmarkFarmHashMapMultiCacheGetIntKeySingle10xItems(b *testing.B) {
	num := int64(numItems * 10)
	c := NewFarmHashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, num)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%num])
	}
}

func BenchmarkTreeMultiCachePutIntKeySingle(b *testing.B) {
	c := NewInMemoryTreeMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := make([]IntCompositeKey, int64(b.N))
	for i := int64(0); i < int64(b.N); i++ {
		var hashes []int64
		for h := int64(0); h < stdDepth; h++ {
			hashes = append(hashes, h)
		}
		key := NewIntCompositeKey(hashes...)
		keys[i] = key
	}
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Put(keys[i], NewInt64Value(i))
	}
}

func BenchmarkTreeMultiCachePutIntKeySingleDepth100(b *testing.B) {
	c := NewInMemoryTreeMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := make([]IntCompositeKey, b.N)
	for i := int64(0); i < int64(b.N); i++ {
		var hashes []int64
		for h := int64(0); h < 100; h++ {
			hashes = append(hashes, h)
		}
		key := NewIntCompositeKey(hashes...)
		keys[i] = key
	}
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Put(keys[i], NewInt64Value(i))
	}
}

func BenchmarkTreeMultiCachePutIntKeyConcurrent(b *testing.B) {
	c := NewInMemoryTreeMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := make([]IntCompositeKey, b.N)
	for i := int64(0); i < int64(b.N); i++ {
		var hashes []int64
		for h := int64(0); h < stdDepth; h++ {
			hashes = append(hashes, h)
		}
		key := NewIntCompositeKey(hashes...)
		keys[i] = key
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Put(keys[int64(b.N-1)], NewInt64Value(int64(b.N)))
		}
	})
}

func BenchmarkTreeMultiCachePutIntKeyConcurrentDepth100(b *testing.B) {
	c := NewInMemoryTreeMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := make([]IntCompositeKey, b.N+1)
	for i := int64(0); i < int64(b.N); i++ {
		var hashes []int64
		for h := int64(0); h < 100; h++ {
			hashes = append(hashes, h)
		}
		key := NewIntCompositeKey(hashes...)
		keys[i] = key
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Put(keys[int64(b.N-1)], NewInt64Value(int64(b.N)))
		}
	})
}

func BenchmarkTreeMultiCacheGetIntKeySingle(b *testing.B) {
	c := NewInMemoryTreeMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%numItems])
	}
}

func BenchmarkTreeMultiCacheGetIntKeySingleDeepDepth100(b *testing.B) {
	c := NewInMemoryTreeMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKeyWithDepth(c, numItems, 100)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%numItems])
	}
}

func BenchmarkTreeMultiCacheGetIntKeyConcurrent(b *testing.B) {
	c := NewInMemoryTreeMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get(keys[int64(b.N)%numItems])
		}
	})
}

func BenchmarkTreeMultiCacheSetIntKeySingle(b *testing.B) {
	c := NewInMemoryTreeMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Set(keys[i%numItems], NewInt64Value(i))
	}
}

func BenchmarkTreeMultiCacheSetIntKeyConcurrent(b *testing.B) {
	c := NewInMemoryTreeMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set(keys[int64(b.N)%numItems], NewInt64Value(int64(b.N)))
		}
	})
}

func BenchmarkTreeMultiCacheGetIntKeySingle10xItems(b *testing.B) {
	num := int64(numItems * 10)
	c := NewInMemoryTreeMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, num)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%num])
	}
}

func BenchmarkMemoryFarmHashMapMultiCache(b *testing.B) {
	c := NewFarmHashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	benchmarkMemoryUsage(b, c)
}

func BenchmarkMemorySha256HashMapMultiCache(b *testing.B) {
	c := NewSha256HashMapMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	benchmarkMemoryUsage(b, c)
}

func BenchmarkMemoryTreeMultiCache(b *testing.B) {
	c := NewInMemoryTreeMultiCache[IntCompositeKey, Comparable](uopt.Null[time.Duration]())
	benchmarkMemoryUsage(b, c)
}

func benchmarkMemoryUsage(b *testing.B, c MultiCache[IntCompositeKey, Comparable]) {
	var m runtime.MemStats

	for i := int64(0); i < numItems; i++ {
		key := NewIntCompositeKey(i)
		c.Put(key, NewInt64Value(i))
	}

	runtime.ReadMemStats(&m)
	b.Logf("Memory Alloc = %v KB", bToKb(m.Alloc))
}

// Helper function to convert bytes to KB.
func bToKb(b uint64) uint64 {
	return b / 1024
}
