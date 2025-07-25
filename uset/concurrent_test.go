/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uset_test

import (
	"sync"
	"testing"

	"github.com/kordax/basic-utils/v2/uset"
	"github.com/stretchr/testify/assert"
)

func TestConcurrentHashSet_Add(t *testing.T) {
	set := uset.NewConcurrentHashSet[int]()

	assert.True(t, set.Add(1))
	assert.True(t, set.Add(2))
	assert.False(t, set.Add(1)) // Adding duplicate

	assert.Equal(t, 2, set.Size())
}

func TestConcurrentHashSet_Contains(t *testing.T) {
	set := uset.NewConcurrentHashSet[int]()

	set.Add(1)
	set.Add(2)

	assert.True(t, set.Contains(1))
	assert.True(t, set.Contains(2))
	assert.False(t, set.Contains(3))
}

func TestConcurrentHashSet_Remove(t *testing.T) {
	set := uset.NewConcurrentHashSet[int]()

	set.Add(1)
	set.Add(2)

	assert.True(t, set.Remove(1))
	assert.False(t, set.Remove(1)) // Removing non-existent item
	assert.Equal(t, 1, set.Size())
	assert.False(t, set.Contains(1))
	assert.True(t, set.Contains(2))
}

func TestConcurrentHashSet_Size(t *testing.T) {
	set := uset.NewConcurrentHashSet[int]()

	set.Add(1)
	set.Add(2)
	set.Add(3)

	assert.Equal(t, 3, set.Size())
}

func TestConcurrentHashSet_Clear(t *testing.T) {
	set := uset.NewConcurrentHashSet[int]()

	set.Add(1)
	set.Add(2)
	set.Add(3)

	set.Clear()

	assert.Equal(t, 0, set.Size())
	assert.False(t, set.Contains(1))
	assert.False(t, set.Contains(2))
	assert.False(t, set.Contains(3))
}

func TestConcurrentHashSet_CustomHasher(t *testing.T) {
	hasher := func(value int) uint64 {
		return uint64(value) * 123456789
	}
	set := uset.NewCustomConcurrentHashSet[int](hasher)

	assert.True(t, set.Add(1))
	assert.True(t, set.Add(2))
	assert.False(t, set.Add(1)) // Adding duplicate

	assert.Equal(t, 2, set.Size())
	assert.True(t, set.Contains(1))
	assert.True(t, set.Contains(2))
	assert.False(t, set.Contains(3))

	assert.True(t, set.Remove(1))
	assert.False(t, set.Contains(1))
	assert.True(t, set.Contains(2))
	assert.Equal(t, 1, set.Size())
}

func TestConcurrentHashSet_ConcurrentAddAndContains(t *testing.T) {
	set := uset.NewConcurrentHashSet[int]()

	const numGoroutines = 100
	const numElements = 1000

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numElements; j++ {
				set.Add(i*numElements + j)
			}
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numElements; j++ {
				set.Contains(i*numElements + j)
			}
		}(i)
	}

	wg.Wait()

	expectedSize := numGoroutines * numElements
	assert.Equal(t, expectedSize, set.Size(), "Expected size does not match actual size")
}

func TestConcurrentHashSet_ConcurrentAddAndRemove(t *testing.T) {
	set := uset.NewConcurrentHashSet[int]()

	const numGoroutines = 100
	const numElements = 1000

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numElements; j++ {
				set.Add(i*numElements + j)
			}
		}(i)
	}

	wg.Wait()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numElements; j++ {
				set.Remove(i*numElements + j)
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, 0, set.Size(), "Expected size after removal does not match actual size")
}

func TestConcurrentHashSet_ConcurrentAddAndClear(t *testing.T) {
	set := uset.NewConcurrentHashSet[int]()

	const numGoroutines = 100
	const numElements = 1000

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numElements; j++ {
				set.Add(i*numElements + j)
			}
		}(i)
	}

	wg.Wait()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			set.Clear()
		}()
	}

	wg.Wait()

	assert.Equal(t, 0, set.Size(), "Expected size after clear does not match actual size")
}

func TestConcurrentHashSet_ConcurrentAddAndSize(t *testing.T) {
	set := uset.NewConcurrentHashSet[int]()

	const numGoroutines = 100
	const numElements = 1000

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numElements; j++ {
				set.Add(i*numElements + j)
			}
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			set.Size()
		}()
	}

	wg.Wait()

	expectedSize := numGoroutines * numElements
	assert.Equal(t, expectedSize, set.Size(), "Expected size does not match actual size")
}
