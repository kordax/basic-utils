/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap

import "iter"

// MultiMap defines an interface for a generic multimap.
// It allows multiple values, including duplicates, to be associated with a single key.
type MultiMap[K, V any] interface {
	// Get retrieves values associated with the key.
	// Returns a slice of values and a boolean indicating whether the key exists.
	// The slice may contain duplicate values.
	Get(key K) (value []V, ok bool)

	// Set replaces all values associated with the key with the provided values.
	// Returns the count of values that were already associated with the key and have been replaced.
	// This method clears any existing values and inserts the new ones, allowing duplicates among the new values.
	Set(key K, value ...V) int

	// Append adds values to the list of values associated with the key.
	// Returns the total count of values that match the appended ones, counting all occurrences including duplicates.
	// This method does not replace existing values but adds new ones, maintaining any existing duplicates.
	Append(key K, value ...V) int

	// Remove deletes values that match the predicate from the key's associated values.
	// Returns the number of deletions made, accounting for each instance of a value that matches the predicate.
	// Each matching value is checked and removed individually, including each duplicate.
	Remove(key K, predicate func(v V) bool) int

	// Clear removes all values associated with the key.
	// Returns true if there were any values associated with the key before clearing; otherwise, returns false.
	// This method effectively deletes all entries for the key, regardless of their count or duplication status.
	Clear(key K) bool

	// Iterator returns an iterator over all items
	Iterator() iter.Seq2[K, []V]
}
