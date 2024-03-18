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
