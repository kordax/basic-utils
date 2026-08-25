# ustream

`ustream` provides lazy, ordered stream helpers for fluent collection processing on Go 1.27.

Use a stream when the pipeline is clearer than a direct loop or a sequence of `uarray` calls:

```go
result := ustream.From(5, 1, 2, 2, 3, 4).
	Filter(func(v int) bool {
		return v >= 2
	}).
	DistinctBy(func(v int) int {
		return v
	}).
	Map(func(v int) string {
		return strconv.Itoa(v * 10)
	}).
	Limit(2).
	Collect()
```

Generic methods allow type-changing operations such as `Map`, `FlatMap`, `MapMulti`, `Reduce`, `ToMap`, and `ToMultiMap` to remain fluent. The old package-level functions remain only as deprecated compatibility wrappers.

## Execution model

Intermediate operations are lazy. They run when a terminal operation such as `Collect`, `Reduce`, `Find`, or `Count` traverses the stream.

Stateless stages are composed without intermediate slices, and short-circuiting operations stop upstream traversal when possible. `Sort`, `Reverse`, and `ParallelMap` are lazy buffering barriers: they do no work while the pipeline is being built, but must materialize their input when traversed.

Streams created from slices are replayable. A stream created with `FromSeq` inherits the replayability and side effects of the supplied `iter.Seq`.

## Sources and iterator interoperability

Create streams with:

- `Empty`, `From`, and `Of` for ordinary values and slices.
- `FromSeq` for an existing `iter.Seq`.
- `Generate` for a fixed number of supplier calls.
- `Iterate` for a fixed number of successive values.

Use `Seq` to expose a stream as a standard iterator.

## Operations

Common intermediate operations include:

- Selection: `Filter`, `FilterOut`, `Limit`, `Skip`, `TakeWhile`, and `DropWhile`.
- Mapping: `Map`, `Transform`, `FlatMap`, and `MapMulti`.
- Composition and observation: `Concat` and `Peek`.
- Stateful operations: `DistinctBy`, `CompactFunc`, `Sort`, and `Reverse`.
- Parallel mapping: `ParallelMap`.

Common terminal operations include:

- Search and matching: `Find`, `FindFirst`, `FindIndex`, `AnyMatch`, `AllMatch`, and `NoneMatch`.
- Aggregation: `Count`, `Min`, `Max`, and `Reduce`.
- Materialization: `Collect`, `CollectCopy`, `ToMap`, `ToMultiMap`, and `CollectWith`.
- Side effects: `ForEach` and `ToTerminal`.

Collector helpers include `ToSlice`, `ToSliceCopy`, `ToMap`, `ToMultiMap`, `GroupingBy`, and `CountingBy`.

## MapMulti

Use `MapMulti` when one input may emit zero or more outputs without allocating a result slice for every input:

```go
values := ustream.From(1, 2, 3).
	MapMulti(func(v int, emit func(string) bool) {
		if !emit(strconv.Itoa(v)) {
			return
		}
		if v%2 == 0 {
			emit("even")
		}
	}).
	Collect()
```

The mapper must call `emit` synchronously, must not retain it, and should stop when `emit` returns `false`.

## Parallel execution

`ParallelMap` preserves encounter order and bounds the number of workers:

```go
result := ustream.From(values...).
	ParallelMap(func(v Input) Output {
		return expensiveTransform(v)
	}, 4).
	Collect()
```

Use it only when each mapping operation is expensive enough to repay scheduling and buffering costs. A plain loop or sequential `Map` is normally faster for cheap work. The mapper runs concurrently and must be safe to call from multiple goroutines.

`ParallelMap` materializes all upstream values before starting workers, so downstream short-circuiting cannot avoid that upstream work. If a mapper panics, it stops scheduling new work, waits for already running workers, and re-panics with the original value in the calling goroutine. No partial mapped result is emitted downstream.

`TerminalStream.ParallelExecute` has the same concurrency and panic-propagation behavior, but mutations completed before a panic are not rolled back.

Context belongs to a terminal execution rather than a reusable pipeline. Use `ParallelExecuteContext` when scheduling must stop on cancellation or a deadline:

```go
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()

err := terminal.ParallelExecuteContext(ctx, func(ctx context.Context, index int, value *Item) {
	process(ctx, index, value)
}, 4)
```

Cancellation stops new work and is passed to running callbacks. The method waits for those callbacks and returns `context.Cause(ctx)`, so callbacks must observe the context for prompt shutdown. A callback panic still takes precedence and is re-panicked in the calling goroutine. Regular `Stream` pipelines intentionally remain context-free.

## Errors and panics

A regular stream has no error channel. Handle expected errors inside the callback, represent them explicitly as values, or use ordinary Go control flow when error handling is central to the operation.

Sequential operations do not recover panics from sources or callbacks. They propagate normally from the terminal operation and can be handled by the application's existing recovery boundary. Passing a nil callback or violating an `iter.Seq` or `MapMulti` contract may also panic.
