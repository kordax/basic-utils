/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uerror_test

import (
	"errors"
	"testing"

	"github.com/kordax/basic-utils/uerror"
	"github.com/stretchr/testify/assert"
)

func TestMust(t *testing.T) {
	t.Run("NoError", func(t *testing.T) {
		assert.NotPanics(t, func() {
			uerror.Must(nil)
		})
	})

	t.Run("WithError", func(t *testing.T) {
		assert.PanicsWithError(t, "test error", func() {
			uerror.Must(errors.New("test error"))
		})
	})

	t.Run("WithMultiReturnNoError", func(t *testing.T) {
		assert.NotPanics(t, func() {
			uerror.Must(successfulOperation())
		})
	})

	t.Run("WithMultiReturnError", func(t *testing.T) {
		assert.PanicsWithError(t, "operation failed", func() {
			uerror.Must(failingOperation())
		})
	})
}

// successfulOperation simulates a function that returns a result and no error
func successfulOperation() (string, error) {
	return "success", nil
}

// failingOperation simulates a function that returns a result and an error
func failingOperation() (string, error) {
	return "", errors.New("operation failed")
}
