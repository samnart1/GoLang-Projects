package mathutils

func Clamp(value, min, max int) int {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}

func InRange(value, min, max int) bool {
	return value >= min && value <= max
}

func PercentageDifference(a, b int) float64 {
	if a == 0 && b == 0 {
		return 0
	}

	diff := float64(abs(a - b))
	avg := float64(a+b) / 2.0
	if avg == 0 {
		return 100.0
	}
	return (diff / avg) * 100.0
}

func abs(x int) int {
	if x < 0 {
		return - x
	}
	return x
}