/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2023.
 */

package uopt

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	basicutils "github.com/kordax/basic-utils/v4/uconst"
	"github.com/kordax/basic-utils/v4/uref"
)

// Opt represents a generic container for optional values.
// An Opt can either contain a value of type T or no value at all.
// It provides methods to check the presence of a value and to retrieve it.
// This struct offers a safer way to handle potentially absent values, avoiding nil dereferences.
// Using Opt ensures that the user must explicitly handle both the present and absent cases,
// thus preventing unintentional null pointer errors.
//
// It's similar in principle to "Optional" in other languages like Java's java.util.Optional.
//
// The internal 'v' field is a pointer to a value of type T.
// If 'v' is nil, it means the Opt contains no value.
// Otherwise, 'v' points to the contained value.
type Opt[T any] struct {
	v *T
}

// Present checks if the Opt contains a value.
func (o Opt[T]) Present() bool {
	return o.v != nil
}

// IfPresent invokes the provided function if the Opt contains a value.
func (o Opt[T]) IfPresent(f func(t T)) {
	if o.Present() {
		f(*o.v)
	}
}

// Map transforms a present optional value, or returns Null when the optional is empty.
func (o Opt[T]) Map[R any](mapper func(T) R) Opt[R] {
	if !o.Present() {
		return Null[R]()
	}

	return Of(mapper(*o.v))
}

// FlatMap transforms a present value into another optional value, or returns Null when the optional is empty.
func (o Opt[T]) FlatMap[R any](mapper func(T) Opt[R]) Opt[R] {
	if !o.Present() {
		return Null[R]()
	}

	return mapper(*o.v)
}

// Filter keeps a present value only when predicate returns true.
func (o Opt[T]) Filter(predicate func(T) bool) Opt[T] {
	if !o.Present() || !predicate(*o.v) {
		return Null[T]()
	}

	return o
}

// Map transforms a present optional value, or returns Null when opt is empty.
// Deprecated: use opt.Map(mapper).
func Map[T, R any](opt Opt[T], mapper func(T) R) Opt[R] {
	return opt.Map(mapper)
}

// FlatMap transforms a present optional value into another optional value, or returns Null when opt is empty.
// Deprecated: use opt.FlatMap(mapper).
func FlatMap[T, R any](opt Opt[T], mapper func(T) Opt[R]) Opt[R] {
	return opt.FlatMap(mapper)
}

// Filter keeps a present value only when predicate returns true.
// Deprecated: use opt.Filter(predicate).
func Filter[T any](opt Opt[T], predicate func(T) bool) Opt[T] {
	return opt.Filter(predicate)
}

// Null creates an Opt with no value.
func Null[T any]() Opt[T] {
	return Opt[T]{v: nil}
}

// Of creates an Opt with a value.
func Of[T any](v T) Opt[T] {
	return Opt[T]{
		v: &v,
	}
}

// OfNullable creates an Opt that may or may not contain a value based on the provided pointer.
func OfNullable[T any](v *T) Opt[T] {
	if v != nil {
		return Opt[T]{
			v: uref.Ref(*v),
		}
	} else {
		return Opt[T]{
			v: nil,
		}
	}
}

// OfString creates an Opt containing a string, or a null Opt if the string is empty.
func OfString(v string) Opt[string] {
	if v == "" {
		return Null[string]()
	}

	return Opt[string]{
		v: &v,
	}
}

// OfBool creates an Opt containing a boolean, or a null Opt if the boolean is false.
func OfBool(v bool) Opt[bool] {
	if !v {
		return Null[bool]()
	}

	return Opt[bool]{
		v: &v,
	}
}

// OfNumeric creates an Opt containing a numeric value, or a null Opt if the value is 0.
func OfNumeric[T basicutils.Numeric](v T) Opt[T] {
	if v == 0 {
		return Null[T]()
	}

	return Opt[T]{
		v: &v,
	}
}

// OfDuration creates an Opt containing a time.Duration, or a null Opt if the duration is 0.
func OfDuration(v time.Duration) Opt[time.Duration] {
	if v == 0 {
		return Null[time.Duration]()
	}

	return Opt[time.Duration]{
		v: &v,
	}
}

// OfTime creates an Opt containing a time.Time, or a null Opt if the time is zero.
func OfTime(v time.Time) Opt[time.Time] {
	if v.IsZero() {
		return Null[time.Time]()
	}

	return Opt[time.Time]{
		v: &v,
	}
}

