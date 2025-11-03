# Golang Hub

A demonstration repository showcasing Go (Golang) best practices and common patterns.

## Features

- **Clean Architecture**: Organized project structure with `cmd`, `pkg`, and `internal` directories
- **Modular Design**: Separate packages for models, utilities, and configuration
- **Unit Tests**: Comprehensive test coverage with table-driven tests
- **Go Modules**: Modern dependency management

## Project Structure

```
golang-hub/
├── cmd/
│   └── server/          # Application entry points
│       └── main.go
├── pkg/                 # Public packages (can be imported by external projects)
│   ├── models/          # Data models
│   │   ├── user.go
│   │   └── user_test.go
│   └── utils/           # Utility functions
│       ├── helpers.go
│       └── helpers_test.go
├── internal/            # Private packages (cannot be imported externally)
│   └── config/          # Configuration management
│       └── config.go
├── go.mod               # Go module definition
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.16 or higher (tested with Go 1.24.9)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/huanchainrgp/golang-hub.git
cd golang-hub
```

2. Download dependencies:
```bash
go mod download
```

### Running the Application

```bash
go run cmd/server/main.go
```

### Building

```bash
go build -o bin/server cmd/server/main.go
```

### Running Tests

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Run tests with verbose output:
```bash
go test -v ./...
```

## Examples

### Using Utils Package

```go
import "github.com/huanchainrgp/golang-hub/pkg/utils"

// Greet a user
message := utils.Greet("Alice")
// Output: Hello, Alice! Welcome to Golang Hub.

// Calculate sum
numbers := []int{1, 2, 3, 4, 5}
total := utils.Sum(numbers)
// Output: 15

// Calculate average
avg := utils.Average(numbers)
// Output: 3.0
```

### Using Models Package

```go
import "github.com/huanchainrgp/golang-hub/pkg/models"

// Create a new user
user := models.NewUser("Bob", "bob@example.com")
fmt.Printf("User: %s (%s)\n", user.Name, user.Email)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the MIT License.
