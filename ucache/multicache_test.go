package ucache_test

import (
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kordax/basic-utils/ucache"
	"github.com/kordax/basic-utils/uopt"
	"github.com/stretchr/testify/assert"
)

type DummyComparable struct {
	Val int
}

func (d DummyComparable) Hash() int {
	return d.Val<<31 + d.Val
}

func (d DummyComparable) Equals(other ucache.Comparable) bool {
	switch o := other.(type) {
	case DummyComparable:
		return d.Val == o.Val
	default:
		return false
	}
}

//goland:noinspection GoUnusedExportedType
type SimpleKey int64

func (s SimpleKey) Equals(other ucache.Comparable) bool {
	return s == other
}

func (s SimpleKey) Key() int64 {
	return int64(s)
}

func (s SimpleKey) String() string {
	return strconv.Itoa(int(s))
}

type SimpleCompositeKey[T ucache.Unique] struct {
	keys []T
}

func (s SimpleCompositeKey[T]) Equals(other ucache.Comparable) bool {
	switch o := other.(type) {
	case SimpleCompositeKey[T]:
		for i, k := range s.keys {
			if len(s.keys) != len(o.keys) {
				return false
			}
			if k.Key() != o.keys[i].Key() || !k.Equals(o.keys[i]) {
				return false
			}
		}

		return true
	default:
		return false
	}
}

func (s SimpleCompositeKey[T]) Keys() []ucache.Unique {
	result := make([]ucache.Unique, len(s.keys))
	for i, key := range s.keys {
		result[i] = ucache.UIntKey(key.Key())
	}

	return result
}

func (s SimpleCompositeKey[T]) String() string {
	rep := make([]string, len(s.keys))
	for i, key := range s.keys {
		rep[i] = strconv.FormatInt(key.Key(), 10)
	}

	return strings.Join(rep, ", ")
}

func NewSimpleCompositeKey[T ucache.Unique](keys ...T) SimpleCompositeKey[T] {
	return SimpleCompositeKey[T]{keys: keys}
}

func TestHashMapMultiCache(t *testing.T) {
	c := ucache.NewDefaultHashMapMultiCache[SimpleCompositeKey[ucache.StringKey], DummyComparable](uopt.Null[time.Duration]())
	key := NewSimpleCompositeKey[ucache.StringKey]("atest")
	key2 := NewSimpleCompositeKey[ucache.StringKey]("bSeCond-keyQ@!%!%#")
	val := 326
	c.Put(key, DummyComparable{Val: val})

	changes := c.Changes()
	assert.NotNil(t, changes)
	cached := c.Get(key)
	assert.Contains(t, cached, DummyComparable{Val: val})

	assert.Empty(t, c.Changes())
	for i := 0; i < 10; i++ {
		c.Put(key, DummyComparable{Val: i})
	}
	c.Put(key2, DummyComparable{Val: 65535})
	changes = c.Changes()
	assert.EqualValues(t, []SimpleCompositeKey[ucache.StringKey]{key, key2}, changes)

	result := c.Get(key2)
	assert.Contains(t, result, DummyComparable{Val: 65535})

	complexKeyBase := []ucache.StringKey{"p1", "p2", "p3"}
	partialComplexKey := NewSimpleCompositeKey[ucache.StringKey](complexKeyBase...)

	for i := 0; i < 10; i++ {
		complexKey := NewSimpleCompositeKey[ucache.StringKey](append(complexKeyBase, ucache.StringKey("number:"+strconv.Itoa(i)))...)
		c.Put(complexKey, DummyComparable{Val: i})
		changes = c.Changes()
		assert.Contains(t, changes, complexKey)
	}

	result = c.Get(partialComplexKey)
	assert.NotEmpty(t, result)
	for i := 0; i < 10; i++ {
		assert.Contains(t, result, DummyComparable{i})
	}
}

