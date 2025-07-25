package ucast

import (
	"fmt"
	"reflect"

	"github.com/kordax/basic-utils/v2/uconst"
)

// String converts the input string to a value of type R.
// It returns the converted value and an error if the conversion fails.
//
// The type R must satisfy the uconst.BasicType constraint.
//
// Example usage:
//
//	value, err := ucast.String[int]("123")
//	if err != nil {
//	    // handle error
//	}
func String[R uconst.BasicType](str string) (R, error) {
	var zero R
	return StringOrDef(str, zero)
}

// Type converts the input value v of type V to its string representation.
//
// The type V must satisfy the uconst.BasicType constraint.
//
// Example usage:
//
//	str := ucast.Type       // "123"
//	str := ucast.Type[float64](3.14)  // "3.14"
func Type[V uconst.BasicType](v V) string {
	return toString[V](v)
}

// StringOrDef attempts to convert the input string to a value of type R.
// If the conversion fails, it returns the default value def and an error.
//
// The type R must satisfy the uconst.BasicType constraint.
//
// Example usage:
//
//	value, err := ucast.StringOrDef[int]("invalid", 42)
//	// value == 42, err != nil
func StringOrDef[R uconst.BasicType](str string, def R) (R, error) {
	result, err := fromString[R](str)
	if err != nil {
		return def, fmt.Errorf("failed to convert string to target type: %v", err)
	}

	return result, nil
}

func toString[V uconst.BasicType](v V) string {
	switch val := any(v).(type) {
	case string:
		return val
	case *string:
		if val == nil {
			return ""
		}
		return *val
	case bool:
		return BoolToString(val)
	case *bool:
		if val == nil {
			return ""
		}
		return BoolToString(*val)
	case int:
		return IntToString(val)
	case *int:
		if val == nil {
			return ""
		}
		return IntToString(*val)
	case int8:
		return Int8ToString(val)
	case *int8:
		if val == nil {
			return ""
		}
		return Int8ToString(*val)
	case int16:
		return Int16ToString(val)
	case *int16:
		if val == nil {
			return ""
		}
		return Int16ToString(*val)
	case int32:
		return Int32ToString(val)
	case *int32:
		if val == nil {
			return ""
		}
		return Int32ToString(*val)
	case int64:
		return Int64ToString(val)
	case *int64:
		if val == nil {
			return ""
		}
		return Int64ToString(*val)
	case uint:
		return UintToString(val)
	case *uint:
		if val == nil {
			return ""
		}
		return UintToString(*val)
	case uint8:
		return Uint8ToString(val)
	case *uint8:
		if val == nil {
			return ""
		}
		return Uint8ToString(*val)
	case uint16:
		return Uint16ToString(val)
	case *uint16:
		if val == nil {
			return ""
		}
		return Uint16ToString(*val)
	case uint32:
		return Uint32ToString(val)
	case *uint32:
		if val == nil {
			return ""
		}
		return Uint32ToString(*val)
	case uint64:
		return Uint64ToString(val)
	case *uint64:
		if val == nil {
			return ""
		}
		return Uint64ToString(*val)
	case float32:
		return Float32ToString(val)
	case *float32:
		if val == nil {
			return ""
		}
		return Float32ToString(*val)
	case float64:
		return Float64ToString(val)
	case *float64:
		if val == nil {
			return ""
		}
		return Float64ToString(*val)
	default:
		return ""
	}
}

func fromString[U uconst.BasicType](s string) (U, error) {
	var zero U
	var uType = reflect.TypeOf(zero)

	isPtr := uType.Kind() == reflect.Ptr
	if isPtr {
		uType = uType.Elem()
	}

	var value interface{}
	var err error

	switch uType.Kind() {
	case reflect.String:
		value = s
	case reflect.Bool:
		value, err = StringToBool(s)
	case reflect.Int:
		value, err = StringToInt(s)
	case reflect.Int8:
		value, err = StringToInt8(s)
	case reflect.Int16:
		value, err = StringToInt16(s)
	case reflect.Int32:
		value, err = StringToInt32(s)
	case reflect.Int64:
		value, err = StringToInt64(s)
	case reflect.Uint:
		value, err = StringToUint(s)
	case reflect.Uint8:
		value, err = StringToUint8(s)
	case reflect.Uint16:
		value, err = StringToUint16(s)
	case reflect.Uint32:
		value, err = StringToUint32(s)
	case reflect.Uint64:
		value, err = StringToUint64(s)
	case reflect.Float32:
		value, err = StringToFloat32(s)
	case reflect.Float64:
		value, err = StringToFloat64(s)
	default:
		return zero, fmt.Errorf("unsupported target type: %v", uType)
	}

	if err != nil {
		return zero, err
	}

	if isPtr {
		ptrValue := reflect.New(uType)
		ptrValue.Elem().Set(reflect.ValueOf(value))
		return ptrValue.Interface().(U), nil
	}

	return value.(U), nil
}
