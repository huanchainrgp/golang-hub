package models

import "time"

// User represents a user in the system
type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

// NewUser creates a new User instance
func NewUser(name, email string) *User {
	return &User{
		ID:        generateID(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}
}

// generateID generates a simple ID (in real apps, use UUID)
func generateID() string {
	return "user-" + time.Now().Format("20060102150405")
}
