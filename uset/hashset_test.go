/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uset_test

import (
	"testing"

	"github.com/kordax/basic-utils/uset"
	"github.com/stretchr/testify/assert"
)

func TestHashSet_Add(t *testing.T) {
	t.Parallel()

	set := uset.NewHashSet[int]()

	// Test adding unique elements
	assert.True(t, set.Add(1))
	assert.True(t, set.Add(2))
	assert.Equal(t, 2, set.Size())

	// Test adding duplicate elements
	assert.False(t, set.Add(1))
	assert.Equal(t, 2, set.Size())
}

func TestHashSet_Contains(t *testing.T) {
	t.Parallel()

	set := uset.NewHashSet[int]()
	set.Add(1)
	set.Add(2)

	// Test contains existing elements
	assert.True(t, set.Contains(1))
	assert.True(t, set.Contains(2))

	// Test contains non-existent element
	assert.False(t, set.Contains(3))
}

func TestHashSet_Remove(t *testing.T) {
	t.Parallel()

	set := uset.NewHashSet[int]()
	set.Add(1)
	set.Add(2)

	// Test removing existing elements
	assert.True(t, set.Remove(1))
	assert.Equal(t, 1, set.Size())

	// Test removing non-existent elements
	assert.False(t, set.Remove(1))
	assert.Equal(t, 1, set.Size())
}

func TestHashSet_Clear(t *testing.T) {
	t.Parallel()

	set := uset.NewHashSet[int]()
	set.Add(1)
	set.Add(2)

	// Test clear
	set.Clear()
	assert.Equal(t, 0, set.Size())
	assert.False(t, set.Contains(1))
	assert.False(t, set.Contains(2))
}
