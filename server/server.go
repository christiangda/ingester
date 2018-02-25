package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
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

	startTime    time.Time
	shutdownTime time.Time

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
	s.startTime = time.Now()

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

		go s.clientHandler(client)
	}
}

// clientHandler
func (s *Server) clientHandler(c *Client) {

	defer func() {
		c.Conn.Close()
		c.Finished = time.Now() // just statistics
		s.delClient(c)
		fmt.Printf("cerrando todo")
	}()

	s.addClient(c) // just statistics

	rw := bufio.NewReadWriter(bufio.NewReader(c.Conn), bufio.NewWriter(c.Conn))

	scanner := bufio.NewScanner(rw)

	// set the max data to read
	readBuffer := make([]byte, s.maxReadBuffer)
	scanner.Buffer(readBuffer, s.maxReadBuffer)

	for {
		if ok := scanner.Scan(); !ok {
			break
		}

		// -->
		//scanner.Bytes() data readed!

		fmt.Printf("data size: %v, buffer size: %v \n", len(scanner.Bytes()), cap(scanner.Bytes()))

		// <---

		rw.WriteString("clients connected: " + strconv.Itoa(s.clients.Count()) + "\n")
		rw.Flush()

		break
	}
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
