/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2025.
 */

package uctx

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetGet(t *testing.T) {
	ctx := GetContext()
	ctx.Set("foo", 42)

	val := ctx.Get("foo")
	require.Equal(t, 42, val)
}

func TestDelete(t *testing.T) {
	ctx := GetContext()
	ctx.Set("bar", "hello")
	ctx.Delete("bar")

	val := ctx.Get("bar")
	assert.Nil(t, val)
}

func TestMustGetPanicsOnMissing(t *testing.T) {
	ctx := GetContext()
	ctx.Delete("missing") // ensure it's missing

	require.PanicsWithValue(t, "uctx: value not found for key", func() {
		ctx.MustGet("missing")
	})
}

func TestMustGetReturnsValue(t *testing.T) {
	ctx := GetContext()
	ctx.Set("baz", "value")
	val := ctx.MustGet("baz")
	assert.Equal(t, "value", val)
}

func TestTypeAssertionPanics(t *testing.T) {
	ctx := GetContext()
	ctx.Set("num", 123)
	val := ctx.Get("num")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on type assertion, got none")
		}
	}()
	_ = val.(string) // will panic: interface {} is int, not string
}

func TestConcurrentAccess(t *testing.T) {
	ctx := GetContext()
	const goroutines = 50
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "k"
			ctx.Set(key, i)
			_ = ctx.Get(key)
		}(i)
	}

	wg.Wait()
	// Should not panic or data race
}
