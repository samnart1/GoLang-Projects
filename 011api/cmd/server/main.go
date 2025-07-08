package main

import (
	"fmt"
	"os"

	"github.com/golang/011api/internal/config"
)

func main() {
	cfg, err := config.LoadConfig(""); if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Config loaded successfully for environment: %s\n", cfg.Environment)
}