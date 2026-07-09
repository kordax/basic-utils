# upair

Pair type helpers.

Use this package when returning or carrying two related values without introducing a local struct.

## Example

Use `Of` for the zero-allocation value form:

```go
pair := upair.Of("key", 42)
left, right := pair.Values()
```

Use `NewPair` only when a pointer is required. Use `COf` or `NewCPair` when the pair must satisfy `uconst.Comparable`.
