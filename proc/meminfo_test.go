package proc

import (
	"strings"
	"testing"
)

func TestReadMemInfo(t *testing.T) {
	// Create a temporary meminfo file
	input := `MemTotal:       16384000 kB
				MemFree:        1000000 kB
				MemAvailable:   5000000 kB
				Buffers:        200000 kB
				Cached:         3000000 kB
				`

	info, err := ParseMemInfo(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.Total != 16384000 {
		t.Errorf("Total = %d, want 16384000", info.Total)
	}
	if info.Available != 5000000 {
		t.Errorf("Available = %d, want 5000000", info.Available)
	}
	expectedUsed := info.Total - info.Available
	if info.Used != expectedUsed {
		t.Errorf("Used = %d, want %d", info.Used, expectedUsed)
	}
}
