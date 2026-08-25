/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package upair

import (
	"fmt"
	"testing"

	"github.com/kordax/basic-utils/v4/uconst"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type comparableValue struct{}

func (comparableValue) Equals(_ uconst.Comparable) bool { return false }

func TestNewPair(t *testing.T) {
	pair := NewPair(1, "test")

	require.NotNil(t, pair, "NewPair should not return nil")
	assert.Equal(t, 1, pair.GetLeft(), "GetLeft should return the correct left value")
	assert.Equal(t, "test", pair.GetRight(), "GetRight should return the correct right value")
}

func TestOf(t *testing.T) {
	pair := Of(1, "test")

	assert.Equal(t, 1, pair.GetLeft(), "GetLeft should return the correct left value")
	assert.Equal(t, "test", pair.GetRight(), "GetRight should return the correct right value")

	left, right := pair.Values()
	assert.Equal(t, 1, left)
	assert.Equal(t, "test", right)
}

func TestPairSwap(t *testing.T) {
	pair := Of(1, "test")

	swapped := pair.Swap()

	assert.Equal(t, "test", swapped.Left)
	assert.Equal(t, 1, swapped.Right)
}

func TestPairMap(t *testing.T) {
	pair := Of(1, "test")

	leftMapped := pair.MapLeft(func(v int) float64 {
		return float64(v) * 2
	})
	rightMapped := pair.MapRight(func(v string) int {
		return len(v)
	})
	mapped := pair.Map(func(v int) float64 {
		return float64(v) * 2
	}, func(v string) int {
		return len(v)
	})

	assert.Equal(t, Of(2.0, "test"), leftMapped)
	assert.Equal(t, Of(1, 4), rightMapped)
	assert.Equal(t, Of(2.0, 4), mapped)
}

func TestNewCPair(t *testing.T) {
	pair := NewCPair(236, "ctest")

	require.NotNil(t, pair, "NewCPair should not return nil")
	assert.Equal(t, 236, pair.GetLeft(), "GetLeft should return the correct left value")
	assert.Equal(t, "ctest", pair.GetRight(), "GetRight should return the correct right value")

	pair2 := NewCPair(pair.Left, pair.Right)
	assert.True(t, pair.Equals(pair2), "pair.Equals(%v) should return %v", pair, pair2)
}

func TestCOf(t *testing.T) {
	pair := COf(236, "ctest")

	assert.Equal(t, 236, pair.GetLeft(), "GetLeft should return the correct left value")
	assert.Equal(t, "ctest", pair.GetRight(), "GetRight should return the correct right value")
	assert.True(t, pair.Equals(COf(236, "ctest")))
	assert.True(t, pair.Equals(NewCPair(236, "ctest")))
	assert.False(t, pair.Equals(COf(236, "other")))
}

func TestPairDeprecatedMapFunctions(t *testing.T) {
	pair := Of(3, "a")

	mappedLeft := MapLeft(pair, func(v int) string {
		return fmt.Sprintf("value-%d", v)
	})
	mappedRight := MapRight(pair, func(v string) int {
		return len(v)
	})

	assert.Equal(t, Of("value-3", "a"), mappedLeft)
	assert.Equal(t, Of(3, 1), mappedRight)
}

func TestCPair_EqualsBranches(t *testing.T) {
	var _ uconst.Comparable = comparableValue{}

	pair := COf("left", "right")

	var ptr *CPair[string, string]
	assert.False(t, pair.Equals(ptr))
	assert.False(t, pair.Equals(comparableValue{}))
}
