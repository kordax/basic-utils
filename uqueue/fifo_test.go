/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uqueue

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFIFOQueueBasic(t *testing.T) {
	q := NewFIFOQueue[int]()

	// Test initial state.
	if v := q.Fetch(); v.Present() {
		t.Fatalf("Expected empty queue, got %v", v.Get())
	}

	// Test Queue and Fetch.
	q.Queue(1)
	if v := q.Fetch(); !v.Present() || v.OrElse(0) != 1 {
		t.Fatalf("Expected 1, got %v", v.OrElse(-1))
	}

	// Test Queue multiple and Fetch.
	q.Queue(2)
	q.Queue(3)
	q.Queue(4)
	if v := q.Fetch(); !v.Present() || v.OrElse(0) != 2 {
		t.Fatalf("Expected 2, got %v", v.OrElse(-1))
	}
	if v := q.Fetch(); !v.Present() || v.OrElse(0) != 3 {
		t.Fatalf("Expected 3, got %v", v.OrElse(-1))
	}

	// Test Poll with timeout.
	if v := q.Poll(10 * time.Millisecond); !v.Present() || v.OrElse(0) != 4 {
		t.Fatalf("Expected 4 from Poll, got %v", v.OrElse(-1))
	}
	if v := q.Poll(10 * time.Millisecond); v.Present() {
		t.Fatalf("Expected empty result from Poll, got %v", v.OrElse(-1))
	}
}

func TestFIFOQueueImpl_Len(t *testing.T) {
	q := NewFIFOQueue[int]()
	const n = 1000
	const n2 = 2347

	for i := range n {
		q.Queue(i)
	}
	assert.EqualValues(t, n, q.Len())

	for i := range n2 {
		q.Queue(i)
	}
	assert.EqualValues(t, n+n2, q.Len())

	for i := 0; i < n+n2; i++ {
		assert.True(t, q.Poll(time.Second*5).Present())
	}

	assert.EqualValues(t, 0, q.Len())
	assert.False(t, q.Fetch().Present())
}

func TestFIFOQueueCopiesInitialElements(t *testing.T) {
	values := []int{1, 2, 3}
	q := NewFIFOQueue(values...)

	values[0] = 99

	assert.Equal(t, 1, q.Fetch().OrElse(0))
	assert.Equal(t, []int{99, 2, 3}, values)
}

func TestFIFOQueuePollWaitsForQueuedItem(t *testing.T) {
	q := NewFIFOQueue[int]()
	result := make(chan int, 1)

	go func() {
		result <- q.Poll(time.Second).OrElse(-1)
	}()

	time.Sleep(10 * time.Millisecond)
	q.Queue(42)

	require.Eventually(t, func() bool {
		return len(result) == 1
	}, time.Second, time.Millisecond)
	assert.Equal(t, 42, <-result)
	assert.EqualValues(t, 0, q.Len())
}

func TestFIFOQueueConcurrentQueueAndPoll(t *testing.T) {
	q := NewFIFOQueue[int]()
	const n = 1000
	results := make(chan int, n)
	var wg sync.WaitGroup

	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if v := q.Poll(time.Second); v.Present() {
				results <- v.OrElse(-1)
			}
		}()
	}

	for i := range n {
		q.Queue(i)
	}

	wg.Wait()
	close(results)

	seen := make(map[int]struct{}, n)
	for v := range results {
		if _, exists := seen[v]; exists {
			t.Fatalf("duplicate value fetched: %d", v)
		}
		seen[v] = struct{}{}
	}
	require.Len(t, seen, n)
	assert.EqualValues(t, 0, q.Len())
}
