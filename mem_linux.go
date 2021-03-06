// +build linux

package gopsutil

import (
	"strconv"
	"strings"
	"syscall"
)

func VirtualMemory() (*VirtualMemoryStat, error) {
	filename := "/proc/meminfo"
	lines, _ := readLines(filename)

	ret := &VirtualMemoryStat{}
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		value = strings.Replace(value, " kB", "", -1)

		t, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return ret, err
		}
		switch key {
		case "MemTotal":
			ret.Total = t * 1000
		case "MemFree":
			ret.Free = t * 1000
		case "Buffers":
			ret.Buffers = t * 1000
		case "Cached":
			ret.Cached = t * 1000
		case "Active":
			ret.Active = t * 1000
		case "Inactive":
			ret.Inactive = t * 1000
		}
	}
	ret.Available = ret.Free + ret.Buffers + ret.Cached
	ret.Used = ret.Total - ret.Free
	ret.UsedPercent = float64(ret.Total-ret.Available) / float64(ret.Total) * 100.0

	return ret, nil
}

func SwapMemory() (*SwapMemoryStat, error) {
	sysinfo := &syscall.Sysinfo_t{}

	if err := syscall.Sysinfo(sysinfo); err != nil {
		return nil, err
	}
	ret := &SwapMemoryStat{
		Total: uint64(sysinfo.Totalswap),
		Free:  uint64(sysinfo.Freeswap),
	}
	ret.Used = ret.Total - ret.Free
	//check Infinity
	if ret.Total != 0 {
		ret.UsedPercent = float64(ret.Total-ret.Free) / float64(ret.Total) * 100.0
	} else {
		ret.UsedPercent = 0
	}

	return ret, nil
}
