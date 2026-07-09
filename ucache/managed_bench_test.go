package ucache_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kordax/basic-utils/v3/ucache"
	"github.com/kordax/basic-utils/v3/uopt"
)

func BenchmarkManagedCacheForceCleanupFreshKeys(b *testing.B) {
	cache := ucache.NewInMemoryComparableMapCache[string, int](uopt.Of(time.Hour))
	managed := ucache.NewManagedCache[string, int](cache, time.Hour)
	defer managed.Stop()

	for i := 0; i < 10_000; i++ {
		managed.Set(fmt.Sprintf("key-%d", i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		managed.ForceCleanup()
	}
}

func BenchmarkManagedCacheForceCleanupExpiredKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cache := ucache.NewInMemoryComparableMapCache[string, int](uopt.Of(time.Nanosecond))
		managed := ucache.NewManagedCache[string, int](cache, time.Hour)

		for j := 0; j < 1_000; j++ {
			managed.Set(fmt.Sprintf("key-%d", j), j)
		}

		managed.ForceCleanup()
		managed.Stop()
	}
}

func BenchmarkManagedMultiCacheForceCleanupFreshKeys(b *testing.B) {
	cache := ucache.NewDefaultHashMapMultiCache[ucache.IntCompositeKey, ucache.Int64Value](uopt.Of(time.Hour))
	managed := ucache.NewManagedMultiCache[ucache.IntCompositeKey, ucache.Int64Value](cache, time.Hour)
	defer managed.Stop()

	for i := int64(0); i < 10_000; i++ {
		managed.Put(ucache.NewIntCompositeKey(i), ucache.NewInt64Value(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		managed.ForceCleanup()
	}
}
