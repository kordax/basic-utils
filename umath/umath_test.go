/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package umath_test

import (
	"math"
	"testing"

	"github.com/kordax/basic-utils/v2/umath"
	"github.com/stretchr/testify/assert"
)

func TestClosestMatch(t *testing.T) {
	slice := []uint32{1, 2, 3, 4, 5}

	// Test case 1: Closest match exists in the slice
	toMatch1 := uint32(3)
	expectedIndex1 := 2
	expectedValue1 := uint32(3)
	index1, value1 := umath.ClosestMatch(toMatch1, slice)
	if index1 != expectedIndex1 || value1 != expectedValue1 {
		t.Errorf("Test case 1 failed: ClosestMatch returned (%d, %d), expected (%d, %d)", index1, value1, expectedIndex1, expectedValue1)
	}

	// Test case 2: Closest match does not exist in the slice (value less than the minimum)
	toMatch2 := uint32(0)
	expectedIndex2 := 0
	expectedValue2 := uint32(1)
	index2, value2 := umath.ClosestMatch(toMatch2, slice)
	if index2 != expectedIndex2 || value2 != expectedValue2 {
		t.Errorf("Test case 2 failed: ClosestMatch returned (%d, %d), expected (%d, %d)", index2, value2, expectedIndex2, expectedValue2)
	}

	// Test case 3: Closest match does not exist in the slice (value greater than the maximum)
	toMatch3 := uint32(6)
	expectedIndex3 := 4
	expectedValue3 := uint32(5)
	index3, value3 := umath.ClosestMatch(toMatch3, slice)
	if index3 != expectedIndex3 || value3 != expectedValue3 {
		t.Errorf("Test case 3 failed: ClosestMatch returned (%d, %d), expected (%d, %d)", index3, value3, expectedIndex3, expectedValue3)
	}
}

func TestValOrMin(t *testing.T) {
	if umath.ValOrMin(5, 10) != 10 {
		t.Error("Expected minimum value of 10")
	}

	if umath.ValOrMin(15, 10) != 15 {
		t.Error("Expected value of 15")
	}
}

func TestAvg(t *testing.T) {
	array := []int{1, 2, 3, 4, 5}
	if umath.Avg(array) != 3 {
		t.Error("Expected average value of 3")
	}
}

func TestAvgFloat(t *testing.T) {
	tests := []struct {
		name     string
		array    []float64
		expected float64
	}{
		{
			name:     "With large values",
			array:    []float64{math.MaxInt32, math.MaxInt32, math.MaxInt32},
			expected: float64(math.MaxInt32),
		},
		{
			name:     "With negative large values",
			array:    []float64{math.MinInt32, math.MinInt32, math.MinInt32},
			expected: float64(math.MinInt32),
		},
		{
			name:     "Mixed large and small values",
			array:    []float64{math.MaxInt32, 0, math.MinInt32},
			expected: -0.3333333333333333,
		},
		{
			name:     "Repeated number sequence",
			array:    []float64{1, 2, 3, 1, 2, 3, 1, 2, 3},
			expected: 2,
		},
		{
			name:     "With zero crossings",
			array:    []float64{-3, -2, -1, 0, 1, 2, 3},
			expected: 0,
		},
		{
			name:     "Array of length just below max int32",
			array:    make([]float64, 100),
			expected: 0, // since all are initialized to 0
		},
		{
			name:     "Array with multiple zeros",
			array:    []float64{0, 0, 0, 0, 0},
			expected: 0,
		},
		{
			name:     "Array with alternating zeros and ones",
			array:    []float64{0, 1, 0, 1, 0, 1, 0, 1, 0},
			expected: 0.4444444444444444, // 4 ones divided by 9
		},
		{
			name:     "Involving very small numbers",
			array:    []float64{1e-10, 1e-10, 1e-10},
			expected: 1e-10,
		},
		{
			name:     "Involving very large numbers",
			array:    []float64{1e10, 1e10, 1e10},
			expected: 1e10,
		},
		{
			name:     "Mix of very small and very large numbers",
			array:    []float64{1e-10, 1e10, 1e-10},
			expected: (1e-10 + 1e10 + 1e-10) / 3,
		},
		{
			name:     "Numbers around machine epsilon",
			array:    []float64{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64},
			expected: math.SmallestNonzeroFloat64,
		},
		{
			name:     "A sequence with expected average close to zero",
			array:    []float64{1.1, 2.2, 3.3, -1.1, -2.2, -3.3},
			expected: 0,
		},
		{
			name:     "Values leading to potential precision loss",
			array:    []float64{0.1, 0.2, 0.3},
			expected: 0.2,
		},
	}

	tolerance := 1e-16
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if result := umath.AvgFloat(test.array); math.Abs(result-test.expected) > tolerance {
				t.Errorf("Expected average of %v to be %v, but got %v", test.array, test.expected, result)
			}
		})
	}
}

