package ucast

import (
	"fmt"
	"reflect"

	"github.com/kordax/basic-utils/uconst"
)

func String[R uconst.BasicType](str string) (R, error) {
	var zero R
	return StringOrDef(str, zero)
}

func Type[V uconst.BasicType](v V) string {
	return toString[V](v)
}

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
		valPtr := &val
		return BoolToString(valPtr)
	case *bool:
		if val == nil {
			return ""
		}
		return BoolToString(val)
	case int:
		valPtr := &val
		return IntToString(valPtr)
	case *int:
		if val == nil {
			return ""
		}
		return IntToString(val)
	case int8:
		valPtr := &val
		return Int8ToString(valPtr)
	case *int8:
		if val == nil {
			return ""
		}
		return Int8ToString(val)
	case int16:
		valPtr := &val
		return Int16ToString(valPtr)
	case *int16:
		if val == nil {
			return ""
		}
		return Int16ToString(val)
	case int32:
		valPtr := &val
		return Int32ToString(valPtr)
	case *int32:
		if val == nil {
			return ""
		}
		return Int32ToString(val)
	case int64:
		valPtr := &val
		return Int64ToString(valPtr)
	case *int64:
		if val == nil {
			return ""
		}
		return Int64ToString(val)
	case uint:
		valPtr := &val
		return UintToString(valPtr)
	case *uint:
		if val == nil {
			return ""
		}
		return UintToString(val)
	case uint8:
		valPtr := &val
		return Uint8ToString(valPtr)
	case *uint8:
		if val == nil {
			return ""
		}
		return Uint8ToString(val)
	case uint16:
		valPtr := &val
		return Uint16ToString(valPtr)
	case *uint16:
		if val == nil {
			return ""
		}
		return Uint16ToString(val)
	case uint32:
		valPtr := &val
		return Uint32ToString(valPtr)
	case *uint32:
		if val == nil {
			return ""
		}
		return Uint32ToString(val)
	case uint64:
		valPtr := &val
		return Uint64ToString(valPtr)
	case *uint64:
		if val == nil {
			return ""
		}
		return Uint64ToString(val)
	case float32:
		valPtr := &val
		return Float32ToString(valPtr)
	case *float32:
		if val == nil {
			return ""
		}
		return Float32ToString(val)
	case float64:
		valPtr := &val
		return Float64ToString(valPtr)
	case *float64:
		if val == nil {
			return ""
		}
		return Float64ToString(val)
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

	strPtr := &s

	switch uType.Kind() {
	case reflect.String:
		value = s
	case reflect.Bool:
		value, err = StringToBool(strPtr)
	case reflect.Int:
		value, err = StringToInt(strPtr)
	case reflect.Int8:
		value, err = StringToInt8(strPtr)
	case reflect.Int16:
		value, err = StringToInt16(strPtr)
	case reflect.Int32:
		value, err = StringToInt32(strPtr)
	case reflect.Int64:
		value, err = StringToInt64(strPtr)
	case reflect.Uint:
		value, err = StringToUint(strPtr)
	case reflect.Uint8:
		value, err = StringToUint8(strPtr)
	case reflect.Uint16:
		value, err = StringToUint16(strPtr)
	case reflect.Uint32:
		value, err = StringToUint32(strPtr)
	case reflect.Uint64:
		value, err = StringToUint64(strPtr)
	case reflect.Float32:
		value, err = StringToFloat32(strPtr)
	case reflect.Float64:
		value, err = StringToFloat64(strPtr)
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
