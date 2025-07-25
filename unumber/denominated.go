/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package unumber

import (
	"errors"

	"github.com/kordax/basic-utils/v2/uconst"
)

type denominated interface {
	~int64 | ~uint64
}

// Denominated represents a value that has been scaled (or "denominated")
// by a certain number of decimal places, effectively converting it to an integer type.
//
// This struct is useful for scenarios where floating-point numbers need to be converted
// to integer types to maintain precision, such as when dealing with monetary values.
//
// Type Parameters:
// - `T`: The underlying integer type that stores the denominated value. This type must satisfy the `denominated` constraint, typically `int64` or `uint64`.
//
// Fields:
// - `v`: The denominated value stored as an integer of type `T`. This value represents the original floating-point number scaled by `10^d`.
// - `d`: The denominator, indicating the number of decimal places the original value was shifted by. This is an integer representing the power of 10 used to scale the original value.
//
// Example Usage:
// Suppose you have a floating-point number `12.345` that you want to store as an integer with three decimal places of precision:
//
// ```go
// var denom Denominated[int64]
// denom.v = 12345  // The original value (12.345) multiplied by 10^3
// denom.d = 3      // Indicates the value was shifted by 3 decimal places
// ```
//
// In this example, `denom.v` stores `12345` and `denom.d` is `3`, indicating that the original value was `12.345`.
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
