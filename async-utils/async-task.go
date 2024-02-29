/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package asyncutils

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	arrayutils "github.com/kordax/basic-utils/array-utils"
)

// AsyncTask represents an asynchronous task that can be executed in the background.
// It provides mechanisms to start, wait for completion, and cancel the task.
// Additionally, it allows for retrying the task a specified number of times in case of failure.
type AsyncTask[R any] struct {
	fn         func() (*R, error) // The function that represents the async task.
	cancelFunc context.CancelFunc // Function to cancel the task.
	retries    int                // Number of times to retry the task on failure.

	f    Future[R]                   // A future object to represent the result or error of the task.
	done *arrayutils.Pair[*R, error] // A pair containing the result and error of the task.
	mtx  sync.RWMutex                // A mutex for thread-safe operations on the struct.
}

// NewAsyncTask creates a new instance of AsyncTask.
func NewAsyncTask[R any](fn func() (*R, error), cancelFunc context.CancelFunc, retries int) *AsyncTask[R] {
	f := NewFuture[R]()
	return &AsyncTask[R]{fn: fn, cancelFunc: cancelFunc, retries: retries, f: f}
}

// ExecuteAsync initiates the execution of the task in a separate goroutine.
// If the task fails, it will be retried up to the specified number of times.
// Once the task completes or fails after all retries, the result or error is stored internally.
func (t *AsyncTask[R]) ExecuteAsync() {
	go func() {
		r, err := tryTask(t.fn, 0, t.retries)
		if err != nil {
			t.f.Fail(err)
			t.mtx.Lock()
			t.done = arrayutils.NewPair[*R, error](nil, err)
			t.mtx.Unlock()
			return
		}

		t.f.Complete(r)
		t.mtx.Lock()
		t.done = arrayutils.NewPair[*R, error](r, nil)
		t.mtx.Unlock()
	}()
}

// Wait allows the caller to wait for the task to complete or fail.
// If the task has already completed or failed, it returns the result or error immediately.
func (t *AsyncTask[R]) Wait(timeout time.Duration) (*R, error) {
	t.mtx.RLock()
	if t.done != nil { // Check if the task is already done.
		return t.done.Left, t.done.Right
	}
	t.mtx.RUnlock()

	r, err := t.f.TimeWait(timeout)
	t.mtx.Lock()
	t.done = arrayutils.NewPair(r, err)
	t.mtx.Unlock()

	return r, err
}

// Cancel attempts to cancel the execution of the task.
// It invokes the provided cancellation function and marks the task as canceled.
func (t *AsyncTask[R]) Cancel() {
	if t.cancelFunc != nil {
		t.cancelFunc()
	}
	t.f.Cancel()
	t.mtx.Lock()
	t.done = arrayutils.NewPair[*R, error](nil, errors.New("cancelled")) // Set the done pair to represent the cancellation.
	t.mtx.Unlock()
}

// tryTask tries to execute a task function up to a maximum number of times.
func tryTask[R any](fn func() (*R, error), try int, max int) (*R, error) {
	r, err := fn() // Execute the task function.
	if err != nil {
		if try < max {
			return tryTask(fn, try+1, max)
		} else {
			return nil, fmt.Errorf("attempt %d has failed: %w", try, err)
		}
	}

	return r, nil
}
