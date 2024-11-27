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
	"strings"
	"testing"

	"github.com/kordax/basic-utils/uarray"
	"github.com/kordax/basic-utils/umath"
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
		index, value := uarray.ContainsStruct[int, MyStruct](slice, test.element)
		if index != test.expectedIndex || !test.expectedValue.Equals(*value) {
			t.Errorf("Expected index %d and value %v, but got index %d and value %v",
				test.expectedIndex, test.expectedValue, index, value)
		}
	}

	// Test case for non-existing struct element
	index, value := uarray.ContainsStruct[int, MyStruct](slice, MyStruct{ID: 4, Name: "Alice"})
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
		index := uarray.Contains(slice, test.element)
		if index != test.expectedIndex {
			t.Errorf("Expected index %d, but got %d for element %s", test.expectedIndex, index, test.element)
		}
	}

	// Test case for non-existing element
	index := uarray.Contains(slice, "grape")
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

func TestHas(t *testing.T) {
	tests := []struct {
		name     string
		values   interface{}
		val      interface{}
		expected bool
	}{
		{
			name:     "String slice - value exists",
			values:   []string{"apple", "banana", "orange"},
			val:      "banana",
			expected: true,
		},
		{
			name:     "String slice - value does not exist",
			values:   []string{"apple", "banana", "orange"},
			val:      "grape",
			expected: false,
		},
		{
			name:     "Int slice - value exists",
			values:   []int{1, 2, 3, 4, 5},
			val:      3,
			expected: true,
		},
		{
			name:     "Int slice - value does not exist",
			values:   []int{1, 2, 3, 4, 5},
			val:      6,
			expected: false,
		},
		{
			name:     "Empty slice",
			values:   []string{},
			val:      "apple",
			expected: false,
		},
		{
			name:     "Slice with duplicates",
			values:   []int{1, 2, 2, 3},
			val:      2,
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch v := tt.values.(type) {
			case []string:
				val := tt.val.(string)
				result := uarray.Has(v, val)
				require.NotNil(t, result, "Result should not be nil")
				assert.Equal(t, tt.expected, result, "Has(%v, %v) = %v; want %v", v, val, result, tt.expected)
			case []int:
				val := tt.val.(int)
				result := uarray.Has(v, val)
				require.NotNil(t, result, "Result should not be nil")
				assert.Equal(t, tt.expected, result, "Has(%v, %v) = %v; want %v", v, val, result, tt.expected)
			default:
				t.Fatalf("Unsupported type for values: %T", tt.values)
			}
		})
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
	filtered := uarray.FilterBySet(values, filter...)
	if !reflect.DeepEqual(filtered, []int{2, 4}) {
		t.Error("FilterBySet function failed")
	}
}

func TestFilterOutBySet(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		filter   []int
		expected []int
	}{
		{
			name:     "Filter out some elements",
			values:   []int{1, 2, 3, 4, 5},
			filter:   []int{2, 4},
			expected: []int{1, 3, 5},
		},
		{
			name:     "Filter is empty",
			values:   []int{1, 2, 3},
			filter:   []int{},
			expected: []int{1, 2, 3},
		},
		{
			name:     "Values is empty",
			values:   []int{},
			filter:   []int{1, 2},
			expected: []int{},
		},
		{
			name:     "All elements filtered out",
			values:   []int{1, 2, 3},
			filter:   []int{1, 2, 3},
			expected: []int{},
		},
		{
			name:     "No elements filtered out",
			values:   []int{1, 2, 3},
			filter:   []int{4, 5},
			expected: []int{1, 2, 3},
		},
		{
			name:     "Values with duplicates",
			values:   []int{1, 2, 2, 3, 4, 5},
			filter:   []int{2, 4},
			expected: []int{1, 3, 5},
		},
		{
			name:     "Filter with duplicates",
			values:   []int{1, 2, 3, 4, 5},
			filter:   []int{2, 2, 4, 4},
			expected: []int{1, 3, 5},
		},
		{
			name:     "Both values and filter are empty",
			values:   []int{},
			filter:   []int{},
			expected: []int{},
		},
		{
			name:     "Non-overlapping values and filter",
			values:   []int{6, 7, 8},
			filter:   []int{1, 2, 3},
			expected: []int{6, 7, 8},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := uarray.FilterOutBySet(tt.values, tt.filter...)
			require.NotNil(t, result, "Result should not be nil")
			assert.Equal(t, tt.expected, result, "Filtered result does not match expected output")
		})
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

