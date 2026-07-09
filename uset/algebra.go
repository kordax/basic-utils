package uset

// Union returns a new set containing values from both sets.
func Union[T comparable](left, right Set[T]) *HashSet[T] {
	result := NewHashSetWithSize[T](left.Size() + right.Size())
	for _, v := range left.Values() {
		result.Add(v)
	}
	for _, v := range right.Values() {
		result.Add(v)
	}

	return result
}

// Intersect returns a new set containing values present in both sets.
func Intersect[T comparable](left, right Set[T]) *HashSet[T] {
	if left.Size() > right.Size() {
		left, right = right, left
	}

	result := NewHashSet[T]()
	for _, v := range left.Values() {
		if right.Contains(v) {
			result.Add(v)
		}
	}

	return result
}

// Difference returns a new set containing values from left that are not present in right.
func Difference[T comparable](left, right Set[T]) *HashSet[T] {
	result := NewHashSetWithSize[T](left.Size())
	for _, v := range left.Values() {
		if !right.Contains(v) {
			result.Add(v)
		}
	}

	return result
}

// IsSubset returns true when every value from left is present in right.
func IsSubset[T comparable](left, right Set[T]) bool {
	if left.Size() > right.Size() {
		return false
	}

	for _, v := range left.Values() {
		if !right.Contains(v) {
			return false
		}
	}

	return true
}

// IsSuperset returns true when every value from right is present in left.
func IsSuperset[T comparable](left, right Set[T]) bool {
	return IsSubset(right, left)
}
