/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uevent_test

import (
	"context"
	"sync"
	"testing"

	"github.com/kordax/basic-utils/v2/uarray"
	"github.com/kordax/basic-utils/v2/uevent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBroadcastWatcher(t *testing.T) {
	inputCh := make(chan int)
	watcher := uevent.NewBroadcastWatcher(inputCh)

	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup
	var m sync.Mutex

	expectedResults := uarray.Range(0, 10)
	expectedCount := len(expectedResults)

	// Register listener 1
	var received1 []int
	wg1.Add(expectedCount)
	watcher.Register(func(ctx context.Context, msg int) {
		defer wg1.Done()
		m.Lock()
		defer m.Unlock()
		received1 = append(received1, msg)
	})

	// Register listener 2
	var received2 []int
	wg2.Add(expectedCount)
	watcher.Register(func(ctx context.Context, msg int) {
		defer wg2.Done()
		m.Lock()
		defer m.Unlock()
		received2 = append(received2, msg)
	})

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	started := watcher.Watch(ctx)
	require.True(t, started, "Watcher should have started successfully")

	startedAgain := watcher.Watch(ctx)
	require.False(t, startedAgain, "Watcher should not start again")

	go func() {
		for i := 0; i <= expectedCount-1; i++ {
			inputCh <- i
		}
		close(inputCh)
	}()

	wg1.Wait()
	wg2.Wait()

	assert.ElementsMatch(t, expectedResults, received1, "listener1 received messages do not match the expected results")
	assert.ElementsMatch(t, expectedResults, received2, "listener2 received messages do not match the expected results")
}

func TestBroadcastWatcherWithNoListeners(t *testing.T) {
	inputCh := make(chan int)
	watcher := uevent.NewBroadcastWatcher(inputCh)

	expectedResults := uarray.Range(0, 10)
	expectedCount := len(expectedResults)

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the watcher
	started := watcher.Watch(ctx)
	require.True(t, started, "Watcher should have started successfully")

	// Ensure the watcher cannot be started again
	startedAgain := watcher.Watch(ctx)
	require.False(t, startedAgain, "Watcher should not start again")

	// Send messages to input channel
	go func() {
		for i := 0; i <= expectedCount-1; i++ {
			inputCh <- i
		}
		close(inputCh)
	}()
}

func TestBroadcastWatcherContextCancel(t *testing.T) {
	inputCh := make(chan int)
	watcher := uevent.NewBroadcastWatcher(inputCh)

	expectedResults := uarray.Range(0, 10)
	expectedCount := len(expectedResults)

	var wg sync.WaitGroup

	// Register listener 1
	var received1 []int
	wg.Add(expectedCount)
	watcher.Register(func(ctx context.Context, msg int) {
		defer wg.Done()
		received1 = append(received1, msg)
	})

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Start the watcher
	started := watcher.Watch(ctx)
	require.True(t, started, "Watcher should have started successfully")

	// Ensure the watcher cannot be started again
	startedAgain := watcher.Watch(ctx)
	require.False(t, startedAgain, "Watcher should not start again")

	// Cancel the context before sending messages
	cancel()

	// Send messages to input channel
	go func() {
		for i := 0; i <= expectedCount-1; i++ {
			inputCh <- i
		}
		close(inputCh)
	}()

	assert.ElementsMatch(t, []int{}, received1, "Listener 1 should not have received any messages due to context cancellation")
}
