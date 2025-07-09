/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2025.
 */

package uonce

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testObj struct {
	n int
}

func TestOnce_Singleton(t *testing.T) {
	counter := 0
	getObj := Once(func() *testObj {
		counter++
		return &testObj{n: 42}
	})

	o1 := getObj()
	o2 := getObj()

	require.Same(t, o1, o2, "expected same instance on every call")
	assert.Equal(t, 42, o1.n)
	assert.Equal(t, 1, counter, "constructor should only be called once")
}

func TestOnce_ThreadSafe(t *testing.T) {
	counter := 0
	getObj := Once(func() *testObj {
		counter++
		return &testObj{n: 99}
	})

	const goroutines = 50
	var wg sync.WaitGroup
	results := make(chan *testObj, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- getObj()
		}()
	}

	wg.Wait()
	close(results)

	var singleton *testObj
	for o := range results {
		if singleton == nil {
			singleton = o
		}
		require.Same(t, singleton, o, "all goroutines should get same instance")
	}
	assert.Equal(t, 1, counter, "constructor should only be called once, even with concurrency")
}
