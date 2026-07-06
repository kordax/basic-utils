/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package usrlz_test

import (
	"testing"
	"unsafe"

	"git.casinomodule.org/casino27/basic-utils/v3/usrlz"
	"github.com/stretchr/testify/assert"
)

type ComplexStruct struct {
	IntField        int
	FloatField      float64
	StringField     string
	BoolField       bool
	ArrayField      [3]int
	SliceField      []int
	MapField        map[string]int
	PointerField    *int
	NilPointerField *int
	StructField     SimpleStruct
}

type SimpleStruct struct {
	A int
	B float64
}

type EmptyStruct struct{}

type FixedStruct struct {
	A int64
	B [2]uint32
}

func TestToBytes(t *testing.T) {
	fixed := FixedStruct{A: 123, B: [2]uint32{456, 789}}

	bytes := usrlz.ToBytes(&fixed)

	var result *FixedStruct
	ptr := unsafe.Pointer(&bytes[0])
	result = (*FixedStruct)(ptr)

	assert.Equal(t, fixed, *result)
}

func TestToBytesReturnsStableCopy(t *testing.T) {
	value := FixedStruct{A: 10, B: [2]uint32{20, 30}}

	bytes := usrlz.ToBytes(&value)
	value.A = 99

	result := (*FixedStruct)(unsafe.Pointer(&bytes[0]))
	assert.Equal(t, int64(10), result.A)
	assert.Equal(t, [2]uint32{20, 30}, result.B)
}

func TestToBytesRejectsReferenceBearingTypes(t *testing.T) {
	t.Run("string field", func(t *testing.T) {
		value := struct {
			Name string
		}{Name: "unsafe"}

		assert.Panics(t, func() {
			_ = usrlz.ToBytes(&value)
		})
	})

	t.Run("slice", func(t *testing.T) {
		value := []int{1, 2, 3}

		assert.Panics(t, func() {
			_ = usrlz.ToBytes(&value)
		})
	})

	t.Run("map", func(t *testing.T) {
		value := map[string]int{"one": 1}

		assert.Panics(t, func() {
			_ = usrlz.ToBytes(&value)
		})
	})
}

func TestToBytesEmptyStruct(t *testing.T) {
	emptyStruct := EmptyStruct{}

	bytes := usrlz.ToBytes(&emptyStruct)

	var result *EmptyStruct
	ptr := unsafe.Pointer(&bytes)
	result = (*EmptyStruct)(ptr)

	assert.Equal(t, emptyStruct, *result, "EmptyStruct does not match")
}

func TestToBytesNilPointer(t *testing.T) {
	var nilPointer *ComplexStruct = nil

	assert.Panics(t, func() {
		_ = usrlz.ToBytes(nilPointer)
	}, "The code did not panic with a nil pointer")
}

func TestToBytesSlice(t *testing.T) {
	slice := []int{1, 2, 3}

	assert.Panics(t, func() {
		_ = usrlz.ToBytes(&slice)
	})
}

func TestToBytesMap(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2}

	assert.Panics(t, func() {
		_ = usrlz.ToBytes(&m)
	})
}
