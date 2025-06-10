/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uasync

import (
	"context"
)

// Execute initiates the execution of the provided operation and waits for its completion up to the given timeout duration.
// If the operation completes within the timeout, the function returns the result of the operation or any error that occurred.
// If the timeout elapses before the operation completes, the function returns context error and the operation is canceled using the provided cancelFunc.
//
// Possible errors returned:
// - context.Canceled: Indicates that the operation was canceled before it could complete.
// - Other generic errors: Represents any other error that might occur during the operation.
//
// It's important for callers to handle these specific error cases, especially if there's a need to distinguish between a genuine operation failure and a timeout or cancellation.
func Execute[R any](ctx context.Context, fn func(ctx context.Context) (*R, error)) (*R, error) {
	task := NewAsyncTask(ctx, fn, 0)
	task.ExecuteAsync()

	return task.Wait()
}
