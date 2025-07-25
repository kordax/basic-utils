/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uopt_test

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/kordax/basic-utils/v2/uopt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPresent tests the Present method.
func TestPresent(t *testing.T) {
	// Test when value is present
	o := uopt.Of(42)
	if !o.Present() {
		t.Error("Expected Present() to return true")
	}

	// Test when value is not present
	o = uopt.Null[int]()
	if o.Present() {
		t.Error("Expected Present() to return false")
	}
}

// TestIfPresent tests the IfPresent method.
func TestIfPresent(t *testing.T) {
	// Test when value is present
	o := uopt.Of(42)
	var result int
	o.IfPresent(func(v int) {
		result = v
	})
	if result != 42 {
		assert.Fail(t, fmt.Sprintf("Expected IfPresent to execute the provided function with value 42, but got %d", result))
	}

	// Test when value is not present
	o = uopt.Null[int]()
	o.IfPresent(func(v int) {
		assert.Fail(t, "IfPresent should not execute the provided function when value is not present")
	})
}

// TestNull tests the Null method.
func TestNull(t *testing.T) {
	o := uopt.Null[int]()
	if o.Present() {
		t.Error("Expected Null() to return an Opt without a value")
	}
}

// TestOf tests the Of method.
func TestOf(t *testing.T) {
	o := uopt.Of(42)
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
	o := uopt.OfNullable(&value)
	if !o.Present() {
		t.Error("Expected OfNullable to create an Opt with a value")
	}
	if *o.Get() != 42 {
		assert.Fail(t, fmt.Sprintf("Expected OfNullable to create an Opt with value 42, but got %v", *o.Get()))
	}

	// Test with a nil value
	o = uopt.OfNullable[int](nil)
	if o.Present() {
		t.Error("Expected OfNullable to create an Opt without a value")
	}
}

// Testopt.OfString tests the opt.OfString method.
func TestOfString(t *testing.T) {
	// Test with a non-empty string
	o := uopt.OfString("hello")
	if !o.Present() {
		t.Error("Expected opt.OfString to create an Opt with a value")
	}
	if *o.Get() != "hello" {
		assert.Fail(t, fmt.Sprintf("Expected opt.OfString to create an Opt with value 'hello', but got %v", *o.Get()))
	}

	// Test with an empty string
	o = uopt.OfString("")
	if o.Present() {
		t.Error("Expected opt.OfString to create an Opt without a value")
	}
}

// TestOfBool tests the OfBool method.
func TestOfBool(t *testing.T) {
	// Test with true
	o := uopt.OfBool(true)
	if !o.Present() {
		t.Error("Expected OfBool to create an Opt with value")
	}

	// Test with false
	o = uopt.OfBool(false)
	if o.Present() {
		t.Error("Expected OfBool to create an Opt without a value")
	}
}

// TestOfNumeric tests the OfNumeric method.
func TestOfNumeric(t *testing.T) {
	// Test with non-zero numeric value
	o := uopt.OfNumeric(42)
	if !o.Present() {
		t.Error("Expected OfNumeric to create an Opt with a value")
	}
	if *o.Get() != 42 {
		assert.Fail(t, fmt.Sprintf("Expected OfNumeric to create an Opt with value 42, but got %v", *o.Get()))
	}

	// Test with zero numeric value
	o = uopt.OfNumeric(0)
	if o.Present() {
		t.Error("Expected OfNumeric to create an Opt without a value")
	}
}

// TestOfCond tests the OfCond method.
func TestOfCond(t *testing.T) {
	// Test with a value that matches the condition
	o := uopt.OfCond(42, func(v int) bool {
		return v > 0
	})
	if !o.Present() {
		t.Error("Expected OfCond to create an Opt with a value")
	}
	if *o.Get() != 42 {
		assert.Fail(t, fmt.Sprintf("Expected OfCond to create an Opt with value 42, but got %v", *o.Get()))
	}

	// Test with a value that doesn't match the condition
	o = uopt.OfCond(0, func(v int) bool {
		return v > 0
	})
	if o.Present() {
		t.Error("Expected OfCond to create an Opt without a value")
	}
}

