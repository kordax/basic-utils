/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uasync

import (
	"context"
	"sync"
	"sync/atomic"
)

// Future is an interface that represents a value or an error that will be available in the future.
// It provides mechanisms to wait for the value, set the value or error, and cancel the operation.
type Future[T any] interface {
	Wait() (*T, error)
	Cancel() bool
	Complete(t *T) bool
	Fail(err error) bool
}

// FutureImpl is an implementation of the Future interface.
// It uses condition variables to allow waiting for the value or error to be set.
type FutureImpl[T any] struct {
	cond *sync.Cond

	v   atomic.Value
	err atomic.Value

	completed bool

	ctx context.Context
}

func NewFuture[T any](ctx context.Context) *FutureImpl[T] {
	return &FutureImpl[T]{
		cond: sync.NewCond(&sync.Mutex{}),
		ctx:  ctx,
	}
}

// Wait blocks until the FutureImpl completes.
// It then returns the value or error stored in the FutureImpl.
func (f *FutureImpl[T]) Wait() (*T, error) {
	done := make(chan struct{})

	go func() {
		select {
		case <-f.ctx.Done():
			f.cond.L.Lock()
			defer f.cond.L.Unlock()
			f.cond.Broadcast()
		case <-done:
		}
	}()

	f.cond.L.Lock()
	defer f.cond.L.Unlock()
	defer close(done)

	for f.getV() == nil && f.getE() == nil {
		if f.ctx.Err() != nil {
			return nil, f.ctx.Err()
		}
		f.cond.Wait()
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
	return f.Fail(context.Canceled)
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
