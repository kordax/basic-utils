/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap_test

import (
	"crypto/sha256"
	"encoding/binary"
	"testing"

	"github.com/kordax/basic-utils/umap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ umap.MultiMap[string, int] = (*umap.UniqueMultiMap[string, int])(nil)

type MyValue struct {
	Value string
}

func (v MyValue) Hash() int64 {
	hasher := sha256.New()
	hasher.Write([]byte(v.Value))
	hashBytes := hasher.Sum(nil)

	return int64(binary.LittleEndian.Uint64(hashBytes[:8]))
}

func TestUniqueMultiMap_Get(t *testing.T) {
	mm := umap.NewUniqueMultiMap[string, MyValue]()
	_, exists := mm.Get("key")
	assert.False(t, exists, "Expected no values for key")

	mm.Append("key", MyValue{Value: "hello"})
	values, exists := mm.Get("key")
	require.True(t, exists, "Expected values to exist")
	require.Len(t, values, 1, "Expected one value")
}

func TestUniqueMultiMap_Set(t *testing.T) {
	mm := umap.NewUniqueMultiMap[string, MyValue]()
	addedCount := mm.Set("key", MyValue{Value: "hello"}, MyValue{Value: "world"})
	assert.Equal(t, 2, addedCount, "Expected two unique values to be added")

	addedCount = mm.Set("key", MyValue{Value: "new"})
	assert.Equal(t, 1, addedCount, "Expected one new value to be added")
	values, _ := mm.Get("key")
	require.Len(t, values, 1, "Expected one value after set")
}

func TestUniqueMultiMap_Append(t *testing.T) {
	mm := umap.NewUniqueMultiMap[string, MyValue]()
	mm.Append("key", MyValue{Value: "hello"})
	mm.Append("key", MyValue{Value: "world"})

	addedCount := mm.Append("key", MyValue{Value: "hello"}, MyValue{Value: "new"})
	assert.Equal(t, 1, addedCount, "Expected one new value, as 'hello' is a duplicate")
}

func TestUniqueMultiMap_Remove(t *testing.T) {
	mm := umap.NewUniqueMultiMap[string, MyValue]()
	mm.Set("key", MyValue{Value: "remove"}, MyValue{Value: "keep"})

	removalCount := mm.Remove("key", func(v MyValue) bool { return v.Value == "remove" })
	assert.Equal(t, 1, removalCount, "Expected one value to be removed")
	values, _ := mm.Get("key")
	assert.Len(t, values, 1, "Expected one value left")
	assert.Equal(t, "keep", values[0].Value, "Expected 'keep' to remain")
}

func TestUniqueMultiMap_Clear(t *testing.T) {
	mm := umap.NewUniqueMultiMap[string, MyValue]()
	mm.Set("key", MyValue{Value: "data"})

	cleared := mm.Clear("key")
	assert.True(t, cleared, "Expected true, indicating values were cleared")

	_, exists := mm.Get("key")
	assert.False(t, exists, "Expected no values after clear")
}

func TestUniqueMultiMap_Iterator(t *testing.T) {
	mm := umap.NewUniqueMultiMap[string, int]()
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
