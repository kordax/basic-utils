/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uarray_test

import (
	"testing"

	"github.com/kordax/basic-utils/uarray"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringToInt(t *testing.T) {
	val := "42"
	result := uarray.StringToInt(&val)
	assert.Equal(t, 42, result)
}

func TestStringToInt32(t *testing.T) {
	val := "32"
	result := uarray.StringToInt32(&val)
	assert.Equal(t, int32(32), result)
}

func TestStringToInt64(t *testing.T) {
	val := "64"
	result := uarray.StringToInt64(&val)
	assert.Equal(t, int64(64), result)
}

func TestStringToFloat32(t *testing.T) {
	val := "0.32"
	result := uarray.StringToFloat32(&val)
	assert.InDelta(t, float32(0.32), result, 0.0001)
}

func TestStringToFloat64(t *testing.T) {
	val := "0.64"
	result := uarray.StringToFloat64(&val)
	assert.InDelta(t, float64(0.64), result, 0.0001)
}

func TestStringToBool(t *testing.T) {
	trueVal := "true"
	falseVal := "false"
	resultTrue := uarray.StringToBool(&trueVal)
	resultFalse := uarray.StringToBool(&falseVal)
	assert.True(t, resultTrue)
	assert.False(t, resultFalse)
}

func TestFloat64ToFloat32(t *testing.T) {
	val := float64(3.14)
	result := uarray.Float64ToFloat32(&val)
	assert.InDelta(t, float32(3.14), result, 0.0001)
}

func TestInt64ToInt32(t *testing.T) {
	val := int64(123)
	result := uarray.Int64ToInt32(&val)
	assert.Equal(t, int32(123), result)
}

func TestIntToString(t *testing.T) {
	val := 42
	result := uarray.IntToString(&val)
	assert.Equal(t, "42", result)
}

func TestInt64ToString(t *testing.T) {
	val := int64(64)
	result := uarray.Int64ToString(&val)
	assert.Equal(t, "64", result)
}

func TestFloat32ToString(t *testing.T) {
	val := float32(3.14)
	result := uarray.Float32ToString(&val)
	assert.Equal(t, "3.14", result)
}

func TestFloat64ToString(t *testing.T) {
	val := float64(6.28)
	result := uarray.Float64ToString(&val)
	assert.Equal(t, "6.28", result)
}

func TestBoolToString(t *testing.T) {
	trueVal := true
	falseVal := false
	resultTrue := uarray.BoolToString(&trueVal)
	resultFalse := uarray.BoolToString(&falseVal)
	assert.Equal(t, "true", resultTrue)
	assert.Equal(t, "false", resultFalse)
}

func TestMapStringToInt64(t *testing.T) {
	values := []string{"10", "20", "30"}
	expected := []int64{10, 20, 30}
	result := uarray.Map(values, uarray.StringToInt64)
	require.Equal(t, expected, result)
}

func TestMapStringToFloat64(t *testing.T) {
	values := []string{"1.1", "2.2", "3.3"}
	expected := []float64{1.1, 2.2, 3.3}
	result := uarray.Map(values, uarray.StringToFloat64)
	assert.InDeltaSlice(t, expected, result, 0.0001)
}

func TestMapStringToBool(t *testing.T) {
	values := []string{"true", "false", "true"}
	expected := []bool{true, false, true}
	result := uarray.Map(values, uarray.StringToBool)
	assert.Equal(t, expected, result)
}

func TestMapFloat64ToFloat32(t *testing.T) {
	values := []float64{1.5, 2.5, 3.5}
	expected := []float32{1.5, 2.5, 3.5}
	result := uarray.Map(values, uarray.Float64ToFloat32)
	assert.InDeltaSlice(t, expected, result, 0.0001)
}

func TestMapInt64ToInt32(t *testing.T) {
	values := []int64{100, 200, 300}
	expected := []int32{100, 200, 300}
	result := uarray.Map(values, uarray.Int64ToInt32)
	assert.Equal(t, expected, result)
}

func TestMapInt64ToString(t *testing.T) {
	values := []int64{123, 456, 789}
	expected := []string{"123", "456", "789"}
	result := uarray.Map(values, uarray.Int64ToString)
	assert.Equal(t, expected, result)
}

func TestMapFloat32ToString(t *testing.T) {
	values := []float32{1.1, 2.2, 3.3}
	expected := []string{"1.1", "2.2", "3.3"}
	result := uarray.Map(values, uarray.Float32ToString)
	assert.Equal(t, expected, result)
}

func TestMapBoolToString(t *testing.T) {
	values := []bool{true, false, true}
	expected := []string{"true", "false", "true"}
	result := uarray.Map(values, uarray.BoolToString)
	assert.Equal(t, expected, result)
}
