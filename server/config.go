package server

import (
	"os"
	"time"
)

// Config store server configuration
type Config struct {
	Version  *Version
	Host     string
	Port     int
	Protocol string

	MaxReadBuffer     int
	MaxIdleTimeout    time.Duration
	MaxConcurrentConn int
}

// GetHostName return the hostname
// First lock at ENV, then command line
func GetHostName() string {

	if os.Getenv("INGESTER_HOST_NAME") != "" {
		return os.Getenv("INGESTER_HOST_NAME")
	}
	return "localhost"
}
