/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2025.
 */

package uonce

import "sync"

// Once returns a function that always returns the same instance of T,
// initializing it only once using the provided constructor.
//
// This utility is intended for implementing singletons in Go—
// it ensures the constructor is called exactly once, and all callers
// receive the same instance, even in the presence of concurrent calls.
//
// Example (singleton pattern):
//
//	// Define a package-level singleton getter
//	getFoo := uonce.Once(func() *Foo {
//		return &Foo{X: 123}
//	})
//
//	f1 := getFoo()
//	f2 := getFoo()
//	// f1 == f2 will be true
//
// The returned function is safe for concurrent use from multiple goroutines.
func Once[T any](constructor func() T) func() T {
	var (
		once sync.Once
		val  T
	)
	return func() T {
		once.Do(func() { val = constructor() })
		return val
	}
}