// OfCond creates an Opt containing a value if a given condition is met.
func OfCond[T any](v T, cond func(v T) bool) Opt[T] {
	if cond(v) {
		return Opt[T]{
			v: &v,
		}
	}

	return Null[T]()
}

// OfUnix creates an Opt containing a time.Time value based on a Unix timestamp.
func OfUnix[T basicutils.SignedNumeric](v T) Opt[time.Time] {
	return Opt[time.Time]{
		v: uref.Ref(time.Unix(int64(v), 0)),
	}
}

// OfBuilder creates an Opt by invoking the provided builder function to generate a value.
func OfBuilder[T any](build func() T) Opt[T] {
	v := build()
	return Opt[T]{
		v: &v,
	}
}

// OrElse retrieves the value within the Opt or a provided default if the Opt is null.
func (o Opt[T]) OrElse(v T) T {
	if o.v == nil {
		return v
	} else {
		return *o.v
	}
}

// OrElseGet retrieves the value within the Opt or invokes fallback when the Opt is null.
func (o Opt[T]) OrElseGet(fallback func() T) T {
	if o.v == nil {
		return fallback()
	}

	return *o.v
}

// Get retrieves the value within the Opt as a pointer.
func (o Opt[T]) Get() *T {
	return o.v
}

// Def behaves as Get, but the returns default value if value is not present,
// Def is an alias to OrElse(*new(T))
func (o Opt[T]) Def() T {
	return o.OrElse(*new(T))
}

// Set sets the value within the Opt.
func (o *Opt[T]) Set(v *T) {
	o.v = v
}

// GetAs retrieves the value within the Opt after applying a mapping function.
func (o Opt[T]) GetAs(mapping func(t T) any) any {
	return mapping(*o.v)
}

// UnmarshalJSON implements the json.Unmarshaler interface for the Opt type.
func (o *Opt[T]) UnmarshalJSON(bytes []byte) error {
	var v T
	if (string)(bytes) == "null" {
		o.v = nil
		return nil
	}
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return err
	}
	o.v = &v

	return nil
}

// MarshalJSON implements the json.Marshaler interface for the Opt type.
func (o Opt[T]) MarshalJSON() ([]byte, error) {
	if !o.Present() {
		return []byte("null"), nil
	}

	return json.Marshal(o.Get())
}

// Value implements the driver.Valuer interface for the Opt type, converting its value to a SQL value.
func (o Opt[T]) Value() (driver.Value, error) {
	if o.v == nil {
		return nil, nil
	}

	return sqlValue(*o.v, o.v)
}

// Scan implements the sql.Scanner interface for the Opt type, reading a SQL value into the Opt.
func (o *Opt[T]) Scan(src interface{}) error {
	if src == nil {
		*o = Null[T]()
		return nil
	}

	v, err := scanSQLValue[T](src)
	if err != nil {
		return err
	}
	*o = Of(v)

	return nil
}

func sqlValue[T any](value T, valuePtr *T) (driver.Value, error) {
	if valuer, ok := any(value).(driver.Valuer); ok {
		return valuer.Value()
	}
	if valuer, ok := any(valuePtr).(driver.Valuer); ok {
		return valuer.Value()
	}

	switch v := any(value).(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return uintToInt64(uint64(v))
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return uintToInt64(v)
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case bool:
		return v, nil
	case string:
		return v, nil
	case []byte:
		return v, nil
	case time.Time:
		return v, nil
	case time.Duration:
		return int64(v), nil
	case json.RawMessage:
		return []byte(v), nil
	}

	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return nil, nil
	}
	switch rv.Kind() {
	case reflect.Map, reflect.Struct, reflect.Slice, reflect.Array:
		return json.Marshal(value)
	default:
		if driver.IsValue(value) {
			return value, nil
		}
		return json.Marshal(value)
	}
}

func uintToInt64(v uint64) (driver.Value, error) {
	const maxInt64 = uint64(1<<63 - 1)
	if v > maxInt64 {
		return nil, fmt.Errorf("uint value %d overflows int64 SQL value", v)
	}
	return int64(v), nil
}

