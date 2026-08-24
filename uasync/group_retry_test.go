package uasync_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"git.casinomodule.org/casino27/basic-utils/v4/uasync"
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
