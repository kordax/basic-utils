/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uos

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kordax/basic-utils/uref"
)

// MappingFunc is a type for functions that convert a string to a pointer of type T, returning an error if the conversion fails.
type MappingFunc[T any] func(value string) (*T, error)

// MapString returns string value
func MapString(value string) (*string, error) {
	return &value, nil
}

// MapStringToTrimmed returns string a value trimmed of space characters.
func MapStringToTrimmed(value string) (*string, error) {
	return uref.Ref(strings.TrimSpace(value)), nil
}

// MapStringToInt maps a string value to an *int.
func MapStringToInt(value string) (*int, error) {
	result, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// MapStringToInt64 maps a string value to an *int64.
func MapStringToInt64(value string) (*int64, error) {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// MapStringToInt32 maps a string value to an *int32.
func MapStringToInt32(value string) (*int32, error) {
	temp, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return nil, err
	}
	result := int32(temp)
	return &result, nil
}

// MapStringToInt16 maps a string value to an *int16.
func MapStringToInt16(value string) (*int16, error) {
	temp, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return nil, err
	}
	result := int16(temp)
	return &result, nil
}

// MapStringToInt8 maps a string value to an *int8.
func MapStringToInt8(value string) (*int8, error) {
	temp, err := strconv.ParseInt(value, 10, 8)
	if err != nil {
		return nil, err
	}
	result := int8(temp)
	return &result, nil
}

// MapStringToUint maps a string value to an *uint.
func MapStringToUint(value string) (*uint, error) {
	temp, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return nil, err
	}
	result := uint(temp)
	return &result, nil
}

// MapStringToUint64 maps a string value to an *uint64.
func MapStringToUint64(value string) (*uint64, error) {
	result, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// MapStringToUint32 maps a string value to an *uint32.
func MapStringToUint32(value string) (*uint32, error) {
	temp, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return nil, err
	}
	result := uint32(temp)
	return &result, nil
}

// MapStringToUint16 maps a string value to an *uint16.
func MapStringToUint16(value string) (*uint16, error) {
	temp, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return nil, err
	}
	result := uint16(temp)
	return &result, nil
}

// MapStringToUint8 maps a string value to an *uint8.
func MapStringToUint8(value string) (*uint8, error) {
	temp, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return nil, err
	}
	result := uint8(temp)
	return &result, nil
}

// MapStringToFloat64 maps a string value to a *float64.
func MapStringToFloat64(value string) (*float64, error) {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// MapStringToFloat32 maps a string value to a *float32.
func MapStringToFloat32(value string) (*float32, error) {
	temp, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return nil, err
	}
	result := float32(temp)
	return &result, nil
}

// MapStringToDuration maps value to time.Duration.
func MapStringToDuration(value string) (*time.Duration, error) {
	d, err := time.ParseDuration(value)
	return &d, err
}

// MapStringToTime creates a function to convert a string to a time.Time using the specified layout.
// This allows for flexibility in parsing different time formats.
func MapStringToTime(layout string) MappingFunc[time.Time] {
	return func(value string) (*time.Time, error) {
		t, err := time.Parse(layout, value)
		if err != nil {
			return nil, err
		}
		return &t, nil
	}
}

// MapStringToBase64 decodes a Base64 encoded string into its original representation.
func MapStringToBase64(value string) (*string, error) {
	bytes, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	result := string(bytes)
	return &result, nil
}

// MapStringToHex decodes a Hex encoded string into its original representation.
func MapStringToHex(value string) (*[]byte, error) {
	bytes, err := hex.DecodeString(value)
	if err != nil {
		return nil, err
	}
	result := bytes
	return &result, nil
}

// MapStringToURL parses a string into a *url.URL.
func MapStringToURL(value string) (*url.URL, error) {
	u, err := url.Parse(value)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// MapStringToBool parses a string into a *bool.
func MapStringToBool(value string) (*bool, error) {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return nil, err
	}

	return &b, nil
}
