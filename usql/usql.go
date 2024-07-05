/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package usql

import (
	"database/sql"
	"math"
	"reflect"
	"time"
)

// NullString constructs a sql.NullString from a regular string.
// It sets the Valid field to true if the provided string is not empty.
func NullString(v string) sql.NullString {
	return sql.NullString{
		String: v,
		Valid:  v != "",
	}
}

// NullBool constructs an sql.NullBool from a regular bool.
// The Valid field is always set to true as a non-null bool is being provided.
func NullBool(v bool) sql.NullBool {
	return sql.NullBool{
		Bool:  v,
		Valid: true, // A bool is always valid, even if false, because it's a non-nullable type in Go
	}
}

// NullTime constructs an sql.NullTime from a regular time.Time.
// It sets the Valid field to true if the provided time is not the zero time.
func NullTime(v time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  v,
		Valid: !v.IsZero(),
	}
}

// NullFloat64 constructs an sql.NullFloat64 from a regular float64.
// It sets the Valid field to true if the provided float64 is not NaN (Not-a-Number).
func NullFloat64(v float64) sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: v,
		Valid:   !math.IsNaN(v),
	}
}

func NewNull[T any](v T) sql.Null[T] {
	// reflect.DeepEqual is used to check if the value is the zero value for its type.
	// reflect.Zero generates a zero value for the type of v, and DeepEqual checks
	// if v is equal to this zero value.
	return sql.Null[T]{
		V:     v,
		Valid: !reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface()),
	}
}