// TestOfUnix tests the OfUnix method.
func TestOfUnix(t *testing.T) {
	o := uopt.OfUnix(time.Now().Unix())
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
	o := uopt.OfBuilder(build)
	if !o.Present() {
		t.Error("Expected OfBuilder to create an Opt with a value")
	}
	if *o.Get() != "hello" {
		assert.Fail(t, fmt.Sprintf("Expected OfBuilder to create an Opt with value 'hello', but got %v", *o.Get()))
	}
}

// TestOrElse tests the OrElse method.
func TestOrElse(t *testing.T) {
	o := uopt.Null[int]()
	result := o.OrElse(42)
	if result != 42 {
		assert.Fail(t, fmt.Sprintf("Expected OrElse to return 42, but got %v", result))
	}

	o = uopt.Of(24)
	result = o.OrElse(42)
	if result != 24 {
		assert.Fail(t, fmt.Sprintf("Expected OrElse to return 24, but got %v", result))
	}
}

// TestOrDef tests the Def method.
func TestOrDef(t *testing.T) {
	o := uopt.Null[int]()
	result := o.Def()
	if result != 0 {
		assert.Fail(t, fmt.Sprintf("Expected OrElse to return 0, but got %v", result))
	}

	o = uopt.Of(24)
	result = o.Def()
	if result != 24 {
		assert.Fail(t, fmt.Sprintf("Expected OrElse to return 24, but got %v", result))
	}
}

// TestGet tests the Get method.
func TestGet(t *testing.T) {
	o := uopt.Of[int](42)
	result := o.Get()
	if result == nil || *result != 42 {
		assert.Fail(t, fmt.Sprintf("Expected Get to return a pointer to value 42, but got %v", result))
	}

	o = uopt.Null[int]()
	result = o.Get()
	if result != nil {
		t.Error("Expected Get to return nil for Opt without a value")
	}
}

// TestSet tests the Set method.
func TestSet(t *testing.T) {
	o := uopt.Null[int]()
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
	o := uopt.Of[int](42)
	result := o.GetAs(func(t int) interface{} {
		return t
	})
	if result != 42 {
		assert.Fail(t, fmt.Sprintf("Expected GetAs to return value 42, but got %v", result))
	}
}

// TestUnmarshalJSON tests the UnmarshalJSON method.
func TestUnmarshalJSON(t *testing.T) {
	// Test with a JSON value
	data := []byte(`42`)
	var o uopt.Opt[int]
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
	o = uopt.Of(42)
	err = o.UnmarshalJSON(data)
	assert.NoError(t, err)
	if o.Present() {
		t.Error("Expected UnmarshalJSON to set Opt value to nil")
	}
}

// TestMarshalJSON tests the MarshalJSON method.
func TestMarshalJSON(t *testing.T) {
	// Test with a value in Opt
	o := uopt.Of[int](42)
	data, err := o.MarshalJSON()
	assert.NoError(t, err)
	expectedData := []byte(`42`)
	if !reflect.DeepEqual(data, expectedData) {
		assert.Fail(t, fmt.Sprintf("Expected MarshalJSON to return %s, but got %s", expectedData, data))
	}

	// Test without a value in Opt
	o = uopt.Null[int]()
	data, err = o.MarshalJSON()
	assert.NoError(t, err)
	expectedData = []byte("null")
	if string(data) != string(expectedData) {
		assert.Fail(t, fmt.Sprintf("Expected MarshalJSON to return %s, but got %s", expectedData, data))
	}
}

