/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2025.
 */

package uasync_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kordax/basic-utils/v2/uasync"
	"github.com/stretchr/testify/require"
)

func TestScheduledTask_ExecutesAfterDelay(t *testing.T) {
	ctx := context.Background()
	start := time.Now().Add(100 * time.Millisecond)

	called := make(chan struct{}, 1)
	task := uasync.NewScheduledTask(ctx, start, func(ctx context.Context) (*int, error) {
		val := 42
		called <- struct{}{}
		return &val, nil
	}, 0)

	fut := task.Schedule()
	select {
	case <-called:
		// ok
	case <-time.After(300 * time.Millisecond):
		t.Fatal("task did not execute in expected time")
	}

	res, err := fut.Wait()
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, 42, *res)
}

func TestScheduledTask_CancelBeforeStart(t *testing.T) {
	ctx := context.Background()
	start := time.Now().Add(1 * time.Second)

	task := uasync.NewScheduledTask(ctx, start, func(ctx context.Context) (*int, error) {
		val := 1
		return &val, nil
	}, 0)

	fut := task.Schedule()
	task.Cancel()

	_, err := fut.Wait()
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
}

func TestScheduleTask_WithExistingTask(t *testing.T) {
	ctx := context.Background()
	at := uasync.NewAsyncTask(ctx, func(ctx context.Context) (*string, error) {
		s := "done"
		return &s, nil
	}, 0)

	start := time.Now().Add(100 * time.Millisecond)
	sched := uasync.AsScheduledTask(start, at)

	fut := sched.Schedule()

	val, err := fut.Wait()
	require.NoError(t, err)
	require.Equal(t, "done", *val)
}

func TestScheduledTask_RetriesOnError(t *testing.T) {
	ctx := context.Background()
	start := time.Now().Add(50 * time.Millisecond)

	callCount := 0
	task := uasync.NewScheduledTask(ctx, start, func(ctx context.Context) (*int, error) {
		callCount++
		if callCount < 3 {
			return nil, errors.New("fail")
		}
		v := 7
		return &v, nil
	}, 3)

	fut := task.Schedule()
	val, err := fut.Wait()
	require.NoError(t, err)
	require.Equal(t, 7, *val)
	require.Equal(t, 3, callCount)
}
