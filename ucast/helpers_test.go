/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ucast_test

import (
	"fmt"
	"testing"

	"github.com/kordax/basic-utils/ucast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringToInt(t *testing.T) {
	val := "42"
	result, err := ucast.StringToInt(&val)
	require.NoError(t, err)
	assert.Equal(t, 42, result)
}

func TestStringToInt8(t *testing.T) {
	val := "42"
	result, err := ucast.StringToInt8(&val)
	require.NoError(t, err)
	assert.Equal(t, int8(42), result)
}

func TestStringToInt16(t *testing.T) {
	val := "42"
	result, err := ucast.StringToInt16(&val)
	require.NoError(t, err)
	assert.Equal(t, int16(42), result)
}

func TestStringToInt32(t *testing.T) {
	val := "32"
	result, err := ucast.StringToInt32(&val)
	require.NoError(t, err)
	assert.Equal(t, int32(32), result)
}

func TestStringToInt64(t *testing.T) {
	val := "64"
	result, _ := ucast.StringToInt64(&val)
	assert.Equal(t, int64(64), result)
}

func TestStringToUint(t *testing.T) {
	val := "42"
	result, err := ucast.StringToUint(&val)
	require.NoError(t, err)
	assert.Equal(t, uint(42), result)
}

func TestStringToUint8(t *testing.T) {
	val := "42"
	result, err := ucast.StringToUint8(&val)
	require.NoError(t, err)
	assert.Equal(t, uint8(42), result)
}

func TestStringToUint16(t *testing.T) {
	val := "42"
	result, err := ucast.StringToUint16(&val)
	require.NoError(t, err)
	assert.Equal(t, uint16(42), result)
}

func TestStringToUint32(t *testing.T) {
	val := "32"
	result, err := ucast.StringToUint32(&val)
	require.NoError(t, err)
	assert.Equal(t, uint32(32), result)
}

func TestStringToUint64(t *testing.T) {
	val := "64"
	result, _ := ucast.StringToUint64(&val)
	assert.Equal(t, uint64(64), result)
}

func TestStringToFloat32(t *testing.T) {
	val := "0.32"
	result, err := ucast.StringToFloat32(&val)
	require.NoError(t, err)
	assert.InDelta(t, float32(0.32), result, 0.0001)
}

func TestStringToFloat64(t *testing.T) {
	val := "0.64"
	result, err := ucast.StringToFloat64(&val)
	require.NoError(t, err)
	assert.InDelta(t, float64(0.64), result, 0.0001)
}

func TestStringToBool(t *testing.T) {
	trueVal := "true"
	falseVal := "false"
	resultTrue, err := ucast.StringToBool(&trueVal)
	require.NoError(t, err)
	resultFalse, err := ucast.StringToBool(&falseVal)
	require.NoError(t, err)
	assert.True(t, resultTrue)
	assert.False(t, resultFalse)
}

func TestFloat64ToFloat32(t *testing.T) {
	val := float64(3.14)
	result := ucast.Float64ToFloat32(&val)
	assert.InDelta(t, float32(3.14), result, 0.0001)
}

func TestIntToString(t *testing.T) {
	testCases := []struct {
		input    int
		expected string
	}{
		{input: 0, expected: "0"},
		{input: 42, expected: "42"},
		{input: -42, expected: "-42"},
		{input: 347347, expected: "347347"},
		{input: -3473718, expected: "-3473718"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("IntToString_%d", tc.input), func(t *testing.T) {
			result := ucast.IntToString(&tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
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
			result := ucast.Int8ToString(&tc.input)
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
			result := ucast.Int16ToString(&tc.input)
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
			result := ucast.Int32ToString(&tc.input)
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
			result := ucast.Int64ToString(&tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestUintToString(t *testing.T) {
	testCases := []struct {
		input    uint
		expected string
	}{
		{input: 0, expected: "0"},
		{input: 255, expected: "255"}, // Max uint8
		{input: 128, expected: "128"},
		{input: 1, expected: "1"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("UintToString_%d", tc.input), func(t *testing.T) {
			result := ucast.UintToString(&tc.input)
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
			result := ucast.Uint8ToString(&tc.input)
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
			result := ucast.Uint16ToString(&tc.input)
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
			result := ucast.Uint32ToString(&tc.input)
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
			result := ucast.Uint64ToString(&tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestInt8ToString_Negative(t *testing.T) {
	val := int8(-100)
	result := ucast.Int8ToString(&val)
	assert.Equal(t, "-100", result)
}

func TestInt16ToString_Negative(t *testing.T) {
	val := int16(-30000)
	result := ucast.Int16ToString(&val)
	assert.Equal(t, "-30000", result)
}

func TestInt32ToString_Negative(t *testing.T) {
	val := int32(-200000)
	result := ucast.Int32ToString(&val)
	assert.Equal(t, "-200000", result)
}

func TestInt64ToString_Negative(t *testing.T) {
	val := int64(-9000000000)
	result := ucast.Int64ToString(&val)
	assert.Equal(t, "-9000000000", result)
}

func TestUint8ToString_Max(t *testing.T) {
	val := uint8(255)
	result := ucast.Uint8ToString(&val)
	assert.Equal(t, "255", result)
}

func TestUint16ToString_Max(t *testing.T) {
	val := uint16(65535)
	result := ucast.Uint16ToString(&val)
	assert.Equal(t, "65535", result)
}

func TestUint32ToString_Max(t *testing.T) {
	val := uint32(4294967295)
	result := ucast.Uint32ToString(&val)
	assert.Equal(t, "4294967295", result)
}

func TestUint64ToString_Max(t *testing.T) {
	val := uint64(18446744073709551615)
	result := ucast.Uint64ToString(&val)
	assert.Equal(t, "18446744073709551615", result)
}

func TestFloat32ToString(t *testing.T) {
	val := float32(3.14)
	result := ucast.Float32ToString(&val)
	assert.Equal(t, "3.14", result)
}

func TestFloat64ToString(t *testing.T) {
	val := float64(6.28)
	result := ucast.Float64ToString(&val)
	assert.Equal(t, "6.28", result)
}

func TestBoolToString(t *testing.T) {
	trueVal := true
	falseVal := false
	resultTrue := ucast.BoolToString(&trueVal)
	resultFalse := ucast.BoolToString(&falseVal)
	assert.Equal(t, "true", resultTrue)
	assert.Equal(t, "false", resultFalse)
}
