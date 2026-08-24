/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uos_test

import (
	"net/url"
	"os"
	"testing"
	"time"

	"git.casinomodule.org/casino27/basic-utils/v4/uos"
	"git.casinomodule.org/casino27/basic-utils/v4/uref"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnvOpt(t *testing.T) {
	key := "TEST_GET_ENV_OPT"

	t.Run("present", func(t *testing.T) {
		value := "some-value"
		t.Setenv(key, value)

		result := uos.GetEnvOpt(key)
		require.True(t, result.Present())
		assert.Equal(t, value, result.OrElse(""))
	})

	t.Run("missing", func(t *testing.T) {
		_ = os.Unsetenv(key)

		result := uos.GetEnvOpt(key)
		require.False(t, result.Present())
		assert.Nil(t, result.Get())
	})
}

func TestGetEnvOptAs(t *testing.T) {
	key := "TEST_GET_ENV_OPT_AS"

	t.Run("present", func(t *testing.T) {
		t.Setenv(key, "42")

		result := uos.GetEnvOptAs[int](key, uos.MapStringToInt)
		require.True(t, result.Present())
		assert.Equal(t, 42, result.OrElse(0))
	})

	t.Run("missing", func(t *testing.T) {
		_ = os.Unsetenv(key)

		result := uos.GetEnvOptAs[int](key, uos.MapStringToInt)
		require.False(t, result.Present())
		assert.Nil(t, result.Get())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Setenv(key, "abc")

		result := uos.GetEnvOptAs[int](key, uos.MapStringToInt)
		require.False(t, result.Present())
		assert.Nil(t, result.Get())
	})
}

func TestGetEnvOptNumeric(t *testing.T) {
	key := "TEST_GET_ENV_OPT_NUMERIC"

	t.Run("present", func(t *testing.T) {
		t.Setenv(key, "123")

		result := uos.GetEnvOptNumeric[int64](key)
		require.True(t, result.Present())
		assert.Equal(t, int64(123), result.OrElse(0))
	})

	t.Run("missing", func(t *testing.T) {
		_ = os.Unsetenv(key)

		result := uos.GetEnvOptNumeric[int64](key)
		require.False(t, result.Present())
		assert.Nil(t, result.Get())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Setenv(key, "abc")

		result := uos.GetEnvOptNumeric[int64](key)
		require.False(t, result.Present())
		assert.Nil(t, result.Get())
	})
}

func TestGetEnvOptDuration(t *testing.T) {
	key := "TEST_GET_ENV_OPT_DURATION"

	t.Run("present", func(t *testing.T) {
		t.Setenv(key, "2m30s")

		result := uos.GetEnvOptDuration(key)
		require.True(t, result.Present())
		assert.Equal(t, 2*time.Minute+30*time.Second, result.OrElse(0))
	})

	t.Run("missing", func(t *testing.T) {
		_ = os.Unsetenv(key)

		result := uos.GetEnvOptDuration(key)
		require.False(t, result.Present())
		assert.Nil(t, result.Get())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Setenv(key, "not-a-duration")

		result := uos.GetEnvOptDuration(key)
		require.False(t, result.Present())
		assert.Nil(t, result.Get())
	})
}

func TestGetEnvOptTime(t *testing.T) {
	key := "TEST_GET_ENV_OPT_TIME"
	layout := time.RFC3339
	expected := time.Date(2024, 10, 15, 12, 30, 45, 0, time.UTC)

	t.Run("present", func(t *testing.T) {
		t.Setenv(key, expected.Format(layout))

		result := uos.GetEnvOptTime(key, layout)
		require.True(t, result.Present())
		assert.True(t, expected.Equal(result.OrElse(time.Time{})))
	})

	t.Run("missing", func(t *testing.T) {
		_ = os.Unsetenv(key)

		result := uos.GetEnvOptTime(key, layout)
		require.False(t, result.Present())
		assert.Nil(t, result.Get())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Setenv(key, "invalid-time")

		result := uos.GetEnvOptTime(key, layout)
		require.False(t, result.Present())
		assert.Nil(t, result.Get())
	})
}

func TestGetEnvOptURL(t *testing.T) {
	key := "TEST_GET_ENV_OPT_URL"
	expected := &url.URL{
		Scheme: "https",
		Host:   "example.com",
		Path:   "/test",
	}

	t.Run("present", func(t *testing.T) {
		t.Setenv(key, expected.String())

		result := uos.GetEnvOptURL(key)
		require.True(t, result.Present())
		assert.Equal(t, expected.String(), uref.Ref(result.OrElse(url.URL{})).String())
	})

	t.Run("missing", func(t *testing.T) {
		_ = os.Unsetenv(key)

		result := uos.GetEnvOptURL(key)
		require.False(t, result.Present())
		assert.Nil(t, result.Get())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Setenv(key, ":// bad url")

		result := uos.GetEnvOptURL(key)
		require.False(t, result.Present())
		assert.Nil(t, result.Get())
	})
}

func TestGetEnvOptBool(t *testing.T) {
	key := "TEST_GET_ENV_OPT_BOOL"

	t.Run("present true", func(t *testing.T) {
		t.Setenv(key, "true")

		result := uos.GetEnvOptBool(key)
		require.True(t, result.Present())
		assert.True(t, result.OrElse(false))
	})

	t.Run("present false", func(t *testing.T) {
		t.Setenv(key, "false")

		result := uos.GetEnvOptBool(key)
		require.True(t, result.Present())
		assert.False(t, result.OrElse(true))
	})

	t.Run("missing", func(t *testing.T) {
		_ = os.Unsetenv(key)

		result := uos.GetEnvOptBool(key)
		require.False(t, result.Present())
		assert.Nil(t, result.Get())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Setenv(key, "not-bool")

		result := uos.GetEnvOptBool(key)
		require.False(t, result.Present())
		assert.Nil(t, result.Get())
	})
}
