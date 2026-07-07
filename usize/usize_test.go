/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2026.
 */

package usize_test

import (
	"testing"

	"git.casinomodule.org/casino27/basic-utils/v3/usize"
	"github.com/stretchr/testify/require"
)

func TestUnits(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected int
	}{
		{name: "Byte", value: usize.Byte, expected: 1},
		{name: "B", value: usize.B, expected: usize.Byte},

		{name: "Kilobyte", value: usize.Kilobyte, expected: 1000},
		{name: "Megabyte", value: usize.Megabyte, expected: 1000 * 1000},
		{name: "Gigabyte", value: usize.Gigabyte, expected: 1000 * 1000 * 1000},
		{name: "Terabyte", value: usize.Terabyte, expected: 1000 * 1000 * 1000 * 1000},
		{name: "Petabyte", value: usize.Petabyte, expected: 1000 * 1000 * 1000 * 1000 * 1000},

		{name: "Kibibyte", value: usize.Kibibyte, expected: 1024},
		{name: "Mebibyte", value: usize.Mebibyte, expected: 1024 * 1024},
		{name: "Gibibyte", value: usize.Gibibyte, expected: 1024 * 1024 * 1024},
		{name: "Tebibyte", value: usize.Tebibyte, expected: 1024 * 1024 * 1024 * 1024},
		{name: "Pebibyte", value: usize.Pebibyte, expected: 1024 * 1024 * 1024 * 1024 * 1024},

		{name: "KB", value: usize.KB, expected: usize.Kilobyte},
		{name: "MB", value: usize.MB, expected: usize.Megabyte},
		{name: "GB", value: usize.GB, expected: usize.Gigabyte},
		{name: "TB", value: usize.TB, expected: usize.Terabyte},
		{name: "PB", value: usize.PB, expected: usize.Petabyte},

		{name: "KiB", value: usize.KiB, expected: usize.Kibibyte},
		{name: "MiB", value: usize.MiB, expected: usize.Mebibyte},
		{name: "GiB", value: usize.GiB, expected: usize.Gibibyte},
		{name: "TiB", value: usize.TiB, expected: usize.Tebibyte},
		{name: "PiB", value: usize.PiB, expected: usize.Pebibyte},

		{name: "Kilobit", value: usize.Kilobit, expected: usize.Kilobyte / 8},
		{name: "Megabit", value: usize.Megabit, expected: usize.Megabyte / 8},
		{name: "Gigabit", value: usize.Gigabit, expected: usize.Gigabyte / 8},
		{name: "Terabit", value: usize.Terabit, expected: usize.Terabyte / 8},
		{name: "Petabit", value: usize.Petabit, expected: usize.Petabyte / 8},

		{name: "Kibibit", value: usize.Kibibit, expected: usize.Kibibyte / 8},
		{name: "Mebibit", value: usize.Mebibit, expected: usize.Mebibyte / 8},
		{name: "Gibibit", value: usize.Gibibit, expected: usize.Gibibyte / 8},
		{name: "Tebibit", value: usize.Tebibit, expected: usize.Tebibyte / 8},
		{name: "Pebibit", value: usize.Pebibit, expected: usize.Pebibyte / 8},

		{name: "Kbit", value: usize.Kbit, expected: usize.Kilobit},
		{name: "Mbit", value: usize.Mbit, expected: usize.Megabit},
		{name: "Gbit", value: usize.Gbit, expected: usize.Gigabit},
		{name: "Tbit", value: usize.Tbit, expected: usize.Terabit},
		{name: "Pbit", value: usize.Pbit, expected: usize.Petabit},

		{name: "Kibit", value: usize.Kibit, expected: usize.Kibibit},
		{name: "Mibit", value: usize.Mibit, expected: usize.Mebibit},
		{name: "Gibit", value: usize.Gibit, expected: usize.Gibibit},
		{name: "Tibit", value: usize.Tibit, expected: usize.Tebibit},
		{name: "Pibit", value: usize.Pibit, expected: usize.Pebibit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.value)
		})
	}
}

func TestConstantsCanBeUsedInConstExpressions(t *testing.T) {
	const databaseLoadChunkSize = 8 * usize.MiB

	require.Equal(t, 8*1024*1024, databaseLoadChunkSize)
}
