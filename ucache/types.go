package ucache

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dgryski/go-farm"
	"github.com/kordax/basic-utils/uarray"
	"github.com/vmihailenco/msgpack/v5"
)

type UIntKey uint64
type IntKey int64
type StringKey string

type keyContainer[K any] struct {
	key       K
	updatedAt time.Time
}

/*
Comparable entity
*/
type Comparable interface {
	Equals(other Comparable) bool
}

/*
Unique specifies an abstract key with an ability to provide hash.
*/
type Unique interface {
	Comparable
	Key() int64 // Key should return a unique item key. It can be a hash or just an index.
}

/*
CompositeKey specifies an abstract key with an ability to provide an ordered list of available keys.
*/
type CompositeKey interface {
	Comparable
	Keys() []Unique // Keys returns an ordered list of keys ordered by priority (ASC), so the first element has the most prio.
}

func (k IntKey) Equals(other Comparable) bool {
	otherKey, ok := other.(IntKey)
	return ok && k == otherKey
}

func (k IntKey) Key() int64 {
	return int64(k)
}

func (k IntKey) Keys() []Unique {
	return []Unique{k}
}

func (k IntKey) String() string {
	return fmt.Sprintf("%d", k)
}

func (k StringKey) Equals(other Comparable) bool {
	otherKey, ok := other.(StringKey)
	return ok && k == otherKey
}

func (k StringKey) Key() int64 {
	hash := int64(0)
	for i := 0; i < len(k); i++ {
		hash = 31*hash + int64(k[i])
	}

	return hash
}

func (k StringKey) Keys() []Unique {
	return []Unique{IntKey(farm.Hash64([]byte(k)))}
}

func (k StringKey) String() string {
	return string(k)
}

func (k UIntKey) Equals(other Comparable) bool {
	otherKey, ok := other.(UIntKey)
	return ok && k == otherKey
}

func (k UIntKey) Key() int64 {
	return int64(k)
}

func (k UIntKey) Keys() []Unique {
	return []Unique{k}
}

func (k UIntKey) String() string {
	return fmt.Sprintf("%d", k)
}

type ComparableKey[T comparable] struct {
	v T
}

func NewComparableKey[T comparable](v T) ComparableKey[T] {
	return ComparableKey[T]{v: v}
}

func (k ComparableKey[T]) Key() int64 {
	hash := int64(0)

	switch value := any(k.v).(type) {
	case int, int8, int16, int32, int64:
		hash = 31*hash + reflect.ValueOf(value).Int()
	case uint, uint8, uint16, uint32, uint64:
		hash = 31*hash + int64(reflect.ValueOf(value).Uint())
	case float32, float64:
		hash = 31*hash + int64(reflect.ValueOf(value).Float())
	case complex64, complex128:
		c := reflect.ValueOf(value).Complex()
		hash = 31*(31*hash+int64(real(c))) + int64(imag(c))
	case string:
		for i := 0; i < len(value); i++ {
			hash = 31*hash + int64(value[i])
		}
	case bool:
		if value {
			hash = 31*hash + 1
		} else {
			hash = 31*hash + 0
		}
	default:
		return hash
	}

	return hash
}

func (k ComparableKey[T]) String() string {
	return convertToString(k.v)
}

func (k ComparableKey[T]) Equals(other Comparable) bool {
	otherKey, ok := other.(ComparableKey[T])
	return ok && k == otherKey
}

type ComparableSlice[T Comparable] struct {
	Data []T
}

func (c ComparableSlice[T]) Equals(other Comparable) bool {
	switch o := other.(type) {
	case ComparableSlice[T]:
		return c.eq(o)
	default:
		return false
	}
}

func (c ComparableSlice[T]) eq(other ComparableSlice[T]) bool {
	if len(c.Data) != len(other.Data) {
		return false
	}

	for i, v := range c.Data {
		if !v.Equals(other.Data[i]) {
			return false
		}
	}

	return true
}

type UIntCompositeKey struct {
	keys []UIntKey
}

func NewUIntCompositeKey(keys ...uint64) UIntCompositeKey {
	var conv []UIntKey
	for _, key := range keys {
		conv = append(conv, UIntKey(key))
	}
	return UIntCompositeKey{keys: conv}
}

func (k UIntCompositeKey) Equals(other Comparable) bool {
	switch o := other.(type) {
	case UIntCompositeKey:
		return uarray.EqualsWithOrder(k.keys, o.keys)
	default:
		return false
	}
}

func (k UIntCompositeKey) Keys() []int64 {
	result := make([]int64, len(k.keys))
	for i, key := range k.keys {
		conv := IntKey(key)
		result[i] = conv.Key()
	}

	return result
}

