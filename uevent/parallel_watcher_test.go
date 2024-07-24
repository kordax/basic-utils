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

	"github.com/kordax/basic-utils/uarray"
	"github.com/kordax/basic-utils/uevent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParallelWatcher(t *testing.T) {
	inputCh := make(chan int)
	expectedResults := uarray.Range(0, 10)
	expectedCount := len(expectedResults)

	var wg sync.WaitGroup
	var m sync.Mutex

	// Register the function that will receive messages
	var received []int
	wg.Add(expectedCount)
	watcherFunc := func(ctx context.Context, msg *int) {
		defer wg.Done()
		m.Lock()
		defer m.Unlock()
		received = append(received, *msg)
	}

	watcher := uevent.NewSingleListenerWatcher(inputCh, watcherFunc)

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

	// Wait for all listeners to finish
	wg.Wait()

	assert.ElementsMatch(t, expectedResults, received, "Received messages do not match the expected results")
}

func TestParallelWatcherWithNoMessages(t *testing.T) {
	inputCh := make(chan int)
	expectedResults := []int{}
	expectedCount := len(expectedResults)

	var wg sync.WaitGroup

	// Register the function that will receive messages
	var received []int
	wg.Add(expectedCount)
	watcherFunc := func(ctx context.Context, msg *int) {
		defer wg.Done()
		received = append(received, *msg)
	}

	watcher := uevent.NewSingleListenerWatcher(inputCh, watcherFunc)

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the watcher
	started := watcher.Watch(ctx)
	require.True(t, started, "Watcher should have started successfully")

	// Ensure the watcher cannot be started again
	startedAgain := watcher.Watch(ctx)
	require.False(t, startedAgain, "Watcher should not start again")

	// No messages to send, just close the channel
	close(inputCh)

	// Wait for all listeners to finish
	wg.Wait()

	assert.ElementsMatch(t, expectedResults, received, "Received messages do not match the expected results")
}

func TestParallelWatcherContextCancel(t *testing.T) {
	inputCh := make(chan int)
	expectedResults := uarray.Range(0, 10)
	expectedCount := len(expectedResults)

	var wg sync.WaitGroup

	// Register the function that will receive messages
	var received []int
	wg.Add(expectedCount)
	watcherFunc := func(ctx context.Context, msg *int) {
		defer wg.Done()
		received = append(received, *msg)
	}

	watcher := uevent.NewSingleListenerWatcher(inputCh, watcherFunc)

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

	assert.ElementsMatch(t, []int{}, received, "Received messages do not match the expected results")
}
