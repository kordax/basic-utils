/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ucache

import (
	"testing"
	"time"

	"git.casinomodule.org/casino27/basic-utils/v3/uopt"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryComparableMapCacheOptions_Defaults(t *testing.T) {
	c := NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())

	assert.Equal(t, defaultComparableMapCacheBufferedWorkers, c.bufferWorkers)
	assert.Equal(t, defaultComparableMapCacheBufferedQueueSize, c.bufferQueueSize)
}

func TestInMemoryComparableMapCacheOptions_CustomBufferedSettings(t *testing.T) {
	c := NewInMemoryComparableMapCacheWithOptions[string, int](InMemoryComparableMapCacheOptions{
		TTL:               uopt.Of(time.Minute),
		BufferedWorkers:   2,
		BufferedQueueSize: 1024,
		BufferedMaxKeys:   128,
	})

	assert.Equal(t, 2, c.bufferWorkers)
	assert.Equal(t, 1024, c.bufferQueueSize)
	assert.Equal(t, int64(128), c.bufferedMaxKeys)
	assert.Equal(t, time.Minute, c.ttl())
}

func TestInMemoryComparableMapCacheOptions_InvalidBufferedSettingsUseDefaults(t *testing.T) {
	c := NewInMemoryComparableMapCacheWithOptions[string, int](InMemoryComparableMapCacheOptions{
		BufferedWorkers:   -1,
		BufferedQueueSize: 0,
		BufferedMaxKeys:   -1,
	})

	assert.Equal(t, defaultComparableMapCacheBufferedWorkers, c.bufferWorkers)
	assert.Equal(t, defaultComparableMapCacheBufferedQueueSize, c.bufferQueueSize)
	assert.Zero(t, c.bufferedMaxKeys)
}
