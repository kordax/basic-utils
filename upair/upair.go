/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package upair

import "github.com/kordax/basic-utils/v4/uconst"

// Pair stores two related values.
type Pair[L, R any] struct {
	Left  L
	Right R
}

// Of creates a pair value without pointer allocation.
func Of[L, R any](l L, r R) Pair[L, R] {
	return Pair[L, R]{Left: l, Right: r}
}

// NewPair creates a pair pointer.
func NewPair[L, R any](l L, r R) *Pair[L, R] {
	pair := Of(l, r)
	return &pair
}

// GetLeft returns the left value.
func (p Pair[L, R]) GetLeft() L {
	return p.Left
}

// GetRight returns the right value.
func (p Pair[L, R]) GetRight() R {
	return p.Right
}

// Values returns both pair values.
func (p Pair[L, R]) Values() (L, R) {
	return p.Left, p.Right
}

// Swap returns a pair with left and right values swapped.
func (p Pair[L, R]) Swap() Pair[R, L] {
	return Of(p.Right, p.Left)
}

// Map transforms both pair values.
func (p Pair[L, R]) Map[NL, NR any](left func(L) NL, right func(R) NR) Pair[NL, NR] {
	return Of(left(p.Left), right(p.Right))
}

// MapLeft maps the left value and keeps the right value unchanged.
func (p Pair[L, R]) MapLeft[NL any](mapper func(L) NL) Pair[NL, R] {
	return Of(mapper(p.Left), p.Right)
}

// MapRight maps the right value and keeps the left value unchanged.
func (p Pair[L, R]) MapRight[NR any](mapper func(R) NR) Pair[L, NR] {
	return Of(p.Left, mapper(p.Right))
}

// MapLeft maps the left value and keeps the right value unchanged.
// Deprecated: use p.MapLeft(mapper).
func MapLeft[L, R, NL any](p Pair[L, R], mapper func(L) NL) Pair[NL, R] {
	return p.MapLeft(mapper)
}

// MapRight maps the right value and keeps the left value unchanged.
// Deprecated: use p.MapRight(mapper).
func MapRight[L, R, NR any](p Pair[L, R], mapper func(R) NR) Pair[L, NR] {
	return p.MapRight(mapper)
}

// CPair is the same struct as Pair, but forces comparable constraints to support uconst.Comparable contract.
type CPair[L, R comparable] struct {
	Pair[L, R]
}

// COf creates a comparable pair value.
func COf[L, R comparable](l L, r R) CPair[L, R] {
	return CPair[L, R]{
		Pair: Of(l, r),
	}
}

// NewCPair creates a comparable pair pointer.
func NewCPair[L, R comparable](l L, r R) *CPair[L, R] {
	pair := COf(l, r)
	return &pair
}

// Equals checks whether other contains the same comparable pair values.
func (p CPair[L, R]) Equals(other uconst.Comparable) bool {
	switch o := other.(type) {
	case CPair[L, R]:
		return p.Left == o.Left && p.Right == o.Right
	case *CPair[L, R]:
		if o == nil {
			return false
		}
		return p.Left == o.Left && p.Right == o.Right
	default:
		return false
	}
}
