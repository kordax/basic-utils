/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package mathutils

import (
	"fmt"
	"math"
	"sort"

	basicutils "github.com/kordax/basic-utils"
)

// ClosestMatch finds the closest match in a slice. This implementation doesn't use Round.
// Returns index of a found element and the element. If no element was found then -1 as index will be returned.
func ClosestMatch[T basicutils.Numeric](toMatch T, slice []T) (int, T) {
	diff := MaxValue[T]()
	var value T
	index := -1
	for i, val := range slice {
		d := AbsDiff(val, toMatch)
		if d < diff {
			diff = d
			index = i
			value = val
		}
	}

	return index, value
}

// AbsDiff returns absolute diff result of a numeric type
func AbsDiff[T basicutils.Numeric](one T, two T) T {
	if one < two {
		return two - one
	}

	return one - two
}

// AbsVal returns absolute value of a numeric type
func AbsVal[T basicutils.Numeric](val T) T {
	if val < 0 {
		return -val
	}

	return val
}

// ValOrMin returns value or mn if value is less than mn
func ValOrMin(val int, mn int) int {
	if val < mn {
		return mn
	}

	return val
}

// Min returns mn numeric value found in array
func Min[T basicutils.Numeric](array []T) T {
	if len(array) == 0 {
		return *new(T)
	}

	mn := array[0]
	for _, a := range array {
		if a < mn {
			mn = a
		}
	}

	return mn
}

// Max returns mx numeric value found in array
func Max[T basicutils.Numeric](array []T) T {
	if len(array) == 0 {
		return *new(T)
	}

	mx := array[0]
	for _, a := range array {
		if a > mx {
			mx = a
		}
	}

	return mx
}

// MinMaxInt returns mn and mx value in array correspondingly
func MinMaxInt(array []int) (int, int) {
	if len(array) == 0 {
		return 0, 0
	}

	mn := array[0]
	mx := array[0]
	for _, a := range array {
		if a < mn {
			mn = a
		}
		if a > mx {
			mx = a
		}
	}

	return mn, mx
}

// Avg returns average value from the array
func Avg[T basicutils.Numeric](array []T) T {
	if len(array) == 0 {
		return 0
	}

	sum := Sum(array)

	return sum / T(len(array))
}

// AvgFloat returns average value from the array as float64
func AvgFloat[T basicutils.Numeric](array []T) float64 {
	if len(array) == 0 {
		return 0
	}

	sum := Sum(array)

	return float64(sum) / float64(T(len(array)))
}

func Med[T basicutils.Numeric](array []T) T {
	if len(array) == 0 {
		return 0.0
	}
	sort.Slice(array, func(i, j int) bool { return array[i] < array[j] })

	ln := len(array)
	if ln == 1 {
		return array[0]
	}

	var med T
	if ln%2 != 0 {
		med = array[(ln / 2)]
	} else {
		med = (array[(ln/2)-1] + array[(ln/2)]) / 2
	}

	return med
}

func Sum[T basicutils.Numeric](array []T) T {
	var sum T
	for _, v := range array {
		sum = sum + v
	}

	return sum
}

func MinMax[T basicutils.Numeric](array []T) (T, T) {
	if len(array) == 0 {
		return *new(T), *new(T)
	}

	mn := array[0]
	mx := array[0]
	for _, a := range array {
		if a < mn {
			mn = a
		}
		if a > mx {
			mx = a
		}
	}

	return mn, mx
}

func MinMaxFromMap[K comparable, T basicutils.Numeric](m map[K]T) (T, T) {
	if len(m) == 0 {
		return *new(T), *new(T)
	}

	var mn T
	mx := *new(T)
	first := true
	for _, a := range m {
		if a < mn || first {
			first = false
			mn = a
		}
		if a > mx {
			mx = a
		}
	}

	return mn, mx
}

