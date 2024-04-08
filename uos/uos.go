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

	basicutils "github.com/kordax/basic-utils"
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

// GetEnvNumeric retrieves an environment variable specified by `key` and converts it
// to the specified basicutils.Numeric type `T`.
// The function panics if the environment variable is not set, cannot be converted to type `T`,
// or if `T` is not an integer type. It uses the appropriate bit size for parsing to ensure
// values fit into the specified type without overflow.
func GetEnvNumeric[T basicutils.Numeric](key string) T {
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

// GetEnvAs retrieves an environment variable specified by `key` and uses a provided
// MappingFunc `f` to convert the environment variable's string value into the desired type `T`.
// The MappingFunc `f` should take a string as input and return a pointer to the type `T` and an error.
// If the environment variable is not set, GetEnvAs panics with an error indicating that
// the expected environment variable was not found.
// If the MappingFunc `f` returns an error, GetEnvAs panics with an error indicating
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
//	duration := GetEnvAs("TIMEOUT", MapStringToDuration(time.RFC3339))
//	url := GetEnvAs("MY_URL", MapStringToURL)
//
// Note: Since this function uses panic for error handling, it should be used in contexts
// where such behavior is acceptable, or it should be recovered using defer and recover mechanisms.
func GetEnvAs[T any](key string, f MappingFunc[T]) T {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Errorf("expected environment variable '%s' was not found", key))
	}

	result, err := f(value)
	if err != nil {
		panic(fmt.Errorf("failed to parse environment variable '%s' as duration'", key))
	}

	return *result
}

// GetEnvDuration retrieves the environment variable specified by key as a time.Duration.
// This function uses GetEnvAs under the hood to convert the environment variable string
// to a time.Duration type. If the environment variable is not found, or if the conversion
// fails (e.g., due to an invalid duration format), GetEnvDuration will panic.
//
// The expected format for the duration is a string accepted by time.ParseDuration,
// which includes any input valid for time.Duration such as "300ms", "1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
//
// Example:
// If you have an environment variable named TIMEOUT with the value "2m30s",
// GetEnvDuration("TIMEOUT") will return a time.Duration of 2 minutes and 30 seconds.
//
//	os.Setenv("TIMEOUT", "2m30s")
//	timeout := GetEnvDuration("TIMEOUT")
//	fmt.Println(timeout) // Prints: 2m30s
//
// Note: Because it can panic, this function should be used in cases where the environment variable
// is expected to be set and correctly formatted. For more flexible error handling, consider
// using the underlying GetEnvAs function directly with appropriate error checks.
func GetEnvDuration(key string) time.Duration {
	return GetEnvAs[time.Duration](key, MapStringToDuration)
}

// GetEnvTime is the same helper as GetEnvDuration, but for time.Time.
func GetEnvTime(key string, layout string) time.Time {
	return GetEnvAs[time.Time](key, MapStringToTime(layout))
}

// GetEnvURL helper for URL.
func GetEnvURL(key string) url.URL {
	return GetEnvAs[url.URL](key, MapStringToURL)
}

// GetEnvBool helper for bool values.
func GetEnvBool(key string) bool {
	return GetEnvAs[bool](key, MapStringToBool)
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
