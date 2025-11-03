package config

// Config holds application configuration
type Config struct {
	AppName string
	Version string
	Port    int
}

// Load returns the application configuration
func Load() *Config {
	return &Config{
		AppName: "Golang Hub",
		Version: "1.0.0",
		Port:    8080,
	}
}
