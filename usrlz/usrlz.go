/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package usrlz

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// ToBytes encodes a fixed-layout object using a deterministic little-endian representation.
// Types containing references, such as strings, slices, maps, pointers, interfaces, funcs, or chans,
// are rejected because their in-memory representation is not stable serialized data.
// Struct padding is omitted, and platform-sized integers are encoded as 64-bit values.
//
// Type Parameters:
// - T: The type of the object.
// - P: The pointer type to the object (should be *T).
//
// Parameters:
// - obj: The object to encode, passed as a pointer of type P.
//
// Returns:
// - A detached byte slice containing the encoded object.
//
// Panics:
// - If obj is nil.
// - If the object contains reference-bearing fields.
// - If encoding fails for an unsupported fixed-layout type.
func ToBytes[T any, P *T](obj P) []byte {
	if obj == nil {
		panic("object is nil")
	}

	t := reflect.TypeOf(obj).Elem()
	if hasReferences(t) {
		panic(fmt.Sprintf("type %s contains reference-bearing fields", t))
	}

	var result bytes.Buffer
	if err := writeValue(&result, reflect.ValueOf(obj).Elem()); err != nil {
		panic(fmt.Sprintf("encode type %s: %v", t, err))
	}
	if result.Len() == 0 {
		return []byte{}
	}

	return result.Bytes()
}

func writeValue(result *bytes.Buffer, value reflect.Value) error {
	switch value.Kind() {
	case reflect.Bool:
		return binary.Write(result, binary.LittleEndian, value.Bool())
	case reflect.Int8:
		return binary.Write(result, binary.LittleEndian, int8(value.Int())) // #nosec G115 -- reflect.Int8 guarantees the value fits int8
	case reflect.Int16:
		return binary.Write(result, binary.LittleEndian, int16(value.Int())) // #nosec G115 -- reflect.Int16 guarantees the value fits int16
	case reflect.Int32:
		return binary.Write(result, binary.LittleEndian, int32(value.Int())) // #nosec G115 -- reflect.Int32 guarantees the value fits int32
	case reflect.Int, reflect.Int64:
		return binary.Write(result, binary.LittleEndian, value.Int())
	case reflect.Uint8:
		return result.WriteByte(byte(value.Uint())) // #nosec G115 -- reflect.Uint8 guarantees the value fits uint8
	case reflect.Uint16:
		return binary.Write(result, binary.LittleEndian, uint16(value.Uint())) // #nosec G115 -- reflect.Uint16 guarantees the value fits uint16
	case reflect.Uint32:
		return binary.Write(result, binary.LittleEndian, uint32(value.Uint())) // #nosec G115 -- reflect.Uint32 guarantees the value fits uint32
	case reflect.Uint, reflect.Uint64, reflect.Uintptr:
		return binary.Write(result, binary.LittleEndian, value.Uint())
	case reflect.Float32:
		return binary.Write(result, binary.LittleEndian, float32(value.Float()))
	case reflect.Float64:
		return binary.Write(result, binary.LittleEndian, value.Float())
	case reflect.Complex64:
		return binary.Write(result, binary.LittleEndian, complex64(value.Complex()))
	case reflect.Complex128:
		return binary.Write(result, binary.LittleEndian, value.Complex())
	case reflect.Array:
		for i := 0; i < value.Len(); i++ {
			if err := writeValue(result, value.Index(i)); err != nil {
				return err
			}
		}
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			if value.Type().Field(i).Name == "_" {
				field = reflect.Zero(field.Type())
			}
			if err := writeValue(result, field); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported kind %s", value.Kind())
	}

	return nil
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