// TestValue tests the Value method.
func TestValue(t *testing.T) {
	// Test with a value in Opt
	o := uopt.Of[int](42)
	value, err := o.Value()
	if err != nil {
		assert.Fail(t, fmt.Sprintf("Value returned an unexpected error: %v", err))
	}
	if value != int64(42) {
		assert.Fail(t, fmt.Sprintf("Expected Value to return int64(42), but got %v", value))
	}

	// Test without a value in Opt
	o = uopt.Null[int]()
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
	t.Run("string source", func(t *testing.T) {
		var o uopt.Opt[int]
		err := o.Scan("42")
		require.NoError(t, err, fmt.Sprintf("Scan returned an unexpected error: %v", err))
		require.True(t, o.Present(), "Expected Scan to set a value in Opt")
		require.EqualValues(t, 42, *o.Get(), "Expected Scan to set value 42 in Opt, but got %v", *o.Get())
	})

	// Test with a []uint8 source
	t.Run("[]uint8 source", func(t *testing.T) {
		o := uopt.Of(42)
		err := o.Scan([]uint8(`24`))
		require.NoError(t, err, fmt.Sprintf("Scan returned an unexpected error: %v", err))
		require.True(t, o.Present(), "Expected Scan to set a value in Opt")
		require.EqualValues(t, 24, *o.Get(), "Expected Scan to set value 24 in Opt, but got %v", *o.Get())
	})

	// Test with a nil source
	t.Run("nil source", func(t *testing.T) {
		o := uopt.Of(42)
		_ = o.Scan(nil)
		require.False(t, o.Present(), "Expected Scan to set Opt value to nil")
	})

	// Test with an incompatible source type
	t.Run("incompatible source", func(t *testing.T) {
		o := uopt.Of[int](42)
		err := o.Scan(true)
		require.Error(t, err, "Scan should return an error for incompatible source type")
		expectedError := "incompatible type for Opt[int]: bool"
		assert.ErrorContains(t, errors.New(expectedError), err.Error())
	})

	// Test with a time.Time source
	t.Run("time.Time source", func(t *testing.T) {
		var timeOpt uopt.Opt[time.Time]
		expectedTime := time.Now()
		err := timeOpt.Scan(expectedTime)
		require.NoError(t, err, "Scan returned an unexpected error")
		assert.True(t, timeOpt.Present(), "Expected Scan to set a value in Opt")
		assert.Equal(t, expectedTime, *timeOpt.Get(), "Expected Scan to set the correct time in Opt")
	})

	// Test with a float64 source
	t.Run("float64 source", func(t *testing.T) {
		var floatOpt uopt.Opt[float64]
		err := floatOpt.Scan("42.56")
		require.NoError(t, err, "Scan returned an unexpected error")
		assert.True(t, floatOpt.Present(), "Expected Scan to set a value in Opt")
		assert.Equal(t, 42.56, *floatOpt.Get(), "Expected Scan to set value 42.56 in Opt")
	})
}

func TestOpt_Present(t *testing.T) {
	opt1 := uopt.Of(123)
	opt2 := uopt.Null[int]()

	if !opt1.Present() {
		t.Error("Expected opt1 to be present")
	}

	if opt2.Present() {
		t.Error("Expected opt2 to not be present")
	}
}

func TestOpt_IfPresent(t *testing.T) {
	o := uopt.Of(123)
	called := false
	o.IfPresent(func(v int) {
		called = true
		if v != 123 {
			t.Error("Unexpected value in IfPresent")
		}
	})

	if !called {
		t.Error("IfPresent function was not called")
	}
}

func TestOpt_Null(t *testing.T) {
	o := uopt.Null[int]()
	if o.Present() {
		t.Error("Expected Null to not be present")
	}
}

func TestOpt_Of(t *testing.T) {
	o := uopt.Of(123)
	if !o.Present() || *o.Get() != 123 {
		t.Error("Of did not set the expected value")
	}
}

func TestOpt_OfNullable(t *testing.T) {
	var nullPtr *int = nil
	val := 123
	opt1 := uopt.OfNullable(nullPtr)
	opt2 := uopt.OfNullable(&val)

	if opt1.Present() {
		t.Error("Expected opt1 to not be present")
	}

	if !opt2.Present() || *opt2.Get() != 123 {
		t.Error("opt2 did not set the expected value")
	}
}

func TestOpt_OfString(t *testing.T) {
	opt1 := uopt.OfString("")
	opt2 := uopt.OfString("hello")

	if opt1.Present() {
		t.Error("Expected opt1 to not be present")
	}

	if !opt2.Present() || *opt2.Get() != "hello" {
		t.Error("opt2 did not set the expected value")
	}
}

func TestOpt_OfBool(t *testing.T) {
	opt1 := uopt.OfBool(false)
	opt2 := uopt.OfBool(true)

	if opt1.Present() {
		t.Error("Expected opt1 to not be present")
	}

	if !opt2.Present() || *opt2.Get() != true {
		t.Error("opt2 did not set the expected value")
	}
}

func TestOpt_OfNumeric(t *testing.T) {
	opt1 := uopt.OfNumeric(0)
	opt2 := uopt.OfNumeric(123)

	if opt1.Present() {
		t.Error("Expected opt1 to not be present")
	}

	if !opt2.Present() || *opt2.Get() != 123 {
		t.Error("opt2 did not set the expected value")
	}
}

