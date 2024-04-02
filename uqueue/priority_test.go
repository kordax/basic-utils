/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uqueue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPriorityQueueBasic(t *testing.T) {
	pq := NewPriorityQueue[int]()

	if pq.Poll(10 * time.Millisecond).Present() {
		t.Fatal("Expected Poll to return nil for empty queue")
	}

	pq.Queue(5, 2)
	if val := pq.Fetch(); !val.Present() || val.OrElse(-1) != 5 {
		t.Fatalf("Expected 5, got %v", val.OrElse(-1))
	}

	if pq.Poll(10 * time.Millisecond).Present() {
		t.Fatal("Expected Poll to return nil after fetching all items")
	}
}

func TestPriorityQueuePriority(t *testing.T) {
	pq := NewPriorityQueue[int]()

	pq.Queue(5, 2)
	pq.Queue(10, 3)
	pq.Queue(15, 1)

	if val := pq.Fetch(); !val.Present() || val.OrElse(-1) != 10 {
		t.Fatalf("Expected 10, got %v", val.OrElse(-1))
	}

	if val := pq.Fetch(); !val.Present() || val.OrElse(-1) != 5 {
		t.Fatalf("Expected 5, got %v", val.OrElse(-1))
	}

	if val := pq.Fetch(); !val.Present() || val.OrElse(-1) != 15 {
		t.Fatalf("Expected 15, got %v", val.OrElse(-1))
	}
}

func TestPriorityQueue10k(t *testing.T) {
	pq := NewPriorityQueue[int]()
	const n = 10000

	// Concurrently enqueuing
	for i := range n {
		pq.Queue(i, i)
	}

	// Concurrently dequeuing
	fetched := make([]bool, n)
	for range n {
		item := pq.Poll(10 * time.Second)
		if !item.Present() {
			t.Fatal("Unexpected nil during Poll")
		}
		fetched[item.OrElse(-1)] = true
	}

	// Ensure every item from 0 to n-1 was fetched
	for i := range n {
		if !fetched[i] {
			t.Fatalf("Did not fetch %d from the queue", i)
		}
	}
}

func TestPriorityQueueImpl_Len(t *testing.T) {
	q := NewPriorityQueue[int]()
	const n = 1000
	const n2 = 2347

	for i := range n {
		q.Queue(i, i)
	}
	assert.EqualValues(t, n, q.Len())

	for i := range n2 {
		q.Queue(i, i)
	}
	assert.EqualValues(t, n+n2, q.Len())

	for i := 0; i < n+n2; i++ {
		assert.True(t, q.Poll(time.Second*5).Present())
	}

	assert.EqualValues(t, 0, q.Len())
	assert.False(t, q.Fetch().Present())
}
