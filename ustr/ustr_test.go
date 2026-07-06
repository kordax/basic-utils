/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ustr_test

import (
	"testing"

	"git.casinomodule.org/casino27/basic-utils/v3/ustr"
	"github.com/stretchr/testify/require"
)

func TestDef(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: "",
		},
		{
			name:     "non-empty input",
			input:    ptr("hello"),
			expected: "hello",
		},
		{
			name:     "empty input",
			input:    ptr(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ustr.Def(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStringHelpers(t *testing.T) {
	require.True(t, ustr.IsBlank(" \t\n"))
	require.False(t, ustr.IsBlank(" value "))
	require.Equal(t, "fallback", ustr.DefaultIfBlank(" ", "fallback"))
	require.Equal(t, "value", ustr.DefaultIfBlank("value", "fallback"))
	require.Equal(t, "value", *ustr.Ptr("value"))
	require.Nil(t, ustr.TrimPtr(nil))
	require.Equal(t, "value", *ustr.TrimPtr(ptr(" value ")))
	require.Equal(t, []string{"a", "b", "c"}, ustr.SplitAndTrim(" a, b ,, c ", ","))
	require.Equal(t, "a b c", ustr.NormalizeSpace(" a\t b\n c "))
}

func ptr(s string) *string {
	return &s
}
