package uopt_test

import (
	"strconv"
	"testing"

	"git.casinomodule.org/casino27/basic-utils/v4/uopt"
	"github.com/stretchr/testify/assert"
)

func TestFunctionalHelpers(t *testing.T) {
	assert.Equal(t, "2", uopt.Of(2).Map(func(v int) string {
		return strconv.Itoa(v)
	}).OrElse(""))

	assert.Equal(t, 4, uopt.Of(2).FlatMap(func(v int) uopt.Opt[int] {
		return uopt.Of(v * 2)
	}).OrElse(0))

	assert.True(t, uopt.Of(2).Filter(func(v int) bool { return v%2 == 0 }).Present())
	assert.False(t, uopt.Of(3).Filter(func(v int) bool { return v%2 == 0 }).Present())
}

func TestOrElseGet(t *testing.T) {
	called := false
	assert.Equal(t, 10, uopt.Of(10).OrElseGet(func() int {
		called = true
		return 20
	}))
	assert.False(t, called)

	assert.Equal(t, 20, uopt.Null[int]().OrElseGet(func() int {
		called = true
		return 20
	}))
	assert.True(t, called)
}
