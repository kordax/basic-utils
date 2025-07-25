/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ucast

import "strconv"

// Please note that these helper methods ignore parsing errors and therefore should be used only if you know value types.

// StringToInt converts a string to int
func StringToInt(v string) (int, error) {
	parsed, err := strconv.Atoi(v)
	return parsed, err
}

// StringToInt8 converts a string to int8
func StringToInt8(v string) (int8, error) {
	parsed, err := strconv.ParseInt(v, 10, 8)
	return int8(parsed), err
}

// StringToInt16 converts a string to int16
func StringToInt16(v string) (int16, error) {
	parsed, err := strconv.ParseInt(v, 10, 16)
	return int16(parsed), err
}

// StringToInt32 converts a string to int32
func StringToInt32(v string) (int32, error) {
	parsed, err := strconv.ParseInt(v, 10, 32)
	return int32(parsed), err
}

// StringToInt64 converts a string to int64
func StringToInt64(v string) (int64, error) {
	return strconv.ParseInt(v, 10, 64)
}

// StringToUint converts a string to uint
func StringToUint(v string) (uint, error) {
	parsed, err := strconv.ParseUint(v, 10, 0)
	return uint(parsed), err
}

// StringToUint8 converts a string to uint8
func StringToUint8(v string) (uint8, error) {
	parsed, err := strconv.ParseUint(v, 10, 8)
	return uint8(parsed), err
}

// StringToUint16 converts a string to uint16
func StringToUint16(v string) (uint16, error) {
	parsed, err := strconv.ParseUint(v, 10, 16)
	return uint16(parsed), err
}

// StringToUint32 converts a string to uint32
func StringToUint32(v string) (uint32, error) {
	parsed, err := strconv.ParseUint(v, 10, 32)
	return uint32(parsed), err
}

// StringToUint64 converts a string to uint64
func StringToUint64(v string) (uint64, error) {
	return strconv.ParseUint(v, 10, 64)
}

// StringToFloat32 converts a string to float32
func StringToFloat32(v string) (float32, error) {
	parsed, err := strconv.ParseFloat(v, 32)
	return float32(parsed), err
}

// StringToFloat64 converts a string to float64
func StringToFloat64(v string) (float64, error) {
	return strconv.ParseFloat(v, 64)
}

// StringToBool converts a string to bool
func StringToBool(v string) (bool, error) {
	return strconv.ParseBool(v)
}

// Float64ToFloat32 converts a float64 to float32
func Float64ToFloat32(v float64) float32 {
	return float32(v)
}

// IntToString converts an int to string
func IntToString(v int) string {
	return strconv.Itoa(v)
}

// Int8ToString converts an int64 to string
func Int8ToString(v int8) string {
	return strconv.FormatInt(int64(v), 10)
}

// Int16ToString converts an int64 to string
func Int16ToString(v int16) string {
	return strconv.FormatInt(int64(v), 10)
}

// Int32ToString converts an int64 to string
func Int32ToString(v int32) string {
	return strconv.FormatInt(int64(v), 10)
}

// Int64ToString converts an int64 to string
func Int64ToString(v int64) string {
	return strconv.FormatInt(v, 10)
}

// UintToString converts a uint to string
func UintToString(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}

// Uint8ToString converts a uint8 to string
func Uint8ToString(v uint8) string {
	return strconv.FormatUint(uint64(v), 10)
}

// Uint16ToString converts a uint16 to string
func Uint16ToString(v uint16) string {
	return strconv.FormatUint(uint64(v), 10)
}

// Uint32ToString converts a uint32 to string
func Uint32ToString(v uint32) string {
	return strconv.FormatUint(uint64(v), 10)
}

// Uint64ToString converts a uint64 to string
func Uint64ToString(v uint64) string {
	return strconv.FormatUint(v, 10)
}

// Float32ToString converts a float32 to string
func Float32ToString(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

// Float64ToString converts a float64 to string
func Float64ToString(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// BoolToString converts a bool to string
func BoolToString(v bool) string {
	return strconv.FormatBool(v)
}
