/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uarray

import (
	"cmp"
	"slices"
	"sort"
	"strings"

	"git.casinomodule.org/casino27/basic-utils/v2/ucast"
	"git.casinomodule.org/casino27/basic-utils/v2/uconst"
	"golang.org/x/exp/maps"
)

var dummy struct{}

type Pair[L any, R any] struct {
	Left  L
	Right R
}

func NewPair[L any, R any](left L, right R) *Pair[L, R] {
	return &Pair[L, R]{Left: left, Right: right}
}

func IndexOfUint32(slice []uint32, value uint32) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}

	return -1
}

type Indexed[T any] interface {
	GetIndex() T
}

// ContainsPredicate checks if slice contains specified struct element by using a predicate.
// Returns its index and value if found, -1 and nil otherwise.
func ContainsPredicate[T any](values []T, predicate func(v T) bool) (int, *T) {
	for i := range values {
		if predicate(values[i]) {
			return i, &values[i]
		}
	}

	return -1, nil
}

// ContainsStruct checks if slice contains specified struct element.
// Returns its index and value if found, -1 and nil otherwise.
func ContainsStruct[K comparable, V Indexed[K]](values []V, val V) (int, *V) {
	for i := range values {
		if equals(values[i].GetIndex(), val.GetIndex()) {
			return i, &values[i]
		}
	}

	return -1, nil
}

// Contains checks if slice contains specified element.
// Returns its index if found, -1 otherwise.
func Contains[V comparable](values []V, val V) int {
	for i, v := range values {
		if equals(v, val) {
			return i
		}
	}

	return -1
}

// ContainsAny checks if slice contains any element from another slice.
// Returns its index if found, -1 otherwise.
func ContainsAny[V comparable](from []V, values ...V) int {
	for i, v := range values {
		for _, f := range from {
			if equals(v, f) {
				return i
			}
		}
	}

	return -1
}

// AllMatch checks if all slice elements match the predicate.
// Returns true if all elements match the predicate, false otherwise.
func AllMatch[T any](values []T, predicate func(v T) bool) bool {
	ind, _ := ContainsPredicate(values, func(v T) bool {
		return !predicate(v)
	})
	return ind == -1
}

// AnyMatch checks if slice has an element that matches predicate.
// Returns true if there's a match, false otherwise.
func AnyMatch[T any](values []T, predicate func(v T) bool) bool {
	ind, _ := ContainsPredicate(values, predicate)
	return ind != -1
}

// Has checks if slice has an element.
// Returns true if there's a match, false otherwise.
func Has[T comparable](values []T, val T) bool {
	return AnyMatch(values, func(v T) bool {
		return v == val
	})
}

// Filter filters values slice and returns a copy with filtered elements matching a predicate.
func Filter[V any](values []V, filter func(v V) bool) []V {
	if len(values) == 0 {
		return []V{}
	}
	result := make([]V, 0)
	for _, v := range values {
		if filter(v) {
			result = append(result, v)
		}
	}

	return result
}

// FilterAll filters values slice and returns a copy with filtered elements matching a predicate and elements that do not match any filter.
func FilterAll[V any](values []V, filter func(v V) bool) ([]V, []V) {
	if len(values) == 0 {
		return []V{}, []V{}
	}

	result := make([]V, 0)
	nonMatching := make([]V, 0)
	for _, v := range values {
		if filter(v) {
			result = append(result, v)
		} else {
			nonMatching = append(nonMatching, v)
		}
	}

	return result, nonMatching
}

// FilterBySet filters values slice and returns a copy with filtered elements matching values from filter.
func FilterBySet[V comparable](values []V, filter ...V) []V {
	if len(values) == 0 || len(filter) == 0 {
		return []V{}
	}

	filterSet := make(map[V]struct{})
	for _, v := range filter {
		filterSet[v] = struct{}{}
	}
	result := make([]V, 0)
	for _, v := range values {
		if _, found := filterSet[v]; found {
			result = append(result, v)
		}
	}

	return result
}

// FilterOut is a macros to Filter, so it acts like Filter, but filters out values.
// That means that only values not matching the filter will be returned.
func FilterOut[V any](values []V, filter func(v V) bool) []V {
	return Filter(values, func(v V) bool {
		return !filter(v)
	})
}

