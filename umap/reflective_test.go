package umap_test

import (
	"testing"

	"github.com/kordax/basic-utils/umap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReflectiveMultiMap_Get_Set(t *testing.T) {
	mm := umap.NewReflectiveMultiMap[string, string]()

	mm.Set("key", "value1")
	values, ok := mm.Get("key")
	require.True(t, ok)
	require.Len(t, values, 1)
	assert.Contains(t, values, "value1")

	mm.Set("key", "value2", "value3")
	values, _ = mm.Get("key")
	require.Len(t, values, 2)
	assert.Contains(t, values, "value2")
	assert.Contains(t, values, "value3")
}

func TestReflectiveMultiMap_Append(t *testing.T) {
	mm := umap.NewReflectiveMultiMap[string, string]()
	count := mm.Append("key", "value1")
	assert.Equal(t, 0, count) // No duplicates initially

	count = mm.Append("key", "value1", "value2")
	assert.Equal(t, 1, count) // One duplicate of value1

	values, ok := mm.Get("key")
	require.True(t, ok)
	assert.Len(t, values, 3) // Includes duplicates
}

func TestReflectiveMultiMap_Remove(t *testing.T) {
	mm := umap.NewReflectiveMultiMap[string, string]()
	mm.Set("key", "value1", "value2", "value1") // Intentional duplicate for testing

	// Remove specific values
	removed := mm.Remove("key", func(v string) bool { return v == "value1" })
	assert.Equal(t, 2, removed) // Two value1s removed

	// Verify removal
	values, _ := mm.Get("key")
	assert.Len(t, values, 1)
	assert.Contains(t, values, "value2")
}

func TestReflectiveMultiMap_Clear(t *testing.T) {
	mm := umap.NewReflectiveMultiMap[string, string]()
	mm.Set("key", "value1")

	// Clear and verify
	existed := mm.Clear("key")
	assert.True(t, existed)

	// Verify clear
	_, ok := mm.Get("key")
	assert.False(t, ok)
}

func TestReflectiveMultiMap_Collisions(t *testing.T) {
	// Simulate hash collision by overriding computeHash function, assuming implementation allows
	mm := umap.NewReflectiveMultiMap[string, struct{ A, B string }]()
	customValue1 := struct{ A, B string }{"one", "two"}
	customValue2 := struct{ A, B string }{"three", "four"} // Assume same hash for testing

	mm.Append("collision", customValue1)
	mm.Append("collision", customValue2)
	values, ok := mm.Get("collision")
	require.True(t, ok)
	require.Len(t, values, 2)

	// Ensure that distinct struct instances are stored despite the collision
	assert.Contains(t, values, customValue1)
	assert.Contains(t, values, customValue2)
}
