package server

import (
	"net"
	"time"

	"github.com/google/uuid"
)

// Connection is all information about Connection connetion
type Connection struct {
	ID   uuid.UUID
	Conn net.Conn

	Created  time.Time
	Finished time.Time

	IdleTimeout time.Duration

	ReadBuffer int
}

func (c *Connection) Write(b []byte) {
}

func (c *Connection) Read(b []byte) {
	c.updateTimeout()
}

func (c *Connection) close() {
	c.Conn.Close()
}

func (c *Connection) updateTimeout() {
	duration := time.Now().Add(c.IdleTimeout)
	c.Conn.SetDeadline(duration)
}
