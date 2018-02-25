package main

import (
	"github.com/christiangda/ingester/server"
)

func main() {

	// portString := os.Getenv("SERVER_PORT")
	// ipString := os.Getenv("SERVER_IP")

	s := server.NewServer("127.0.0.1", 8080, "tcp")
	s.SetReadBuffer(100)
	s.SetIdleTimeout(100)
	s.SetMaxConnections(100)

	err := s.ListenAndServer()
	if err != nil {

	}

}
