package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	ErrMissingCredentials = errors.New("missing required credentials (CLIENT_ID or CLIENT_SECRET)")
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	collector, err := NewCollector(config)
	if err != nil {
		log.Fatalf("Failed to create collector: %v", err)
	}
	
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	
	collector.Start()
	log.Printf("BikeBoxes Collector started with schedule: %s", config.JobSchedule)
	
	<-sigCh
	log.Println("Shutdown signal received")
	
	if err := collector.Stop(); err != nil {
		log.Fatalf("Failed to stop collector: %v", err)
	}
	
	log.Println("BikeBoxes Collector shutdown complete")
}