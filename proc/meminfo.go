package proc

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// MemInfo holds memory statistics in kilobytes.
type MemInfo struct {
	Total         uint64
	Available     uint64
	Used          uint64 // calculated as Total - Available
}

// ReadMemInfo parses /proc/meminfo and returns a MemInfo struct.
func ReadMemInfo() (*MemInfo, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info := &MemInfo{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}

		switch key {
		case "MemTotal":
			info.Total = value
		case "MemAvailable":
			info.Available = value
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	info.Used = info.Total - info.Available
	return info, nil
}