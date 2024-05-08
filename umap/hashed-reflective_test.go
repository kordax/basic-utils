/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap_test

import (
	"crypto/sha256"
	"testing"

	"github.com/kordax/basic-utils/umap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	ID   int
	Name string
}

func TestHashedReflectiveMultiMap_Get(t *testing.T) {
	mm := umap.NewHashedReflectiveMultiMap[string, TestStruct](sha256.New())
	_, exists := mm.Get("key")
	assert.False(t, exists, "Should return false for non-existent keys")

	mm.Append("key", TestStruct{ID: 1, Name: "Test"})
	values, exists := mm.Get("key")
	require.True(t, exists, "Expected values to exist")
	require.Len(t, values, 1, "Expected one value")
	assert.Equal(t, TestStruct{ID: 1, Name: "Test"}, values[0], "Expected value to match")
}

func TestHashedReflectiveMultiMap_Set(t *testing.T) {
	mm := umap.NewHashedReflectiveMultiMap[string, TestStruct](sha256.New())
	addedCount := mm.Set("key", TestStruct{ID: 1, Name: "Hello"}, TestStruct{ID: 2, Name: "World"})
	assert.Equal(t, 2, addedCount, "Expected two unique values to be added")

	addedCount = mm.Set("key", TestStruct{ID: 3, Name: "New"})
	assert.Equal(t, 1, addedCount, "Expected one new value to be added")
	values, _ := mm.Get("key")
	require.Len(t, values, 1, "Expected one value after set")
	assert.Equal(t, TestStruct{ID: 3, Name: "New"}, values[0], "Expected value to be 'New'")
}

func TestHashedReflectiveMultiMap_Append(t *testing.T) {
	mm := umap.NewHashedReflectiveMultiMap[string, TestStruct](sha256.New())
	mm.Append("key", TestStruct{ID: 1, Name: "Hello"})
	mm.Append("key", TestStruct{ID: 2, Name: "World"})

	addedCount := mm.Append("key", TestStruct{ID: 1, Name: "Hello"}, TestStruct{ID: 3, Name: "New"})
	assert.Equal(t, 1, addedCount, "Expected one new value, as 'Hello' is a duplicate")
}

func TestHashedReflectiveMultiMap_Remove(t *testing.T) {
	mm := umap.NewHashedReflectiveMultiMap[string, TestStruct](sha256.New())
	mm.Set("key", TestStruct{ID: 1, Name: "Remove"}, TestStruct{ID: 2, Name: "Keep"})

	removalCount := mm.Remove("key", func(v TestStruct) bool { return v.Name == "Remove" })
	assert.Equal(t, 1, removalCount, "Expected one value to be removed")
	values, _ := mm.Get("key")
	assert.Len(t, values, 1, "Expected one value left")
	assert.Equal(t, TestStruct{ID: 2, Name: "Keep"}, values[0], "Expected 'Keep' to remain")
}

func TestHashedReflectiveMultiMap_Clear(t *testing.T) {
	mm := umap.NewHashedReflectiveMultiMap[string, TestStruct](sha256.New())
	mm.Set("key", TestStruct{ID: 1, Name: "Data"})

	cleared := mm.Clear("key")
	assert.True(t, cleared, "Expected true, indicating values were cleared")

	_, exists := mm.Get("key")
	assert.False(t, exists, "Expected no values after clear")
}
