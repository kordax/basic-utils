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
