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
	"time"
)

// MappingFunc is a type for functions that convert a string to a pointer of type T, returning an error if the conversion fails.
type MappingFunc[T any] func(value string) (*T, error)

// MapString returns string value
func MapString(value string) (*string, error) {
	return &value, nil
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
