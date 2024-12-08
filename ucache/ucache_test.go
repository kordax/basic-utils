/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ucache_test

import (
	"sync"
	"testing"
	"time"

	"github.com/kordax/basic-utils/ucache"
	"github.com/kordax/basic-utils/uopt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// StatusType is a custom type representing the status.
type StatusType string

// Predefined statuses for StatusType.
const (
	StatusActive   StatusType = "active"
	StatusInactive StatusType = "inactive"
	StatusPending  StatusType = "pending"
)

// NestedMeta is a nested comparable struct within ComplexComparableStruct.
type NestedMeta struct {
	Category string
	Level    int
}

// Coordinates represents geographical coordinates.
type Coordinates struct {
	Latitude  float64
	Longitude float64
}

// ComplexComparableStruct is a complex struct used for unit testing.
// It includes various comparable field types, nested structs, arrays, and custom types.
type ComplexComparableStruct struct {
	ID               int
	Name             string
	Timestamp        time.Time
	Active           bool
	Scores           [5]int
	Meta             NestedMeta
	Location         Coordinates
	Permissions      [3]bool
	Rating           float32
	Status           StatusType
	RetryCounts      [2]int8
	Flags            [4]uint8
	Attributes       [3]string
	Dimension        Dimension
	DeviceStatus     DeviceStatus
	RegistrationData RegistrationData
	OptionalValue    *int
	OptionalCategory *string
}

// Dimension is another nested comparable struct.
type Dimension struct {
	Width  float64
	Height float64
	Depth  float64
}

// DeviceStatus represents the status of a device with multiple fields.
type DeviceStatus struct {
	Online       bool
	BatteryLevel uint8
	ErrorCode    uint16
}

// RegistrationData holds registration details in a nested struct.
type RegistrationData struct {
	RegisteredAt time.Time
	Verified     bool
	Notes        [2]string
}

func TestHashMapCache_CompositeKey(t *testing.T) {
	c := ucache.NewInMemoryHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
	key := ucache.StringKey("category")
	key2 := ucache.StringKey("category2")
	val := 10
	val2 := 236261

	c.SetQuietly(key, val)
	c.SetQuietly(key2, val2)

	result, ok := c.Get(key)
	require.True(t, ok, "value was expected to be cached")
	result2, ok := c.Get(key2)
	require.True(t, ok, "value was expected to be cached")
	assert.EqualValues(t, val, *result)
	assert.EqualValues(t, val2, *result2)
}

func TestHashMapCache_SetNil(t *testing.T) {
	c := ucache.NewInMemoryHashMapCache[ucache.StringKey, *int](uopt.Null[time.Duration]())
	key := ucache.StringKey("category")
	key2 := ucache.StringKey("category2")
	key3 := ucache.StringKey("category3")
	val := 10
	val2 := 236261
	var val3 *int = nil

	c.Set(key, &val)
	c.Set(key2, &val2)
	c.Set(key3, val3)

	result, ok := c.Get(key)
	require.True(t, ok, "value was expected to be cached")
	result2, ok := c.Get(key2)
	require.True(t, ok, "value was expected to be cached")
	result3, ok := c.Get(key3)
	require.True(t, ok, "value was expected to be cached")
	assert.EqualValues(t, val, **result)
	assert.EqualValues(t, val2, **result2)
	assert.EqualValues(t, val3, *result3)
}

func TestHashMapCache_PutQuietly(t *testing.T) {
	c := ucache.NewInMemoryHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
	key := ucache.StringKey("kp_1")
	val := 10
	val2 := 15

	c.SetQuietly(key, val)
	c.SetQuietly(key, val)
	c.SetQuietly(key, val)

	result, ok := c.Get(key)
	require.True(t, ok, "value was expected to be cached")
	assert.Equal(t, *result, val)

	c.SetQuietly(key, val2)
	result, ok = c.Get(key)
	require.True(t, ok, "value was expected to be cached")
	assert.Equal(t, *result, val2)

	c.SetQuietly(key, val)
	result, ok = c.Get(key)
	require.True(t, ok, "value was expected to be cached")
	assert.Equal(t, *result, val)
}

func TestHashMapCache_TTLExpiry(t *testing.T) {
	ttl := 100 * time.Millisecond
	c := ucache.NewInMemoryHashMapCache[ucache.StringKey, int](uopt.Of(ttl))
	key := ucache.StringKey("ttlKey")
	val := 42

	c.Set(key, val)
	time.Sleep(2 * ttl)
	outdated := c.Outdated(uopt.Of(key))
	assert.True(t, outdated, "key should be marked as outdated")
}

func TestHashMapCache_Concurrency(t *testing.T) {
	c := ucache.NewInMemoryHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
	key := ucache.StringKey("concurrencyKey")
	val := 42

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Set(key, val+i)
			result, _ := c.Get(key)
			assert.NotNil(t, result)
		}(i)
	}
	wg.Wait()
}

