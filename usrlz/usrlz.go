/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package usrlz

import (
	"reflect"
	"unsafe"
)

// ToBytes converts an object of type T, passed as a pointer P, into a byte slice.
// It uses reflection and unsafe package to achieve this conversion.
//
// Type Parameters:
// - T: The type of the object.
// - P: The pointer type to the object (should be *T).
//
// Parameters:
// - obj: The object to convert, passed as a pointer of type P.
//
// Returns:
// - A byte slice representing the raw bytes of the object.
//
// Panics:
// - If the object is not addressable (cannot take the address of the value).
//
// Note:
//   - The object must be passed as a pointer, and it must be addressable.
//   - This function uses unsafe operations, which can lead to undefined behavior
//     if not used carefully. Ensure that the object remains valid for the duration
//     of the byte slice usage.
func ToBytes[T any, P *T](obj P) []byte {
	val := reflect.ValueOf(obj).Elem() // Dereference the pointer

	if !val.CanAddr() {
		panic("value is not addressable")
	}

	size := val.Type().Size()
	ptr := unsafe.Pointer(val.UnsafeAddr())
	slice := unsafe.Slice((*byte)(ptr), size)

	return slice
}
