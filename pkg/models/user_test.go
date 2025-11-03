package models

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	name := "John Doe"
	email := "john@example.com"

	user := NewUser(name, email)

	if user == nil {
		t.Fatal("NewUser returned nil")
	}

	if user.Name != name {
		t.Errorf("Name = %q; want %q", user.Name, name)
	}

	if user.Email != email {
		t.Errorf("Email = %q; want %q", user.Email, email)
	}

	if user.ID == "" {
		t.Error("ID should not be empty")
	}

	if user.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	// Check that CreatedAt is recent (within last second)
	if time.Since(user.CreatedAt) > time.Second {
		t.Error("CreatedAt should be recent")
	}
}

func TestGenerateID(t *testing.T) {
	id := generateID()

	if id == "" {
		t.Error("generateID returned empty string")
	}

	// ID should have the expected prefix
	expectedPrefix := "user-"
	if len(id) < len(expectedPrefix) {
		t.Errorf("ID %q is too short", id)
	}

	if id[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("ID %q should start with %q", id, expectedPrefix)
	}
}
