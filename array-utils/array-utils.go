/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package arrayutils

import (
	"sort"

	maputils "github.com/kordax/basic-utils/map-utils"
	"golang.org/x/exp/constraints"
)

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
func ContainsPredicate[T any](values []T, predicate func(v *T) bool) (int, *T) {
	for i, v := range values {
		if predicate(&v) {
			return i, &v
		}
	}

	return -1, nil
}

// ContainsStruct checks if slice contains specified struct element.
// Returns its index and value if found, -1 and nil otherwise.
func ContainsStruct[K comparable, V Indexed[K]](val V, values []V) (int, *V) {
	for i, v := range values {
		if equals(v.GetIndex(), val.GetIndex()) {
			return i, &v
		}
	}

	return -1, nil
}

// Contains checks if slice contains specified element.
// Returns its index if found, -1 otherwise.
func Contains[V comparable](val V, values []V) int {
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

// AnyMatch checks if slice has an element that matches predicate.
// Returns true if there's a match, -1 otherwise.
func AnyMatch[T any](values []T, predicate func(v *T) bool) bool {
	ind, _ := ContainsPredicate(values, predicate)
	return ind != -1
}

// Filter filters values slice and returns a copy with filtered elements matching a predicate.
// Returns its index if found, -1 otherwise.
func Filter[V any](values []V, filter func(v *V) bool) []V {
	if len(values) == 0 {
		return []V{}
	}
	result := make([]V, 0)
	for _, v := range values {
		if filter(&v) {
			result = append(result, v)
		}
	}

	return result
}

// FilterAll filters values slice and returns a copy with filtered elements matching a predicate and elements that do not match any filter.
// Returns its index if found, -1 otherwise.
func FilterAll[V any](values []V, filter func(v *V) bool) ([]V, []V) {
	if len(values) == 0 {
		return []V{}, []V{}
	}
	result := make([]V, 0)
	nonMatching := make([]V, 0)
	for _, v := range values {
		if filter(&v) {
			result = append(result, v)
		} else {
			nonMatching = append(nonMatching, v)
		}
	}

	return result, nonMatching
}

// FilterBySet filters values slice and returns a copy with filtered elements matching values from filter.
// Returns its index if found, -1 otherwise.
func FilterBySet[V comparable](values []V, filter []V) []V {
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

// Find finds first match in provided slice.
// TODO: Improve performance/sort slice
func Find[V any](values []V, filter func(v *V) bool) *V {
	if len(values) == 0 {
		return nil
	}
	for _, v := range values {
		if filter(&v) {
			return &v
		}
	}

	return nil
}

// MapAggr maps a func to each set of elements and returns an aggregated result.
func MapAggr[V, R any](values []V, aggr func(v *V) []R) []R {
	result := make([]R, 0)
	for _, v := range values {
		result = append(result, aggr(&v)...)
	}

	return result
}

// Map maps a func and returns a result.
func Map[V, R any](values []V, m func(v *V) R) []R {
	result := make([]R, 0)
	for _, v := range values {
		result = append(result, m(&v))
	}

	return result
}

// FlatMap applies the Map method and the Flat method consequently.
func FlatMap[V, R any](values [][]V, m func(v *V) R) []R {
	flatten := Flat(values)
	result := make([]R, 0)
	for _, v := range flatten {
		result = append(result, m(&v))
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

// ToMap collects a stream using collector func to a map.
func ToMap[V any, K comparable, R any](values []V, m func(v *V) (K, R)) map[K]R {
	result := make(map[K]R)
	for _, v := range values {
		k, v := m(&v)
		result[k] = v
	}

	return result
}

// ToMultiMap collects a stream using collector func to a multimap.
func ToMultiMap[V any, K comparable, R any](values []V, m func(v *V) (K, R)) map[K][]R {
	result := make(map[K][]R)
	for _, v := range values {
		k, v := m(&v)
		result[k] = append(result[k], v)
	}

	return result
}

// Uniq filters unique elements by predicate that returns any comparable value
func Uniq[V any, F comparable](values []V, getter func(v *V) F) []V {
	result := make([]V, 0)
	for _, v := range values {
		if !AnyMatch(result, func(i *V) bool { return getter(i) == getter(&v) }) {
			result = append(result, v)
		}
	}

	return result
}

// GroupBy groups and aggregates elements with aggregator method func
func GroupBy[V any, G comparable](values []V, group func(v *V) G, aggregator func(v1, v2 *V) V) []V {
	result := make(map[G]V)
	for _, v := range values {
		g := group(&v)
		if existing, contains := result[g]; contains {
			result[g] = aggregator(&existing, &v)
		} else {
			result[g] = v
		}
	}

	return maputils.Values(result)
}

// GroupToMapBy groups elements with group method func
func GroupToMapBy[V any, G comparable](values []V, group func(v *V) G) map[G][]V {
	result := make(map[G][]V)
	for _, v := range values {
		g := group(&v)
		result[g] = append(result[g], v)
	}

	return result
}

// CopyWithoutIndex copies a slice ignored the element at specific index
func CopyWithoutIndex[T any](src []T, index int) []T {
	cpy := make([]T, 0)
	cpy = append(cpy, src[:index]...)

	return append(cpy, src[index+1:]...)
}

// CollectAsMap collects corresponding values to a map.
func CollectAsMap[K comparable, V, R any](values []V, key func(v *V) K, val func(v V) R) map[K]R {
	result := make(map[K]R)
	for _, v := range values {
		result[key(&v)] = val(v)
	}

	return result
}

// EqualsWithOrder compares two slices taking into consideration elements order
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

// EqualsCompareWithOrder compares two slices taking into consideration elements order
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

// EqualValues compares values of two slices regardless of elements order
func EqualValues[T constraints.Ordered](left []T, right []T) bool {
	if len(left) != len(right) {
		return false
	}

	sort.SliceStable(left, func(i, j int) bool {
		return left[i] < left[j]
	})
	sort.SliceStable(right, func(i, j int) bool {
		return right[i] < right[j]
	})

	for i, v := range left {
		if right[i] != v {
			return false
		}
	}

	return true
}

// EqualValuesCompare compares values of two slices regardless of elements order
func EqualValuesCompare[T any](left []T, right []T, compare func(t1, t2 T) bool, less func(t1, t2 T) bool) bool {
	if len(left) != len(right) {
		return false
	}

	sort.SliceStable(left, func(i, j int) bool {
		return less(left[i], left[j])
	})
	sort.SliceStable(right, func(i, j int) bool {
		return less(right[i], right[j])
	})

	for i, v := range left {
		if !compare(v, right[i]) {
			return false
		}
	}

	return true
}

func MapKeys[K comparable, V any](m map[K]V) []K {
	k := make([]K, len(m))
	i := 0
	for key := range m {
		k[i] = key
		i++
	}

	return k
}

func MapValues[K comparable, V any](m map[K]V) []V {
	v := make([]V, len(m))
	i := 0
	for _, value := range m {
		v[i] = value
		i++
	}

	return v
}

// Merge merges two slices with t1 elements prioritized against elements of t2.
func Merge[K comparable, T any](t1 []T, t2 []T, key func(t1 *T) K) []T {
	var hashes map[K]struct{}
	var result []T
	for _, t := range t1 {
		if _, ok := hashes[key(&t)]; !ok {
			result = append(result, t)
		}
	}
	for _, t := range t2 {
		if _, ok := hashes[key(&t)]; !ok {
			result = append(result, t)
		}
	}

	return result
}

func equals[T comparable](t1, t2 T) bool {
	return t1 == t2
}