func TestHashMapMultiCache_CompositeKey(t *testing.T) {
	c := ucache.NewDefaultHashMapMultiCache[ucache.StrCompositeKey, DummyComparable](uopt.Null[time.Duration]())
	key := ucache.NewStrCompositeKey("category", "kp_2")
	key2 := ucache.NewStrCompositeKey("category2", "kp_2")
	val := DummyComparable{Val: 10}
	val2 := DummyComparable{Val: 236261}

	c.PutQuietly(key, val)
	c.PutQuietly(key2, val2)

	results := c.Get(key)
	results2 := c.Get(key2)
	assert.Len(t, results, 1)
	assert.Len(t, results, 1)
	assert.EqualValues(t, results[0], val)
	assert.EqualValues(t, results2[0], val2)
}

func TestHashMapMultiCache_DropKey(t *testing.T) {
	c := ucache.NewDefaultHashMapMultiCache[ucache.StrCompositeKey, DummyComparable](uopt.Null[time.Duration]())
	categoryKey := ucache.NewStrCompositeKey("category")
	key := ucache.NewStrCompositeKey("category", "kp_232626")
	key2 := ucache.NewStrCompositeKey("category2", "kp_232626")
	catVal := DummyComparable{Val: rand.Int()}
	val := DummyComparable{Val: rand.Int()}
	val2 := DummyComparable{Val: rand.Int()}

	c.Put(categoryKey, catVal)
	c.Put(key, val)
	c.Put(key2, val2)

	catRes := c.Get(categoryKey)
	res := c.Get(key)
	res2 := c.Get(key2)
	assert.Len(t, catRes, 2)
	assert.Len(t, res, 1)
	assert.Len(t, res2, 1)

	c.DropKey(key)
	catRes = c.Get(categoryKey)
	res = c.Get(key)
	res2 = c.Get(key2)
	assert.Len(t, catRes, 2)
	assert.Len(t, res, 0)
	assert.Len(t, res2, 1)
}

func TestHashMapMultiCache_PutQuietly(t *testing.T) {
	c := ucache.NewDefaultHashMapMultiCache[SimpleCompositeKey[ucache.StringKey], DummyComparable](uopt.Null[time.Duration]())
	key := NewSimpleCompositeKey[ucache.StringKey]("kp_1", "kp_2")
	val := DummyComparable{Val: 10}
	val2 := DummyComparable{Val: 15}

	c.PutQuietly(key, val)
	c.PutQuietly(key, val)
	c.PutQuietly(key, val)

	results := c.Get(key)
	assert.Len(t, results, 1)

	c.PutQuietly(key, val2)
	results = c.Get(key)
	assert.Len(t, results, 2)

	c.PutQuietly(key, val)
	results = c.Get(key)
	assert.Len(t, results, 2)
}

func TestTreeMultiCache(t *testing.T) {
	c := ucache.NewInMemoryTreeMultiCache[SimpleCompositeKey[ucache.StringKey], DummyComparable](uopt.Null[time.Duration]())
	key := NewSimpleCompositeKey[ucache.StringKey]("atest")
	key2 := NewSimpleCompositeKey[ucache.StringKey]("bSeCond-keyQ@!%!%#")
	val := 326
	c.Put(key, DummyComparable{Val: val})

	changes := c.Changes()
	assert.NotNil(t, changes)
	cached := c.Get(key)
	assert.Contains(t, cached, DummyComparable{Val: val})

	assert.Empty(t, c.Changes())
	for i := 0; i < 10; i++ {
		c.Put(key, DummyComparable{Val: i})
	}
	c.Put(key2, DummyComparable{Val: 65535})
	changes = c.Changes()
	assert.EqualValues(t, []SimpleCompositeKey[ucache.StringKey]{key, key2}, changes)

	result := c.Get(key2)
	assert.Contains(t, result, DummyComparable{Val: 65535})

	complexKeyBase := []ucache.StringKey{"p1", "p2", "p3"}
	partialComplexKey := NewSimpleCompositeKey[ucache.StringKey](complexKeyBase...)

	for i := 0; i < 10; i++ {
		complexKey := NewSimpleCompositeKey[ucache.StringKey](append(complexKeyBase, ucache.StringKey("number:"+strconv.Itoa(i)))...)
		c.Put(complexKey, DummyComparable{Val: i})
		changes = c.Changes()
		assert.Contains(t, changes, complexKey)
	}

	result = c.Get(partialComplexKey)
	assert.NotEmpty(t, result)
	for i := 0; i < 10; i++ {
		assert.Contains(t, result, DummyComparable{i})
	}
}

