/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uopt_test

import (
	"testing"

	"github.com/kordax/basic-utils/uopt"
	"github.com/stretchr/testify/assert"
)

func TestNullBool(t *testing.T) {
	v := uopt.NullBool()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullInt(t *testing.T) {
	v := uopt.NullInt()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullInt8(t *testing.T) {
	v := uopt.NullInt8()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullInt16(t *testing.T) {
	v := uopt.NullInt16()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullInt32(t *testing.T) {
	v := uopt.NullInt32()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullInt64(t *testing.T) {
	v := uopt.NullInt64()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullUint(t *testing.T) {
	v := uopt.NullUint()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullUint8(t *testing.T) {
	v := uopt.NullUint8()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullUint16(t *testing.T) {
	v := uopt.NullUint16()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullUint32(t *testing.T) {
	v := uopt.NullUint32()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullUint64(t *testing.T) {
	v := uopt.NullUint64()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullFloat32(t *testing.T) {
	v := uopt.NullFloat32()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullFloat64(t *testing.T) {
	v := uopt.NullFloat64()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullString(t *testing.T) {
	v := uopt.NullString()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullByte(t *testing.T) {
	v := uopt.NullByte()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullRune(t *testing.T) {
	v := uopt.NullRune()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullComplex64(t *testing.T) {
	v := uopt.NullComplex64()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}

func TestNullComplex128(t *testing.T) {
	v := uopt.NullComplex128()
	assert.False(t, v.Present())
	assert.Nil(t, v.Get())
}
