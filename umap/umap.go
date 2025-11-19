/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package umap

// Contains checks if map contains specified element.
// Returns value if found, nil otherwise.
func Contains[K comparable, T comparable](values map[K]T, e T) *T {
	for _, v := range values {
		if v == e {
			val := v
			return &val
		}
	}

	return nil
}

// ContainsPredicate checks if map contains specified struct element matching a predicate.
// Returns value if found, nil otherwise.
func ContainsPredicate[K comparable, T any](values map[K]T, predicate func(k K, v *T) bool) *T {
	for k, v := range values {
		val := v
		if predicate(k, &val) {
			return &val
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

// EqualsP returns true if maps are equal. Map order is ignored. Values are compared using predicate.
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

// Copy returns a copy of a map.
func Copy[K comparable, T any](values map[K]T) map[K]T {
	r := make(map[K]T, len(values))
	for k, v := range values {
		r[k] = v
	}

	return r
}

func Keys[K comparable, T any](values map[K]T) []K {
	result := make([]K, 0, len(values))
	for k := range values {
		result = append(result, k)
	}

	return result
}

// Merge merges two maps and returns the result. Existing keys from src map will be used.
func Merge[K comparable, T any](src map[K]T, add map[K]T) map[K]T {
	result := make(map[K]T, len(src)+len(add))
	for k, v := range add {
		result[k] = v
	}
	for k, v := range src {
		result[k] = v
	}

	return result
}

func Values[K comparable, T any](values map[K]T) []T {
	result := make([]T, 0, len(values))
	for _, v := range values {
		result = append(result, v)
	}

	return result
}

// IfPresent executes action(value) if key exists in the map.
func IfPresent[K comparable, V any](values map[K]V, key K, action func(value V)) {
	if v, ok := values[key]; ok {
		action(v)
	}
}

// GetOrDef returns m[key] if present, otherwise def.
func GetOrDef[K comparable, V any](values map[K]V, key K, def V) V {
	if v, ok := values[key]; ok {
		return v
	}

	return def
}
