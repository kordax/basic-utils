package cache

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	arrayutils "github.com/kordax/basic-utils/array-utils"
)

type UIntKey uint64
type IntKey int64
type StringKey string

/*
Comparable entity
*/
type Comparable interface {
	Equals(other Comparable) bool
}

/*
Hashed specifies an abstract key with an ability to provide keys for hash method. Shouldn't implement hashing at all.
*/
type Hashed interface {
	Comparable
	Key() int // Key should return a unique item key for hash use.
	String() string
}

/*
CompositeKey specifies an abstract key with an ability to provide an ordered list of available keys.
*/
type CompositeKey interface {
	Comparable
	Keys() []int // Keys returns an ordered list of keys ordered by priority (ASC), so the first element has the most prio.
	String() string
}

func (s IntKey) Equals(other Comparable) bool {
	switch other.(type) {
	case IntKey:
		return s == other
	default:
		return false
	}
}

func (s IntKey) Key() int {
	return int(s)
}

func (s StringKey) Equals(other Comparable) bool {
	switch other.(type) {
	case StringKey:
		return s == other
	default:
		return false
	}
}

func (s StringKey) Key() int {
	hash := 0
	for i := 0; i < len(s); i++ {
		hash = 31*hash + int(s[i])
	}

	return hash
}

func (s StringKey) String() string {
	return string(s)
}

type ComparableKey[T comparable] struct {
	v T
}

func NewComparableKey[T comparable](v T) ComparableKey[T] {
	return ComparableKey[T]{v: v}
}

func (k ComparableKey[T]) Hash() int {
	hash := 0

	switch value := any(k.v).(type) {
	case int, int8, int16, int32, int64:
		hash = 31*hash + int(reflect.ValueOf(value).Int())
	case uint, uint8, uint16, uint32, uint64:
		hash = 31*hash + int(reflect.ValueOf(value).Uint())
	case float32, float64:
		hash = 31*hash + int(reflect.ValueOf(value).Float())
	case complex64, complex128:
		c := reflect.ValueOf(value).Complex()
		hash = 31*(31*hash+int(real(c))) + int(imag(c))
	case string:
		for i := 0; i < len(value); i++ {
			hash = 31*hash + int(value[i])
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
		return arrayutils.EqualsWithOrder(k.keys, o.keys)
	default:
		return false
	}
}

func (k UIntCompositeKey) Keys() []int {
	result := make([]int, len(k.keys))
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
		return arrayutils.EqualsWithOrder(k.keys, o.keys)
	default:
		return false
	}
}

func (k IntCompositeKey) Keys() []int {
	result := make([]int, len(k.keys))
	for i, key := range k.keys {
		result[i] = key.Key()
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
		return arrayutils.EqualsWithOrder(k.keys, o.keys)
	default:
		return false
	}
}

func (k StrCompositeKey) Keys() []int {
	result := make([]int, len(k.keys))
	for i, key := range k.keys {
		result[i] = key.Key()
	}

	return result
}

func (k StrCompositeKey) String() string {
	return strings.Join(arrayutils.Map(k.keys, func(v *StringKey) string {
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
	switch other.(type) {
	case StrCompositeKey:
		return arrayutils.EqualsWithOrder(k.keys, other.(GenericCompositeKey).keys)
	default:
		return false
	}
}

func (k GenericCompositeKey) Keys() []int {
	result := make([]int, len(k.keys))
	for i, key := range k.keys {
		result[i] = key.Hash()
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

		return arrayutils.EqualValues(s.v, otherValue.v)
	}

	return arrayutils.EqualValues(s.v, otherValuePtr.v)
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

type GenKey interface {
	int | ~string
}

func convertToString[T comparable](input T) string {
	switch value := any(input).(type) {
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