func TestHashMapCache_EmptyCache(t *testing.T) {
	c := ucache.NewInMemoryHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
	key := ucache.StringKey("emptyKey")

	_, ok := c.Get(key)
	assert.False(t, ok, "key should not be found in an empty cache")

	c.DropKey(key)
	_, ok = c.Get(key)
	assert.False(t, ok, "key should not be found in an empty cache")
}

func TestHashMapCache_DropAll(t *testing.T) {
	c := ucache.NewInMemoryHashMapCache[ucache.StringKey, int](uopt.Null[time.Duration]())
	key := ucache.StringKey("key1")
	key2 := ucache.StringKey("key2")
	c.Set(key, 1)
	c.Set(key2, 2)

	c.Drop()
	_, ok1 := c.Get(key)
	_, ok2 := c.Get(key2)

	assert.False(t, ok1, "key1 should be dropped")
	assert.False(t, ok2, "key2 should be dropped")
}

func TestInMemoryHashMapCache(t *testing.T) {
	cache := ucache.NewInMemoryHashMapCache[ucache.IntKey, string](uopt.Null[time.Duration]())

	// Define multiple keys and values
	key1 := ucache.IntKey(1)
	value1 := "MyValue"

	key2 := ucache.IntKey(2)
	value2 := "AnotherValue"

	key3 := ucache.IntKey(3)
	value3 := "ThirdValue"

	key4 := ucache.IntKey(4)
	value4 := "FourthValue"

	key5 := ucache.IntKey(5)
	value5 := "FifthValue"

	// Test setting and getting multiple keys
	cache.Set(key1, value1)
	cache.Set(key2, value2)
	cache.Set(key3, value3)
	cache.Set(key4, value4)
	cache.Set(key5, value5)

	// Verify all keys return correct values
	retrievedValue, ok := cache.Get(key1)
	require.True(t, ok, "Expected to retrieve value for key1")
	assert.Equal(t, value1, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key2)
	require.True(t, ok, "Expected to retrieve value for key2")
	assert.Equal(t, value2, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key3)
	require.True(t, ok, "Expected to retrieve value for key3")
	assert.Equal(t, value3, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key4)
	require.True(t, ok, "Expected to retrieve value for key4")
	assert.Equal(t, value4, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key5)
	require.True(t, ok, "Expected to retrieve value for key5")
	assert.Equal(t, value5, *retrievedValue, "Retrieved value should match the set value")

	// Test updating values for existing keys
	updatedValue1 := "UpdatedMyValue"
	cache.Set(key1, updatedValue1)
	retrievedValue, ok = cache.Get(key1)
	require.True(t, ok, "Expected to retrieve updated value for key1")
	assert.Equal(t, updatedValue1, *retrievedValue, "Retrieved value should match the updated value")

	// Test removing keys
	cache.DropKey(key1)
	retrievedValue, ok = cache.Get(key1)
	assert.False(t, ok, "Expected key1 to be removed from cache")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key1 should be nil")

	// Ensure other keys are still retrievable and correct after removing key1
	retrievedValue, ok = cache.Get(key2)
	require.True(t, ok, "Expected to retrieve value for key2")
	assert.Equal(t, value2, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key3)
	require.True(t, ok, "Expected to retrieve value for key3")
	assert.Equal(t, value3, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key4)
	require.True(t, ok, "Expected to retrieve value for key4")
	assert.Equal(t, value4, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key5)
	require.True(t, ok, "Expected to retrieve value for key5")
	assert.Equal(t, value5, *retrievedValue, "Retrieved value should match the set value")

	// Test SetQuietly
	cache.SetQuietly(key1, updatedValue1)
	retrievedValue, ok = cache.Get(key1)
	require.True(t, ok, "Expected to retrieve value for key1 after SetQuietly")
	assert.Equal(t, updatedValue1, *retrievedValue, "Retrieved value should match the set value")

	// Test Drop (clearing the entire cache)
	cache.Drop()
	retrievedValue, ok = cache.Get(key1)
	assert.False(t, ok, "Expected key1 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key1 should be nil after Drop")

	retrievedValue, ok = cache.Get(key2)
	assert.False(t, ok, "Expected key2 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key2 should be nil after Drop")

	retrievedValue, ok = cache.Get(key3)
	assert.False(t, ok, "Expected key3 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key3 should be nil after Drop")

	retrievedValue, ok = cache.Get(key4)
	assert.False(t, ok, "Expected key4 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key4 should be nil after Drop")

	retrievedValue, ok = cache.Get(key5)
	assert.False(t, ok, "Expected key5 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key5 should be nil after Drop")
}

