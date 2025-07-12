package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang/011api/internal/config"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Server will run on: %s\n", cfg.GetServerAddress())
	fmt.Printf("Database URL: %s\n", cfg.GetDatabaseURL())
	fmt.Printf("Redis URL: %s\n", cfg.GetRedisAddress())
}
