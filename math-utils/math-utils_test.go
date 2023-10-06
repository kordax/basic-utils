/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package mathutils_test

import (
	"math"
	"testing"

	mathutils "github.com/kordax/basic-utils/math-utils"
)

func TestClosestMatch(t *testing.T) {
	slice := []uint32{1, 2, 3, 4, 5}

	// Test case 1: Closest match exists in the slice
	toMatch1 := uint32(3)
	expectedIndex1 := 2
	expectedValue1 := uint32(3)
	index1, value1 := mathutils.ClosestMatch(toMatch1, slice)
	if index1 != expectedIndex1 || value1 != expectedValue1 {
		t.Errorf("Test case 1 failed: ClosestMatch returned (%d, %d), expected (%d, %d)", index1, value1, expectedIndex1, expectedValue1)
	}

	// Test case 2: Closest match does not exist in the slice (value less than the minimum)
	toMatch2 := uint32(0)
	expectedIndex2 := 0
	expectedValue2 := uint32(1)
	index2, value2 := mathutils.ClosestMatch(toMatch2, slice)
	if index2 != expectedIndex2 || value2 != expectedValue2 {
		t.Errorf("Test case 2 failed: ClosestMatch returned (%d, %d), expected (%d, %d)", index2, value2, expectedIndex2, expectedValue2)
	}

	// Test case 3: Closest match does not exist in the slice (value greater than the maximum)
	toMatch3 := uint32(6)
	expectedIndex3 := 4
	expectedValue3 := uint32(5)
	index3, value3 := mathutils.ClosestMatch(toMatch3, slice)
	if index3 != expectedIndex3 || value3 != expectedValue3 {
		t.Errorf("Test case 3 failed: ClosestMatch returned (%d, %d), expected (%d, %d)", index3, value3, expectedIndex3, expectedValue3)
	}
}

func TestValOrMin(t *testing.T) {
	if mathutils.ValOrMin(5, 10) != 10 {
		t.Error("Expected minimum value of 10")
	}

	if mathutils.ValOrMin(15, 10) != 15 {
		t.Error("Expected value of 15")
	}
}

// ... Add tests for other functions in a similar fashion ...

func TestAvgInt(t *testing.T) {
	array := []int{1, 2, 3, 4, 5}
	if mathutils.AvgInt(array) != 3 {
		t.Error("Expected average value of 3")
	}
}

func TestMin(t *testing.T) {
	arrayInt := []int{10, 30, 50, 20, 40}
	if mathutils.Min(arrayInt) != 10 {
		t.Error("Expected minimum value of 10 for int slice")
	}

	arrayFloat64 := []float64{10.5, 30.5, 50.5, 20.5, 40.5}
	if mathutils.Min(arrayFloat64) != 10.5 {
		t.Error("Expected minimum value of 10.5 for float64 slice")
	}
}

func TestMax(t *testing.T) {
	arrayInt := []int{10, 30, 50, 20, 40}
	if mathutils.Max(arrayInt) != 50 {
		t.Error("Expected maximum value of 50 for int slice")
	}

	arrayFloat64 := []float64{10.5, 30.5, 50.5, 20.5, 40.5}
	if mathutils.Max(arrayFloat64) != 50.5 {
		t.Error("Expected maximum value of 50.5 for float64 slice")
	}
}

func TestMinMaxInt(t *testing.T) {
	array := []int{10, 30, 50, 20, 40}
	mn, mx := mathutils.MinMaxInt(array)
	if mn != 10 || mx != 50 {
		t.Errorf("Expected min 10 and max 50, got min %d and max %d", mn, mx)
	}
}

func TestMedInt(t *testing.T) {
	array := []int{10, 20, 30, 40, 50}
	if mathutils.Med(array) != 30 {
		t.Error("Expected median value of 30")
	}
}

