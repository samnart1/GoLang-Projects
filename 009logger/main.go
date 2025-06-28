package main

import (
	"log"
	"os"

	"github.com/samnart1/golang/009l9gger/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Printf("Error executing command: %v", err)
		os.Exit(1)
	}
}