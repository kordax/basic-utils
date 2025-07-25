package ucache_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kordax/basic-utils/v2/ucache"
	"github.com/kordax/basic-utils/v2/uconst"
	"github.com/kordax/basic-utils/v2/uopt"
	"github.com/kordax/basic-utils/v2/uref"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	expectedUniqueKeys := []uconst.Unique{ucache.IntKey(1), ucache.IntKey(2), ucache.IntKey(3)}
	assert.EqualValues(t, expectedUniqueKeys, key.Keys())
	assert.Equal(t, "1, 2, 3", key.String())
}

func TestStrCompositeKey_Keys_String(t *testing.T) {
	key := ucache.NewStrCompositeKey("a", "b", "c")
	expectedUniqueKeys := []uconst.Unique{ucache.IntKey(97), ucache.IntKey(98), ucache.IntKey(99)} // ASCII values of 'a', 'b', 'c'
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

func TestFarmHash64Entity(t *testing.T) {
	cache := ucache.NewInMemoryHashMapCache[ucache.FarmHash64Entity, ucache.StringValue](uopt.NullDuration())

	type nested struct {
		ExportedNestedStringField string
		ExportedNestedIntField    int
	}

	type custom struct {
		ExportedStringField string
		exportedIntField    int
		ExportedFloatField  float64
		ExportedBoolField   bool
		ExportedPointer     *int
		ExportedMapField    map[string]int
		ExportedSliceField  []string
		ExportedNested      nested
		ExportedNestedPtr   *nested
		ExportedInterface   interface{}
		unexportedField     string // This field will not be serialized
	}

	val1 := 42
	nestedVal := nested{
		ExportedNestedStringField: "nested",
		ExportedNestedIntField:    100,
	}

	key1 := custom{
		ExportedStringField: "value1",
		exportedIntField:    1,
		ExportedFloatField:  1.1,
		ExportedBoolField:   true,
		ExportedPointer:     &val1,
		ExportedMapField:    map[string]int{"one": 1, "two": 2},
		ExportedSliceField:  []string{"slice1", "slice2"},
		ExportedNested:      nestedVal,
		ExportedNestedPtr:   &nestedVal,
		ExportedInterface:   "interface1",
		unexportedField:     "unexported",
	}
	value1 := ucache.NewStringValue("data1")

	key2 := custom{
		ExportedStringField: "value2",
		exportedIntField:    2,
		ExportedFloatField:  2.2,
		ExportedBoolField:   false,
		ExportedPointer:     &val1,
		ExportedMapField:    map[string]int{"three": 3, "four": 4},
		ExportedSliceField:  []string{"slice3", "slice4"},
		ExportedNested:      nestedVal,
		ExportedNestedPtr:   &nestedVal,
		ExportedInterface:   "interface2",
		unexportedField:     "unexported",
	}
	value2 := ucache.NewStringValue("data2")

	wrappedKey1 := ucache.Hashed(key1)
	wrappedKey2 := ucache.Hashed(key2)

	cache.Set(wrappedKey1, value1)
	cache.Set(wrappedKey2, value2)

	found1, ok1 := cache.Get(wrappedKey1)
	found2, ok2 := cache.Get(wrappedKey2)

	require.True(t, ok1)
	assert.Equal(t, value1, *found1)

	require.True(t, ok2)
	assert.Equal(t, value2, *found2)
}

func TestFarmHash64EntityConsistency(t *testing.T) {
	obj := "test object"
	entity1 := ucache.Hashed(obj)
	entity2 := ucache.Hashed(obj)

	hash1 := entity1.Key()
	hash2 := entity2.Key()

	require.True(t, entity1.Equals(entity2), "Entities with the same object should be equal")
	assert.Equal(t, hash1, hash2, "Hash values should be identical for the same object")
}

func TestFarmHash64EntityInequality(t *testing.T) {
	obj1 := "test object 1"
	obj2 := "test object 2"
	entity1 := ucache.Hashed(obj1)
	entity2 := ucache.Hashed(obj2)

	// It's highly unlikely that different objects produce the same hash
	assert.NotEqual(t, entity1.Key(), entity2.Key(), "Different objects should have different hash values")
	require.False(t, entity1.Equals(entity2), "Entities with different objects should not be equal")
}

func TestFarmHash64EntityEdgeCases(t *testing.T) {
	entityNil1 := ucache.Hashed(nil)
	entityNil2 := ucache.Hashed(nil)

	require.True(t, entityNil1.Equals(entityNil2), "Entities with nil objects should be equal")
	assert.Equal(t, entityNil1.Key(), entityNil2.Key(), "Hash values should be identical for nil objects")

	obj := "pointer test"
	entityPtr1 := ucache.Hashed(&obj)
	entityPtr2 := ucache.Hashed(&obj)

	require.True(t, entityPtr1.Equals(entityPtr2), "Entities with the same pointer should be equal")
	assert.Equal(t, entityPtr1.Key(), entityPtr2.Key(), "Hash values should be identical for the same pointer")
}

func TestFarmHash64EntityCollisions(t *testing.T) {
	cache := ucache.NewInMemoryHashMapCache[*ucache.FarmHash64Entity, ucache.StringValue](uopt.NullDuration())

	rand.New(rand.NewSource(time.Now().UnixNano()))

	type nested struct {
		ExportedNestedStringField string
		ExportedNestedIntField    int
	}

	type custom struct {
		ExportedStringField string
		ExportedIntField    int
		ExportedFloatField  float64
		ExportedBoolField   bool
		ExportedPointer     *int
		ExportedMapField    map[string]int
		ExportedSliceField  []string
		ExportedNested      nested
		ExportedNestedPtr   *nested
		ExportedInterface   interface{}
		unexportedField     string
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	val1 := 42
	nestedVal := nested{
		ExportedNestedStringField: "nested",
		ExportedNestedIntField:    100,
	}

	for i := 0; i < 1000; i++ {
		key := custom{
			ExportedStringField: fmt.Sprintf("value%d", i),
			ExportedIntField:    i,
			ExportedFloatField:  r.Float64(),
			ExportedBoolField:   i%2 == 0,
			ExportedPointer:     &val1,
			ExportedMapField:    map[string]int{fmt.Sprintf("key%d", i): i},
			ExportedSliceField:  []string{fmt.Sprintf("slice%d", i)},
			ExportedNested:      nestedVal,
			ExportedNestedPtr:   &nestedVal,
			ExportedInterface:   fmt.Sprintf("interface%d", i),
			unexportedField:     fmt.Sprintf("unexported%d", i),
		}
		value := ucache.NewStringValue(fmt.Sprintf("data%d", i))

		wrappedKey := ucache.Hashed(key)
		cache.Set(&wrappedKey, value)

		found, ok := cache.Get(&wrappedKey)
		require.True(t, ok)
		assert.Equal(t, value, *found)
	}
}

func TestNewFarmHashCompositeKey(t *testing.T) {
	key1 := ucache.Hashed("key1")
	key2 := ucache.Hashed("key2")

	compositeKey := ucache.NewFarmHashCompositeKey("key1", "key2")

	require.Equal(t, 2, len(compositeKey.Keys()))

	assert.Equal(t, key1.Key(), compositeKey.Keys()[0].Key())
	assert.Equal(t, key2.Key(), compositeKey.Keys()[1].Key())
}

func TestFarmHash64CompositeKey_Equals(t *testing.T) {
	key1 := ucache.Hashed("key1")
	key2 := ucache.Hashed("key2")

	compositeKey1 := ucache.NewFarmHashCompositeKey(key1, key2)
	compositeKey2 := ucache.NewFarmHashCompositeKey(key1, key2)

	key3 := ucache.Hashed("key3")
	compositeKey3 := ucache.NewFarmHashCompositeKey(key1, key3)

	assert.True(t, compositeKey1.Equals(compositeKey2))
	assert.False(t, compositeKey1.Equals(compositeKey3))
}

func TestFarmHash64CompositeKey_Keys(t *testing.T) {
	key1 := ucache.Hashed("key1")
	key2 := ucache.Hashed("key2")

	compositeKey := ucache.NewFarmHashCompositeKey("key1", "key2")
	keys := compositeKey.Keys()

	assert.Equal(t, 2, len(keys))

	assert.Equal(t, key1.Key(), keys[0].Key())
	assert.Equal(t, key2.Key(), keys[1].Key())
}

func TestIntKey_Equals_2(t *testing.T) {
	var key1 ucache.IntKey = 123
	var key2 ucache.IntKey = 123
	key3 := ucache.IntKey(456)
	keyPtr := uref.Ref(ucache.IntKey(123))
	var nilKeyPtr *ucache.IntKey = nil

	assert.True(t, key1.Equals(key2), "Expected keys to be equal (value to value)")
	assert.True(t, key1.Equals(keyPtr), "Expected keys to be equal (value to pointer)")
	assert.True(t, keyPtr.Equals(key1), "Expected keys to be equal (pointer to value)")
	assert.True(t, keyPtr.Equals(keyPtr), "Expected keys to be equal (pointer to pointer)")
	assert.False(t, key1.Equals(key3), "Expected keys not to be equal (different values)")
	assert.False(t, key1.Equals(nilKeyPtr), "Expected keys not to be equal (value to nil pointer)")
}

func TestStringKey_Equals_2(t *testing.T) {
	key1 := ucache.StringKey("test")
	key2 := ucache.StringKey("test")
	key3 := ucache.StringKey("different")
	keyPtr := uref.Ref(ucache.StringKey("test"))
	var nilKeyPtr *ucache.StringKey = nil

	assert.True(t, key1.Equals(key2), "Expected keys to be equal (value to value)")
	assert.True(t, key1.Equals(keyPtr), "Expected keys to be equal (value to pointer)")
	assert.True(t, keyPtr.Equals(key1), "Expected keys to be equal (pointer to value)")
	assert.True(t, keyPtr.Equals(keyPtr), "Expected keys to be equal (pointer to pointer)")
	assert.False(t, key1.Equals(key3), "Expected keys not to be equal (different values)")
	assert.False(t, key1.Equals(nilKeyPtr), "Expected keys not to be equal (value to nil pointer)")
}

func TestUIntKey_Equals_2(t *testing.T) {
	key1 := ucache.UIntKey(100)
	key2 := ucache.UIntKey(100)
	key3 := ucache.UIntKey(200)
	keyPtr := uref.Ref(ucache.UIntKey(100))
	var nilKeyPtr *ucache.UIntKey = nil

	assert.True(t, key1.Equals(key2), "Expected keys to be equal (value to value)")
	assert.True(t, key1.Equals(keyPtr), "Expected keys to be equal (value to pointer)")
	assert.True(t, keyPtr.Equals(key1), "Expected keys to be equal (pointer to value)")
	assert.True(t, keyPtr.Equals(keyPtr), "Expected keys to be equal (pointer to pointer)")
	assert.False(t, key1.Equals(key3), "Expected keys not to be equal (different values)")
	assert.False(t, key1.Equals(nilKeyPtr), "Expected keys not to be equal (value to nil pointer)")
}

func TestIntCompositeKey_Equals_PointerCases(t *testing.T) {
	key1 := ucache.NewIntCompositeKey(1, 2, 3)
	key2 := ucache.NewIntCompositeKey(1, 2, 3)
	key3 := ucache.NewIntCompositeKey(4, 5, 6)
	keyPtr := &key1
	var nilKeyPtr *ucache.IntCompositeKey = nil

	assert.True(t, key1.Equals(key2), "Expected keys to be equal (value to value)")
	assert.True(t, key1.Equals(keyPtr), "Expected keys to be equal (value to pointer)")
	assert.True(t, keyPtr.Equals(key1), "Expected keys to be equal (pointer to value)")
	assert.True(t, keyPtr.Equals(keyPtr), "Expected keys to be equal (pointer to pointer)")
	assert.False(t, key1.Equals(key3), "Expected keys not to be equal (different values)")
	assert.False(t, key1.Equals(nilKeyPtr), "Expected keys not to be equal (value to nil pointer)")
}

func TestUIntCompositeKey_Equals_PointerCases(t *testing.T) {
	key1 := ucache.NewUIntCompositeKey(1, 2, 3)
	key2 := ucache.NewUIntCompositeKey(1, 2, 3)
	key3 := ucache.NewUIntCompositeKey(4, 5, 6)
	keyPtr := &key1
	var nilKeyPtr *ucache.UIntCompositeKey = nil

	assert.True(t, key1.Equals(key2), "Expected keys to be equal (value to value)")
	assert.True(t, key1.Equals(keyPtr), "Expected keys to be equal (value to pointer)")
	assert.True(t, keyPtr.Equals(key1), "Expected keys to be equal (pointer to value)")
	assert.True(t, keyPtr.Equals(keyPtr), "Expected keys to be equal (pointer to pointer)")
	assert.False(t, key1.Equals(key3), "Expected keys not to be equal (different values)")
	assert.False(t, key1.Equals(nilKeyPtr), "Expected keys not to be equal (value to nil pointer)")
}

func TestStrCompositeKey_Equals_PointerCases(t *testing.T) {
	key1 := ucache.NewStrCompositeKey("a", "b", "c")
	key2 := ucache.NewStrCompositeKey("a", "b", "c")
	key3 := ucache.NewStrCompositeKey("d", "e", "f")
	keyPtr := &key1
	var nilKeyPtr *ucache.StrCompositeKey = nil

	assert.True(t, key1.Equals(key2), "Expected keys to be equal (value to value)")
	assert.True(t, key1.Equals(keyPtr), "Expected keys to be equal (value to pointer)")
	assert.True(t, keyPtr.Equals(key1), "Expected keys to be equal (pointer to value)")
	assert.True(t, keyPtr.Equals(keyPtr), "Expected keys to be equal (pointer to pointer)")
	assert.False(t, key1.Equals(key3), "Expected keys not to be equal (different values)")
	assert.False(t, key1.Equals(nilKeyPtr), "Expected keys not to be equal (value to nil pointer)")
}
