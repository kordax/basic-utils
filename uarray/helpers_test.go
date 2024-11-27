/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uarray_test

import (
	"fmt"
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

func TestInt8ToString(t *testing.T) {
	testCases := []struct {
		input    int8
		expected string
	}{
		{input: 0, expected: "0"},
		{input: 42, expected: "42"},
		{input: -42, expected: "-42"},
		{input: 127, expected: "127"},   // Max int8
		{input: -128, expected: "-128"}, // Min int8
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Int8ToString_%d", tc.input), func(t *testing.T) {
			result := uarray.Int8ToString(&tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestInt16ToString(t *testing.T) {
	testCases := []struct {
		input    int16
		expected string
	}{
		{input: 0, expected: "0"},
		{input: 32000, expected: "32000"},
		{input: -32000, expected: "-32000"},
		{input: 32767, expected: "32767"},   // Max int16
		{input: -32768, expected: "-32768"}, // Min int16
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Int16ToString_%d", tc.input), func(t *testing.T) {
			result := uarray.Int16ToString(&tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestInt32ToString(t *testing.T) {
	testCases := []struct {
		input    int32
		expected string
	}{
		{input: 0, expected: "0"},
		{input: 2147483647, expected: "2147483647"},   // Max int32
		{input: -2147483648, expected: "-2147483648"}, // Min int32
		{input: 123456, expected: "123456"},
		{input: -123456, expected: "-123456"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Int32ToString_%d", tc.input), func(t *testing.T) {
			result := uarray.Int32ToString(&tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestInt64ToString(t *testing.T) {
	testCases := []struct {
		input    int64
		expected string
	}{
		{input: 0, expected: "0"},
		{input: 9223372036854775807, expected: "9223372036854775807"},   // Max int64
		{input: -9223372036854775808, expected: "-9223372036854775808"}, // Min int64
		{input: 123456789012345, expected: "123456789012345"},
		{input: -123456789012345, expected: "-123456789012345"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Int64ToString_%d", tc.input), func(t *testing.T) {
			result := uarray.Int64ToString(&tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestUint8ToString(t *testing.T) {
	testCases := []struct {
		input    uint8
		expected string
	}{
		{input: 0, expected: "0"},
		{input: 255, expected: "255"}, // Max uint8
		{input: 128, expected: "128"},
		{input: 1, expected: "1"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Uint8ToString_%d", tc.input), func(t *testing.T) {
			result := uarray.Uint8ToString(&tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestUint16ToString(t *testing.T) {
	testCases := []struct {
		input    uint16
		expected string
	}{
		{input: 0, expected: "0"},
		{input: 65535, expected: "65535"}, // Max uint16
		{input: 32768, expected: "32768"},
		{input: 1, expected: "1"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Uint16ToString_%d", tc.input), func(t *testing.T) {
			result := uarray.Uint16ToString(&tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestUint32ToString(t *testing.T) {
	testCases := []struct {
		input    uint32
		expected string
	}{
		{input: 0, expected: "0"},
		{input: 4294967295, expected: "4294967295"}, // Max uint32
		{input: 2147483648, expected: "2147483648"},
		{input: 1, expected: "1"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Uint32ToString_%d", tc.input), func(t *testing.T) {
			result := uarray.Uint32ToString(&tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestUint64ToString(t *testing.T) {
	testCases := []struct {
		input    uint64
		expected string
	}{
		{input: 0, expected: "0"},
		{input: 18446744073709551615, expected: "18446744073709551615"}, // Max uint64
		{input: 9223372036854775808, expected: "9223372036854775808"},
		{input: 1, expected: "1"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Uint64ToString_%d", tc.input), func(t *testing.T) {
			result := uarray.Uint64ToString(&tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestInt8ToString_Negative(t *testing.T) {
	val := int8(-100)
	result := uarray.Int8ToString(&val)
	assert.Equal(t, "-100", result)
}

func TestInt16ToString_Negative(t *testing.T) {
	val := int16(-30000)
	result := uarray.Int16ToString(&val)
	assert.Equal(t, "-30000", result)
}

func TestInt32ToString_Negative(t *testing.T) {
	val := int32(-200000)
	result := uarray.Int32ToString(&val)
	assert.Equal(t, "-200000", result)
}

func TestInt64ToString_Negative(t *testing.T) {
	val := int64(-9000000000)
	result := uarray.Int64ToString(&val)
	assert.Equal(t, "-9000000000", result)
}

func TestUint8ToString_Max(t *testing.T) {
	val := uint8(255)
	result := uarray.Uint8ToString(&val)
	assert.Equal(t, "255", result)
}

func TestUint16ToString_Max(t *testing.T) {
	val := uint16(65535)
	result := uarray.Uint16ToString(&val)
	assert.Equal(t, "65535", result)
}

func TestUint32ToString_Max(t *testing.T) {
	val := uint32(4294967295)
	result := uarray.Uint32ToString(&val)
	assert.Equal(t, "4294967295", result)
}

func TestUint64ToString_Max(t *testing.T) {
	val := uint64(18446744073709551615)
	result := uarray.Uint64ToString(&val)
	assert.Equal(t, "18446744073709551615", result)
}

func TestMapInt8ToString(t *testing.T) {
	values := []int8{10, -20, 30}
	expected := []string{"10", "-20", "30"}
	result := uarray.Map(values, uarray.Int8ToString)
	require.Equal(t, expected, result)
}

func TestMapInt16ToString(t *testing.T) {
	values := []int16{-1000, 0, 1000}
	expected := []string{"-1000", "0", "1000"}
	result := uarray.Map(values, uarray.Int16ToString)
	require.Equal(t, expected, result)
}

func TestMapInt32ToString(t *testing.T) {
	values := []int32{-100000, 0, 100000}
	expected := []string{"-100000", "0", "100000"}
	result := uarray.Map(values, uarray.Int32ToString)
	require.Equal(t, expected, result)
}

func TestMapUint8ToString(t *testing.T) {
	values := []uint8{0, 128, 255}
	expected := []string{"0", "128", "255"}
	result := uarray.Map(values, uarray.Uint8ToString)
	require.Equal(t, expected, result)
}

func TestMapUint16ToString(t *testing.T) {
	values := []uint16{0, 32768, 65535}
	expected := []string{"0", "32768", "65535"}
	result := uarray.Map(values, uarray.Uint16ToString)
	require.Equal(t, expected, result)
}

func TestMapUint32ToString(t *testing.T) {
	values := []uint32{0, 2147483648, 4294967295}
	expected := []string{"0", "2147483648", "4294967295"}
	result := uarray.Map(values, uarray.Uint32ToString)
	require.Equal(t, expected, result)
}

func TestMapUint64ToString(t *testing.T) {
	values := []uint64{0, 9223372036854775808, 18446744073709551615}
	expected := []string{"0", "9223372036854775808", "18446744073709551615"}
	result := uarray.Map(values, uarray.Uint64ToString)
	require.Equal(t, expected, result)
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
