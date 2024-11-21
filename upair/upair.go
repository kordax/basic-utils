/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package upair

import "github.com/kordax/basic-utils/uconst"

type Pair[L, R any] struct {
	Left  L
	Right R
}

func NewPair[L, R any](l L, r R) *Pair[L, R] {
	return &Pair[L, R]{Left: l, Right: r}
}

func (p Pair[L, R]) GetLeft() L {
	return p.Left
}

func (p Pair[L, R]) GetRight() R {
	return p.Right
}

// CPair is the same struct as Pair, but forces comparable constraints to support uconst.Comparable contract.
type CPair[L, R comparable] struct {
	Pair[L, R]
}

func NewCPair[L, R comparable](l L, r R) *CPair[L, R] {
	return &CPair[L, R]{
		Pair: Pair[L, R]{Left: l, Right: r},
	}
}

func (p CPair[L, R]) Equals(other uconst.Comparable) bool {
	if o, ok := other.(CPair[L, R]); ok {
		return p.Left == o.Left &&
			p.Right == o.Right
	}

	return false
}
