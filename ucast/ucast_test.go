package ucast_test

import (
	"testing"

	"github.com/kordax/basic-utils/ucast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		wantErr  bool
	}{
		// Integers
		{"Int", "123", int(123), false},
		{"IntPtr", "123", new(int), false},
		{"Int8", "127", int8(127), false},
		{"Int8Ptr", "-128", new(int8), false},
		{"Int16", "32767", int16(32767), false},
		{"Int16Ptr", "-32768", new(int16), false},
		{"Int32", "2147483647", int32(2147483647), false},
		{"Int32Ptr", "-2147483648", new(int32), false},
		{"Int64", "9223372036854775807", int64(9223372036854775807), false},
		{"Int64Ptr", "-9223372036854775808", new(int64), false},
		{"InvalidInt", "invalid", int(0), true},

		// Unsigned Integers
		{"Uint", "123", uint(123), false},
		{"UintPtr", "123", new(uint), false},
		{"Uint8", "255", uint8(255), false},
		{"Uint8Ptr", "0", new(uint8), false},
		{"Uint16", "65535", uint16(65535), false},
		{"Uint16Ptr", "0", new(uint16), false},
		{"Uint32", "4294967295", uint32(4294967295), false},
		{"Uint32Ptr", "0", new(uint32), false},
		{"Uint64", "18446744073709551615", uint64(18446744073709551615), false},
		{"Uint64Ptr", "0", new(uint64), false},
		{"InvalidUint", "-1", uint(0), true},

		// Floats
		{"Float32", "3.14", float32(3.14), false},
		{"Float32Ptr", "-3.14", new(float32), false},
		{"Float64", "3.1415926535", float64(3.1415926535), false},
		{"Float64Ptr", "-3.1415926535", new(float64), false},
		{"InvalidFloat", "invalid", float64(0), true},

		// Booleans
		{"BoolTrue", "true", true, false},
		{"BoolFalse", "false", false, false},
		{"BoolPtr", "true", new(bool), false},
		{"InvalidBool", "invalid", false, true},

		// Strings
		{"String", "hello", "hello", false},
		{"StringPtr", "world", new(string), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch expected := tt.expected.(type) {
			case int:
				result, err := ucast.String[int](tt.input)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, expected, result)
				}
			case *int:
				result, err := ucast.String[*int](tt.input)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, result)
				} else {
					require.NoError(t, err)
					require.NotNil(t, result)
					val, _ := ucast.String[int](tt.input)
					assert.Equal(t, val, *result)
				}
			case int8:
				result, err := ucast.String[int8](tt.input)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, expected, result)
				}
			case *int8:
				result, err := ucast.String[*int8](tt.input)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, result)
				} else {
					require.NoError(t, err)
					require.NotNil(t, result)
					val, _ := ucast.String[int8](tt.input)
					assert.Equal(t, val, *result)
				}
			case int16:
				result, err := ucast.String[int16](tt.input)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, expected, result)
				}
			case *int16:
				result, err := ucast.String[*int16](tt.input)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, result)
				} else {
					require.NoError(t, err)
					require.NotNil(t, result)
					val, _ := ucast.String[int16](tt.input)
					assert.Equal(t, val, *result)
				}
			case int32:
				result, err := ucast.String[int32](tt.input)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, expected, result)
				}
			case *int32:
				result, err := ucast.String[*int32](tt.input)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, result)
				} else {
					require.NoError(t, err)
					require.NotNil(t, result)
					val, _ := ucast.String[int32](tt.input)
					assert.Equal(t, val, *result)
				}
			case int64:
				result, err := ucast.String[int64](tt.input)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, expected, result)
				}
			case *int64:
				result, err := ucast.String[*int64](tt.input)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, result)
				} else {
					require.NoError(t, err)
					require.NotNil(t, result)
					val, _ := ucast.String[int64](tt.input)
					assert.Equal(t, val, *result)
				}
			case uint:
				result, err := ucast.String[uint](tt.input)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, expected, result)
				}
			case *uint:
				result, err := ucast.String[*uint](tt.input)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, result)
				} else {
					require.NoError(t, err)
					require.NotNil(t, result)
					val, _ := ucast.String[uint](tt.input)
					assert.Equal(t, val, *result)
				}
			case uint8:
				result, err := ucast.String[uint8](tt.input)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, expected, result)
				}
			case *uint8:
				result, err := ucast.String[*uint8](tt.input)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, result)
				} else {
					require.NoError(t, err)
					require.NotNil(t, result)
					val, _ := ucast.String[uint8](tt.input)
					assert.Equal(t, val, *result)
				}
			case uint16:
				result, err := ucast.String[uint16](tt.input)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, expected, result)
				}
			case *uint16:
				result, err := ucast.String[*uint16](tt.input)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, result)
				} else {
					require.NoError(t, err)
					require.NotNil(t, result)
					val, _ := ucast.String[uint16](tt.input)
					assert.Equal(t, val, *result)
				}
			case uint32:
				result, err := ucast.String[uint32](tt.input)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, expected, result)
				}
			case *uint32:
				result, err := ucast.String[*uint32](tt.input)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, result)
				} else {
					require.NoError(t, err)
					require.NotNil(t, result)
					val, _ := ucast.String[uint32](tt.input)
					assert.Equal(t, val, *result)
				}
			case uint64:
				result, err := ucast.String[uint64](tt.input)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, expected, result)
				}
			case *uint64:
				result, err := ucast.String[*uint64](tt.input)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, result)
				} else {
					require.NoError(t, err)
					require.NotNil(t, result)
					val, _ := ucast.String[uint64](tt.input)
					assert.Equal(t, val, *result)
				}
			case float32:
				result, err := ucast.String[float32](tt.input)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.InDelta(t, expected, result, 0.0001)
				}
			case *float32:
				result, err := ucast.String[*float32](tt.input)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, result)
				} else {
					require.NoError(t, err)
					require.NotNil(t, result)
					val, _ := ucast.String[float32](tt.input)
					assert.InDelta(t, val, *result, 0.0001)
				}
			case float64:
				result, err := ucast.String[float64](tt.input)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, expected, result)
				}
			case *float64:
				result, err := ucast.String[*float64](tt.input)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, result)
				} else {
					require.NoError(t, err)
					require.NotNil(t, result)
					val, _ := ucast.String[float64](tt.input)
					assert.Equal(t, val, *result)
				}
			case bool:
				result, err := ucast.String[bool](tt.input)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, expected, result)
				}
			case *bool:
				result, err := ucast.String[*bool](tt.input)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, result)
				} else {
					require.NoError(t, err)
					require.NotNil(t, result)
					val, _ := ucast.String[bool](tt.input)
					assert.Equal(t, val, *result)
				}
			case string:
				result, err := ucast.String[string](tt.input)
				require.NoError(t, err)
				assert.Equal(t, expected, result)
			case *string:
				result, err := ucast.String[*string](tt.input)
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.input, *result)
			default:
				t.Fatalf("Unsupported type: %T", expected)
			}
		})
	}
}

