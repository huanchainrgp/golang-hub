package utils

import "testing"

func TestGreet(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"basic", "Alice", "Hello, Alice! Welcome to Golang Hub."},
		{"empty", "", "Hello, ! Welcome to Golang Hub."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Greet(tt.input)
			if result != tt.expected {
				t.Errorf("Greet(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{"positive numbers", []int{1, 2, 3, 4, 5}, 15},
		{"with zero", []int{0, 1, 2}, 3},
		{"empty slice", []int{}, 0},
		{"single element", []int{42}, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sum(tt.input)
			if result != tt.expected {
				t.Errorf("Sum(%v) = %d; want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAverage(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected float64
	}{
		{"positive numbers", []int{2, 4, 6}, 4.0},
		{"empty slice", []int{}, 0.0},
		{"single element", []int{10}, 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Average(tt.input)
			if result != tt.expected {
				t.Errorf("Average(%v) = %f; want %f", tt.input, result, tt.expected)
			}
		})
	}
}
