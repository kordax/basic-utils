package ucache_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kordax/basic-utils/ucache"
	"github.com/kordax/basic-utils/uopt"
	"github.com/stretchr/testify/assert"
)

func TestStringKey_Key(t *testing.T) {
	uuid.EnableRandPool()

	value := ucache.StringKey("MyTestString")
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
			uKey := ucache.StringKey(uuid.NewString()).Key()
			assert.NotNil(t, uKey)
			assert.NotZero(t, uKey)
			assert.NotEqualValues(t, key, uKey)
			assert.NotEqualValues(t, key2, uKey)
		})
	}
}

func TestStringKey_String(t *testing.T) {
	key := ucache.StringKey("hello")
	assert.Equal(t, "hello", key.String())
}

func TestStrCompositeKey_Equals(t *testing.T) {
	key1 := ucache.NewStrCompositeKey("a", "b", "c")
	key2 := ucache.NewStrCompositeKey("a", "b", "c")
	key3 := ucache.NewStrCompositeKey("d", "e", "f")

	assert.True(t, key1.Equals(key2))
	assert.False(t, key1.Equals(key3))
}

func TestIntKey_Key(t *testing.T) {
	value := ucache.IntKey(123)
	key := value.Key()
	assert.EqualValues(t, 123, key)

	value2 := ucache.IntKey(456)
	key2 := value2.Key()
	assert.NotEqual(t, key, key2)
}

func TestUIntKey_Equals(t *testing.T) {
	key1 := ucache.UIntKey(100)
	key2 := ucache.UIntKey(100)
	key3 := ucache.UIntKey(200)

	assert.True(t, key1.Equals(key2))
	assert.False(t, key1.Equals(key3))
}

func TestUIntKey_Key(t *testing.T) {
	value := ucache.UIntKey(123)
	key := value.Key()
	assert.EqualValues(t, 123, key)

	value2 := ucache.UIntKey(456)
	key2 := value2.Key()
	assert.NotEqual(t, key, key2)
}

func TestIntCompositeKey_Keys(t *testing.T) {
	value := ucache.NewIntCompositeKey(123, 456)
	keys := value.Keys()
	assert.Len(t, keys, 2)
	assert.Contains(t, keys, ucache.IntKey(123))
	assert.Contains(t, keys, ucache.IntKey(456))

	value2 := ucache.NewIntCompositeKey(789)
	keys2 := value2.Keys()
	assert.Len(t, keys2, 1)
	assert.NotEqual(t, keys, keys2)
}

func TestIntCompositeKey_Equals(t *testing.T) {
	key1 := ucache.NewIntCompositeKey(1, 2, 3)
	key2 := ucache.NewIntCompositeKey(1, 2, 3)
	key3 := ucache.NewIntCompositeKey(4, 5, 6)

	assert.True(t, key1.Equals(key2))
	assert.False(t, key1.Equals(key3))
}

func TestUIntCompositeKey_Equals(t *testing.T) {
	key1 := ucache.NewUIntCompositeKey(1, 2, 3)
	key2 := ucache.NewUIntCompositeKey(1, 2, 3)
	key3 := ucache.NewUIntCompositeKey(4, 5, 6)

	assert.True(t, key1.Equals(key2))
	assert.False(t, key1.Equals(key3))
}

func TestUIntCompositeKey_Keys(t *testing.T) {
	key := ucache.NewUIntCompositeKey(1, 2, 3)
	assert.EqualValues(t, []int64{1, 2, 3}, key.Keys())
}

func TestGenericCompositeKey_Equals(t *testing.T) {
	key1 := ucache.NewGenericCompositeKey(ucache.StringKey("test"), ucache.IntKey(123))
	key2 := ucache.NewGenericCompositeKey(ucache.StringKey("test"), ucache.IntKey(123))
	key3 := ucache.NewGenericCompositeKey(ucache.StringKey("different"), ucache.IntKey(456))

	assert.True(t, key1.Equals(key2))
	assert.False(t, key1.Equals(key3))
}

func TestGenericCompositeKey_Keys(t *testing.T) {
	key := ucache.NewGenericCompositeKey(ucache.StringKey("test"), ucache.IntKey(123))
	assert.NotNil(t, key.Keys()) // More specific tests might depend on the hashing algorithm
}

