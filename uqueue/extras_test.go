package uqueue

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFIFOQueueExtras(t *testing.T) {
	q := NewFIFOQueue(1, 2, 3)

	require.True(t, q.Peek().Present())
	assert.Equal(t, 1, q.Peek().OrElse(0))
	assert.Equal(t, []int{1, 2}, q.Drain(2))
	assert.Equal(t, 3, q.Peek().OrElse(0))

	q.Clear()
	assert.True(t, q.Empty())
	assert.False(t, q.Peek().Present())
	assert.Empty(t, q.Drain())
}

func TestPriorityQueueExtras(t *testing.T) {
	q := NewPriorityQueue[int]()
	q.Queue(1, 1)
	q.Queue(3, 3)
	q.Queue(2, 2)

	assert.Equal(t, 3, q.Peek().OrElse(0))
	assert.Equal(t, []int{3, 2}, q.Drain(2))
	assert.Equal(t, 1, q.Peek().OrElse(0))

	q.Clear()
	assert.True(t, q.Empty())
	assert.False(t, q.Fetch().Present())
}

func TestConcurrentFIFOQueueExtras(t *testing.T) {
	q := NewConcurrentFIFOQueueImpl[int]()
	q.Queue(1)
	q.Queue(2)
	q.Queue(3)

	assert.Equal(t, 1, q.Peek().OrElse(0))
	assert.Equal(t, []int{1, 2}, q.Drain(2))
	assert.Equal(t, 3, q.Peek().OrElse(0))

	q.Clear()
	assert.True(t, q.Empty())
	assert.False(t, q.Fetch().Present())
}

func TestQueuePollContext(t *testing.T) {
	q := NewFIFOQueue[int]()
	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan int, 1)

	go func() {
		result <- q.PollContext(ctx).OrElse(-1)
	}()

	q.Queue(42)
	require.Eventually(t, func() bool { return len(result) == 1 }, time.Second, time.Millisecond)
	assert.Equal(t, 42, <-result)
	cancel()
}

func TestQueuePollContextCancellation(t *testing.T) {
	q := NewFIFOQueue[int]()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	assert.False(t, q.PollContext(ctx).Present())
}
