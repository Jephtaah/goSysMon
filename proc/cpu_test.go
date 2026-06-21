package proc

import (
	"strings"
	"testing"
)

func TestParseCPUStats_Valid(t *testing.T) {
	input := "cpu 100 200 300 400 500 600 700 800\n"
	stats, err := ParseCPUSats(strings.NewReader(input))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.User != 100 {
		t.Errorf("User = %d, want 100", stats.User)
	}

	expectedTotal := uint64(100 + 200 + 300 + 400 + 500 + 600 + 700 + 800)
	if stats.Total != expectedTotal {
		t.Errorf("Total = %d, want %d", stats.Total, expectedTotal)
	}
}

func TestParseCPUStats_BadFormat(t *testing.T) {
	_, err := ParseCPUSats(strings.NewReader("notcpu 1 2 3"))
	if err == nil {
		t.Error("expected error for non-cpu line")
	}
}
