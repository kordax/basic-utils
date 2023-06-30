/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package number_test

import (
	"math/big"
	"testing"

	"github.com/kordax/basic-utils/number"
	"github.com/stretchr/testify/assert"
)

func TestFromString(t *testing.T) {
	tests := []struct {
		input       string
		isBig       bool
		numType     number.ValueType
		expectedErr bool
	}{
		{"123", false, number.Uint, false},
		{"-123", false, number.Int, false},
		{"123.456", false, number.Float, false},
		{"1234567890123456789012345678901234567890", true, number.BigInt, false},
		{"123.4567890123456789012345678901234567890", true, number.BigFloat, false},
		{"abc", false, 0, true},
		{"123.abc", false, 0, true},
	}

	for _, tt := range tests {
		got, err := number.FromString(tt.input, tt.isBig)
		if (err != nil) != tt.expectedErr {
			t.Errorf("FromString(%q, %v) unexpected error: %v", tt.input, tt.isBig, err)
			continue
		}

		if err == nil && got.T() != tt.numType {
			t.Errorf("FromString(%q, %v) expected type %v, got %v", tt.input, tt.isBig, tt.numType, got.T())
		}
	}

	// Test values for specific types
	num, err := number.FromString("123", false)
	if err != nil || num.Ui() != 123 {
		t.Errorf("Expected integer value 123, got %v with error %v", num.Ui(), err)
	}
	assert.Equal(t, num.T(), number.Uint)

	num, err = number.FromString("-123", false)
	if err != nil || num.I() != -123 {
		t.Errorf("Expected integer value -123, got %v with error %v", num.I(), err)
	}
	assert.Equal(t, num.T(), number.Int)

	num, err = number.FromString("123.456", false)
	if err != nil || num.F() != 123.456 {
		t.Errorf("Expected float value 123.456, got %v with error %v", num, err)
	}
	assert.Equal(t, num.T(), number.Float)

	num, err = number.FromString("1234567890123456789012345678901234567890", true)
	bigInt, _ := big.NewInt(0).SetString("1234567890123456789012345678901234567890", 10)
	if err != nil || num.Bi().Cmp(bigInt) != 0 {
		t.Errorf("Expected big integer value, got %v with error %v", num, err)
	}
	assert.Equal(t, num.T(), number.BigInt)

	num, err = number.FromString("123.4567890123456789012345678901234567890", true)
	expectedBigFloat, _ := new(big.Float).SetString("123.4567890123456789012345678901234567890")
	if err != nil || num.Bf().Cmp(expectedBigFloat) != 0 {
		t.Errorf("Expected big float value, got %v with error %v", num, err)
	}
	assert.Equal(t, num.T(), number.BigFloat)
}
