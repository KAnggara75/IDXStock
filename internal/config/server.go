package config

import "os"

// GetPort retrieves the port from environment variable or defaults to 3000.
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return port
}