func TestHashMapCacheHighCollisionProbability(t *testing.T) {
	c := ucache.NewInMemoryHashMapCache[CollisionTestKey, ucache.Int64Value](uopt.Null[time.Duration]())

	// Define a set of keys that all produce the same hash code
	keys := []CollisionTestKey{
		{id: 1, hash: []int64{1, 2, 3}},
		{id: 2, hash: []int64{1, 2, 3}},
		{id: 3, hash: []int64{1, 2, 3}},
	}

	// Add values to the c for each key
	for i, key := range keys {
		c.Set(key, ucache.NewInt64Value(int64(i)))
	}

	// Ensure that all values can be retrieved despite the high collision probability
	for i, key := range keys {
		value, ok := c.Get(key)
		require.True(t, ok, "Expected to retrieve value for key")
		assert.EqualValues(t, ucache.NewInt64Value(int64(i)), *value, "Expected value for key %s", key.String())
	}
}

func TestComparableMapCache_CompositeKey(t *testing.T) {
	// Using string as a comparable key
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())
	key := "category"
	key2 := "category2"
	val := 10
	val2 := 236261

	c.SetQuietly(key, val)
	c.SetQuietly(key2, val2)

	result, ok := c.Get(key)
	require.True(t, ok, "value was expected to be cached")
	result2, ok := c.Get(key2)
	require.True(t, ok, "value was expected to be cached")
	assert.EqualValues(t, val, *result)
	assert.EqualValues(t, val2, *result2)
}

func TestComparableMapCache_SetNil(t *testing.T) {
	// Using string as a comparable key and *int as value type
	c := ucache.NewInMemoryComparableMapCache[string, *int](uopt.Null[time.Duration]())
	key := "category"
	key2 := "category2"
	key3 := "category3"
	val := 10
	val2 := 236261
	var val3 *int = nil

	c.Set(key, &val)
	c.Set(key2, &val2)
	c.Set(key3, val3)

	result, ok := c.Get(key)
	require.True(t, ok, "value was expected to be cached")
	result2, ok := c.Get(key2)
	require.True(t, ok, "value was expected to be cached")
	result3, ok := c.Get(key3)
	require.True(t, ok, "value was expected to be cached")
	assert.EqualValues(t, val, **result)
	assert.EqualValues(t, val2, **result2)
	assert.EqualValues(t, val3, *result3)
}

func TestComparableMapCache_PutQuietly(t *testing.T) {
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())
	key := "kp_1"
	val := 10
	val2 := 15

	c.SetQuietly(key, val)
	c.SetQuietly(key, val)
	c.SetQuietly(key, val)

	result, ok := c.Get(key)
	require.True(t, ok, "value was expected to be cached")
	assert.Equal(t, *result, val)

	c.SetQuietly(key, val2)
	result, ok = c.Get(key)
	require.True(t, ok, "value was expected to be cached")
	assert.Equal(t, *result, val2)

	c.SetQuietly(key, val)
	result, ok = c.Get(key)
	require.True(t, ok, "value was expected to be cached")
	assert.Equal(t, *result, val)
}

func TestComparableMapCache_TTLExpiry(t *testing.T) {
	ttl := 100 * time.Millisecond
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Of(ttl))
	key := "ttlKey"
	val := 42

	c.Set(key, val)
	time.Sleep(2 * ttl)
	outdated := c.Outdated(uopt.Of(key))
	assert.True(t, outdated, "key should be marked as outdated")
}

func TestComparableMapCache_Concurrency(t *testing.T) {
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())
	key := "concurrencyKey"
	val := 42

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Set(key, val+i)
			result, _ := c.Get(key)
			assert.NotNil(t, result)
		}(i)
	}
	wg.Wait()
}

func TestComparableMapCache_EmptyCache(t *testing.T) {
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())
	key := "emptyKey"

	_, ok := c.Get(key)
	assert.False(t, ok, "key should not be found in an empty cache")

	c.DropKey(key)
	_, ok = c.Get(key)
	assert.False(t, ok, "key should not be found in an empty cache")
}

