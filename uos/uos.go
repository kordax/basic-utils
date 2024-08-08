/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uos

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/kordax/basic-utils/uarray"
	basicutils "github.com/kordax/basic-utils/uconst"
	"github.com/kordax/basic-utils/umath"
)

// GetCPUs calculates and returns the number of CPU cores available to the application.
// This function is designed to provide a more accurate count of available CPUs, especially when running
// in a containerized environment such as Docker or Kubernetes on Linux, where CPU resources can be limited.
//
// The function operates as follows:
//
//  1. Checks the operating system where the application is running. If it's Linux, it proceeds to check
//     the cgroup settings that are commonly used in containerized environments to limit CPU resources.
//
//  2. For Linux:
//     a. Attempts to read the CPU allocation from the cgroup settings, specifically the cpu.cfs_quota_us and
//     cpu.cfs_period_us files, which define the CPU quota and period for the container.
//     b. If the cgroup CPU quota and period can be successfully read, it calculates the number of CPUs
//     as the quotient of the quota and period. This gives a decimal representation of the number of CPU cores
//     allocated to the container. This value is then returned as an integer.
//     c. If the cgroup files are not accessible or an error occurs while reading them (which might happen if the
//     application is not running in a containerized environment), the error is ignored, and the function falls back
//     to the next step.
//
//  3. If the operating system is not Linux or if the cgroup CPU information could not be obtained, the function
//     falls back to using runtime.NumCPU(). This standard Go library function returns the number of logical CPUs
//     available to the current process, as seen by the Go runtime. This count reflects the total number of CPU cores
//     available to the process, which may be the total number of cores on the machine or limited by the OS scheduler,
//     depending on the environment and OS settings.
//
// Return value:
//   - The function returns an integer representing the number of CPU cores that the application should consider available.
//     This number is either derived from the cgroup settings (in containerized Linux environments) or from the runtime
//     information provided by the Go runtime.NumCPU() function (in non-containerized environments or non-Linux operating systems).
//
// Note:
//   - This function is particularly useful in containerized environments where CPU resources might be limited
//     and different from the physical host's total CPU resources. By accounting for these limits, applications
//     can make more informed decisions about resource allocation, concurrency, and parallel processing.
func GetCPUs() int {
	switch runtime.GOOS {
	case "linux":
		if cgroupCPUs, err := getCGroupCPUs(); err == nil {
			return cgroupCPUs
		}
	}

	return runtime.NumCPU()
}

// RequireEnvNumeric retrieves an environment variable specified by `key` and converts it
// to the specified basicutils.Numeric type `T`.
// The function panics if the environment variable is not set, cannot be converted to type `T`,
// or if `T` is not an integer type. It uses the appropriate bit size for parsing to ensure
// values fit into the specified type without overflow.
func RequireEnvNumeric[T basicutils.Numeric](key string) T {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Errorf("expected environment variable '%s' was not found", key))
	}

	var result any
	var err error

	switch reflect.TypeFor[T]().Kind() {
	case reflect.Int:
		var parsed int64
		parsed, err = strconv.ParseInt(value, 10, 0)
		result = int(parsed)
	case reflect.Int8:
		var parsed int64
		parsed, err = strconv.ParseInt(value, 10, 8)
		result = int8(parsed)
	case reflect.Int16:
		var parsed int64
		parsed, err = strconv.ParseInt(value, 10, 16)
		result = int16(parsed)
	case reflect.Int32:
		var parsed int64
		parsed, err = strconv.ParseInt(value, 10, 32)
		result = int32(parsed)
	case reflect.Int64:
		result, err = strconv.ParseInt(value, 10, 64)
	case reflect.Uint:
		var parsed uint64
		parsed, err = strconv.ParseUint(value, 10, 0)
		result = uint(parsed)
	case reflect.Uint8:
		var parsed uint64
		parsed, err = strconv.ParseUint(value, 10, 8)
		result = uint8(parsed)
	case reflect.Uint16:
		var parsed uint64
		parsed, err = strconv.ParseUint(value, 10, 16)
		result = uint16(parsed)
	case reflect.Uint32:
		var parsed uint64
		parsed, err = strconv.ParseUint(value, 10, 32)
		result = uint32(parsed)
	case reflect.Uint64:
		result, err = strconv.ParseUint(value, 10, 64)
	case reflect.Float32:
		var parsed float64
		parsed, err = strconv.ParseFloat(value, 32)
		result = float32(parsed)
	case reflect.Float64:
		result, err = strconv.ParseFloat(value, 64)
	default:
		panic(fmt.Errorf("failed to parse environment variable '%s', unsupported type for Numeric: %s", key, reflect.TypeFor[T]().Kind()))
	}

	if err != nil {
		panic(fmt.Errorf("failed to parse environment variable '%s' as type %s: %s", key, reflect.TypeOf(*new(T)).Kind(), err))
	}

	return result.(T)
}

