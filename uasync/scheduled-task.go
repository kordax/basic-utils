/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2025.
 */

package uasync

import (
	"context"
	"sync"
	"time"
)

// ScheduledTask wraps an AsyncTask and schedules it for execution at a specific time.
type ScheduledTask[R any] struct {
	task    *AsyncTask[R]
	startAt time.Time

	started bool
	mtx     sync.Mutex
}

// NewScheduledTask creates a new ScheduledTask that will start the given async function at `startAt`.
func NewScheduledTask[R any](ctx context.Context, startAt time.Time, fn func(ctx context.Context) (*R, error), retries int) *ScheduledTask[R] {
	at := &AsyncTask[R]{
		ctx:     ctx,
		fn:      fn,
		retries: retries,
		f:       NewFuture[R](ctx),
	}
	return &ScheduledTask[R]{
		task:    at,
		startAt: startAt,
	}
}

// AsScheduledTask creates a new ScheduledTask from existing task.
func AsScheduledTask[R any](startAt time.Time, task *AsyncTask[R]) *ScheduledTask[R] {
	return &ScheduledTask[R]{
		task:    task,
		startAt: startAt,
	}
}

// Schedule schedules the task to run at the predefined time.
// Returns the Future so callers can wait for the result or nil, if task was already srated.
func (s *ScheduledTask[R]) Schedule() Future[R] {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.started {
		return nil
	}

	s.started = true
	go func() {
		delay := time.Until(s.startAt)
		if delay > 0 {
			select {
			case <-time.After(delay):
			case <-s.task.ctx.Done():
				// Task was canceled before it could start, cancel will be called in a task itself.
				return
			}
		}
		s.task.ExecuteAsync()
	}()

	return s.task.f
}

// Cancel cancels the scheduled task before it starts (if still pending) or does nothing if task wasn't started.
func (s *ScheduledTask[R]) Cancel() {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if !s.started {
		return
	}

	s.task.cancel()
}
