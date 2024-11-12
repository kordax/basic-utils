/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uasync

import (
	"errors"
)

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