// RequireEnv is an alias to RequireEnvAs[string](key, MapString)
func RequireEnv(key string) string {
	return RequireEnvAs[string](key, MapString)
}

// RequireEnvAs retrieves an environment variable specified by `key` and uses a provided
// MappingFunc `f` to convert the environment variable's string value into the desired type `T`.
// The MappingFunc `f` should take a string as input and return a pointer to the type `T` and an error.
// If the environment variable is not set, RequireEnvAs panics with an error indicating that
// the expected environment variable was not found.
// If the MappingFunc `f` returns an error, RequireEnvAs panics with an error indicating
// that the environment variable could not be parsed as the desired type.
//
// Parameters:
//
//	key: the name of the environment variable to retrieve.
//	f: a MappingFunc that converts a string value to the desired type `T`, returning a pointer to `T` and an error.
//
// Returns:
//
//	The value of the environment variable converted to type `T`.
//
// Example usage:
//
//	// Assuming MapStringToDuration and MapStringToInt are defined MappingFuncs for time.Duration and int respectively
//	duration := RequireEnvAs("TIMEOUT", MapStringToDuration(time.RFC3339))
//	url := RequireEnvAs("MY_URL", MapStringToURL)
//
// Note: Since this function uses panic for error handling, it should be used in contexts
// where such behavior is acceptable, or it should be recovered using defer and recover mechanisms.
func RequireEnvAs[T any](key string, f MappingFunc[T]) T {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Errorf("expected environment variable '%s' was not found", key))
	}

	result, err := f(value)
	if err != nil {
		panic(fmt.Errorf("failed to map environment variable '%s'", key))
	}

	return *result
}

// RequireEnvSlice retrieves an environment variable specified by `key` and returns it as a slice of strings.
// The environment variable should contain a list of strings separated by commas, possibly surrounded by spaces.
// This function uses RequireEnvAs to handle the retrieval and conversion of the environment variable,
// leveraging its error handling to manage scenarios where the environment variable is not set or cannot be properly parsed.
//
// Parameters:
//
//	key: the name of the environment variable to retrieve, which should contain a list of strings separated by commas.
//
// Returns:
//
//	A slice of strings parsed from the environment variable, with leading and trailing spaces around each item removed.
//
// Panics:
//   - If the environment variable is not set, this function will panic, indicating that the expected
//     environment variable was not found.
//   - If there is an error during the parsing process, a panic will occur, indicating that the
//     environment variable could not be parsed into a slice of strings.
//
// Example usage:
//
//	// Assuming an environment variable "COLORS" is set to "red, green, blue"
//	colors := RequireEnvSlice("COLORS")
//	fmt.Println(colors) // Output: [red green blue]
//
// Note:
//
//	 This function trims all values, thus removing space characters.
//		This function also uses panic for error handling, which can halt the application unless handled properly.
//		It is recommended to use this function in contexts where such behavior is acceptable, or to employ
//		defer and recover mechanisms to gracefully manage errors and prevent application termination.
func RequireEnvSlice(key string) []string {
	raw := RequireEnvAs[string](key, MapString)
	return uarray.Map(strings.Split(raw, ","), func(v *string) string {
		return strings.TrimSpace(*v)
	})
}