func TestUIntCompositeKey_Keys_String(t *testing.T) {
	key := ucache.NewUIntCompositeKey(1, 2, 3)
	assert.EqualValues(t, []int64{1, 2, 3}, key.Keys())
	assert.Equal(t, "1, 2, 3", key.String())
}

func TestIntCompositeKey_Keys_String(t *testing.T) {
	key := ucache.NewIntCompositeKey(1, 2, 3)
	expectedUniqueKeys := []ucache.Unique{ucache.IntKey(1), ucache.IntKey(2), ucache.IntKey(3)}
	assert.EqualValues(t, expectedUniqueKeys, key.Keys())
	assert.Equal(t, "1, 2, 3", key.String())
}

func TestStrCompositeKey_Keys_String(t *testing.T) {
	key := ucache.NewStrCompositeKey("a", "b", "c")
	expectedUniqueKeys := []ucache.Unique{ucache.IntKey(97), ucache.IntKey(98), ucache.IntKey(99)} // ASCII values of 'a', 'b', 'c'
	assert.EqualValues(t, expectedUniqueKeys, key.Keys())
	assert.Equal(t, "a, b, c", key.String())
}

func TestNewStringSliceValue(t *testing.T) {
	input := []string{"c", "a", "b"}
	sorted := []string{"a", "b", "c"}
	value := ucache.NewStringSliceValue(input)

	assert.Equal(t, sorted, value.Values())
}

func TestStringSliceValue_Values(t *testing.T) {
	input := []string{"a", "b", "c"}
	value := ucache.NewStringSliceValue(input)

	assert.Equal(t, input, value.Values())
}

func TestStringSliceValue_Equals(t *testing.T) {
	value1 := ucache.NewStringSliceValue([]string{"a", "b", "c"})
	value2 := ucache.NewStringSliceValue([]string{"a", "b", "c"})
	value3 := ucache.NewStringSliceValue([]string{"d", "e", "f"})

	assert.True(t, value1.Equals(value2))
	assert.False(t, value1.Equals(value3))
	assert.False(t, value1.Equals(ucache.StringKey("test string")))
}

func TestComparableSlice_Equals(t *testing.T) {
	slice1 := ucache.ComparableSlice[ucache.IntKey]{Data: []ucache.IntKey{1, 2, 3}}
	slice2 := ucache.ComparableSlice[ucache.IntKey]{Data: []ucache.IntKey{1, 2, 3}}
	slice3 := ucache.ComparableSlice[ucache.IntKey]{Data: []ucache.IntKey{4, 5, 6}}

	assert.True(t, slice1.Equals(slice2))
	assert.False(t, slice1.Equals(slice3))
}

func TestComparableSlice_Equals_NotComparableSlice(t *testing.T) {
	c := ucache.ComparableSlice[ucache.IntKey]{Data: []ucache.IntKey{1, 2, 3}}
	other := ucache.NewComparableKey[int](1)

	assert.False(t, c.Equals(other))
}