func TestComparableMapCache_DropAll(t *testing.T) {
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())
	key := "key1"
	key2 := "key2"
	c.Set(key, 1)
	c.Set(key2, 2)

	c.Drop()
	_, ok1 := c.Get(key)
	_, ok2 := c.Get(key2)

	assert.False(t, ok1, "key1 should be dropped")
	assert.False(t, ok2, "key2 should be dropped")
}

func TestInMemoryComparableMapCache(t *testing.T) {
	cache := ucache.NewInMemoryComparableMapCache[string, string](uopt.Null[time.Duration]())

	// Define multiple keys and values
	key1 := "key1"
	value1 := "MyValue"

	key2 := "key2"
	value2 := "AnotherValue"

	key3 := "key3"
	value3 := "ThirdValue"

	key4 := "key4"
	value4 := "FourthValue"

	key5 := "key5"
	value5 := "FifthValue"

	// Test setting and getting multiple keys
	cache.Set(key1, value1)
	cache.Set(key2, value2)
	cache.Set(key3, value3)
	cache.Set(key4, value4)
	cache.Set(key5, value5)

	// Verify all keys return correct values
	retrievedValue, ok := cache.Get(key1)
	require.True(t, ok, "Expected to retrieve value for key1")
	assert.Equal(t, value1, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key2)
	require.True(t, ok, "Expected to retrieve value for key2")
	assert.Equal(t, value2, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key3)
	require.True(t, ok, "Expected to retrieve value for key3")
	assert.Equal(t, value3, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key4)
	require.True(t, ok, "Expected to retrieve value for key4")
	assert.Equal(t, value4, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key5)
	require.True(t, ok, "Expected to retrieve value for key5")
	assert.Equal(t, value5, *retrievedValue, "Retrieved value should match the set value")

	// Test updating values for existing keys
	updatedValue1 := "UpdatedMyValue"
	cache.Set(key1, updatedValue1)
	retrievedValue, ok = cache.Get(key1)
	require.True(t, ok, "Expected to retrieve updated value for key1")
	assert.Equal(t, updatedValue1, *retrievedValue, "Retrieved value should match the updated value")

	// Test removing keys
	cache.DropKey(key1)
	retrievedValue, ok = cache.Get(key1)
	assert.False(t, ok, "Expected key1 to be removed from cache")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key1 should be nil")

	// Ensure other keys are still retrievable and correct after removing key1
	retrievedValue, ok = cache.Get(key2)
	require.True(t, ok, "Expected to retrieve value for key2")
	assert.Equal(t, value2, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key3)
	require.True(t, ok, "Expected to retrieve value for key3")
	assert.Equal(t, value3, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key4)
	require.True(t, ok, "Expected to retrieve value for key4")
	assert.Equal(t, value4, *retrievedValue, "Retrieved value should match the set value")

	retrievedValue, ok = cache.Get(key5)
	require.True(t, ok, "Expected to retrieve value for key5")
	assert.Equal(t, value5, *retrievedValue, "Retrieved value should match the set value")

	// Test SetQuietly
	cache.SetQuietly(key1, updatedValue1)
	retrievedValue, ok = cache.Get(key1)
	require.True(t, ok, "Expected to retrieve value for key1 after SetQuietly")
	assert.Equal(t, updatedValue1, *retrievedValue, "Retrieved value should match the set value")

	// Test Drop (clearing the entire cache)
	cache.Drop()
	retrievedValue, ok = cache.Get(key1)
	assert.False(t, ok, "Expected key1 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key1 should be nil after Drop")

	retrievedValue, ok = cache.Get(key2)
	assert.False(t, ok, "Expected key2 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key2 should be nil after Drop")

	retrievedValue, ok = cache.Get(key3)
	assert.False(t, ok, "Expected key3 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key3 should be nil after Drop")

	retrievedValue, ok = cache.Get(key4)
	assert.False(t, ok, "Expected key4 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key4 should be nil after Drop")

	retrievedValue, ok = cache.Get(key5)
	assert.False(t, ok, "Expected key5 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key5 should be nil after Drop")
}

