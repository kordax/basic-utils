/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uref_test

import (
	"testing"

	"github.com/kordax/basic-utils/uref"
	"github.com/stretchr/testify/assert"
)

func TestRef(t *testing.T) {
	val := 5
	ptr := uref.Ref(val)
	if *ptr != val {
		t.Errorf("Expected %d, but got %d", val, *ptr)
	}
}

func TestCompare(t *testing.T) {
	a, b := 5, 5
	ptrA, ptrB := &a, &b
	if !uref.Compare(ptrA, ptrB) {
		t.Errorf("Compare function failed for equal values")
	}

	b = 6
	if uref.Compare(ptrA, ptrB) {
		t.Errorf("Compare function failed for unequal values")
	}

	if !uref.Compare[any](nil, nil) {
		t.Errorf("Compare function failed for nil values")
	}

	if uref.Compare(ptrA, nil) {
		t.Errorf("Compare function failed for one nil value")
	}
}

func TestCompareF(t *testing.T) {
	strA, strB := "hello", "hello"
	ptrA, ptrB := &strA, &strB
	if !uref.CompareF(ptrA, ptrB, func(t1, t2 *string) bool {
		return *t1 == *t2
	}) {
		t.Errorf("CompareF function failed for equal strings")
	}

	strB = "world"
	if uref.CompareF(ptrA, ptrB, func(t1, t2 *string) bool {
		return *t1 == *t2
	}) {
		t.Errorf("CompareF function failed for unequal strings")
	}
}

func TestOr(t *testing.T) {
	var valRef *string
	res := uref.Or(valRef, "default")
	assert.Equal(t, "default", res)

	val := "value"
	res = uref.Or(&val, "default")
	assert.Equal(t, val, res)
}

func TestOrRef(t *testing.T) {
	var value *string
	other := uref.Ref("Test")
	res := uref.OrRef(value, other)
	assert.EqualValues(t, other, res)

	value = uref.Ref("Value")
	res = uref.OrRef(value, other)
	assert.Equal(t, value, res)
}

func TestDo(t *testing.T) {
	val := 5
	ptr := &val
	res := uref.Do(ptr, func(v int) *int {
		v = v + 5
		return &v
	})

	if *res != 10 {
		t.Errorf("Expected 10, but got %d", *res)
	}

	var nilPtr *int
	res = uref.Do(nilPtr, func(v int) *int {
		v = v + 5
		return &v
	})

	if res != nil {
		t.Errorf("Expected nil, but got a value")
	}
}

func TestDef(t *testing.T) {
	var f *float64
	assert.Equal(t, 0.0, uref.Def(f))
	assert.Equal(t, 5, uref.Def(uref.Ref(5)))
	assert.Equal(t, 0, uref.Def(uref.Ref(0)))
	assert.Equal(t, "hello", uref.Def(uref.Ref("hello")))
	assert.Equal(t, "", uref.Def(uref.Ref("")))
	assert.Equal(t, true, uref.Def(uref.Ref(true)))
	assert.Equal(t, false, uref.Def(uref.Ref(false)))
	assert.Equal(t, []int{1, 2}, uref.Def(uref.Ref([]int{1, 2})))
	assert.Equal(t, []int{}, uref.Def(uref.Ref([]int{})))
}