// FilterOutBySet filters out values slice and returns a copy without filtered elements matching values from filter.
func FilterOutBySet[V comparable](values []V, filter ...V) []V {
	if len(values) == 0 || len(filter) == 0 {
		return values
	}

	filterSet := make(map[V]struct{})
	for _, v := range filter {
		filterSet[v] = struct{}{}
	}
	result := make([]V, 0)
	for _, v := range values {
		if _, found := filterSet[v]; !found {
			result = append(result, v)
		}
	}

	return result
}

// SortFind sorts the given slice using the provided less function and then finds the first match
// using a binary search with the filter function. This approach is efficient for large slices
// and repeated searches, as it leverages the speed of binary search.
// NOTE: please note that the original slice will be modified.
//
// Parameters:
//   - values: the slice of elements to search through.
//   - less: a function that defines the order of elements for sorting.
//   - filter: a function that tests each element to find a match.
//
// Returns:
//   - A pointer to the found element, or nil if no match is found.
func SortFind[V any](values []V, less func(a, b V) bool, filter func(V) bool) *V {
	if len(values) == 0 {
		return nil
	}

	sort.Slice(values, func(i, j int) bool {
		return less(values[i], values[j])
	})

	return FindBinary(values, filter)
}

// Find finds the first match in a slice using simple looping.
func Find[V any](values []V, filter func(v V) bool) *V {
	for i := range values {
		if filter(values[i]) {
			return &values[i]
		}
	}
	return nil
}

// FindIndex returns the index of the first element matching predicate, or -1 when no element matches.
func FindIndex[T any](values []T, predicate func(v T) bool) int {
	for i, v := range values {
		if predicate(v) {
			return i
		}
	}

	return -1
}

// FindBinary finds the first match in a sorted slice using binary search.
// The filter function should implement a comparison suitable for binary search.
func FindBinary[V any](values []V, filter func(v V) bool) *V {
	start, end := 0, len(values)-1

	for start <= end {
		mid := start + (end-start)/2
		if filter(values[mid]) {
			if mid == 0 || !filter(values[mid-1]) {
				return &values[mid]
			}

			end = mid - 1
		} else {
			start = mid + 1
		}
	}

	return nil
}

// MapAggr maps a func to each set of elements and returns an aggregated result.
func MapAggr[V, R any](values []V, aggr func(v V) []R) []R {
	result := make([]R, 0)
	for _, v := range values {
		result = append(result, aggr(v)...)
	}

	return result
}

// Map maps a func and returns a result.
func Map[V, R any](values []V, m func(v V) R) []R {
	result := make([]R, 0)
	for _, v := range values {
		result = append(result, m(v))
	}

	return result
}

// FlatMap applies the Map method and the Flat method consequently.
func FlatMap[V, R any](values [][]V, m func(v V) R) []R {
	flatten := Flat(values)
	result := make([]R, 0)
	for _, v := range flatten {
		result = append(result, m(v))
	}

	return result
}

// Flat flattens the stream (slice).
func Flat[V any](values [][]V) []V {
	result := make([]V, 0)
	for _, v := range values {
		result = append(result, v...)
	}

	return result
}

// Reduce folds values into a single result by applying reducer from left to right.
func Reduce[T, R any](values []T, initial R, reducer func(acc R, v T) R) R {
	result := initial
	for _, v := range values {
		result = reducer(result, v)
	}

	return result
}

// ToMap collects elements of a slice into a map using a collector function.
//
// If the mapping function produces the same key for multiple elements, the resulting
// map will contain only the last value associated with that key, as the map does not
// behave like a multimap. Each key in the returned map corresponds to a single value,
// and any previous value for the same key will be overwritten.
func ToMap[V any, K comparable, R any](values []V, m func(v V) (K, R)) map[K]R {
	result := make(map[K]R)
	for _, v := range values {
		k, nv := m(v)
		result[k] = nv
	}

	return result
}

