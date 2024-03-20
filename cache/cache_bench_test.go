//go:build integration_test

/*
 * @Author: kordax, 10/5/23, 6:29 PM
 */

package cache

import (
	"testing"
	"time"

	"github.com/kordax/basic-utils/opt"
)

const (
	numItems = 1000
	stdDepth = 3
)

// prepareCacheIntKey populates the cache with IntKey items.
func prepareCacheIntKey(c Cache[IntCompositeKey, Comparable], num int64) []IntCompositeKey {
	keys := make([]IntCompositeKey, num)
	for i := int64(0); i < num; i++ {
		key := NewIntCompositeKey(i)
		keys[i] = key
		c.Put(key, NewInt64Value(i))
	}

	return keys
}

// prepareCacheIntKey populates the cache with IntKey items.
func prepareCacheIntKeyWithDepth(c Cache[IntCompositeKey, Comparable], num, maxDepth int64) []IntCompositeKey {
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

func BenchmarkSha256HashMapCachePutIntKeySingle(b *testing.B) {
	c := NewSha256HashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
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

func BenchmarkSha256HashMapCachePutIntKeySingleDepth100(b *testing.B) {
	c := NewSha256HashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
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

func BenchmarkSha256HashMapCachePutIntKeyConcurrent(b *testing.B) {
	c := NewSha256HashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
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

func BenchmarkSha256HashMapCachePutIntKeyConcurrentDepth100(b *testing.B) {
	c := NewSha256HashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
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

func BenchmarkSha256HashMapCacheGetIntKeySingle(b *testing.B) {
	c := NewSha256HashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%numItems])
	}
}

func BenchmarkSha256HashMapCacheGetIntKeySingleDeepDepth100(b *testing.B) {
	c := NewSha256HashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKeyWithDepth(c, numItems, 100)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%numItems])
	}
}

func BenchmarkSha256HashMapCacheGetIntKeyConcurrent(b *testing.B) {
	c := NewSha256HashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get(keys[int64(b.N)%numItems])
		}
	})
}

func BenchmarkSha256HashMapCacheSetIntKeySingle(b *testing.B) {
	c := NewSha256HashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Set(keys[i%numItems], NewInt64Value(i))
	}
}

func BenchmarkSha256HashMapCacheSetIntKeyConcurrent(b *testing.B) {
	c := NewSha256HashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set(keys[int64(b.N)%numItems], NewInt64Value(int64(b.N)))
		}
	})
}

func BenchmarkSha256HashMapCacheGetIntKeySingle10xItems(b *testing.B) {
	num := int64(numItems * 10)
	c := NewSha256HashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, num)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%num])
	}
}

func BenchmarkFarmHashMapCachePutIntKeySingle(b *testing.B) {
	c := NewFarmHashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
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

func BenchmarkFarmHashMapCachePutIntKeySingleDepth100(b *testing.B) {
	c := NewFarmHashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
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

func BenchmarkFarmHashMapCachePutIntKeyConcurrent(b *testing.B) {
	c := NewFarmHashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
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

func BenchmarkFarmHashMapCachePutIntKeyConcurrentDepth100(b *testing.B) {
	c := NewFarmHashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
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

func BenchmarkFarmHashMapCacheGetIntKeySingle(b *testing.B) {
	c := NewFarmHashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%numItems])
	}
}

func BenchmarkFarmHashMapCacheGetIntKeySingleDeepDepth100(b *testing.B) {
	c := NewFarmHashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKeyWithDepth(c, numItems, 100)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%numItems])
	}
}

func BenchmarkFarmHashMapCacheGetIntKeyConcurrent(b *testing.B) {
	c := NewFarmHashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get(keys[int64(b.N)%numItems])
		}
	})
}

func BenchmarkFarmHashMapCacheSetIntKeySingle(b *testing.B) {
	c := NewFarmHashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Set(keys[i%numItems], NewInt64Value(i))
	}
}

func BenchmarkFarmHashMapCacheSetIntKeyConcurrent(b *testing.B) {
	c := NewFarmHashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set(keys[int64(b.N)%numItems], NewInt64Value(int64(b.N)))
		}
	})
}

func BenchmarkFarmHashMapCacheGetIntKeySingle10xItems(b *testing.B) {
	num := int64(numItems * 10)
	c := NewFarmHashMapCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, num)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%num])
	}
}

func BenchmarkTreeCachePutIntKeySingle(b *testing.B) {
	c := NewInMemoryTreeCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
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

func BenchmarkTreeCachePutIntKeySingleDepth100(b *testing.B) {
	c := NewInMemoryTreeCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
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

func BenchmarkTreeCachePutIntKeyConcurrent(b *testing.B) {
	c := NewInMemoryTreeCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
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

func BenchmarkTreeCachePutIntKeyConcurrentDepth100(b *testing.B) {
	c := NewInMemoryTreeCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
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

func BenchmarkTreeCacheGetIntKeySingle(b *testing.B) {
	c := NewInMemoryTreeCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%numItems])
	}
}

func BenchmarkTreeCacheGetIntKeySingleDeepDepth100(b *testing.B) {
	c := NewInMemoryTreeCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKeyWithDepth(c, numItems, 100)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%numItems])
	}
}

func BenchmarkTreeCacheGetIntKeyConcurrent(b *testing.B) {
	c := NewInMemoryTreeCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get(keys[int64(b.N)%numItems])
		}
	})
}

func BenchmarkTreeCacheSetIntKeySingle(b *testing.B) {
	c := NewInMemoryTreeCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Set(keys[i%numItems], NewInt64Value(i))
	}
}

func BenchmarkTreeCacheSetIntKeyConcurrent(b *testing.B) {
	c := NewInMemoryTreeCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, numItems)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set(keys[int64(b.N)%numItems], NewInt64Value(int64(b.N)))
		}
	})
}

func BenchmarkTreeCacheGetIntKeySingle10xItems(b *testing.B) {
	num := int64(numItems * 10)
	c := NewInMemoryTreeCache[IntCompositeKey, Comparable](opt.Null[time.Duration]())
	keys := prepareCacheIntKey(c, num)
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		c.Get(keys[i%num])
	}
}
