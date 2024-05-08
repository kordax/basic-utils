/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package upair

type Pair[L, R any] struct {
	Left  L
	Right R
}

func NewPair[L, R any](l L, r R) *Pair[L, R] {
	return &Pair[L, R]{Left: l, Right: r}
}

func (p *Pair[L, R]) GetLeft() L {
	return p.Left
}

func (p *Pair[L, R]) GetRight() R {
	return p.Right
}
