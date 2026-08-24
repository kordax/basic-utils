// Package ustream provides ordered, reusable lazy pipelines over Go values.
//
// Constructors and intermediate operations do not traverse their source. A
// terminal operation such as Collect, Count, Find, or ForEach starts traversal.
// Stateless stages such as Map, Filter, FlatMap, Limit, and TakeWhile pass
// values directly to the next stage without allocating intermediate slices.
//
// Sort, Reverse, and ParallelMap are buffering barriers: they remain lazy, but
// materialize their complete upstream when a terminal operation reaches them.
// ParallelMap preserves encounter order and may call its mapper concurrently,
// so the mapper must be safe for concurrent use.
//
// Streams built from slices are replayable. A stream built with FromSeq is
// replayable only when the supplied iter.Seq is itself replayable. This differs
// deliberately from Java streams, which are normally single-use.
package ustream
