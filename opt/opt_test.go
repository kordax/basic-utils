/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package opt

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestPresent tests the Present method.
func TestPresent(t *testing.T) {
	// Test when value is present
	o := Of(42)
	if !o.Present() {
		t.Error("Expected Present() to return true")
	}

	// Test when value is not present
	o = Null[int]()
	if o.Present() {
		t.Error("Expected Present() to return false")
	}
}

// TestIfPresent tests the IfPresent method.
func TestIfPresent(t *testing.T) {
	// Test when value is present
	o := Of(42)
	var result int
	o.IfPresent(func(v int) {
		result = v
	})
	if result != 42 {
		assert.Fail(t, fmt.Sprintf("Expected IfPresent to execute the provided function with value 42, but got %d", result))
	}

	// Test when value is not present
	o = Null[int]()
	o.IfPresent(func(v int) {
		assert.Fail(t, "IfPresent should not execute the provided function when value is not present")
	})
}

// TestNull tests the Null method.
func TestNull(t *testing.T) {
	o := Null[int]()
	if o.Present() {
		t.Error("Expected Null() to return an Opt without a value")
	}
}

// TestOf tests the Of method.
func TestOf(t *testing.T) {
	o := Of(42)
	if !o.Present() {
		t.Error("Expected Of to create an Opt with a value")
	}
	if *o.Get() != 42 {
		assert.Fail(t, fmt.Sprintf("Expected Of to create an Opt with value 42, but got %v", *o.Get()))
	}
}

// TestOfNullable tests the OfNullable method.
func TestOfNullable(t *testing.T) {
	// Test with a non-nil value
	value := 42
	o := OfNullable(&value)
	if !o.Present() {
		t.Error("Expected OfNullable to create an Opt with a value")
	}
	if *o.Get() != 42 {
		assert.Fail(t, fmt.Sprintf("Expected OfNullable to create an Opt with value 42, but got %v", *o.Get()))
	}

	// Test with a nil value
	o = OfNullable[int](nil)
	if o.Present() {
		t.Error("Expected OfNullable to create an Opt without a value")
	}
}

// TestOfString tests the OfString method.
func TestOfString(t *testing.T) {
	// Test with a non-empty string
	o := OfString("hello")
	if !o.Present() {
		t.Error("Expected OfString to create an Opt with a value")
	}
	if *o.Get() != "hello" {
		assert.Fail(t, fmt.Sprintf("Expected OfString to create an Opt with value 'hello', but got %v", *o.Get()))
	}

	// Test with an empty string
	o = OfString("")
	if o.Present() {
		t.Error("Expected OfString to create an Opt without a value")
	}
}

// TestOfBool tests the OfBool method.
func TestOfBool(t *testing.T) {
	// Test with true
	o := OfBool(true)
	if o.Present() {
		t.Error("Expected OfBool to create an Opt without a value")
	}

	// Test with false
	o = OfBool(false)
	if !o.Present() {
		t.Error("Expected OfBool to create an Opt with a value")
	}
	if *o.Get() != false {
		assert.Fail(t, fmt.Sprintf("Expected OfBool to create an Opt with value false, but got %v", *o.Get()))
	}
}

// TestOfNumeric tests the OfNumeric method.
func TestOfNumeric(t *testing.T) {
	// Test with non-zero numeric value
	o := OfNumeric(42)
	if !o.Present() {
		t.Error("Expected OfNumeric to create an Opt with a value")
	}
	if *o.Get() != 42 {
		assert.Fail(t, fmt.Sprintf("Expected OfNumeric to create an Opt with value 42, but got %v", *o.Get()))
	}

	// Test with zero numeric value
	o = OfNumeric(0)
	if o.Present() {
		t.Error("Expected OfNumeric to create an Opt without a value")
	}
}

// TestOfCond tests the OfCond method.
func TestOfCond(t *testing.T) {
	// Test with a value that matches the condition
	o := OfCond(42, func(v *int) bool {
		return *v > 0
	})
	if !o.Present() {
		t.Error("Expected OfCond to create an Opt with a value")
	}
	if *o.Get() != 42 {
		assert.Fail(t, fmt.Sprintf("Expected OfCond to create an Opt with value 42, but got %v", *o.Get()))
	}

	// Test with a value that doesn't match the condition
	o = OfCond(0, func(v *int) bool {
		return *v > 0
	})
	if o.Present() {
		t.Error("Expected OfCond to create an Opt without a value")
	}
}

// TestOfUnix tests the OfUnix method.
func TestOfUnix(t *testing.T) {
	o := OfUnix(int64(time.Now().Unix()))
	if !o.Present() {
		t.Error("Expected OfUnix to create an Opt with a value")
	}
	if o.Get() == nil || reflect.TypeOf(o.Get()).Elem().Name() != "Time" {
		assert.Fail(t, fmt.Sprintf("Expected OfUnix to create an Opt with value of type time.Time, but got %v", o.Get()))
	}
}

// TestOfBuilder tests the OfBuilder method.
func TestOfBuilder(t *testing.T) {
	build := func() string {
		return "hello"
	}
	o := OfBuilder(build)
	if !o.Present() {
		t.Error("Expected OfBuilder to create an Opt with a value")
	}
	if *o.Get() != "hello" {
		assert.Fail(t, fmt.Sprintf("Expected OfBuilder to create an Opt with value 'hello', but got %v", *o.Get()))
	}
}

