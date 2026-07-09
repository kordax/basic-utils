# uconst

Shared constraints and small common type contracts used by the utility packages.

Use this package when building generic helpers that need numeric, comparable, stringable, or unique-key constraints.

## Example

```go
func ClampMin[T uconst.Numeric](value, min T) T {
	if value < min {
		return min
	}
	return value
}
```
