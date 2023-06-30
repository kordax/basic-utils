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

// Add more tests for other functions...