func scanSQLValue[T any](src any) (T, error) {
	var zero T
	if valuer, ok := src.(driver.Valuer); ok {
		value, err := valuer.Value()
		if err != nil {
			return zero, fmt.Errorf("incompatible type for Opt[%T]: %T, failed to retrieve value", zero, reflect.TypeOf(src))
		}
		src = value
	}
	if src == nil {
		return zero, nil
	}

	var result T
	if scanner, ok := any(&result).(sql.Scanner); ok {
		if err := scanner.Scan(src); err != nil {
			return zero, err
		}
		return result, nil
	}

	if converted, ok := directConvert[T](src); ok {
		return converted, nil
	}

	target := reflect.TypeOf((*T)(nil)).Elem()
	if target == reflect.TypeOf(time.Time{}) {
		value, err := scanTime(src)
		if err != nil {
			return zero, err
		}
		return any(value).(T), nil
	}

	if target == reflect.TypeOf(time.Duration(0)) {
		value, err := scanDuration(src)
		if err != nil {
			return zero, err
		}
		return any(value).(T), nil
	}

	if converted, ok, err := scanScalar[T](src, target); ok || err != nil {
		return converted, err
	}

	if text, ok := sqlText(src); ok {
		return scanText[T](text, target, src)
	}

	return zero, fmt.Errorf("incompatible type for Opt[%T]: %T", zero, src)
}

func directConvert[T any](src any) (T, bool) {
	var zero T
	target := reflect.TypeOf((*T)(nil)).Elem()
	source := reflect.ValueOf(src)
	if !source.IsValid() {
		return zero, false
	}
	if target.Kind() == reflect.String && source.Kind() != reflect.String {
		return zero, false
	}
	var result T
	resultValue := reflect.ValueOf(&result).Elem()
	if source.Type().AssignableTo(target) {
		resultValue.Set(source)
		return result, true
	}
	if source.Type().ConvertibleTo(target) {
		resultValue.Set(source.Convert(target))
		return result, true
	}
	return zero, false
}

func scanScalar[T any](src any, target reflect.Type) (T, bool, error) {
	var zero T
	source := reflect.ValueOf(src)
	if !source.IsValid() {
		return zero, false, nil
	}

	switch target.Kind() {
	case reflect.String:
		value, ok := scalarString(source)
		if !ok {
			return zero, false, nil
		}
		return any(value).(T), true, nil
	case reflect.Bool:
		value, ok := scalarBool(source)
		if !ok {
			return zero, false, nil
		}
		return reflect.ValueOf(value).Convert(target).Interface().(T), true, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, ok := scalarInt(source)
		if !ok {
			return zero, false, nil
		}
		result := reflect.New(target).Elem()
		if result.OverflowInt(value) {
			return zero, true, fmt.Errorf("value %d overflows Opt[%T]", value, zero)
		}
		result.SetInt(value)
		return result.Interface().(T), true, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, ok := scalarUint(source)
		if !ok {
			return zero, false, nil
		}
		result := reflect.New(target).Elem()
		if result.OverflowUint(value) {
			return zero, true, fmt.Errorf("value %d overflows Opt[%T]", value, zero)
		}
		result.SetUint(value)
		return result.Interface().(T), true, nil
	case reflect.Float32, reflect.Float64:
		value, ok := scalarFloat(source)
		if !ok {
			return zero, false, nil
		}
		result := reflect.New(target).Elem()
		if result.OverflowFloat(value) {
			return zero, true, fmt.Errorf("value %f overflows Opt[%T]", value, zero)
		}
		result.SetFloat(value)
		return result.Interface().(T), true, nil
	default:
		return zero, false, nil
	}
}

func scalarString(value reflect.Value) (string, bool) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10), true
	case reflect.Float32:
		return strconv.FormatFloat(value.Float(), 'f', -1, 32), true
	case reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64), true
	case reflect.Bool:
		return strconv.FormatBool(value.Bool()), true
	default:
		return "", false
	}
}

func scalarBool(value reflect.Value) (bool, bool) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() != 0, true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint() != 0, true
	case reflect.Float32, reflect.Float64:
		return value.Float() != 0, true
	default:
		return false, false
	}
}

func scalarInt(value reflect.Value) (int64, bool) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int(), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue := value.Uint()
		const maxInt64 = uint64(1<<63 - 1)
		if uintValue > maxInt64 {
			return 0, false
		}

		return int64(uintValue), true
	case reflect.Float32, reflect.Float64:
		floatValue := value.Float()
		if floatValue != floatValue || floatValue < -0x1p63 || floatValue >= 0x1p63 {
			return 0, false
		}

		return int64(floatValue), true
	default:
		return 0, false
	}
}

