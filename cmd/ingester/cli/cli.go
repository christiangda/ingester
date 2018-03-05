package cli

import (
	"os"
	"time"
)

// CommandLine store the CLI configurtion
type CommandLine struct {
	version string
	Host    string
	Port    int

	MaxReadBuffer     int
	MaxIdleTimeout    time.Duration
	MaxConcurrentConn int
}

// New return an instance of CommandLine with
// the client version
func New(version string) *CommandLine {
	return &CommandLine{
		version: version,
	}
}

// GetHostName return the hostname
// First lock at ENV, then command line
func GetHostName() string {

	if os.Getenv("INGESTER_HOST_NAME") != "" {
		return os.Getenv("INGESTER_HOST_NAME")
	}
	return "localhost"
}
