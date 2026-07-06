package uset_test

import (
	"testing"

	"git.casinomodule.org/casino27/basic-utils/v2/uset"
	"github.com/stretchr/testify/assert"
)

func TestSetAlgebra(t *testing.T) {
	left := uset.NewHashSet(1, 2, 3)
	right := uset.NewHashSet(3, 4)

	assert.ElementsMatch(t, []int{1, 2, 3, 4}, uset.Union[int](left, right).Values())
	assert.ElementsMatch(t, []int{3}, uset.Intersect[int](left, right).Values())
	assert.ElementsMatch(t, []int{1, 2}, uset.Difference[int](left, right).Values())
	assert.True(t, uset.IsSubset(uset.NewHashSet(1, 2), left))
	assert.True(t, uset.IsSuperset(left, uset.NewHashSet(1, 2)))
	assert.False(t, uset.IsSubset(right, left))
}