func TestTreeMultiCache_CompositeKey(t *testing.T) {
	c := ucache.NewInMemoryTreeMultiCache[ucache.StrCompositeKey, DummyComparable](uopt.Null[time.Duration]())
	key := ucache.NewStrCompositeKey("category", "kp_2")
	key2 := ucache.NewStrCompositeKey("category2", "kp_2")
	val := DummyComparable{Val: 10}
	val2 := DummyComparable{Val: 236261}

	c.PutQuietly(key, val)
	c.PutQuietly(key2, val2)

	results := c.Get(key)
	results2 := c.Get(key2)
	assert.Len(t, results, 1)
	assert.Len(t, results, 1)
	assert.EqualValues(t, results[0], val)
	assert.EqualValues(t, results2[0], val2)
}

func TestTreeMultiCache_DropKey(t *testing.T) {
	c := ucache.NewInMemoryTreeMultiCache[ucache.StrCompositeKey, DummyComparable](uopt.Null[time.Duration]())
	categoryKey := ucache.NewStrCompositeKey("category")
	key := ucache.NewStrCompositeKey("category", "kp_232626")
	key2 := ucache.NewStrCompositeKey("category2", "kp_232626")
	catVal := DummyComparable{Val: rand.Int()}
	val := DummyComparable{Val: rand.Int()}
	val2 := DummyComparable{Val: rand.Int()}

	c.Put(categoryKey, catVal)
	c.Put(key, val)
	c.Put(key2, val2)

	catRes := c.Get(categoryKey)
	res := c.Get(key)
	res2 := c.Get(key2)
	assert.Len(t, catRes, 2)
	assert.Len(t, res, 1)
	assert.Len(t, res2, 1)

	c.DropKey(key)
	catRes = c.Get(categoryKey)
	res = c.Get(key)
	res2 = c.Get(key2)
	assert.Len(t, catRes, 1)
	assert.Len(t, res, 0)
	assert.Len(t, res2, 1)
}

func TestTreeMultiCache_AddTransparent(t *testing.T) {
	c := ucache.NewInMemoryTreeMultiCache[SimpleCompositeKey[ucache.StringKey], DummyComparable](uopt.Null[time.Duration]())
	key := NewSimpleCompositeKey[ucache.StringKey]("kp_1", "kp_2")
	val := DummyComparable{Val: 10}
	val2 := DummyComparable{Val: 15}

	c.PutQuietly(key, val)
	c.PutQuietly(key, val)
	c.PutQuietly(key, val)

	results := c.Get(key)
	assert.Len(t, results, 1)

	c.PutQuietly(key, val2)
	results = c.Get(key)
	assert.Len(t, results, 2)

	c.PutQuietly(key, val)
	results = c.Get(key)
	assert.Len(t, results, 2)
}