func (k UIntCompositeKey) String() string {
	rep := make([]string, len(k.keys))
	for i, key := range k.keys {
		rep[i] = strconv.FormatUint(uint64(key), 10)
	}

	return strings.Join(rep, ", ")
}

type IntCompositeKey struct {
	keys []IntKey
}

func NewIntCompositeKey(keys ...int64) IntCompositeKey {
	var conv []IntKey
	for _, key := range keys {
		conv = append(conv, IntKey(key))
	}
	return IntCompositeKey{keys: conv}
}

func (k IntCompositeKey) Equals(other Comparable) bool {
	switch o := other.(type) {
	case IntCompositeKey:
		return uarray.EqualsWithOrder(k.keys, o.keys)
	default:
		return false
	}
}

func (k IntCompositeKey) Keys() []Unique {
	result := make([]Unique, len(k.keys))
	for i, key := range k.keys {
		result[i] = IntKey(key.Key())
	}

	return result
}

func (k IntCompositeKey) String() string {
	rep := make([]string, len(k.keys))
	for i, key := range k.keys {
		rep[i] = strconv.FormatInt(int64(key), 10)
	}

	return strings.Join(rep, ", ")
}

type StrCompositeKey struct {
	keys []StringKey
}

//goland:noinspection GoUnusedExportedFunction
func NewStrCompositeKey(keys ...string) StrCompositeKey {
	var conv []StringKey
	for _, key := range keys {
		conv = append(conv, StringKey(key))
	}
	return StrCompositeKey{keys: conv}
}

func (k StrCompositeKey) Equals(other Comparable) bool {
	switch o := other.(type) {
	case StrCompositeKey:
		return uarray.EqualsWithOrder(k.keys, o.keys)
	default:
		return false
	}
}

func (k StrCompositeKey) Keys() []Unique {
	result := make([]Unique, len(k.keys))
	for i, key := range k.keys {
		result[i] = IntKey(key.Key())
	}

	return result
}

func (k StrCompositeKey) String() string {
	return strings.Join(uarray.Map(k.keys, func(v *StringKey) string {
		return string(*v)
	}), ", ")
}

type GenericCompositeKey struct {
	keys []ComparableKey[any]
}

// NewGenericCompositeKey creates a GenericCompositeKey that supports only 'comparable' keys
func NewGenericCompositeKey(keys ...any) GenericCompositeKey {
	var conv []ComparableKey[any]
	for _, key := range keys {
		if !isComparable(key) {
			panic(fmt.Errorf("unsupported key type passed to NewGenericCompositeKey: %s", reflect.TypeOf(key).Kind().String()))
		}
		conv = append(conv, NewComparableKey(key))
	}
	return GenericCompositeKey{keys: any(conv).([]ComparableKey[any])}
}

func (k GenericCompositeKey) Equals(other Comparable) bool {
	switch o := other.(type) {
	case GenericCompositeKey:
		return uarray.EqualsWithOrder(k.keys, o.keys)
	default:
		return false
	}
}

func (k GenericCompositeKey) Keys() []Unique {
	result := make([]Unique, len(k.keys))
	for i, key := range k.keys {
		result[i] = IntKey(key.Key())
	}

	return result
}

func (k GenericCompositeKey) String() string {
	builder := strings.Builder{}
	for _, key := range k.keys {
		builder.WriteString(convertToString(key))
	}

	return builder.String()
}

type StringValue struct {
	v string
}

func NewStringValue(v string) StringValue {
	return StringValue{v: v}
}

func (s StringValue) Value() string {
	return s.v
}

func (s StringValue) Equals(other Comparable) bool {
	otherValuePtr, pok := other.(*StringValue)
	if !pok {
		otherValue, ok := other.(StringValue)
		if !ok {
			return false
		}

		return s.v == otherValue.v
	}

	return s.v == otherValuePtr.v
}

type StringSliceValue struct {
	v []string
}

func NewStringSliceValue(v []string) StringSliceValue {
	sort.Strings(v)
	return StringSliceValue{v: v}
}

func (s StringSliceValue) Values() []string {
	return s.v
}

func (s StringSliceValue) Equals(other Comparable) bool {
	otherValuePtr, pok := other.(*StringSliceValue)
	if !pok {
		otherValue, ok := other.(StringSliceValue)
		if !ok {
			return false
		}

		return uarray.EqualValues(s.v, otherValue.v)
	}

	return uarray.EqualValues(s.v, otherValuePtr.v)
}

type Int64Value struct {
	v int64
}

func NewInt64Value(v int64) Int64Value {
	return Int64Value{v: v}
}

func (s Int64Value) Value() string {
	return strconv.FormatInt(s.v, 10)
}

