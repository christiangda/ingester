package server

import (
	"bufio"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Server return new server
type Server struct {
	sync.RWMutex

	ip       string
	port     int
	protocol string

	listener net.Listener

	maxReadBuffer     int
	maxIdleTimeout    time.Duration
	maxConcurrentConn int

	running  bool
	shutdown bool

	clients *ClientStore
}

const (
	// DefaultIdleTimeout for connections
	DefaultIdleTimeout = (time.Second * 2)

	// DefaultReadBuffer for connections
	DefaultReadBuffer = 1024

	// DefaultConcurrentConn for server
	DefaultConcurrentConn = 1024
)

// NewServer return
func NewServer(ip string, port int, protocol string) *Server {

	return &Server{
		ip:       ip,
		port:     port,
		protocol: protocol,

		maxReadBuffer:     DefaultReadBuffer,
		maxIdleTimeout:    DefaultIdleTimeout,
		maxConcurrentConn: DefaultConcurrentConn,

		clients: NewClientStore(),
	}
}

//ListenAndServer start
func (s *Server) ListenAndServer() (err error) {
	ipAndPort := s.ip + ":" + strconv.Itoa(s.port)

	listener, err := net.Listen(s.protocol, ipAndPort)
	if err != nil {
		return err
	}
	defer listener.Close()

	s.listener = listener

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}

		client := &Client{
			ID:          uuid.New(),
			Conn:        conn,
			Created:     time.Now(),
			IdleTimeout: s.maxIdleTimeout,
			ReadBuffer:  s.maxReadBuffer,
		}

		s.handleClient(client)
	}
}

// handleClient
func (s *Server) handleClient(c *Client) error {
	defer func() {
		c.Conn.Close()
		s.delClient(c)
	}()

	s.addClient(c)

	r := bufio.NewReader(c.Conn)
	w := bufio.NewWriter(c.Conn)

	scanr := bufio.NewScanner(r)
	for {
		scanned := scanr.Scan()
		if !scanned {
			if err := scanr.Err(); err != nil {
				return err
			}
			break
		}
		w.WriteString(strings.ToUpper(scanr.Text()) + "\n")
		w.Flush()
	}
	return nil
}

func (s *Server) addClient(c *Client) {
	s.clients.Set(c.ID, c)
}

func (s *Server) delClient(c *Client) {
	s.clients.Delete(c.ID)
}

// Clients return the number of clients
func (s *Server) Clients() int {
	return s.clients.Count()
}

// SetReadBuffer set the max read buffer size for connection
func (s *Server) SetReadBuffer(limit int) {
	s.maxReadBuffer = limit
}

// SetIdleTimeout set the max idle timeout for connection
func (s *Server) SetIdleTimeout(limit int) {

	s.maxIdleTimeout = (time.Second * 2)
}

// SetMaxConnections set the max num of concurrent conections
func (s *Server) SetMaxConnections(limit int) {
	s.maxConcurrentConn = limit
}

// IsRunning validate if server is running
func (s *Server) IsRunning() bool {
	s.Lock()
	defer s.Unlock()
	return s.running
}
