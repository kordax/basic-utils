/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uasync

import (
	"context"
	"fmt"
	"sync"

	"github.com/kordax/basic-utils/v2/uarray"
)

// AsyncTask represents an asynchronous task that can be executed in the background.
// It provides mechanisms to start, wait for completion, and cancel the task.
// Additionally, it allows for retrying the task a specified number of times in case of failure.
type AsyncTask[R any] struct {
	fn      func(ctx context.Context) (*R, error) // The function that represents the async task.
	retries int                                   // Number of times to retry the task on failure. Value less than 0 means retry forever.

	f    Future[R]               // A future object to represent the result or error of the task.
	done *uarray.Pair[*R, error] // A pair containing the result and error of the task.
	mtx  sync.RWMutex            // A mutex for thread-safe operations on the struct.

	ctx context.Context
}

// NewAsyncTask creates a new instance of AsyncTask.
func NewAsyncTask[R any](ctx context.Context, fn func(ctx context.Context) (*R, error), retries int) *AsyncTask[R] {
	return &AsyncTask[R]{ctx: ctx, fn: fn, retries: retries, f: NewFuture[R](ctx)}
}

// ExecuteAsync initiates the execution of the task in a separate goroutine.
// If the task fails, it will be retried up to the specified number of times.
// Once the task completes or fails after all retries, the result or error is stored internally.
func (t *AsyncTask[R]) ExecuteAsync() {
	go func() {
		resultChan := make(chan *R)
		errChan := make(chan error)

		r, err := tryTask(t.ctx, t.fn, 0, t.retries)

		go func() {
			select {
			case <-t.ctx.Done():
				// If context is canceled, set the task as canceled
				t.cancel()
				t.mtx.Lock()
				t.done = uarray.NewPair[*R, error](nil, t.ctx.Err())
				t.mtx.Unlock()

			case r = <-resultChan:
				// If the task completes successfully
				t.f.Complete(r)
				t.mtx.Lock()
				t.done = uarray.NewPair[*R, error](r, nil)
				t.mtx.Unlock()

			case err = <-errChan:
				// If the task fails after all retries
				t.f.Fail(err)
				t.mtx.Lock()
				t.done = uarray.NewPair[*R, error](nil, err)
				t.mtx.Unlock()
			}
		}()

		if err != nil {
			errChan <- err
		} else {
			resultChan <- r
		}
	}()
}

// Wait allows the caller to wait for the task to complete or fail.
// If the task has already completed or failed, it returns the result or error immediately.
func (t *AsyncTask[R]) Wait() (*R, error) {
	t.mtx.RLock()
	if t.done != nil { // Check if the task is already done.
		t.mtx.RUnlock()
		return t.done.Left, t.done.Right
	}
	t.mtx.RUnlock()

	r, err := t.f.Wait()
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if t.done == nil { // Check if t.done is still nil to avoid overwriting by other calls to Wait()
		t.done = uarray.NewPair(r, err)
	}

	return r, err
}

// Cancel attempts to cancel the execution of the task.
// It invokes the provided cancellation function and marks the task as canceled.
func (t *AsyncTask[R]) cancel() {
	t.f.Cancel()
	t.mtx.Lock()
	t.done = uarray.NewPair[*R, error](nil, context.Canceled) // Set the done pair to represent the cancellation.
	t.mtx.Unlock()
}

// tryTask tries to execute a task function up to a maximum number of times.
// If max < 0, it retries forever. Context cancellation is respected.
func tryTask[R any](ctx context.Context, fn func(ctx context.Context) (*R, error), try int, max int) (*R, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			r, err := fn(ctx)
			if err == nil {
				return r, nil
			}

			if max >= 0 && try >= max {
				return nil, fmt.Errorf("attempt %d has failed: %w", try, err)
			}
			try++
		}
	}
}
