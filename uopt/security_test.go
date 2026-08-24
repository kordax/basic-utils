package uopt

import (
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScalarIntRejectsOutOfRangeValues(t *testing.T) {
	value, ok := scalarInt(reflect.ValueOf(uint64(1 << 63)))
	assert.False(t, ok)
	assert.Zero(t, value)

	value, ok = scalarInt(reflect.ValueOf(float64(0x1p63)))
	assert.False(t, ok)
	assert.Zero(t, value)

	value, ok = scalarInt(reflect.ValueOf(math.NaN()))
	assert.False(t, ok)
	assert.Zero(t, value)
}

func TestScalarIntAcceptsBoundaryValues(t *testing.T) {
	value, ok := scalarInt(reflect.ValueOf(uint64(1<<63 - 1)))
	assert.True(t, ok)
	assert.Equal(t, int64(1<<63-1), value)

	value, ok = scalarInt(reflect.ValueOf(float64(-0x1p63)))
	assert.True(t, ok)
	assert.Equal(t, int64(-1<<63), value)
}

func TestScalarUintRejectsOutOfRangeValues(t *testing.T) {
	value, ok := scalarUint(reflect.ValueOf(int64(-1)))
	assert.False(t, ok)
	assert.Zero(t, value)

	value, ok = scalarUint(reflect.ValueOf(float64(0x1p64)))
	assert.False(t, ok)
	assert.Zero(t, value)

	value, ok = scalarUint(reflect.ValueOf(math.NaN()))
	assert.False(t, ok)
	assert.Zero(t, value)
}

func TestScalarUintPreservesValidConversions(t *testing.T) {
	value, ok := scalarUint(reflect.ValueOf(int64(42)))
	assert.True(t, ok)
	assert.Equal(t, uint64(42), value)

	value, ok = scalarUint(reflect.ValueOf(float64(42.75)))
	assert.True(t, ok)
	assert.Equal(t, uint64(42), value)
}
