/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package sql_utils_test

import (
	"math"
	"testing"
	"time"

	sqlutils "github.com/kordax/basic-utils/sql-utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNullString(t *testing.T) {
	result := sqlutils.NullString("hello")
	require.True(t, result.Valid)
	assert.Equal(t, "hello", result.String)

	result = sqlutils.NullString("")
	require.False(t, result.Valid)
}

func TestNullBool(t *testing.T) {
	result := sqlutils.NullBool(true)
	require.True(t, result.Valid)
	assert.True(t, result.Bool)

	result = sqlutils.NullBool(false)
	require.True(t, result.Valid)
	assert.False(t, result.Bool)
}

func TestNullTime(t *testing.T) {
	now := time.Now()
	result := sqlutils.NullTime(now)
	require.True(t, result.Valid)
	assert.Equal(t, now, result.Time)

	result = sqlutils.NullTime(time.Time{})
	require.False(t, result.Valid)
}

func TestNullFloat64(t *testing.T) {
	result := sqlutils.NullFloat64(1.23)
	require.True(t, result.Valid)
	assert.Equal(t, 1.23, result.Float64)

	result = sqlutils.NullFloat64(math.NaN())
	require.False(t, result.Valid)
}

// For generic Null function, testing with one type should suffice as it uses reflection
func TestNewNull(t *testing.T) {
	result := sqlutils.NewNull("hello")
	require.True(t, result.Valid)
	assert.Equal(t, "hello", result.V)

	result = sqlutils.NewNull("")
	require.False(t, result.Valid)
}
