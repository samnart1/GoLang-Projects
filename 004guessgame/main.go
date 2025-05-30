package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("We are building a guessing game!")
	if err := something(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func something() error {
	fmt.Println("haha got you")
	return nil
}