package proc

import (
	"bufio"
	"fmt"
	"io"
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
	return ParseCPUSats(f)
}

func ParseCPUSats(r io.Reader) (*CPUStats, error) {
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		return nil, scanner.Err()
	}

	line := scanner.Text() // firt line
	fields := strings.Fields(line)
	if len(fields) < 8 || fields[0] != "cpu" {
		return nil, fmt.Errorf("unexpected /proc/stat format: %q", line)
	}

	values := make([]uint64, 8)
	for i := 0; i < 8; i++ {
		val, err := strconv.ParseUint(fields[i+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing field %d: %w", i+1, err)
		}
		values[i] = val
	}

	stats := &CPUStats{
		User:    values[0],
		Nice:    values[1],
		System:  values[2],
		Idle:    values[3],
		IOWait:  values[4],
		IRQ:     values[5],
		SoftIRQ: values[6],
		Steal:   values[7],
	}

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
	if prev == nil || curr == nil {
		return CPUUsage{Percent: 0.0}
	}
	totalDelta := curr.Total - prev.Total
	idleDelta := curr.Idle - prev.Idle
	usage := 0.0

	if totalDelta > 0 {
		usage = (1.0 - float64(idleDelta)/float64(totalDelta)) * 100.0
	}

	return CPUUsage{Percent: usage}
}