func TestOpt_OfCond(t *testing.T) {
	opt1 := uopt.OfCond(123, func(v int) bool {
		return v > 200
	})

	opt2 := uopt.OfCond(123, func(v int) bool {
		return v < 200
	})

	if opt1.Present() {
		t.Error("Expected opt1 to not be present")
	}

	if !opt2.Present() || *opt2.Get() != 123 {
		t.Error("opt2 did not set the expected value")
	}
}

func TestOpt_OfUnix(t *testing.T) {
	unixTime := int64(1633506600)
	date := time.Unix(unixTime, 0)
	o := uopt.OfUnix(unixTime)

	if !o.Present() || !reflect.DeepEqual(*o.Get(), date) {
		t.Error("OfUnix did not set the expected value")
	}
}

func TestOpt_OfBuilder(t *testing.T) {
	o := uopt.OfBuilder(func() int {
		return 123
	})

	if !o.Present() || *o.Get() != 123 {
		t.Error("OfBuilder did not set the expected value")
	}
}

func TestOpt_OrElse(t *testing.T) {
	opt1 := uopt.Null[int]()
	opt2 := uopt.Of(123)

	if opt1.OrElse(456) != 456 {
		t.Error("OrElse did not return the expected value for opt1")
	}

	if opt2.OrElse(456) != 123 {
		t.Error("OrElse did not return the expected value for opt2")
	}
}

func TestOpt_Get(t *testing.T) {
	o := uopt.Of(123)
	if *o.Get() != 123 {
		t.Error("Get did not return the expected value")
	}
}

func TestOpt_Set(t *testing.T) {
	o := uopt.Of(123)
	newVal := 456
	o.Set(&newVal)
	if *o.Get() != 456 {
		t.Error("Set did not update the value as expected")
	}
}

func TestOpt_GetAs(t *testing.T) {
	o := uopt.Of(123)
	result := o.GetAs(func(t int) any {
		return t + 1
	})

	if result.(int) != 124 {
		t.Error("GetAs did not return the expected value")
	}
}

func TestOpt_UnmarshalJSON(t *testing.T) {
	// Test for integer type
	optInt := uopt.Opt[int]{}
	err := json.Unmarshal([]byte("123"), &optInt)
	if err != nil || !optInt.Present() || *optInt.Get() != 123 {
		t.Errorf("UnmarshalJSON failed for integer type: %v", err)
	}

	// Test for null value
	optInt = uopt.Opt[int]{}
	err = json.Unmarshal([]byte("null"), &optInt)
	if err != nil || optInt.Present() {
		t.Errorf("UnmarshalJSON failed for null value: %v", err)
	}

	// Test for string type
	optStr := uopt.Opt[string]{}
	err = json.Unmarshal([]byte(`"hello"`), &optStr)
	if err != nil || !optStr.Present() || *optStr.Get() != "hello" {
		t.Errorf("UnmarshalJSON failed for string type: %v", err)
	}
}

func TestOpt_MarshalJSON(t *testing.T) {
	// Test for integer type
	optInt := uopt.Of(123)
	data, err := json.Marshal(optInt)
	if err != nil || string(data) != "123" {
		t.Errorf("MarshalJSON failed for integer type: %v", err)
	}

	// Test for null value
	optInt = uopt.Null[int]()
	data, err = json.Marshal(optInt)
	if err != nil || string(data) != "null" {
		t.Errorf("MarshalJSON failed for 0 value: %v", err)
	}

	// Test for string type
	optStr := uopt.Of("hello")
	data, err = json.Marshal(optStr)
	if err != nil || string(data) != `"hello"` {
		t.Errorf("MarshalJSON failed for string type: %v", err)
	}
}

// Test for driver.Valuer type
type customType struct{}

func (c customType) Value() (driver.Value, error) {
	return "customValue", nil
}

