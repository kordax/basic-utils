/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uasync_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/kordax/basic-utils/v2/uasync"
	"github.com/stretchr/testify/require"
)

func TestAsyncTaskPotentialRace(t *testing.T) {
	expectedInt := rand.Int()
	taskFunc := func(ctx context.Context) (*int, error) {
		time.Sleep(time.Millisecond)
		res := expectedInt
		return &res, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	task := uasync.NewAsyncTask[int](ctx, taskFunc, 3)
	task.ExecuteAsync()

	// This sleep mimics a possible race window.

	result1, err1 := task.Wait()
	//result2, err2 := task.Wait()

	if err1 != nil {
		t.Errorf("Call 1: Expected no error, got: %v", err1)
	}
	//if err2 != nil {
	//	t.Errorf("Call 2: Expected no error, got: %v", err2)
	//}
	if *result1 != expectedInt {
		t.Errorf("Call 1: Expected %d, got: %v", expectedInt, *result1)
	}
	//if *result2 != expectedInt {
	//	t.Errorf("Call 2: Expected %d, got: %v", expectedInt, *result2)
	//}
}

func TestAsyncTask_Success(t *testing.T) {
	task := uasync.NewAsyncTask(context.Background(), func(ctx context.Context) (*int, error) {
		val := 5
		return &val, nil
	}, 3)

	task.ExecuteAsync()

	result, err := task.Wait()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if *result != 5 {
		t.Fatalf("Expected result to be 5, but got: %v", *result)
	}
}

func TestAsyncTask_RetryOnFailure(t *testing.T) {
	attempts := 0
	task := uasync.NewAsyncTask(context.Background(), func(ctx context.Context) (*int, error) {
		attempts++
		if attempts < 3 {
			return nil, errors.New("failure")
		}
		val := 5
		return &val, nil
	}, 3)

	task.ExecuteAsync()

	result, err := task.Wait()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if *result != 5 {
		t.Fatalf("Expected result to be 5, but got: %v", *result)
	}
}

func TestAsyncTask_FailAfterRetries(t *testing.T) {
	task := uasync.NewAsyncTask(context.Background(), func(ctx context.Context) (*int, error) {
		return nil, errors.New("failure")
	}, 3)

	task.ExecuteAsync()

	_, err := task.Wait()
	if err == nil {
		t.Fatal("Expected an error, but got none")
	}

	if err.Error() != "attempt 3 has failed: failure" {
		t.Fatalf("Unexpected error message: %v", err.Error())
	}
}

func TestAsyncTask_TimeoutBeforeCompletion(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	task := uasync.NewAsyncTask(ctx, func(ctx context.Context) (*int, error) {
		time.Sleep(15 * time.Second)
		val := 5
		return &val, nil
	}, 3)

	task.ExecuteAsync()

	_, err := task.Wait()
	require.Error(t, err)
}

func TestAsyncTask_TimeoutAfterCompletion(t *testing.T) {
	task := uasync.NewAsyncTask(context.Background(), func(ctx context.Context) (*int, error) {
		val := 5
		return &val, nil
	}, 3)

	task.ExecuteAsync()

	result, err := task.Wait()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if *result != 5 {
		t.Fatalf("Expected result to be 5, but got: %v", *result)
	}
}

func TestAsyncTask_ConcurrentWaits(t *testing.T) {
	task := uasync.NewAsyncTask(context.Background(), func(ctx context.Context) (*int, error) {
		val := 5
		return &val, nil
	}, 3)

	task.ExecuteAsync()

	done := make(chan bool)

	for i := 0; i < 5; i++ {
		go func() {
			result, err := task.Wait()
			if err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}

			if *result != 5 {
				t.Errorf("Expected result to be 5, but got: %v", *result)
			}
			done <- true
		}()
	}

	for i := 0; i < 5; i++ {
		<-done
	}
}

func TestAsyncTask_InfiniteRetryUntilSuccess(t *testing.T) {
	attempts := 0
	task := uasync.NewAsyncTask(context.Background(), func(ctx context.Context) (*int, error) {
		attempts++
		if attempts < 5 {
			return nil, errors.New("temporary failure")
		}
		val := 42
		return &val, nil
	}, -1)

	task.ExecuteAsync()

	result, err := task.Wait()
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 42, *result)
	require.GreaterOrEqual(t, attempts, 5)
}
