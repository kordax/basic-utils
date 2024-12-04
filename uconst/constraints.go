/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023-2024.
 */

package uconst

type Numeric interface {
	Integer | Float
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Float interface {
	~float32 | ~float64
}

type SignedNumeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

type Stringable interface {
	Numeric | ~bool | ~string
}

type BasicType interface {
	~string | ~bool | ~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 |
		*string | *bool | *int | *int8 | *int16 | *int32 | *int64 |
		*uint | *uint8 | *uint16 | *uint32 | *uint64 | *float32 | *float64
}