func TestOpt_Value(t *testing.T) {
	// Test for integer type
	optInt := uopt.Of(123)
	val, err := optInt.Value()
	if err != nil || val != int64(123) {
		t.Errorf("Value method failed for integer type: %v", err)
	}

	// Test for null value
	optInt = uopt.Null[int]()
	val, err = optInt.Value()
	if err != nil || val != nil {
		t.Errorf("Value method failed for null value: %v", err)
	}

	// Test for time type
	optTime := uopt.Of(time.Now())
	val, err = optTime.Value()
	if err != nil || val.(time.Time).Unix() != optTime.Get().Unix() {
		t.Errorf("Value method failed for time type: %v", err)
	}

	// Test for int8 type
	optInt8 := uopt.Of(int8(123))
	val, err = optInt8.Value()
	if err != nil || val != int64(123) {
		t.Errorf("Value method failed for int8 type: %v", err)
	}

	// Test for int16 type
	optInt16 := uopt.Of(int16(12345))
	val, err = optInt16.Value()
	if err != nil || val != int64(12345) {
		t.Errorf("Value method failed for int16 type: %v", err)
	}

	// Test for int32 type
	optInt32 := uopt.Of(int32(123456789))
	val, err = optInt32.Value()
	if err != nil || val != int64(123456789) {
		t.Errorf("Value method failed for int32 type: %v", err)
	}

	// Test for uint type
	optUint := uopt.Of(uint(123))
	val, err = optUint.Value()
	if err != nil || val != int64(123) {
		t.Errorf("Value method failed for uint type: %v", err)
	}

	// Test for uint8 type
	optUint8 := uopt.Of(uint8(123))
	val, err = optUint8.Value()
	if err != nil || val != int64(123) {
		t.Errorf("Value method failed for uint8 type: %v", err)
	}

	// Test for uint16 type
	optUint16 := uopt.Of(uint16(12345))
	val, err = optUint16.Value()
	if err != nil || val != int64(12345) {
		t.Errorf("Value method failed for uint16 type: %v", err)
	}

	// Test for uint32 type
	optUint32 := uopt.Of(uint32(123456789))
	val, err = optUint32.Value()
	if err != nil || val != int64(123456789) {
		t.Errorf("Value method failed for uint32 type: %v", err)
	}

	// Test for uint64 type
	optUint64 := uopt.Of(uint64(1234567890123456789))
	val, err = optUint64.Value()
	if err != nil || val != int64(1234567890123456789) {
		t.Errorf("Value method failed for uint64 type: %v", err)
	}

	// Test for float32 type
	optFloat32 := uopt.Of(float32(123.45))
	val, err = optFloat32.Value()
	if err != nil || math.Abs(val.(float64)-123.4) <= 1e-16 {
		t.Errorf("Value method failed for float32 type: %v", err)
	}

	// Test for float64 type
	optFloat64 := uopt.Of(123.45)
	val, err = optFloat64.Value()
	if err != nil || val != 123.45 {
		t.Errorf("Value method failed for float64 type: %v", err)
	}

	optValuer := uopt.Of(customType{})
	val, err = optValuer.Value()
	if err != nil || val != "customValue" {
		t.Errorf("Value method failed for driver.Valuer type: %v", err)
	}

	// Test for generic type not explicitly handled
	optString := uopt.Of("testString")
	val, err = optString.Value()
	if err != nil || val != "testString" {
		t.Errorf("Value method failed for generic type: %v", err)
	}
}

