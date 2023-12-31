/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package refutils_test

import (
	"testing"

	refutils "github.com/kordax/basic-utils/ref-utils"
)

func TestRef(t *testing.T) {
	val := 5
	ptr := refutils.Ref(val)
	if *ptr != val {
		t.Errorf("Expected %d, but got %d", val, *ptr)
	}
}

func TestCompare(t *testing.T) {
	a, b := 5, 5
	ptrA, ptrB := &a, &b
	if !refutils.Compare(ptrA, ptrB) {
		t.Errorf("Compare function failed for equal values")
	}

	b = 6
	if refutils.Compare(ptrA, ptrB) {
		t.Errorf("Compare function failed for unequal values")
	}

	if !refutils.Compare[any](nil, nil) {
		t.Errorf("Compare function failed for nil values")
	}

	if refutils.Compare(ptrA, nil) {
		t.Errorf("Compare function failed for one nil value")
	}
}

func TestCompareF(t *testing.T) {
	strA, strB := "hello", "hello"
	ptrA, ptrB := &strA, &strB
	if !refutils.CompareF(ptrA, ptrB, func(t1, t2 *string) bool {
		return *t1 == *t2
	}) {
		t.Errorf("CompareF function failed for equal strings")
	}

	strB = "world"
	if refutils.CompareF(ptrA, ptrB, func(t1, t2 *string) bool {
		return *t1 == *t2
	}) {
		t.Errorf("CompareF function failed for unequal strings")
	}
}

func TestOr(t *testing.T) {
	val := ""
	res := refutils.Or(val, "default")
	if res != "default" {
		t.Errorf("Expected 'default', but got '%s'", res)
	}

	val = "value"
	res = refutils.Or(val, "default")
	if res != "value" {
		t.Errorf("Expected 'value', but got '%s'", res)
	}
}

func TestDo(t *testing.T) {
	val := 5
	ptr := &val
	res := refutils.Do(ptr, func(v int) *int {
		v = v + 5
		return &v
	})

	if *res != 10 {
		t.Errorf("Expected 10, but got %d", *res)
	}

	var nilPtr *int
	res = refutils.Do(nilPtr, func(v int) *int {
		v = v + 5
		return &v
	})

	if res != nil {
		t.Errorf("Expected nil, but got a value")
	}
}