func TestMinMax(t *testing.T) {
	arrayInt := []int{10, 20, 30, 40, 50}
	minInt, maxInt := mathutils.MinMax(arrayInt)
	if minInt != 10 || maxInt != 50 {
		t.Errorf("Expected min 10 and max 50 for int slice, got min %d and max %d", minInt, maxInt)
	}

	arrayFloat64 := []float64{10.5, 20.5, 30.5, 40.5, 50.5}
	minFloat64, maxFloat64 := mathutils.MinMax(arrayFloat64)
	if minFloat64 != 10.5 || maxFloat64 != 50.5 {
		t.Errorf("Expected min 10.5 and max 50.5 for float64 slice, got min %f and max %f", minFloat64, maxFloat64)
	}
}

func TestMinMaxFromMap(t *testing.T) {
	mapInt := map[string]int{"a": 10, "b": 20, "c": 30, "d": 40, "e": 50}
	minInt, maxInt := mathutils.MinMaxFromMap(mapInt)
	if minInt != 10 || maxInt != 50 {
		t.Errorf("Expected min 10 and max 50 for int map, got min %d and max %d", minInt, maxInt)
	}

	mapFloat64 := map[string]float64{"a": 10.5, "b": 20.5, "c": 30.5, "d": 40.5, "e": 50.5}
	minFloat64, maxFloat64 := mathutils.MinMaxFromMap(mapFloat64)
	if minFloat64 != 10.5 || maxFloat64 != 50.5 {
		t.Errorf("Expected min 10.5 and max 50.5 for float64 map, got min %f and max %f", minFloat64, maxFloat64)
	}
}

func TestRoundWithPrecision(t *testing.T) {
	if mathutils.RoundWithPrecision(10.4567, 2) != 10.46 {
		t.Error("Expected rounded value of 10.46")
	}

	if mathutils.RoundWithPrecision(10.453, 2) != 10.45 {
		t.Error("Expected rounded value of 10.45")
	}
}

func TestMaxValue(t *testing.T) {
	// float32
	if val := mathutils.MaxValue[float32](); val != math.MaxFloat32 {
		t.Errorf("Expected %v for float32, got %v", math.Float32frombits(^uint32(0)>>1), val)
	}

	// float64
	if val := mathutils.MaxValue[float64](); val != math.MaxFloat64 {
		t.Errorf("Expected %v for float64, got %v", math.Float64frombits(^uint64(0)>>1), val)
	}

	// int
	if val := mathutils.MaxValue[int](); val != int(math.MaxInt) {
		t.Errorf("Expected %v for int, got %v", int(math.MaxInt), val)
	}

	// int8
	if val := mathutils.MaxValue[int8](); val != int8(math.MaxInt8) {
		t.Errorf("Expected %v for int8, got %v", int8(math.MaxInt8), val)
	}

	// int16
	if val := mathutils.MaxValue[int16](); val != int16(math.MaxInt16) {
		t.Errorf("Expected %v for int16, got %v", int16(math.MaxInt16), val)
	}

	// int32
	if val := mathutils.MaxValue[int32](); val != int32(math.MaxInt32) {
		t.Errorf("Expected %v for int32, got %v", int32(math.MaxInt32), val)
	}

	// int64
	if val := mathutils.MaxValue[int64](); val != int64(math.MaxInt64) {
		t.Errorf("Expected %v for int64, got %v", int64(math.MaxInt64), val)
	}

	// uint
	if val := mathutils.MaxValue[uint](); val != ^uint(0) {
		t.Errorf("Expected %v for uint, got %v", ^uint(0), val)
	}

	// uint8
	if val := mathutils.MaxValue[uint8](); val != uint8(math.MaxUint8) {
		t.Errorf("Expected %v for uint8, got %v", uint8(math.MaxUint8), val)
	}

	// uint16
	if val := mathutils.MaxValue[uint16](); val != uint16(math.MaxUint16) {
		t.Errorf("Expected %v for uint16, got %v", uint16(math.MaxUint16), val)
	}

	// uint32
	if val := mathutils.MaxValue[uint32](); val != uint32(math.MaxUint32) {
		t.Errorf("Expected %v for uint32, got %v", uint32(math.MaxUint32), val)
	}

	// uint64
	if val := mathutils.MaxValue[uint64](); val != ^uint64(0) {
		t.Errorf("Expected %v for uint64, got %v", ^uint64(0), val)
	}
}
