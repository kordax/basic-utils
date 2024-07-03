/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uarray_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/kordax/basic-utils/uarray"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		index := uarray.IndexOfUint32(slice, test.value)
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
		index := uarray.IndexOfUint32(slice, test.value)
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
		index, value := uarray.ContainsPredicate(slice, func(v *int) bool {
			return *v == test.expectedValue
		})
		if index != test.expectedIndex || *value != test.expectedValue {
			t.Errorf("Expected index %d and value %d, but got index %d and value %d",
				test.expectedIndex, test.expectedValue, index, *value)
		}
	}

	// Test case for predicate that returns false
	index, value := uarray.ContainsPredicate(slice, func(v *int) bool {
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
		index, value := uarray.ContainsStruct[int, MyStruct](test.element, slice)
		if index != test.expectedIndex || !test.expectedValue.Equals(*value) {
			t.Errorf("Expected index %d and value %v, but got index %d and value %v",
				test.expectedIndex, test.expectedValue, index, value)
		}
	}

	// Test case for non-existing struct element
	index, value := uarray.ContainsStruct[int, MyStruct](MyStruct{ID: 4, Name: "Alice"}, slice)
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
		index := uarray.Contains(test.element, slice)
		if index != test.expectedIndex {
			t.Errorf("Expected index %d, but got %d for element %s", test.expectedIndex, index, test.element)
		}
	}

	// Test case for non-existing element
	index := uarray.Contains("grape", slice)
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
		index := uarray.ContainsAny(slice, test.elements...)
		if index == -1 {
			t.Errorf("Expected index element at index %d to be found", test.expectedIndex)
		}
	}

	// Test case for non-existing elements
	index := uarray.ContainsAny(slice, 6, 7, 8)
	if index != -1 {
		t.Errorf("Expected index -1, but got %d for non-existing elements", index)
	}
}

func TestAllMatch(t *testing.T) {
	slice := []string{"apple", "banana", "orange", "grape", "watermelon"}

	// Test case for matching predicate
	predicate := func(v *string) bool {
		return len(*v) > 3
	}
	match := uarray.AllMatch(slice, predicate)
	if !match {
		t.Errorf("Expected AllMatch to return true, but got false")
	}

	// Test case for non-matching predicate
	predicate = func(v *string) bool {
		return len(*v) > 5
	}
	match = uarray.AllMatch(slice, predicate)
	if match {
		t.Errorf("Expected AllMatch to return false, but got true")
	}

	// Additional test case for empty slice
	emptySlice := []string{}
	match = uarray.AllMatch(emptySlice, predicate)
	if !match {
		t.Errorf("Expected AllMatch to return true for empty slice, but got false")
	}

	// Test case with different predicate
	predicate = func(v *string) bool {
		return *v != "banana"
	}
	match = uarray.AllMatch(slice, predicate)
	if match {
		t.Errorf("Expected AllMatch to return false, but got true")
	}
}

func TestAnyMatch(t *testing.T) {
	slice := []string{"apple", "banana", "orange"}

	// Test case for matching predicate
	predicate := func(v *string) bool {
		return len(*v) > 5
	}
	match := uarray.AnyMatch(slice, predicate)
	if !match {
		t.Errorf("Expected AnyMatch to return true, but got false")
	}

	// Test case for non-matching predicate
	predicate = func(v *string) bool {
		return len(*v) > 10
	}
	match = uarray.AnyMatch(slice, predicate)
	if match {
		t.Errorf("Expected AnyMatch to return false, but got true")
	}
}

func TestFilter(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	filtered := uarray.Filter(values, func(v *int) bool {
		return *v%2 == 0
	})
	if !reflect.DeepEqual(filtered, []int{2, 4}) {
		t.Error("Filter function failed")
	}
}

func TestFilterOut(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	filteredOut := uarray.FilterOut(values, func(v *int) bool {
		return *v%2 == 0 // filter out even numbers
	})
	expected := []int{1, 3, 5} // odd numbers should remain

	if !reflect.DeepEqual(filteredOut, expected) {
		t.Errorf("FilterOut function failed, expected %v, got %v", expected, filteredOut)
	}
}

