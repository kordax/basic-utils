# ustream

Stream-style helpers for fluent slice processing.

Use this package when a pipeline is clearer than direct `uarray` calls:

```go
result := ustream.From(5, 1, 2, 2, 3, 4).
	Filter(func(v int) bool {
		return v >= 2
	}).
	DistinctBy(func(v int) string {
		return strconv.Itoa(v)
	}).
	Transform(func(v int) int {
		return v * 10
	}).
	Sort(func(a, b int) bool {
		return a < b
	}).
	Skip(1).
	Limit(2).
	Reverse().
	Collect()
```

For same-type chains use methods: `Filter`, `FilterOut`, `Transform`, `FlatMap`, `Peek`, `Limit`, `Skip`, `TakeWhile`, `DropWhile`, `Reverse`, `Sort`, `DistinctBy`, `CompactFunc`, `Find`, `FindIndex`, `AnyMatch`, `AllMatch`, `NoneMatch`, `Reduce`, `Collect`.

For type-changing operations Go does not allow method-level type parameters, so use typed package functions:

```go
stream := ustream.From(1, 2, 3)

mapped := ustream.Map(stream, func(v int) string {
	return fmt.Sprintf("n-%d", v)
})

grouped := ustream.Collect(mapped, ustream.GroupingBy(func(v string) int {
	return len(v)
}))
```

Collectors are available through `Collect`, `ToSlice`, `ToSliceCopy`, `ToMap`, `ToMultiMap`, `GroupingBy`, and `CountingBy`.
