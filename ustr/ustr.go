/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ustr

import (
	"strings"
)

// Def behaves as Or(val, "")), so it returns default empty string if value is not present,
func Def(val *string) string {
	if val == nil {
		return ""
	}

	return *val
}

func Concat(vals ...string) string {
	return strings.Join(vals, "")
}

func IsBlank(value string) bool {
	return strings.TrimSpace(value) == ""
}

func DefaultIfBlank(value string, def string) string {
	if IsBlank(value) {
		return def
	}

	return value
}

func Ptr(value string) *string {
	return &value
}

func TrimPtr(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	return &trimmed
}

func SplitAndTrim(value string, sep string) []string {
	parts := strings.Split(value, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

func NormalizeSpace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
