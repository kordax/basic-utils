/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package maputils

// Contains checks if map contains specified element.
// Returns value if found, nil otherwise.
func Contains[K comparable, T comparable](e T, values map[K]T) *T {
	for _, v := range values {
		if v == e {
			return &v
		}
	}

	return nil
}

// ContainsPredicate checks if map contains specified struct element matching a predicate.
// Returns value if found, nil otherwise.
func ContainsPredicate[K comparable, T any](predicate func(k K, v *T) bool, values map[K]T) *T {
	for k, v := range values {
		if predicate(k, &v) {
			return &v
		}
	}

	return nil
}

// Equals returns true if maps are equal. Map order is ignored.
func Equals[K comparable, T comparable](m1 map[K]T, m2 map[K]T) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k1, v1 := range m1 {
		if v2, ok := m2[k1]; !ok || v2 != v1 {
			return false
		}
	}

	return true
}

// EqualsP returns true if maps are equal. Map order is ignored. Values are compared using predicate
func EqualsP[K comparable, T any](m1 map[K]T, m2 map[K]T, equals func(t1 T, t2 T) bool) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k1, v1 := range m1 {
		if v2, ok := m2[k1]; !ok || !equals(v2, v1) {
			return false
		}
	}

	return true
}

// Copy returns a copy of a map
func Copy[K comparable, T any](m1 map[K]T) map[K]T {
	r := make(map[K]T)
	for k, v := range m1 {
		r[k] = v
	}

	return r
}

func Keys[K comparable, T any](m map[K]T) []K {
	result := make([]K, 0)
	for k := range m {
		result = append(result, k)
	}

	return result
}

// Merge merges two maps and returns the result. Existing keys from src map will be used.
func Merge[K comparable, T any](src map[K]T, add map[K]T) map[K]T {
	result := make(map[K]T)
	for k, v := range add {
		result[k] = v
	}
	for k, v := range src {
		result[k] = v
	}

	return result
}

func Values[K comparable, T any](m map[K]T) []T {
	result := make([]T, 0)
	for _, v := range m {
		result = append(result, v)
	}

	return result
}
