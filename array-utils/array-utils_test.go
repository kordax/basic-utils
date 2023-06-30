/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package arrayutils_test

import (
	"testing"

	arrayutils "github.com/kordax/basic-utils/array-utils"
)

type MyStruct struct {
	ID   int
	Name string
}

func (s MyStruct) GetIndex() int {
	return s.ID
}

func (s MyStruct) Equals(other any) bool {
	if o, ok := other.(MyStruct); ok {
		return s.ID == o.ID && s.Name == o.Name
	}
	return false
}

func TestIndexOfUint32(t *testing.T) {
	slice := []uint32{1, 2, 3, 4, 5}

	// Test cases for existing values
	existingValueTests := []struct {
		value    uint32
		expected int
	}{
		{1, 0},
		{3, 2},
		{5, 4},
	}
	for _, test := range existingValueTests {
		index := arrayutils.IndexOfUint32(slice, test.value)
		if index != test.expected {
			t.Errorf("Expected index %d, but got %d for value %d", test.expected, index, test.value)
		}
	}

	// Test cases for non-existing values
	nonExistingValueTests := []struct {
		value    uint32
		expected int
	}{
		{0, -1},
		{6, -1},
	}
	for _, test := range nonExistingValueTests {
		index := arrayutils.IndexOfUint32(slice, test.value)
		if index != test.expected {
			t.Errorf("Expected index %d, but got %d for value %d", test.expected, index, test.value)
		}
	}
}

func TestContainsPredicate(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	// Test cases for predicate that returns true
	truePredicateTests := []struct {
		expectedIndex int
		expectedValue int
	}{
		{0, 1},
		{2, 3},
		{4, 5},
	}
	for _, test := range truePredicateTests {
		index, value := arrayutils.ContainsPredicate(slice, func(v *int) bool {
			return *v == test.expectedValue
		})
		if index != test.expectedIndex || *value != test.expectedValue {
			t.Errorf("Expected index %d and value %d, but got index %d and value %d",
				test.expectedIndex, test.expectedValue, index, *value)
		}
	}

	// Test case for predicate that returns false
	index, value := arrayutils.ContainsPredicate(slice, func(v *int) bool {
		return *v == 0
	})
	if index != -1 || value != nil {
		t.Errorf("Expected index -1 and nil value, but got index %d and value %v", index, value)
	}
}

func TestContainsStruct(t *testing.T) {
	slice := []MyStruct{
		{ID: 1, Name: "John"},
		{ID: 2, Name: "Jane"},
		{ID: 3, Name: "Bob"},
	}

	// Test cases for existing struct elements
	existingElementTests := []struct {
		element       MyStruct
		expectedIndex int
		expectedValue *MyStruct
	}{
		{MyStruct{ID: 2, Name: "Jane"}, 1, &slice[1]},
		{MyStruct{ID: 3, Name: "Bob"}, 2, &slice[2]},
	}
	for _, test := range existingElementTests {
		index, value := arrayutils.ContainsStruct[int, MyStruct](test.element, slice)
		if index != test.expectedIndex || !test.expectedValue.Equals(*value) {
			t.Errorf("Expected index %d and value %v, but got index %d and value %v",
				test.expectedIndex, test.expectedValue, index, value)
		}
	}

	// Test case for non-existing struct element
	index, value := arrayutils.ContainsStruct[int, MyStruct](MyStruct{ID: 4, Name: "Alice"}, slice)
	if index != -1 || value != nil {
		t.Errorf("Expected index -1 and nil value, but got index %d and value %v", index, value)
	}
}

func TestContains(t *testing.T) {
	slice := []string{"apple", "banana", "orange"}

	// Test cases for existing elements
	existingElementTests := []struct {
		element       string
		expectedIndex int
	}{
		{"apple", 0},
		{"banana", 1},
		{"orange", 2},
	}
	for _, test := range existingElementTests {
		index := arrayutils.Contains(test.element, slice)
		if index != test.expectedIndex {
			t.Errorf("Expected index %d, but got %d for element %s", test.expectedIndex, index, test.element)
		}
	}

	// Test case for non-existing element
	index := arrayutils.Contains("grape", slice)
	if index != -1 {
		t.Errorf("Expected index -1, but got %d for non-existing element", index)
	}
}

func TestContainsAny(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	// Test cases for existing elements
	existingElementTests := []struct {
		elements      []int
		expectedIndex int
	}{
		{[]int{2, 4, 6}, 1},
		{[]int{3, 6, 9}, 2},
		{[]int{5, 8, 10}, 4},
	}
	for _, test := range existingElementTests {
		index := arrayutils.ContainsAny(slice, test.elements...)
		if index == -1 {
			t.Errorf("Expected index element at index %d to be found", test.expectedIndex)
		}
	}

	// Test case for non-existing elements
	index := arrayutils.ContainsAny(slice, 6, 7, 8)
	if index != -1 {
		t.Errorf("Expected index -1, but got %d for non-existing elements", index)
	}
}

func TestAnyMatch(t *testing.T) {
	slice := []string{"apple", "banana", "orange"}

	// Test case for matching predicate
	predicate := func(v *string) bool {
		return len(*v) > 5
	}
	match := arrayutils.AnyMatch(slice, predicate)
	if !match {
		t.Errorf("Expected AnyMatch to return true, but got false")
	}

	// Test case for non-matching predicate
	predicate = func(v *string) bool {
		return len(*v) > 10
	}
	match = arrayutils.AnyMatch(slice, predicate)
	if match {
		t.Errorf("Expected AnyMatch to return false, but got true")
	}
}