// RequireEnvSliceAs is the same as RequireEnvSlice, but supports any types.
// This func retrieves an environment variable specified by `key`, splits it by commas, trims any spaces,
// and maps each element to a type T using a provided MappingFunc. This function leverages RequireEnvAs for the initial
// retrieval and error handling, ensuring robust handling of missing or improperly formatted environment variables.
//
// Parameters:
//
//	key: the name of the environment variable to retrieve, expected to contain a comma-separated list of strings.
//	f: a MappingFunc that converts a trimmed string value to the desired type `T`, returning `T` and an error.
//
// Returns:
//
//	A slice of `T` representing the parsed and converted elements of the environment variable.
//
// Panics:
//   - If the environment variable is not set, or if any element cannot be successfully converted to type `T`,
//     this function will panic, indicating the specific error encountered.
//
// Example usage:
//
//	// Assuming an environment variable "RATES" is set to "0.5, 2.3, 3.8"
//	rates := RequireEnvSlice[float64]("RATES", MapStringToFloat64)
//	fmt.Println(rates) // Output might be [0.5 2.3 3.8]
//
// Note:
//
//	This function uses panic for error handling, which can halt the application unless handled properly.
//	It is recommended to use this function in contexts where such behavior is acceptable, or to employ
//	defer and recover mechanisms to gracefully manage errors and prevent application termination.
func RequireEnvSliceAs[T any](key string, f MappingFunc[T]) []T {
	parts := RequireEnvSlice(key)
	result := make([]T, 0, len(parts))
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		mappedValue, err := f(trimmedPart)
		if err != nil {
			panic(fmt.Errorf("failed to map environment variable '%s'", key))
		}

		result = append(result, *mappedValue)
	}

	return result
}

// RequireEnvDuration retrieves the environment variable specified by key as a time.Duration.
// This function uses RequireEnvAs under the hood to convert the environment variable string
// to a time.Duration type. If the environment variable is not found, or if the conversion
// fails (e.g., due to an invalid duration format), RequireEnvDuration will panic.
//
// The expected format for the duration is a string accepted by time.ParseDuration,
// which includes any input valid for time.Duration such as "300ms", "1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
//
// Example:
// If you have an environment variable named TIMEOUT with the value "2m30s",
// RequireEnvDuration("TIMEOUT") will return a time.Duration of 2 minutes and 30 seconds.
//
//	os.Setenv("TIMEOUT", "2m30s")
//	timeout := RequireEnvDuration("TIMEOUT")
//	fmt.Println(timeout) // Prints: 2m30s
//
// Note: Because it can panic, this function should be used in cases where the environment variable
// is expected to be set and correctly formatted. For more flexible error handling, consider
// using the underlying RequireEnvAs function directly with appropriate error checks.
func RequireEnvDuration(key string) time.Duration {
	return RequireEnvAs[time.Duration](key, MapStringToDuration)
}

// RequireEnvTime is the same helper as RequireEnvDuration, but for time.Time.
func RequireEnvTime(key string, layout string) time.Time {
	return RequireEnvAs[time.Time](key, MapStringToTime(layout))
}

// RequireEnvURL helper for URL.
func RequireEnvURL(key string) url.URL {
	return RequireEnvAs[url.URL](key, MapStringToURL)
}

// RequireEnvBool helper for bool values.
func RequireEnvBool(key string) bool {
	return RequireEnvAs[bool](key, MapStringToBool)
}

// CheckEnvBool helper is the same as RequireEnvBool, but doesn't panic as all RequireEnv or RequireEnvAs functions.
// It returns true if env variable is set to 'true' (case-insensitive)  returns false otherwise.
func CheckEnvBool(key string) bool {
	value := os.Getenv(key)
	if value == "" {
		return false
	}

	result, err := MapStringToBool(value)
	if err != nil {
		return false
	}

	return *result
}

func getCGroupCPUs() (int, error) { // coverage-ignore
	quota, err := readCgroupValue("/sys/fs/cgroup/cpu/cpu.cfs_quota_us")
	if err != nil {
		return 0, err
	}

	period, err := readCgroupValue("/sys/fs/cgroup/cpu/cpu.cfs_period_us")
	if err != nil {
		return 0, err
	}

	if quota <= 0 {
		return 0, fmt.Errorf("no cgroup limit found or quota is not set")
	}

	cpus := int(umath.RoundUp(quota / period))
	if cpus <= 0 {
		return 0, fmt.Errorf("no quota limit available")
	}

	return cpus, nil
}

func readCgroupValue(path string) (float64, error) { // coverage-ignore
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}
