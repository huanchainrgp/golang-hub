package utils

import "fmt"

// Greet returns a greeting message for the given name
func Greet(name string) string {
	return fmt.Sprintf("Hello, %s! Welcome to Golang Hub.", name)
}

// Sum calculates the sum of integers in a slice
func Sum(numbers []int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

// Average calculates the average of integers in a slice
func Average(numbers []int) float64 {
	if len(numbers) == 0 {
		return 0
	}
	return float64(Sum(numbers)) / float64(len(numbers))
}
