/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2025.
 */

package uctx

import (
	"sync"

	"github.com/kordax/basic-utils/v2/uonce"
)

var getGlobalContext = uonce.Once(func() *UContext {
	return &UContext{
		data: make(map[any]any),
	}
})

// UContext provides a concurrent-safe, app-global key-value store.
//
// Use GetContext() to access the singleton instance from anywhere in your application.
// Values can be stored and retrieved using any comparable key type.
//
// Example usage:
//
//	ctx := uctx.GetContext()
//	ctx.Set("config", myConfig)
//	cfg := ctx.Get("config").(*Config)
type UContext struct {
	mu   sync.RWMutex
	data map[any]any
}

// GetContext returns the singleton app-global context instance.
// This instance is safe for concurrent access and should be used
// throughout your application to store and retrieve global values.
func GetContext() *UContext {
	return getGlobalContext()
}

// Set stores the given value under the provided key.
// Overwrites any existing value for that key.
func (g *UContext) Set(key, value any) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.data[key] = value
}

// Get retrieves the value associated with the provided key.
// Returns nil if the key does not exist.
func (g *UContext) Get(key any) any {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.data[key]
}

// MustGet retrieves the value for the given key.
// Panics if the key is not present in the context.
func (g *UContext) MustGet(key any) any {
	val := g.Get(key)
	if val == nil {
		panic("uctx: value not found for key")
	}
	return val
}

// Delete removes the value associated with the provided key.
func (g *UContext) Delete(key any) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.data, key)
}
