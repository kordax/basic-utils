/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uos_test

import (
	"math"
	"os"
	"strconv"
	"testing"

	"github.com/kordax/basic-utils/uos"
	"github.com/stretchr/testify/assert"
)

func TestGetCPUs_Stub(t *testing.T) {
	assert.NotZero(t, uos.GetCPUs())
}

func TestGetEnvNumeric(t *testing.T) {
	// Setting environment variables with boundary values
	os.Setenv("TEST_INT", strconv.Itoa(math.MaxInt))
	os.Setenv("TEST_INT8", strconv.Itoa(math.MinInt8))
	os.Setenv("TEST_INT16", strconv.Itoa(math.MinInt16))
	os.Setenv("TEST_INT32", strconv.Itoa(math.MaxInt32))
	os.Setenv("TEST_INT64", strconv.FormatInt(math.MinInt64, 10))
	os.Setenv("TEST_UINT", strconv.FormatUint(math.MaxUint, 10))
	os.Setenv("TEST_UINT8", strconv.Itoa(math.MaxUint8))
	os.Setenv("TEST_UINT16", strconv.Itoa(math.MaxUint16))
	os.Setenv("TEST_UINT32", strconv.FormatUint(math.MaxUint32, 10))
	os.Setenv("TEST_UINT64", strconv.FormatUint(math.MaxUint64, 10))
	os.Setenv("TEST_FLOAT32", strconv.FormatFloat(math.MaxFloat32, 'f', -1, 64))
	os.Setenv("TEST_FLOAT64", strconv.FormatFloat(-math.MaxFloat64, 'f', -1, 64))

	tests := []struct {
		name    string
		key     string
		want    any
		wantErr bool
	}{
		{"Int", "TEST_INT", math.MaxInt, false},
		{"Int8", "TEST_INT8", int8(math.MinInt8), false},
		{"Int16", "TEST_INT16", int16(math.MinInt16), false},
		{"Int32", "TEST_INT32", int32(math.MaxInt32), false},
		{"Int64", "TEST_INT64", int64(math.MinInt64), false},
		{"Uint", "TEST_UINT", uint(math.MaxUint), false},
		{"Uint8", "TEST_UINT8", uint8(math.MaxUint8), false},
		{"Uint16", "TEST_UINT16", uint16(math.MaxUint16), false},
		{"Uint32", "TEST_UINT32", uint32(math.MaxUint32), false},
		{"Uint64", "TEST_UINT64", uint64(math.MaxUint64), false},
		{"Float32", "TEST_FLOAT32", float32(math.MaxFloat32), false},
		{"Float64", "TEST_FLOAT64", -math.MaxFloat64, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantErr {
						t.Errorf("GetEnvNumeric() for %s panicked unexpectedly: %v", tt.key, r)
					}
				}
			}()

			var result any
			switch tt.name {
			case "Int":
				result = uos.GetEnvNumeric[int](tt.key)
			case "Int8":
				result = uos.GetEnvNumeric[int8](tt.key)
			case "Int16":
				result = uos.GetEnvNumeric[int16](tt.key)
			case "Int32":
				result = uos.GetEnvNumeric[int32](tt.key)
			case "Int64":
				result = uos.GetEnvNumeric[int64](tt.key)
			case "Uint":
				result = uos.GetEnvNumeric[uint](tt.key)
			case "Uint8":
				result = uos.GetEnvNumeric[uint8](tt.key)
			case "Uint16":
				result = uos.GetEnvNumeric[uint16](tt.key)
			case "Uint32":
				result = uos.GetEnvNumeric[uint32](tt.key)
			case "Uint64":
				result = uos.GetEnvNumeric[uint64](tt.key)
			case "Float32":
				result = uos.GetEnvNumeric[float32](tt.key)
			case "Float64":
				result = uos.GetEnvNumeric[float64](tt.key)
			default:
				t.Fatalf("Unhandled type in test cases")
			}

			if !tt.wantErr {
				assert.Equal(t, tt.want, result, "GetEnvNumeric() returned incorrect value for %s", tt.key)
			}
		})
	}
}

func TestGetEnvNumeric_Panic(t *testing.T) {
	os.Setenv("TEST_INVALID_INT", "invalid_int")     // Invalid int value
	os.Setenv("TEST_OVERFLOW_INT", "2147483648")     // Overflow int32 value
	os.Setenv("TEST_INVALID_FLOAT", "invalid_float") // Invalid float value

	tests := []struct {
		name    string
		key     string
		want    any
		wantErr bool
	}{
		// Previous positive test cases...

		// Negative test cases
		{"InvalidInt", "TEST_INVALID_INT", nil, true},
		{"OverflowInt", "TEST_OVERFLOW_INT", nil, true},
		{"InvalidFloat", "TEST_INVALID_FLOAT", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() {
				var result any
				switch tt.name {
				case "InvalidInt":
					result = uos.GetEnvNumeric[int](tt.key)
				case "OverflowInt":
					result = uos.GetEnvNumeric[int32](tt.key)
				case "InvalidFloat":
					result = uos.GetEnvNumeric[float64](tt.key)
				default:
					result = uos.GetEnvNumeric[int](tt.key) // Default case
				}
				if !tt.wantErr {
					assert.Equal(t, tt.want, result, "GetEnvNumeric() returned incorrect value for %s", tt.key)
				}
			}, "Expected panic for invalid input in GetEnvNumeric()")
		})
	}
}
