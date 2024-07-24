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

func TestOrderedHashSet_Add(t *testing.T) {
	t.Parallel()

	set := uset.NewOrderedHashSet[testElement, int]()

	// Test adding unique elements
	assert.True(t, set.Add(testElement{key: 1}))
	assert.True(t, set.Add(testElement{key: 2}))
	assert.Equal(t, 2, set.Size())

	// Test adding duplicate elements
	assert.False(t, set.Add(testElement{key: 1}))
	assert.Equal(t, 2, set.Size())
}

func TestOrderedHashSet_Contains(t *testing.T) {
	t.Parallel()

	set := uset.NewOrderedHashSet[testElement, int]()
	set.Add(testElement{key: 1})
	set.Add(testElement{key: 2})

	// Test contains existing elements
	assert.True(t, set.Contains(testElement{key: 1}))
	assert.True(t, set.Contains(testElement{key: 2}))

	// Test contains non-existent element
	assert.False(t, set.Contains(testElement{key: 3}))
}

func TestOrderedHashSet_Remove(t *testing.T) {
	t.Parallel()

	set := uset.NewOrderedHashSet[testElement, int]()
	set.Add(testElement{key: 1})
	set.Add(testElement{key: 2})

	// Test removing existing elements
	assert.True(t, set.Remove(testElement{key: 1}))
	assert.Equal(t, 1, set.Size())

	// Test removing non-existent elements
	assert.False(t, set.Remove(testElement{key: 1}))
	assert.Equal(t, 1, set.Size())
}

func TestOrderedHashSet_Clear(t *testing.T) {
	t.Parallel()

	set := uset.NewOrderedHashSet[testElement, int]()
	set.Add(testElement{key: 1})
	set.Add(testElement{key: 2})

	// Test clear
	set.Clear()
	assert.Equal(t, 0, set.Size())
	assert.False(t, set.Contains(testElement{key: 1}))
	assert.False(t, set.Contains(testElement{key: 2}))
}

func TestOrderedHashSet_AsSlice(t *testing.T) {
	t.Parallel()

	t.Run("Initial set of elements", func(t *testing.T) {
		set := uset.NewOrderedHashSet[testElement, int]()
		set.Add(testElement{key: 1})
		set.Add(testElement{key: 2})
		set.Add(testElement{key: 3})

		expected := []testElement{
			{key: 1},
			{key: 2},
			{key: 3},
		}
		assert.Equal(t, expected, set.AsSlice())
	})

	t.Run("Duplicate elements", func(t *testing.T) {
		set := uset.NewOrderedHashSet[testElement, int]()
		set.Add(testElement{key: 1})
		set.Add(testElement{key: 1})
		set.Add(testElement{key: 2})

		expected := []testElement{
			{key: 1},
			{key: 2},
		}
		assert.Equal(t, expected, set.AsSlice())
	})

	t.Run("Removing an element", func(t *testing.T) {
		set := uset.NewOrderedHashSet[testElement, int]()
		set.Add(testElement{key: 1})
		set.Add(testElement{key: 2})
		set.Add(testElement{key: 3})
		set.Remove(testElement{key: 2})

		expected := []testElement{
			{key: 1},
			{key: 3},
		}
		assert.Equal(t, expected, set.AsSlice())
	})

	t.Run("Multiple operations", func(t *testing.T) {
		set := uset.NewOrderedHashSet[testElement, int]()
		set.Add(testElement{key: 1})
		set.Add(testElement{key: 2})
		set.Remove(testElement{key: 1})
		set.Add(testElement{key: 3})
		set.Add(testElement{key: 4})
		set.Remove(testElement{key: 3})
		set.Add(testElement{key: 1})

		expected := []testElement{
			{key: 2},
			{key: 4},
			{key: 1},
		}
		assert.Equal(t, expected, set.AsSlice())
	})

	t.Run("No elements", func(t *testing.T) {
		set := uset.NewOrderedHashSet[testElement, int]()
		assert.Empty(t, set.AsSlice())
	})

	t.Run("Clearing the set", func(t *testing.T) {
		set := uset.NewOrderedHashSet[testElement, int]()
		set.Add(testElement{key: 1})
		set.Add(testElement{key: 2})
		set.Clear()

		assert.Empty(t, set.AsSlice())
	})

	t.Run("Elements added after clearing the set", func(t *testing.T) {
		set := uset.NewOrderedHashSet[testElement, int]()
		set.Add(testElement{key: 1})
		set.Add(testElement{key: 2})
		set.Clear()

		set.Add(testElement{key: 3})
		set.Add(testElement{key: 4})

		expected := []testElement{
			{key: 3},
			{key: 4},
		}
		assert.Equal(t, expected, set.AsSlice())
	})

	t.Run("Interleaved operations", func(t *testing.T) {
		set := uset.NewOrderedHashSet[testElement, int]()
		set.Add(testElement{key: 5})
		set.Add(testElement{key: 6})
		set.Remove(testElement{key: 5})
		set.Add(testElement{key: 7})
		set.Remove(testElement{key: 6})
		set.Add(testElement{key: 8})

		expected := []testElement{
			{key: 7},
			{key: 8},
		}
		assert.Equal(t, expected, set.AsSlice())
	})
}
