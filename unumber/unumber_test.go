/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package unumber_test

import (
	"math/big"
	"testing"

	"github.com/kordax/basic-utils/unumber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromString(t *testing.T) {
	tests := []struct {
		input       string
		isBig       bool
		numType     unumber.ValueType
		expectedErr bool
	}{
		{"123", false, unumber.Uint, false},
		{"-123", false, unumber.Int, false},
		{"123.456", false, unumber.Float, false},
		{"1234567890123456789012345678901234567890", true, unumber.BigInt, false},
		{"123.4567890123456789012345678901234567890", true, unumber.BigFloat, false},
		{"abc", false, 0, true},
		{"123.abc", false, 0, true},
	}

	for _, tt := range tests {
		got, err := unumber.FromString(tt.input, tt.isBig)
		if (err != nil) != tt.expectedErr {
			t.Errorf("FromString(%q, %v) unexpected error: %v", tt.input, tt.isBig, err)
			continue
		}

		if err == nil && got.T() != tt.numType {
			t.Errorf("FromString(%q, %v) expected type %v, got %v", tt.input, tt.isBig, tt.numType, got.T())
		}
	}

	// Test values for specific types
	num, err := unumber.FromString("123", false)
	if err != nil || num.Ui() != 123 {
		t.Errorf("Expected integer value 123, got %v with error %v", num.Ui(), err)
	}
	assert.Equal(t, num.T(), unumber.Uint)

	num, err = unumber.FromString("-123", false)
	if err != nil || num.I() != -123 {
		t.Errorf("Expected integer value -123, got %v with error %v", num.I(), err)
	}
	assert.Equal(t, num.T(), unumber.Int)

	num, err = unumber.FromString("123.456", false)
	if err != nil || num.F() != 123.456 {
		t.Errorf("Expected float value 123.456, got %v with error %v", num, err)
	}
	assert.Equal(t, num.T(), unumber.Float)

	num, err = unumber.FromString("1234567890123456789012345678901234567890", true)
	bigInt, _ := big.NewInt(0).SetString("1234567890123456789012345678901234567890", 10)
	if err != nil || num.Bi().Cmp(bigInt) != 0 {
		t.Errorf("Expected big integer value, got %v with error %v", num, err)
	}
	assert.Equal(t, num.T(), unumber.BigInt)

	num, err = unumber.FromString("123.4567890123456789012345678901234567890", true)
	expectedBigFloat, _ := new(big.Float).SetString("123.4567890123456789012345678901234567890")
	if err != nil || num.Bf().Cmp(expectedBigFloat) != 0 {
		t.Errorf("Expected big float value, got %v with error %v", num, err)
	}
	assert.Equal(t, num.T(), unumber.BigFloat)
}

func TestNewInt(t *testing.T) {
	val := 5
	num := unumber.NewInt(val)
	if num.T() != unumber.Int || num.I() != val {
		t.Errorf("Expected Int type with value %d, but got %v with value %d", val, num.T(), num.I())
	}
}

func TestNewFloat(t *testing.T) {
	val := 5.5
	num := unumber.NewFloat(val)
	if num.T() != unumber.Float || num.F() != val {
		t.Errorf("Expected Float type with value %f, but got %v with value %f", val, num.T(), num.F())
	}
}

func TestNewUint(t *testing.T) {
	val := uint64(5)
	num := unumber.NewUint(val)
	if num.T() != unumber.Uint || num.Ui() != val {
		t.Errorf("Expected Uint type with value %d, but got %v with value %d", val, num.T(), num.Ui())
	}
}

func TestNewBigInt(t *testing.T) {
	val := big.NewInt(5)
	num := unumber.NewBigInt(val)
	if num.T() != unumber.BigInt || num.Bi().Cmp(val) != 0 {
		t.Errorf("Expected BigInt type with value %v, but got %v with value %v", val, num.T(), num.Bi())
	}
}

func TestNewBigFloat(t *testing.T) {
	val := big.NewFloat(5.5)
	num := unumber.NewBigFloat(val)
	if num.T() != unumber.BigFloat || num.Bf().Cmp(val) != 0 {
		t.Errorf("Expected BigFloat type with value %v, but got %v with value %v", val, num.T(), num.Bf())
	}
}

