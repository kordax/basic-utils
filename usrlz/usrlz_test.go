/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package usrlz_test

import (
	"testing"
	"unsafe"

	"github.com/kordax/basic-utils/v2/usrlz"
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

func TestToBytes(t *testing.T) {
	simpleStruct := SimpleStruct{A: 1, B: 2.0}
	intVal := 42
	complexStruct := ComplexStruct{
		IntField:        123,
		FloatField:      456.789,
		StringField:     "test",
		BoolField:       true,
		ArrayField:      [3]int{1, 2, 3},
		SliceField:      []int{4, 5, 6},
		MapField:        map[string]int{"one": 1, "two": 2},
		PointerField:    &intVal,
		NilPointerField: nil,
		StructField:     simpleStruct,
	}

	bytes := usrlz.ToBytes(&complexStruct)

	var result *ComplexStruct
	ptr := unsafe.Pointer(&bytes[0])
	result = (*ComplexStruct)(ptr)

	assert.Equal(t, complexStruct.IntField, result.IntField, "IntField does not match")
	assert.Equal(t, complexStruct.FloatField, result.FloatField, "FloatField does not match")
	assert.Equal(t, complexStruct.StringField, result.StringField, "StringField does not match")
	assert.Equal(t, complexStruct.BoolField, result.BoolField, "BoolField does not match")
	assert.Equal(t, complexStruct.ArrayField, result.ArrayField, "ArrayField does not match")
	assert.Equal(t, complexStruct.SliceField, result.SliceField, "SliceField does not match")
	assert.Equal(t, complexStruct.MapField, result.MapField, "MapField does not match")
	assert.Equal(t, complexStruct.PointerField, result.PointerField, "PointerField does not match")
	assert.Equal(t, complexStruct.NilPointerField, result.NilPointerField, "NilPointerField does not match")
	assert.Equal(t, complexStruct.StructField, result.StructField, "StructField does not match")
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

	bytes := usrlz.ToBytes(&slice)

	var result *[]int
	ptr := unsafe.Pointer(&bytes[0])
	result = (*[]int)(ptr)

	assert.EqualValues(t, slice, *result)
}

func TestToBytesMap(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2}

	bytes := usrlz.ToBytes(&m)

	var result *map[string]int
	ptr := unsafe.Pointer(&bytes[0])
	result = (*map[string]int)(ptr)

	assert.EqualValues(t, m, *result)
}
