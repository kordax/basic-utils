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

var _ umap.MultiMap[string, int] = (*umap.HashMultiMap[string, int])(nil)

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

func TestHashMultiMap_Iterator(t *testing.T) {
	mm := umap.NewHashMultiMap[string, int]()
	mm.Append("key1", 1, 2, 3)
	mm.Append("key2", 4, 5)
	mm.Append("key3", 6)

	collected := make(map[string][]int)
	mm.Iterator()(func(key string, values []int) bool {
		collected[key] = append(collected[key], values...)
		return true
	})

	require.Len(t, collected, 3, "Expected three keys in the map")
	assert.ElementsMatch(t, []int{1, 2, 3}, collected["key1"], "Values for key1 should match")
	assert.ElementsMatch(t, []int{4, 5}, collected["key2"], "Values for key2 should match")
	assert.ElementsMatch(t, []int{6}, collected["key3"], "Values for key3 should match")

	collected = make(map[string][]int)
	for k, v := range mm.Iterator() {
		collected[k] = v
	}
	require.Len(t, collected, 3, "Expected three keys in the map")
	assert.ElementsMatch(t, []int{1, 2, 3}, collected["key1"], "Values for key1 should match")
	assert.ElementsMatch(t, []int{4, 5}, collected["key2"], "Values for key2 should match")
	assert.ElementsMatch(t, []int{6}, collected["key3"], "Values for key3 should match")
}
