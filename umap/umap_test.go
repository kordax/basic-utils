/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package umap_test

import (
	"testing"

	"github.com/kordax/basic-utils/uarray"
	"github.com/kordax/basic-utils/umap"
)

type MyStruct struct {
	ID   int
	Name string
}

func TestContainsPredicate(t *testing.T) {
	// Test case 1: Map contains element matching the predicate
	m1 := map[int]MyStruct{
		1: {ID: 1, Name: "John"},
		2: {ID: 2, Name: "Jane"},
		3: {ID: 3, Name: "Bob"},
	}
	predicate := func(k int, v *MyStruct) bool {
		return v.Name == "Jane"
	}
	result := umap.ContainsPredicate(predicate, m1)
	expected := &MyStruct{ID: 2, Name: "Jane"}
	if result == nil || *result != *expected {
		t.Errorf("Test case 1 failed: Expected %v, but got %v", expected, result)
	}

	// Test case 2: Map does not contain element matching the predicate
	m2 := map[int]MyStruct{
		1: {ID: 1, Name: "John"},
		2: {ID: 2, Name: "Jane"},
		3: {ID: 3, Name: "Bob"},
	}
	predicate2 := func(k int, v *MyStruct) bool {
		return v.Name == "Alice"
	}
	result2 := umap.ContainsPredicate(predicate2, m2)
	if result2 != nil {
		t.Errorf("Test case 2 failed: Expected nil, but got %v", result2)
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		element  int
		values   map[string]int
		expected *int
	}{
		{
			element:  5,
			values:   map[string]int{"a": 5, "b": 6, "c": 7},
			expected: &[]int{5}[0], // workaround to get a pointer to int literal
		},
		{
			element:  8,
			values:   map[string]int{"a": 5, "b": 6, "c": 7},
			expected: nil,
		},
		{
			element:  6,
			values:   map[string]int{"d": 6},
			expected: &[]int{6}[0],
		},
		{
			element:  9,
			values:   map[string]int{},
			expected: nil,
		},
	}

	for _, test := range tests {
		result := umap.Contains(test.element, test.values)
		if (result == nil && test.expected != nil) || (result != nil && test.expected == nil) || (result != nil && *result != *test.expected) {
			t.Errorf("Expected %v for element %d in map %v, but got %v", test.expected, test.element, test.values, result)
		}
	}
}

func TestEquals(t *testing.T) {
	// Test case 1: Maps are equal
	m1 := map[string]int{"a": 1, "b": 2, "c": 3}
	m2 := map[string]int{"a": 1, "b": 2, "c": 3}
	result := umap.Equals(m1, m2)
	if !result {
		t.Error("Test case 1 failed: Expected true, but got false")
	}

	// Test case 2: Maps have different values
	m3 := map[string]int{"a": 1, "b": 2, "c": 3}
	m4 := map[string]int{"a": 1, "b": 2, "c": 4}
	result2 := umap.Equals(m3, m4)
	if result2 {
		t.Error("Test case 2 failed: Expected false, but got true")
	}

	// Test case 3: Maps have different lengths
	m5 := map[string]int{"a": 1, "b": 2, "c": 3}
	m6 := map[string]int{"a": 1, "b": 2}
	result3 := umap.Equals(m5, m6)
	if result3 {
		t.Error("Test case 3 failed: Expected false, but got true")
	}
}

