package display

import (
	"fmt"
	"goSysMon/monitor"
	"goSysMon/proc"
	"strings"
)

// ClearScreen moves cursor to top-left and clear the terminal
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// ShowSnapshot prints the sytem snapshot to shout
func ShowSnapshot(snap *monitor.SystemSnapshot, numProcs int) {
	ClearScreen()

	pageSize := proc.GetPageSize()

	fmt.Println("=== goSysMon ===")
	fmt.Printf("CPU Usage:  %.2f%%\n", snap.CPUUsage.Percent)
	fmt.Printf("Memory:     Total: %d kB | Available: %d kB | Used: %d kB\n", snap.Memory.Total, snap.Memory.Available, snap.Memory.Used)

	// Print header for process table
	fmt.Printf("\n%-7s %-20s %-8s %-10s\n", "PID", "COMMAND", "CPU%", "MEM (kB)")
	fmt.Printf(strings.Repeat("-", 60))

	//Display top N processes (by memory usage, simple sort)
	procs := snap.Processes
	limit := numProcs

	if len(procs) < limit {
		limit = len(procs)
	}

	for i := 0; i < limit; i++ {
		p := procs[i]

		// Convert RSS pages to KB
		memKB := p.RSS * pageSize / 1024

		// CPU% would need previous per-process data, we'll show 0 now
		fmt.Printf("%-7d %-20 %-8.1f %-10d\n",
			p.PID, truncateString(p.Command, 20), 0.0, memKB)
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}
