/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package opt

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	basic_utils "github.com/kordax/basic-utils"
	refutils "github.com/kordax/basic-utils/ref-utils"
)

// Opt represents a generic container for optional values.
// An Opt can either contain a value of type T or no value at all.
// It provides methods to check the presence of a value and to retrieve it.
// This struct offers a safer way to handle potentially absent values, avoiding nil dereferences.
// Using Opt ensures that the user must explicitly handle both the present and absent cases,
// thus preventing unintentional null pointer errors.
//
// It's similar in principle to "Optional" in other languages like Java's java.util.Optional.
//
// The internal 'v' field is a pointer to a value of type T.
// If 'v' is nil, it means the Opt contains no value.
// Otherwise, 'v' points to the contained value.
type Opt[T any] struct {
	v *T
}

// Present checks if the Opt contains a value.
func (o Opt[T]) Present() bool {
	return o.v != nil
}

// IfPresent invokes the provided function if the Opt contains a value.
func (o Opt[T]) IfPresent(f func(t T)) {
	if o.Present() {
		f(*o.v)
	}
}

// Null creates an Opt with no value.
func Null[T any]() Opt[T] {
	return Opt[T]{v: nil}
}

// Of creates an Opt with a value.
func Of[T any](v T) Opt[T] {
	return Opt[T]{
		v: &v,
	}
}

// OfNullable creates an Opt that may or may not contain a value based on the provided pointer.
func OfNullable[T any](v *T) Opt[T] {
	if v != nil {
		return Opt[T]{
			v: refutils.Ref(*v),
		}
	} else {
		return Opt[T]{
			v: nil,
		}
	}
}

// OfString creates an Opt containing a string, or a null Opt if the string is empty.
func OfString(v string) Opt[string] {
	if v == "" {
		return Null[string]()
	}

	return Opt[string]{
		v: &v,
	}
}

// OfBool creates an Opt containing a boolean, or a null Opt if the boolean is false.
func OfBool(v bool) Opt[bool] {
	if !v {
		return Null[bool]()
	}

	return Opt[bool]{
		v: &v,
	}
}

// OfNumeric creates an Opt containing a numeric value, or a null Opt if the value is 0.
func OfNumeric[T basic_utils.Numeric](v T) Opt[T] {
	if v == 0 {
		return Null[T]()
	}

	return Opt[T]{
		v: &v,
	}
}

// OfCond creates an Opt containing a value if a given condition is met.
func OfCond[T any](v T, cond func(v *T) bool) Opt[T] {
	if cond(&v) {
		return Opt[T]{
			v: &v,
		}
	}

	return Null[T]()
}

// OfUnix creates an Opt containing a time.Time value based on a Unix timestamp.
func OfUnix[T basic_utils.SignedNumeric](v T) Opt[time.Time] {
	return Opt[time.Time]{
		v: refutils.Ref(time.Unix(int64(v), 0)),
	}
}

// OfBuilder creates an Opt by invoking the provided builder function to generate a value.
func OfBuilder[T any](build func() T) Opt[T] {
	v := build()
	return Opt[T]{
		v: &v,
	}
}

// OrElse retrieves the value within the Opt or a provided default if the Opt is null.
func (o Opt[T]) OrElse(v T) T {
	if o.v == nil {
		return v
	} else {
		return *o.v
	}
}

// Get retrieves the value within the Opt as a pointer.
func (o Opt[T]) Get() *T {
	return o.v
}

// Set sets the value within the Opt.
func (o *Opt[T]) Set(v *T) {
	o.v = v
}

// GetAs retrieves the value within the Opt after applying a mapping function.
func (o Opt[T]) GetAs(mapping func(t *T) any) any {
	return mapping(o.v)
}

// UnmarshalJSON implements the json.Unmarshaler interface for the Opt type.
func (o *Opt[T]) UnmarshalJSON(bytes []byte) error {
	var v T
	if (string)(bytes) == "null" {
		o.v = nil
		return nil
	}
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return err
	}
	o.v = &v

	return nil
}

// MarshalJSON implements the json.Marshaler interface for the Opt type.
func (o Opt[T]) MarshalJSON() ([]byte, error) {
	if !o.Present() {
		return json.Marshal(new(T))
	}

	return json.Marshal(o.Get())
}