func TestEqualsP(t *testing.T) {
	// Test case 1: Maps are equal using custom equality function
	m1 := map[string]MyStruct{
		"a": {ID: 1, Name: "John"},
		"b": {ID: 2, Name: "Jane"},
	}
	m2 := map[string]MyStruct{
		"a": {ID: 1, Name: "John"},
		"b": {ID: 2, Name: "Jane"},
	}
	equalsFunc := func(t1, t2 MyStruct) bool {
		return t1.ID == t2.ID && t1.Name == t2.Name
	}
	result := umap.EqualsP(m1, m2, equalsFunc)
	if !result {
		t.Error("Test case 1 failed: Expected true, but got false")
	}

	// Test case 2: Maps have different values using custom equality function
	m3 := map[string]MyStruct{
		"a": {ID: 1, Name: "John"},
		"b": {ID: 2, Name: "Jane"},
	}
	m4 := map[string]MyStruct{
		"a": {ID: 1, Name: "John"},
		"b": {ID: 2, Name: "Alice"},
	}
	result2 := umap.EqualsP(m3, m4, equalsFunc)
	if result2 {
		t.Error("Test case 2 failed: Expected false, but got true")
	}

	// Test case 3: Maps have different lengths using custom equality function
	m5 := map[string]MyStruct{
		"a": {ID: 1, Name: "John"},
		"b": {ID: 2, Name: "Jane"},
	}
	m6 := map[string]MyStruct{
		"a": {ID: 1, Name: "John"},
	}
	result3 := umap.EqualsP(m5, m6, equalsFunc)
	if result3 {
		t.Error("Test case 3 failed: Expected false, but got true")
	}
}

func TestCopy(t *testing.T) {
	// Test case 1: Copying a non-empty map
	m1 := map[int]string{1: "one", 2: "two", 3: "three"}
	copy := umap.Copy(m1)
	if !umap.Equals(m1, copy) {
		t.Error("Test case 1 failed: The copied map is not equal to the original map")
	}

	// Test case 2: Copying an empty map
	m2 := map[int]string{}
	copy2 := umap.Copy(m2)
	if !umap.Equals(m2, copy2) {
		t.Error("Test case 2 failed: The copied map is not equal to the original map")
	}
}

func TestKeys(t *testing.T) {
	// Test case 1: Getting keys from a non-empty map
	m1 := map[int]string{1: "one", 2: "two", 3: "three"}
	keys := umap.Keys(m1)
	expected := []int{1, 2, 3}
	if !uarray.EqualValues(keys, expected) {
		t.Errorf("Test case 1 failed: Expected keys %v, but got %v", expected, keys)
	}

	// Test case 2: Getting keys from an empty map
	m2 := map[int]string{}
	keys2 := umap.Keys(m2)
	if len(keys2) != 0 {
		t.Errorf("Test case 2 failed: Expected empty keys slice, but got %v", keys2)
	}
}

func TestMerge(t *testing.T) {
	// Test case 1: Merging two non-empty maps
	m1 := map[int]string{1: "one", 2: "two"}
	m2 := map[int]string{3: "three", 4: "four"}
	merged := umap.Merge(m1, m2)
	expected := map[int]string{1: "one", 2: "two", 3: "three", 4: "four"}
	if !umap.Equals(merged, expected) {
		t.Errorf("Test case 1 failed: Merged map is not equal to the expected map")
	}

	// Test case 2: Merging a non-empty map with an empty map
	m3 := map[int]string{1: "one", 2: "two"}
	m4 := map[int]string{}
	merged2 := umap.Merge(m3, m4)
	if !umap.Equals(merged2, m3) {
		t.Errorf("Test case 2 failed: Merged map is not equal to the original map")
	}

	// Test case 3: Merging two empty maps
	m5 := map[int]string{}
	m6 := map[int]string{}
	merged3 := umap.Merge(m5, m6)
	if !umap.Equals(merged3, m5) {
		t.Errorf("Test case 3 failed: Merged map is not equal to the original map")
	}
}

func TestValues(t *testing.T) {
	// Test case 1: Getting values from a non-empty map
	m1 := map[int]string{1: "one", 2: "two", 3: "three"}
	values := umap.Values(m1)
	expected := []string{"one", "two", "three"}
	if !uarray.EqualValues(values, expected) {
		t.Errorf("Test case 1 failed: Expected values %v, but got %v", expected, values)
	}

	// Test case 2: Getting values from an empty map
	m2 := map[int]string{}
	values2 := umap.Values(m2)
	if len(values2) != 0 {
		t.Errorf("Test case 2 failed: Expected empty values slice, but got %v", values2)
	}
}