func TestMapAndGroupToMapBy(t *testing.T) {
	values := []string{"apple", "banana", "cherry", "avocado"}
	expected := map[int][]string{
		len("APPLE"):   {"APPLE"},
		len("BANANA"):  {"BANANA", "CHERRY"},
		len("AVOCADO"): {"AVOCADO"},
	}

	result := uarray.MapAndGroupToMapBy(values, func(v *string) (int, string) {
		return len(*v), strings.ToUpper(*v)
	})
	assert.Equal(t, expected, result, "MapAndGroupToMapBy function failed for strings mapped to uppercase")

	intValues := []int{1, 2, 3, 4, 5}
	expectedInt := map[bool][]int{
		true:  {4, 16},    // even numbers squared
		false: {1, 9, 25}, // odd numbers squared
	}
	intResult := uarray.MapAndGroupToMapBy(intValues, func(v *int) (bool, int) {
		return *v%2 == 0, (*v) * (*v)
	})

	assert.Equal(t, expectedInt, intResult, "MapAndGroupToMapBy function failed for integers squared")
}

func TestCopyWithoutIndex(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	cpy := uarray.CopyWithoutIndex(src, 2)
	if !reflect.DeepEqual(cpy, []int{1, 2, 4, 5}) {
		t.Error("CopyWithoutIndex function failed")
	}
}

