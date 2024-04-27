/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uos_test

import (
	"encoding/base64"
	"encoding/hex"
	"math"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/kordax/basic-utils/uos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCPUs_Stub(t *testing.T) {
	cpus := uos.GetCPUs()
	assert.NotZero(t, cpus)
	assert.Less(t, cpus, 100)

	for i := 0; i < 100; i++ {
		next := uos.GetCPUs()
		assert.NotZero(t, next)
		assert.Less(t, next, 100)
		assert.Equal(t, cpus, next)
	}
}

func TestRequireEnv(t *testing.T) {
	key := "TESTENVVALUE"
	value := "someRandomValue12431!@#$!@#%^^"
	require.NoError(t, os.Setenv(key, value))
	defer func() {
		_ = os.Unsetenv("TESTENVVALUE")
	}()

	assert.EqualValues(t, value, uos.RequireEnv(key))
	_ = os.Unsetenv("TESTENVVALUE")
	assert.Panics(t, func() {
		uos.RequireEnv(key)
	})
}

func TestRequireEnvNumeric(t *testing.T) {
	require.NoError(t, os.Setenv("TEST_INT", strconv.Itoa(math.MaxInt)))
	require.NoError(t, os.Setenv("TEST_INT8", strconv.Itoa(math.MinInt8)))
	require.NoError(t, os.Setenv("TEST_INT16", strconv.Itoa(math.MinInt16)))
	require.NoError(t, os.Setenv("TEST_INT32", strconv.Itoa(math.MaxInt32)))
	require.NoError(t, os.Setenv("TEST_INT64", strconv.FormatInt(math.MinInt64, 10)))
	require.NoError(t, os.Setenv("TEST_UINT", strconv.FormatUint(math.MaxUint, 10)))
	require.NoError(t, os.Setenv("TEST_UINT8", strconv.Itoa(math.MaxUint8)))
	require.NoError(t, os.Setenv("TEST_UINT16", strconv.Itoa(math.MaxUint16)))
	require.NoError(t, os.Setenv("TEST_UINT32", strconv.FormatUint(math.MaxUint32, 10)))
	require.NoError(t, os.Setenv("TEST_UINT64", strconv.FormatUint(math.MaxUint64, 10)))
	require.NoError(t, os.Setenv("TEST_FLOAT32", strconv.FormatFloat(math.MaxFloat32, 'f', -1, 64)))
	require.NoError(t, os.Setenv("TEST_FLOAT64", strconv.FormatFloat(-math.MaxFloat64, 'f', -1, 64)))

	defer func() {
		_ = os.Unsetenv("TEST_INT")
		_ = os.Unsetenv("TEST_INT8")
		_ = os.Unsetenv("TEST_INT16")
		_ = os.Unsetenv("TEST_INT32")
		_ = os.Unsetenv("TEST_INT64")
		_ = os.Unsetenv("TEST_UINT")
		_ = os.Unsetenv("TEST_UINT8")
		_ = os.Unsetenv("TEST_UINT16")
		_ = os.Unsetenv("TEST_UINT32")
		_ = os.Unsetenv("TEST_UINT64")
		_ = os.Unsetenv("TEST_FLOAT32")
		_ = os.Unsetenv("TEST_FLOAT64")
	}()

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
						t.Errorf("RequireEnvNumeric() for %s panicked unexpectedly: %v", tt.key, r)
					}
				}
			}()

			var result any
			switch tt.name {
			case "Int":
				result = uos.RequireEnvNumeric[int](tt.key)
			case "Int8":
				result = uos.RequireEnvNumeric[int8](tt.key)
			case "Int16":
				result = uos.RequireEnvNumeric[int16](tt.key)
			case "Int32":
				result = uos.RequireEnvNumeric[int32](tt.key)
			case "Int64":
				result = uos.RequireEnvNumeric[int64](tt.key)
			case "Uint":
				result = uos.RequireEnvNumeric[uint](tt.key)
			case "Uint8":
				result = uos.RequireEnvNumeric[uint8](tt.key)
			case "Uint16":
				result = uos.RequireEnvNumeric[uint16](tt.key)
			case "Uint32":
				result = uos.RequireEnvNumeric[uint32](tt.key)
			case "Uint64":
				result = uos.RequireEnvNumeric[uint64](tt.key)
			case "Float32":
				result = uos.RequireEnvNumeric[float32](tt.key)
			case "Float64":
				result = uos.RequireEnvNumeric[float64](tt.key)
			default:
				t.Fatalf("Unhandled type in test cases")
			}

			if !tt.wantErr {
				assert.Equal(t, tt.want, result, "RequireEnvNumeric() returned incorrect value for %s", tt.key)
			}
		})
	}
}

