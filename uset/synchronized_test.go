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

const (
	numGoroutines = 1000
	numElements   = 100
)

func TestSynchronizedHashSet_Add(t *testing.T) {
	t.Parallel()

	s := uset.NewSynchronizedHashSet[int]()
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numElements; j++ {
				s.Add(i*numElements + j)
			}
		}(i)
	}

	wg.Wait()
	assert.Equal(t, numGoroutines*numElements, s.Size())
}

func TestSynchronizedHashSet_Contains(t *testing.T) {
	t.Parallel()

	s := uset.NewSynchronizedHashSet[int]()
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < numElements; j++ {
			s.Add(i*numElements + j)
		}
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numElements; j++ {
				assert.True(t, s.Contains(i*numElements+j))
			}
		}(i)
	}

	wg.Wait()
}

func TestSynchronizedHashSet_Remove(t *testing.T) {
	t.Parallel()

	s := uset.NewSynchronizedHashSet[int]()
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < numElements; j++ {
			s.Add(i*numElements + j)
		}
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numElements; j++ {
				s.Remove(i*numElements + j)
			}
		}(i)
	}

	wg.Wait()
	assert.Equal(t, 0, s.Size())
}

func TestSynchronizedHashSet_Clear(t *testing.T) {
	t.Parallel()

	s := uset.NewSynchronizedHashSet[int]()
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < numElements; j++ {
			s.Add(i*numElements + j)
		}
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Clear()
		}()
	}

	wg.Wait()
	assert.Equal(t, 0, s.Size())
}