func TestOpt_Scan(t *testing.T) {
	// Test for integer type
	optInt := uopt.Opt[int]{}
	err := optInt.Scan(int64(123))
	if err != nil || !optInt.Present() || *optInt.Get() != 123 {
		t.Errorf("Scan method failed for integer type: %v", err)
	}

	// Test for null value
	optInt = uopt.Opt[int]{}
	err = optInt.Scan(nil)
	if err != nil || optInt.Present() {
		t.Errorf("Scan method failed for null value: %v", err)
	}

	// Test for string type
	optStr := uopt.Opt[string]{}
	err = optStr.Scan("hello")
	if err != nil || !optStr.Present() || *optStr.Get() != "hello" {
		t.Errorf("Scan method failed for string type: %v", err)
	}

	// Test for byte slice type
	optStr = uopt.Opt[string]{}
	err = optStr.Scan([]byte("hello"))
	if err != nil || !optStr.Present() || *optStr.Get() != "hello" {
		t.Errorf("Scan method failed for byte slice type: %v", err)
	}

	// Test for int8 type
	optInt8 := uopt.Opt[int8]{}
	err = optInt8.Scan(int64(12))
	if err != nil || !optInt8.Present() || *optInt8.Get() != 12 {
		t.Errorf("Scan method failed for int8 type: %v", err)
	}

	// Test for int16 type
	optInt16 := uopt.Opt[int16]{}
	err = optInt16.Scan(int64(1234))
	if err != nil || !optInt16.Present() || *optInt16.Get() != 1234 {
		t.Errorf("Scan method failed for int16 type: %v", err)
	}

	// Test for int32 type
	optInt32 := uopt.Opt[int32]{}
	err = optInt32.Scan(int64(12345678))
	if err != nil || !optInt32.Present() || *optInt32.Get() != 12345678 {
		t.Errorf("Scan method failed for int32 type: %v", err)
	}

	// Test for int64 type
	optInt64 := uopt.Opt[int64]{}
	err = optInt64.Scan(int64(1234567890123456))
	if err != nil || !optInt64.Present() || *optInt64.Get() != 1234567890123456 {
		t.Errorf("Scan method failed for int64 type: %v", err)
	}

	// Test for uint8 type
	optUint8 := uopt.Opt[uint8]{}
	err = optUint8.Scan(int64(25))
	if err != nil || !optUint8.Present() || *optUint8.Get() != 25 {
		t.Errorf("Scan method failed for uint8 type: %v", err)
	}

	// Test for uint16 type
	optUint16 := uopt.Opt[uint16]{}
	err = optUint16.Scan(int64(12345))
	if err != nil || !optUint16.Present() || *optUint16.Get() != 12345 {
		t.Errorf("Scan method failed for uint16 type: %v", err)
	}

	// Test for uint32 type
	optUint32 := uopt.Opt[uint32]{}
	err = optUint32.Scan(int64(1234567890))
	if err != nil || !optUint32.Present() || *optUint32.Get() != 1234567890 {
		t.Errorf("Scan method failed for uint32 type: %v", err)
	}

	// Test for uint32 type
	optUint64 := uopt.Opt[uint64]{}
	err = optUint64.Scan(int64(1234567890))
	if err != nil || !optUint64.Present() || *optUint64.Get() != 1234567890 {
		t.Errorf("Scan method failed for uint64 type: %v", err)
	}

	// Test for float32 type
	optFloat32 := uopt.Opt[float32]{}
	err = optFloat32.Scan(float64(1.234))
	if err != nil || !optFloat32.Present() || *optFloat32.Get() != float32(1.234) {
		t.Errorf("Scan method failed for float32 type: %v", err)
	}

	// Test for float64 type
	optFloat64 := uopt.Opt[float64]{}
	err = optFloat64.Scan(1.2345678)
	if err != nil || !optFloat64.Present() || *optFloat64.Get() != 1.2345678 {
		t.Errorf("Scan method failed for float64 type: %v", err)
	}

	// Test for int8 type
	optInt8Str := uopt.Opt[int8]{}
	err = optInt8Str.Scan("12")
	if err != nil || !optInt8Str.Present() || *optInt8Str.Get() != 12 {
		t.Errorf("Scan method failed for int8 type with string input: %v", err)
	}

	// Test for int16 type
	optInt16Str := uopt.Opt[int16]{}
	err = optInt16Str.Scan("1234")
	if err != nil || !optInt16Str.Present() || *optInt16Str.Get() != 1234 {
		t.Errorf("Scan method failed for int16 type with string input: %v", err)
	}

	// Test for int32 type
	optInt32Str := uopt.Opt[int32]{}
	err = optInt32Str.Scan("12345678")
	if err != nil || !optInt32Str.Present() || *optInt32Str.Get() != 12345678 {
		t.Errorf("Scan method failed for int32 type with string input: %v", err)
	}

	// Test for int64 type
	optInt64Str := uopt.Opt[int64]{}
	err = optInt64Str.Scan("1234567890123456")
	if err != nil || !optInt64Str.Present() || *optInt64Str.Get() != 1234567890123456 {
		t.Errorf("Scan method failed for int64 type with string input: %v", err)
	}

	// Test for uint8 type
	optUint8Str := uopt.Opt[uint8]{}
	err = optUint8Str.Scan("25")
	if err != nil || !optUint8Str.Present() || *optUint8Str.Get() != 25 {
		t.Errorf("Scan method failed for uint8 type with string input: %v", err)
	}

	// Test for uint16 type
	optUint16Str := uopt.Opt[uint16]{}
	err = optUint16Str.Scan("12345")
	if err != nil || !optUint16Str.Present() || *optUint16Str.Get() != 12345 {
		t.Errorf("Scan method failed for uint16 type with string input: %v", err)
	}

	// Test for uint32 type
	optUint32Str := uopt.Opt[uint32]{}
	err = optUint32Str.Scan("1234567890")
	if err != nil || !optUint32Str.Present() || *optUint32Str.Get() != 1234567890 {
		t.Errorf("Scan method failed for uint32 type with string input: %v", err)
	}

	// Test for uint64 type
	optUint64Str := uopt.Opt[uint64]{}
	err = optUint64Str.Scan("1234567890123456")
	if err != nil || !optUint64Str.Present() || *optUint64Str.Get() != 1234567890123456 {
		t.Errorf("Scan method failed for uint64 type with string input: %v", err)
	}

	// Test for float32 type
	optFloat32Str := uopt.Opt[float32]{}
	err = optFloat32Str.Scan("3.14")
	if err != nil || !optFloat32Str.Present() || *optFloat32Str.Get() != float32(3.14) {
		t.Errorf("Scan method failed for float32 type with string input: %v", err)
	}

	// Test for float64 type
	optFloat64Str := uopt.Opt[float64]{}
	err = optFloat64Str.Scan("3.1415926535")
	if err != nil || !optFloat64Str.Present() || *optFloat64Str.Get() != 3.1415926535 {
		t.Errorf("Scan method failed for float64 type with string input: %v", err)
	}

	// Test for bool type
	optBool := uopt.Opt[bool]{}
	err = optBool.Scan("true")
	if err != nil || !optBool.Present() || *optBool.Get() != true {
		t.Errorf("Scan method failed for bool type with string input: %v", err)
	}

	// Test for complex64 type
	optComplex64 := uopt.Opt[complex64]{}
	err = optComplex64.Scan("3+4i")
	if err != nil || !optComplex64.Present() || *optComplex64.Get() != complex(float32(3), float32(4)) {
		t.Errorf("Scan method failed for complex64 type with string input: %v", err)
	}

	// Test for complex128 type
	optComplex128 := uopt.Opt[complex128]{}
	err = optComplex128.Scan("3.14+4.15i")
	if err != nil || !optComplex128.Present() || *optComplex128.Get() != complex(3.14, 4.15) {
		t.Errorf("Scan method failed for complex128 type with string input: %v", err)
	}

	// Test for bool type
	optBoolInt := uopt.Opt[bool]{}
	err = optBoolInt.Scan(int64(1))
	if err != nil || !optBoolInt.Present() || *optBoolInt.Get() != true {
		t.Errorf("Scan method failed for bool type with string input: %v", err)
	}

	// Test for string type
	optStrFloat32 := uopt.Opt[string]{}
	err = optStrFloat32.Scan(float32(3.1415926535))
	if err != nil || !optStrFloat32.Present() || *optStrFloat32.Get() != "3.1415927" {
		t.Errorf("Scan method failed for string type with string input: %v", err)
	}

	// Test for float32 type
	optFloat32Float := uopt.Opt[float32]{}
	err = optFloat32Float.Scan(float32(3.1415926535))
	if err != nil || !optFloat32Float.Present() || *optFloat32Float.Get() != 3.1415926535 {
		t.Errorf("Scan method failed for float32 type with float32 input: %v", err)
	}

	// Test for float32 type
	optFloat64Float := uopt.Opt[float64]{}
	err = optFloat64Float.Scan(3.1415926535)
	if err != nil || !optFloat64Float.Present() || *optFloat64Float.Get() != 3.1415926535 {
		t.Errorf("Scan method failed for float64 type with float32 input: %v", err)
	}

	// Test for string type
	optStrFloat64 := uopt.Opt[string]{}
	err = optStrFloat64.Scan(float32(3.1415926535))
	if err != nil || !optStrFloat64.Present() || *optStrFloat64.Get() != "3.1415927" {
		t.Errorf("Scan method failed for string type with float32 input: %v", err)
	}

	// Test for float64 type
	optFloat64Float32 := uopt.Opt[float64]{}
	err = optFloat64Float32.Scan(float32(3.1415926535))
	if err != nil || !optFloat64Float32.Present() || math.Abs(*optFloat64Float32.Get()-3.1415926535) <= 1e-16 {
		t.Errorf("Scan method failed for float64 type with float64 input: %v", err)
	}
}