func TestComparableSlice_Equals_Case2(t *testing.T) {
	slice1 := ucache.ComparableSlice[DummyComparable]{
		Data: []DummyComparable{
			{1},
			{2},
			{3},
			{4},
		},
	}
	slice2 := ucache.ComparableSlice[DummyComparable]{
		Data: []DummyComparable{
			{1},
			{2},
			{3},
			{4},
		},
	}
	assert.EqualValues(t, slice1, slice2)
	slice2 = ucache.ComparableSlice[DummyComparable]{
		Data: []DummyComparable{
			{1},
			{2},
			{5},
			{4},
		},
	}
	assert.NotEqualValues(t, slice1, slice2)
	slice2 = ucache.ComparableSlice[DummyComparable]{
		Data: []DummyComparable{
			{1},
			{2},
			{4},
			{3},
		},
	}
	assert.NotEqualValues(t, slice1, slice2)
	slice2 = ucache.ComparableSlice[DummyComparable]{
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
	key1 := ucache.NewComparableKey(123)
	key2 := ucache.NewComparableKey(123)
	key3 := ucache.NewComparableKey(456)

	assert.True(t, key1.Equals(key2), "keys with the same value should be equal")
	assert.False(t, key1.Equals(key3), "keys with different values should not be equal")
}

func TestComparableKey_String(t *testing.T) {
	key := ucache.NewComparableKey(123)
	assert.Equal(t, "123", key.String(), "String() should return the correct string representation")

	key2 := ucache.NewComparableKey("abc")
	assert.Equal(t, "abc", key2.String(), "String() should handle string types correctly")
}

func TestComparableKey_Equals_CoverAllTypes(t *testing.T) {
	assert.True(t, ucache.NewComparableKey("abc").Equals(ucache.NewComparableKey("abc")))
	assert.True(t, ucache.NewComparableKey(42).Equals(ucache.NewComparableKey(42)))
	assert.True(t, ucache.NewComparableKey(int8(42)).Equals(ucache.NewComparableKey(int8(42))))
	assert.True(t, ucache.NewComparableKey(int16(42)).Equals(ucache.NewComparableKey(int16(42))))
	assert.True(t, ucache.NewComparableKey(int32(42)).Equals(ucache.NewComparableKey(int32(42))))
	assert.True(t, ucache.NewComparableKey(int64(42)).Equals(ucache.NewComparableKey(int64(42))))
	assert.True(t, ucache.NewComparableKey(uint(42)).Equals(ucache.NewComparableKey(uint(42))))
	assert.True(t, ucache.NewComparableKey(uint8(42)).Equals(ucache.NewComparableKey(uint8(42))))
	assert.True(t, ucache.NewComparableKey(uint16(42)).Equals(ucache.NewComparableKey(uint16(42))))
	assert.True(t, ucache.NewComparableKey(uint32(42)).Equals(ucache.NewComparableKey(uint32(42))))
	assert.True(t, ucache.NewComparableKey(uint64(42)).Equals(ucache.NewComparableKey(uint64(42))))
	assert.True(t, ucache.NewComparableKey(true).Equals(ucache.NewComparableKey(true)))
	assert.True(t, ucache.NewComparableKey(false).Equals(ucache.NewComparableKey(false)))
	assert.True(t, ucache.NewComparableKey(float32(42.1)).Equals(ucache.NewComparableKey(float32(42.1))))
	assert.True(t, ucache.NewComparableKey(float64(42.1)).Equals(ucache.NewComparableKey(float64(42.1))))
	assert.True(t, ucache.NewComparableKey(complex64(1+2i)).Equals(ucache.NewComparableKey(complex64(1+2i))))
	assert.True(t, ucache.NewComparableKey(complex128(1+2i)).Equals(ucache.NewComparableKey(complex128(1+2i))))
	n := time.Now()
	assert.True(t, ucache.NewComparableKey(n).Equals(ucache.NewComparableKey(n)))
}

func TestNewStringValue(t *testing.T) {
	v := "test string"
	stringValue := ucache.NewStringValue(v)
	assert.Equal(t, v, stringValue.Value())
}

func TestStringValue_Equals(t *testing.T) {
	stringValue1 := ucache.NewStringValue("123")
	stringValue2 := ucache.NewStringValue("123")
	stringValue3 := ucache.NewStringValue("world")

	assert.True(t, stringValue1.Equals(stringValue2), "should be equal for the same string value")
	assert.False(t, stringValue1.Equals(stringValue3), "should not be equal for different string values")
	assert.False(t, stringValue1.Equals(ucache.NewInt64Value(123)), "should be false when compared with a different type")
}

func TestKeysWithCache(t *testing.T) {
	cache := ucache.NewInMemoryTreeMultiCache[ucache.CompositeKey, ucache.StringValue](uopt.NullDuration())

	intKey := ucache.IntKey(42)
	cache.Put(intKey, ucache.NewStringValue("value for int"))
	retrievedInt := cache.Get(intKey)
	assert.Equal(t, []ucache.StringValue{ucache.NewStringValue("value for int")}, retrievedInt, "IntKey retrieval failed")

	uintKey := ucache.UIntKey(42)
	cache.Set(uintKey, ucache.NewStringValue("value for uint"))
	retrievedUInt := cache.Get(uintKey)
	assert.Equal(t, []ucache.StringValue{ucache.NewStringValue("value for uint")}, retrievedUInt, "UIntKey retrieval failed")

	stringKey := ucache.StringKey("key")
	cache.Put(stringKey, ucache.NewStringValue("value for string"))
	retrievedString := cache.Get(stringKey)
	assert.Equal(t, []ucache.StringValue{ucache.NewStringValue("value for string")}, retrievedString, "StringKey retrieval failed")
}
