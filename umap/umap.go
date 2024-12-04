/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package umap

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

// IfPresent checks if the specified key exists in the map `m`.
// If the key is present, it executes the provided `action` function with the associated value.
//
// This function helps avoid boilerplate code when working with maps,
// especially when you need to perform an action only if a key is present.
//
// Type Parameters:
//   - K: The type of the keys in the map. It must be a comparable type.
//   - V: The type of the values in the map.
//
// Parameters:
//   - m map[K]V: The map to check for the presence of the key.
//   - key K: The key to look for in the map.
//   - action func(value V): The function to execute if the key is present.
//     It receives the value associated with the key.
//
// Example Usage:
//
//	package main
//
//	import (
//	    "fmt"
//	    "your_module/ucast"
//	)
//
//	func main() {
//	    m := map[string]int{
//	        "apple":  1,
//	        "banana": 2,
//	    }
//
//	    ucast.IfPresent(m, "apple", func(value int) {
//	        fmt.Println("Value:", value)
//	    })
//
//	    // Output:
//	    // Value: 1
//	}
func IfPresent[K comparable, V any](m map[K]V, key K, action func(value V)) {
	if v, ok := m[key]; ok {
		action(v)
	}
}

// GetOrDef retrieves the value associated with the specified key from the map `m`.
// If the key is not present, it returns the provided `defaultValue`.
//
// This function simplifies accessing map values by providing a default value when a key is absent,
// reducing the need for explicit checks in your code.
//
// Type Parameters:
//   - K: The type of the keys in the map. It must be a comparable type.
//   - V: The type of the values in the map.
//
// Parameters:
//   - m map[K]V: The map from which to retrieve the value.
//   - key K: The key whose associated value is to be returned.
//   - defaultValue V: The value to return if the key is not present in the map.
//
// Returns:
//   - V: The value associated with the key, or `defaultValue` if the key is not found.
//
// Example Usage:
//
//	package main
//
//	import (
//	    "fmt"
//	    "your_module/ucast"
//	)
//
//	func main() {
//	    m := map[string]int{
//	        "apple":  1,
//	        "banana": 2,
//	    }
//
//	    value := ucast.GetOrDef(m, "orange", 0)
//	    fmt.Println("Value:", value) // Output: Value: 0
//	}
func GetOrDef[K comparable, V any](m map[K]V, key K, def V) V {
	if v, ok := m[key]; ok {
		return v
	}

	return def
}
