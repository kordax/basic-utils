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

type MyValue struct {
	Value string
}

func (v MyValue) Hash() int64 {
	hasher := sha256.New()
	hasher.Write([]byte(v.Value))
	hashBytes := hasher.Sum(nil)

	return int64(binary.LittleEndian.Uint64(hashBytes[:8]))
}

func TestHashedMultiMap_Get(t *testing.T) {
	mm := umap.NewHashedMultiMap[string, MyValue]()
	_, exists := mm.Get("key")
	assert.False(t, exists, "Expected no values for key")

	mm.Append("key", MyValue{Value: "hello"})
	values, exists := mm.Get("key")
	require.True(t, exists, "Expected values to exist")
	require.Len(t, values, 1, "Expected one value")
}

func TestHashedMultiMap_Set(t *testing.T) {
	mm := umap.NewHashedMultiMap[string, MyValue]()
	addedCount := mm.Set("key", MyValue{Value: "hello"}, MyValue{Value: "world"})
	assert.Equal(t, 2, addedCount, "Expected two unique values to be added")

	addedCount = mm.Set("key", MyValue{Value: "new"})
	assert.Equal(t, 1, addedCount, "Expected one new value to be added")
	values, _ := mm.Get("key")
	require.Len(t, values, 1, "Expected one value after set")
}

func TestHashedMultiMap_Append(t *testing.T) {
	mm := umap.NewHashedMultiMap[string, MyValue]()
	mm.Append("key", MyValue{Value: "hello"})
	mm.Append("key", MyValue{Value: "world"})

	addedCount := mm.Append("key", MyValue{Value: "hello"}, MyValue{Value: "new"})
	assert.Equal(t, 1, addedCount, "Expected one new value, as 'hello' is a duplicate")
}

func TestHashedMultiMap_Remove(t *testing.T) {
	mm := umap.NewHashedMultiMap[string, MyValue]()
	mm.Set("key", MyValue{Value: "remove"}, MyValue{Value: "keep"})

	removalCount := mm.Remove("key", func(v MyValue) bool { return v.Value == "remove" })
	assert.Equal(t, 1, removalCount, "Expected one value to be removed")
	values, _ := mm.Get("key")
	assert.Len(t, values, 1, "Expected one value left")
	assert.Equal(t, "keep", values[0].Value, "Expected 'keep' to remain")
}

func TestHashedMultiMap_Clear(t *testing.T) {
	mm := umap.NewHashedMultiMap[string, MyValue]()
	mm.Set("key", MyValue{Value: "data"})

	cleared := mm.Clear("key")
	assert.True(t, cleared, "Expected true, indicating values were cleared")

	_, exists := mm.Get("key")
	assert.False(t, exists, "Expected no values after clear")
}
