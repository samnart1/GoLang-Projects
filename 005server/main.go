package main

import (
	"fmt"
	"os"

	"github.com/samnart1/GoLang-Projects/005server/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}