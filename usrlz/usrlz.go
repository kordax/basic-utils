/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package usrlz

import (
	"fmt"
	"reflect"
	"unsafe"
)

// ToBytes converts a fixed-layout object of type T, passed as a pointer P, into a byte slice copy.
// Types containing references, such as strings, slices, maps, pointers, interfaces, funcs, or chans,
// are rejected because their in-memory representation is not stable serialized data.
//
// Type Parameters:
// - T: The type of the object.
// - P: The pointer type to the object (should be *T).
//
// Parameters:
// - obj: The object to convert, passed as a pointer of type P.
//
// Returns:
// - A byte slice copy representing the raw bytes of the object.
//
// Panics:
// - If the object is not addressable (cannot take the address of the value).
// - If the object contains reference-bearing fields.
//
// Note:
//   - The object must be passed as a pointer, and it must be addressable.
//   - The returned slice is detached from the source object.
func ToBytes[T any, P *T](obj P) []byte {
	val := reflect.ValueOf(obj).Elem() // Dereference the pointer

	if !val.CanAddr() {
		panic("value is not addressable")
	}

	if hasReferences(val.Type()) {
		panic(fmt.Sprintf("type %s contains reference-bearing fields", val.Type()))
	}

	size := val.Type().Size()
	if size == 0 {
		return []byte{}
	}

	ptr := unsafe.Pointer(val.UnsafeAddr())
	slice := unsafe.Slice((*byte)(ptr), size)
	result := make([]byte, len(slice))
	copy(result, slice)

	return result
}

func hasReferences(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.String, reflect.UnsafePointer:
		return true
	case reflect.Array:
		return hasReferences(t.Elem())
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if hasReferences(t.Field(i).Type) {
				return true
			}
		}
	}

	return false
}
