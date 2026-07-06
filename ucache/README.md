# ucache

In-memory cache implementations with optional TTL support, composite keys, multi-cache structures, and managed cleanup wrappers.

The package includes comparable-key caches, hash-map caches, tree-based multi-cache variants, and composite key/value helpers.

Use `InMemoryComparableMapCache` for the fastest hot path when keys are regular comparable Go values. Use `GetValue`
when pointer allocation from `Get` is not needed.

Optional Ristretto comparison benchmarks live in a separate module under `benchmarks/ristretto`, so Ristretto is not a
dependency of the library itself.

Managed wrappers can be configured with `ManagedCacheOptions` to set cleanup interval and TTL in one place:

```go
cache := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())
managed := ucache.NewManagedCacheWithOptions(cache, ucache.ManagedCacheOptions{
	CleanupInterval: time.Second,
	TTL:             uopt.Of(time.Minute),
})
defer managed.Stop()
```