func TestRequireEnvNumeric_Panic(t *testing.T) {
	require.NoError(t, os.Setenv("TEST_INVALID_INT", "invalid_int"))     // Invalid int value
	require.NoError(t, os.Setenv("TEST_OVERFLOW_INT", "2147483648"))     // Overflow int32 value
	require.NoError(t, os.Setenv("TEST_INVALID_FLOAT", "invalid_float")) // Invalid float value

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
					result = uos.RequireEnvNumeric[int](tt.key)
				case "OverflowInt":
					result = uos.RequireEnvNumeric[int32](tt.key)
				case "InvalidFloat":
					result = uos.RequireEnvNumeric[float64](tt.key)
				default:
					result = uos.RequireEnvNumeric[int](tt.key) // Default case
				}
				if !tt.wantErr {
					assert.Equal(t, tt.want, result, "RequireEnvNumeric() returned incorrect value for %s", tt.key)
				}
			}, "Expected panic for invalid input in RequireEnvNumeric()")
		})
	}
}

func TestRequireEnvAs(t *testing.T) {
	require.NoError(t, os.Setenv("TEST_TIME", "2023-01-01T15:04:05Z"))
	require.NoError(t, os.Setenv("TEST_BASE64", base64.StdEncoding.EncodeToString([]byte("hello"))))
	require.NoError(t, os.Setenv("TEST_HEX", hex.EncodeToString([]byte("hello"))))
	require.NoError(t, os.Setenv("TEST_URL", "https://www.example.com"))

	defer t.Cleanup(func() {
		_ = os.Unsetenv("TEST_TIME")
		_ = os.Unsetenv("TEST_BASE64")
		_ = os.Unsetenv("TEST_HEX")
		_ = os.Unsetenv("TEST_URL")
	})

	t.Run("Time", func(t *testing.T) {
		expectedTime, _ := time.Parse(time.RFC3339, os.Getenv("TEST_TIME"))
		result := uos.RequireEnvAs("TEST_TIME", uos.MapStringToTime(time.RFC3339))
		assert.Equal(t, expectedTime, result)
	})

	t.Run("Base64", func(t *testing.T) {
		expectedText := "hello"
		result := uos.RequireEnvAs("TEST_BASE64", uos.MapStringToBase64)
		assert.Equal(t, expectedText, result)
	})

	t.Run("Hex", func(t *testing.T) {
		expectedText := []byte("hello")
		result := uos.RequireEnvAs("TEST_HEX", uos.MapStringToHex)
		assert.Equal(t, expectedText, result)
	})

	t.Run("URL", func(t *testing.T) {
		expectedURL, _ := url.Parse(os.Getenv("TEST_URL"))
		result := uos.RequireEnvAs("TEST_URL", uos.MapStringToURL)
		assert.Equal(t, *expectedURL, result)
	})
}

func TestRequireEnvHelpers(t *testing.T) {
	require.NoError(t, os.Setenv("TEST_DURATION", "2h45m"))
	require.NoError(t, os.Setenv("TEST_TIME", "2023-01-02T15:04:05Z"))
	require.NoError(t, os.Setenv("TEST_URL", "https://www.example.com"))
	require.NoError(t, os.Setenv("TEST_BOOL", "true"))

	defer t.Cleanup(func() {
		_ = os.Unsetenv("TEST_DURATION")
		_ = os.Unsetenv("TEST_TIME")
		_ = os.Unsetenv("TEST_URL")
		_ = os.Unsetenv("TEST_BOOL")
	})

	t.Run("Duration", func(t *testing.T) {
		expectedDuration, _ := time.ParseDuration(os.Getenv("TEST_DURATION"))
		result := uos.RequireEnvDuration("TEST_DURATION")
		assert.Equal(t, expectedDuration, result)
	})

	t.Run("Time", func(t *testing.T) {
		expectedTime, _ := time.Parse(time.RFC3339, os.Getenv("TEST_TIME"))
		result := uos.RequireEnvTime("TEST_TIME", time.RFC3339)
		assert.Equal(t, expectedTime, result)
	})

	t.Run("URL", func(t *testing.T) {
		expectedURL, _ := url.Parse(os.Getenv("TEST_URL"))
		result := uos.RequireEnvURL("TEST_URL")
		assert.Equal(t, *expectedURL, result)
	})

	t.Run("Bool", func(t *testing.T) {
		expectedBool, _ := strconv.ParseBool(os.Getenv("TEST_BOOL"))
		result := uos.RequireEnvBool("TEST_BOOL")
		assert.Equal(t, expectedBool, result)
	})
}
