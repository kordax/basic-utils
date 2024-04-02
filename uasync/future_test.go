/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uasync

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFutureSimpleCompletion(t *testing.T) {
	f := NewFuture[int]()
	value := 42
	f.Complete(&value)

	result, err := f.Wait()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if *result != value {
		t.Fatalf("Expected %d, got: %d", value, *result)
	}
}

func TestFutureSimpleFailure(t *testing.T) {
	f := NewFuture[int]()
	expectedErr := errors.New("an error")
	f.Fail(expectedErr)

	_, err := f.Wait()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, got: %v", expectedErr, err)
	}
}

func TestFutureTimeoutBeforeCompletion(t *testing.T) {
	f := NewFuture[int]()

	go func() {
		time.Sleep(3 * time.Second)
		val := 100
		f.Complete(&val)
	}()

	_, err := f.TimeWait(100 * time.Millisecond)
	if err == nil {
		t.Fatal("Expected timeout error, got none")
	}
}

func TestFutureTimeoutAfterCompletion(t *testing.T) {
	f := NewFuture[int]()
	val := 123
	f.Complete(&val)

	result, err := f.TimeWait(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if *result != val {
		t.Fatalf("Expected %d, got: %d", val, *result)
	}
}

func TestFutureConcurrentCompletionAndFailure(t *testing.T) {
	f := NewFuture[int]()
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		val := 10
		f.Complete(&val)
	}()

	go func() {
		defer wg.Done()
		f.Fail(errors.New("an error"))
	}()

	wg.Wait()

	// Only one of the above goroutines should succeed, so either result will be non-nil, or error will be non-nil
	result, err := f.Wait()
	if result == nil && err == nil {
		t.Fatal("Expected either result or error, got neither")
	}
}

func TestFutureMultipleCompletions(t *testing.T) {
	f := NewFuture[int]()
	val1 := 1
	val2 := 2
	f.Complete(&val1)
	if f.Complete(&val2) {
		t.Fatal("Expected subsequent completion to be unsuccessful")
	}

	result, err := f.Wait()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if *result != val1 {
		t.Fatalf("Expected %d, got: %d", val1, *result)
	}
}

func TestFutureConcurrency(t *testing.T) {
	f := NewFuture[int]()
	expected := 0

	// Spawn 1000 goroutines trying to complete/fail the future at the same time
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				if f.Complete(&i) {
					expected = i
				}
			} else {
				if f.Fail(fmt.Errorf("error at iteration #%d", i)) {
					expected = -1
				}
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Check if the result is consistent
	val, err := f.Wait()
	if val != nil && err != nil {
		t.Fatalf("Both value and error are set. Value: %v, Error: %v", val, err)
	}

	if expected == -1 {
		if !assert.Error(t, err) {
			t.Fatal("Fail call with error was expected on")
		}
	} else {
		if !assert.NotNil(t, val) {
			t.Fatal("Value expected, but nil received")
		}
		assert.Equal(t, expected, *val)
	}

	if val == nil && err == nil {
		t.Fatal("Neither value nor error is set")
	}
}

func TestDeadlockRisk(t *testing.T) {
	f := NewFuture[int]()

	// Channel to signal when the lock is obtained in the second goroutine
	lockObtained := make(chan struct{})

	// Goroutine that calls TimeWait with a long timeout
	go func() {
		_, _ = f.TimeWait(15 * time.Second)
	}()

	// Small sleep to ensure the above goroutine starts first
	time.Sleep(100 * time.Millisecond)

	// Goroutine that tries to lock f.cond.L
	go func() {
		f.cond.L.Lock()
		close(lockObtained)
		f.cond.L.Unlock()
	}()

	// We wait up to 4 seconds for the lock to be obtained, which is less than the TimeWait timeout
	select {
	case <-lockObtained:
		// Lock was obtained before the timeout, no deadlock
	case <-time.After(1 * time.Second):
		t.Fatal("Potential deadlock detected!")
	}
}

func TestFutureCancelBeforeCompletion(t *testing.T) {
	f := NewFuture[int]()

	// Launch a goroutine that tries to complete the future after a delay.
	go func() {
		time.Sleep(3 * time.Second)
		val := 100
		f.Complete(&val)
	}()

	// Cancel the future before it can complete.
	cancelled := f.Cancel()
	if !cancelled {
		t.Fatal("Expected future to be cancelled, but it wasn't")
	}

	// Check if the future returns the CancelledOperationError after being cancelled.
	_, err := f.Wait()
	if err == nil {
		t.Fatal("Expected cancellation error, got none")
	}

	if !IsCancelledOperationError(err) {
		t.Fatalf("Expected CancelledOperationError, got: %v", err)
	}
}
