/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package maputils_test

import (
	"testing"

	arrayutils "github.com/kordax/basic-utils/array-utils"
	maputils "github.com/kordax/basic-utils/map-utils"
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
	result := maputils.ContainsPredicate(predicate, m1)
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
	result2 := maputils.ContainsPredicate(predicate2, m2)
	if result2 != nil {
		t.Errorf("Test case 2 failed: Expected nil, but got %v", result2)
	}
}

func TestEquals(t *testing.T) {
	// Test case 1: Maps are equal
	m1 := map[string]int{"a": 1, "b": 2, "c": 3}
	m2 := map[string]int{"a": 1, "b": 2, "c": 3}
	result := maputils.Equals(m1, m2)
	if !result {
		t.Error("Test case 1 failed: Expected true, but got false")
	}

	// Test case 2: Maps have different values
	m3 := map[string]int{"a": 1, "b": 2, "c": 3}
	m4 := map[string]int{"a": 1, "b": 2, "c": 4}
	result2 := maputils.Equals(m3, m4)
	if result2 {
		t.Error("Test case 2 failed: Expected false, but got true")
	}

	// Test case 3: Maps have different lengths
	m5 := map[string]int{"a": 1, "b": 2, "c": 3}
	m6 := map[string]int{"a": 1, "b": 2}
	result3 := maputils.Equals(m5, m6)
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
	result := maputils.EqualsP(m1, m2, equalsFunc)
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
	result2 := maputils.EqualsP(m3, m4, equalsFunc)
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
	result3 := maputils.EqualsP(m5, m6, equalsFunc)
	if result3 {
		t.Error("Test case 3 failed: Expected false, but got true")
	}
}

func TestCopy(t *testing.T) {
	// Test case 1: Copying a non-empty map
	m1 := map[int]string{1: "one", 2: "two", 3: "three"}
	copy := maputils.Copy(m1)
	if !maputils.Equals(m1, copy) {
		t.Error("Test case 1 failed: The copied map is not equal to the original map")
	}

	// Test case 2: Copying an empty map
	m2 := map[int]string{}
	copy2 := maputils.Copy(m2)
	if !maputils.Equals(m2, copy2) {
		t.Error("Test case 2 failed: The copied map is not equal to the original map")
	}
}

func TestKeys(t *testing.T) {
	// Test case 1: Getting keys from a non-empty map
	m1 := map[int]string{1: "one", 2: "two", 3: "three"}
	keys := maputils.Keys(m1)
	expected := []int{1, 2, 3}
	if !arrayutils.EqualValues(keys, expected) {
		t.Errorf("Test case 1 failed: Expected keys %v, but got %v", expected, keys)
	}

	// Test case 2: Getting keys from an empty map
	m2 := map[int]string{}
	keys2 := maputils.Keys(m2)
	if len(keys2) != 0 {
		t.Errorf("Test case 2 failed: Expected empty keys slice, but got %v", keys2)
	}
}

func TestMerge(t *testing.T) {
	// Test case 1: Merging two non-empty maps
	m1 := map[int]string{1: "one", 2: "two"}
	m2 := map[int]string{3: "three", 4: "four"}
	merged := maputils.Merge(m1, m2)
	expected := map[int]string{1: "one", 2: "two", 3: "three", 4: "four"}
	if !maputils.Equals(merged, expected) {
		t.Errorf("Test case 1 failed: Merged map is not equal to the expected map")
	}

	// Test case 2: Merging a non-empty map with an empty map
	m3 := map[int]string{1: "one", 2: "two"}
	m4 := map[int]string{}
	merged2 := maputils.Merge(m3, m4)
	if !maputils.Equals(merged2, m3) {
		t.Errorf("Test case 2 failed: Merged map is not equal to the original map")
	}

	// Test case 3: Merging two empty maps
	m5 := map[int]string{}
	m6 := map[int]string{}
	merged3 := maputils.Merge(m5, m6)
	if !maputils.Equals(merged3, m5) {
		t.Errorf("Test case 3 failed: Merged map is not equal to the original map")
	}
}

func TestValues(t *testing.T) {
	// Test case 1: Getting values from a non-empty map
	m1 := map[int]string{1: "one", 2: "two", 3: "three"}
	values := maputils.Values(m1)
	expected := []string{"one", "two", "three"}
	if !arrayutils.EqualValues(values, expected) {
		t.Errorf("Test case 1 failed: Expected values %v, but got %v", expected, values)
	}

	// Test case 2: Getting values from an empty map
	m2 := map[int]string{}
	values2 := maputils.Values(m2)
	if len(values2) != 0 {
		t.Errorf("Test case 2 failed: Expected empty values slice, but got %v", values2)
	}
}