func TestInMemoryComparableMapCacheWithComplexStruct(t *testing.T) {
	cache := ucache.NewInMemoryComparableMapCache[string, ComplexComparableStruct](uopt.Null[time.Duration]())

	key1 := "key1"
	value1 := ComplexComparableStruct{
		ID:        1,
		Name:      "Test Entity 1",
		Timestamp: time.Now(),
		Active:    true,
		Scores:    [5]int{90, 85, 95, 80, 100},
		Meta: NestedMeta{
			Category: "A",
			Level:    2,
		},
		Location: Coordinates{
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
		Permissions: [3]bool{true, false, true},
		Rating:      4.5,
		Status:      StatusActive,
		RetryCounts: [2]int8{1, 2},
		Flags:       [4]uint8{1, 0, 1, 1},
		Attributes:  [3]string{"fast", "reliable", "scalable"},
		Dimension: Dimension{
			Width:  1920.0,
			Height: 1080.0,
			Depth:  500.0,
		},
		DeviceStatus: DeviceStatus{
			Online:       true,
			BatteryLevel: 85,
			ErrorCode:    0,
		},
		RegistrationData: RegistrationData{
			RegisteredAt: time.Now().Add(-24 * time.Hour),
			Verified:     true,
			Notes:        [2]string{"Initial registration", "No issues"},
		},
		OptionalValue:    nil,
		OptionalCategory: nil,
	}

	key2 := "key2"
	value2 := ComplexComparableStruct{
		ID:        2,
		Name:      "Test Entity 2",
		Timestamp: time.Now(),
		Active:    false,
		Scores:    [5]int{80, 75, 85, 70, 90},
		Meta: NestedMeta{
			Category: "B",
			Level:    1,
		},
		Location: Coordinates{
			Latitude:  40.7128,
			Longitude: -74.0060,
		},
		Permissions: [3]bool{false, true, false},
		Rating:      3.5,
		Status:      StatusInactive,
		RetryCounts: [2]int8{0, 1},
		Flags:       [4]uint8{0, 1, 0, 0},
		Attributes:  [3]string{"slow", "unreliable", "non-scalable"},
		Dimension: Dimension{
			Width:  1280.0,
			Height: 720.0,
			Depth:  300.0,
		},
		DeviceStatus: DeviceStatus{
			Online:       false,
			BatteryLevel: 50,
			ErrorCode:    404,
		},
		RegistrationData: RegistrationData{
			RegisteredAt: time.Now().Add(-48 * time.Hour),
			Verified:     false,
			Notes:        [2]string{"Delayed verification", "Issues encountered"},
		},
		OptionalValue:    nil,
		OptionalCategory: nil,
	}

	key3 := "key3"
	value3 := ComplexComparableStruct{
		ID:        3,
		Name:      "Test Entity 3",
		Timestamp: time.Now(),
		Active:    true,
		Scores:    [5]int{100, 100, 100, 100, 100},
		Meta: NestedMeta{
			Category: "C",
			Level:    3,
		},
		Location: Coordinates{
			Latitude:  34.0522,
			Longitude: -118.2437,
		},
		Permissions: [3]bool{true, true, true},
		Rating:      5.0,
		Status:      StatusPending,
		RetryCounts: [2]int8{2, 3},
		Flags:       [4]uint8{1, 1, 1, 1},
		Attributes:  [3]string{"efficient", "robust", "scalable"},
		Dimension: Dimension{
			Width:  2560.0,
			Height: 1440.0,
			Depth:  600.0,
		},
		DeviceStatus: DeviceStatus{
			Online:       true,
			BatteryLevel: 100,
			ErrorCode:    0,
		},
		RegistrationData: RegistrationData{
			RegisteredAt: time.Now().Add(-72 * time.Hour),
			Verified:     true,
			Notes:        [2]string{"Verified", "All checks passed"},
		},
		OptionalValue:    nil,
		OptionalCategory: nil,
	}

	cache.Set(key1, value1)
	cache.Set(key2, value2)
	cache.Set(key3, value3)

	retrievedValue, ok := cache.Get(key1)
	require.True(t, ok, "Expected to retrieve value for key1")
	assert.Equal(t, value1, *retrievedValue, "Retrieved value should match the set value for key1")

	retrievedValue, ok = cache.Get(key2)
	require.True(t, ok, "Expected to retrieve value for key2")
	assert.Equal(t, value2, *retrievedValue, "Retrieved value should match the set value for key2")

	retrievedValue, ok = cache.Get(key3)
	require.True(t, ok, "Expected to retrieve value for key3")
	assert.Equal(t, value3, *retrievedValue, "Retrieved value should match the set value for key3")

	updatedValue1 := ComplexComparableStruct{
		ID:        1,
		Name:      "Updated Test Entity 1",
		Timestamp: value1.Timestamp.Add(1 * time.Hour),
		Active:    false,
		Scores:    [5]int{85, 80, 90, 75, 95},
		Meta: NestedMeta{
			Category: "A1",
			Level:    3,
		},
		Location: Coordinates{
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
		Permissions: [3]bool{false, false, true},
		Rating:      4.0,
		Status:      StatusInactive,
		RetryCounts: [2]int8{2, 3},
		Flags:       [4]uint8{0, 0, 1, 1},
		Attributes:  [3]string{"moderate", "reliable", "scalable"},
		Dimension: Dimension{
			Width:  1920.0,
			Height: 1080.0,
			Depth:  500.0,
		},
		DeviceStatus: DeviceStatus{
			Online:       false,
			BatteryLevel: 75,
			ErrorCode:    1,
		},
		RegistrationData: RegistrationData{
			RegisteredAt: value1.RegistrationData.RegisteredAt.Add(1 * time.Hour),
			Verified:     false,
			Notes:        [2]string{"Re-registered", "Minor issues"},
		},
		OptionalValue:    nil,
		OptionalCategory: nil,
	}

	cache.Set(key1, updatedValue1)
	retrievedValue, ok = cache.Get(key1)
	require.True(t, ok, "Expected to retrieve updated value for key1")
	assert.Equal(t, updatedValue1, *retrievedValue, "Retrieved value should match the updated value for key1")

	cache.DropKey(key1)
	retrievedValue, ok = cache.Get(key1)
	assert.False(t, ok, "Expected key1 to be removed from cache")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key1 should be nil")

	retrievedValue, ok = cache.Get(key2)
	require.True(t, ok, "Expected to retrieve value for key2")
	assert.Equal(t, value2, *retrievedValue, "Retrieved value should match the set value for key2")

	retrievedValue, ok = cache.Get(key3)
	require.True(t, ok, "Expected to retrieve value for key3")
	assert.Equal(t, value3, *retrievedValue, "Retrieved value should match the set value for key3")

	quietUpdatedValue2 := ComplexComparableStruct{
		ID:        2,
		Name:      "Quietly Updated Test Entity 2",
		Timestamp: value2.Timestamp.Add(2 * time.Hour),
		Active:    true,
		Scores:    [5]int{95, 90, 100, 85, 105},
		Meta: NestedMeta{
			Category: "B1",
			Level:    2,
		},
		Location: Coordinates{
			Latitude:  40.7128,
			Longitude: -74.0060,
		},
		Permissions: [3]bool{true, true, false},
		Rating:      4.8,
		Status:      StatusPending,
		RetryCounts: [2]int8{3, 4},
		Flags:       [4]uint8{1, 1, 0, 0},
		Attributes:  [3]string{"efficient", "reliable", "scalable"},
		Dimension: Dimension{
			Width:  1280.0,
			Height: 720.0,
			Depth:  300.0,
		},
		DeviceStatus: DeviceStatus{
			Online:       true,
			BatteryLevel: 90,
			ErrorCode:    2,
		},
		RegistrationData: RegistrationData{
			RegisteredAt: value2.RegistrationData.RegisteredAt.Add(2 * time.Hour),
			Verified:     true,
			Notes:        [2]string{"Verified again", "No new issues"},
		},
		OptionalValue:    nil,
		OptionalCategory: nil,
	}

	cache.SetQuietly(key2, quietUpdatedValue2)
	retrievedValue, ok = cache.Get(key2)
	require.True(t, ok, "Expected to retrieve value for key2 after SetQuietly")
	assert.Equal(t, quietUpdatedValue2, *retrievedValue, "Retrieved value should match the set value for key2 after SetQuietly")

	cache.Drop()
	retrievedValue, ok = cache.Get(key2)
	assert.False(t, ok, "Expected key2 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key2 should be nil after Drop")

	retrievedValue, ok = cache.Get(key3)
	assert.False(t, ok, "Expected key3 to be removed from cache after Drop")
	assert.Nil(t, retrievedValue, "Retrieved value for removed key3 should be nil after Drop")
}

func TestComparableMapCacheHighCollisionProbability(t *testing.T) {
	// Since Go's native maps handle collisions internally, this test will focus on ensuring
	// that multiple distinct keys can coexist and be retrieved correctly.

	c := ucache.NewInMemoryComparableMapCache[int, ucache.Int64Value](uopt.Null[time.Duration]())

	// Define a set of keys that are distinct
	keys := []int{1, 2, 3, 4, 5}

	// Add values to the cache for each key
	for i, key := range keys {
		c.Set(key, ucache.NewInt64Value(int64(i)))
	}

	// Ensure that all values can be retrieved correctly
	for i, key := range keys {
		value, ok := c.Get(key)
		require.True(t, ok, "Expected to retrieve value for key %d", key)
		assert.EqualValues(t, ucache.NewInt64Value(int64(i)), *value, "Expected value for key %d", key)
	}
}

func TestComparableMapCache_MultipleSetAndGet(t *testing.T) {
	// Test multiple sets and gets to ensure consistency
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())

	keys := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
	values := []int{1, 2, 3, 4, 5}

	// Set multiple keys
	for i, key := range keys {
		c.Set(key, values[i])
	}

	// Get and verify multiple keys
	for i, key := range keys {
		val, ok := c.Get(key)
		require.True(t, ok, "Expected to find key %s", key)
		assert.Equal(t, values[i], *val, "Expected value %d for key %s", values[i], key)
	}
}

func TestComparableMapCache_ChangeTracking(t *testing.T) {
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())

	c.Set("key1", 100)
	c.Set("key2", 200)
	c.SetQuietly("key3", 300)

	changes := c.Changes()
	expected := []string{"key1", "key2"}
	assert.ElementsMatch(t, expected, changes, "Changes should include only key1 and key2")

	c.Set("key4", 400)
	changes = c.Changes()
	expected = append(expected, "key4")
	assert.ElementsMatch(t, expected, changes, "Changes should include key4")
}