func TestNewGenericCompositeKey(t *testing.T) {
	c := ucache.NewInMemoryTreeMultiCache[ucache.GenericCompositeKey, DummyComparable](uopt.Null[time.Duration]())
	key1 := ucache.NewGenericCompositeKey("KeyString", 2, 3.0, uint8(4), int16(-5))
	key2 := ucache.NewGenericCompositeKey("KeyString", 2, 3.0, int8(100), uint16(5))
	subKey1 := ucache.NewGenericCompositeKey("KeyString", 2, 3.0, uint8(4))
	subKey2 := ucache.NewGenericCompositeKey("KeyString", 2, 3.0, int8(100))
	baseKey := ucache.NewGenericCompositeKey("KeyString", 2, 3.0)

	var firstKeyDummies []DummyComparable
	for i := 0; i < 2; i++ {
		firstKeyDummies = append(firstKeyDummies, DummyComparable{Val: i})
	}
	var secondKeyDummies []DummyComparable
	for i := 2; i < 7; i++ {
		secondKeyDummies = append(secondKeyDummies, DummyComparable{Val: i})
	}

	sort.Slice(firstKeyDummies, func(i, j int) bool {
		return firstKeyDummies[i].Val < (firstKeyDummies[j].Val)
	})
	sort.Slice(secondKeyDummies, func(i, j int) bool {
		return secondKeyDummies[i].Val < (secondKeyDummies[j].Val)
	})

	c.Put(key1, firstKeyDummies...)
	c.Put(key2, secondKeyDummies...)

	changes := c.Changes()
	assert.NotNil(t, changes)
	key1Result := c.Get(key1)
	assert.EqualValues(t, firstKeyDummies, key1Result)
	key2Result := c.Get(key2)
	assert.EqualValues(t, secondKeyDummies, key2Result)
	subKey1Result := c.Get(subKey1)
	assert.EqualValues(t, firstKeyDummies, subKey1Result)
	subKey2Result := c.Get(subKey2)
	assert.EqualValues(t, secondKeyDummies, subKey2Result)
	baseKeyResult := c.Get(baseKey)

	combinedDummies := append(firstKeyDummies, secondKeyDummies...)
	sort.Slice(combinedDummies, func(i, j int) bool {
		return combinedDummies[i].Val < (combinedDummies[j].Val)
	})
	sort.Slice(baseKeyResult, func(i, j int) bool {
		return baseKeyResult[i].Val < (baseKeyResult[j].Val)
	})

	assert.EqualValues(t, combinedDummies, baseKeyResult)
}

type CollisionTestKey struct {
	id   int
	hash []int64
}

// Implement the CompositeKey interface for TestKey
func (k CollisionTestKey) Keys() []ucache.Unique {
	result := make([]ucache.Unique, len(k.hash))
	for i, h := range k.hash {
		result[i] = ucache.IntKey(h)
	}
	return result
}

func (k CollisionTestKey) String() string {
	return strconv.Itoa(k.id)
}

func (k CollisionTestKey) Equals(other ucache.Comparable) bool {
	ok, _ := other.(CollisionTestKey)
	return k.id == ok.id
}

func TestTreeMultiCacheHighCollisionProbability(t *testing.T) {
	c := ucache.NewInMemoryTreeMultiCache[CollisionTestKey, ucache.Int64Value](uopt.Null[time.Duration]())

	// Define a set of keys that all produce the same hash code
	keys := []CollisionTestKey{
		{id: 1, hash: []int64{1, 2, 3}},
		{id: 2, hash: []int64{1, 2, 3}},
		{id: 3, hash: []int64{1, 2, 3}},
	}

	// Add values to the c for each key
	for i, key := range keys {
		c.Put(key, ucache.NewInt64Value(int64(i)))
	}

	// Ensure that all values can be retrieved despite the high collision probability
	for i, key := range keys {
		values := c.Get(key)
		assert.Contains(t, values, ucache.NewInt64Value(int64(i))) // Check if the expected value is present in the retrieved values
	}
}

func TestInMemoryTreeMultiCache_Set(t *testing.T) {
	c := ucache.NewInMemoryTreeMultiCache[ucache.StrCompositeKey, ucache.StringValue](uopt.Null[time.Duration]())
	key := ucache.NewStrCompositeKey("key1")
	value1 := ucache.NewStringValue("value1")
	value2 := ucache.NewStringValue("value2")

	c.Set(key, value1)
	retrieved := c.Get(key)
	assert.Len(t, retrieved, 1)
	assert.Equal(t, value1, retrieved[0])

	c.Set(key, value2)
	retrieved = c.Get(key)
	assert.Len(t, retrieved, 1)
	assert.Equal(t, value2, retrieved[0])
}

