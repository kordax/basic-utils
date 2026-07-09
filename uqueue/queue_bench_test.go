package uqueue

import (
	"context"
	"testing"
)

func BenchmarkFIFOQueueDrain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := NewFIFOQueue[int]()
		for v := 0; v < 1000; v++ {
			q.Queue(v)
		}
		_ = q.Drain()
	}
}

func BenchmarkPriorityQueueDrain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := NewPriorityQueue[int]()
		for v := 0; v < 1000; v++ {
			q.Queue(v, v)
		}
		_ = q.Drain()
	}
}

func BenchmarkConcurrentFIFOQueueDrain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := NewConcurrentFIFOQueueImpl[int]()
		for v := 0; v < 1000; v++ {
			q.Queue(v)
		}
		_ = q.Drain()
	}
}

func BenchmarkFIFOQueuePeek(b *testing.B) {
	q := NewFIFOQueue(1)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = q.Peek()
	}
}

func BenchmarkFIFOQueuePollContextReady(b *testing.B) {
	q := NewFIFOQueue[int]()
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		q.Queue(i)
		_ = q.PollContext(ctx)
	}
}
