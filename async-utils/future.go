/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package async_utils

import (
	"sync"
	"sync/atomic"
	"time"
)

// Future is an interface that represents a value or an error that will be available in the future.
// It provides mechanisms to wait for the value, set the value or error, and cancel the operation.
type Future[T any] interface {
	Wait() (*T, error)
	TimeWait(timeout time.Duration) (*T, error)
	Cancel() bool
	Complete(t *T) bool
	Fail(err error) bool
}

// FutureImpl is an implementation of the Future interface.
// It uses condition variables to allow waiting for the value or error to be set.
type FutureImpl[T any] struct {
	cond *sync.Cond
	once sync.Once

	v   atomic.Value
	err atomic.Value

	completed bool
}

func NewFuture[T any]() *FutureImpl[T] {
	return &FutureImpl[T]{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

// Wait blocks until the FutureImpl completes.
// It then returns the value or error stored in the FutureImpl.
//
// Possible errors returned:
// - TimeoutError: Indicates that the operation did not complete within the specified timeout duration.
// - CancelledOperationError: Indicates that the operation was canceled before it could complete.
// - Other generic errors: Represents any other error that might occur during the operation.
func (f *FutureImpl[T]) Wait() (*T, error) {
	f.cond.L.Lock()
	defer f.cond.L.Unlock()

	for f.getV() == nil && f.getE() == nil {
		f.cond.Wait()
	}

	return f.getV(), f.getE()
}

// TimeWait blocks the calling goroutine for the specified timeout duration or until the FutureImpl completes,
// whichever comes first. If the FutureImpl completes within the timeout, it returns the value or any error
// that occurred during the operation. If the timeout elapses before the FutureImpl completes, it returns a TimeoutError.
//
// Possible errors returned:
// - TimeoutError: Indicates that the operation did not complete within the specified timeout duration.
// - CancelledOperationError: Indicates that the operation was canceled before it could complete.
// - Other generic errors: Represents any other error that might occur during the operation.
//
// It's important for callers to handle these specific error cases, especially if there's a need to distinguish between
// a genuine operation failure and a timeout or cancellation.
func (f *FutureImpl[T]) TimeWait(timeout time.Duration) (*T, error) {
	f.cond.L.Lock()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	ch := make(chan struct{})
	go func() {
		f.cond.Wait()
		f.cond.L.Unlock()
		close(ch)
	}()
	for f.getV() == nil && f.getE() == nil {
		select {
		case <-timer.C:
			return nil, NewTimeoutError(timeout)
		case <-ch:
			return f.getV(), f.getE()
		}
	}

	return f.getV(), f.getE()
}

// Complete sets the value for the FutureImpl and marks it as completed.
// It then wakes up all goroutines waiting on the FutureImpl.
// It returns false if the FutureImpl is already completed.
func (f *FutureImpl[T]) Complete(t *T) bool {
	f.cond.L.Lock()
	defer f.cond.L.Unlock()

	if f.completed {
		return false
	}

	f.v.Store(t)
	f.cond.Broadcast()
	f.completed = true

	return true
}

// Fail sets an error for the FutureImpl and marks it as completed.
// It then wakes up all goroutines waiting on the FutureImpl.
// It returns false if the FutureImpl is already completed.
func (f *FutureImpl[T]) Fail(err error) bool {
	f.cond.L.Lock()
	defer f.cond.L.Unlock()

	if f.completed {
		return false
	}

	f.err.Store(err)
	f.cond.Broadcast()
	f.completed = true

	return true
}

// Cancel attempts to cancel the FutureImpl.
// It sets an internal cancellation error and marks the FutureImpl as completed.
// It then wakes up all goroutines waiting on the FutureImpl.
// It returns false if the FutureImpl is already completed.
func (f *FutureImpl[T]) Cancel() bool {
	return f.Fail(NewCancelledOperationError())
}

func (f *FutureImpl[T]) getV() *T {
	v := f.v.Load()
	if v == nil {
		return nil
	} else {
		return v.(*T)
	}
}

func (f *FutureImpl[T]) getE() error {
	err := f.err.Load()
	if err == nil {
		return nil
	} else {
		return err.(error)
	}
}