func TestType(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		// Integers
		{"Int", int(123), "123"},
		{"IntPtr", func() *int { v := 123; return &v }(), "123"},
		{"IntNilPtr", (*int)(nil), ""},
		{"Int8", int8(127), "127"},
		{"Int8Ptr", func() *int8 { v := int8(127); return &v }(), "127"},
		{"Int8NilPtr", (*int8)(nil), ""},
		{"Int16", int16(32767), "32767"},
		{"Int16Ptr", func() *int16 { v := int16(32767); return &v }(), "32767"},
		{"Int16NilPtr", (*int16)(nil), ""},
		{"Int32", int32(2147483647), "2147483647"},
		{"Int32Ptr", func() *int32 { v := int32(2147483647); return &v }(), "2147483647"},
		{"Int32NilPtr", (*int32)(nil), ""},
		{"Int64", int64(9223372036854775807), "9223372036854775807"},
		{"Int64Ptr", func() *int64 { v := int64(9223372036854775807); return &v }(), "9223372036854775807"},
		{"Int64NilPtr", (*int64)(nil), ""},

		// Unsigned Integers
		{"Uint", uint(123), "123"},
		{"UintPtr", func() *uint { v := uint(123); return &v }(), "123"},
		{"UintNilPtr", (*uint)(nil), ""},
		{"Uint8", uint8(255), "255"},
		{"Uint8Ptr", func() *uint8 { v := uint8(255); return &v }(), "255"},
		{"Uint8NilPtr", (*uint8)(nil), ""},
		{"Uint16", uint16(65535), "65535"},
		{"Uint16Ptr", func() *uint16 { v := uint16(65535); return &v }(), "65535"},
		{"Uint16NilPtr", (*uint16)(nil), ""},
		{"Uint32", uint32(4294967295), "4294967295"},
		{"Uint32Ptr", func() *uint32 { v := uint32(4294967295); return &v }(), "4294967295"},
		{"Uint32NilPtr", (*uint32)(nil), ""},
		{"Uint64", uint64(18446744073709551615), "18446744073709551615"},
		{"Uint64Ptr", func() *uint64 { v := uint64(18446744073709551615); return &v }(), "18446744073709551615"},
		{"Uint64NilPtr", (*uint64)(nil), ""},

		// Floats
		{"Float32", float32(3.14), "3.14"},
		{"Float32Ptr", func() *float32 { v := float32(3.14); return &v }(), "3.14"},
		{"Float32NilPtr", (*float32)(nil), ""},
		{"Float64", float64(3.1415926535), "3.1415926535"},
		{"Float64Ptr", func() *float64 { v := float64(3.1415926535); return &v }(), "3.1415926535"},
		{"Float64NilPtr", (*float64)(nil), ""},

		// Booleans
		{"BoolTrue", true, "true"},
		{"BoolFalse", false, "false"},
		{"BoolPtr", func() *bool { v := true; return &v }(), "true"},
		{"BoolNilPtr", (*bool)(nil), ""},

		// Strings
		{"String", "hello", "hello"},
		{"StringPtr", func() *string { v := "world"; return &v }(), "world"},
		{"StringNilPtr", (*string)(nil), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch input := tt.input.(type) {
			case int:
				result := ucast.Type[int](input)
				assert.Equal(t, tt.expected, result)
			case *int:
				result := ucast.Type[*int](input)
				assert.Equal(t, tt.expected, result)
			case int8:
				result := ucast.Type[int8](input)
				assert.Equal(t, tt.expected, result)
			case *int8:
				result := ucast.Type[*int8](input)
				assert.Equal(t, tt.expected, result)
			case int16:
				result := ucast.Type[int16](input)
				assert.Equal(t, tt.expected, result)
			case *int16:
				result := ucast.Type[*int16](input)
				assert.Equal(t, tt.expected, result)
			case int32:
				result := ucast.Type[int32](input)
				assert.Equal(t, tt.expected, result)
			case *int32:
				result := ucast.Type[*int32](input)
				assert.Equal(t, tt.expected, result)
			case int64:
				result := ucast.Type[int64](input)
				assert.Equal(t, tt.expected, result)
			case *int64:
				result := ucast.Type[*int64](input)
				assert.Equal(t, tt.expected, result)
			case uint:
				result := ucast.Type[uint](input)
				assert.Equal(t, tt.expected, result)
			case *uint:
				result := ucast.Type[*uint](input)
				assert.Equal(t, tt.expected, result)
			case uint8:
				result := ucast.Type[uint8](input)
				assert.Equal(t, tt.expected, result)
			case *uint8:
				result := ucast.Type[*uint8](input)
				assert.Equal(t, tt.expected, result)
			case uint16:
				result := ucast.Type[uint16](input)
				assert.Equal(t, tt.expected, result)
			case *uint16:
				result := ucast.Type[*uint16](input)
				assert.Equal(t, tt.expected, result)
			case uint32:
				result := ucast.Type[uint32](input)
				assert.Equal(t, tt.expected, result)
			case *uint32:
				result := ucast.Type[*uint32](input)
				assert.Equal(t, tt.expected, result)
			case uint64:
				result := ucast.Type[uint64](input)
				assert.Equal(t, tt.expected, result)
			case *uint64:
				result := ucast.Type[*uint64](input)
				assert.Equal(t, tt.expected, result)
			case float32:
				result := ucast.Type[float32](input)
				assert.Equal(t, tt.expected, result)
			case *float32:
				result := ucast.Type[*float32](input)
				assert.Equal(t, tt.expected, result)
			case float64:
				result := ucast.Type[float64](input)
				assert.Equal(t, tt.expected, result)
			case *float64:
				result := ucast.Type[*float64](input)
				assert.Equal(t, tt.expected, result)
			case bool:
				result := ucast.Type[bool](input)
				assert.Equal(t, tt.expected, result)
			case *bool:
				result := ucast.Type[*bool](input)
				assert.Equal(t, tt.expected, result)
			case string:
				result := ucast.Type[string](input)
				assert.Equal(t, tt.expected, result)
			case *string:
				result := ucast.Type[*string](input)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestStringOrDef(t *testing.T) {
	t.Run("ValidInt", func(t *testing.T) {
		result, err := ucast.StringOrDef[int]("123", 0)
		require.NoError(t, err)
		assert.Equal(t, 123, result)
	})

	t.Run("InvalidIntWithDefault", func(t *testing.T) {
		result, err := ucast.StringOrDef[int]("invalid", 42)
		assert.Error(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("ValidFloat64", func(t *testing.T) {
		result, err := ucast.StringOrDef[float64]("3.14", 0.0)
		require.NoError(t, err)
		assert.Equal(t, 3.14, result)
	})

	t.Run("InvalidFloat64WithDefault", func(t *testing.T) {
		result, err := ucast.StringOrDef[float64]("invalid", 1.23)
		assert.Error(t, err)
		assert.Equal(t, 1.23, result)
	})

	t.Run("ValidBool", func(t *testing.T) {
		result, err := ucast.StringOrDef[bool]("true", false)
		require.NoError(t, err)
		assert.Equal(t, true, result)
	})

	t.Run("InvalidBoolWithDefault", func(t *testing.T) {
		result, err := ucast.StringOrDef[bool]("invalid", false)
		assert.Error(t, err)
		assert.Equal(t, false, result)
	})
}