// ToMultiMap collects a stream using collector func to a multimap.
func ToMultiMap[V any, K comparable, R any](values []V, m func(v V) (K, R)) map[K][]R {
	result := make(map[K][]R)
	for _, v := range values {
		k, vv := m(v)
		result[k] = append(result[k], vv)
	}

	return result
}

// IndexBy indexes values by key. If multiple values produce the same key, the last value wins.
func IndexBy[T any, K comparable](values []T, key func(v T) K) map[K]T {
	result := make(map[K]T, len(values))
	for _, v := range values {
		result[key(v)] = v
	}

	return result
}

// CountBy counts values grouped by key.
func CountBy[T any, K comparable](values []T, key func(v T) K) map[K]int {
	result := make(map[K]int)
	for _, v := range values {
		result[key(v)]++
	}

	return result
}

// Uniq filters unique elements by predicate that returns any comparable value.
func Uniq[V any, F comparable](values []V, getter func(v V) F) []V {
	set := make(map[F]struct{})
	result := make([]V, 0)

	for _, v := range values {
		key := getter(v)
		if _, exists := set[key]; !exists {
			set[key] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

// Unique filters unique elements from a slice.
func Unique[V comparable](values []V, transform ...func(v V) V) []V {
	set := make(map[V]struct{})
	result := make([]V, 0)

	for _, v := range values {
		transformed := v
		for _, t := range transform {
			transformed = t(transformed)
		}

		if _, exists := set[transformed]; !exists {
			set[transformed] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

// Compact returns a copy without zero-value elements.
func Compact[T comparable](values []T) []T {
	var zero T
	result := make([]T, 0, len(values))
	for _, v := range values {
		if v != zero {
			result = append(result, v)
		}
	}

	return result
}

// CompactFunc returns a copy without elements that match empty predicate.
func CompactFunc[T any](values []T, empty func(v T) bool) []T {
	result := make([]T, 0, len(values))
	for _, v := range values {
		if !empty(v) {
			result = append(result, v)
		}
	}

	return result
}

// GroupBy groups and aggregates elements with aggregator method func.
func GroupBy[V any, G comparable](values []V, group func(v V) G, aggregator func(v1, v2 V) V) []V {
	result := make(map[G]V)
	for _, v := range values {
		g := group(v)
		if existing, contains := result[g]; contains {
			result[g] = aggregator(existing, v)
		} else {
			result[g] = v
		}
	}

	return maps.Values(result)
}

// GroupToMapBy groups elements with group method func.
func GroupToMapBy[V any, G comparable](values []V, group func(v V) G) map[G][]V {
	result := make(map[G][]V)
	for _, v := range values {
		g := group(v)
		result[g] = append(result[g], v)
	}

	return result
}

// MapAndGroupToMapBy same as GroupToMapBy, but allows elements to be mapped to a different type.
func MapAndGroupToMapBy[V any, G comparable, R any](values []V, group func(v V) (G, R)) map[G][]R {
	result := make(map[G][]R)
	for _, v := range values {
		g, r := group(v)
		result[g] = append(result[g], r)
	}

	return result
}

// CopyWithoutIndex copies a slice while ignoring an element at specific index.
func CopyWithoutIndex[T any](src []T, index int) []T {
	cpy := make([]T, 0)
	cpy = append(cpy, src[:index]...)

	return append(cpy, src[index+1:]...)
}

// CopyWithoutIndexes copies a slice while ignoring elements at specific indexes. Duplicate values for indexes are ignored.
func CopyWithoutIndexes[T any](src []T, indexes []int) []T {
	indexMap := make(map[int]struct{})
	for _, index := range indexes {
		if index >= 0 && index < len(src) {
			indexMap[index] = dummy
		}
	}

	result := make([]T, 0, len(src)-len(indexMap))
	for i, v := range src {
		if _, remove := indexMap[i]; !remove {
			result = append(result, v)
		}
	}

	return result
}

// CollectAsMap collects corresponding values to a map.
func CollectAsMap[K comparable, V, R any](values []V, key func(v V) K, val func(v V) R) map[K]R {
	result := make(map[K]R)
	for _, v := range values {
		result[key(v)] = val(v)
	}

	return result
}

// Difference returns a copy of left values that are not present in right.
// The order and duplicate values from left are preserved.
func Difference[T comparable](left, right []T) []T {
	if len(left) == 0 {
		return []T{}
	}
	if len(right) == 0 {
		return slices.Clone(left)
	}

	rightSet := make(map[T]struct{}, len(right))
	for _, v := range right {
		rightSet[v] = dummy
	}

	result := make([]T, 0, len(left))
	for _, v := range left {
		if _, exists := rightSet[v]; !exists {
			result = append(result, v)
		}
	}

	return result
}

// Intersect returns a copy of left values that are present in right.
// The order and duplicate values from left are preserved.
func Intersect[T comparable](left, right []T) []T {
	if len(left) == 0 || len(right) == 0 {
		return []T{}
	}

	rightSet := make(map[T]struct{}, len(right))
	for _, v := range right {
		rightSet[v] = dummy
	}

	result := make([]T, 0, len(left))
	for _, v := range left {
		if _, exists := rightSet[v]; exists {
			result = append(result, v)
		}
	}

	return result
}

// EqualsWithOrder compares two slices taking into consideration elements order.
func EqualsWithOrder[T comparable](left []T, right []T) bool {
	if len(left) != len(right) {
		return false
	}

	for i, v := range left {
		if right[i] != v {
			return false
		}
	}

	return true
}

// EqualsCompareWithOrder compares two slices taking into consideration elements order.
func EqualsCompareWithOrder[T any](left []T, right []T, compare func(t1 T, t2 T) bool) bool {
	if len(left) != len(right) {
		return false
	}

	for i, v := range left {
		if !compare(right[i], v) {
			return false
		}
	}

	return true
}

// EqualValues compares values of two slices regardless of elements order.
func EqualValues[T cmp.Ordered](left []T, right []T) bool {
	if len(left) != len(right) {
		return false
	}

	leftSorted := slices.Clone(left)
	rightSorted := slices.Clone(right)

	sort.SliceStable(leftSorted, func(i, j int) bool {
		return leftSorted[i] < leftSorted[j]
	})
	sort.SliceStable(rightSorted, func(i, j int) bool {
		return rightSorted[i] < rightSorted[j]
	})

	for i, v := range leftSorted {
		if rightSorted[i] != v {
			return false
		}
	}

	return true
}

// EqualValuesCompare compares values of two slices regardless of elements order.
func EqualValuesCompare[T any](left []T, right []T, compare func(t1, t2 T) bool, less func(t1, t2 T) bool) bool {
	if len(left) != len(right) {
		return false
	}

	leftSorted := slices.Clone(left)
	rightSorted := slices.Clone(right)

	sort.SliceStable(leftSorted, func(i, j int) bool {
		return less(leftSorted[i], leftSorted[j])
	})
	sort.SliceStable(rightSorted, func(i, j int) bool {
		return less(rightSorted[i], rightSorted[j])
	})

	for i, v := range leftSorted {
		if !compare(v, rightSorted[i]) {
			return false
		}
	}

	return true
}

// Merge merges two slices with t1 elements prioritized against elements of t2.
func Merge[K comparable, T any](t1 []T, t2 []T, key func(t T) K) []T {
	hashes := make(map[K]struct{})
	var result []T
	for _, t := range t1 {
		k := key(t)
		if _, ok := hashes[k]; !ok {
			hashes[k] = struct{}{}
			result = append(result, t)
		}
	}
	for _, t := range t2 {
		k := key(t)
		if _, ok := hashes[k]; !ok {
			hashes[k] = struct{}{}
			result = append(result, t)
		}
	}

	return result
}

// Reverse returns a reversed copy of values.
func Reverse[T any](values []T) []T {
	result := make([]T, len(values))
	for i, v := range values {
		result[len(values)-1-i] = v
	}

	return result
}

// ReverseInPlace reverses values in place.
func ReverseInPlace[T any](values []T) {
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}
}

// Range generates a slice of integers from 'from' to 'to' (exclusive).
// The type T must be an integer type (e.g., int, int64, uint, etc.).
// The returned slice includes 'from', but is exclusive to 'to'.
// Example usage: FromRange(1, 5) returns []int{1, 2, 3, 4}.
func Range[T uconst.Integer](from, to T) []T {
	result := make([]T, to-from)
	for i := from; i < to; i++ {
		result[i-from] = i
	}

	return result
}

// RangeWithStep generates a slice of integers starting from 'from' up to and including 'to' with a specified step.
func RangeWithStep(from, to, step int) []int {
	if step <= 0 {
		panic("RangeWithStep step must be a positive value")
	}

	size := (to-from)/step + 1
	if (to-from)%step != 0 {
		size++ // Adjust size if 'to' is not divisible by 'step'
	}

	result := make([]int, size)
	for i := 0; i < size; i++ {
		result[i] = from + step*i
	}

	return result
}

// BestMatchBy iterates over the slice and selects the element that satisfies the predicate function as the best match.
// The predicate function compares the current best match with the candidate and determines if the candidate is better.
// Returns a pointer to the selected element or nil if the slice is empty.
func BestMatchBy[T any](values []T, predicate func(currentBest, candidate T) bool) *T {
	if len(values) == 0 {
		return nil
	}

	bestIdx := 0
	for i := 1; i < len(values); i++ {
		if predicate(values[bestIdx], values[i]) {
			bestIdx = i
		}
	}

	return &values[bestIdx]
}

// ClampIndex clamps index to the valid range of values. It returns -1 for empty slices.
func ClampIndex[T any](values []T, index int) int {
	if len(values) == 0 {
		return -1
	}
	if index < 0 {
		return 0
	}
	if index >= len(values) {
		return len(values) - 1
	}

	return index
}

// At returns a pointer to the element at index, or nil when index is out of range.
func At[T any](values []T, index int) *T {
	if index < 0 || index >= len(values) {
		return nil
	}

	return &values[index]
}

// AtOr returns the element at index, or fallback when index is out of range.
func AtOr[T any](values []T, index int, fallback T) T {
	if index < 0 || index >= len(values) {
		return fallback
	}

	return values[index]
}

// Split divides a slice into multiple smaller slices (chunks) of a specified size and returns a slice of these chunks.
//
// If chunkSize is less than or equal to zero, the function returns a slice containing the original slice as its only element.
// If the length of the input slice isn't perfectly divisible by chunkSize, the last chunk will contain the remaining elements.
//
// Example:
//
//	input:    []int{1, 2, 3, 4, 5}, chunkSize: 2
//	output:   [][]int{{1, 2}, {3, 4}, {5}}
//
// Type Parameters:
//
//	T - The type of elements in the slice.
//
// Parameters:
//
//	slice - The input slice to be split.
//	chunkSize - The desired size of each chunk.
//
// Returns:
//
//	A two-dimensional slice where each inner slice is a chunk of the original slice.
func Split[T any](slice []T, chunkSize int) [][]T {
	var chunks [][]T
	if chunkSize <= 0 {
		return append(chunks, slice)
	}

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// AsString converts any supported stringable value to a string and joins them with the specified delimiter.
func AsString[T uconst.Stringable](delimiter string, values ...T) string {
	var parts []string
	for _, v := range values {
		var s string
		switch val := any(v).(type) {
		case int:
			s = ucast.IntToString(val)
		case int8:
			s = ucast.Int8ToString(val)
		case int16:
			s = ucast.Int16ToString(val)
		case int32:
			s = ucast.Int32ToString(val)
		case int64:
			s = ucast.Int64ToString(val)
		case uint:
			s = ucast.UintToString(val)
		case uint8:
			s = ucast.Uint8ToString(val)
		case uint16:
			s = ucast.Uint16ToString(val)
		case uint32:
			s = ucast.Uint32ToString(val)
		case uint64:
			s = ucast.Uint64ToString(val)
		case float32:
			s = ucast.Float32ToString(val)
		case float64:
			s = ucast.Float64ToString(val)
		case bool:
			s = ucast.BoolToString(val)
		case string:
			s = val
		}
		parts = append(parts, s)
	}

	return strings.Join(parts, delimiter)
}

func equals[T comparable](t1, t2 T) bool {
	return t1 == t2
}