func TestFilterAll(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	matching, nonMatching := uarray.FilterAll(values, func(v *int) bool {
		return *v%2 == 0
	})
	if !reflect.DeepEqual(matching, []int{2, 4}) || !reflect.DeepEqual(nonMatching, []int{1, 3, 5}) {
		t.Error("FilterAll function failed")
	}
}

func TestFilterBySet(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	filter := []int{2, 4}
	filtered := uarray.FilterBySet(values, filter)
	if !reflect.DeepEqual(filtered, []int{2, 4}) {
		t.Error("FilterBySet function failed")
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		target   int
		expected *int
	}{
		{
			name:     "finds element",
			values:   []int{1, 2, 3, 4, 5},
			target:   2,
			expected: &[]int{2}[0],
		},
		{
			name:     "element not found",
			values:   []int{1, 2, 3, 4, 5},
			target:   6,
			expected: nil,
		},
		{
			name:     "empty slice",
			values:   []int{},
			target:   1,
			expected: nil,
		},
		{
			name:     "single element found",
			values:   []int{1},
			target:   1,
			expected: &[]int{1}[0],
		},
		{
			name:     "single element not found",
			values:   []int{2},
			target:   1,
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			found := uarray.Find(test.values, func(v *int) bool {
				return *v == test.target
			})
			assert.Equal(t, test.expected, found, "Find returned an unexpected result")
		})
	}
}

func TestSortFind(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		less     func(a, b int) bool
		target   int
		expected *int
	}{
		{
			name:   "sorted find element",
			values: []int{5, 3, 1, 4, 2},
			less: func(a, b int) bool {
				return a < b
			},
			target:   3,
			expected: &[]int{3}[0],
		},
		{
			name:   "element not found",
			values: []int{5, 3, 1, 4, 2},
			less: func(a, b int) bool {
				return a < b
			},
			target:   6,
			expected: nil,
		},
		{
			name:   "empty slice",
			values: []int{},
			less: func(a, b int) bool {
				return a < b
			},
			target:   1,
			expected: nil,
		},
		{
			name:   "single element found",
			values: []int{1},
			less: func(a, b int) bool {
				return a < b
			},
			target:   1,
			expected: &[]int{1}[0],
		},
		{
			name:   "single element not found",
			values: []int{2},
			less: func(a, b int) bool {
				return a < b
			},
			target:   1,
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			found := uarray.SortFind(test.values, test.less, func(v *int) bool {
				return *v == test.target
			})
			assert.Equal(t, test.expected, found, "SortFind returned the wrong value")
		})
	}
}

func TestMapAggr(t *testing.T) {
	values := []int{1, 2, 3}
	result := uarray.MapAggr(values, func(v *int) []string {
		return []string{fmt.Sprintf("%d-item", *v)}
	})
	expected := []string{"1-item", "2-item", "3-item"}
	if !reflect.DeepEqual(result, expected) {
		t.Error("MapAggr function failed")
	}
}

func TestMap(t *testing.T) {
	values := []int{1, 2, 3}
	result := uarray.Map(values, func(v *int) string {
		return fmt.Sprintf("%d-item", *v)
	})
	expected := []string{"1-item", "2-item", "3-item"}
	if !reflect.DeepEqual(result, expected) {
		t.Error("Map function failed")
	}
}

func TestFlatMap(t *testing.T) {
	values := [][]int{{1, 2}, {3, 4}}
	result := uarray.FlatMap(values, func(v *int) string {
		return fmt.Sprintf("%d-item", *v)
	})
	expected := []string{"1-item", "2-item", "3-item", "4-item"}
	if !reflect.DeepEqual(result, expected) {
		t.Error("FlatMap function failed")
	}
}

func TestFlat(t *testing.T) {
	values := [][]int{{1, 2}, {3, 4}}
	result := uarray.Flat(values)
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Error("Flat function failed")
	}
}

func TestToMap(t *testing.T) {
	values := []string{"apple", "banana"}
	result := uarray.ToMap(values, func(v *string) (int, string) {
		return len(*v), *v
	})
	expected := map[int]string{5: "apple", 6: "banana"}
	if !reflect.DeepEqual(result, expected) {
		t.Error("ToMap function failed")
	}
}