// TestOrElse tests the OrElse method.
func TestOrElse(t *testing.T) {
	o := Null[int]()
	result := o.OrElse(42)
	if result != 42 {
		assert.Fail(t, fmt.Sprintf("Expected OrElse to return 42, but got %v", result))
	}

	o = Of(24)
	result = o.OrElse(42)
	if result != 24 {
		assert.Fail(t, fmt.Sprintf("Expected OrElse to return 24, but got %v", result))
	}
}

// TestGet tests the Get method.
func TestGet(t *testing.T) {
	o := Of[int](42)
	result := o.Get()
	if result == nil || *result != 42 {
		assert.Fail(t, fmt.Sprintf("Expected Get to return a pointer to value 42, but got %v", result))
	}

	o = Null[int]()
	result = o.Get()
	if result != nil {
		t.Error("Expected Get to return nil for Opt without a value")
	}
}

// TestSet tests the Set method.
func TestSet(t *testing.T) {
	o := Null[int]()
	value := 42
	o.Set(&value)
	if !o.Present() {
		t.Error("Expected Set to set a value in Opt")
	}
	if *o.Get() != 42 {
		assert.Fail(t, fmt.Sprintf("Expected Set to set value 42 in Opt, but got %v", *o.Get()))
	}
}

// TestGetAs tests the GetAs method.
func TestGetAs(t *testing.T) {
	o := Of[int](42)
	result := o.GetAs(func(t *int) interface{} {
		return *t
	})
	if result != 42 {
		assert.Fail(t, fmt.Sprintf("Expected GetAs to return value 42, but got %v", result))
	}
}

// TestUnmarshalJSON tests the UnmarshalJSON method.
func TestUnmarshalJSON(t *testing.T) {
	// Test with a JSON value
	data := []byte(`42`)
	var o Opt[int]
	err := o.UnmarshalJSON(data)
	assert.NoError(t, err)
	if !o.Present() {
		t.Error("Expected UnmarshalJSON to set a value in Opt")
	}
	if *o.Get() != 42 {
		assert.Fail(t, fmt.Sprintf("Expected UnmarshalJSON to set value 42 in Opt, but got %v", *o.Get()))
	}

	// Test with a null JSON value
	data = []byte(`null`)
	o = Of(42)
	err = o.UnmarshalJSON(data)
	assert.NoError(t, err)
	if o.Present() {
		t.Error("Expected UnmarshalJSON to set Opt value to nil")
	}
}

// TestMarshalJSON tests the MarshalJSON method.
func TestMarshalJSON(t *testing.T) {
	// Test with a value in Opt
	o := Of[int](42)
	data, err := o.MarshalJSON()
	assert.NoError(t, err)
	expectedData := []byte(`42`)
	if !reflect.DeepEqual(data, expectedData) {
		assert.Fail(t, fmt.Sprintf("Expected MarshalJSON to return %s, but got %s", expectedData, data))
	}

	// Test without a value in Opt
	o = Null[int]()
	data, err = o.MarshalJSON()
	assert.NoError(t, err)
	expectedData = []byte(`0`)
	if string(data) != string(expectedData) {
		assert.Fail(t, fmt.Sprintf("Expected MarshalJSON to return %s, but got %s", expectedData, data))
	}
}

// TestValue tests the Value method.
func TestValue(t *testing.T) {
	// Test with a value in Opt
	o := Of[int](42)
	value, err := o.Value()
	if err != nil {
		assert.Fail(t, fmt.Sprintf("Value returned an unexpected error: %v", err))
	}
	if value != int64(42) {
		assert.Fail(t, fmt.Sprintf("Expected Value to return int64(42), but got %v", value))
	}

	// Test without a value in Opt
	o = Null[int]()
	value, err = o.Value()
	if err != nil {
		assert.Fail(t, fmt.Sprintf("Value returned an unexpected error: %v", err))
	}
	if value != nil {
		assert.Fail(t, fmt.Sprintf("Expected Value to return nil, but got %v", value))
	}
}

// TestScan tests the Scan method.
func TestScan(t *testing.T) {
	// Test with a string source
	var o Opt[int]
	err := o.Scan("42")
	if err != nil {
		assert.Fail(t, fmt.Sprintf("Scan returned an unexpected error: %v", err))
	}
	if !o.Present() {
		t.Error("Expected Scan to set a value in Opt")
	}
	if *o.Get() != 42 {
		assert.Fail(t, fmt.Sprintf("Expected Scan to set value 42 in Opt, but got %v", *o.Get()))
	}

	// Test with a []uint8 source
	o = Of(42)
	err = o.Scan([]uint8(`24`))
	if err != nil {
		assert.Fail(t, fmt.Sprintf("Scan returned an unexpected error: %v", err))
	}
	if !o.Present() {
		t.Error("Expected Scan to set a value in Opt")
	}
	if *o.Get() != 24 {
		assert.Fail(t, fmt.Sprintf("Expected Scan to set value 24 in Opt, but got %v", *o.Get()))
	}

	// Test with a nil source
	o = Of(42)
	err = o.Scan(nil)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("Scan returned an unexpected error: %v", err))
	}
	if o.Present() {
		t.Error("Expected Scan to set Opt value to nil")
	}

	// Test with an incompatible source type
	o = Of[int](42)
	err = o.Scan(true)
	if err == nil {
		t.Error("Scan should return an error for incompatible source type")
	}
	expectedError := "incompatible type for Opt[*int]: bool"
	if err.Error() != expectedError {
		assert.Fail(t, fmt.Sprintf("Expected Scan to return error '%s', but got '%s'", expectedError, err.Error()))
	}
}
