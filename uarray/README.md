# uarray

Slice utilities for searching, filtering, mapping, grouping, set-like operations, safe indexing, chunking, and string conversion.

Use this package when the standard `slices` package is too low-level for common collection workflows.

## Example

```go
values := []int{1, 2, 3, 4}
even := uarray.Filter(values, func(v int) bool { return v%2 == 0 })
sum := uarray.Reduce(values, 0, func(acc, v int) int { return acc + v })
```
