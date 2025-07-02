package main

import (
	"log"
	"os"

	"github.com/samnart1/golang/010timer/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}