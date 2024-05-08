/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uasync_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kordax/basic-utils/uasync"
)

func TestExecute(t *testing.T) {
	t.Run("success within timeout", func(t *testing.T) {
		fn := func() (*int, error) {
			val := 5
			return &val, nil
		}
		cancelFunc := context.CancelFunc(func() {})

		result, err := uasync.Execute(fn, cancelFunc, 1*time.Second)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if *result != 5 {
			t.Fatalf("Expected result to be 5, got %v", *result)
		}
	})

	t.Run("timeout before completion", func(t *testing.T) {
		fn := func() (*int, error) {
			time.Sleep(2 * time.Second)
			val := 5
			return &val, nil
		}
		cancelFunc := context.CancelFunc(func() {})

		_, err := uasync.Execute(fn, cancelFunc, 1*time.Second)
		if err == nil {
			t.Fatal("Expected timeout error, got none")
		}
		if !uasync.IsTimeoutError(err) {
			t.Fatalf("Expected TimeoutError, got: %v", err)
		}
	})

	t.Run("function returns error", func(t *testing.T) {
		fn := func() (*int, error) {
			return nil, errors.New("sample error")
		}
		cancelFunc := context.CancelFunc(func() {})

		_, err := uasync.Execute(fn, cancelFunc, 1*time.Second)
		if err == nil {
			t.Fatal("Expected an error, got none")
		}
		if err.Error() != "attempt 0 has failed: sample error" {
			t.Fatalf("Expected error 'sample error', got: %v", err)
		}
	})
}
