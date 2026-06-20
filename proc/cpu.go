package proc

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// CPUStats hold aggregated CPU time counters (in USER_HZ, usually 100 ticks per second
type CPUStats struct {
	User    uint64
	Nice    uint64
	System  uint64
	Idle    uint64
	IOWait  uint64
	IRQ     uint64
	SoftIRQ uint64
	Steal   uint64
	Total   uint64 // sum of all above
}

// ReadCPUStats reads the first line of /proc/stat (aggregate CPU).
func ReadCPUStats() (*CPUStats, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}
	line := scanner.Text() // firt line
	fields := strings.Fields(line)
	if len(fields) < 8 || fields[0] != "cpu" {
		return nil, nil // should not happen
	}

	stats := &CPUStats{}
	// fields[1] is user, fields[2] nice, fields[3] system, fields[4] idle, etc.
	values := make([]uint64, 8)
	for i := 0; i < 8; i++ {
		val, err := strconv.ParseUint(fields[i+1], 10, 64)
		if err != nil {
			return nil, err
		}
		values[i] = val
	}

	stats.User = values[0]
	stats.Nice = values[1]
	stats.System = values[2]
	stats.Idle = values[3]
	stats.IOWait = values[4]
	stats.IRQ = values[5]
	stats.SoftIRQ = values[6]
	stats.Steal = values[7]

	// Total is sum of all
	var total uint64
	for _, v := range values {
		total += v
	}

	stats.Total = total
	return stats, nil

}

// CPUUsage represents a percentage usage snapshot
type CPUUsage struct {
	Percent float64
}

// CalculateCPUUsage takes two CPUStats snapshots and returns the usage percentage.
func CalculateCPUUsage(prev, curr *CPUStats) CPUUsage {
	totalDelta := curr.Total - prev.Total
	idleDelta := curr.Idle - prev.Idle
	usage := 0.0

	if totalDelta > 0 {
		usage = (1.0 - float64(idleDelta)/float64(totalDelta)) * 100.0
	}

	return CPUUsage{Percent: usage}
}
