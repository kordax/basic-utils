/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package os_utils

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	mathutils "github.com/kordax/basic-utils/math-utils"
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

	cpus := int(mathutils.RoundUp(quota / period))
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
