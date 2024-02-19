/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package strutils_test

import (
	"testing"

	strutils "github.com/kordax/basic-utils/str-utils"
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
			result := strutils.Def(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func ptr(s string) *string {
	return &s
}
