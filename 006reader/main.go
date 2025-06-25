package main

import (
	"fmt"
	"os"

	"github.com/samnart1/GoLang/006reader/cmd"
	"github.com/samnart1/GoLang/006reader/internal/logger"
)

var (
	version	= "dev"
	commit	= "none"
	date	= "unknown"
)

func main() {
	log, err := logger.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	cmd.SetVersion(version, commit, date)

	if err := cmd.Execute(); err != nil {
		log.Error("Application failed", logger.Error(err))
		os.Exit(1)
	}
}