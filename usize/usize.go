/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2026.
 */

// Package usize provides byte-sized data unit constants.
package usize

const (
	Byte = 1
	B    = Byte
)

const (
	Kilobyte = 1000 * Byte
	Megabyte = 1000 * Kilobyte
	Gigabyte = 1000 * Megabyte
	Terabyte = 1000 * Gigabyte
	Petabyte = 1000 * Terabyte
)

const (
	Kibibyte = 1024 * Byte
	Mebibyte = 1024 * Kibibyte
	Gibibyte = 1024 * Mebibyte
	Tebibyte = 1024 * Gibibyte
	Pebibyte = 1024 * Tebibyte
)

const (
	KB = Kilobyte
	MB = Megabyte
	GB = Gigabyte
	TB = Terabyte
	PB = Petabyte
)

const (
	KiB = Kibibyte
	MiB = Mebibyte
	GiB = Gibibyte
	TiB = Tebibyte
	PiB = Pebibyte
)

const (
	Kilobit = Kilobyte / 8
	Megabit = Megabyte / 8
	Gigabit = Gigabyte / 8
	Terabit = Terabyte / 8
	Petabit = Petabyte / 8
)

const (
	Kibibit = Kibibyte / 8
	Mebibit = Mebibyte / 8
	Gibibit = Gibibyte / 8
	Tebibit = Tebibyte / 8
	Pebibit = Pebibyte / 8
)

const (
	Kbit = Kilobit
	Mbit = Megabit
	Gbit = Gigabit
	Tbit = Terabit
	Pbit = Petabit
)

const (
	Kibit = Kibibit
	Mibit = Mebibit
	Gibit = Gibibit
	Tibit = Tebibit
	Pibit = Pebibit
)
