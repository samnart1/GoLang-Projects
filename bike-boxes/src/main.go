package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Define custom errors
var (
	ErrMissingCredentials = errors.New("missing required credentials (CLIENT_ID or CLIENT_SECRET)")
)

func main() {
	// Load configuration from environment
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	// Create and start collector
	collector, err := NewCollector(config)
	if err != nil {
		log.Fatalf("Failed to create collector: %v", err)
	}
	
	// Setup signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	
	// Start the collector
	collector.Start()
	log.Printf("BikeBoxes Collector started with schedule: %s", config.JobSchedule)
	
	// Wait for shutdown signal
	<-sigCh
	log.Println("Shutdown signal received")
	
	// Clean shutdown
	if err := collector.Stop(); err != nil {
		log.Fatalf("Failed to stop collector: %v", err)
	}
	
	log.Println("BikeBoxes Collector shutdown complete")
}