// RoundWithPrecision rounds a numeric value to a specified number of decimal places.
// This function is generic and can operate on any type that satisfies the basicutils.Numeric constraint,
// which typically includes integer and floating-point types like int, float32, and float64.
//
// Parameters:
//   - value: The numeric value to be rounded. This can be any type that satisfies the basicutils.Numeric constraint.
//   - precision: The number of decimal places to round to. It specifies how many digits should appear after the decimal point.
//     If precision is negative, the function will round off digits to the left of the decimal point.
//
// Returns:
//   - The rounded value of the same type as the input `value`. The rounding is done according to the standard rounding rules,
//     where numbers halfway between two possible outcomes are rounded to the nearest even number (also known as "bankers rounding").
//
// Behavior and Edge Cases:
//   - If precision is 0, the function rounds to the nearest integer.
//   - Positive precision values round to the specified number of digits after the decimal point.
//   - Negative precision values round to the nearest 10^(-precision) place. For example, a precision of -1 rounds to the nearest 10,
//     -2 rounds to the nearest 100, and so forth, effectively reducing the precision of the input number.
//   - For floating-point values, the function handles edge cases like rounding at 0.5, where the standard rounding method is applied.
//   - For integer types, the precision parameter is less meaningful beyond 0, but the function will still operate without error,
//     effectively leaving the integer value unchanged for positive precision values, and rounding to multiples of 10 for negative precision values.
//   - The function may lose precision for very large numbers or for numbers requiring more precision than the type can represent.
//     For example, rounding a float64 to a very high precision can result in a loss of accuracy due to the inherent limitations of floating-point representation.
//   - If the type T does not have enough precision to represent the result accurately, the behavior is dependent on the underlying type's precision and range.
//
// Usage:
// The function is versatile and can be used with various numeric types.
// It is especially useful in scenarios where the precision of numerical computation needs to be controlled explicitly.
func RoundWithPrecision[T basicutils.Numeric](value T, precision int) T {
	ratio := math.Pow(10, float64(precision))

	return T(math.Round(float64(value)*ratio) / ratio)
}

// RoundUp rounds the given numeric value up to the nearest integer, regardless of the fractional part.
// This function is akin to RoundWithPrecision but differs in its rounding strategy: while RoundWithPrecision
// rounds to a specified number of decimal places based on the usual rounding rules, RoundUp always rounds up
// to the nearest integer.
//
// This function is generic and can work with any type that satisfies the basicutils.Numeric constraint,
// which includes integer and floating-point numeric types.
//
// Parameters:
// - value: The numeric value to be rounded up. This can be any type that satisfies the basicutils.Numeric constraint.
//
// Returns:
// - The value rounded up to the nearest integer, with the same type as the input `value`.
//
// Example:
// - RoundUp(0.1) returns 1 (whereas RoundWithPrecision(0.1, 0) would return 0)
// - RoundUp(1.5) returns 2 (whereas RoundWithPrecision(1.5, 0) would return 2)
// - RoundUp(-1.5) returns -1 (whereas RoundWithPrecision(-1.5, 0) would return -2)
//
// Note:
// This function always rounds numbers up to the next highest integer. This is different from RoundWithPrecision,
// which rounds to the nearest integer based on the fractional part and specified precision. If you need traditional
// rounding to the nearest value or rounding down, you should use RoundWithPrecision or a similar function.
func RoundUp[T basicutils.Numeric](value T) T {
	return T(math.Ceil(float64(value)))
}

func MaxValue[T basicutils.Numeric]() T {
	switch v := any(*new(T)).(type) {
	case float32:
		return any(float32(math.MaxFloat32)).(T)
	case float64:
		return any(math.MaxFloat64).(T)
	case int:
		return any(math.MaxInt).(T)
	case int8:
		return any(int8(math.MaxInt8)).(T)
	case int16:
		return any(int16(math.MaxInt16)).(T)
	case int32:
		return any(int32(math.MaxInt32)).(T)
	case int64:
		return any(int64(math.MaxInt64)).(T)
	case uint:
		return any(^uint(0)).(T)
	case uint8:
		return any(uint8(math.MaxUint8)).(T)
	case uint16:
		return any(uint16(math.MaxUint16)).(T)
	case uint32:
		return any(uint32(math.MaxUint32)).(T)
	case uint64:
		return any(^uint64(0)).(T)
	default:
		panic(fmt.Sprintf("Unhandled type: %T", v))
	}
}