func TestComparableMapCache_Outdated_NoTTL(t *testing.T) {
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())
	c.Set("key1", 1)

	// Without TTL, no key should be outdated
	outdated := c.Outdated(uopt.Of("key1"))
	assert.False(t, outdated, "With no TTL set, key should not be outdated")
}

func TestComparableMapCache_Outdated_WithTTL(t *testing.T) {
	ttl := 500 * time.Millisecond
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Of(ttl))
	c.Set("key1", 1)

	// Immediately, key should not be outdated
	outdated := c.Outdated(uopt.Of("key1"))
	assert.False(t, outdated, "Key should not be outdated immediately after setting")

	// After TTL, key should be outdated
	time.Sleep(600 * time.Millisecond)
	outdated = c.Outdated(uopt.Of("key1"))
	assert.True(t, outdated, "Key should be outdated after TTL")

	// Key that doesn't exist should be considered outdated
	outdated = c.Outdated(uopt.Of("nonexistent"))
	assert.True(t, outdated, "Non-existent key should be considered outdated")
}

func TestComparableMapCache_Outdated_PartialTTL(t *testing.T) {
	ttl := 1 * time.Second
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Of(ttl))
	c.Set("key1", 1)
	c.Set("key2", 2)

	// Wait for half the TTL
	time.Sleep(500 * time.Millisecond)

	// Both keys should not be outdated
	outdated1 := c.Outdated(uopt.Of("key1"))
	outdated2 := c.Outdated(uopt.Of("key2"))
	assert.False(t, outdated1, "Key1 should not be outdated yet")
	assert.False(t, outdated2, "Key2 should not be outdated yet")

	// Wait for another 600ms (total 1.1s > 1s TTL)
	time.Sleep(600 * time.Millisecond)

	// Both keys should now be outdated
	outdated1 = c.Outdated(uopt.Of("key1"))
	outdated2 = c.Outdated(uopt.Of("key2"))
	assert.True(t, outdated1, "Key1 should be outdated after TTL")
	assert.True(t, outdated2, "Key2 should be outdated after TTL")
}

