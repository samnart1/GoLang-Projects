package random

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func IntRange(min, max int) (int, error) {
	if min > max {
		return 0, fmt.Errorf("min (%d) cannot be greater than max (%d)", min, max)
	}

	if min == max {
		return min, nil
	}

	diff := int64(max - min + 1)
	n, err := rand.Int(rand.Reader, big.NewInt(diff))
	if err != nil {
		return 0, fmt.Errorf("failed to generate a random number: %w", err)
	}

	return int(n.Int64()) + min, nil
}