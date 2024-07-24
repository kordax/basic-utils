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

type testElement struct {
	key int
}

func (e testElement) Key() int {
	return e.key
}

func TestComparableHashSet_Add(t *testing.T) {
	t.Parallel()

	set := uset.NewComparableHashSet[testElement, int]()

	assert.True(t, set.Add(testElement{key: 1}))
	assert.True(t, set.Add(testElement{key: 2}))
	assert.Equal(t, 2, set.Size())

	assert.False(t, set.Add(testElement{key: 1}))
	assert.Equal(t, 2, set.Size())
}

func TestComparableHashSet_Contains(t *testing.T) {
	t.Parallel()

	set := uset.NewComparableHashSet[testElement, int]()
	set.Add(testElement{key: 1})
	set.Add(testElement{key: 2})

	assert.True(t, set.Contains(testElement{key: 1}))
	assert.True(t, set.Contains(testElement{key: 2}))
	assert.False(t, set.Contains(testElement{key: 3}))
}

func TestComparableHashSet_Remove(t *testing.T) {
	t.Parallel()

	set := uset.NewComparableHashSet[testElement, int]()
	set.Add(testElement{key: 1})
	set.Add(testElement{key: 2})

	assert.True(t, set.Remove(testElement{key: 1}))
	assert.Equal(t, 1, set.Size())
	assert.False(t, set.Remove(testElement{key: 1}))
	assert.Equal(t, 1, set.Size())
}

func TestComparableHashSet_Clear(t *testing.T) {
	t.Parallel()

	set := uset.NewComparableHashSet[testElement, int]()
	set.Add(testElement{key: 1})
	set.Add(testElement{key: 2})

	set.Clear()
	assert.Equal(t, 0, set.Size())
	assert.False(t, set.Contains(testElement{key: 1}))
	assert.False(t, set.Contains(testElement{key: 2}))
}
