/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ucache

import (
	"testing"
	"time"

	"github.com/kordax/basic-utils/v3/uopt"
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

func TestInMemoryComparableMapCache_BufferInternals(t *testing.T) {
	c := NewInMemoryComparableMapCacheWithOptions[string, int](InMemoryComparableMapCacheOptions{
		BufferedMaxKeys: 2,
	})

	assert.False(t, c.latestBufferedEntry(bufferedComparableMapEntry[string, int]{key: "missing", seq: 1}))

	c.bufferedLatest.Store("key", "bad-seq")
	assert.False(t, c.latestBufferedEntry(bufferedComparableMapEntry[string, int]{key: "key", seq: 1}))

	c.bufferedLatest.Store("key", uint64(2))
	assert.False(t, c.latestBufferedEntry(bufferedComparableMapEntry[string, int]{key: "key", seq: 1}))
	assert.True(t, c.latestBufferedEntry(bufferedComparableMapEntry[string, int]{key: "key", seq: 2}))

	assert.True(t, c.admitBufferedKey("key"))
	assert.True(t, c.admitBufferedKey("second"))
	assert.False(t, c.admitBufferedKey("third"))

	c.clearBufferedKeys()
	assert.Zero(t, c.bufferedKeyLen.Load())
	assert.True(t, c.admitBufferedKey("third"))
}

func TestInMemoryComparableMapCache_BufferCloseAndRejectedWrites(t *testing.T) {
	c := NewInMemoryComparableMapCacheWithOptions[string, int](InMemoryComparableMapCacheOptions{
		BufferedWorkers: 1,
	})

	c.closeBuffered()
	assert.False(t, c.enqueueBuffered("key", 1, true))

	c = NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())
	c.Set("key", 1)
	c.closeBuffered()
	c.closeBuffered()

	value, ok := c.GetValue("key")
	assert.True(t, ok)
	assert.Equal(t, 1, value)
	assert.False(t, c.enqueueBuffered("after-close", 2, true))
}
