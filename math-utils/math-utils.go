/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package mathutils

import (
	"math"
	"sort"
)

// ClosestMatch finds the closest match in a slice. This implementation doesn't use Round.
func ClosestMatch(toMatch uint32, slice []uint32) (int, uint32) {
	var diff uint32 = math.MaxUint32
	var value uint32 = 0
	index := 0
	for i, val := range slice {
		d := AbsDiffUInt32(val, toMatch)
		if d < diff {
			diff = d
			index = i
			value = val
		}
	}

	return index, value
}

// AbsDiffUInt32 returns absolute diff result of a numeric type
func AbsDiffUInt32(one uint32, two uint32) uint32 {
	if one < two {
		return two - one
	}

	return one - two
}

// AbsValInt returns absolute value of a numeric type
func AbsValInt(val int) int {
	if val < 0 {
		return -val
	}

	return val
}

// ValOrMin returns value or min if value is less than min
func ValOrMin(val int, min int) int {
	if val < min {
		return min
	}

	return val
}

// MaxInt returns max value in array correspondingly
func MaxInt(array []int) int {
	if len(array) == 0 {
		return 0
	}

	sort.Ints(array)

	return array[len(array)-1]
}

// Min returns min numeric value found in array
func Min[T Numeric](array []T) T {
	if len(array) == 0 {
		return *new(T)
	}

	min := array[0]
	for _, a := range array {
		if a < min {
			min = a
		}
	}

	return min
}

// Max returns max numeric value found in array
func Max[T Numeric](array []T) T {
	if len(array) == 0 {
		return *new(T)
	}

	max := array[0]
	for _, a := range array {
		if a > max {
			max = a
		}
	}

	return max
}

// MinMaxInt returns min and max value in array correspondingly
func MinMaxInt(array []int) (int, int) {
	if len(array) == 0 {
		return 0, 0
	}

	min := array[0]
	max := array[0]
	for _, a := range array {
		if a < min {
			min = a
		}
		if a > max {
			max = a
		}
	}

	return min, max
}

// MinMaxUInt32 returns min and max value in array correspondingly
func MinMaxUInt32(array []uint32) (uint32, uint32) {
	if len(array) == 0 {
		return 0, 0
	}

	min := array[0]
	max := array[0]
	for _, a := range array {
		if a < min {
			min = a
		}
		if a > max {
			max = a
		}
	}

	return min, max
}

// MinMaxUInt64 returns min and max value in array correspondingly
func MinMaxUInt64(array []uint64) (uint64, uint64) {
	if len(array) == 0 {
		return 0, 0
	}

	min := array[0]
	max := array[0]
	for _, a := range array {
		if a < min {
			min = a
		}
		if a > max {
			max = a
		}
	}

	return min, max
}

// MinMaxFloat64 returns min and max value in array correspondingly
func MinMaxFloat64(array []float64) (float64, float64) {
	if len(array) == 0 {
		return 0, 0
	}

	min := array[0]
	max := array[0]
	for _, a := range array {
		if a < min {
			min = a
		}
		if a > max {
			max = a
		}
	}

	return min, max
}

// AvgInt returns average value from the array
func AvgInt(array []int) int {
	if len(array) == 0 {
		return 0
	}

	sum := SumInt(array)

	return sum / len(array)
}

// AvgUInt64 returns average value from the array
func AvgUInt64(array []uint64) uint64 {
	if len(array) == 0 {
		return 0
	}

	sum := SumUInt64(array)

	return sum / uint64(len(array))
}

// AvgInt64 returns average value from the array
func AvgInt64(array []int64) int64 {
	if len(array) == 0 {
		return 0
	}

	sum := SumInt64(array)

	return sum / int64(len(array))
}

// AvgFloat64 returns average value from the array
func AvgFloat64(array []float64) float64 {
	if len(array) == 0 {
		return 0
	}

	sum := SumFloat64(array)

	return sum / float64(len(array))
}

func MedInt(array []int) int {
	if len(array) == 0 {
		return 0.0
	}
	sort.Slice(array, func(i, j int) bool { return array[i] < array[j] })

	ln := len(array)
	if ln == 1 {
		return array[0]
	}

	var med int
	if ln%2 != 0 {
		med = array[(ln / 2)]
	} else {
		med = (array[(ln/2)-1] + array[(ln/2)]) / 2
	}

	return med
}

func MedUInt64(array []uint64) uint64 {
	if len(array) == 0 {
		return 0.0
	}
	sort.Slice(array, func(i, j int) bool { return array[i] < array[j] })

	ln := len(array)
	if ln == 1 {
		return array[0]
	}

	var med uint64
	if ln%2 != 0 {
		med = array[(ln / 2)]
	} else {
		med = (array[(ln/2)-1] + array[(ln/2)]) / 2
	}

	return med
}

func MedFloat64(array []float64) float64 {
	if len(array) == 0 {
		return 0.0
	}
	sort.Float64s(array)

	ln := len(array)
	if ln == 1 {
		return array[0]
	}

	var med float64
	if ln%2 != 0 {
		med = array[(ln / 2)]
	} else {
		med = (array[(ln/2)-1] + array[(ln/2)]) / 2
	}

	return med
}

func SumInt(array []int) int {
	sum := 0
	for _, v := range array {
		sum = sum + v
	}

	return sum
}

func SumInt64(array []int64) int64 {
	sum := int64(0)
	for _, v := range array {
		sum = sum + v
	}

	return sum
}

func SumUInt64(array []uint64) uint64 {
	sum := uint64(0)
	for _, v := range array {
		sum = sum + v
	}

	return sum
}

func SumFloat64(array []float64) float64 {
	sum := 0.0
	for _, a := range array {
		sum = sum + a
	}

	return sum
}

func MinMax[T Ordered](array []T) (T, T) {
	if len(array) == 0 {
		return *new(T), *new(T)
	}

	min := array[0]
	max := array[0]
	for _, a := range array {
		if a < min {
			min = a
		}
		if a > max {
			max = a
		}
	}

	return min, max
}

func MinMaxFromMap[K comparable, T Ordered](m map[K]T) (T, T) {
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

func RoundWithPrecision[T Numeric](value T, precision int) T {
	ratio := math.Pow(10, float64(precision))

	return T(math.Round(float64(value)*ratio) / ratio)
}

type Numeric interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type Float interface {
	float32 | float64
}

type SignedNumeric interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

type Ordered interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}