// Value implements the driver.Valuer interface for the Opt type, converting its value to a SQL value.
func (o Opt[T]) Value() (driver.Value, error) {
	if o.v != nil {
		switch p := any(o.v).(type) {
		case *int:
			return driver.Value(int64(*p)), nil
		case *int8:
			return driver.Value(int64(*p)), nil
		case *int16:
			return driver.Value(int64(*p)), nil
		case *int32:
			return driver.Value(int64(*p)), nil
		case *uint:
			return driver.Value(int64(*p)), nil
		case *uint8:
			return driver.Value(int64(*p)), nil
		case *uint16:
			return driver.Value(int64(*p)), nil
		case *uint32:
			return driver.Value(int64(*p)), nil
		case *uint64:
			return driver.Value(int64(*p)), nil
		case *float32:
			return driver.Value(float64(*p)), nil
		case driver.Valuer:
			return p.Value()
		}
		return driver.Value(*o.v), nil
	}
	return nil, nil
}

// Scan implements the sql.Scanner interface for the Opt type, reading a SQL value into the Opt.
func (o *Opt[T]) Scan(src interface{}) error {
	if src == nil {
		*o = Null[T]()
		return nil
	}

	var v *T
	switch src.(type) {
	case []uint8:
		switch ptr := any(&v).(type) {
		case **string:
			*ptr = refutils.Ref(string(src.([]uint8)))
		case **T:
			err := json.Unmarshal(src.([]byte), &ptr)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("incompatible type for Opt[%s]: %s, failed to retrieve value", reflect.TypeOf(o.v), reflect.TypeOf(src).String())
		}
	case string:
		switch ptr := any(&v).(type) {
		case **string:
			if src != "" {
				*ptr = refutils.Ref(src.(string))
			}
		case **uint:
			val, err := strconv.ParseUint(src.(string), 10, 32)
			*ptr = refutils.Ref(uint(val))
			if err != nil {
				return fmt.Errorf("failed to parse varchar sql value to numeric opt: %s", src)
			}
		case **uint8:
			val, err := strconv.ParseUint(src.(string), 10, 8)
			*ptr = refutils.Ref(uint8(val))
			if err != nil {
				return fmt.Errorf("failed to parse varchar sql value to numeric opt: %s", src)
			}
		case **uint16:
			val, err := strconv.ParseUint(src.(string), 10, 16)
			*ptr = refutils.Ref(uint16(val))
			if err != nil {
				return fmt.Errorf("failed to parse varchar sql value to numeric opt: %s", src)
			}
		case **uint32:
			val, err := strconv.ParseUint(src.(string), 10, 32)
			*ptr = refutils.Ref(uint32(val))
			if err != nil {
				return fmt.Errorf("failed to parse varchar sql value to numeric opt: %s", src)
			}
		case **uint64:
			val, err := strconv.ParseUint(src.(string), 10, 64)
			*ptr = refutils.Ref(val)
			if err != nil {
				return fmt.Errorf("failed to parse varchar sql value to numeric opt: %s", src)
			}
		case **int:
			val, err := strconv.ParseInt(src.(string), 10, 32)
			*ptr = refutils.Ref(int(val))
			if err != nil {
				return fmt.Errorf("failed to parse varchar sql value to numeric opt: %s", src)
			}
		case **int8:
			val, err := strconv.ParseInt(src.(string), 10, 8)
			*ptr = refutils.Ref(int8(val))
			if err != nil {
				return fmt.Errorf("failed to parse varchar sql value to numeric opt: %s", src)
			}
		case **int16:
			val, err := strconv.ParseInt(src.(string), 10, 16)
			*ptr = refutils.Ref(int16(val))
			if err != nil {
				return fmt.Errorf("failed to parse varchar sql value to numeric opt: %s", src)
			}
		case **int32:
			val, err := strconv.ParseInt(src.(string), 10, 32)
			*ptr = refutils.Ref(int32(val))
			if err != nil {
				return fmt.Errorf("failed to parse varchar sql value to numeric opt: %s", src)
			}
		case **int64:
			val, err := strconv.ParseInt(src.(string), 10, 64)
			*ptr = refutils.Ref(val)
			if err != nil {
				return fmt.Errorf("failed to parse varchar sql value to numeric opt: %s", src)
			}
		case **float32:
			val, err := strconv.ParseFloat(src.(string), 10)
			*ptr = refutils.Ref(float32(val))
			if err != nil {
				return fmt.Errorf("failed to parse varchar sql value to float opt: %s", src)
			}
		case **float64:
			val, err := strconv.ParseFloat(src.(string), 10)
			*ptr = refutils.Ref(val)
			if err != nil {
				return fmt.Errorf("failed to parse varchar sql value to float opt: %s", src)
			}
		case **bool:
			val, err := strconv.ParseBool(src.(string))
			*ptr = refutils.Ref(val)
			if err != nil {
				return fmt.Errorf("failed to parse varchar sql value to bool opt: %s", src)
			}
		case **complex64:
			val, err := strconv.ParseComplex(src.(string), 64)
			*ptr = refutils.Ref(complex64(val))
			if err != nil {
				return fmt.Errorf("failed to parse varchar sql value to complex opt: %s", src)
			}
		case **complex128:
			val, err := strconv.ParseComplex(src.(string), 128)
			*ptr = refutils.Ref(val)
			if err != nil {
				return fmt.Errorf("failed to parse varchar sql value to complex opt: %s", src)
			}
		default:
			return fmt.Errorf("incompatible type for Opt[%s]: %s, failed to retrieve value", reflect.TypeOf(o.v), reflect.TypeOf(src).String())
		}
	case int64:
		switch ptr := any(&v).(type) {
		case **string:
			*ptr = refutils.Ref(src.(string))
		case **uint:
			*ptr = refutils.Ref(uint(src.(int64)))
		case **uint8:
			*ptr = refutils.Ref(uint8(src.(int64)))
		case **uint16:
			*ptr = refutils.Ref(uint16(src.(int64)))
		case **uint32:
			*ptr = refutils.Ref(uint32(src.(int64)))
		case **uint64:
			*ptr = refutils.Ref(uint64(src.(int64)))
		case **int:
			*ptr = refutils.Ref(int(src.(int64)))
		case **int8:
			*ptr = refutils.Ref(int8(src.(int64)))
		case **int16:
			*ptr = refutils.Ref(int16(src.(int64)))
		case **int32:
			*ptr = refutils.Ref(int32(src.(int64)))
		case **int64:
			*ptr = refutils.Ref(src.(int64))
		case **bool:
			if src.(int64) >= 1 {
				*ptr = refutils.Ref(true)
			} else {
				*ptr = refutils.Ref(false)
			}
		default:
			return fmt.Errorf("incompatible type for Opt[%s]: %s, failed to retrieve value", reflect.TypeOf(o.v), reflect.TypeOf(src).String())
		}
	case float32:
		switch ptr := any(&v).(type) {
		case **string:
			*ptr = refutils.Ref(strconv.FormatFloat(float64(src.(float32)), 'f', -1, 32))
		case **float32:
			*ptr = refutils.Ref(src.(float32))
		case **float64:
			*ptr = refutils.Ref(float64(src.(float32)))
		default:
			return fmt.Errorf("incompatible type for Opt[%s]: %s, failed to retrieve value", reflect.TypeOf(o.v), reflect.TypeOf(src).String())
		}
	case float64:
		switch ptr := any(&v).(type) {
		case **string:
			*ptr = refutils.Ref(strconv.FormatFloat(src.(float64), 'f', -1, 64))
		case **float32:
			*ptr = refutils.Ref(float32(src.(float64)))
		case **float64:
			*ptr = refutils.Ref(src.(float64))
		default:
			return fmt.Errorf("incompatible type for Opt[%s]: %s, failed to retrieve value", reflect.TypeOf(o.v), reflect.TypeOf(src).String())
		}
	case nil:
		return nil
	case driver.Valuer:
		valSql, err := src.(driver.Valuer).Value()
		if err != nil {
			return fmt.Errorf("incompatible type for Opt[%s]: %s, failed to retrieve value", reflect.TypeOf(o.v), reflect.TypeOf(src).String())
		}
		v = valSql.(*T)
	default:
		return fmt.Errorf("incompatible type for Opt[%s]: %s", reflect.TypeOf(o.v), reflect.TypeOf(src).String())
	}

	*o = OfNullable(v)

	return nil
}
