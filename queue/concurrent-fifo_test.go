/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package queue

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentFIFOQueueRaceConditions(t *testing.T) {
	q := NewConcurrentFIFOQueueImpl[int]()
	const n = 1000

	var wg sync.WaitGroup
	for i := range n {
		q.Queue(i)
	}

	// Concurrently fetch numbers.
	results := make(chan int, n)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range n {
			if v := q.Fetch(); v.Present() {
				results <- v.OrElse(-i)
			}
		}
		close(results)
	}()

	// Wait for both goroutines to finish.
	wg.Wait()

	// Verify the results.
	seen := make(map[int]bool)
	for r := range results {
		if seen[r] {
			t.Fatalf("Duplicate value fetched: %v", r)
		}
		seen[r] = true
	}
	if len(seen) != n {
		t.Fatalf("Expected %v unique numbers, got %v", n, len(seen))
	}
}

func TestConcurrentFIFOQueueConcurrentPolls(t *testing.T) {
	q := NewConcurrentFIFOQueueImpl[int]()
	const n = 1000

	var wg sync.WaitGroup
	for i := range n {
		q.Queue(i)
		time.Sleep(time.Microsecond) // Simulate some workload.
	}

	// Concurrently poll numbers.
	results := make(chan int, n)
	for _ = range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if v := q.Poll(1 * time.Minute); v.Present() {
				results <- v.OrElse(0)
			}
		}()
	}

	// Wait for both goroutines to finish.
	wg.Wait()
	close(results)

	// Verify the results.
	seen := make(map[int]bool)
	for r := range results {
		if seen[r] {
			t.Fatalf("Duplicate value fetched: %v", r)
		}
		seen[r] = true
	}
	if len(seen) != n {
		t.Fatalf("Expected %v unique numbers, got %v", n, len(seen))
	}
}

func TestConcurrentFIFOQueueImpl_Len(t *testing.T) {
	q := NewConcurrentFIFOQueueImpl[int]()
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
