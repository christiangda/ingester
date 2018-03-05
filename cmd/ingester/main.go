package main

import (
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/christiangda/ingester/server"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ver, err := server.NewVersion("1.0.0")
	if err != nil {

	}

	config := server.Config{
		Version: ver,
	}

	s := server.NewServer(config)

	err := s.ListenAndServer()
	if err != nil {
		log.Fatal("server error")
	}

}

// Main represents the program execution.
type Main struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewMain return a new instance of Main.
func NewMain() *Main {
	return &Main{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}
