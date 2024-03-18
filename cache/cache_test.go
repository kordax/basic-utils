package cache_test

import (
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kordax/basic-utils/cache"
	"github.com/kordax/basic-utils/opt"
	"github.com/stretchr/testify/assert"
)

type DummyComparable struct {
	Val int
}

func (d DummyComparable) Hash() int {
	return d.Val<<31 + d.Val
}

func (d DummyComparable) Equals(other cache.Comparable) bool {
	switch o := other.(type) {
	case DummyComparable:
		return d.Val == o.Val
	default:
		return false
	}
}

//goland:noinspection GoUnusedExportedType
type SimpleKey int

func (s SimpleKey) Equals(other cache.Comparable) bool {
	return s == other
}

func (s SimpleKey) Key() int {
	return int(s)
}

func (s SimpleKey) String() string {
	return strconv.Itoa(int(s))
}

type SimpleCompositeKey[T cache.Hashed] struct {
	keys []T
}

func (s SimpleCompositeKey[T]) Equals(other cache.Comparable) bool {
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

func (s SimpleCompositeKey[T]) Keys() []int {
	result := make([]int, len(s.keys))
	for i, key := range s.keys {
		result[i] = key.Key()
	}

	return result
}

func (s SimpleCompositeKey[T]) String() string {
	rep := make([]string, len(s.keys))
	for i, key := range s.keys {
		rep[i] = key.String()
	}

	return strings.Join(rep, ", ")
}

func NewSimpleCompositeKey[T cache.Hashed](keys ...T) SimpleCompositeKey[T] {
	return SimpleCompositeKey[T]{keys: keys}
}

func TestHashMapCache(t *testing.T) {
	c := cache.NewDefaultHashMapCache[SimpleCompositeKey[cache.StringKey], DummyComparable](opt.Null[time.Duration]())
	key := NewSimpleCompositeKey[cache.StringKey]("atest")
	key2 := NewSimpleCompositeKey[cache.StringKey]("bSeCond-keyQ@!%!%#")
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
	assert.EqualValues(t, []SimpleCompositeKey[cache.StringKey]{key, key2}, changes)

	result := c.Get(key2)
	assert.Contains(t, result, DummyComparable{Val: 65535})

	complexKeyBase := []cache.StringKey{"p1", "p2", "p3"}
	partialComplexKey := NewSimpleCompositeKey[cache.StringKey](complexKeyBase...)

	for i := 0; i < 10; i++ {
		complexKey := NewSimpleCompositeKey[cache.StringKey](append(complexKeyBase, cache.StringKey("number:"+strconv.Itoa(i)))...)
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

func TestHashMapCache_CompositeKey(t *testing.T) {
	c := cache.NewDefaultHashMapCache[cache.StrCompositeKey, DummyComparable](opt.Null[time.Duration]())
	key := cache.NewStrCompositeKey("category", "kp_2")
	key2 := cache.NewStrCompositeKey("category2", "kp_2")
	val := DummyComparable{Val: 10}
	val2 := DummyComparable{Val: 236261}

	c.AddSilently(key, val)
	c.AddSilently(key2, val2)

	results := c.Get(key)
	results2 := c.Get(key2)
	assert.Len(t, results, 1)
	assert.Len(t, results, 1)
	assert.EqualValues(t, results[0], val)
	assert.EqualValues(t, results2[0], val2)
}

func TestHashMapCache_DropKey(t *testing.T) {
	c := cache.NewDefaultHashMapCache[cache.StrCompositeKey, DummyComparable](opt.Null[time.Duration]())
	categoryKey := cache.NewStrCompositeKey("category")
	key := cache.NewStrCompositeKey("category", "kp_232626")
	key2 := cache.NewStrCompositeKey("category2", "kp_232626")
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

func TestHashMapCache_AddTransparent(t *testing.T) {
	c := cache.NewDefaultHashMapCache[SimpleCompositeKey[cache.StringKey], DummyComparable](opt.Null[time.Duration]())
	key := NewSimpleCompositeKey[cache.StringKey]("kp_1", "kp_2")
	val := DummyComparable{Val: 10}
	val2 := DummyComparable{Val: 15}

	c.AddSilently(key, val)
	c.AddSilently(key, val)
	c.AddSilently(key, val)

	results := c.Get(key)
	assert.Len(t, results, 1)

	c.AddSilently(key, val2)
	results = c.Get(key)
	assert.Len(t, results, 2)

	c.AddSilently(key, val)
	results = c.Get(key)
	assert.Len(t, results, 2)
}

func TestTreeCache(t *testing.T) {
	c := cache.NewInMemoryTreeCache[SimpleCompositeKey[cache.StringKey], DummyComparable](opt.Null[time.Duration]())
	key := NewSimpleCompositeKey[cache.StringKey]("atest")
	key2 := NewSimpleCompositeKey[cache.StringKey]("bSeCond-keyQ@!%!%#")
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
	assert.EqualValues(t, []SimpleCompositeKey[cache.StringKey]{key, key2}, changes)

	result := c.Get(key2)
	assert.Contains(t, result, DummyComparable{Val: 65535})

	complexKeyBase := []cache.StringKey{"p1", "p2", "p3"}
	partialComplexKey := NewSimpleCompositeKey[cache.StringKey](complexKeyBase...)

	for i := 0; i < 10; i++ {
		complexKey := NewSimpleCompositeKey[cache.StringKey](append(complexKeyBase, cache.StringKey("number:"+strconv.Itoa(i)))...)
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

func TestTreeCache_CompositeKey(t *testing.T) {
	c := cache.NewInMemoryTreeCache[cache.StrCompositeKey, DummyComparable](opt.Null[time.Duration]())
	key := cache.NewStrCompositeKey("category", "kp_2")
	key2 := cache.NewStrCompositeKey("category2", "kp_2")
	val := DummyComparable{Val: 10}
	val2 := DummyComparable{Val: 236261}

	c.AddSilently(key, val)
	c.AddSilently(key2, val2)

	results := c.Get(key)
	results2 := c.Get(key2)
	assert.Len(t, results, 1)
	assert.Len(t, results, 1)
	assert.EqualValues(t, results[0], val)
	assert.EqualValues(t, results2[0], val2)
}

func TestTreeCache_DropKey(t *testing.T) {
	c := cache.NewInMemoryTreeCache[cache.StrCompositeKey, DummyComparable](opt.Null[time.Duration]())
	categoryKey := cache.NewStrCompositeKey("category")
	key := cache.NewStrCompositeKey("category", "kp_232626")
	key2 := cache.NewStrCompositeKey("category2", "kp_232626")
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

func TestTreeCache_AddTransparent(t *testing.T) {
	c := cache.NewInMemoryTreeCache[SimpleCompositeKey[cache.StringKey], DummyComparable](opt.Null[time.Duration]())
	key := NewSimpleCompositeKey[cache.StringKey]("kp_1", "kp_2")
	val := DummyComparable{Val: 10}
	val2 := DummyComparable{Val: 15}

	c.AddSilently(key, val)
	c.AddSilently(key, val)
	c.AddSilently(key, val)

	results := c.Get(key)
	assert.Len(t, results, 1)

	c.AddSilently(key, val2)
	results = c.Get(key)
	assert.Len(t, results, 2)

	c.AddSilently(key, val)
	results = c.Get(key)
	assert.Len(t, results, 2)
}

func TestComparableSlice_Equals(t *testing.T) {
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

func TestNewGenericCompositeKey(t *testing.T) {
	c := cache.NewInMemoryTreeCache[cache.GenericCompositeKey, DummyComparable](opt.Null[time.Duration]())
	key1 := cache.NewGenericCompositeKey("KeyString", 2, 3.0, uint8(4), int16(-5))
	key2 := cache.NewGenericCompositeKey("KeyString", 2, 3.0, int8(100), uint16(5))
	subKey1 := cache.NewGenericCompositeKey("KeyString", 2, 3.0, uint8(4))
	subKey2 := cache.NewGenericCompositeKey("KeyString", 2, 3.0, int8(100))
	baseKey := cache.NewGenericCompositeKey("KeyString", 2, 3.0)

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
	hash []int
}

// Implement the CompositeKey interface for TestKey
func (k CollisionTestKey) Keys() []int {
	return k.hash
}

func (k CollisionTestKey) String() string {
	return strconv.Itoa(k.id)
}

func (k CollisionTestKey) Equals(other cache.Comparable) bool {
	ok, _ := other.(CollisionTestKey)
	return k.id == ok.id
}

func TestTreeCacheHighCollisionProbability(t *testing.T) {
	c := cache.NewInMemoryTreeCache[CollisionTestKey, cache.Int64Value](opt.Null[time.Duration]())

	// Define a set of keys that all produce the same hash code
	keys := []CollisionTestKey{
		{id: 1, hash: []int{1, 2, 3}},
		{id: 2, hash: []int{1, 2, 3}},
		{id: 3, hash: []int{1, 2, 3}},
	}

	// Add values to the c for each key
	for i, key := range keys {
		c.Put(key, cache.NewInt64Value(int64(i)))
	}

	// Ensure that all values can be retrieved despite the high collision probability
	for i, key := range keys {
		values := c.Get(key)
		assert.Contains(t, values, cache.NewInt64Value(int64(i))) // Check if the expected value is present in the retrieved values
	}
}
