# ucache

In-memory cache implementations with optional TTL support, composite keys, multi-cache structures, and managed cleanup wrappers.

The package includes comparable-key caches, hash-map caches, tree-based multi-cache variants, and composite key/value helpers.

Use `InMemoryComparableMapCache` for the fastest hot path when keys are regular comparable Go values. Use `GetValue`
when pointer allocation from `Get` is not needed.

Use `InMemoryBufferedComparableMapCache` for write-heavy paths where immediate visibility is not required. It has the
same cache methods, but `Set` and `SetQuietly` use buffered admission writes by default:

## Example

```go
var cache ucache.ComparableCache[string, int] = ucache.NewInMemoryBufferedComparableMapCacheWithOptions[string, int](ucache.InMemoryComparableMapCacheOptions{
	TTL:               uopt.Null[time.Duration](),
	BufferedWorkers:   4,
	BufferedQueueSize: 65536,
	BufferedMaxKeys:   10000,
})
```

When deterministic visibility or shutdown is needed, keep the concrete value and call `Wait` or `CloseBuffered`:

```go
cache := ucache.NewInMemoryBufferedComparableMapCacheWithOptions[string, int](ucache.InMemoryComparableMapCacheOptions{
	BufferedMaxKeys: 10000,
})
defer cache.CloseBuffered()

cache.Set("key", 42)
cache.Wait()
```

`BufferedMaxKeys` enables bounded/admission-style behavior: new keys above the limit can be ignored by `Set` and
`SetQuietly`. Use `InMemoryComparableMapCache` when every write must be applied synchronously.

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
