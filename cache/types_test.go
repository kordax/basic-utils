package cache_test

import (
	"testing"
	"time"

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

func TestStringKey_String(t *testing.T) {
	key := cache.StringKey("hello")
	assert.Equal(t, "hello", key.String())
}

func TestStrCompositeKey_Equals(t *testing.T) {
	key1 := cache.NewStrCompositeKey("a", "b", "c")
	key2 := cache.NewStrCompositeKey("a", "b", "c")
	key3 := cache.NewStrCompositeKey("d", "e", "f")

	assert.True(t, key1.Equals(key2))
	assert.False(t, key1.Equals(key3))
}

func TestIntKey_Key(t *testing.T) {
	value := cache.IntKey(123)
	key := value.Key()
	assert.EqualValues(t, 123, key)

	value2 := cache.IntKey(456)
	key2 := value2.Key()
	assert.NotEqual(t, key, key2)
}

func TestUIntKey_Equals(t *testing.T) {
	key1 := cache.UIntKey(100)
	key2 := cache.UIntKey(100)
	key3 := cache.UIntKey(200)

	assert.True(t, key1.Equals(key2))
	assert.False(t, key1.Equals(key3))
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

func TestIntCompositeKey_Equals(t *testing.T) {
	key1 := cache.NewIntCompositeKey(1, 2, 3)
	key2 := cache.NewIntCompositeKey(1, 2, 3)
	key3 := cache.NewIntCompositeKey(4, 5, 6)

	assert.True(t, key1.Equals(key2))
	assert.False(t, key1.Equals(key3))
}

func TestUIntCompositeKey_Equals(t *testing.T) {
	key1 := cache.NewUIntCompositeKey(1, 2, 3)
	key2 := cache.NewUIntCompositeKey(1, 2, 3)
	key3 := cache.NewUIntCompositeKey(4, 5, 6)

	assert.True(t, key1.Equals(key2))
	assert.False(t, key1.Equals(key3))
}

func TestUIntCompositeKey_Keys(t *testing.T) {
	key := cache.NewUIntCompositeKey(1, 2, 3)
	assert.EqualValues(t, []int64{1, 2, 3}, key.Keys())
}

func TestGenericCompositeKey_Equals(t *testing.T) {
	key1 := cache.NewGenericCompositeKey(cache.StringKey("test"), cache.IntKey(123))
	key2 := cache.NewGenericCompositeKey(cache.StringKey("test"), cache.IntKey(123))
	key3 := cache.NewGenericCompositeKey(cache.StringKey("different"), cache.IntKey(456))

	assert.True(t, key1.Equals(key2))
	assert.False(t, key1.Equals(key3))
}

func TestGenericCompositeKey_Keys(t *testing.T) {
	key := cache.NewGenericCompositeKey(cache.StringKey("test"), cache.IntKey(123))
	assert.NotNil(t, key.Keys()) // More specific tests might depend on the hashing algorithm
}

func TestUIntCompositeKey_Keys_String(t *testing.T) {
	key := cache.NewUIntCompositeKey(1, 2, 3)
	assert.EqualValues(t, []int64{1, 2, 3}, key.Keys())
	assert.Equal(t, "1, 2, 3", key.String())
}

func TestIntCompositeKey_Keys_String(t *testing.T) {
	key := cache.NewIntCompositeKey(1, 2, 3)
	assert.EqualValues(t, []int64{1, 2, 3}, key.Keys())
	assert.Equal(t, "1, 2, 3", key.String())
}

func TestStrCompositeKey_Keys_String(t *testing.T) {
	key := cache.NewStrCompositeKey("a", "b", "c")
	expectedKeys := []int64{97, 98, 99} // ASCII values of 'a', 'b', 'c'
	assert.EqualValues(t, expectedKeys, key.Keys())
	assert.Equal(t, "a, b, c", key.String())
}

func TestNewStringSliceValue(t *testing.T) {
	input := []string{"c", "a", "b"}
	sorted := []string{"a", "b", "c"}
	value := cache.NewStringSliceValue(input)

	assert.Equal(t, sorted, value.Values())
}

func TestStringSliceValue_Values(t *testing.T) {
	input := []string{"a", "b", "c"}
	value := cache.NewStringSliceValue(input)

	assert.Equal(t, input, value.Values())
}

func TestStringSliceValue_Equals(t *testing.T) {
	value1 := cache.NewStringSliceValue([]string{"a", "b", "c"})
	value2 := cache.NewStringSliceValue([]string{"a", "b", "c"})
	value3 := cache.NewStringSliceValue([]string{"d", "e", "f"})

	assert.True(t, value1.Equals(value2))
	assert.False(t, value1.Equals(value3))
	assert.False(t, value1.Equals(cache.StringKey("test string")))
}

func TestComparableSlice_Equals(t *testing.T) {
	slice1 := cache.ComparableSlice[cache.IntKey]{Data: []cache.IntKey{1, 2, 3}}
	slice2 := cache.ComparableSlice[cache.IntKey]{Data: []cache.IntKey{1, 2, 3}}
	slice3 := cache.ComparableSlice[cache.IntKey]{Data: []cache.IntKey{4, 5, 6}}

	assert.True(t, slice1.Equals(slice2))
	assert.False(t, slice1.Equals(slice3))
}

func TestComparableSlice_Equals_NotComparableSlice(t *testing.T) {
	c := cache.ComparableSlice[cache.IntKey]{Data: []cache.IntKey{1, 2, 3}}
	other := cache.NewComparableKey[int](1)

	assert.False(t, c.Equals(other))
}

func TestComparableSlice_Equals_Case2(t *testing.T) {
	slice1 := cache.ComparableSlice[DummyComparable]{
		Data: []DummyComparable{
			{1},
			{2},
			{3},
			{4},
		},
	}
	slice2 := cache.ComparableSlice[DummyComparable]{
		Data: []DummyComparable{
			{1},
			{2},
			{3},
			{4},
		},
	}
	assert.EqualValues(t, slice1, slice2)
	slice2 = cache.ComparableSlice[DummyComparable]{
		Data: []DummyComparable{
			{1},
			{2},
			{5},
			{4},
		},
	}
	assert.NotEqualValues(t, slice1, slice2)
	slice2 = cache.ComparableSlice[DummyComparable]{
		Data: []DummyComparable{
			{1},
			{2},
			{4},
			{3},
		},
	}
	assert.NotEqualValues(t, slice1, slice2)
	slice2 = cache.ComparableSlice[DummyComparable]{
		Data: []DummyComparable{
			{1},
			{2},
			{3},
			{4},
		},
	}
	assert.EqualValues(t, slice1, slice2)
}

func TestComparableKey_Equals(t *testing.T) {
	key1 := cache.NewComparableKey(123)
	key2 := cache.NewComparableKey(123)
	key3 := cache.NewComparableKey(456)

	assert.True(t, key1.Equals(key2), "keys with the same value should be equal")
	assert.False(t, key1.Equals(key3), "keys with different values should not be equal")
}

func TestComparableKey_String(t *testing.T) {
	key := cache.NewComparableKey(123)
	assert.Equal(t, "123", key.String(), "String() should return the correct string representation")

	key2 := cache.NewComparableKey("abc")
	assert.Equal(t, "abc", key2.String(), "String() should handle string types correctly")
}

func TestComparableKey_Equals_CoverAllTypes(t *testing.T) {
	assert.True(t, cache.NewComparableKey("abc").Equals(cache.NewComparableKey("abc")))
	assert.True(t, cache.NewComparableKey(42).Equals(cache.NewComparableKey(42)))
	assert.True(t, cache.NewComparableKey(int8(42)).Equals(cache.NewComparableKey(int8(42))))
	assert.True(t, cache.NewComparableKey(int16(42)).Equals(cache.NewComparableKey(int16(42))))
	assert.True(t, cache.NewComparableKey(int32(42)).Equals(cache.NewComparableKey(int32(42))))
	assert.True(t, cache.NewComparableKey(int64(42)).Equals(cache.NewComparableKey(int64(42))))
	assert.True(t, cache.NewComparableKey(uint(42)).Equals(cache.NewComparableKey(uint(42))))
	assert.True(t, cache.NewComparableKey(uint8(42)).Equals(cache.NewComparableKey(uint8(42))))
	assert.True(t, cache.NewComparableKey(uint16(42)).Equals(cache.NewComparableKey(uint16(42))))
	assert.True(t, cache.NewComparableKey(uint32(42)).Equals(cache.NewComparableKey(uint32(42))))
	assert.True(t, cache.NewComparableKey(uint64(42)).Equals(cache.NewComparableKey(uint64(42))))
	assert.True(t, cache.NewComparableKey(true).Equals(cache.NewComparableKey(true)))
	assert.True(t, cache.NewComparableKey(false).Equals(cache.NewComparableKey(false)))
	assert.True(t, cache.NewComparableKey(float32(42.1)).Equals(cache.NewComparableKey(float32(42.1))))
	assert.True(t, cache.NewComparableKey(float64(42.1)).Equals(cache.NewComparableKey(float64(42.1))))
	assert.True(t, cache.NewComparableKey(complex64(1+2i)).Equals(cache.NewComparableKey(complex64(1+2i))))
	assert.True(t, cache.NewComparableKey(complex128(1+2i)).Equals(cache.NewComparableKey(complex128(1+2i))))
	n := time.Now()
	assert.True(t, cache.NewComparableKey(n).Equals(cache.NewComparableKey(n)))
}

func TestNewStringValue(t *testing.T) {
	v := "test string"
	stringValue := cache.NewStringValue(v)
	assert.Equal(t, v, stringValue.Value())
}

func TestStringValue_Equals(t *testing.T) {
	stringValue1 := cache.NewStringValue("123")
	stringValue2 := cache.NewStringValue("123")
	stringValue3 := cache.NewStringValue("world")

	assert.True(t, stringValue1.Equals(stringValue2), "should be equal for the same string value")
	assert.False(t, stringValue1.Equals(stringValue3), "should not be equal for different string values")
	assert.False(t, stringValue1.Equals(cache.NewInt64Value(123)), "should be false when compared with a different type")
}