func TestToMultiMap(t *testing.T) {
	values := []string{"apple", "banana", "cherry"}
	result := uarray.ToMultiMap(values, func(v *string) (int, string) {
		return len(*v), *v
	})
	expected := map[int][]string{5: {"apple"}, 6: {"banana", "cherry"}}
	if !reflect.DeepEqual(result, expected) {
		t.Error("ToMultiMap function failed")
	}
}

func TestUniq(t *testing.T) {
	values := []int{1, 2, 2, 3, 3, 3}
	unique := uarray.Uniq(values, func(v *int) int {
		return *v
	})
	if !reflect.DeepEqual(unique, []int{1, 2, 3}) {
		t.Error("Uniq function failed")
	}
}

func TestGroupBy(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	grouped := uarray.GroupBy(values, func(v *int) bool {
		return (*v)%2 == 0
	}, func(v1, v2 *int) int {
		return *v1 + *v2
	})
	sort.Slice(grouped, func(i, j int) bool {
		return grouped[i] < grouped[j]
	})
	require.Equal(t, []int{6, 9}, grouped)
}

func TestGroupToMapBy(t *testing.T) {
	values := []string{"apple", "banana", "cherry"}
	result := uarray.GroupToMapBy(values, func(v *string) int {
		return len(*v)
	})
	expected := map[int][]string{5: {"apple"}, 6: {"banana", "cherry"}}
	if !reflect.DeepEqual(result, expected) {
		t.Error("GroupToMapBy function failed")
	}
}

func TestCopyWithoutIndex(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	cpy := uarray.CopyWithoutIndex(src, 2)
	if !reflect.DeepEqual(cpy, []int{1, 2, 4, 5}) {
		t.Error("CopyWithoutIndex function failed")
	}
}

func TestCollectAsMap(t *testing.T) {
	values := []string{"apple", "banana"}
	result := uarray.CollectAsMap(values, func(v *string) int {
		return len(*v)
	}, func(v string) string {
		return v
	})
	expected := map[int]string{5: "apple", 6: "banana"}
	if !reflect.DeepEqual(result, expected) {
		t.Error("CollectAsMap function failed")
	}
}

func TestEqualsWithOrder(t *testing.T) {
	left := []int{1, 2, 3}
	right := []int{1, 2, 3}
	if !uarray.EqualsWithOrder(left, right) {
		t.Error("EqualsWithOrder function failed")
	}
}

func TestEqualValues(t *testing.T) {
	left := []int{3, 1, 2}
	right := []int{1, 2, 3}
	if !uarray.EqualValues(left, right) {
		t.Error("EqualValues function failed")
	}
}

func TestMerge(t *testing.T) {
	t1 := []int{1, 2, 3}
	t2 := []int{3, 4, 5}
	merged := uarray.Merge(t1, t2, func(t1 *int) int {
		return *t1
	})
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(merged, expected) {
		t.Error("Merge function failed")
	}
}

func TestRange(t *testing.T) {
	// Test the Range function
	expected := []int{1, 2, 3, 4}
	result := uarray.Range(1, 5)
	assert.Equal(t, expected, result)
}

func TestRangeWithStep(t *testing.T) {
	expected := []int{1, 3, 5, 7, 9}
	result := uarray.RangeWithStep(1, 9, 2)
	assert.Equal(t, expected, result)
}

func TestRangeWithHugeStep(t *testing.T) {
	expected := []int{1, 101}
	result := uarray.RangeWithStep(1, 9, 100)
	assert.Equal(t, expected, result)
}

func TestRangeWithStep_ZeroStep(t *testing.T) {
	assert.Panics(t, func() {
		_ = uarray.RangeWithStep(1, 5, 0)
	}, "RangeWithStep should panic with zero step")
}

func TestRangeWithStep_NegativeStep(t *testing.T) {
	assert.Panics(t, func() {
		_ = uarray.RangeWithStep(1, 5, -1)
	}, "RangeWithStep should panic with negative step")
}

func TestRangeWithStep_UnalignedRange(t *testing.T) {
	expected := []int{1, 3, 5, 7, 9, 11}
	result := uarray.RangeWithStep(1, 10, 2)
	assert.Equal(t, expected, result)
}
