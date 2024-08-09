/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package unumber

import (
	"errors"

	"github.com/kordax/basic-utils/uconst"
)

type denominated interface {
	~int64 | ~uint64
}

type Denominated[T denominated] struct {
	v T
	d int
}

func (d *Denominated[T]) Denominator() int {
	return d.d
}

func (d *Denominated[T]) Value() T {
	return d.v
}

func (d *Denominated[T]) IsValid() bool {
	return d.d != 0
}

// AsDenom converts a floating-point value to a Denominated struct by applying the specified denomination.
// The function uses Denominate to convert the value and then returns a Denominated struct.
//
// Parameters:
// - `value`: The floating-point value to be denominated.
// - `denomination`: The number of decimal places to shift the value by. Must be between 0 and 15.
//
// Returns:
// - A pointer to a Denominated struct containing the denominated value and the denomination.
// - An error if the denomination exceeds 15 or if the denomination operation results in overflow or precision loss.
func AsDenom[T denominated, V uconst.Float](value V, denomination int) (*Denominated[T], error) {
	if denomination > 15 {
		return nil, errors.New("denomination cannot be greater than 15")
	}

	denominatedValue, err := Denominate[V, T](value, denomination)
	if err != nil {
		return nil, err
	}

	return &Denominated[T]{
		v: denominatedValue,
		d: denomination,
	}, nil
}
