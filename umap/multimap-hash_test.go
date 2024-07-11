/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap_test

import (
	"testing"

	"github.com/kordax/basic-utils/umap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashMultiMap_Get(t *testing.T) {
	mm := umap.NewHashMultiMap[string, int]()
	_, exists := mm.Get("key")
	assert.False(t, exists, "Should return false for non-existent keys")

	mm.Append("key", 1)
	values, exists := mm.Get("key")
	require.True(t, exists, "Expected values to exist")
	require.Len(t, values, 1, "Expected one value")
	assert.Equal(t, 1, values[0], "Expected value to match")
}

func TestHashMultiMap_Set(t *testing.T) {
	mm := umap.NewHashMultiMap[string, int]()
	addedCount := mm.Set("key", 1, 2)
	assert.Equal(t, 2, addedCount, "Expected two values to be added")

	addedCount = mm.Set("key", 3)
	assert.Equal(t, 1, addedCount, "Expected one new value to be added")
	values, _ := mm.Get("key")
	require.Len(t, values, 1, "Expected one value after set")
	assert.Equal(t, 3, values[0], "Expected value to be '3'")
}

func TestHashMultiMap_Append(t *testing.T) {
	mm := umap.NewHashMultiMap[string, int]()
	mm.Append("key", 1)
	mm.Append("key", 2)

	addedCount := mm.Append("key", 1, 3)
	assert.Equal(t, 2, addedCount, "Expected two new values, as '1' is allowed to be duplicated")
	values, _ := mm.Get("key")
	require.Len(t, values, 4, "Expected four values after append")
	assert.ElementsMatch(t, []int{1, 2, 1, 3}, values, "Expected values to be [1, 2, 1, 3]")
}

func TestHashMultiMap_Remove(t *testing.T) {
	mm := umap.NewHashMultiMap[string, int]()
	mm.Set("key", 1, 2, 3)

	removalCount := mm.Remove("key", func(v int) bool { return v == 2 })
	assert.Equal(t, 1, removalCount, "Expected one value to be removed")
	values, _ := mm.Get("key")
	assert.Len(t, values, 2, "Expected two values left")
	assert.ElementsMatch(t, []int{1, 3}, values, "Expected '1' and '3' to remain")
}

func TestHashMultiMap_Clear(t *testing.T) {
	mm := umap.NewHashMultiMap[string, int]()
	mm.Set("key", 1)

	cleared := mm.Clear("key")
	assert.True(t, cleared, "Expected true, indicating values were cleared")

	_, exists := mm.Get("key")
	assert.False(t, exists, "Expected no values after clear")
}
