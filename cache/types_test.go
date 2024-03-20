package cache_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kordax/basic-utils/cache"
	"github.com/stretchr/testify/assert"
)

func TestStringKey_Key(t *testing.T) {
	uuid.EnableRandPool()

	value := cache.StringKey("MyTestString")
	key := value.Key()
	assert.NotNil(t, key)
	assert.NotZero(t, key)
	assert.EqualValues(t, 2057938083392025135, key)

	for i := 0; i < 10; i++ {
		t.Run("parallel key test", func(t *testing.T) {
			key = value.Key()
			assert.NotNil(t, key)
			assert.NotZero(t, key)
			assert.EqualValues(t, 2057938083392025135, key)
		})
	}

	value2 := value + "2"
	key2 := value2.Key()
	assert.NotNil(t, key2)
	assert.NotZero(t, key2)
	assert.NotEqualValues(t, key, key2)

	for i := 0; i < 100; i++ {
		t.Run("uuid key test", func(t *testing.T) {
			uKey := cache.StringKey(uuid.NewString()).Key()
			assert.NotNil(t, uKey)
			assert.NotZero(t, uKey)
			assert.NotEqualValues(t, key, uKey)
			assert.NotEqualValues(t, key2, uKey)
		})
	}
}

func TestIntKey_Key(t *testing.T) {
	value := cache.IntKey(123)
	key := value.Key()
	assert.EqualValues(t, 123, key)

	value2 := cache.IntKey(456)
	key2 := value2.Key()
	assert.NotEqual(t, key, key2)
}

func TestUIntKey_Key(t *testing.T) {
	value := cache.UIntKey(123)
	key := value.Key()
	assert.Equal(t, 123, key)

	value2 := cache.UIntKey(456)
	key2 := value2.Key()
	assert.NotEqual(t, key, key2)
}

func TestIntCompositeKey_Keys(t *testing.T) {
	value := cache.NewIntCompositeKey(123, 456)
	keys := value.Keys()
	assert.Len(t, keys, 2)
	assert.Contains(t, keys, int64(123))
	assert.Contains(t, keys, int64(456))

	value2 := cache.NewIntCompositeKey(789)
	keys2 := value2.Keys()
	assert.Len(t, keys2, 1)
	assert.NotEqual(t, keys, keys2)
}
