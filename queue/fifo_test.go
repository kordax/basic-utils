/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