func TestNumberString(t *testing.T) {
	tests := []struct {
		num      *unumber.Number
		expected string
	}{
		{unumber.NewInt(5), "5"},
		{unumber.NewFloat(5.5), "5.5"},
		{unumber.NewUint(5), "5"},
		{unumber.NewBigInt(big.NewInt(5)), "5"},
		{unumber.NewBigFloat(big.NewFloat(5.5)), "5.5"},
	}

	for _, test := range tests {
		str := test.num.String()
		if str != test.expected {
			t.Errorf("Expected %s, but got %s", test.expected, str)
		}
	}
}

func TestDenominate(t *testing.T) {
	tests := []struct {
		name        string
		amount      float64
		denominator int
		expected    int64
		expectError bool
	}{
		{"Normal case with small value", 1.0, 6, 1000000, false},
		{"Normal case with larger value", 12345.6789, 4, 123456789, false},
		{"Small float with denominator", 2.5, 1, 25, false},
		{"Small float with larger denominator", 1.23, 2, 123, false},
		{"Edge case with zero amount", 0, 6, 0, false},
		{"Edge case with minimum int64", -9223372036854775808, 0, -9223372036854775808, false},
		{"Denominator greater than 15", 1.0, 16, 0, true}, // Denominator greater than 15 should fail
		{"Edge case with float 15 signs overflow", 9223372036854.775, 3, 9223372036854776, false},
		{"Edge case with float 14 signs and no overflow", 9223372036854.75, 3, 9223372036854750, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := unumber.Denominate[float64, int64](tt.amount, tt.denominator)

			if tt.expectError {
				require.Error(t, err, "Test case '%s' should have returned an error", tt.name)
			} else {
				require.NoError(t, err, "Test case '%s' should not have returned an error", tt.name)
			}

			assert.IsType(t, int64(0), result, "Test case '%s' failed: result should be of type int64", tt.name)
			assert.Equal(t, tt.expected, result, "Test case '%s' failed: expected %v, got %v", tt.name, tt.expected, result)
		})
	}
}

func TestParseDenominated(t *testing.T) {
	tests := []struct {
		name        string
		amount      float64
		denominator int
		expected    float64
		expectError bool
	}{
		{"Normal case with small value", 1000000, 6, 1.0, false},
		{"Normal case with decimal result", 123456789, 4, 12345.6789, false},
		{"Small amount with denominator", 25, 1, 2.5, false},
		{"Zero amount with denominator", 0, 6, 0.0, false},
		{"Small amount with large denominator", 123, 2, 1.23, false},
		{"Edge case with minimum positive amount", 1, 6, 0.000001, false},
		{"Edge case with potential overflow", 1e15, 15, 1.0, false},
		{"Invalid large denominator", 1.0, 308, 0, true}, // Should return an error due to overflow
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := unumber.ParseDenominated(tt.amount, tt.denominator)

			if tt.expectError {
				require.Error(t, err, "Test case '%s' should have returned an error", tt.name)
			} else {
				require.NoError(t, err, "Test case '%s' should not have returned an error", tt.name)
				assert.Equal(t, tt.expected, result, "Test case '%s' failed: expected %v, got %v", tt.name, tt.expected, result)
			}
		})
	}
}

func TestAsDenominated(t *testing.T) {
	tests := []struct {
		name         string
		value        float64
		denomination int
		expected     int64
		expectError  bool
	}{
		{
			name:         "Valid case with 2 decimal places",
			value:        12345.6789,
			denomination: 2,
			expected:     1234567,
			expectError:  false,
		},
		{
			name:         "Valid case with no decimal places",
			value:        12345.6789,
			denomination: 0,
			expected:     12345,
			expectError:  false,
		},
		{
			name:         "Valid case with rounding",
			value:        12345.6789,
			denomination: 3,
			expected:     12345678,
			expectError:  false,
		},
		{
			name:         "Denomination greater than 15",
			value:        12345.6789,
			denomination: 16,
			expected:     0,
			expectError:  true,
		},
		{
			name:         "Large number with potential overflow",
			value:        9223372036854.775,
			denomination: 3,
			expected:     9223372036854776,
			expectError:  false,
		},
		{
			name:         "Edge case with smallest positive value",
			value:        0.0001,
			denomination: 6,
			expected:     100,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := unumber.AsDenom[int64, float64](tt.value, tt.denomination)

			if tt.expectError {
				require.Error(t, err, "Expected an error but got none")
			} else {
				require.NoError(t, err, "Expected no error but got one")
				assert.Equal(t, tt.expected, result.Value(), "Expected value %v, got %v", tt.expected, result.Value())
				assert.Equal(t, tt.denomination, result.Denominator(), "Expected denomination %v, got %v", tt.denomination, result.Denominator())
			}
		})
	}
}
