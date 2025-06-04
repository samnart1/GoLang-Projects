package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetGuess(min, max int) (int, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("Enter your guess (%d-%d): ", min, max)
		input, _, err := reader.ReadLine()
		if err != nil {
			return 0, err
		}

		guess, err := strconv.Atoi(strings.TrimSpace(string(input)))
		if err != nil {
			fmt.Println("Please enter a valid guess number")
			continue
		}

		if guess < min || guess > max {
			fmt.Printf("Please enter a number between %d and %d.\n", min, max)
			continue
		}

		return guess, nil

	}
}