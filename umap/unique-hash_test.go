/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap_test

import (
	"crypto/sha256"
	"testing"

	"github.com/kordax/basic-utils/v2/umap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ umap.MultiMap[string, int] = (*umap.UniqueHashMultiMap[string, int])(nil)

type TestStruct struct {
	ID   int
	Name string
}

func TestUniqueHashMultiMap_Get(t *testing.T) {
	mm := umap.NewUniqueHashMultiMap[string, TestStruct](sha256.New())
	_, exists := mm.Get("key")
	assert.False(t, exists, "Should return false for non-existent keys")

	mm.Append("key", TestStruct{ID: 1, Name: "Test"})
	values, exists := mm.Get("key")
	require.True(t, exists, "Expected values to exist")
	require.Len(t, values, 1, "Expected one value")
	assert.Equal(t, TestStruct{ID: 1, Name: "Test"}, values[0], "Expected value to match")
}

func TestUniqueHashMultiMap_Set(t *testing.T) {
	mm := umap.NewUniqueHashMultiMap[string, TestStruct](sha256.New())
	addedCount := mm.Set("key", TestStruct{ID: 1, Name: "Hello"}, TestStruct{ID: 2, Name: "World"})
	assert.Equal(t, 2, addedCount, "Expected two unique values to be added")

	addedCount = mm.Set("key", TestStruct{ID: 3, Name: "New"})
	assert.Equal(t, 1, addedCount, "Expected one new value to be added")
	values, _ := mm.Get("key")
	require.Len(t, values, 1, "Expected one value after set")
	assert.Equal(t, TestStruct{ID: 3, Name: "New"}, values[0], "Expected value to be 'New'")
}

func TestUniqueHashMultiMap_Append(t *testing.T) {
	mm := umap.NewUniqueHashMultiMap[string, TestStruct](sha256.New())
	mm.Append("key", TestStruct{ID: 1, Name: "Hello"})
	mm.Append("key", TestStruct{ID: 2, Name: "World"})

	addedCount := mm.Append("key", TestStruct{ID: 1, Name: "Hello"}, TestStruct{ID: 3, Name: "New"})
	assert.Equal(t, 1, addedCount, "Expected one new value, as 'Hello' is a duplicate")

	values, _ := mm.Get("key")
	require.Len(t, values, 3, "Expected three values after append")
	expectedValues := []TestStruct{
		{ID: 1, Name: "Hello"},
		{ID: 2, Name: "World"},
		{ID: 3, Name: "New"},
	}
	assert.ElementsMatch(t, expectedValues, values, "Expected values to be [{ID: 1, Name: \"Hello\"}, {ID: 2, Name: \"World\"}, {ID: 3, Name: \"New\"}]")
}

func TestUniqueHashMultiMap_Remove(t *testing.T) {
	mm := umap.NewUniqueHashMultiMap[string, TestStruct](sha256.New())
	mm.Set("key", TestStruct{ID: 1, Name: "Remove"}, TestStruct{ID: 2, Name: "Keep"})

	removalCount := mm.Remove("key", func(v TestStruct) bool { return v.Name == "Remove" })
	assert.Equal(t, 1, removalCount, "Expected one value to be removed")
	values, _ := mm.Get("key")
	assert.Len(t, values, 1, "Expected one value left")
	assert.Equal(t, TestStruct{ID: 2, Name: "Keep"}, values[0], "Expected 'Keep' to remain")
}

func TestUniqueHashMultiMap_Clear(t *testing.T) {
	mm := umap.NewUniqueHashMultiMap[string, TestStruct](sha256.New())
	mm.Set("key", TestStruct{ID: 1, Name: "Data"})

	cleared := mm.Clear("key")
	assert.True(t, cleared, "Expected true, indicating values were cleared")

	_, exists := mm.Get("key")
	assert.False(t, exists, "Expected no values after clear")
}

func TestUniqueHashMultiMap_Iterator(t *testing.T) {
	mm := umap.NewUniqueHashMultiMap[string, int](sha256.New())
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
		collected[k] = append(collected[k], v...)
	}

	require.Len(t, collected, 3, "Expected three keys in the map")
	assert.ElementsMatch(t, []int{1, 2, 3}, collected["key1"], "Values for key1 should match")
	assert.ElementsMatch(t, []int{4, 5}, collected["key2"], "Values for key2 should match")
	assert.ElementsMatch(t, []int{6}, collected["key3"], "Values for key3 should match")
}
