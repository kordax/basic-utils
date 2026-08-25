package uos

import (
	"net/url"
	"time"

	"github.com/kordax/basic-utils/v4/uconst"
	"github.com/kordax/basic-utils/v4/uopt"
)

// GetEnvOpt is an alias to GetEnvOptAs[string](key, MapString).
// Behaves like RequireEnv, but returns an empty Opt instead of panicking.
func GetEnvOpt(key string) uopt.Opt[string] {
	return GetEnvOptAs[string](key, MapString)
}

// GetEnvOptAs retrieves an environment variable and maps it to type T using MappingFunc.
// If the variable is missing or mapping fails, an empty Opt is returned.
//
// This is the non-panicking variant of RequireEnvAs.
func GetEnvOptAs[T any](key string, f MappingFunc[T]) uopt.Opt[T] {
	as, err := GetEnvAs[T](key, f)
	if err != nil {
		return uopt.Null[T]()
	}

	return uopt.OfNullable(as)
}

// GetEnvOptNumeric retrieves an environment variable as a numeric type.
// Returns an empty Opt if the variable is missing or cannot be parsed.
//
// This is the non-panicking variant of RequireEnvNumeric.
func GetEnvOptNumeric[T uconst.Numeric](key string) uopt.Opt[T] {
	return GetEnvOptAs(key, MapStringToNumeric[T])
}

// GetEnvOptDuration retrieves an environment variable as time.Duration.
// Returns an empty Opt if the variable is missing or cannot be parsed.
//
// This is the non-panicking variant of RequireEnvDuration.
func GetEnvOptDuration(key string) uopt.Opt[time.Duration] {
	return GetEnvOptAs[time.Duration](key, MapStringToDuration)
}

// GetEnvOptTime retrieves an environment variable as time.Time using the provided layout.
// Returns an empty Opt if the variable is missing or cannot be parsed.
//
// This is the non-panicking variant of RequireEnvTime.
func GetEnvOptTime(key string, layout string) uopt.Opt[time.Time] {
	return GetEnvOptAs[time.Time](key, MapStringToTime(layout))
}

// GetEnvOptURL retrieves an environment variable as url.URL.
// Returns an empty Opt if the variable is missing or cannot be parsed.
//
// This is the non-panicking variant of RequireEnvURL.
func GetEnvOptURL(key string) uopt.Opt[url.URL] {
	return GetEnvOptAs[url.URL](key, MapStringToURL)
}

// GetEnvOptBool retrieves an environment variable as bool.
// Returns an empty Opt if the variable is missing or cannot be parsed.
//
// This is the non-panicking variant of RequireEnvBool.
func GetEnvOptBool(key string) uopt.Opt[bool] {
	return GetEnvOptAs[bool](key, MapStringToBool)
}
