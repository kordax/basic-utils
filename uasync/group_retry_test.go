package uasync_test

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kordax/basic-utils/v4/uasync"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrySucceedsAfterFailure(t *testing.T) {
	var attempts atomic.Int32

	result, err := uasync.Retry(context.Background(), 3, 0, func(context.Context) (*int, error) {
		if attempts.Add(1) < 2 {
			return nil, errors.New("temporary")
		}
		v := 42
		return &v, nil
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 42, *result)
	assert.EqualValues(t, 2, attempts.Load())
}

func TestRetryStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := uasync.Retry[int](ctx, 0, time.Millisecond, func(context.Context) (*int, error) {
		t.Fatal("fn should not be called")
		return nil, nil
	})

	require.ErrorIs(t, err, context.Canceled)
}

func TestRunGroup(t *testing.T) {
	var running atomic.Int32
	var maxRunning atomic.Int32

	tasks := make([]func(context.Context) error, 10)
	for i := range tasks {
		tasks[i] = func(ctx context.Context) error {
			current := running.Add(1)
			for {
				max := maxRunning.Load()
				if current <= max || maxRunning.CompareAndSwap(max, current) {
					break
				}
			}
			time.Sleep(time.Millisecond)
			running.Add(-1)
			return ctx.Err()
		}
	}

	require.NoError(t, uasync.RunGroup(context.Background(), 3, tasks...))
	assert.LessOrEqual(t, maxRunning.Load(), int32(3))
}

func TestRunGroupCancelsOnError(t *testing.T) {
	expected := errors.New("failed")

	err := uasync.RunGroup(context.Background(), 2,
		func(context.Context) error { return expected },
		func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		},
	)

	require.ErrorIs(t, err, expected)
}

func TestRetryStopsAfterAttemptsExhausted(t *testing.T) {
	var attempts atomic.Int32

	result, err := uasync.Retry[int](context.Background(), 3, 0, func(context.Context) (*int, error) {
		attempts.Add(1)
		return nil, errors.New("failure")
	})

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualValues(t, 3, attempts.Load())
	require.ErrorContains(t, err, "retry attempts exhausted")
}

func TestRetryUsesNilContextDefaults(t *testing.T) {
	called := atomic.Bool{}

	result, err := uasync.Retry[int](nil, 1, 0, func(ctx context.Context) (*int, error) {
		require.NotNil(t, ctx)
		called.Store(true)
		value := 42
		return &value, nil
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, called.Load())
	assert.Equal(t, 42, *result)
}

func TestRetryCancelledDuringDelay(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var attempts atomic.Int32

	result, err := uasync.Retry[int](ctx, 5, time.Millisecond*20, func(context.Context) (*int, error) {
		if attempts.Add(1) == 1 {
			cancel()
		}

		return nil, fmt.Errorf("failure")
	})

	require.ErrorIs(t, err, context.Canceled)
	require.Nil(t, result)
	require.EqualValues(t, 1, attempts.Load())
}