func TestCopyWithoutIndexes(t *testing.T) {
	tests := []struct {
		name     string
		src      []int
		indexes  []int
		expected []int
	}{
		{
			name:     "single index",
			src:      []int{1, 2, 3, 4, 5},
			indexes:  []int{2},
			expected: []int{1, 2, 4, 5},
		},
		{
			name:     "multiple indexes",
			src:      []int{1, 2, 3, 4, 5},
			indexes:  []int{1, 3},
			expected: []int{1, 3, 5},
		},
		{
			name:     "indexes out of order",
			src:      []int{1, 2, 3, 4, 5},
			indexes:  []int{3, 1},
			expected: []int{1, 3, 5},
		},
		{
			name:     "duplicate indexes",
			src:      []int{1, 2, 3, 4, 5},
			indexes:  []int{1, 1},
			expected: []int{1, 3, 4, 5},
		},
		{
			name:     "no indexes",
			src:      []int{1, 2, 3, 4, 5},
			indexes:  []int{},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "indexes out of bounds",
			src:      []int{1, 2, 3, 4, 5},
			indexes:  []int{10},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "all indexes",
			src:      []int{1, 2, 3, 4, 5},
			indexes:  []int{0, 1, 2, 3, 4},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpy := uarray.CopyWithoutIndexes(tt.src, tt.indexes)
			require.Equal(t, tt.expected, cpy, "result should match expected slice")
		})
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

func TestBestMatchBy(t *testing.T) {
	t.Run("FindLargestInteger", func(t *testing.T) {
		bestOne := 45
		values := []int{10, 20, bestOne, 25, 15}
		best := uarray.BestMatchBy(values, func(current, candidate *int) bool {
			return *candidate > *current
		})
		assert.NotNil(t, best, "BestMatchBy should return a non-nil result for non-empty slice")
		assert.Equal(t, bestOne, *best, "BestMatchBy should return the largest integer")
	})

	t.Run("FindSmallestInteger", func(t *testing.T) {
		bestOne := 10
		values := []int{25, 20, 45, bestOne, 15}
		best := uarray.BestMatchBy(values, func(current, candidate *int) bool {
			return *candidate < *current
		})
		assert.NotNil(t, best, "BestMatchBy should return a non-nil result for non-empty slice")
		assert.Equal(t, bestOne, *best, "BestMatchBy should return the smallest integer")
	})

	t.Run("FindLongestString", func(t *testing.T) {
		values := []string{"apple", "banana", "cherry", "watermelon"}
		best := uarray.BestMatchBy(values, func(current, candidate *string) bool {
			return len(*candidate) > len(*current)
		})
		assert.NotNil(t, best, "BestMatchBy should return a non-nil result for non-empty slice")
		assert.Equal(t, "watermelon", *best, "BestMatchBy should return the longest string")
	})

	t.Run("FindShortestString", func(t *testing.T) {
		values := []string{"apple", "banana", "cherry", "fig"}
		best := uarray.BestMatchBy(values, func(current, candidate *string) bool {
			return len(*candidate) < len(*current)
		})
		assert.NotNil(t, best, "BestMatchBy should return a non-nil result for non-empty slice")
		assert.Equal(t, "fig", *best, "BestMatchBy should return the shortest string")
	})

	t.Run("FindHighestValueInStruct", func(t *testing.T) {
		type TestStruct struct {
			ID    int
			Value int
		}
		values := []TestStruct{
			{ID: 1, Value: 100},
			{ID: 2, Value: 200},
			{ID: 3, Value: 150},
		}
		best := uarray.BestMatchBy(values, func(current, candidate *TestStruct) bool {
			return candidate.Value > current.Value
		})
		assert.NotNil(t, best, "BestMatchBy should return a non-nil result for non-empty slice")
		assert.Equal(t, 200, best.Value, "BestMatchBy should return the struct with the highest Value field")
	})

	t.Run("EmptySlice", func(t *testing.T) {
		values := []int{}
		best := uarray.BestMatchBy(values, func(current, candidate *int) bool {
			return true // Arbitrary predicate
		})
		assert.Nil(t, best, "BestMatchBy should return nil for an empty slice")
	})

	t.Run("SingleElement", func(t *testing.T) {
		values := []int{42}
		best := uarray.BestMatchBy(values, func(current, candidate *int) bool {
			return true // Always true
		})
		assert.NotNil(t, best, "BestMatchBy should return a non-nil result for a single-element slice")
		assert.Equal(t, 42, *best, "BestMatchBy should return the single element itself")
	})

	t.Run("IdenticalElements", func(t *testing.T) {
		values := []int{5, 5, 5, 5}
		best := uarray.BestMatchBy(values, func(current, candidate *int) bool {
			return *candidate > *current
		})
		assert.NotNil(t, best, "BestMatchBy should return a non-nil result for a slice with identical elements")
		assert.Equal(t, 5, *best, "BestMatchBy should return one of the identical elements")
	})

	t.Run("FindLargestAbsoluteValue", func(t *testing.T) {
		values := []int{-10, -20, 15, 5, -30}
		best := uarray.BestMatchBy(values, func(current, candidate *int) bool {
			return umath.AbsVal(*candidate) > umath.AbsVal(*current)
		})
		assert.NotNil(t, best, "BestMatchBy should return a non-nil result for non-empty slice")
		assert.Equal(t, -30, *best, "BestMatchBy should return the number with the largest absolute value")
	})
}

func TestSplit_EmptySlice(t *testing.T) {
	var slice []int
	chunkSize := 3
	result := uarray.Split(slice, chunkSize)
	assert.Equal(t, 0, len(result), "Expected empty result for empty slice")
}

func TestSplit_NilSlice(t *testing.T) {
	var slice []int = nil
	chunkSize := 3
	result := uarray.Split(slice, chunkSize)
	assert.Equal(t, 0, len(result), "Expected empty result for nil slice")
}

func TestSplit_ChunkSizeZero(t *testing.T) {
	slice := []int{1, 2, 3}
	chunkSize := 0
	result := uarray.Split(slice, chunkSize)
	assert.Equal(t, [][]int{{1, 2, 3}}, result, "Expected original slice when chunkSize is zero")
}

func TestSplit_ChunkSizeNegative(t *testing.T) {
	slice := []int{1, 2, 3}
	chunkSize := -1
	result := uarray.Split(slice, chunkSize)
	assert.Equal(t, [][]int{{1, 2, 3}}, result, "Expected original slice when chunkSize is negative")
}

func TestSplit_ChunkSizeOne(t *testing.T) {
	slice := []int{1, 2, 3}
	chunkSize := 1
	result := uarray.Split(slice, chunkSize)
	expected := [][]int{{1}, {2}, {3}}
	assert.Equal(t, expected, result, "Expected each element in its own chunk when chunkSize is one")
}

func TestSplit_ChunkSizeGreaterThanLength(t *testing.T) {
	slice := []int{1, 2, 3}
	chunkSize := 5
	result := uarray.Split(slice, chunkSize)
	assert.Equal(t, [][]int{{1, 2, 3}}, result, "Expected original slice in one chunk when chunkSize exceeds slice length")
}

func TestSplit_ExactDivision(t *testing.T) {
	slice := []int{1, 2, 3, 4}
	chunkSize := 2
	result := uarray.Split(slice, chunkSize)
	expected := [][]int{{1, 2}, {3, 4}}
	assert.Equal(t, expected, result, "Expected chunks of exact size when slice length is divisible by chunkSize")
}

func TestSplit_NonExactDivision(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	chunkSize := 2
	result := uarray.Split(slice, chunkSize)
	expected := [][]int{{1, 2}, {3, 4}, {5}}
	assert.Equal(t, expected, result, "Expected last chunk to contain remaining elements when slice length isn't divisible by chunkSize")
}

func TestSplit_StringSlice(t *testing.T) {
	slice := []string{"a", "b", "c", "d"}
	chunkSize := 2
	result := uarray.Split(slice, chunkSize)
	expected := [][]string{{"a", "b"}, {"c", "d"}}
	assert.Equal(t, expected, result, "Expected correct chunks for string slice")
}

func TestSplit_StructSlice(t *testing.T) {
	type Item struct{ Value int }
	slice := []Item{{1}, {2}, {3}}
	chunkSize := 2
	result := uarray.Split(slice, chunkSize)
	expected := [][]Item{{Item{1}, Item{2}}, {Item{3}}}
	assert.Equal(t, expected, result, "Expected correct chunks for slice of structs")
}

func TestSplit_LargeChunkSize(t *testing.T) {
	slice := []int{1, 2, 3}
	chunkSize := 1000
	result := uarray.Split(slice, chunkSize)
	assert.Equal(t, [][]int{{1, 2, 3}}, result, "Expected original slice in one chunk when chunkSize is very large")
}

func TestSplit_NilSlice_ChunkSizeZero(t *testing.T) {
	var slice []int = nil
	chunkSize := 0
	result := uarray.Split(slice, chunkSize)
	assert.Equal(t, 1, len(result), "Expected result length of 1 for nil slice and chunkSize zero")
	assert.Nil(t, result[0], "Expected first element to be nil")
}

func TestAsString(t *testing.T) {
	tests := []struct {
		name      string
		delimiter string
		input     interface{}
		expected  string
	}{
		{
			name:      "Empty int slice",
			delimiter: ",",
			input:     []int{},
			expected:  "",
		},
		{
			name:      "Single int",
			delimiter: ",",
			input:     []int{42},
			expected:  "42",
		},
		{
			name:      "Multiple ints",
			delimiter: ",",
			input:     []int{1, 2, 3},
			expected:  "1,2,3",
		},
		{
			name:      "Multiple int8s",
			delimiter: ";",
			input:     []int8{10, 20, 30},
			expected:  "10;20;30",
		},
		{
			name:      "Multiple int16s with different delimiter",
			delimiter: "|",
			input:     []int16{100, 200, 300},
			expected:  "100|200|300",
		},
		{
			name:      "Multiple int32s",
			delimiter: ",",
			input:     []int32{1000, 2000, 3000},
			expected:  "1000,2000,3000",
		},
		{
			name:      "Multiple int64s",
			delimiter: ",",
			input:     []int64{10000, 20000, 30000},
			expected:  "10000,20000,30000",
		},
		{
			name:      "Multiple uint8s",
			delimiter: ",",
			input:     []uint8{255, 128, 64},
			expected:  "255,128,64",
		},
		{
			name:      "Multiple uint16s",
			delimiter: ",",
			input:     []uint16{65535, 32768, 16384},
			expected:  "65535,32768,16384",
		},
		{
			name:      "Multiple uint32s",
			delimiter: ",",
			input:     []uint32{4294967295, 2147483648, 1073741824},
			expected:  "4294967295,2147483648,1073741824",
		},
		{
			name:      "Multiple uint64s",
			delimiter: ",",
			input:     []uint64{18446744073709551615, 9223372036854775808, 4611686018427387904},
			expected:  "18446744073709551615,9223372036854775808,4611686018427387904",
		},
		{
			name:      "Multiple float32s",
			delimiter: ",",
			input:     []float32{3.14, 2.71, 1.61},
			expected:  "3.14,2.71,1.61",
		},
		{
			name:      "Multiple float64s",
			delimiter: ",",
			input:     []float64{6.28, 5.55, 4.44},
			expected:  "6.28,5.55,4.44",
		},
		{
			name:      "Multiple bools",
			delimiter: ",",
			input:     []bool{true, false, true},
			expected:  "true,false,true",
		},
		{
			name:      "Mixed types with different delimiters",
			delimiter: "-",
			input:     []int{1, 2, 3},
			expected:  "1-2-3",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var result string
			switch input := tt.input.(type) {
			case []int:
				result = uarray.AsString(tt.delimiter, input...)
			case []int8:
				result = uarray.AsString(tt.delimiter, input...)
			case []int16:
				result = uarray.AsString(tt.delimiter, input...)
			case []int32:
				result = uarray.AsString(tt.delimiter, input...)
			case []int64:
				result = uarray.AsString(tt.delimiter, input...)
			case []uint8:
				result = uarray.AsString(tt.delimiter, input...)
			case []uint16:
				result = uarray.AsString(tt.delimiter, input...)
			case []uint32:
				result = uarray.AsString(tt.delimiter, input...)
			case []uint64:
				result = uarray.AsString(tt.delimiter, input...)
			case []float32:
				result = uarray.AsString(tt.delimiter, input...)
			case []float64:
				result = uarray.AsString(tt.delimiter, input...)
			case []bool:
				result = uarray.AsString(tt.delimiter, input...)
			default:
				t.Fatalf("Unsupported input type: %T", tt.input)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}