func scalarUint(value reflect.Value) (uint64, bool) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue := value.Int()
		if intValue < 0 {
			return 0, false
		}

		return uint64(intValue), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint(), true
	case reflect.Float32, reflect.Float64:
		floatValue := value.Float()
		if floatValue != floatValue || floatValue < 0 || floatValue >= 0x1p64 {
			return 0, false
		}

		return uint64(floatValue), true
	default:
		return 0, false
	}
}

func scalarFloat(value reflect.Value) (float64, bool) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(value.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(value.Uint()), true
	case reflect.Float32, reflect.Float64:
		return value.Float(), true
	default:
		return 0, false
	}
}

func scanText[T any](text string, target reflect.Type, src any) (T, error) {
	var zero T
	switch target.Kind() {
	case reflect.String:
		return any(text).(T), nil
	case reflect.Bool:
		value, err := parseSQLBool(text)
		if err != nil {
			return zero, parseError(src, "bool", err)
		}
		return reflect.ValueOf(value).Convert(target).Interface().(T), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(text, 10, target.Bits())
		if err != nil {
			return zero, parseError(src, "numeric", err)
		}
		return reflect.ValueOf(value).Convert(target).Interface().(T), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.ParseUint(text, 10, target.Bits())
		if err != nil {
			return zero, parseError(src, "numeric", err)
		}
		return reflect.ValueOf(value).Convert(target).Interface().(T), nil
	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(text, target.Bits())
		if err != nil {
			return zero, parseError(src, "float", err)
		}
		return reflect.ValueOf(value).Convert(target).Interface().(T), nil
	case reflect.Complex64, reflect.Complex128:
		value, err := strconv.ParseComplex(text, target.Bits())
		if err != nil {
			return zero, parseError(src, "complex", err)
		}
		return reflect.ValueOf(value).Convert(target).Interface().(T), nil
	case reflect.Slice, reflect.Array:
		if value, ok, err := scanPostgresArrayText[T](text, target); ok || err != nil {
			return value, err
		}

		var result T
		if err := json.Unmarshal([]byte(text), &result); err != nil {
			return zero, err
		}
		return result, nil
	case reflect.Map, reflect.Struct, reflect.Pointer, reflect.Interface:
		var result T
		if err := json.Unmarshal([]byte(text), &result); err != nil {
			return zero, err
		}
		return result, nil
	default:
		return zero, fmt.Errorf("incompatible type for Opt[%T]: %T, failed to retrieve value", zero, reflect.TypeOf(src))
	}
}

type postgresArrayElement struct {
	value string
	null  bool
}

func scanPostgresArrayText[T any](text string, target reflect.Type) (T, bool, error) {
	var zero T
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "{") {
		return zero, false, nil
	}

	elements, err := parsePostgresArrayElements(text)
	if err != nil {
		return zero, true, err
	}

	if target.Kind() == reflect.Array && len(elements) != target.Len() {
		return zero, true, fmt.Errorf("postgres array length %d does not match Opt[%T] array length", len(elements), zero)
	}

	elemType := target.Elem()
	resultType := target
	if target.Kind() == reflect.Array {
		resultType = reflect.SliceOf(elemType)
	}
	result := reflect.MakeSlice(resultType, len(elements), len(elements))
	for i, element := range elements {
		value, err := scanPostgresArrayElement(element, elemType)
		if err != nil {
			return zero, true, err
		}
		result.Index(i).Set(value)
	}

	if target.Kind() == reflect.Array {
		array := reflect.New(target).Elem()
		reflect.Copy(array.Slice(0, target.Len()), result)
		return array.Interface().(T), true, nil
	}

	return result.Interface().(T), true, nil
}

