package monitor

import (
	"goSysMon/proc"
)

// SystemSnapshot holds all data for one refresh cycle
type SystemSnapshot struct {
	Memory    proc.MemInfo
	CPUUsage  proc.CPUUsage
	Processes []proc.Process
}

// Monitor holds state (previous CPU stats) needed to compute CPU percentage
type Monitor struct {
	prevCPUStats *proc.CPUStats
	prevProcCPU  map[int]cpuTicks
}

type cpuTicks struct {
	utime uint64
	stime uint64
}

// NewMonitor creates a new Monitor
func NewMonitor() *Monitor {
	return &Monitor{
		prevProcCPU: make(map[int]cpuTicks),
	}
}

// Collect gathers a new SystemSnapshot
func (m *Monitor) Collect() (*SystemSnapshot, error) {
	snap := &SystemSnapshot{}

	// Memory
	mem, err := proc.ReadMemInfo()
	if err != nil {
		return nil, err
	}
	snap.Memory = *mem

	// CPU usage (aggregate)
	curCPU, err := proc.ReadCPUStats()
	if err != nil {
		return nil, err
	}

	if m.prevCPUStats != nil {
		snap.CPUUsage = proc.CalculateCPUUsage(m.prevCPUStats, curCPU)
	}
	m.prevCPUStats = curCPU

	// Processed
	pids, err := proc.ListPIDs()
	if err != nil {
		return nil, err
	}
	processes := make([]proc.Process, 0, len(pids))
	for _, pid := range pids {
		p, err := proc.ReadProcess(pid)
		if err != nil {
			continue
		}
		processes = append(processes, *p)
	}
	snap.Processes = processes

	// Update per-process CPU ticks for future usage if needed (stretch goal). We'll just store them.

	return snap, nil
}