func TestComparableMapCache_Changes_AfterDrop(t *testing.T) {
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())

	c.Set("key1", 100)
	c.Set("key2", 200)
	c.Drop()

	changes := c.Changes()
	assert.Empty(t, changes, "Changes should be empty after Drop()")
}

func TestComparableMapCache_Changes_AfterDropKey(t *testing.T) {
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())

	c.Set("key1", 100)
	c.Set("key2", 200)
	c.Set("key3", 300)

	c.DropKey("key2")

	changes := c.Changes()
	expected := []string{"key1", "key3"}
	assert.ElementsMatch(t, expected, changes, "Changes should include only key2 after DropKey()")

	_, ok := c.Get("key2")
	assert.False(t, ok, "key2 should be removed from cache")
}

func TestComparableMapCache_OverwriteValue(t *testing.T) {
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())

	key := "key1"
	c.Set(key, 1)
	c.Set(key, 2)
	c.Set(key, 3)

	// Changes should include key1 only once
	changes := c.Changes()
	expected := []string{"key1"}
	assert.ElementsMatch(t, expected, changes, "Changes should include key1 once")

	// Verify that the value is the last set value
	val, ok := c.Get(key)
	require.True(t, ok, "Expected to retrieve value for key1")
	assert.Equal(t, 3, *val, "Expected value to be the last set value")
}

