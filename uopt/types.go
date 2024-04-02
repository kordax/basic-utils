/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uopt

type Bool = Opt[bool]
type String = Opt[string]

type Int = Opt[int]
type Int8 = Opt[int8]
type Int16 = Opt[int16]
type Int32 = Opt[int32]
type Int64 = Opt[int64]

type Uint = Opt[uint]
type Uint8 = Opt[uint8]
type Uint16 = Opt[uint16]
type Uint32 = Opt[uint32]
type Uint64 = Opt[uint64]

type Float32 = Opt[float32]
type Float64 = Opt[float64]

type Complex64 = Opt[complex64]
type Complex128 = Opt[complex128]

type Byte = Opt[byte]
type Rune = Opt[rune]

func NullBool() Bool {
	return Bool{v: nil}
}

func NullString() String {
	return String{v: nil}
}

func NullInt() Int {
	return Int{v: nil}
}

func NullInt8() Int8 {
	return Int8{v: nil}
}

func NullInt16() Int16 {
	return Int16{v: nil}
}

func NullInt32() Int32 {
	return Int32{v: nil}
}

func NullInt64() Int64 {
	return Int64{v: nil}
}

func NullUint() Uint {
	return Uint{v: nil}
}

func NullUint8() Uint8 {
	return Uint8{v: nil}
}

func NullUint16() Uint16 {
	return Uint16{v: nil}
}

func NullUint32() Uint32 {
	return Uint32{v: nil}
}

func NullUint64() Uint64 {
	return Uint64{v: nil}
}

func NullFloat32() Float32 {
	return Float32{v: nil}
}

func NullFloat64() Float64 {
	return Float64{v: nil}
}

func NullComplex64() Complex64 {
	return Complex64{v: nil}
}

func NullComplex128() Complex128 {
	return Complex128{v: nil}
}

func NullByte() Byte {
	return Byte{v: nil}
}

func NullRune() Rune {
	return Rune{v: nil}
}
