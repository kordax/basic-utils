/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package asyncutils

import (
	"errors"
	"fmt"
	"time"
)

// TimeoutError is used to indicate that a specific operation has exceeded the allocated time to complete.
// It provides information about the duration after which the timeout occurred.It's particularly useful in scenarios where you have operations that can block for an extended period, and you want to ensure they complete within a given timeframe. For example, waiting for a response from a remote service or waiting for a Future to complete.
type TimeoutError struct {
	Duration time.Duration
}

func NewTimeoutError(duration time.Duration) *TimeoutError {
	return &TimeoutError{Duration: duration}
}

func (e TimeoutError) Error() string {
	return fmt.Sprintf("timeout after %v", e.Duration)
}

func IsTimeoutError(err error) bool {
	var timeoutError *TimeoutError
	ok := errors.As(err, &timeoutError)
	return ok
}

// CancelledOperationError represents an error that occurs when an operation is canceled before it completes.
// This can be due to user intervention or a system decision.
// It's useful in scenarios where long-running tasks can be canceled by external triggers or user input,
// such as canceling an async task or an HTTP request.
type CancelledOperationError struct{}

func NewCancelledOperationError() *CancelledOperationError {
	return &CancelledOperationError{}
}

func (e CancelledOperationError) Error() string {
	return "cancelled"
}

func IsCancelledOperationError(err error) bool {
	var cancelledOperationError *CancelledOperationError
	ok := errors.As(err, &cancelledOperationError)
	return ok
}
