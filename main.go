package main

import (
	"flag"
	"fmt"
	"goSysMon/display"
	"goSysMon/monitor"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Parse command‑line flags
	interval := flag.Int("interval", 2, "Refresh interval in seconds")
	numProcs := flag.Int("procs", 10, "Number of processes to display")
	flag.Parse()

	// Set up channel to catch SIGINT and SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create the monitor
	mon := monitor.NewMonitor()

	// We'll use a ticker for periodic refresh
	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()

	// Initial collection and display
	snap, err := mon.Collect()
	if err != nil {
		log.Fatalf("Initial collection failed: %v", err)
	}
	display.ShowSnapshot(snap, *numProcs)

	// Main loop: wait for tick or signal
	for {
		select {
		case <-ticker.C:
			snap, err := mon.Collect()
			if err != nil {
				log.Printf("Error collecting: %v", err)
				continue
			}
			display.ShowSnapshot(snap, *numProcs)
		case sig := <-sigChan:
			// Graceful shutdown
			fmt.Printf("\nReceived signal: %v. Shutting down.\n", sig)
			return
		}
	}
}