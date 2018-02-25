package server

import (
	"net"
	"time"

	"github.com/google/uuid"
)

// Client is all information about Client connetion
type Client struct {
	ID      uuid.UUID
	Conn    net.Conn
	Created time.Time

	IdleTimeout time.Duration

	ReadBuffer int
}

func (c *Client) Write(b []byte) {
}

func (c *Client) Read(b []byte) {
	c.updateTimeout()
}

func (c *Client) close() {
	c.Conn.Close()
}

func (c *Client) updateTimeout() {
	duration := time.Now().Add(c.IdleTimeout)
	c.Conn.SetDeadline(duration)
}
