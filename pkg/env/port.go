package env

import (
	"os"
)

const (
	// EnvPort is the name of the env var for port.
	EnvPort = "PORT"
	// DefaultPort to be listened on.
	DefaultPort = "8080"
)

// GetPort returns the port to serve on.
func GetPort() string {
	port := os.Getenv(EnvPort)

	if len(port) == 0 {
		port = DefaultPort
	}

	return port
}