func parsePostgresArrayElements(text string) ([]postgresArrayElement, error) {
	if len(text) < 2 || text[0] != '{' || text[len(text)-1] != '}' {
		return nil, fmt.Errorf("invalid postgres array literal %q", text)
	}
	if text == "{}" {
		return []postgresArrayElement{}, nil
	}

	var elements []postgresArrayElement
	var builder strings.Builder
	quoted := false
	wasQuoted := false
	escaped := false
	for i := 1; i < len(text)-1; i++ {
		ch := text[i]
		switch {
		case escaped:
			builder.WriteByte(ch)
			escaped = false
		case quoted && ch == '\\':
			escaped = true
		case ch == '"':
			quoted = !quoted
			wasQuoted = true
		case !quoted && ch == ',':
			elements = append(elements, postgresArrayElement{
				value: builder.String(),
				null:  !wasQuoted && builder.String() == "NULL",
			})
			builder.Reset()
			wasQuoted = false
		default:
			builder.WriteByte(ch)
		}
	}
	if quoted || escaped {
		return nil, fmt.Errorf("invalid postgres array literal %q", text)
	}
	elements = append(elements, postgresArrayElement{
		value: builder.String(),
		null:  !wasQuoted && builder.String() == "NULL",
	})

	return elements, nil
}

func scanPostgresArrayElement(element postgresArrayElement, target reflect.Type) (reflect.Value, error) {
	result := reflect.New(target).Elem()
	if element.null {
		return result, nil
	}

	if target == reflect.TypeOf(time.Time{}) {
		value, err := scanTime(element.value)
		if err != nil {
			return result, err
		}
		result.Set(reflect.ValueOf(value))
		return result, nil
	}

	switch target.Kind() {
	case reflect.String:
		result.SetString(element.value)
	case reflect.Bool:
		value, err := parseSQLBool(element.value)
		if err != nil {
			return result, err
		}
		result.SetBool(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(element.value, 10, target.Bits())
		if err != nil {
			return result, err
		}
		result.SetInt(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.ParseUint(element.value, 10, target.Bits())
		if err != nil {
			return result, err
		}
		result.SetUint(value)
	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(element.value, target.Bits())
		if err != nil {
			return result, err
		}
		result.SetFloat(value)
	case reflect.Interface:
		result.Set(reflect.ValueOf(element.value))
	case reflect.Pointer:
		value, err := scanPostgresArrayElement(element, target.Elem())
		if err != nil {
			return result, err
		}
		pointer := reflect.New(target.Elem())
		pointer.Elem().Set(value)
		result.Set(pointer)
	default:
		if err := json.Unmarshal([]byte(element.value), result.Addr().Interface()); err != nil {
			return result, err
		}
	}

	return result, nil
}

func sqlText(src any) (string, bool) {
	switch value := src.(type) {
	case string:
		return value, true
	case []byte:
		return string(value), true
	default:
		return "", false
	}
}

func parseSQLBool(text string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(text)) {
	case "1", "t", "true", "y", "yes", "on":
		return true, nil
	case "0", "f", "false", "n", "no", "off":
		return false, nil
	default:
		return strconv.ParseBool(text)
	}
}

func scanTime(src any) (time.Time, error) {
	if value, ok := src.(time.Time); ok {
		return value, nil
	}
	text, ok := sqlText(src)
	if !ok {
		return time.Time{}, fmt.Errorf("incompatible type for Opt[%T]: %T", time.Time{}, src)
	}
	for _, layout := range []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
		"2006-01-02",
		time.DateOnly,
		time.DateTime,
	} {
		value, err := time.Parse(layout, text)
		if err == nil {
			return value, nil
		}
	}
	return time.Time{}, fmt.Errorf("failed to parse sql time value %q", text)
}

func scanDuration(src any) (time.Duration, error) {
	switch value := src.(type) {
	case int64:
		return time.Duration(value), nil
	case float64:
		return time.Duration(value), nil
	}
	text, ok := sqlText(src)
	if !ok {
		return 0, fmt.Errorf("incompatible type for Opt[%T]: %T", time.Duration(0), src)
	}
	if value, err := time.ParseDuration(text); err == nil {
		return value, nil
	}
	value, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return 0, parseError(src, "numeric", err)
	}
	return time.Duration(value), nil
}

func parseError(src any, kind string, err error) error {
	if _, ok := src.([]byte); ok {
		if kind == "float" {
			return fmt.Errorf("failed to parse bytes/blob sql value to float opt: %s", err)
		}
		return fmt.Errorf("failed to parse bytes/blob sql value to %s opt: %s", kind, err)
	}
	if kind == "float" {
		return fmt.Errorf("failed to parse varchar sql value to float opt: %s", err)
	}
	return fmt.Errorf("failed to parse varchar sql value to %s opt: %s", kind, err)
}
