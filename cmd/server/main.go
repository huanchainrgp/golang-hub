package main

import (
	"fmt"
	"log"

	"github.com/huanchainrgp/golang-hub/internal/config"
	"github.com/huanchainrgp/golang-hub/pkg/models"
	"github.com/huanchainrgp/golang-hub/pkg/utils"
)

func main() {
	fmt.Println("Welcome to Golang Hub!")

	// Load configuration
	cfg := config.Load()
	fmt.Printf("Application: %s v%s\n", cfg.AppName, cfg.Version)

	// Demonstrate utils
	greeting := utils.Greet("Developer")
	fmt.Println(greeting)

	// Demonstrate models
	user := models.NewUser("Alice", "alice@example.com")
	fmt.Printf("User created: %s (%s)\n", user.Name, user.Email)

	// Example slice operations
	numbers := []int{1, 2, 3, 4, 5}
	sum := utils.Sum(numbers)
	fmt.Printf("Sum of %v = %d\n", numbers, sum)

	log.Println("Golang Hub started successfully!")
}
