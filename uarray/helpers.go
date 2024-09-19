/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uarray

import "strconv"

// Please note that these helper methods ignore parsing errors and therefore should be used only if you know value types.

// StringToInt converts a string to int
func StringToInt(v *string) int {
	parsed, _ := strconv.Atoi(*v)
	return parsed
}

// StringToInt32 converts a string to int32
func StringToInt32(v *string) int32 {
	parsed, _ := strconv.ParseInt(*v, 10, 32)
	return int32(parsed)
}

// StringToInt64 converts a string to int64
func StringToInt64(v *string) int64 {
	parsed, _ := strconv.ParseInt(*v, 10, 64)
	return parsed
}

// StringToFloat32 converts a string to float32
func StringToFloat32(v *string) float32 {
	parsed, _ := strconv.ParseFloat(*v, 32)
	return float32(parsed)
}

// StringToFloat64 converts a string to float64
func StringToFloat64(v *string) float64 {
	parsed, _ := strconv.ParseFloat(*v, 64)
	return parsed
}

// StringToBool converts a string to bool
func StringToBool(v *string) bool {
	parsed, _ := strconv.ParseBool(*v)
	return parsed
}

// Float64ToFloat32 converts a float64 to float32
func Float64ToFloat32(v *float64) float32 {
	return float32(*v)
}

// Int64ToInt32 converts an int64 to int32
func Int64ToInt32(v *int64) int32 {
	return int32(*v)
}

// IntToString converts an int to string
func IntToString(v *int) string {
	return strconv.Itoa(*v)
}

// Int64ToString converts an int64 to string
func Int64ToString(v *int64) string {
	return strconv.FormatInt(*v, 10)
}

// Float32ToString converts a float32 to string
func Float32ToString(v *float32) string {
	return strconv.FormatFloat(float64(*v), 'f', -1, 32)
}

// Float64ToString converts a float64 to string
func Float64ToString(v *float64) string {
	return strconv.FormatFloat(*v, 'f', -1, 64)
}

// BoolToString converts a bool to string
func BoolToString(v *bool) string {
	return strconv.FormatBool(*v)
}