func TestComparableMapCache_SetQuietly_DoesNotTrackChanges(t *testing.T) {
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())

	c.SetQuietly("key1", 100)
	c.SetQuietly("key2", 200)

	// Changes should be empty
	changes := c.Changes()
	assert.Empty(t, changes, "SetQuietly should not track changes")

	// Verify that keys are set correctly
	val, ok := c.Get("key1")
	require.True(t, ok, "Expected to retrieve value for key1")
	assert.Equal(t, 100, *val, "Expected value for key1 to be 100")

	val, ok = c.Get("key2")
	require.True(t, ok, "Expected to retrieve value for key2")
	assert.Equal(t, 200, *val, "Expected value for key2 to be 200")
}

func TestComparableMapCache_DuplicateKeys(t *testing.T) {
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())

	key := "duplicateKey"
	c.Set(key, 1)
	c.Set(key, 2)
	c.Set(key, 3)

	// Changes should include duplicateKey only once
	changes := c.Changes()
	expected := []string{"duplicateKey"}
	assert.ElementsMatch(t, expected, changes, "Changes should include duplicateKey once")

	// Verify the final value
	val, ok := c.Get(key)
	require.True(t, ok, "Expected to retrieve value for duplicateKey")
	assert.Equal(t, 3, *val, "Expected value for duplicateKey to be 3")
}

func TestComparableMapCache_MultipleKeysSameValue(t *testing.T) {
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())

	keys := []string{"key1", "key2", "key3"}
	val := 42

	for _, key := range keys {
		c.Set(key, val)
	}

	// Verify that all keys have the same value
	for _, key := range keys {
		value, ok := c.Get(key)
		require.True(t, ok, "Expected to retrieve value for %s", key)
		assert.Equal(t, val, *value, "Expected value for %s to be %d", key, val)
	}

	// Verify changes
	changes := c.Changes()
	expected := keys
	assert.ElementsMatch(t, expected, changes, "Changes should include all keys")
}

func TestComparableMapCache_MultipleSetQuietly(t *testing.T) {
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())

	keys := []string{"key1", "key2", "key3"}
	vals := []int{10, 20, 30}

	for i, key := range keys {
		c.SetQuietly(key, vals[i])
	}

	// Changes should be empty
	changes := c.Changes()
	assert.Empty(t, changes, "SetQuietly should not track changes")

	// Verify that all keys have the correct values
	for i, key := range keys {
		value, ok := c.Get(key)
		require.True(t, ok, "Expected to retrieve value for %s", key)
		assert.Equal(t, vals[i], *value, "Expected value for %s to be %d", key, vals[i])
	}
}

func TestComparableMapCache_SetAndDropKeys(t *testing.T) {
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())

	keys := []string{"key1", "key2", "key3", "key4"}
	vals := []int{100, 200, 300, 400}

	for i, key := range keys {
		c.Set(key, vals[i])
	}

	c.DropKey("key2")
	c.DropKey("key4")

	_, ok := c.Get("key2")
	assert.False(t, ok, "key2 should be dropped")
	_, ok = c.Get("key4")
	assert.False(t, ok, "key4 should be dropped")

	val, ok := c.Get("key1")
	require.True(t, ok, "Expected to retrieve value for key1")
	assert.Equal(t, 100, *val, "Expected value for key1 to be 100")

	val, ok = c.Get("key3")
	require.True(t, ok, "Expected to retrieve value for key3")
	assert.Equal(t, 300, *val, "Expected value for key3 to be 300")

	changes := c.Changes()
	expected := []string{"key1", "key3"}
	assert.ElementsMatch(t, expected, changes, "Changes should include all set and dropped keys")
}

func TestComparableMapCache_OverwriteSetQuietly(t *testing.T) {
	c := ucache.NewInMemoryComparableMapCache[string, int](uopt.Null[time.Duration]())

	key := "key1"
	c.SetQuietly(key, 1)
	c.SetQuietly(key, 2)
	c.SetQuietly(key, 3)

	// Changes should be empty
	changes := c.Changes()
	assert.Empty(t, changes, "SetQuietly should not track changes")

	// Verify that the final value is correct
	val, ok := c.Get(key)
	require.True(t, ok, "Expected to retrieve value for key1")
	assert.Equal(t, 3, *val, "Expected value for key1 to be 3")
}
