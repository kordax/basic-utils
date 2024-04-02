/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uref

import "reflect"

func Ref[T any](t T) *T {
	return &t
}

func Compare[T comparable](ptr1, ptr2 *T) bool {
	if ptr1 == nil || ptr2 == nil {
		return ptr1 == ptr2
	}

	return *ptr1 == *ptr2
}

// CompareF compares two types using their pointers. Values passed to compare func never take nil/zero vales.
func CompareF[T any](ptr1, ptr2 *T, compare func(t1, t2 *T) bool) bool {
	if ptr1 == nil || ptr2 == nil {
		return ptr1 == ptr2
	}

	return compare(ptr1, ptr2)
}

// Or returns value of val or 'other', but supports any other type
func Or[R any](val any, other R) R {
	v := reflect.ValueOf(val)
	if !v.IsValid() || v.IsZero() {
		return other
	}

	if v.Type() != reflect.TypeOf(other) {
		return other
	}

	return val.(R)
}

// Do safely executes function 'do' in case ptr is not nil or returns 'nil' otherwise.
func Do[T any](ptr *T, do func(v T) *T) *T {
	if ptr == nil {
		return nil
	} else {
		return do(*ptr)
	}
}

// Def behaves as Or(val, *new(R)), so it returns default value if value is not present or types are different
func Def[R any](val R) R {
	v := reflect.ValueOf(val)
	if !v.IsValid() || v.IsZero() {
		newR := reflect.New(v.Type()).Elem()

		if v.Kind() == reflect.Ptr {
			newPtr := reflect.New(v.Type().Elem())
			return newPtr.Interface().(R)
		}

		return newR.Interface().(R)
	}

	other := *new(R)
	if v.Type() != reflect.TypeOf(other) &&
		v.Type() != reflect.TypeOf(&other) {
		return other
	}

	return val
}
