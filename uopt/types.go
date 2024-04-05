/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uopt

import "time"

type OptBool = Opt[bool]
type OptString = Opt[string]

type OptInt = Opt[int]
type OptInt8 = Opt[int8]
type OptInt16 = Opt[int16]
type OptInt32 = Opt[int32]
type OptInt64 = Opt[int64]

type OptUint = Opt[uint]
type OptUint8 = Opt[uint8]
type OptUint16 = Opt[uint16]
type OptUint32 = Opt[uint32]
type OptUint64 = Opt[uint64]

type OptFloat32 = Opt[float32]
type OptFloat64 = Opt[float64]

type OptComplex64 = Opt[complex64]
type OptComplex128 = Opt[complex128]

type OptByte = Opt[byte]
type OptRune = Opt[rune]

type OptDuration = Opt[time.Duration]

func NullBool() OptBool {
	return OptBool{v: nil}
}

func NullString() OptString {
	return OptString{v: nil}
}

func NullInt() OptInt {
	return OptInt{v: nil}
}

func NullInt8() OptInt8 {
	return OptInt8{v: nil}
}

func NullInt16() OptInt16 {
	return OptInt16{v: nil}
}

func NullInt32() OptInt32 {
	return OptInt32{v: nil}
}

func NullInt64() OptInt64 {
	return OptInt64{v: nil}
}

func NullUint() OptUint {
	return OptUint{v: nil}
}

func NullUint8() OptUint8 {
	return OptUint8{v: nil}
}

func NullUint16() OptUint16 {
	return OptUint16{v: nil}
}

func NullUint32() OptUint32 {
	return OptUint32{v: nil}
}

func NullUint64() OptUint64 {
	return OptUint64{v: nil}
}

func NullFloat32() OptFloat32 {
	return OptFloat32{v: nil}
}

func NullFloat64() OptFloat64 {
	return OptFloat64{v: nil}
}

func NullComplex64() OptComplex64 {
	return OptComplex64{v: nil}
}

func NullComplex128() OptComplex128 {
	return OptComplex128{v: nil}
}

func NullByte() OptByte {
	return OptByte{v: nil}
}

func NullRune() OptRune {
	return OptRune{v: nil}
}

func NullDuration() OptDuration {
	return OptDuration{v: nil}
}