func TestMin(t *testing.T) {
	arrayInt := []int{10, 30, 50, 20, 40}
	if umath.Min(arrayInt) != 10 {
		t.Error("Expected minimum value of 10 for int slice")
	}

	arrayFloat64 := []float64{10.5, 30.5, 50.5, 20.5, 40.5}
	if umath.Min(arrayFloat64) != 10.5 {
		t.Error("Expected minimum value of 10.5 for float64 slice")
	}
}

func TestMax(t *testing.T) {
	arrayInt := []int{10, 30, 50, 20, 40}
	if umath.Max(arrayInt) != 50 {
		t.Error("Expected maximum value of 50 for int slice")
	}

	arrayFloat64 := []float64{10.5, 30.5, 50.5, 20.5, 40.5}
	if umath.Max(arrayFloat64) != 50.5 {
		t.Error("Expected maximum value of 50.5 for float64 slice")
	}
}

func TestMinMaxInt(t *testing.T) {
	array := []int{10, 30, 50, 20, 40}
	mn, mx := umath.MinMaxInt(array)
	if mn != 10 || mx != 50 {
		t.Errorf("Expected min 10 and max 50, got min %d and max %d", mn, mx)
	}
}

func TestMedInt(t *testing.T) {
	array := []int{10, 20, 30, 40, 50}
	if umath.Med(array) != 30 {
		t.Error("Expected median value of 30")
	}
}

func TestMinMax(t *testing.T) {
	arrayInt := []int{10, 20, 30, 40, 50}
	minInt, maxInt := umath.MinMax(arrayInt)
	if minInt != 10 || maxInt != 50 {
		t.Errorf("Expected min 10 and max 50 for int slice, got min %d and max %d", minInt, maxInt)
	}

	arrayFloat64 := []float64{10.5, 20.5, 30.5, 40.5, 50.5}
	minFloat64, maxFloat64 := umath.MinMax(arrayFloat64)
	if minFloat64 != 10.5 || maxFloat64 != 50.5 {
		t.Errorf("Expected min 10.5 and max 50.5 for float64 slice, got min %f and max %f", minFloat64, maxFloat64)
	}
}

func TestMinMaxFromMap(t *testing.T) {
	mapInt := map[string]int{"a": 10, "b": 20, "c": 30, "d": 40, "e": 50}
	minInt, maxInt := umath.MinMaxFromMap(mapInt)
	if minInt != 10 || maxInt != 50 {
		t.Errorf("Expected min 10 and max 50 for int map, got min %d and max %d", minInt, maxInt)
	}

	mapFloat64 := map[string]float64{"a": 10.5, "b": 20.5, "c": 30.5, "d": 40.5, "e": 50.5}
	minFloat64, maxFloat64 := umath.MinMaxFromMap(mapFloat64)
	if minFloat64 != 10.5 || maxFloat64 != 50.5 {
		t.Errorf("Expected min 10.5 and max 50.5 for float64 map, got min %f and max %f", minFloat64, maxFloat64)
	}
}

func TestRoundWithPrecision(t *testing.T) {
	assert.Equal(t, 10.46, umath.RoundWithPrecision(10.4567, 2), "Expected rounded value of 10.46")
	assert.Equal(t, 10.45, umath.RoundWithPrecision(10.453, 2), "Expected rounded value of 10.45")
	assert.Equal(t, 0.0, umath.RoundWithPrecision(0.1, 0), "Expected rounded value of 1")
}