func (s Int64Value) Equals(other Comparable) bool {
	otherValuePtr, pok := other.(*Int64Value)
	if !pok {
		otherValue, ok := other.(Int64Value)
		if !ok {
			return false
		}

		return s.v == otherValue.v
	}

	return s.v == otherValuePtr.v
}

/*
FarmHash64Entity wraps any object and provides a Unique implementation
using farm's 64-bit hash function to be used in cache.
This hashed entity uses internal hash to avoid redundant rehashing operations.

IMPORTANT: The object must have exported fields and only those fields will be considered for the hashing uniqueness operation.
IMPORTANT: If the object is a pointer, the hash will compare pointer values. If the object is not a pointer, the hash will compare contents.

  - Equals method compares the hash values of the wrapped objects.
  - Key method uses farm.Hash64 to generate a 64-bit hash of the object and
    returns it as an int64 value for the key.

This can be used for uniquely identifying objects and comparing them
based on their content rather than their memory address.
*/
type FarmHash64Entity struct {
	obj       any
	hashValue int64
	hashReady bool
}

/*
Hashed is a constructor function that creates and returns a new instance
of FarmHash64Entity, wrapping the provided object. This instance provides
a Unique implementation using farm's 64-bit hash function.

Usage:
  - To uniquely identify objects based on their content rather than their
    memory address.
  - To generate a unique key for any object by hashing its content.

Example:

	obj := Hashed("example object")
	fmt.Println("Key:", obj.Key())
	// Outputs the unique key generated by farm.Hash64 for the string "example object".
*/
func Hashed(obj any) *FarmHash64Entity {
	return &FarmHash64Entity{obj: obj}
}

func (e *FarmHash64Entity) calculateHash() int64 {
	if e.hashReady {
		return e.hashValue
	}

	val := reflect.ValueOf(e.obj)
	if val.Kind() == reflect.Ptr {
		e.hashValue = int64(val.Pointer())
	} else {
		b, err := msgpack.Marshal(e.obj)
		if err != nil {
			panic(err)
		}
		e.hashValue = int64(farm.Hash64(b))
	}

	e.hashReady = true

	return e.hashValue
}

func (e *FarmHash64Entity) Equals(other Comparable) bool {
	otherFH, ok := other.(*FarmHash64Entity)
	if !ok {
		return false
	}
	return e.calculateHash() == otherFH.calculateHash()
}

func (e *FarmHash64Entity) Key() int64 {
	return e.calculateHash()
}

type FarmHash64CompositeKey struct {
	keys []ComparableKey[FarmHash64Entity]
}

// NewFarmHashCompositeKey creates a GenericCompositeKey that supports only 'comparable' keys
func NewFarmHashCompositeKey(keys ...any) FarmHash64CompositeKey {
	var conv []ComparableKey[any]
	for _, key := range keys {
		if !isComparable(key) {
			panic(fmt.Errorf("unsupported key type passed to NewGenericCompositeKey: %s", reflect.TypeOf(key).Kind().String()))
		}
		conv = append(conv, NewComparableKey(key))
	}
	return FarmHash64CompositeKey{keys: any(conv).([]ComparableKey[FarmHash64Entity])}
}

func (k FarmHash64CompositeKey) Equals(other Comparable) bool {
	switch o := other.(type) {
	case FarmHash64CompositeKey:
		return uarray.EqualsWithOrder(k.keys, o.keys)
	default:
		return false
	}
}

func (k FarmHash64CompositeKey) Keys() []Unique {
	result := make([]Unique, len(k.keys))
	for i, key := range k.keys {
		result[i] = IntKey(key.Key())
	}

	return result
}

func (k FarmHash64CompositeKey) String() string {
	builder := strings.Builder{}
	for _, key := range k.keys {
		builder.WriteString(convertToString(key))
	}

	return builder.String()
}

func convertToString[T comparable](input T) string {
	switch value := any(input).(type) {
	case string:
		return value
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.FormatInt(int64(value), 10)
	case int16:
		return strconv.FormatInt(int64(value), 10)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case bool:
		return strconv.FormatBool(value)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case complex64:
		return strconv.FormatComplex(complex128(value), 'f', -1, 64)
	case complex128:
		return strconv.FormatComplex(value, 'f', -1, 128)
	default:
		return ""
	}
}

func isComparable(value interface{}) bool {
	t := reflect.TypeOf(value)

	// Handling pointers by getting the underlying element type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Slice, reflect.Map, reflect.Func:
		return false
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if !isComparable(reflect.Zero(field.Type).Interface()) {
				return false
			}
		}
	default:
		return true
	}

	return true
}