// mockValuer is a mock type that implements the driver.Valuer interface.
type mockValuer struct {
	val interface{}
}

// Value implements the driver.Valuer interface for mockValuer.
func (m mockValuer) Value() (driver.Value, error) {
	return m.val, nil
}

func TestOpt_ValueDriverValuer(t *testing.T) {
	expectedValue := "test string"
	m := mockValuer{val: expectedValue}
	o := uopt.Of(m)

	val, err := o.Value()
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, val)

	now := time.Now()
	oTime := uopt.Of[time.Time](now)

	val, err = oTime.Value()
	assert.NoError(t, err)
	assert.Equal(t, now, val)
}

func TestOpt_Scan_Errors(t *testing.T) {
	// Case: Incompatible type for Opt, failed to retrieve value
	t.Run("Incompatible type retrieval", func(t *testing.T) {
		var o uopt.Opt[int]
		err := o.Scan([]byte("invalid"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse bytes/blob sql value to numeric opt: strconv.ParseInt: parsing \"invalid\": invalid syntax")

		// Parse invalid value for Opt[int8]
		var o8 uopt.Opt[int8]
		err = o8.Scan([]byte("invalid"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse bytes/blob sql value to numeric opt: strconv.ParseInt: parsing \"invalid\": invalid syntax")

		// Parse invalid value for Opt[int16]
		var o16 uopt.Opt[int16]
		err = o16.Scan([]byte("invalid"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse bytes/blob sql value to numeric opt: strconv.ParseInt: parsing \"invalid\": invalid syntax")

		// Parse invalid value for Opt[int32]
		var o32 uopt.Opt[int32]
		err = o32.Scan([]byte("invalid"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse bytes/blob sql value to numeric opt: strconv.ParseInt: parsing \"invalid\": invalid syntax")

		// Parse invalid value for Opt[int64]
		var o64 uopt.Opt[int64]
		err = o64.Scan([]byte("invalid"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse bytes/blob sql value to numeric opt: strconv.ParseInt: parsing \"invalid\": invalid syntax")

		// Parse invalid value for Opt[uint]
		var ou uopt.Opt[uint]
		err = ou.Scan([]byte("invalid"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse bytes/blob sql value to numeric opt: strconv.ParseUint: parsing \"invalid\": invalid syntax")

		// Parse invalid value for Opt[uint8]
		var ou8 uopt.Opt[uint8]
		err = ou8.Scan([]byte("invalid"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse bytes/blob sql value to numeric opt: strconv.ParseUint: parsing \"invalid\": invalid syntax")

		// Parse invalid value for Opt[uint16]
		var ou16 uopt.Opt[uint16]
		err = ou16.Scan([]byte("invalid"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse bytes/blob sql value to numeric opt: strconv.ParseUint: parsing \"invalid\": invalid syntax")

		// Parse invalid value for Opt[uint32]
		var ou32 uopt.Opt[uint32]
		err = ou32.Scan([]byte("invalid"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse bytes/blob sql value to numeric opt: strconv.ParseUint: parsing \"invalid\": invalid syntax")

		// Parse invalid value for Opt[uint64]
		var ou64 uopt.Opt[uint64]
		err = ou64.Scan([]byte("invalid"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse bytes/blob sql value to numeric opt: strconv.ParseUint: parsing \"invalid\": invalid syntax")

		// Parse invalid value for Opt[float32]
		var of32 uopt.Opt[float32]
		err = of32.Scan([]byte("invalid"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse bytes/blob sql value to float opt: strconv.ParseFloat: parsing \"invalid\": invalid syntax")

		// Parse invalid value for Opt[float64]
		var of64 uopt.Opt[float64]
		err = of64.Scan([]byte("invalid"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse bytes/blob sql value to float opt: strconv.ParseFloat: parsing \"invalid\": invalid syntax")
	})

	// Case: Failed to parse varchar sql value to numeric opt
	t.Run("Failed parse varchar to numeric", func(t *testing.T) {
		var o uopt.Opt[int]
		err := o.Scan("not-a-number")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse varchar sql value to numeric opt")
	})

	// Case: Failed to parse varchar sql value to bool opt
	t.Run("Failed parse varchar to bool", func(t *testing.T) {
		var o uopt.Opt[bool]
		err := o.Scan("not-a-bool")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse varchar sql value to bool opt")
	})
}