func TestInMemoryTreeMultiCache_Outdated_WithStringKeyAndValue(t *testing.T) {
	ttl := 1 * time.Millisecond
	longTTL := 1 * time.Hour
	c := ucache.NewInMemoryTreeMultiCache[ucache.StrCompositeKey, ucache.StringValue](uopt.Of(ttl))
	cLong := ucache.NewInMemoryTreeMultiCache[ucache.StrCompositeKey, ucache.StringValue](uopt.Of(longTTL))

	key1 := ucache.NewStrCompositeKey("key1")
	key2 := ucache.NewStrCompositeKey("key2")
	value1 := ucache.NewStringValue("value1")

	// Test with long TTL
	assert.True(t, cLong.Outdated(uopt.Null[ucache.StrCompositeKey]()))
	cLong.Put(key1, value1)
	assert.False(t, cLong.Outdated(uopt.Of(key1)))
	time.Sleep(10 * time.Millisecond)
	assert.False(t, cLong.Outdated(uopt.Of(key1)))

	// Test immediate expiration
	assert.True(t, c.Outdated(uopt.Null[ucache.StrCompositeKey]()))
	c.Put(key1, value1)
	assert.False(t, c.Outdated(uopt.Of(key1)))
	time.Sleep(ttl + 10*time.Millisecond)
	assert.True(t, c.Outdated(uopt.Of(key1)))

	// Test overwriting key resets TTL
	c.Put(key1, value1)
	time.Sleep(ttl / 2)
	c.Put(key1, value1) // Reset TTL
	assert.False(t, c.Outdated(uopt.Of(key1)))
	assert.False(t, c.Outdated(uopt.Of(key1)))
	time.Sleep(ttl)
	assert.True(t, c.Outdated(uopt.Of(key1)))

	// Test Drop() method
	c.Put(key1, value1)
	c.Put(key2, value1)
	c.Drop()
	assert.True(t, c.Outdated(uopt.Of(key1)))
	assert.True(t, c.Outdated(uopt.Of(key2)))
}

func TestInMemoryTreeMultiCache_Outdated_WithDifferentTTLs(t *testing.T) {
	shortTTL := 10 * time.Millisecond
	mediumTTL := 20 * time.Millisecond
	longTTL := 30 * time.Millisecond

	cShort := ucache.NewInMemoryTreeMultiCache[ucache.StrCompositeKey, ucache.StringValue](uopt.Of(shortTTL))
	cMedium := ucache.NewInMemoryTreeMultiCache[ucache.StrCompositeKey, ucache.StringValue](uopt.Of(mediumTTL))
	cLong := ucache.NewInMemoryTreeMultiCache[ucache.StrCompositeKey, ucache.StringValue](uopt.Of(longTTL))

	key := ucache.NewStrCompositeKey("key")

	value := ucache.NewStringValue("value")

	cShort.Put(key, value)
	cMedium.Put(key, value)
	cLong.Put(key, value)

	time.Sleep(shortTTL + 1*time.Millisecond)
	assert.True(t, cShort.Outdated(uopt.Of(key)))
	assert.False(t, cMedium.Outdated(uopt.Of(key)))
	assert.False(t, cLong.Outdated(uopt.Of(key)))

	time.Sleep(mediumTTL - shortTTL)
	assert.True(t, cMedium.Outdated(uopt.Of(key)))
	assert.False(t, cLong.Outdated(uopt.Of(key)))

	time.Sleep(longTTL - mediumTTL)
	assert.True(t, cLong.Outdated(uopt.Of(key)))
}