func TestRoundUp(t *testing.T) {
	assert.Equal(t, 2.0, umath.RoundUp(1.1), "1.1 should round up to 2")
	assert.Equal(t, -1.0, umath.RoundUp(-1.1), "-1.1 should round up to -1")
	assert.Equal(t, 1, umath.RoundUp(1), "Integers should remain unchanged")
	assert.Equal(t, -1, umath.RoundUp(-1), "Negative integers should remain unchanged")
	assert.Equal(t, 1.0, umath.RoundUp(0.9999), "0.9999 should round up to 1")
	assert.Equal(t, 1.0, umath.RoundUp(0.0001), "0.0001 should round up to 1")
	assert.Equal(t, 0.0, umath.RoundUp(math.SmallestNonzeroFloat64*0.1), "Expected SmallestNonzeroFloat64 to round up to 1")
	assert.Equal(t, 1000001.0, umath.RoundUp(1000000.1), "1000000.1 should round up to 1000001")
	assert.Equal(t, 0.0, umath.RoundUp(0.0), "0 should remain as 0")
	assert.Equal(t, 2.0, umath.RoundUp(1.9999), "1.9999 should round up to 2")
	assert.Equal(t, 1.0, umath.RoundUp(math.SmallestNonzeroFloat64), "Smallest positive float should round up to 1")
	assert.Equal(t, -2.0, umath.RoundUp(-2.1), "-2.1 should round up to -2")
	assert.Equal(t, -2.0, umath.RoundUp(-2.9), "-2.9 should round up to -2")
	assert.Equal(t, 10.0, umath.RoundUp(9.1), "9.1 should round up to 10")
}

func TestAbsVal(t *testing.T) {
	// Test for int type
	intResult := umath.AbsVal[int](-5)
	if intResult != 5 {
		t.Errorf("Expected AbsVal of -5 to be 5 for int type, but got %v", intResult)
	}

	intResultPos := umath.AbsVal[int](5)
	if intResultPos != 5 {
		t.Errorf("Expected AbsVal of 5 to be 5 for int type, but got %v", intResultPos)
	}

	// Test for float32 type
	floatResult := umath.AbsVal[float32](-5.5)
	if floatResult != 5.5 {
		t.Errorf("Expected AbsVal of -5.5 to be 5.5 for float32 type, but got %v", floatResult)
	}

	floatResultPos := umath.AbsVal[float32](5.5)
	if floatResultPos != 5.5 {
		t.Errorf("Expected AbsVal of 5.5 to be 5.5 for float32 type, but got %v", floatResultPos)
	}
}

func TestMaxValue(t *testing.T) {
	// float32
	if val := umath.MaxValue[float32](); val != math.MaxFloat32 {
		t.Errorf("Expected %v for float32, got %v", math.Float32frombits(^uint32(0)>>1), val)
	}

	// float64
	if val := umath.MaxValue[float64](); val != math.MaxFloat64 {
		t.Errorf("Expected %v for float64, got %v", math.Float64frombits(^uint64(0)>>1), val)
	}

	// int
	if val := umath.MaxValue[int](); val != math.MaxInt {
		t.Errorf("Expected %v for int, got %v", math.MaxInt, val)
	}

	// int8
	if val := umath.MaxValue[int8](); val != int8(math.MaxInt8) {
		t.Errorf("Expected %v for int8, got %v", int8(math.MaxInt8), val)
	}

	// int16
	if val := umath.MaxValue[int16](); val != int16(math.MaxInt16) {
		t.Errorf("Expected %v for int16, got %v", int16(math.MaxInt16), val)
	}

	// int32
	if val := umath.MaxValue[int32](); val != int32(math.MaxInt32) {
		t.Errorf("Expected %v for int32, got %v", int32(math.MaxInt32), val)
	}

	// int64
	if val := umath.MaxValue[int64](); val != int64(math.MaxInt64) {
		t.Errorf("Expected %v for int64, got %v", int64(math.MaxInt64), val)
	}

	// uint
	if val := umath.MaxValue[uint](); val != ^uint(0) {
		t.Errorf("Expected %v for uint, got %v", ^uint(0), val)
	}

	// uint8
	if val := umath.MaxValue[uint8](); val != uint8(math.MaxUint8) {
		t.Errorf("Expected %v for uint8, got %v", uint8(math.MaxUint8), val)
	}

	// uint16
	if val := umath.MaxValue[uint16](); val != uint16(math.MaxUint16) {
		t.Errorf("Expected %v for uint16, got %v", uint16(math.MaxUint16), val)
	}

	// uint32
	if val := umath.MaxValue[uint32](); val != uint32(math.MaxUint32) {
		t.Errorf("Expected %v for uint32, got %v", uint32(math.MaxUint32), val)
	}

	// uint64
	if val := umath.MaxValue[uint64](); val != ^uint64(0) {
		t.Errorf("Expected %v for uint64, got %v", ^uint64(0), val)
	}
}
