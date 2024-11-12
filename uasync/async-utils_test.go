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
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	t.Run("success within timeout", func(t *testing.T) {
		fn := func(ctx context.Context) (*int, error) {
			val := 5
			return &val, nil
		}
		ctx := context.Background()

		result, err := uasync.Execute(ctx, fn)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if *result != 5 {
			t.Fatalf("Expected result to be 5, got %v", *result)
		}
	})

	t.Run("timeout before completion", func(t *testing.T) {
		fn := func(ctx context.Context) (*int, error) {
			time.Sleep(2 * time.Second)
			val := 5
			return &val, nil
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
		defer cancel()

		_, err := uasync.Execute(ctx, fn)
		require.Error(t, err)
	})

	t.Run("function returns error", func(t *testing.T) {
		fn := func(ctx context.Context) (*int, error) {
			return nil, errors.New("sample error")
		}
		ctx := context.Background()

		_, err := uasync.Execute(ctx, fn)
		if err == nil {
			t.Fatal("Expected an error, got none")
		}
		if err.Error() != "attempt 0 has failed: sample error" {
			t.Fatalf("Expected error 'sample error', got: %v", err)
		}
	})
}
