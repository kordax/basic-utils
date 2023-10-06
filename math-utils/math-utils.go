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

	basic_utils "github.com/kordax/basic-utils"
)

// ClosestMatch finds the closest match in a slice. This implementation doesn't use Round.
// Returns index of a found element and the element. If no element was found then -1 as index will be returned.
func ClosestMatch[T basic_utils.Numeric](toMatch T, slice []T) (int, T) {
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
func AbsDiff[T basic_utils.Numeric](one T, two T) T {
	if one < two {
		return two - one
	}

	return one - two
}

// AbsVal returns absolute value of a numeric type
func AbsVal[T basic_utils.Numeric](val T) T {
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
func Min[T basic_utils.Numeric](array []T) T {
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
func Max[T basic_utils.Numeric](array []T) T {
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

// AvgInt returns average value from the array
func AvgInt[T basic_utils.Numeric](array []T) T {
	if len(array) == 0 {
		return 0
	}

	sum := Sum(array)

	return sum / T(len(array))
}

func Med[T basic_utils.Numeric](array []T) T {
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

func Sum[T basic_utils.Numeric](array []T) T {
	var sum T
	for _, v := range array {
		sum = sum + v
	}

	return sum
}

func MinMax[T basic_utils.Numeric](array []T) (T, T) {
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

func MinMaxFromMap[K comparable, T basic_utils.Numeric](m map[K]T) (T, T) {
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

func RoundWithPrecision[T basic_utils.Numeric](value T, precision int) T {
	ratio := math.Pow(10, float64(precision))

	return T(math.Round(float64(value)*ratio) / ratio)
}

func MaxValue[T basic_utils.Numeric]() T {
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
