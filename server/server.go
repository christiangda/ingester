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
	cfg Config

	sync.RWMutex

	listener net.Listener

	running  bool
	shutdown bool

	startTime time.Time

	connStore *ConnectionStore
}

var (

	// DefaultIdleTimeout for connections
	DefaultIdleTimeout = (time.Second * 2)

	// DefaultReadBuffer for connections
	DefaultReadBuffer = 1024

	// DefaultConcurrentConn for server
	DefaultConcurrentConn = 1024

	startTime time.Time
)

func init() {
	startTime = time.Now().UTC()
}

// NewServer return
func NewServer(config Config) *Server {

	return &Server{
		cfg: config,

		connStore: NewConnectionStore(),
	}
}

//ListenAndServer start
func (s *Server) ListenAndServer() (err error) {
	ipAndPort := s.cfg.Host + ":" + strconv.Itoa(s.cfg.Port)

	listener, err := net.Listen(s.cfg.Protocol, ipAndPort)
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

		Connection := &Connection{
			ID:          uuid.New(),
			Conn:        conn,
			Created:     time.Now(),
			IdleTimeout: s.cfg.MaxIdleTimeout,
			ReadBuffer:  s.cfg.MaxReadBuffer,
		}

		go s.ConnectionHandler(Connection)
	}
}

// ConnectionHandler
func (s *Server) ConnectionHandler(c *Connection) {

	defer func() {
		c.Conn.Close()
		c.Finished = time.Now() // just statistics
		s.delConnection(c)
		fmt.Printf("cerrando todo")
	}()

	s.addConnection(c) // just statistics

	rw := bufio.NewReadWriter(bufio.NewReader(c.Conn), bufio.NewWriter(c.Conn))

	scanner := bufio.NewScanner(rw)

	// set the max data to read
	readBuffer := make([]byte, s.cfg.MaxReadBuffer)
	scanner.Buffer(readBuffer, s.cfg.MaxReadBuffer)

	for {
		if ok := scanner.Scan(); !ok {
			break
		}

		// -->
		//scanner.Bytes() data readed!

		fmt.Printf("data size: %v, buffer size: %v \n", len(scanner.Bytes()), cap(scanner.Bytes()))

		// <---

		rw.WriteString("Connections connected: " + strconv.Itoa(s.connStore.Count()) + "\n")
		rw.Flush()

		break
	}
}

func (s *Server) addConnection(c *Connection) {
	s.connStore.Set(c.ID, c)
}

func (s *Server) delConnection(c *Connection) {
	s.connStore.Delete(c.ID)
}

// Connections return the number of Connections
func (s *Server) Connections() int {
	return s.connStore.Count()
}

// SetReadBuffer set the max read buffer size for connection
func (s *Server) SetReadBuffer(limit int) {
	s.cfg.MaxReadBuffer = limit
}

// SetIdleTimeout set the max idle timeout for connection
func (s *Server) SetIdleTimeout(limit int) {

	s.cfg.MaxIdleTimeout = (time.Second * 2)
}

// SetMaxConnections set the max num of concurrent conections
func (s *Server) SetMaxConnections(limit int) {
	s.cfg.MaxConcurrentConn = limit
}

// IsRunning validate if server is running
func (s *Server) IsRunning() bool {
	s.Lock()
	defer s.Unlock()
	return s.running
}
