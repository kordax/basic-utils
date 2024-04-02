/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package unumber

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type ValueType int

const (
	Int ValueType = iota
	Float
	Uint
	BigInt
	BigFloat
)

func (n ValueType) String() string {
	switch n {
	case Int:
		return "Int"
	case Float:
		return "Float"
	case Uint:
		return "Uint"
	case BigInt:
		return "BigInt"
	case BigFloat:
		return "BigFloat"
	default:
		return "Unknown"
	}
}

// Number is a versatile numeric representation that can hold different types
// of numeric values, such as integers, floating-point numbers, unsigned integers,
// as well as arbitrary-precision numbers (big integers and big floats).
//
// The actual type of the numeric value stored is determined by the field 't'
// which is of type ValueType. Depending on this type:
// - For ValueType 'Int', the 'i' field (of type int) stores the value.
// - For ValueType 'Float', the 'f' field (of type float64) stores the value.
// - For ValueType 'Uint', the 'ui' field (of type uint64) stores the value.
// - For ValueType 'BigInt', the 'bi' field (a pointer to big.Int) stores the value.
// - For ValueType 'BigFloat', the 'bf' field (a pointer to big.Float) stores the value.
//
// The appropriate field should be accessed based on the ValueType to get the
// correct numeric value. Methods like I(), F(), Ui(), Bi(), and Bf() are
// provided to retrieve these values safely.
//
// It's also worth noting that the zero value of this struct isn't a valid
// representation and would need initialization before use. The FromString
// function can be used to initialize a Number from its string representation.
type Number struct {
	t  ValueType
	i  int
	f  float64
	ui uint64
	bi *big.Int
	bf *big.Float
}

// NewInt creates a new Number from an int value.
func NewInt(value int) *Number {
	return &Number{
		t: Int,
		i: value,
	}
}

// NewFloat creates a new Number from a float64 value.
func NewFloat(value float64) *Number {
	return &Number{
		t: Float,
		f: value,
	}
}

// NewUint creates a new Number from a uint64 value.
func NewUint(value uint64) *Number {
	return &Number{
		t:  Uint,
		ui: value,
	}
}

// NewBigInt creates a new Number from a *big.Int value.
func NewBigInt(value *big.Int) *Number {
	return &Number{
		t:  BigInt,
		bi: new(big.Int).Set(value), // Create a new copy of the big.Int
	}
}

// NewBigFloat creates a new Number from a *big.Float value.
func NewBigFloat(value *big.Float) *Number {
	return &Number{
		t:  BigFloat,
		bf: new(big.Float).Set(value), // Create a new copy of the big.Float
	}
}

// FromString creates a new Number from a string representation.
// It determines the appropriate type of the number based on the provided string and a flag indicating whether it's a big number.
func FromString(s string, isBig bool) (*Number, error) {
	var num *Number

	if !isBig {
		if strings.ContainsRune(s, '.') {
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				num = new(Number)
				num.t = Float
				num.f = f

				return num, nil
			} else {
				return num, fmt.Errorf("failed to parse Number ```%s``` as float:\n%s", s, err)
			}
		}

		if strings.ContainsRune(s, '-') {
			if i, err := strconv.Atoi(s); err == nil {
				num = new(Number)
				num.t = Int
				num.i = i
				return num, nil
			} else {
				return num, fmt.Errorf("failed to parse Number ```%s``` as int\n%s", s, err)
			}
		}

		if i, err := strconv.ParseUint(s, 10, 64); err == nil {
			num = new(Number)
			num.t = Uint
			num.ui = i

			return num, nil
		} else {
			return num, fmt.Errorf("failed to parse Number ```%s``` as uint\n%s", s, err)
		}
	} else {
		if strings.ContainsRune(s, '.') {
			num = new(Number)
			num.t = BigFloat
			num.bf = new(big.Float)
			num.bf, _ = num.bf.SetString(s)
			if num.bf == nil {
				return nil, fmt.Errorf("failed to parse Number ```%s``` as big float", s)
			}

			return num, nil
		} else {
			num = new(Number)
			num.t = BigInt
			num.bi = new(big.Int)
			num.bi, _ = num.bi.SetString(s, 10)
			if num.bi == nil {
				return nil, fmt.Errorf("failed to parse Number ```%s``` as big int", s)
			}

			return num, nil
		}
	}
}

func (n *Number) T() ValueType {
	return n.t
}

func (n *Number) I() int {
	return n.i
}

func (n *Number) F() float64 {
	return n.f
}

func (n *Number) Ui() uint64 {
	return n.ui
}

func (n *Number) Bi() *big.Int {
	return n.bi
}

func (n *Number) Bf() *big.Float {
	return n.bf
}

func (n *Number) String() string {
	switch n.t {
	case Int:
		return strconv.FormatInt(int64(n.i), 10)
	case Float:
		return strconv.FormatFloat(n.f, 'f', -1, 64)
	case Uint:
		return strconv.FormatUint(n.ui, 10)
	case BigInt:
		if n.bi != nil {
			return n.bi.String()
		}
	case BigFloat:
		if n.bf != nil {
			return n.bf.String()
		}
	}

	return "Unknown"
}
