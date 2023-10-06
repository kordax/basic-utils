/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package mathutils_test

import (
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

func TestAbsDiffUInt32(t *testing.T) {
	// Test case 1: Positive difference
	one := uint32(5)
	two := uint32(10)
	expected1 := uint32(5)
	result1 := mathutils.AbsDiffUInt32(one, two)
	if result1 != expected1 {
		t.Errorf("Test case 1 failed: AbsDiffUInt32 returned %d, expected %d", result1, expected1)
	}

	// Test case 2: Negative difference
	three := uint32(10)
	four := uint32(5)
	expected2 := uint32(5)
	result2 := mathutils.AbsDiffUInt32(three, four)
	if result2 != expected2 {
		t.Errorf("Test case 2 failed: AbsDiffUInt32 returned %d, expected %d", result2, expected2)
	}

	// Test case 3: Zero difference
	five := uint32(10)
	six := uint32(10)
	expected3 := uint32(0)
	result3 := mathutils.AbsDiffUInt32(five, six)
	if result3 != expected3 {
		t.Errorf("Test case 3 failed: AbsDiffUInt32 returned %d, expected %d", result3, expected3)
	}
}

func TestAbsValInt(t *testing.T) {
	if mathutils.AbsValInt(-5) != 5 {
		t.Error("Expected absolute value of 5")
	}

	if mathutils.AbsValInt(5) != 5 {
		t.Error("Expected absolute value of 5")
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

func TestMaxInt(t *testing.T) {
	array := []int{1, 3, 5, 2, 4}
	if mathutils.MaxInt(array) != 5 {
		t.Error("Expected maximum value of 5")
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

func TestMinMaxUInt32(t *testing.T) {
	array := []uint32{10, 30, 50, 20, 40}
	mn, mx := mathutils.MinMaxUInt32(array)
	if mn != 10 || mx != 50 {
		t.Errorf("Expected min 10 and max 50, got min %d and max %d", mn, mx)
	}
}

func TestMinMaxUInt64(t *testing.T) {
	array := []uint64{10, 30, 50, 20, 40}
	mn, mx := mathutils.MinMaxUInt64(array)
	if mn != 10 || mx != 50 {
		t.Errorf("Expected min 10 and max 50, got min %d and max %d", mn, mx)
	}
}

func TestMinMaxFloat64(t *testing.T) {
	array := []float64{10.5, 30.5, 50.5, 20.5, 40.5}
	mn, mx := mathutils.MinMaxFloat64(array)
	if mn != 10.5 || mx != 50.5 {
		t.Errorf("Expected min 10.5 and max 50.5, got min %f and max %f", mn, mx)
	}
}

func TestAvgUInt64(t *testing.T) {
	array := []uint64{10, 20, 30, 40, 50}
	if mathutils.AvgUInt64(array) != 30 {
		t.Error("Expected average value of 30")
	}
}

func TestAvgInt64(t *testing.T) {
	array := []int64{10, 20, 30, 40, 50}
	if mathutils.AvgInt64(array) != 30 {
		t.Error("Expected average value of 30")
	}
}

func TestAvgFloat64(t *testing.T) {
	array := []float64{10.5, 20.5, 30.5, 40.5, 50.5}
	if mathutils.AvgFloat64(array) != 30.5 {
		t.Error("Expected average value of 30.5")
	}
}

func TestMedInt(t *testing.T) {
	array := []int{10, 20, 30, 40, 50}
	if mathutils.MedInt(array) != 30 {
		t.Error("Expected median value of 30")
	}
}

func TestMedUInt64(t *testing.T) {
	array := []uint64{10, 20, 30, 40, 50}
	if mathutils.MedUInt64(array) != 30 {
		t.Error("Expected median value of 30")
	}
}

func TestMedFloat64(t *testing.T) {
	array := []float64{10.5, 20.5, 30.5, 40.5, 50.5}
	if mathutils.MedFloat64(array) != 30.5 {
		t.Error("Expected median value of 30.5")
	}
}

func TestSumInt(t *testing.T) {
	array := []int{10, 20, 30, 40, 50}
	if mathutils.SumInt(array) != 150 {
		t.Error("Expected sum of 150")
	}
}

func TestSumInt64(t *testing.T) {
	array := []int64{10, 20, 30, 40, 50}
	if mathutils.SumInt64(array) != 150 {
		t.Error("Expected sum of 150")
	}
}

func TestSumUInt64(t *testing.T) {
	array := []uint64{10, 20, 30, 40, 50}
	if mathutils.SumUInt64(array) != 150 {
		t.Error("Expected sum of 150")
	}
}

func TestSumFloat64(t *testing.T) {
	array := []float64{10.5, 20.5, 30.5, 40.5, 50.5}
	if mathutils.SumFloat64(array) != 152.5 {
		t.Error("Expected sum of 152.5")
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
