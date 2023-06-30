/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package async_utils

import (
	"errors"
	"math/rand"
	"testing"
	"time"
)

func TestAsyncTaskPotentialRace(t *testing.T) {
	expectedInt := rand.Int()
	taskFunc := func() (*int, error) {
		time.Sleep(100 * time.Millisecond)
		res := expectedInt
		return &res, nil
	}

	task := NewAsyncTask[int](taskFunc, nil, 3)
	task.ExecuteAsync()

	// This sleep mimics a possible race window.
	time.Sleep(50 * time.Millisecond)

	result1, err1 := task.Wait(200 * time.Millisecond)
	result2, err2 := task.Wait(200 * time.Millisecond)

	if err1 != nil {
		t.Errorf("Call 1: Expected no error, got: %v", err1)
	}
	if err2 != nil {
		t.Errorf("Call 2: Expected no error, got: %v", err2)
	}
	if *result1 != expectedInt {
		t.Errorf("Call 1: Expected %d, got: %v", expectedInt, *result1)
	}
	if *result2 != expectedInt {
		t.Errorf("Call 2: Expected %d, got: %v", expectedInt, *result2)
	}
}

func TestAsyncTask_Success(t *testing.T) {
	task := NewAsyncTask(func() (*int, error) {
		time.Sleep(50 * time.Millisecond)
		val := 5
		return &val, nil
	}, nil, 3)

	task.ExecuteAsync()

	result, err := task.Wait(1 * time.Second)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if *result != 5 {
		t.Fatalf("Expected result to be 5, but got: %v", *result)
	}
}

func TestAsyncTask_RetryOnFailure(t *testing.T) {
	attempts := 0
	task := NewAsyncTask(func() (*int, error) {
		attempts++
		if attempts < 3 {
			return nil, errors.New("failure")
		}
		val := 5
		return &val, nil
	}, nil, 3)

	task.ExecuteAsync()

	result, err := task.Wait(1 * time.Second)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if *result != 5 {
		t.Fatalf("Expected result to be 5, but got: %v", *result)
	}
}

func TestAsyncTask_FailAfterRetries(t *testing.T) {
	task := NewAsyncTask(func() (*int, error) {
		return nil, errors.New("failure")
	}, nil, 3)

	task.ExecuteAsync()

	_, err := task.Wait(1 * time.Second)
	if err == nil {
		t.Fatal("Expected an error, but got none")
	}

	if err.Error() != "attempt 3 has failed: failure" {
		t.Fatalf("Unexpected error message: %v", err.Error())
	}
}

func TestAsyncTask_TimeoutBeforeCompletion(t *testing.T) {
	task := NewAsyncTask(func() (*int, error) {
		time.Sleep(15 * time.Second)
		val := 5
		return &val, nil
	}, nil, 3)

	task.ExecuteAsync()

	_, err := task.Wait(100 * time.Millisecond)
	if !IsTimeoutError(err) {
		t.Fatalf("Expected a timeout error, but got: %v", err)
	}
}

func TestAsyncTask_TimeoutAfterCompletion(t *testing.T) {
	task := NewAsyncTask(func() (*int, error) {
		time.Sleep(50 * time.Millisecond)
		val := 5
		return &val, nil
	}, nil, 3)

	task.ExecuteAsync()

	result, err := task.Wait(200 * time.Millisecond)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if *result != 5 {
		t.Fatalf("Expected result to be 5, but got: %v", *result)
	}
}

func TestAsyncTask_ConcurrentWaits(t *testing.T) {
	task := NewAsyncTask(func() (*int, error) {
		time.Sleep(50 * time.Millisecond)
		val := 5
		return &val, nil
	}, nil, 3)

	task.ExecuteAsync()

	done := make(chan bool)

	for i := 0; i < 5; i++ {
		go func() {
			result, err := task.Wait(200 * time.Millisecond)
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
