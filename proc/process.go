package proc

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Process holds per-process information.
type Process struct {
	PID     int
	Command string // from /proc/[pid]/cmdline
	Name    string // from /proc/[pid]/stat (comm field)
	UTime   uint64 // user CPU ticks
	STime   uint64 // kernel DPU tick
	RSS     uint64 // resident set size in pages (need page size to convert to bytes)
}

// ListPIDs returns all numeric directory names in /proc.
func ListPIDs() ([]int, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	var pids []int
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(e.Name())
		if err == nil {
			pids = append(pids, pid)
		}
	}
	return pids, nil
}

// ReadProcess reads info for a single PID. Returns nil if the process disappears
func ReadProcess(pid int) (*Process, error) {

	// Read /proc/[pid]/stat
	statPath := filepath.Join("/proc", strconv.Itoa(pid), "stat")
	statData, err := os.ReadFile(statPath)
	if err != nil {
		return nil, err // may be os.ErrNotExist if process died
	}

	// The stat file format is tricky: the 'comm' field is in parentheses and can conta in space
	// We find the closing parenthesis.
	line := string(statData)
	closeParen := strings.LastIndex(line, ")")
	if closeParen == -1 {
		return nil, nil
	}
	comm := line[strings.Index(line, "(")+1 : closeParen]

	// The part after ')' starts with a space, then fields seperated by space.
	rest := line[closeParen+2:] // skip ")"
	fields := strings.Fields(rest)
	if len(fields) < 22 {
		return nil, nil
	}

	utime, err := strconv.ParseUint(fields[13], 10, 64)
	if err != nil {
		return nil, err
	}
	stime, err := strconv.ParseUint(fields[14], 10, 64)
	if err != nil {
		return nil, err
	}
	rss, err := strconv.ParseUint(fields[23], 10, 64)
	if err != nil {
		return nil, err
	}

	// Read cmdline
	cmdlinePath := filepath.Join("/proc", strconv.Itoa(pid), "cmdline")
	cmdData, err := os.ReadFile(cmdlinePath)
	command := ""
	if err == nil {
		// cmdline uses null bytes as seperators. Replace with spaces.
		command = string(bytes.ReplaceAll(cmdData, []byte{0}, []byte(" ")))
		command = strings.TrimSpace(command)
	}

	if command == "" {
		command = comm // fallback to process name
	}

	return &Process{
		PID:     pid,
		Command: command,
		Name:    comm,
		UTime:   utime,
		STime:   stime,
		RSS:     rss,
	}, nil
}

// GetPageSize returns the system page size in bytes (needed to convert RSS pages to bytes)
func GetPageSize() uint64 {
	return uint64(os.Getpagesize())
}
